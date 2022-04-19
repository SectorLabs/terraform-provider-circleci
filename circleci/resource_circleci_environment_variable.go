package circleci

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCircleCIEnvironmentVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCIEnvironmentVariableCreate,
		Read:   resourceCircleCIEnvironmentVariableRead,
		Delete: resourceCircleCIEnvironmentVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCircleCIEnvironmentVariableImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name of the CircleCI project to create the variable in",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description:  "The name of the environment variable",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateEnvironmentVariableNameFunc,
			},
			"value": {
				Description: "The value of the environment variable",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				StateFunc: func(value interface{}) string {
					/* To avoid storing the value of the environment variable in the state
					but still be able to know when the value change, we store a hash of the value.
					*/
					return hashString(value.(string))
				},
			},
		},
	}
}

// hashString do a sha256 checksum, encode it in base64 and return it as string
// The choice of sha256 for checksum is arbitrary.
func hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func resourceCircleCIEnvironmentVariableCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	has, err := c.HasProjectEnvironmentVariable(project, name)
	if err != nil {
		return err
	}

	if has {
		return fmt.Errorf("environment variable already exists: %s", name)
	}

	if err := c.CreateProjectEnvironmentVariable(project, name, value); err != nil {
		return fmt.Errorf("failed to create environment variable: %w", err)
	}

	id, _ := c.ComposeElementId([]string{project, name})

	d.SetId(id)

	return resourceCircleCIEnvironmentVariableRead(d, m)
}

func resourceCircleCIEnvironmentVariableRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	name := d.Get("name").(string)

	has, err := c.HasProjectEnvironmentVariable(project, name)
	if err != nil {
		return fmt.Errorf("failed to get project environment variable: %w", err)
	}

	if !has {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceCircleCIEnvironmentVariableDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	project := d.Get("project").(string)
	name := d.Get("name").(string)

	if err := c.DeleteProjectEnvironmentVariable(project, name); err != nil {
		return fmt.Errorf("failed to delete project environment variable: %w", err)
	}

	d.SetId("")

	return nil
}

func resourceCircleCIEnvironmentVariableImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*client.Client)

	parts, err := c.DecomposeElementId(d.Id(), []string{"project", "name"})
	if err != nil {
		return nil, err
	}

	project := parts["project"]
	name := parts["name"]

	if has, err := c.HasProjectEnvironmentVariable(project, name); !has || err != nil {
		return nil, fmt.Errorf("environment variable does not exist: %w", err)
	}

	_ = d.Set("project", project)
	_ = d.Set("name", name)

	return []*schema.ResourceData{d}, nil
}
