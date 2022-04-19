package circleci

import (
	"fmt"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceCircleCICheckoutKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCICheckoutKeyCreate,
		Read:   resourceCircleCICheckoutKeyRead,
		Delete: resourceCircleCICheckoutKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCircleCICheckoutKeyImport,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name of the CircleCI project to create the checkout key in.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description:  "The type of the checkout key. Can be either \"user-key\" or \"deploy-key\".",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"user-key", "deploy-key"}, false),
			},
			"fingerprint": {
				Description: "The fingerprint of the checkout key.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"public_key": {
				Description: "The public SSH key of the checkout key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"preferred": {
				Description: "A boolean value that indicates if this key is preferred.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "The date and time the checkout key was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		SchemaVersion: 1,
	}
}

func resourceCircleCICheckoutKeyCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	keyType := d.Get("type").(string)

	checkoutKey, err := c.CreateCheckoutKey(project, keyType)
	if err != nil {
		return err
	}

	id, _ := c.ComposeElementId([]string{project, checkoutKey.Fingerprint})

	d.SetId(id)
	d.Set("fingerprint", checkoutKey.Fingerprint)
	d.Set("public_key", checkoutKey.PublicKey)
	d.Set("preferred", checkoutKey.Preferred)
	d.Set("created_at", checkoutKey.CreatedAt)

	return resourceCircleCICheckoutKeyRead(d, m)
}

func resourceCircleCICheckoutKeyRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	fingerprint := d.Get("fingerprint").(string)

	checkoutKey, err := c.GetCheckoutKey(project, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to get project checkout key: %w", err)
	}

	d.Set("fingerprint", checkoutKey.Fingerprint)
	d.Set("public_key", checkoutKey.PublicKey)
	d.Set("preferred", checkoutKey.Preferred)
	d.Set("created_at", checkoutKey.CreatedAt)

	return nil
}

func resourceCircleCICheckoutKeyDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	fingerprint := d.Get("fingerprint").(string)

	if err := c.DeleteCheckoutKey(project, fingerprint); err != nil {
		return fmt.Errorf("failed to delete project checkout key: %w", err)
	}

	d.SetId("")

	return nil
}

func resourceCircleCICheckoutKeyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*client.Client)

	parts, err := c.DecomposeElementId(d.Id(), []string{"project", "fingerprint"})
	if err != nil {
		return nil, err
	}

	project := parts["project"]
	fingerprint := parts["fingerprint"]

	checkoutKey, err := c.GetCheckoutKey(project, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("checkout key not exist: %w", err)
	}

	_ = d.Set("project", project)
	_ = d.Set("type", checkoutKey.Type)
	_ = d.Set("fingerprint", fingerprint)
	_ = d.Set("public_key", checkoutKey.PublicKey)
	_ = d.Set("preferred", checkoutKey.Preferred)
	_ = d.Set("created_at", checkoutKey.CreatedAt)

	return []*schema.ResourceData{d}, nil
}
