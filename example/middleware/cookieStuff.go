package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"forum-authentication/config"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateCookie(w http.ResponseWriter, r *http.Request, userID int) http.Cookie {

	expiration := time.Now().Add(1 * time.Hour)
	random_string := GenerateRandomString(32)
	user_agent := r.Header.Get("User-Agent")

	encodedData := base64.StdEncoding.EncodeToString([]byte(random_string + "::" + user_agent))

	cookie := http.Cookie{
		Name:     "session-1",
		Value:    encodedData,
		Expires:  expiration,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	_, err := config.DB.Exec("INSERT INTO sessions (user_id, name, value, expiration) VALUES (?, ?, ?, ?)", userID, cookie.Name, random_string, cookie.Expires)
	if err != nil {
		fmt.Println("failed", err)
		return http.Cookie{}
	}

	return cookie
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "session-1",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}
