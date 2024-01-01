package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"forum-authentication/middleware"
	"forum-authentication/types"
)

const (
	google_ClientID     = "89143036124-72c155p6tigsp9l1ch3ud520i9bho94f.apps.googleusercontent.com"
	google_ClientSecret = "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"
	google_RedirectURI  = "http://localhost:8080/auth/google/callback"

	github_ClientID     = "2ac7d35edf087740ae48"
	github_ClientSecret = "425b0e0bacb5616b076f37f2be164b88698767d5"
	github_RedirectURI  = "http://localhost:8080/auth/github/callback"
)

func exchangeCodeForToken(code, provider string) (access_token string, err error) {
	var tokenExchangeUrl string
	var format string

	params := url.Values{}

	switch provider {
	case "google":
		tokenExchangeUrl = "https://oauth2.googleapis.com/token"
		params.Add("client_id", google_ClientID)
		params.Add("client_secret", google_ClientSecret)
		params.Add("redirect_uri", google_RedirectURI)
		params.Add("grant_type", "authorization_code")
		params.Add("code", code)
		format = "application/x-www-form-urlencoded"

		resp, _ := http.Post(tokenExchangeUrl, format, strings.NewReader(params.Encode()))
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		bodyMap, err := convertBodyToJson(resp)
		if err != nil {
			return "", err

		}

		access_token = bodyMap["access_token"].(string)
		return access_token, nil

	case "github":
		tokenExchangeUrl = "https://github.com/login/oauth/access_token"
		params.Add("client_id", github_ClientID)
		params.Add("client_secret", github_ClientSecret)
		params.Add("redirect_uri", github_RedirectURI)
		params.Add("code", code)
		format = "application/x-www-form-urlencoded"

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
			return "", errors.New("Access token not found in response")
		}

		return access_token, nil

	default:
		return "", errors.New("Provider not found")

	}
}

func GoogleGetUserProfile(access_token string, w http.ResponseWriter, r *http.Request) (user types.User, err error) {
	userInfoUrl := "https://www.googleapis.com/oauth2/v1/userinfo"
	params := url.Values{}
	params.Add("access_token", access_token)

	response, err := http.Get(userInfoUrl + "?" + params.Encode())

	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := convertBodyToJson(response)

	if err != nil {
		fmt.Println(err)
		return
	}

	user = types.User{
		Username:  data["name"].(string),
		Email:     data["email"].(string),
		FirstName: data["given_name"].(string),
		LastName:  data["family_name"].(string),
		Password:  "FALSE",
		Provider:  "google",
	}

	return user, nil
}
func GithubRequest(url, access_token string, target interface{}) error {
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

func GithubRequestEmail(url, access_token string) (string, error) {
	var data []map[string]interface{}
	if err := GithubRequest(url, access_token, &data); err != nil {
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

func GithubRequestUser(url, access_token string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := GithubRequest(url, access_token, &data); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

func GithubGetUserProfile(access_token string, w http.ResponseWriter, r *http.Request) (user types.User, err error) {
	var email string
	var user_raw map[string]interface{}

	email, err = GithubRequestEmail("https://api.github.com/user/emails", access_token)
	if err != nil {
		fmt.Println(err)
		return
	}

	user_raw, err = GithubRequestUser("https://api.github.com/user", access_token)
	user = types.User{
		Username: user_raw["login"].(string),
		Email:    email,
		Password: "FALSE",
		Provider: "github",
	}

	return user, nil
}

func (_ *UserController) AuthCallback(w http.ResponseWriter, r *http.Request) {
	// endpoint: auth

	var provider string
	path := r.URL.Path

	fmt.Println(path)
	if strings.Contains(path, "google") {
		provider = "google"
	} else if strings.Contains(path, "github") {
		provider = "github"
	} else {
		fmt.Println("Provider not found")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		fmt.Println("Code not found")
		return
	}

	// exchange code for token
	access_token, err := exchangeCodeForToken(code, provider)
	if err != nil {
		fmt.Println("Error exchanging code for access_token:", err)
		return
	}

	var user types.User

	switch provider {
	case "google":
		user, err = GoogleGetUserProfile(access_token, w, r)
		break

	case "github":
		fmt.Println(access_token)
		user, err = GithubGetUserProfile(access_token, w, r)
		break
	}

	if err != nil {
		fmt.Println("Error getting user profile:", err)
		return
	}

	var userId int
	existingUser, err := user.GetUserByEmail(user.Email)

	if err != nil {
		// if user does not exist, create user
		userId, err = user.CreateUser(user, provider)

		if err != nil || userId == 0 {
			fmt.Println("Error creating user:", err)
			return
		}
	} else {
		// if user exists, get user id
		userId = existingUser.Id
	}

	cookie := middleware.GenerateCookie(w, r, userId)

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (_ *UserController) GoogleAuth(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		// endpoint: auth/google
		authURL := "https://accounts.google.com/o/oauth2/auth"
		params := url.Values{}
		fmt.Println(google_RedirectURI)
		params.Add("client_id", google_ClientID)
		params.Add("redirect_uri", google_RedirectURI)
		params.Add("scope", "profile email")
		params.Add("response_type", "code")

		http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
	}
}

func (_ *UserController) GithubAuth(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		// endpoint: auth/github
		authURL := "https://github.com/login/oauth/authorize"
		params := url.Values{}
		params.Add("client_id", github_ClientID)
		params.Add("redirect_uri", github_RedirectURI)
		params.Add("scope", "user user:email")

		http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
	}
}
