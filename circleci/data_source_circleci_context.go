package circleci

import (
	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCircleCIContext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCircleCIContextRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the context",
			},
		},
	}
}

func dataSourceCircleCIContextRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	name := d.Get("name").(string)

	ctx, err := c.GetContextByName(name)
	if err != nil {
		return err
	}

	d.SetId(ctx.ID)
	return nil
}
