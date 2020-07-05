package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeRequest(t *testing.T) {

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health-check", nil)

	r := MakeRequest(recorder, req)

	if r.writer != recorder {
		t.Errorf("Mismatched writer in req")
	}

	if r.reader != req {
		t.Errorf("Mismatched reader in req, should match supplied go request")
	}

	if r.isHandled {
		t.Errorf("Req isHandled filed must be false on new requests")
	}

}

func TestReply(t *testing.T) {

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health-check", nil)

	r := MakeRequest(recorder, req)

	const body = "Test"
	var headers = []string{"X-Test", "X-Kek"}

	r.Reply(http.StatusOK, map[string]string{headers[0]: body, headers[1]: body}, body)

	// check status
	if r.status != http.StatusOK {
		t.Errorf("Unexpected request status. Got %v wanted %v", r.status, http.StatusOK)
	}

	if recorder.Code != r.status {
		t.Errorf("Unexpected writer status. Got %v wanted %v", recorder.Code, http.StatusOK)
	}

	// check headers
	for i := 0; i < len(headers); i++ {
		if recorder.Header().Get(headers[i]) == "" {
			t.Errorf("Missing header: %s", headers[i])
		}

		if recorder.Header().Get(headers[i]) != body {
			t.Errorf("Header mismatch. Got %s wanted %s", recorder.Header().Get(headers[i]), body)
		}

	}

	// check body
	b, _ := ioutil.ReadAll(recorder.Body)
	rb := string(b)
	if rb != body {
		t.Errorf("Mismatched body. Expected %s, got %s", body, rb)
	}

}
