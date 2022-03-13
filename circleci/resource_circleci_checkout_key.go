package circleci

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	client "github.com/mrolla/terraform-provider-circleci/circleci/client"
)

func resourceCircleCICheckoutKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCICheckoutKeyCreate,
		Read:   resourceCircleCICheckoutKeyRead,
		Delete: resourceCircleCICheckoutKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The CircleCI organization.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
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

	organization, err := c.Organization(d.Get("organization").(string))
	if err != nil {
		return err
	}

	project := d.Get("project").(string)
	keyType := d.Get("type").(string)

	checkoutKey, err := c.CreateCheckoutKey(organization, project, keyType)
	if err != nil {
		return err
	}

	d.SetId(generateCheckoutKeyId(organization, project, keyType, checkoutKey.Fingerprint))

	d.Set("fingerprint", checkoutKey.Fingerprint)
	d.Set("public_key", checkoutKey.PublicKey)
	d.Set("preferred", checkoutKey.Preferred)
	d.Set("created_at", checkoutKey.CreatedAt)

	return resourceCircleCICheckoutKeyRead(d, m)
}

func resourceCircleCICheckoutKeyRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	// If we don't have a checkout fingerprint we're doing an import. Parse it from the ID.
	if _, ok := d.GetOk("fingerprint"); !ok {
		if err := setOrgProjectNameFromCheckoutKeyId(d); err != nil {
			return err
		}
	}

	organization, err := c.Organization(d.Get("organization").(string))
	if err != nil {
		return err
	}

	project := d.Get("project").(string)
	fingerprint := d.Get("fingerprint").(string)

	checkoutKey, err := c.GetCheckoutKey(organization, project, fingerprint)
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

	organization, err := c.Organization(d.Get("organization").(string))
	if err != nil {
		return err
	}

	project := d.Get("project").(string)
	fingerprint := d.Get("fingerprint").(string)

	err = c.DeleteCheckoutKey(organization, project, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to delete project checkout key: %w", err)
	}

	d.SetId("")

	return nil
}

func setOrgProjectNameFromCheckoutKeyId(d *schema.ResourceData) error {
	organization, projectName, keyType, fingerprint := parseCheckoutKeyId(d.Id())
	// Validate that he have values for all the ID segments. This should be at least 3
	if organization == "" || projectName == "" || keyType == "" || fingerprint == "" {
		return fmt.Errorf("error calculating circleci_checkout_key. Please make sure the ID is in the form ORGANIZATION.PROJECTNAME.TYPE.FINGERPRINT (i.e. foo.bar.type.my_fingerprint)")
	}

	_ = d.Set("organization", organization)
	_ = d.Set("project", projectName)
	_ = d.Set("fingerprint", fingerprint)
	_ = d.Set("type", keyType)
	return nil
}

func parseCheckoutKeyId(id string) (organization, projectName, keyType, fingerprint string) {
	parts := strings.Split(id, ".")

	if len(parts) >= 3 {
		organization = parts[0]
		projectName = strings.Join(parts[1:len(parts)-2], ".")
		keyType = parts[len(parts)-2]
		fingerprint = parts[len(parts)-1]
	}

	return organization, projectName, keyType, fingerprint
}

func generateCheckoutKeyId(organization, projectName, keyType, fingerprint string) string {
	vars := []string{
		organization,
		projectName,
		keyType,
		fingerprint,
	}
	return strings.Join(vars, ".")
}
