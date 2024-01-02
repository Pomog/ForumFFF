package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/google/uuid"
)

func (m *Repository) LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		// endpoint: auth/google
		authURL := "https://accounts.google.com/o/oauth2/auth"
		params := url.Values{}
		params.Add("client_id", m.App.GoogleClientID)
		params.Add("redirect_uri", m.App.GoogleRedirectURL)
		params.Add("scope", "profile email")
		params.Add("response_type", "code")

		http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
	}
}

func (m *Repository) CallbackGoogleHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		setErrorAndRedirect(w, r, "Could not found code", "/error-page")
		return
	}

	// exchange code for token
	access_token, err := exchangeCodeForTokenGoogle(code, m)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get access token from: exchangeCodeForTokenGoogle(code, m) - "+err.Error(), "/error-page")
		return
	}
	data, err := googleGetUserProfile(access_token, w, r)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get data from: googleGetUserProfile(access_token,w,r) - "+err.Error(), "/error-page")
		return
	}
	user, err := processUserDataGoogle(data)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user from: processUserDataGoogle(data) - "+err.Error(), "/error-page")
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
		setErrorAndRedirect(w, r, "Unfortunately you cannot log in with your Google account.", "/error-page")
		return
	}
}

func processUserDataGoogle(data map[string]interface{}) (user models.User, err error) {
	avatar, err := processAvatarURL(data["picture"].(string), data["name"].(string))
	if err != nil {
		return user, err
	}
	googleLoginPassword, err := generatePassword(data["email"].(string))
	if err != nil {
		return user, err
	}
	// Check for key existence before type assertion
	name, nameExists := data["name"].(string)
	givenName, givenNameExists := data["given_name"].(string)
	familyName, familyNameExists := data["family_name"].(string)
	email, emailExists := data["email"].(string)

	// Verify if the required keys exist and have the expected types
	if !nameExists || !givenNameExists || !familyNameExists || !emailExists {
		return user, errors.New("missing or invalid data keys")
	}

	user = models.User{
		UserName:  name,
		Password:  googleLoginPassword,
		FirstName: givenName,
		LastName:  familyName,
		Email:     email,
		Picture:   avatar,
	}

	return user, nil
}

func convertBodyToJson(r *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// convert body to map
	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}

func exchangeCodeForTokenGoogle(code string, m *Repository) (access_token string, err error) {
	var tokenExchangeUrl string
	var format string
	var access_tokenExist bool

	params := url.Values{}

	tokenExchangeUrl = "https://oauth2.googleapis.com/token"
	params.Add("client_id", m.App.GoogleClientID)
	params.Add("client_secret", m.App.GoogleClientSecret)
	params.Add("redirect_uri", m.App.GoogleRedirectURL)
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

	access_token, access_tokenExist = bodyMap["access_token"].(string)
	if !access_tokenExist {
		return access_token, errors.New("missing or invalid access_token data keys")
	}
	return access_token, nil

}

func GoogleGetUserProfile(access_token string, w http.ResponseWriter, r *http.Request) (user models.User, err error) {
	userInfoUrl := "https://www.googleapis.com/oauth2/v1/userinfo"
	params := url.Values{}
	params.Add("access_token", access_token)

	response, err := http.Get(userInfoUrl + "?" + params.Encode())

	if err != nil {
		return
	}

	data, err := convertBodyToJson(response)

	if err != nil {
		return
	}

	user = models.User{
		UserName:  data["name"].(string),
		Email:     data["email"].(string),
		FirstName: data["given_name"].(string),
		LastName:  data["family_name"].(string),
		Password:  "FALSE",
	}

	return user, nil
}

func googleGetUserProfile(access_token string, w http.ResponseWriter, r *http.Request) (data map[string]interface{}, err error) {
	userInfoUrl := "https://www.googleapis.com/oauth2/v1/userinfo"
	params := url.Values{}
	params.Add("access_token", access_token)

	response, err := http.Get(userInfoUrl + "?" + params.Encode())

	if err != nil {
		return
	}

	data, err = convertBodyToJson(response)

	if err != nil {
		return
	}

	return data, nil
}
