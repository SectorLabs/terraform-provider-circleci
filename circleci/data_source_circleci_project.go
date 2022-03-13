package circleci

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	client "github.com/mrolla/terraform-provider-circleci/circleci/client"
)

func dataSourceCircleCIProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCircleCIProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the project",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The organization where the project is defined",
			},
		},
	}
}

func dataSourceCircleCIProjectRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	org := d.Get("organization").(string)
	name := d.Get("name").(string)

	project, err := c.GetProject(org, name)
	if err != nil {
		return err
	}

	d.SetId(project.ID)
	return nil
}
