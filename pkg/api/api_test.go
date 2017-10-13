package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type TestResponseWriter struct {
	Response string
	header   http.Header
}

func (wr *TestResponseWriter) Header() http.Header {
	if nil == wr.header {
		wr.header = http.Header{}
	}
	return wr.header
}

func (wr *TestResponseWriter) Write(str []byte) (int, error) {
	wr.Response += string(str)
	return 0, nil
}

func (wr *TestResponseWriter) WriteHeader(code int) {
	return
}

func generateTestRequest() *http.Request {
	request := http.Request{}
	ctx := context.WithValue(context.Background(), CKey("Verbose"), false)
	return request.WithContext(ctx)
}

func generateTestWriter() TestResponseWriter {
	return TestResponseWriter{
		Response: "",
	}
}

func mockAPIHandlerSuccess(request *http.Request) (*executor.Command, error) {
	c := executor.NewProcess("echo", "{\"out\": 0}")
	return c, nil
}

func mockAPIHandlerErrorEmptyOut(request *http.Request) (*executor.Command, error) {
	c := executor.NewProcess("true")
	return c, nil
}

func mockAPIHandlerErrorCorruptedOut(request *http.Request) (*executor.Command, error) {
	c := executor.NewProcess("echo", "{\"out: 0}")
	return c, nil
}

func mockAPIHandlerErrorCmdFail(request *http.Request) (*executor.Command, error) {
	c := executor.NewProcess("ls", "/non-existing-file")
	return c, nil
}

func TestGenericResponseHandler(t *testing.T) {
	request := generateTestRequest()
	writer := generateTestWriter()

	// Scenarios:
	// No output
	// Corrupted output
	// Cmd fail
	// Success
	// Check content type for all?
	//compare output with expected output
	var scenarios = []struct {
		fn      func(request *http.Request) (*executor.Command, error)
		success bool
		code    int
	}{
		{
			fn:      mockAPIHandlerSuccess,
			success: true,
			code:    0,
		},
		{
			fn:      mockAPIHandlerErrorEmptyOut,
			success: false,
			code:    3,
		},
		{
			fn:      mockAPIHandlerErrorCorruptedOut,
			success: false,
			code:    3,
		},
		{
			fn:      mockAPIHandlerErrorCmdFail,
			success: false,
			code:    2,
		},
	}

	for i, s := range scenarios {
		response := &writer.Response

		m := genericResponseHandler(s.fn)
		m(&writer, request)

		var tmp interface{}
		if err := json.Unmarshal([]byte(*response), &tmp); nil != err {
			t.Errorf("Unexpected data received during scenario execution")
		}
		result := tmp.(map[string]interface{})

		if _, ok := result["data"]; ok && s.success {
		} else if errListInterface, ok := result["errors"]; ok && !s.success {
			errorsList := errListInterface.([]interface{})
			// Only one code is expected although the errors is containing a list
			for _, eInterface := range errorsList {
				e := eInterface.(map[string]interface{})
				ecode := int(e["code"].(float64))
				if ecode != s.code {
					t.Errorf("Expected and received error code does not match. Expect %d got %d", s.code, ecode)
				}
			}
		} else {
			t.Errorf("Unexpected result in scenario %d", i)
		}

		*response = ""
	}
}
