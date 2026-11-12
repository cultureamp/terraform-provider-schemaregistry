package subjectconfig

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNormalize(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    *bool
		wantErr bool
	}{
		{"false", 200, `{"compatibilityLevel":"FULL_TRANSITIVE","normalize":false}`, boolPtr(false), false},
		{"absent", 200, `{"compatibilityLevel":"FULL_TRANSITIVE"}`, nil, false},
		{"not found", 404, `{"error_code":40408,"message":"Subject not found"}`, nil, false},
		{"server error", 500, `{"error_code":50001,"message":"boom"}`, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %s, want GET", r.Method)
				}
				if r.URL.Path != "/config/my-topic-value" {
					t.Errorf("path = %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("defaultToGlobal"); got != "false" {
					t.Errorf("defaultToGlobal = %q", got)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			got, err := New(srv.URL+"/", nil, "", "").GetNormalize(context.Background(), "my-topic-value")
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if (got == nil) != (tt.want == nil) || (got != nil && *got != *tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNormalize(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr bool
	}{
		{"ok", 200, `{"normalize":false}`, false},
		{"unprocessable", 422, `{"error_code":42203,"message":"bad"}`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("method = %s, want PUT", r.Method)
				}
				if r.URL.Path != "/config/my-topic-value" {
					t.Errorf("path = %s", r.URL.Path)
				}
				if got := r.Header.Get("Content-Type"); got != "application/vnd.schemaregistry.v1+json" {
					t.Errorf("content type = %q", got)
				}
				b, _ := io.ReadAll(r.Body)
				if string(b) != `{"normalize":false}` {
					t.Errorf("body = %s", b)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			err := New(srv.URL, nil, "", "").SetNormalize(context.Background(), "my-topic-value", false)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != "user" || p != "pass" {
			t.Errorf("basic auth = %q, %q, %v", u, p, ok)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	if _, err := New(srv.URL, nil, "user", "pass").GetNormalize(context.Background(), "s"); err != nil {
		t.Fatal(err)
	}
}

func TestSubjectIsPathEscaped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/config/a%2Fb-value" {
			t.Errorf("escaped path = %s", r.URL.EscapedPath())
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	if _, err := New(srv.URL, nil, "", "").GetNormalize(context.Background(), "a/b-value"); err != nil {
		t.Fatal(err)
	}
}

func boolPtr(b bool) *bool { return &b }
