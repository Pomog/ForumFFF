package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/google/uuid"
)

func (m *Repository) LoginWithGitHubHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		m.App.GitHubClientID, m.App.GitHubRedirectURL)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (m *Repository) CallbackGitHubHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		setErrorAndRedirect(w, r, "Code not provided", "/error-page")
		return
	}

	token, err := exchangeCodeForToken(code, m)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to exchange code for token: "+err.Error(), "/error-page")
		return
	}

	user_email, err := githubRequestEmail("https://api.github.com/user/emails", token)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to get email:  "+err.Error(), "/error-page")
		return
	}

	user_data, err := githubRequestUserData("https://api.github.com/user", token)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to get User Data: "+err.Error(), "/error-page")
		return
	}

	user, err := parseUserData(user_data, user_email)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to parse User: "+err.Error(), "/error-page")
		return
	}

	// generate password based on email
	user.Password, err = generatePassword(user.Email)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to generate password based on GitHub data: "+err.Error(), "/error-page")
		return
	}

	userExist, err := m.DB.UserPresent(user.UserName, user.Email)
	if err != nil {
		setErrorAndRedirect(w, r, "Failed to check User: "+err.Error(), "/error-page")
		return
	}
	if !userExist {
		err := m.DB.CreateUser(user)
		if err != nil {
			setErrorAndRedirect(w, r, "Failed to create user based on GitHub data: "+err.Error(), "/error-page")
			return
		}
	}

	// Check if User is Present in the DB, ERR should be handled
	userID, _ := m.DB.UserPresentLogin(user.Email, user.Password)
	if userID != 0 {
		m.App.UserLogin = uuid.New()
		err := m.DB.InsertSessionintoDB(m.App.UserLogin.String(), userID)
		if err != nil {
			setErrorAndRedirect(w, r, err.Error(), "/error-page")
			return
		}

		cookie := &http.Cookie{
			Name:  strconv.Itoa(userID),
			Value: m.App.UserLogin.String(),
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		setErrorAndRedirect(w, r, "Wrong email or password", "/error-page")
		return
	}

}

func githubRequest(url, access_token string, target interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access_token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}

func exchangeCodeForToken(code string, m *Repository) (access_token string, err error) {
	tokenExchangeUrl := "https://github.com/login/oauth/access_token"
	params := url.Values{}

	params.Add("client_id", m.App.GitHubClientID)
	params.Add("client_secret", m.App.GitHubClientSecret)
	params.Add("redirect_uri", m.App.GitHubRedirectURL)
	params.Add("code", code)
	format := "application/x-www-form-urlencoded"

	resp, err := http.Post(tokenExchangeUrl, format, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body := string(bodyBytes)
	parts := strings.Split(body, "&")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 2 && kv[0] == "access_token" {
			access_token = kv[1]
			break
		}
	}

	if access_token == "" {
		return "", errors.New("access token not found in response")
	}

	return access_token, nil
}

func githubRequestEmail(url, access_token string) (string, error) {
	var data []map[string]interface{}
	if err := githubRequest(url, access_token, &data); err != nil {
		return "", err
	}

	user_email := ""
	for _, v := range data {
		if v["primary"].(bool) {
			user_email = v["email"].(string)
			break
		}
	}

	return user_email, nil
}

func githubRequestUserData(url, access_token string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := githubRequest(url, access_token, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func parseUserData(data map[string]interface{}, email string) (models.User, error) {
	var user models.User

	avatar, err := processAvatarURL(data["avatar_url"].(string), data["login"].(string))
	if err != nil {
		return user, err
	}

	// Map the fields from the data to the User struct
	user.UserName = data["login"].(string)
	user.FirstName, user.LastName = splitName(data["name"].(string))
	user.Email = email
	user.Picture = avatar

	return user, nil
}

func splitName(name string) (string, string) {
	names := strings.Fields(name)
	if len(names) >= 2 {
		return names[0], strings.Join(names[1:], " ")
	}
	return names[0], ""
}

func processAvatarURL(url string, username string) (string, error) {
	// Fetch the image from the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Extract content type from the response headers
	contentType := resp.Header.Get("Content-Type")

	// Get extension from the content type
	extension, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extension) == 0 {
		// Use a default extension or handle the error as needed
		extension = []string{".jpg"}
	}

	newFileName := fmt.Sprintf("%s%s", username, extension[0])
	newFilePath := filepath.Join("static/ava", newFileName)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// Copy the fetched image to the new file
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		return "", err
	}

	return path.Join("/", newFilePath), nil
}

func generatePassword(email string) (string, error) {
	// Use the email as input
	password := []byte(email)

	// Hash the password using sha256
	hasher := sha256.New()
	hasher.Write(password)
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	return hashedPassword, nil
}
