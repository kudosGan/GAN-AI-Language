package main

// Checks ask() against a FAKE Ollama server, so the results do not depend on what a real
// model happens to say. Run with: go test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeOllama starts a pretend Ollama and points ask at it until the test ends.
func fakeOllama(t *testing.T, handler http.HandlerFunc) {
	server := httptest.NewServer(handler)
	old := ollamaURL
	ollamaURL = server.URL + "/api/generate"
	t.Cleanup(func() {
		ollamaURL = old
		server.Close()
	})
}

// ask runs the built-in the same way a GAN program would.
func ask(args ...Value) (Value, error) {
	v, _, err := builtins["ask"].run(args)
	return v, err
}

func wantError(t *testing.T, err error, piece string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), piece) {
		t.Errorf("got error %v, want one containing %q", err, piece)
	}
}

func TestAskSendsTheQuestionAndReturnsTheAnswer(t *testing.T) {
	var sent []map[string]any
	fakeOllama(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/generate" {
			t.Errorf("wrong request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		sent = append(sent, body)
		json.NewEncoder(w).Encode(map[string]any{"response": "  He said \"hi\"\nthen left.  \n"})
	})
	question := "Quotes \" and\nnew lines \\ é"
	got, err := ask(question)
	if err != nil || got != "He said \"hi\"\nthen left." {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := ask("again", "other-model"); err != nil {
		t.Fatal(err)
	}
	if sent[0]["model"] != "qwen2.5:7b" || sent[0]["prompt"] != question || sent[0]["stream"] != false {
		t.Errorf("first request was %v", sent[0])
	}
	if sent[1]["model"] != "other-model" || sent[1]["prompt"] != "again" || len(sent[1]) != 3 {
		t.Errorf("second request was %v (each ask must stand alone)", sent[1])
	}
}

func TestAskErrorsFromOllama(t *testing.T) {
	cases := []struct {
		status      int
		body, wants string
	}{
		{404, `{"error":"model 'nope' not found"}`, `does not have the model "nope"`},
		{500, `{"error":"out of memory"}`, "Ollama reported an error: out of memory"},
		{500, `oops`, "status 500"},
		{404, `404 page not found`, "status 404"},
		{200, `{"error":"something odd"}`, "something odd"},
		{200, `not json`, "could not read"},
		{200, `{}`, "could not read"},
	}
	for _, c := range cases {
		fakeOllama(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(c.status)
			w.Write([]byte(c.body))
		})
		_, err := ask("hi", "nope")
		wantError(t, err, c.wants)
	}
}

func TestAskAllowsAnEmptyAnswer(t *testing.T) {
	fakeOllama(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"response":""}`)) })
	if got, err := ask("hi"); err != nil || got != "" {
		t.Errorf("got %q, %v", got, err)
	}
}

func TestAskDoesNotFollowRedirects(t *testing.T) {
	reached := false
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true }))
	defer other.Close()
	for _, status := range []int{301, 302, 307, 308} {
		fakeOllama(t, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, other.URL, status)
		})
		if _, err := ask("secret"); err == nil {
			t.Errorf("status %d: expected an error", status)
		}
	}
	if reached {
		t.Error("the question was sent on to another address")
	}
}

func TestAskRefusesVeryLongAnswers(t *testing.T) {
	fakeOllama(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"response":"` + strings.Repeat("a", maxAnswerBytes) + `"}`))
	})
	_, err := ask("hi")
	wantError(t, err, "too long")
}

func TestAskStopsWaiting(t *testing.T) {
	release := make(chan struct{})
	fakeOllama(t, func(w http.ResponseWriter, r *http.Request) { <-release })
	old := askTimeout
	askTimeout = 200 * time.Millisecond
	defer func() { askTimeout = old; close(release) }()
	_, err := ask("hi")
	wantError(t, err, "took longer than")
	if waitText(5*time.Minute) != "5 minutes" {
		t.Errorf("waitText gave %q", waitText(5*time.Minute))
	}
}

func TestAskWhenOllamaIsNotRunning(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	address := server.URL
	server.Close() // now nothing is listening there
	old := ollamaURL
	ollamaURL = address
	defer func() { ollamaURL = old }()
	_, err := ask("hi")
	wantError(t, err, "check that Ollama is running")
}

func TestAskChecksItsArguments(t *testing.T) {
	called := false
	fakeOllama(t, func(w http.ResponseWriter, r *http.Request) { called = true })
	_, err := ask(42.0)
	wantError(t, err, "ask needs text as the question, but got number")
	_, err = ask("  ")
	wantError(t, err, "the question, but it is empty")
	_, err = ask("hi", true)
	wantError(t, err, "the model name, but got bool")
	_, err = ask("hi", "")
	wantError(t, err, "the model name, but it is empty")
	if called {
		t.Error("nothing should be sent when an argument is wrong")
	}
}

func TestAskAddressIsThisComputer(t *testing.T) {
	if ollamaURL != "http://127.0.0.1:11434/api/generate" || askTimeout != 5*time.Minute {
		t.Errorf("address %q, time limit %v", ollamaURL, askTimeout)
	}
}
