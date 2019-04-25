package gen

import (
	"strings"
	"context"
	"strconv"
	"io/ioutil"
	"net/url"
	"net/http"
	"encoding/json"
	"time"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathInfo corresponds to READ gen/info.
func (b *backend) pathToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	refresh := d.Get("refresh").(bool)

	apiUrl := "https://login.microsoftonline.com/common/oauth2/token"

	var config Config
	rawConfig, err := req.Storage.Get(ctx, "config/" + name)
	json.Unmarshal(rawConfig.Value, &config)

	if err != nil {
		return nil, err
	}

	if refresh == false {
		raw, tokenErr := req.Storage.Get(ctx, "token/" + name)
		if tokenErr == nil && raw != nil {
			var accessToken AccessToken
			json.Unmarshal(raw.Value, &accessToken)

			//  return cached token if its TTL is more than 30 mins
			if (accessToken.ExpiresOn - int(time.Now().Unix())) > 60*30 {
				return &logical.Response{
					Data: map[string]interface{}{
						"access_token": accessToken.Token,
						"expires_on": accessToken.ExpiresOn,
					},
				}, nil
			}
		}
	}

	if err != nil {
		return nil, err
	}

	params := map[string]string{
	    "grant_type": "password",
	    "resource": "https://analysis.windows.net/powerbi/api",
	    "client_id": config.ClientId,
	    "client_secret": config.ClientSecret,
	    "username": config.Username,
	    "password": config.Password,
	}
	data := url.Values{}

	for k, v := range params { 
		data.Add(k, v)
	}

	client := &http.Client{}
    r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, _ := client.Do(r)
    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    bodyMap := make(map[string]string)
	json.Unmarshal(bodyBytes, &bodyMap)

	expiresOn, _ := strconv.Atoi(bodyMap["expires_on"])
	accessToken := AccessToken{bodyMap["access_token"], expiresOn}

	entry, err := logical.StorageEntryJSON("token/" + name, accessToken)
	req.Storage.Put(ctx, entry)

	return &logical.Response{
		Data: map[string]interface{}{
			"access_token": bodyMap["access_token"],
			"expires_on": bodyMap["expires_on"],
		},
	}, nil
}


type AccessToken struct {
	Token string
	ExpiresOn int
}
