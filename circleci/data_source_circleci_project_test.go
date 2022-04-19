package circleci

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCircleCIProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCircleCIProjectDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_project.foo", "name", "terraform-test"),
				),
			},
		},
	})
}

const testAccCircleCIProjectDataSource = `
data "circleci_project" "foo" {
  name = "terraform-test"
}
`
