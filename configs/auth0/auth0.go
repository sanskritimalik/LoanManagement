package auth0

import (
	"fmt"
	"net/url"
	"os"
)

var TokenFetchURL string
var CallbackEndpoint string
var Secret string
var GetDataForTokenFetchWithCode func(code string) url.Values

func Init() error {
	TokenFetchURL = fmt.Sprintf("https://%s/oauth/token", os.Getenv("Auth0Domain"))
	CallbackEndpoint = os.Getenv("Auth0CallbackEndpoint")
	Secret = os.Getenv("Auth0ClientSecret")
	defaultData := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("Auth0ClientID")},
		"client_secret": {os.Getenv("Auth0ClientSecret")},
		"redirect_uri":  {os.Getenv("Auth0CallbackURL")},
	}

	GetDataForTokenFetchWithCode = func(code string) url.Values {
		data := defaultData
		data["code"] = []string{code}
		return data
	}

	return nil
}
