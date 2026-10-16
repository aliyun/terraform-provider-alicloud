# Terraform Provider for Alibaba Cloud

The Alibaba Cloud Terraform Provider allows managing resources and data sources within Alibaba Cloud.

The 2.x line is the actively maintained line; the 1.x line is end-of-life. Terraform CLI 1.14 or later is recommended.

* [Alibaba Cloud Provider Documentation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs)
* [Tutorials](https://developer.hashicorp.com/terraform/tutorials)
* [Usage Examples](https://github.com/aliyun/terraform-provider-alicloud/tree/master/examples)
* [Community](https://developer.hashicorp.com/terraform/community)

## Usage Example

```hcl
# 1. Specify the version of the Alibaba Cloud Provider to use
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 2.0"
    }
  }
}

# 2. Configure the Alibaba Cloud Provider
provider "alicloud" {
  region = "cn-hangzhou"
}

# 3. Create a VPC
resource "alicloud_vpc" "example" {
  vpc_name   = "example-vpc"
  cidr_block = "10.0.0.0/16"
}
```

* [Usage documentation for the Alibaba Cloud Provider can be found in the Terraform Registry](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](https://go.dev/doc/install) (version 1.25.8+ is *required*) installed on your machine.

```sh
$ git clone https://github.com/aliyun/terraform-provider-alicloud
$ cd terraform-provider-alicloud
$ make tools
$ make build
```

`make build` compiles the provider for macOS, Linux and Windows and places the archives under `bin/`. To build and install a development provider that replaces the one used by your local Terraform CLI, run `make dev` (macOS), `make devlinux` (Linux) or `make devwin` (Windows).

Before you open a pull request, confirm that your commit identity is correct. Git records both an author and a committer for every commit; inspect them with `git var GIT_AUTHOR_IDENT` and `git var GIT_COMMITTER_IDENT`. Use the email address associated with your GitHub account so that the contributor license agreement (CLA) check running on the pull request can identify you. If the recorded identity is wrong, set it for this repository before you commit rather than rewriting existing commits.

In order to test the provider, you can simply run `make test`. In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

Bug reports and feature requests are tracked in the [issue tracker](https://github.com/aliyun/terraform-provider-alicloud/issues).
