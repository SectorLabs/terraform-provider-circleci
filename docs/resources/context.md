---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "circleci_context Resource - terraform-provider-circleci"
subcategory: ""
description: |-
  
---

# circleci_context (Resource)

## Usage
```hcl
resource "circleci_context" "my_context" {
  name = "my-awesome-context"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the context

### Read-Only

- `id` (String) The ID of this resource.

## Import

Contexts can be imported using their names:
```bash
$ terraform import circleci_context.context my-context
```