---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "circleci_checkout_key Resource - terraform-provider-circleci"
subcategory: ""
description: |-
  
---

# circleci_checkout_key (Resource)

## Usage

Creating a `deploy-key`:
```hcl
resource "circleci_checkout_key" "my_deploy_key" {
  project = "my_project"
  type    = "deploy-key"
}
```

Creating a `user-key`:
```hcl
resource "circleci_checkout_key" "my_user_key" {
  project = data.circleci_project.my_project.name
  type    = "user-key"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project` (String) The name of the CircleCI project to create the checkout key in.
- `type` (String) The type of the checkout key. Can be either `user-key` or `deploy-key`.

### Read-Only

- `created_at` (String) The date and time the checkout key was created.
- `fingerprint` (String) The fingerprint of the checkout key.
- `id` (String) The ID of this resource.
- `preferred` (Boolean) A boolean value that indicates if this key is preferred.
- `public_key` (String) The public SSH key of the checkout key.

~> The `preferred` flag is automatically set to true on the most recent key created.

~> For `deploy-key` type, the resource will also create a deploy key in your VCS repository, which will not be deleted in case of Terraform destroy. Requires manual clean up.

## Import

Checkout Keys can be imported using the project name and the fingerprint.
```bash
$ terraform import circleci_checkout_key.key "my-project/12:34:56:78:90:12:34:56:78:90:12:34:56:78:90:12"
```