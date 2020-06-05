package amadeus

import (
	"amadeus-trip-parser/internal/domain/model"
	"os"
	"testing"
)

var config = map[string]string{
	"url":    "https://api.amadeus.com",
	"key":    os.Getenv(key),
	"secret": os.Getenv(secret),
}

const secret = "amadeus-secret"
const key = "amadeus-key"

func checkEnv(t *testing.T) bool {
	if os.Getenv(key) == "" || os.Getenv(secret) == "" {
		t.Errorf("missing amadeus environment configuration for running test")
		return false
	}
	return true
}

func compareJob(a *model.EmailParsingJob, b *model.EmailParsingJob) bool {
	if a == nil && b == nil {
		return true
	}
	return a.Subject == b.Subject && a.Status == b.Status
}

func TestNewAmadeusTripAPI(t *testing.T) {
	if testing.Short() {
		return
	}
	if !checkEnv(t) {
		return
	}

	type args struct {
		config map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "connect with no configuration",
			args:    args{map[string]string{}},
			wantErr: true,
		},
		{
			name:    "connect",
			args:    args{config},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			if _, err := NewAmadeusTripAPI(tt.args.config["url"], tt.args.config["key"], tt.args.config["secret"]);
				(err != nil) != tt.wantErr {
				t1.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tripAPI_CreateJob(t *testing.T) {
	if testing.Short() {
		return
	}
	if !checkEnv(t) {
		return
	}
	var content string
	if content = string(readTestData("testdata/msg-encoded", t)); content == "" {
		return
	}

	api, err := NewAmadeusTripAPI(config["url"], config["key"], config["secret"])
	if err != nil {
		t.Errorf("cannot create amadeus API: %s", err)
	}

	type args struct {
		mail *model.Email
	}
	tests := []struct {
		name    string
		args    args
		want    *model.EmailParsingJob
		wantErr bool
	}{
		{
			"email with empty content",
			args{
				&model.Email{
					Subject: "TEST EMAIL",
				},
			},
			nil,
			true,
		},
		{
			"email with bad content",
			args{
				&model.Email{
					Subject: "TEST EMAIL",
					Content: "fake data",
				},
			},
			nil,
			true,
		},
		{
			"email with ok content",
			args{
				&model.Email{
					Subject: "TEST EMAIL",
					Content: content,
				},
			},
			&model.EmailParsingJob{Status: model.MailParsingStatusPending, Subject: "TEST EMAIL"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			got, err := api.CreateJob(tt.args.mail)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				tt.want.ID = got.ID
			}
			if !compareJob(got, tt.want) {
				t1.Errorf("CreateJob() got = %v, want %v", got, tt.want)
			}
		})
	}
}
