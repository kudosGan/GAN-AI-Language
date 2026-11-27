package main

// Talking to an AI model: ask("question") sends text to Ollama, a program that runs AI
// models on this computer, and gives back the model's answer as text.
//
// Safety rules:
//   - The address is fixed to this computer. There is no setting to send questions elsewhere,
//     no proxy is used, and redirects are not followed. (If you set up a "cloud" model inside
//     Ollama, Ollama itself would send questions online.)
//   - The answer is only TEXT. GAN never runs anything a model writes.
//   - GAN stops waiting after a time limit.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// ollamaURL and askTimeout are variables only so that the tests can use a fake server.
var (
	ollamaURL  = "http://127.0.0.1:11434/api/generate"
	askTimeout = 5 * time.Minute
)

// defaultModel is used when ask is not told which model to use.
const defaultModel = "qwen2.5:7b"

// maxAnswerBytes stops a runaway answer from using up all the memory.
const maxAnswerBytes = 10_000_000

// askModel sends one question to a model and returns its answer. Each call stands alone:
// the model does not remember earlier questions.
func askModel(question, model string) (string, error) {
	request, err := json.Marshal(map[string]any{
		"model":  model,
		"prompt": question,
		"stream": false, // send the whole answer at once, not word by word
	})
	if err != nil {
		return "", fmt.Errorf("ask could not prepare the question")
	}

	client := &http.Client{
		Timeout:   askTimeout,
		Transport: &http.Transport{Proxy: nil}, // go straight to this computer
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse // never pass the question on to another address
		},
	}
	response, err := client.Post(ollamaURL, "application/json", bytes.NewReader(request))
	if err != nil {
		return "", askFailed(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maxAnswerBytes+1))
	if err != nil {
		return "", askFailed(err)
	}
	if len(body) > maxAnswerBytes {
		return "", fmt.Errorf("ask stopped: the answer from Ollama was too long")
	}

	// Ollama answers with {"response": "..."} or, when something is wrong, {"error": "..."}
	var reply struct {
		Response *string `json:"response"`
		Error    string  `json:"error"`
	}
	readErr := json.Unmarshal(body, &reply)

	if response.StatusCode == http.StatusNotFound && strings.Contains(reply.Error, "not found") {
		return "", fmt.Errorf("Ollama does not have the model %q; see your models with: ollama list", model)
	}
	if reply.Error != "" {
		return "", fmt.Errorf("Ollama reported an error: %s", reply.Error)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama reported an error (status %d)", response.StatusCode)
	}
	if readErr != nil || reply.Response == nil {
		return "", fmt.Errorf("ask got a reply from Ollama that it could not read")
	}
	return strings.TrimSpace(*reply.Response), nil
}

// askFailed turns a network problem into plain words.
func askFailed(err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Errorf("the model took longer than %s to answer, so ask stopped waiting", waitText(askTimeout))
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" { // could not even connect
		return fmt.Errorf("ask could not reach Ollama on this computer; check that Ollama is running")
	}
	return fmt.Errorf("ask lost the connection to Ollama")
}

// waitText writes a time limit the way people say it.
func waitText(d time.Duration) string {
	if d >= time.Minute {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	return d.String()
}

// askArgument checks that an argument of ask is text that is not empty.
// what is "the question" or "the model name".
func askArgument(v Value, what string) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("ask needs text as %s, but got %s", what, typeName(v))
	}
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("ask needs text as %s, but it is empty", what)
	}
	return s, nil
}
