package radicle

import "testing"

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
			name: "Some URL return the same URL",
			fields: fields{
				url:         "some_url",
				secretToken: "some_token",
			},
			want: "some_url",
		},
		{
			name: "Empty URL return empty string",
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
