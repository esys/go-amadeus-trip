package gmail

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
)

type client struct {
	service *gmail.Service
}

func credentialsFromFile(file string) (*oauth2.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

func NewGMailClient(credFile string, tokenFile string) (domain.EmailProvider, error) {
	g := &client{}
	cred, err := credentialsFromFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read credentials: %w", err)
	}
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read token: %w", err)
	}
	http := cred.Client(context.Background(), tok)
	svc, err := gmail.NewService(context.Background(), option.WithHTTPClient(http))
	if err != nil {
		return nil, fmt.Errorf("unable to create gmail service: %w", err)
	}
	g.service = svc
	return g, nil
}

func (g *client) GetEmails(filter string) []*model.Email {
	var ms []*model.Email
	pageToken := ""
	for {
		req := g.service.Users.Messages.List("me").Q(filter)
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		r, err := req.Do()
		if err != nil {
			log.Error().Msgf("unable to retrieve messages: %v", err)
		}

		log.Debug().Msgf("getting %v messages", len(r.Messages))
		for _, m := range r.Messages {
			//first get only meta to have parsed headers
			msg, err := g.service.Users.Messages.Get("me", m.Id).Format("metadata").Do()
			if err != nil {
				log.Error().Msgf("Unable to retrieve message %v: %v", m.Id, err)
			}
			date := ""
			subject := ""
			for _, h := range msg.Payload.Headers {
				switch h.Name {
				case "Date":
					date = h.Value
				case "Subject":
					subject = h.Value
				}
			}
			//then to have all email as raw content for ulterior parsing purpose
			rawMail, err := g.service.Users.Messages.Get("me", m.Id).Format("raw").Do()
			if err != nil {
				log.Error().Msgf("Unable to retrieve message raw content %v: %v", m.Id, err)
			}
			ms = append(ms, &model.Email{
				Subject: subject,
				Size:    msg.SizeEstimate,
				ID:      msg.Id,
				Date:    date,
				Snippet: msg.Snippet,
				Content: rawMail.Raw,
			})
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
	return ms
}


