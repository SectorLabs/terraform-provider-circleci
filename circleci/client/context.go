package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/google/uuid"
)

var ErrContextNotFound = errors.New("context not found")

// GetContext gets an existing context by its ID (UUID)
func (c *Client) GetContext(id string) (*api.Context, error) {
	req, err := c.rest.NewRequest("GET", &url.URL{Path: fmt.Sprintf("context/%s", id)}, nil)
	if err != nil {
		return nil, err
	}

	ctx := &api.Context{}

	status, err := c.rest.DoRequest(req, ctx)
	if err != nil {
		if status == 404 {
			return nil, ErrContextNotFound
		}

		return nil, err
	}

	return ctx, nil
}

// GetContextByName gets an existing context by its name
func (c *Client) GetContextByName(name string) (*api.Context, error) {
	return c.contexts.ContextByName(c.vcs, c.organization, name)
}

// GetContextByIDOrName gets a context by ID if a UUID is specified, and by name otherwise
func (c *Client) GetContextByIDOrName(id string) (*api.Context, error) {
	if _, uuidErr := uuid.Parse(id); uuidErr == nil {
		return c.GetContext(id)
	} else {
		return c.contexts.ContextByName(c.vcs, c.organization, id)
	}
}

type createContextRequest struct {
	Name  string        `json:"name"`
	Owner *contextOwner `json:"owner"`
}

type contextOwner struct {
	Slug string `json:"slug"`
	Type string `json:"type"`
}

// CreateContext creates a new context and returns the created context object
func (c *Client) CreateContext(name string) (*api.Context, error) {
	req, err := c.rest.NewRequest("POST", &url.URL{Path: "context"}, &createContextRequest{
		Name: name,
		Owner: &contextOwner{
			Slug: fmt.Sprintf("%s/%s", c.vcs, c.organization),
			Type: "organization",
		},
	})
	if err != nil {
		return nil, err
	}

	ctx := &api.Context{}
	_, err = c.rest.DoRequest(req, ctx)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (c *Client) DeleteContext(id string) error {
	return c.contexts.DeleteContext(id)
}
