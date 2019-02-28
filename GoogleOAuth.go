package gopcp_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleUser struct {
	ID            string
	Email         string
	VerifiedEmail string
	Name          string
	GivenName     string
	FamilyName    string
	Link          string
	Picture       string
	Locale        string
	HD            string
}

func GoogleOAuthLogin(config *oauth2.Config, w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type UserHandler = func(string, error)

func GoogleOAuthCallback(config *oauth2.Config, w http.ResponseWriter, r *http.Request, handleUser UserHandler) {
	user, err := getUserInfo(config, r.FormValue("state"), r.FormValue("code"))
	handleUser(user.Email, err)
}

func getUserInfo(conf *oauth2.Config, state string, code string) (GoogleUser, error) {

	var googleUser GoogleUser
	if state != "state" {
		return googleUser, fmt.Errorf("Invalid OAuth state")
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return googleUser, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return googleUser, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return googleUser, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	json.Unmarshal([]byte(content), &googleUser)

	return googleUser, nil
}
