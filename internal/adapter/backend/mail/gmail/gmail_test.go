package gmail

import (
	"amadeus-trip-parser/internal/domain/model"
	"testing"
)

const credentials = "../../../../client_credentials.json"
const token = "../../../../gmail_token.json"

func TestNewGMailClient(t *testing.T) {
	if testing.Short() {
		return
	}

	type args struct {
		credFile  string
		tokenFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty params",
			args{},
			true,
		},
		{
			"correct params",
			args{credentials, token},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGMailClient(tt.args.credFile, tt.args.tokenFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGMailClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_client_GetEmails(t *testing.T) {
	if testing.Short() {
		return
	}

	g, err := NewGMailClient(credentials, token)
	if err != nil {
		t.Errorf("cannot connect: %v", err)
	}
	type args struct {
		filter string
	}
	tests := []struct {
		name string
		args args
		want func(*testing.T, []*model.Email)
	}{
		{name: "get at least 1 email",
			args: args{""},
			want: func(t *testing.T, ems []*model.Email) {
				if len(ems) <= 0 {
					t.Errorf("GetEmails() got %d emails, want at least 1", len(ems))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.GetEmails(tt.args.filter)
			tt.want(t, got)
		})
	}
}
