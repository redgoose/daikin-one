package daikin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Token struct {
	AccessToken          string `json:"accessToken"`
	AccessTokenExpiresIn int    `json:"accessTokenExpiresIn"`
	TokenType            string `json:"tokenType"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}
var urlBase string = "https://integrator-api.daikinskyport.com"

func GetToken() string {
	body := []byte(`{
		"email": "` + viper.GetString("email") + `",
		"integratorToken": "` + viper.GetString("integratorToken") + `"
	}`)

	r, err := http.NewRequest("POST", urlBase+"/v1/token", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-api-key", viper.GetString("apiKey"))

	res, err := httpClient.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	token := &Token{}
	derr := json.NewDecoder(res.Body).Decode(token)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	return token.AccessToken
}
