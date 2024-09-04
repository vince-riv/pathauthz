package pathauthz

import (
	"context"
	"net/http"
	"strings"
)

// Config holds the plugin configuration.
type Config struct {
	BasePath      string   `json:"basePath,omitempty"`     // Configurable base path (e.g., "/v2")
	UserHeader    string   `json:"userHeader,omitempty"`   // Name of the header that contains the authenticated user
	SuperUsers    []string `json:superUsers,omitempty"`    // List of users who can write to all paths
	ReadOnlyUsers []string `json:readonlyUsers,omitempty"` // List of users can only read
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		BasePath:   "/v2",                  // Default value for the base path
		UserHeader: "X-Authenticated-User", // Default header for user identification
	}
}

// RestrictMethodMiddleware is a middleware that restricts certain HTTP methods on specific paths.
type RestrictMethodMiddleware struct {
	next          http.Handler
	basePath      string
	userHeader    string
	superUsers    map[string]struct{}
	readOnlyUsers map[string]struct{}
}

// New creates a new RestrictMethodMiddleware.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	superUsers := make(map[string]struct{})
	readOnlyUsers := make(map[string]struct{})

	for _, user := range config.SuperUsers {
		superUsers[user] = struct{}{}
	}
	for _, user := range config.ReadOnlyUsers {
		readOnlyUsers[user] = struct{}{}
	}

	return &RestrictMethodMiddleware{
		next:          next,
		basePath:      config.BasePath,
		userHeader:    config.UserHeader,
		superUsers:    superUsers,
		readOnlyUsers: readOnlyUsers,
	}, nil
}

func isWriteReq(method string) bool {
	return method == http.MethodPut || method == http.MethodPost || method == http.MethodDelete || method == http.MethodPatch
}

func (m *RestrictMethodMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := strings.TrimRight(req.URL.Path, "/")
	basePath := strings.TrimRight(m.basePath, "/")

	// If not base path, move on
	if reqPath != basePath && !strings.HasPrefix(reqPath, basePath+"/") {
		m.next.ServeHTTP(rw, req)
		return
	}

	// Extract the username from the specified header
	user := req.Header.Get(m.userHeader)

	// If super user, allow & move on
	if _, isSuperUser := m.superUsers[user]; isSuperUser {
		m.next.ServeHTTP(rw, req)
		return
	}

	// Determine if it is a write request
	isWriteReq := isWriteReq(req.Method)

	// If a write request, and a read-only user, deny access
	if _, isReadOnlyUser := m.readOnlyUsers[user]; isReadOnlyUser && isWriteReq {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	// If request exactly matches base path, then restrict writes
	if isWriteReq && reqPath == basePath {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	// Build the restricted path
	restrictedPath := basePath + "/" + user
	restrictedPathPrefix := basePath + "/" + user + "/"

	// If path doesn't match auth'd username, deny access
	if isWriteReq && reqPath != restrictedPath && !strings.HasPrefix(reqPath, restrictedPathPrefix) {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	// Continue to the next handler
	m.next.ServeHTTP(rw, req)
}
