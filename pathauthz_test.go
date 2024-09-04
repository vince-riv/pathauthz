package pathauthz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) http.Handler {
	cfg := CreateConfig()
	cfg.BasePath = "/restricted"
	cfg.UserHeader = "X-User"
	cfg.SuperUsers = []string{"admin-user"}
	cfg.ReadOnlyUsers = []string{"read-only-user"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "pathauthz-test")
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func newRequest(t *testing.T, user string, method string, path string) *http.Request {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, "http://localhost"+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if user != "" {
		req.Header.Add("X-User", user)
	}
	return req
}

func TestUnRestrictedPath(t *testing.T) {
	testCases := []struct {
		path   string
		method string
		user   string
	}{
		{"/", http.MethodGet, ""},
		{"/", http.MethodPost, ""},
		{"/", http.MethodGet, "testuser"},
		{"/", http.MethodPost, "testuser"},
		{"/api/testuser", http.MethodGet, "testuser2"},
		{"/api/testuser", http.MethodPut, "testuser2"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			handler := newHandler(t)
			recorder := httptest.NewRecorder()
			req := newRequest(t, tc.user, tc.method, tc.path)
			handler.ServeHTTP(recorder, req)
			assertStatusCode(t, recorder, 200, tc.method, tc.path, tc.user)
		})
	}
}

func TestRestrictedPathes(t *testing.T) {
	testCases := []struct {
		path     string
		method   string
		user     string
		expected int
	}{
		{"/restricted", http.MethodGet, "testuser", 200},
		{"/restricted", http.MethodGet, "admin-user", 200},
		{"/restricted", http.MethodGet, "read-only-user", 200},
		{"/restricted", http.MethodPut, "testuser", 403},
		{"/restricted", http.MethodPut, "admin-user", 200},
		{"/restricted", http.MethodPut, "read-only-user", 403},
		{"/restricted/testuser", http.MethodGet, "testuser", 200},
		{"/restricted/testuser", http.MethodGet, "testuser2", 200},
		{"/restricted/testuser", http.MethodGet, "admin-user", 200},
		{"/restricted/testuser", http.MethodGet, "read-only-user", 200},
		{"/restricted/testuser", http.MethodPut, "testuser", 200},
		{"/restricted/testuser", http.MethodPut, "testuser2", 403},
		{"/restricted/testuser", http.MethodPut, "admin-user", 200},
		{"/restricted/testuser", http.MethodPut, "read-only-user", 403},
		{"/restricted/testuser/123", http.MethodGet, "testuser", 200},
		{"/restricted/testuser/123", http.MethodGet, "testuser2", 200},
		{"/restricted/testuser/123", http.MethodGet, "admin-user", 200},
		{"/restricted/testuser/123", http.MethodGet, "read-only-user", 200},
		{"/restricted/testuser/123", http.MethodPut, "testuser", 200},
		{"/restricted/testuser/123", http.MethodPut, "testuser2", 403},
		{"/restricted/testuser/123", http.MethodPut, "admin-user", 200},
		{"/restricted/testuser/123", http.MethodPut, "read-only-user", 403},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			handler := newHandler(t)
			recorder := httptest.NewRecorder()
			req := newRequest(t, tc.user, tc.method, tc.path)
			handler.ServeHTTP(recorder, req)
			assertStatusCode(t, recorder, tc.expected, tc.method, tc.path, tc.user)
		})
	}
}

func assertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, status int, method, path, user string) {
	t.Helper()
	if recorder.Result().StatusCode != status {
		t.Errorf("Expected HTTP %d, got %d: %s %s [user %s]", status, recorder.Result().StatusCode, method, path, user)
	}
}

func TestIsWriteReq(t *testing.T) {
	// write methods
	if !isWriteReq(http.MethodDelete) {
		t.Error("DELETE is a write request")
	}
	if !isWriteReq(http.MethodPatch) {
		t.Error("Patch is a write request")
	}
	if !isWriteReq(http.MethodPost) {
		t.Error("POST is a write request")
	}
	if !isWriteReq(http.MethodPut) {
		t.Error("PUT is a write request")
	}
	// read methods
	if isWriteReq(http.MethodGet) {
		t.Error("GET is a read request")
	}
	if isWriteReq(http.MethodHead) {
		t.Error("HEAD is a read request")
	}
	if isWriteReq(http.MethodOptions) {
		t.Error("OPTIONS is a read request")
	}
	// "read" methods
	if isWriteReq(http.MethodConnect) {
		t.Error("CONNECT is a read request")
	}
	if isWriteReq(http.MethodTrace) {
		t.Error("TRACE is a read request")
	}
}
