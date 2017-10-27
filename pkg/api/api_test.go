package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

type params struct {
	data   interface{}
	status int
	err    error
}

func TestResponseHandler(t *testing.T) {
	tt := []struct {
		name     string
		input    params
		expected params
	}{
		{
			name: "internal error",
			input: params{
				data:   nil,
				status: 0,
				err:    errors.New("this is a random error"),
			},
			expected: params{
				data:   "Internal error",
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "bad input error",
			input: params{
				data:   nil,
				status: http.StatusBadRequest,
				err:    newAPIError(nil, errBadInput, "bad input msg"),
			},
			expected: params{
				data:   fmt.Sprintf(`{"errors":[{"code":%d,"message":%q}]}`, errBadInput, "bad input msg"),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "correct data",
			input: params{
				data:   map[string]string{"id": "123", "name": "test"},
				status: http.StatusOK,
				err:    nil,
			},
			expected: params{
				data:   fmt.Sprintf(`{"data":{"id":"123","name":"test"}}`),
				status: http.StatusOK,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			fn := func(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
				return tc.input.data, tc.input.status, tc.input.err
			}
			hn := respHandler(fn)

			req, err := http.NewRequest("GET", "sample_url", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			hn(rec, req)

			res := rec.Result()

			b, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}

			if data := string(bytes.TrimSpace(b)); data != tc.expected.data {
				t.Errorf("expected data %v; got %v", tc.expected.data, data)
			}

			if res.StatusCode != tc.expected.status {
				t.Errorf("expected status %v; got %v", tc.expected.status, res.StatusCode)
			}

		})
	}

}
