package radicle

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		args    Opts
		wantURL string
		wantErr bool
	}{
		{
			name: "creating radicle forge with some data returns them",
			args: Opts{
				URL:         "http://some_url/without/trail",
				SecretToken: "a_super_secret_token",
			},
			wantURL: "http://some_url/without/trail",
			wantErr: false,
		},
		{
			name: "creating radicle forge with some data returns them without trailing / at URL",
			args: Opts{
				URL:         "http://some_url/with/trail/",
				SecretToken: "a_super_secret_token",
			},
			wantURL: "http://some_url/with/trail",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.URL() != tt.wantURL {
				t.Errorf("New() URL got = %v, want %v", got.URL(), tt.wantURL)
			}
			if got.Name() != "radicle" {
				t.Errorf("New() name got = %v, want %v", got.Name(), "radicle")
			}
		})
	}
}

func Test_radicle_URL(t *testing.T) {
	type fields struct {
		url         string
		secretToken string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "some URL return the same URL",
			fields: fields{
				url:         "some_url",
				secretToken: "some_token",
			},
			want: "some_url",
		},
		{
			name: "empty URL return empty string",
			fields: fields{
				url:         "",
				secretToken: "some_token",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rad := radicle{
				url:         tt.fields.url,
				secretToken: tt.fields.secretToken,
			}
			if got := rad.URL(); got != tt.want {
				t.Errorf("URL() = %v, want %v", got, tt.want)
			}
		})
	}
}
