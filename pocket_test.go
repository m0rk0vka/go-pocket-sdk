package pocket

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func newClient(t *testing.T, statusCode int, path, body string) *Client {
	return &Client{
		client: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, path, r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(strings.NewReader(body)),
				}, nil
			}),
		},
		consumerKey: "key",
	}
}

func TestClient_GetRequestToken(t *testing.T) {
	type args struct {
		ctx         context.Context
		redirectUrl string
	}

	tests := []struct {
		name                 string
		args                 args
		expectedStatusCode   int
		expectedResponse     string
		expectedErrorMessage string
		want                 string
		wantErr              bool
	}{
		{
			name: "Ok",
			args: args{
				ctx:         context.Background(),
				redirectUrl: "http://localhost",
			},
			expectedStatusCode: 200,
			expectedResponse:   "code=qwe-rty-123",
			want:               "qwe-rty-123",
			wantErr:            false,
		},
		{
			name: "Empty redirect URL",
			args: args{
				ctx:         context.Background(),
				redirectUrl: "",
			},
			wantErr: true,
		},
		{
			name: "Empty response code",
			args: args{
				ctx:         context.Background(),
				redirectUrl: "http://localhost",
			},
			expectedStatusCode: 200,
			expectedResponse:   "code=",
			wantErr:            true,
		},
		{
			name: "Non-2XX Response",
			args: args{
				ctx:         context.Background(),
				redirectUrl: "http://localhost",
			},
			expectedStatusCode: 400,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t, tt.expectedStatusCode, "/v3/oauth/request", tt.expectedResponse)

			got, err := c.getRequestToken(tt.args.ctx, tt.args.redirectUrl)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}

}
