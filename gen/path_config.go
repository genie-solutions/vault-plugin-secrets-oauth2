package gen

import (
	"context"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"encoding/json"

)


// pathInfo corresponds to READ gen/info.
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	var config Config
	rawConfig, err := req.Storage.Get(ctx, "config/" + name)
	json.Unmarshal(rawConfig.Value, &config)

	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name": name,
			"username": config.Username,
			"password": config.Password,
			"client_id": config.ClientId,
			"client_secret": config.ClientSecret,
		},
	}, nil
}


func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	client_id := d.Get("client_id").(string)
	client_secret := d.Get("client_secret").(string)
	
	config := Config{username, password, client_id, client_secret}
	entry, err := logical.StorageEntryJSON("config/" + name, config)

	if err != nil {
		return nil, err
	}

	req.Storage.Put(ctx, entry)

	// Verify saved data
	var newConfig Config
	raw, _ := req.Storage.Get(ctx, "config/" + name)
	json.Unmarshal(raw.Value, &newConfig)

	return &logical.Response{
		Data: map[string]interface{}{
			"name": name,
			"username": newConfig.Username,
			"password": newConfig.Password,
			"client_id": newConfig.ClientId,
			"client_secret": newConfig.ClientSecret,
		},
	}, nil
}


type Config struct {
	Username string
	Password string
	ClientId string
	ClientSecret string
}

