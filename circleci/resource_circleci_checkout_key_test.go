package circleci

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	client "github.com/mrolla/terraform-provider-circleci/circleci/client"
)

var (
	keyTypes = []string{"deploy-key", "user-key"}
	idRegEx  = regexp.MustCompile(`(?m).+\..+\..+\.(?:[0-9a-f]{2}\:){15}[0-9a-f]{2}`)
)

func TestAccCircleCICheckoutKeyOrganizationNotSet(t *testing.T) {
	for _, keyType := range keyTypes {
		project := "TEST_" + acctest.RandString(8)

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers: testAccNoOrgProviders,
			Steps: []resource.TestStep{
				{
					Config:      testAccCircleCICheckoutKeyConfigProviderOrg(project, keyType),
					ExpectError: regexp.MustCompile("organization is required"),
				},
			},
		})
	}
}

func TestAccCircleCICheckoutKeyCreateThenUpdateProviderOrg(t *testing.T) {
	for _, keyType := range keyTypes {
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccOrgProviders,
			CheckDestroy: testAccCircleCICheckoutKeyProviderOrgCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfigProviderOrg(project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
				{
					Config: testAccCircleCICheckoutKeyConfigProviderOrg(project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
			},
		})
	}
}

func TestAccCircleCICheckoutKeyCreateThenUpdateResourceOrg(t *testing.T) {
	for _, keyType := range keyTypes {
		organization := os.Getenv("TEST_CIRCLECI_ORGANIZATION")
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccOrgProviders,
			CheckDestroy: testAccCircleCICheckoutKeyResourceOrgCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfigResourceOrg(organization, project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
				{
					Config: testAccCircleCICheckoutKeyConfigResourceOrg(organization, project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
			},
		})
	}
}

func TestAccCircleCICheckoutKeyImportProviderOrg(t *testing.T) {
	for _, keyType := range keyTypes {
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccOrgProviders,
			CheckDestroy: testAccCircleCICheckoutKeyProviderOrgCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfigProviderOrg(project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
				{
					ResourceName:      fmt.Sprintf("circleci_environment_variable.%s", keyType),
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"value",
					},
				},
			},
		})
	}
}

func TestAccCircleCICheckoutKeyImportResourceOrg(t *testing.T) {
	for _, keyType := range keyTypes {
		organization := os.Getenv("TEST_CIRCLECI_ORGANIZATION")
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccOrgProviders,
			CheckDestroy: testAccCircleCICheckoutKeyProviderOrgCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfigResourceOrg(organization, project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
				{
					ResourceName:      fmt.Sprintf("circleci_environment_variable.%s", keyType),
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"value",
					},
				},
			},
		})
	}
}

func TestParseCheckoutKeyId(t *testing.T) {
	organization := acctest.RandString(8)
	projectNames := []string{
		"TEST_" + acctest.RandString(8),
		"TEST-" + acctest.RandString(8),
		"TEST." + acctest.RandString(8),
		"TEST_" + acctest.RandString(8) + "." + acctest.RandString(8),
		"TEST-" + acctest.RandString(8) + "." + acctest.RandString(8),
		"TEST." + acctest.RandString(8) + "." + acctest.RandString(8),
	}

	for _, keyType := range keyTypes {
		for _, name := range projectNames {
			fingerprint := generateFingerprint()
			expectedId := fmt.Sprintf("%s.%s.%s.%s", organization, name, keyType, fingerprint)
			actualOrganization, actualProjectName, actualKeyType, actualFingerprint := parseCheckoutKeyId(expectedId)
			assert.Equal(t, organization, actualOrganization)
			assert.Equal(t, name, actualProjectName)
			assert.Equal(t, keyType, actualKeyType)
			assert.Equal(t, fingerprint, actualFingerprint)
		}
	}
}

func testAccCircleCICheckoutKeyResourceOrgCheckDestroy(s *terraform.State) error {
	c := testAccNoOrgProvider.Meta().(*client.Client)
	return testAccCircleCICheckoutKeyCheckDestroy(c, s)
}

func testAccCircleCICheckoutKeyProviderOrgCheckDestroy(s *terraform.State) error {
	c := testAccOrgProvider.Meta().(*client.Client)
	return testAccCircleCICheckoutKeyCheckDestroy(c, s)
}

func testAccCircleCICheckoutKeyCheckDestroy(c *client.Client, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circleci_checkout_key" {
			continue
		}

		organization := rs.Primary.Attributes["organization"]
		if organization == "" {
			v, err := c.Organization(organization)
			if err != nil {
				return err
			}

			organization = v
		}

		has, err := c.HasProjectCheckoutKey(organization, rs.Primary.Attributes["project"], rs.Primary.Attributes["fingerprint"])
		if err != nil {
			return err
		}

		if has {
			return errors.New("Checkout key should have been destroyed")
		}
	}

	return nil
}

func testAccCircleCICheckoutKeyConfigProviderOrg(project, keyType string) string {
	return fmt.Sprintf(`
resource "circleci_checkout_key" "%[2]s" {
  project = "%[1]s"
  type    = "%[2]s"
}`, project, keyType)
}

func testAccCircleCICheckoutKeyConfigResourceOrg(organization, project, keyType string) string {
	return fmt.Sprintf(`
resource "circleci_checkout_key" "%[2]s" {
  organization = "%[3]s"
  project      = "%[1]s"
  type         = "%[2]s"
}`, project, keyType, organization)
}

func generateFingerprint() string {
	blocks := [16]string{}
	for i := range blocks {
		blocks[i] = acctest.RandStringFromCharSet(2, "0123456789abcdef")
	}

	return strings.Join(blocks[:], ":")
}
