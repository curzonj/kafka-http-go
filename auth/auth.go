package auth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
)

const (
	DEFAULT_PASSWORD = "changeme"
)

var (
	secretkey string
)

// Partially copied from github.com/ryandotsmith/l2met/auth/auth.go

func AuthenticateHttpRequest(w http.ResponseWriter, r *http.Request) bool {
	l := r.Header.Get("Authorization")
	user, err := Parse(l)
	if err != nil {
		http.Error(w, "Unable to parse headers.", 400)
		return false
	}

	matched := false
	if user == secretkey {
		matched = true
	}

	if !matched {
		http.Error(w, "Authentication failed.", 401)
		return false
	}

	return true
}

func init() {
	secretkey = os.Getenv("BASIC_AUTH")
	if secretkey == "" {
		secretkey = DEFAULT_PASSWORD
	}
}

// Extract the username and password from the authorization
// line of an HTTP header. This function will handle the
// parsing and decoding of the line.
func Parse(authLine string) (string, error) {
	parts := strings.SplitN(authLine, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("Authorization header malformed.")
	}
	method := parts[0]
	if method != "Basic" {
		return "", errors.New("Authorization must be basic.")
	}
	payload := parts[1]
	decodedPayload, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(decodedPayload), ":"), nil
}
