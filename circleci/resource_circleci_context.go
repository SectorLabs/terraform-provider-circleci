package circleci

import (
	"errors"
	"fmt"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCircleCIContext() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCIContextCreate,
		Read:   resourceCircleCIContextRead,
		Delete: resourceCircleCIContextDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCircleCIContextImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the context",
			},
		},
	}
}

func resourceCircleCIContextCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	name := d.Get("name").(string)

	ctx, err := c.CreateContext(name)
	if err != nil {
		return fmt.Errorf("error creating context: %w", err)
	}

	d.SetId(ctx.ID)

	return resourceCircleCIContextRead(d, m)
}

func resourceCircleCIContextRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	id := d.Id()

	ctx, err := c.GetContext(id)
	if err != nil {
		if errors.Is(err, client.ErrContextNotFound) {
			d.SetId("")
			return nil
		}

		return err
	}

	_ = d.Set("name", ctx.Name)

	return nil
}

func resourceCircleCIContextDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	if err := c.DeleteContext(d.Id()); err != nil {
		return fmt.Errorf("error deleting context: %w", err)
	}

	return nil
}

func resourceCircleCIContextImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*client.Client)

	context_name := d.Id()

	ctx, err := c.GetContextByName(context_name)
	if err != nil {
		return nil, fmt.Errorf("context does not exist: %w", err)
	}

	d.SetId(ctx.ID)
	_ = d.Set("name", ctx.Name)

	return []*schema.ResourceData{d}, nil
}
