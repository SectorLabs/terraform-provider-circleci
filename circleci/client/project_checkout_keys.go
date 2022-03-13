package client

import (
	"fmt"
	"net/url"
)

type CheckoutKey struct {
	PublicKey   string `json:"public_key"`
	Type        string `json:"type"`
	Fingerprint string `json:"fingerprint"`
	Preferred   bool   `json:"preferred"`
	CreatedAt   string `json:"created_at"`
}

// HasProjectCheckoutKey checks if an existing project contains checkout key by its fingerprint
func (c *Client) HasProjectCheckoutKey(org, project, fingerprint string) (bool, error) {
	slug, err := c.Slug(org, project)
	if err != nil {
		return false, err
	}

	req, err := c.rest.NewRequest("GET", &url.URL{Path: fmt.Sprintf("project/%s/checkout-key/%s", slug, fingerprint)}, nil)
	if err != nil {
		return false, err
	}

	key := &CheckoutKey{}
	if _, err := c.rest.DoRequest(req, key); err != nil {
		return false, nil
	}

	return true, nil
}

// GetCheckoutKey gets an existing project's checkout key by its fingerprint
func (c *Client) GetCheckoutKey(org, project, fingerprint string) (*CheckoutKey, error) {
	slug, err := c.Slug(org, project)
	if err != nil {
		return nil, err
	}

	req, err := c.rest.NewRequest("GET", &url.URL{Path: fmt.Sprintf("project/%s/checkout-key/%s", slug, fingerprint)}, nil)
	if err != nil {
		return nil, err
	}

	key := &CheckoutKey{}
	if _, err := c.rest.DoRequest(req, key); err != nil {
		return nil, err
	}

	return key, nil
}

type createCheckoutKey struct {
	Type string `json:"type"`
}

// CreateCheckoutKey creates a new checkout key and returns the created object
func (c *Client) CreateCheckoutKey(org, project, keyType string) (*CheckoutKey, error) {
	slug, err := c.Slug(org, project)
	if err != nil {
		return nil, err
	}

	req, err := c.rest.NewRequest("POST", &url.URL{Path: fmt.Sprintf("project/%s/checkout-key", slug)}, &createCheckoutKey{
		Type: keyType,
	})

	if err != nil {
		return nil, err
	}

	key := &CheckoutKey{}
	if _, err := c.rest.DoRequest(req, key); err != nil {
		return nil, err
	}

	return key, nil
}

// DeleteCheckoutKey deletes an existing checkout key and returns the created object
func (c *Client) DeleteCheckoutKey(org, project, fingerprint string) error {
	slug, err := c.Slug(org, project)
	if err != nil {
		return err
	}

	req, err := c.rest.NewRequest("DELETE", &url.URL{Path: fmt.Sprintf("project/%s/checkout-key/%s", slug, fingerprint)}, nil)
	if err != nil {
		return err
	}

	_, err = c.rest.DoRequest(req, nil)
	return err
}
