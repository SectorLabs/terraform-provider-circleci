package circleci

import (
	"fmt"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCircleCIContextEnvironmentVariable() *schema.Resource {
	return &schema.Resource{
		// Create and Update have the same implementation, since the upstream API uses PUT
		Create: resourceCircleCIContextEnvironmentVariableCreate,
		Update: resourceCircleCIContextEnvironmentVariableCreate,
		Read:   resourceCircleCIContextEnvironmentVariableRead,
		Delete: resourceCircleCIContextEnvironmentVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCircleCIContextEnvironmentVariableImport,
		},

		Schema: map[string]*schema.Schema{
			"context": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the context where the environment variable is defined",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the environment variable",
				ValidateFunc: validateEnvironmentVariableNameFunc,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				StateFunc: func(value interface{}) string {
					return hashString(value.(string))
				},
				Description: "The value that will be set for the environment variable.",
			},
		},
	}
}

func resourceCircleCIContextEnvironmentVariableCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	context := d.Get("context").(string)
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	if err := c.CreateOrUpdateContextEnvironmentVariable(context, name, value); err != nil {
		return fmt.Errorf("error storing environment variable: %w", err)
	}

	id, _ := c.ComposeElementId([]string{context, name})

	d.SetId(id)

	return resourceCircleCIContextEnvironmentVariableRead(d, m)
}

func resourceCircleCIContextEnvironmentVariableRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	context := d.Get("context").(string)
	name := d.Get("name").(string)

	has, err := c.HasContextEnvironmentVariable(context, name)
	if err != nil {
		return fmt.Errorf("failed to get context environment variables: %w", err)
	}

	if !has {
		d.SetId("")
	}

	return nil
}

func resourceCircleCIContextEnvironmentVariableDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	context := d.Get("context").(string)
	name := d.Get("name").(string)

	if err := c.DeleteContextEnvironmentVariable(context, name); err != nil {
		return fmt.Errorf("error deleting environment variable: %w", err)
	}

	return nil
}

func resourceCircleCIContextEnvironmentVariableImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*client.Client)

	parts, err := c.DecomposeElementId(d.Id(), []string{"context", "name"})
	if err != nil {
		return nil, err
	}

	context := parts["context"]
	name := parts["name"]

	if has, err := c.HasContextEnvironmentVariable(context, name); !has || err != nil {
		return nil, fmt.Errorf("environment variable does not exist: %v", err)
	}

	_ = d.Set("context", parts["context"])
	_ = d.Set("name", parts["name"])

	return []*schema.ResourceData{d}, nil
}
