---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "circleci_context Data Source - terraform-provider-circleci"
subcategory: ""
description: |-
  
---

# circleci_context (Data Source)

## Usage
```hcl
data "circleci_context" "my_context" {
  name = "my-context"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the context

### Read-Only

- `id` (String) The ID of this resource.

