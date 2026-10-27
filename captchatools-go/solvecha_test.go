package captchatoolsgo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func solvechaTestServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	old := solvechaBaseURL
	solvechaBaseURL = srv.URL
	t.Cleanup(func() { solvechaBaseURL = old; srv.Close() })
}

// go test -v -run ^TestSolvecha ./...
func TestSolvechaImageReadyInCreateTask(t *testing.T) {
	solvechaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var p struct {
			ClientKey string `json:"clientKey"`
			Task      struct{ Type, Body string }
		}
		json.NewDecoder(r.Body).Decode(&p)
		if r.URL.Path != "/createTask" || p.ClientKey != "k" || p.Task.Type != "ImageToTextTask" || p.Task.Body != "AAAA" {
			t.Errorf("unexpected request %s %+v", r.URL.Path, p)
		}
		w.Write([]byte(`{"errorId":0,"taskId":"1","status":"ready","solution":{"text":"AB12"}}`))
	})
	h, err := NewHarvester(SolvechaSite, &Config{Api_key: "k", CaptchaType: ImageCaptcha})
	if err != nil {
		t.Fatal(err)
	}
	a, err := h.GetToken(&AdditionalData{B64Img: "AAAA"})
	if err != nil || a.Token != "AB12" {
		t.Fatalf("got %v, %v", a, err)
	}
}

func TestSolvechaPolling(t *testing.T) {
	old := solvechaPollInterval
	solvechaPollInterval = 0
	defer func() { solvechaPollInterval = old }()
	solvechaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/createTask" {
			w.Write([]byte(`{"errorId":0,"taskId":"7","status":"processing"}`))
			return
		}
		w.Write([]byte(`{"errorId":0,"status":"ready","solution":{"text":"XY"}}`))
	})
	s := Solvecha{&Config{Api_key: "k", CaptchaType: ImageCaptcha}}
	a, err := s.GetTokenWithContext(context.Background(), &AdditionalData{B64Img: "AAAA"})
	if err != nil || a.Token != "XY" {
		t.Fatalf("got %v, %v", a, err)
	}
}

func TestSolvechaKeepsAPIErrorMessage(t *testing.T) {
	solvechaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"errorId":1,"errorCode":"ERROR_KEY_DENIED_ACCESS","errorDescription":"The clientKey is missing or wrong."}`))
	})
	s := Solvecha{&Config{Api_key: "bad", CaptchaType: ImageCaptcha}}
	_, err := s.GetToken(&AdditionalData{B64Img: "AAAA"})
	if err == nil || !strings.Contains(err.Error(), "The clientKey is missing or wrong.") {
		t.Fatalf("got %v", err)
	}
}

func TestSolvechaTokenCaptchasNotSupported(t *testing.T) {
	s := Solvecha{&Config{Api_key: "k", CaptchaType: V2Captcha}}
	if _, err := s.GetToken(); !errors.Is(err, ErrNotSupported) {
		t.Fatalf("got %v", err)
	}
}

func TestSolvechaBalance(t *testing.T) {
	solvechaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"errorId":0,"balance":4.5}`))
	})
	s := Solvecha{&Config{Api_key: "k"}}
	if b, err := s.GetBalance(); err != nil || b != 4.5 {
		t.Fatalf("got %v, %v", b, err)
	}
}
