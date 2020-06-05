package amadeus

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
Amadeus constants
*/

const (
	ResourceAuthorize string = "v1/security/oauth2/token"
	ResourceJobs      string = "v2/travel/trip-parser-jobs"
)

type Status string

const (
	StatusStarted    Status = "STARTED"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusError      Status = "ERROR"
)

const (
	ContentType string = "application/vnd.amadeus+json"
	APIType     string = "trip-parser-job"
)

type amadeusConfig struct {
	url    string
	key    string
	secret string
}

type tripAPI struct {
	cfg       amadeusConfig
	client    *http.Client
	token     string
	expiresAt time.Time
}

func NewAmadeusTripAPI(url string, key string, secret string) (domain.EmailParser, error) {
	t := tripAPI{
		client: &http.Client{
			Timeout: time.Second * 15,
		},
		cfg: amadeusConfig{url, key, secret},
	}
	if err := t.authorize(); err != nil {
		return nil, fmt.Errorf("cannot create amadeus api: %w", err)
	}
	return &t, nil
}

func (t *tripAPI) convertAmadeusStatus(s Status) model.MailParsingStatus {
	switch s {
	case StatusError:
		return model.MailParsingStatusError
	case StatusCompleted:
		return model.MailParsingStatusDone
	case StatusInProgress:
		return model.MailParsingStatusPending
	case StatusStarted:
		return model.MailParsingStatusPending
	default:
		log.Debug().Msgf("unknown Amadeus status %s was converted to %s", s, model.MailParsingStatusError)
		return model.MailParsingStatusError
	}
}

func (t *tripAPI) buildRequest(method string, resource string, body io.Reader) (*http.Request, error) {
	resURL := fmt.Sprintf("%s/%s", t.cfg.url, resource)
	req, err := http.NewRequest(method, resURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", ContentType)
	req.Header.Add("Authorization", "Bearer "+t.token)
	return req, nil
}

func (t *tripAPI) authorize() error {
	// token is not expired, do nothing
	if t.token != "" && !t.expiresAt.IsZero() && time.Now().Before(t.expiresAt) {
		return nil
	}
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", t.cfg.key, t.cfg.secret)
	url := fmt.Sprintf("%s/%s", t.cfg.url, ResourceAuthorize)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("cannot create authorize request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var resp authorizeResponse
	code, err := t.doRequest(req, &resp)
	if err != nil {
		return fmt.Errorf("failed authorize request: %w", err)
	}

	if code/100 != 2 {
		return fmt.Errorf("failed to authorize, got status code %d", code)
	}

	t.token = resp.AccessToken
	t.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	return nil
}

func (t *tripAPI) doRequest(req *http.Request, payload interface{}) (int, error) {
	httpResp, err := t.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed request %v: %w", req, err)
	}
	defer httpResp.Body.Close()

	byt, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return httpResp.StatusCode, fmt.Errorf("cannot read response body %v with code %d: %w",
			httpResp,
			httpResp.StatusCode,
			err)
	}

	if err := json.Unmarshal(byt, &payload); err != nil {
		return httpResp.StatusCode, fmt.Errorf("cannot decode response json body %s with code %d: %w",
			string(byt),
			httpResp.StatusCode,
			err)
	}

	return httpResp.StatusCode, nil
}

func (t *tripAPI) getWarnings(errors []apiError) []string {
	if errors == nil {
		return []string{}
	}
	var all []string
	for _, e := range errors {
		all = append(all, e.Detail)
	}
	return all
}

func (t *tripAPI) CreateJob(mail *model.Email) (*model.EmailParsingJob, error) {
	content := strings.ReplaceAll(mail.Content, "-", "+")
	content = strings.ReplaceAll(content, "_", "/")
	payload := createRequest{
		Data: createRequestData{
			TypeP:   APIType,
			Content: content,
		},
	}
	byt, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode create request for %v: %w", mail, err)
	}

	req, err := t.buildRequest(http.MethodPost, ResourceJobs, bytes.NewReader(byt))
	if err != nil {
		return nil, fmt.Errorf("cannot create request for %v: %w", mail, err)
	}

	var body createResponse
	code, err := t.doRequest(req, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to process create request for %v: %w", mail, err)
	}

	if code/100 != 2 {
		return nil, fmt.Errorf("unexpected create request response code: got %d with errors %v", code, body.Errors)
	}

	return &model.EmailParsingJob{
		ID:      body.Data.ID,
		Subject: mail.Subject,
		Status:  t.convertAmadeusStatus(body.Data.Status),
	}, nil
}

func (t *tripAPI) GetJobStatus(job model.EmailParsingJob) (*model.EmailParsingJob, error) {
	res := fmt.Sprintf("%s/%s", ResourceJobs, job.ID)
	req, err := t.buildRequest(http.MethodGet, res, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request for %s: %w", job, err)
	}

	var body statusResponse
	code, err := t.doRequest(req, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to process status request for %v: %w", job, err)
	}

	if code/100 != 2 {
		return nil, fmt.Errorf("unexpected create request response code: got %d with errors %v", code, body.Errors)
	}

	return &model.EmailParsingJob{
		ID:       body.Data.ID,
		Subject:  job.Subject,
		Status:   t.convertAmadeusStatus(body.Data.Status),
		Warnings: t.getWarnings(body.Warnings),
		Detail:   body.Data.Detail, // filled when status is ERROR
	}, nil
}

func (t *tripAPI) GetJobResult(job model.EmailParsingJob) (*model.EmailParsingJob, error) {
	res := fmt.Sprintf("%s/%s/result", ResourceJobs, job.ID)
	req, err := t.buildRequest(http.MethodGet, res, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request for %s: %w", job, err)
	}

	var body resultResponse
	code, err := t.doRequest(req, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to process status request for %v: %w", job, err)
	}

	if code/100 != 2 {
		return nil, fmt.Errorf("unexpected create request response code: got %d with errors %v", code, body.Errors)
	}

	trip, err := newTripConverter().getTrip(body.Data)
	if err != nil {
		return nil, fmt.Errorf("cannot convert data to trip: %w", err)
	}
	return &model.EmailParsingJob{
		ID:       body.Data.ID,
		Subject:  job.Subject,
		Warnings: t.getWarnings(body.Warnings),
		Trip:     trip,
	}, nil
}
