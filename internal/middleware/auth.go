package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

var secretKey = []byte("esg-visa-auth-secret-key-2026")

type TokenPayload struct {
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// GenerateToken generates a signed token for the given email
func GenerateToken(email string) string {
	header := tokenHeader{Alg: "HS256", Typ: "JWT"}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	payload := TokenPayload{
		Email: email,
		Exp:   time.Now().Add(72 * time.Hour).Unix(),
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signingInput := headerB64 + "." + payloadB64
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig
}

// ValidateToken validates a token and returns the email
func ValidateToken(tokenStr string) (string, bool) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return "", false
	}

	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return "", false
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}

	var payload TokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return "", false
	}

	if time.Now().Unix() > payload.Exp {
		return "", false
	}

	return payload.Email, true
}

// GetEmailFromRequest extracts authenticated email from request
// Priority: 1. Authorization header (Bearer token)  2. Query param (backward compat)
func GetEmailFromRequest(r *ghttp.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if email, valid := ValidateToken(token); valid {
			return email
		}
	}

	// Fallback to query/form parameter (backward compatibility)
	email := r.Get("email").String()
	if email == "" {
		email = r.GetForm("email").String()
	}
	return email
}

// AuthRequired is middleware that requires a valid token
func AuthRequired(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if email, valid := ValidateToken(token); valid {
			r.SetParam("auth_email", email)
			r.Middleware.Next()
			return
		}
	}
	r.Response.WriteJson(map[string]interface{}{
		"code":    401,
		"message": "Authentication required",
	})
}
