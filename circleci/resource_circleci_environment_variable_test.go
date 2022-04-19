package circleci

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccCircleCIEnvironmentVariableCreateThenUpdate(t *testing.T) {
	project := os.Getenv("CIRCLECI_PROJECT")
	envName := "TEST_" + acctest.RandString(8)
	resourceName := "circleci_environment_variable." + envName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCircleCIEnvironmentVariableCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircleCIEnvironmentVariableConfig(project, envName, "value-for-the-test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s/%s", project, envName)),
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", hashString("value-for-the-test")),
				),
			},
			{
				Config: testAccCircleCIEnvironmentVariableConfig(project, envName, "value-for-the-test-again"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s/%s", project, envName)),
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", hashString("value-for-the-test-again")),
				),
			},
		},
	})
}

func TestAccCircleCIEnvironmentVariableCreateAlreadyExists(t *testing.T) {
	project := os.Getenv("CIRCLECI_PROJECT")
	envName := "TEST_" + acctest.RandString(8)
	envValue := acctest.RandString(8)

	resourceName := "circleci_environment_variable." + envName

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCircleCIEnvironmentVariableCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircleCIEnvironmentVariableConfig(project, envName, envValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s/%s", project, envName)),
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", hashString(envValue)),
				),
			},
			{
				Config:      testAccCircleCIEnvironmentVariableConfigIdentical(project, envName, envValue),
				ExpectError: regexp.MustCompile("already exists"),
			},
		},
	})
}

func TestParseEnvironmentVariableId(t *testing.T) {
	orgs := []string{
		acctest.RandString(8),
		"",
	}

	testAccOrgProvider, _ := client.New(client.Config{
		Organization: orgs[0],
	})

	clients := []*client.Client{
		testAccOrgProvider,
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
			envName := acctest.RandString(8)
			expectedId := fmt.Sprintf("%s/%s", name, envName)

			parts, err := c.DecomposeElementId(expectedId, []string{"project", "name"})
			assert.Equal(t, 2, len(parts))
			assert.Equal(t, nil, err)

			actualProjectName := parts["project"]
			actualEnvName := parts["name"]

			assert.Equal(t, name, actualProjectName)
			assert.Equal(t, envName, actualEnvName)
		}
	}
}

func testAccCircleCIEnvironmentVariableCheckDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circleci_environment_variable" {
			continue
		}

		has, err := c.HasProjectEnvironmentVariable(rs.Primary.Attributes["project"], rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}

		if has {
			return errors.New("Environment variable should have been destroyed")
		}
	}

	return nil
}

func testAccCircleCIEnvironmentVariableConfig(project, name, value string) string {
	return fmt.Sprintf(`
resource "circleci_environment_variable" "%[2]s" {
  project = "%[1]s"
  name    = "%[2]s"
  value   = "%[3]s"
}`, project, name, value)
}

func testAccCircleCIEnvironmentVariableConfigIdentical(project, name, value string) string {
	return fmt.Sprintf(`
resource "circleci_environment_variable" "%[2]s" {
  project = "%[1]s"
  name    = "%[2]s"
  value   = "%[3]s"
}

resource "circleci_environment_variable" "%[2]s_2" {
  project = "%[1]s"
  name    = "%[2]s"
  value   = "%[3]s"
}`, project, name, value)
}
