package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/SectorLabs/terraform-provider-circleci/circleci/client/rest"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/settings"
)

// Client provides access to the CircleCI REST API
// It uses upstream client functionality where possible and defines its own methods as needed
type Client struct {
	contexts     *api.ContextRestClient
	rest         *rest.Client
	vcs          string
	organization string
}

// Config configures a Client
type Config struct {
	URL   string
	Token string

	VCS          string
	Organization string
}

// New initializes a client object for the provider
func New(config Config) (*Client, error) {
	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	rootURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	contexts, err := api.NewContextRestClient(settings.Config{
		Host:         rootURL,
		RestEndpoint: u.Path,
		Token:        config.Token,
		HTTPClient:   http.DefaultClient,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		rest:     rest.New(rootURL, u.Path, config.Token),
		contexts: contexts,

		vcs:          config.VCS,
		organization: config.Organization,
	}, nil
}

// Organization returns the organization for a request. The organization configured
// in the provider is returned.
func (c *Client) Organization() string {
	return c.organization
}

// Slug returns a project slug, including the VCS, organization, and project names
func (c *Client) Slug(project string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", c.vcs, c.organization, project), nil
}

func (c *Client) DecomposeElementId(id string, identifiers []string) (map[string]string, error) {
	parts := strings.Split(id, "/")

	parent := identifiers[0]
	identifiers = identifiers[1:]

	out := map[string]string{}
	if len(parts) >= 2 {
		out[parent] = strings.Join(parts[0:len(parts)-len(identifiers)], "/")

		for i, identifier := range identifiers {
			out[identifier] = parts[len(parts)-len(identifiers)+i]
		}
	}

	idSyntax := strings.ToUpper(
		fmt.Sprintf("%s/%s", parent, strings.Join(identifiers, "/")),
	)
	composeError := fmt.Errorf("error computing the id. Please make sure the ID is in the form %s", idSyntax)

	if out[parent] == "" {
		return nil, composeError
	}

	for _, identifier := range identifiers {
		if out[identifier] == "" {
			return nil, composeError
		}
	}

	return out, nil
}

func (c *Client) ComposeElementId(identifiers []string) (string, error) {
	return strings.Join(identifiers, "/"), nil
}

func isNotFound(err error) bool {
	var httpError *rest.HTTPError
	if errors.As(err, &httpError) && httpError.Code == 404 {
		return true
	}

	return false
}
