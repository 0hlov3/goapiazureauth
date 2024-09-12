package webserver

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func (c apiConfig) verifyToken(tokenString string, r *http.Request) bool {
	keySet, err := jwk.Fetch(r.Context(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		c.config.Log.Error("Failed to fetch key set", zap.Error(err))
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("Key %v not found", kid)
		}

		var publicKey rsa.PublicKey
		if err := keys.Raw(&publicKey); err != nil {
			return nil, fmt.Errorf("Could not parse public key")
		}

		return &publicKey, nil
	})

	if err != nil {
		c.config.Log.Error("Token verification failed", zap.String("Token", tokenString), zap.Error(err))
		return false
	}

	// Validate claims (e.g., exp, aud, iss, etc.)
	if !token.Valid {
		c.config.Log.Error("Invalid token")
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	// Example claim checks (modify as needed for your use case)
	if claims["iss"] != fmt.Sprintf("https://sts.windows.net/%s/", c.config.Azure.TenantID) {
		c.config.Log.Error("Invalid issuer", zap.String("Issuer", claims["iss"].(string)))
		return false
	}

	if claims["aud"] != "api://example.jwt.application.auth" {
		c.config.Log.Error("Invalid audience", zap.String("Audience", claims["aud"].(string)))
		return false
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			c.config.Log.Error("Token has expired", zap.Time("Expiration", time.Unix(int64(exp), 0)))
			return false
		}
	}

	return true
}

func (c apiConfig) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString = extractBearerToken(tokenString)
		if tokenString == "" {
			http.Error(w, "Bearer token missing", http.StatusUnauthorized)
			return
		}

		if !c.verifyToken(tokenString, r) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
