# CircleCI Terraform Provider

Terraform provider plugin to manage CircleCI resources using CircleCI API V2.

### Maintainers

This provider plugin is maintained by the DevOps team at [SectorLabs](https://www.sectorlabs.ro).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.18+ (to build the provider plugin)

## Using the provider

### Building the provider

Clone repository to: `$GOPATH/src/github.com/SectorLabs/terraform-provider-circleci`

```
$ mkdir -p $GOPATH/src/github.com/SectorLabs; cd $GOPATH/src/github.com/SectorLabs
$ git clone git@github.com:SectorLabs/terraform-provider-circleci.git
```

Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/SectorLabs/terraform-provider-circleci
$ go build -v -o terraform-provider-circleci main.go
```

#### Testing the provider

Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/SectorLabs/terraform-provider-circleci
$ go test -v ./...
```


### Downloading the provider

Download the latest release to your OS from the [release page](https://github.com/SectorLabs/terraform-provider-circleci/releases).

### Installing the provider

Move the provider binary (either built or downloaded) to your local Terraform directory

```
$ mkdir -p ~/.terraform.d/plugins/registry.terraform.io/SectorLabs/circleci/1.0.0/$(go env GOOS)_$(go env GOARCH)
$ mv terraform-provider-circleci ~/.terraform.d/registry.terraform.io/SectorLabs/circleci/1.0.0/$(go env GOOS)_$(go env GOARCH)/terraform-provider-circleci
```

Replace the `1.0.0` directory name with the version you're currently using.

## Resources

```hcl
terraform {
  required_providers {
    circleci = {
      source  = "SectorLabs/circleci"
      version = "1.0.0"
    }
  }
}

provider "circleci" {
  api_token    = "${file("circleci_token")}"
  vcs_type     = "github"
  organization = "MyOrganization"
}

data "circleci_project" "project" {
  name = "myProject"
}

resource "circleci_environment_variable" "variable" {
  project = data.circleci_project.project.name
  name    = "DUMMY"
  value   = "VALUE"
}

resource "circleci_checkout_key" "key" {
  project = data.circleci_project.project.name
  type    = "deploy-key"
}

resource "circleci_context" "context" {
  name = "myContext"
}

resource "circleci_context_environment_variable" "variable" {
  context = circleci_context.context.name

  name  = "DUMMY"
  value = "VALUE"
}
```