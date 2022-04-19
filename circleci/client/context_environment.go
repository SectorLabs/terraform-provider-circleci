package client

import (
	"fmt"

	"github.com/CircleCI-Public/circleci-cli/api"
)

// CreateOrUpdateContextEnvironmentVariable creates a new context environment variable
func (c *Client) CreateOrUpdateContextEnvironmentVariable(context_name, variable, value string) error {
	// Find context ID
	ctx, err := c.GetContextByName(context_name)
	if err != nil {
		return fmt.Errorf("could not find context by name: %v", err)
	}

	// CreateEnvironmentVariable calls PUT and can be used to update an existing variable with a matching context/name
	return c.contexts.CreateEnvironmentVariable(ctx.ID, variable, value)
}

// ListContextEnvironmentVariables lists all environment variables for a given context
func (c *Client) ListContextEnvironmentVariables(context_name string) (*[]api.EnvironmentVariable, error) {
	// Find context ID
	ctx, err := c.GetContextByName(context_name)
	if err != nil {
		return nil, fmt.Errorf("could not find context by name: %v", err)
	}

	return c.contexts.EnvironmentVariables(ctx.ID)
}

// HasContextEnvironmentVariable lists all environment variables for a given context and checks whether the specified variable is defined.
// If either the context or the variable does not exist, it returns false.
func (c *Client) HasContextEnvironmentVariable(context_name, variable string) (bool, error) {
	envs, err := c.ListContextEnvironmentVariables(context_name)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}

		return false, err
	}

	for _, env := range *envs {
		if env.Variable == variable {
			return true, nil
		}
	}

	return false, nil
}

// DeleteContextEnvironmentVariable deletes a context environment variable by context ID and name
func (c *Client) DeleteContextEnvironmentVariable(context_name, variable string) error {
	// Find context ID
	ctx, err := c.GetContextByName(context_name)
	if err != nil {
		return fmt.Errorf("could not find context by name: %v", err)
	}

	return c.contexts.DeleteEnvironmentVariable(ctx.ID, variable)
}
