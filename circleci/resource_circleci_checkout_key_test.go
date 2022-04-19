package circleci

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

var (
	keyTypes = []string{"deploy-key", "user-key"}
	idRegEx  = regexp.MustCompile(`(?m).+\/(?:[0-9a-f]{2}\:){15}[0-9a-f]{2}`)
)

func TestAccCircleCICheckoutKeyCreateThenUpdate(t *testing.T) {
	for _, keyType := range keyTypes {
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCircleCICheckoutKeyCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfig(project, keyType),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(resourceName, "id", idRegEx),
						resource.TestCheckResourceAttr(resourceName, "project", project),
						resource.TestCheckResourceAttr(resourceName, "type", keyType),
					),
				},
				{
					Config: testAccCircleCICheckoutKeyConfig(project, keyType),
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

func TestAccCircleCICheckoutKeyImport(t *testing.T) {
	for _, keyType := range keyTypes {
		project := os.Getenv("CIRCLECI_PROJECT")
		resourceName := "circleci_checkout_key." + keyType

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCircleCICheckoutKeyCheckDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCircleCICheckoutKeyConfig(project, keyType),
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
	orgs := []string{
		acctest.RandString(8),
		"",
	}

	testAccProvider, _ := client.New(client.Config{
		Organization: orgs[0],
	})

	clients := []*client.Client{
		testAccProvider,
	}

	for i := range clients {
		c := clients[i]
		projectNames := []string{
			"TEST_" + acctest.RandString(8),
			"TEST-" + acctest.RandString(8),
			"TEST." + acctest.RandString(8),
			"TEST_" + acctest.RandString(8) + "." + acctest.RandString(8),
			"TEST-" + acctest.RandString(8) + "." + acctest.RandString(8),
			"TEST." + acctest.RandString(8) + "." + acctest.RandString(8),
		}

		for _, name := range projectNames {
			fingerprint := generateFingerprint()
			expectedId := fmt.Sprintf("%s/%s", name, fingerprint)

			parts, err := c.DecomposeElementId(expectedId, []string{"project", "fingerprint"})
			assert.Equal(t, 2, len(parts))
			assert.Equal(t, nil, err)

			actualProjectName := parts["project"]
			actualFingerprint := parts["fingerprint"]

			assert.Equal(t, name, actualProjectName)
			assert.Equal(t, fingerprint, actualFingerprint)
		}
	}
}

func testAccCircleCICheckoutKeyCheckDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circleci_checkout_key" {
			continue
		}

		has, err := c.HasProjectCheckoutKey(rs.Primary.Attributes["project"], rs.Primary.Attributes["fingerprint"])
		if err != nil {
			return err
		}

		if has {
			return errors.New("Checkout key should have been destroyed")
		}
	}

	return nil
}

func testAccCircleCICheckoutKeyConfig(project, keyType string) string {
	return fmt.Sprintf(`
resource "circleci_checkout_key" "%[2]s" {
  project = "%[1]s"
  type    = "%[2]s"
}`, project, keyType)
}

func generateFingerprint() string {
	blocks := [16]string{}
	for i := range blocks {
		blocks[i] = acctest.RandStringFromCharSet(2, "0123456789abcdef")
	}

	return strings.Join(blocks[:], ":")
}
