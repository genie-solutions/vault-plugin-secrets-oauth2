package gen

import (
	"context"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pkg/errors"
)

// Factory creates a new usable instance of this secrets engine.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, errors.Wrap(err, "failed to create factory")
	}
	return b, nil
}

// backend is the actual backend.
type backend struct {
	*framework.Backend
	logger log.Logger
}

// Backend creates a new backend.
func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.logger = c.Logger

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        backendHelp,
		Paths: []*framework.Path{
			// gen/info
			&framework.Path{
				Pattern:      "info",
				HelpSynopsis: "Display information about this plugin",
				HelpDescription: `

Displays information about the plugin, such as the plugin version and where to
get help.

`,
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInfo,
				},
			},

			// oauth2/powerbi/config/<name>
			&framework.Path{
				Pattern:      "powerbi/config/" + framework.GenericNameRegex("name"),
				HelpSynopsis: "Configure credentials to get access token",
				HelpDescription: `

Configure Power BI credentials to retrieve access token

`,
				Fields: map[string]*framework.FieldSchema{
					"name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Credentials' name. i.e: staging/production.",
						Default:     "-",
					},
					"username": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Username",
						Default:     "-",
					},
					"password": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Pasword",
						Default:     "-",
					},
					"client_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Client Id",
						Default:     "-",
					},
					"client_secret": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Client secret",
						Default:     "-",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathConfigWrite,
					logical.ReadOperation: b.pathConfigRead,
				},
			},

			// oauth2/powerbi/token/<name>
			&framework.Path{
				Pattern:      "powerbi/token/" + framework.GenericNameRegex("name"),
				HelpSynopsis: "Generate and return an OAuth access token for Power BI Service",
				HelpDescription: `

Invoke API to generate an access token with the existing credentials

`,
				Fields: map[string]*framework.FieldSchema{
					"name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Credentials' name. i.e: staging/production.",
						Default:     "-",
					},
					"refresh": &framework.FieldSchema{
						Type:        framework.TypeBool,
						Description: "We are caching access token for 30 mins to reduce API call to MS OAuth Server. Set refresh=true if you want to get a new access token.",
						Default:     false,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathToken,
				},
			},
		},
	}

	return &b
}

const backendHelp = `
The oauth2 secrets engine generate access token for OAuth2-authorized services. It currently supports MS PowerBI only.
`
