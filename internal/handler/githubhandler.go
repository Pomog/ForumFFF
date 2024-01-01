package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (m *Repository) LoginWithGitHubHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		m.App.GitHubClientID, m.App.GitHubRedirectURL)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (m *Repository) CallbackGitHubHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		fmt.Println("Code not provided")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := exchangeCodeForToken(code, m)
	if err != nil {
		fmt.Println("Failed to exchange code for token:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Token: %s", token)

	email, err := githubRequestEmail("https://api.github.com/user/emails", token)
	if err != nil {
		fmt.Println("Failed to get Email: ", err)
		return
	}

	fmt.Println("**************************************************")
	fmt.Println(email)

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

	body, err := ioutil.ReadAll(resp.Body)
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
		fmt.Println(err)
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
