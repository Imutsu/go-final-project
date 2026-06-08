package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

var todoPassword string

func SetPassword(pass string) {
	todoPassword = pass
}

func hashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func createToken(pass string) (string, error) {
	claims := jwt.MapClaims{
		"hash": hashPassword(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(pass))
}

func validateToken(tokenString, pass string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	return hash == hashPassword(pass)
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(SigninResponse{
			Error: "invalid request",
		})
		return
	}

	if req.Password != todoPassword {
		json.NewEncoder(w).Encode(SigninResponse{
			Error: "Неверный пароль",
		})
		return
	}

	token, err := createToken(todoPassword)
	if err != nil {
		json.NewEncoder(w).Encode(SigninResponse{
			Error: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(SigninResponse{
		Token: token,
	})
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if todoPassword != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !validateToken(cookie.Value, todoPassword) {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	}
}
