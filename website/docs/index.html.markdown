---
layout: "alicloud"
page_title: "Provider: alicloud"
sidebar_current: "docs-alicloud-index"
description: |-
  The Alicloud provider is used to interact with many resources supported by Alicloud. The provider needs to be configured with the proper credentials before it can be used.
---

# Alibaba Cloud Provider

The Alibaba Cloud provider is used to interact with the
many resources supported by [Alibaba Cloud](https://www.alibabacloud.com). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.50.0, the provider start to support Terraform 0.12.x.

-> **Note:** When you use terraform on a `Windows` computer, please install [golang](https://golang.org/dl/) first.
Otherwise, you may encounter an issue that occurs from the version 1.8.1 to 1.10.0. For more information, please read the [Crash Error](https://github.com/alibaba/terraform-provider/issues/469).


## Example Usage

```hcl
# Configure the Alicloud Provider
provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "alicloud_instance_types" "c2g4" {
  cpu_core_count = 2
  memory_size = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}

# Create a web server
resource "alicloud_instance" "web" {
  image_id          = "${data.alicloud_images.default.images.0.id}"
  internet_charge_type  = "PayByBandwidth"

  instance_type        = "${data.alicloud_instance_types.c2g4.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${alicloud_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id = "vsw-abc12345"
}

# Create security group
resource "alicloud_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id = "vpc-abc12345"
}
```

## Authentication

The Alicloud provider accepts several ways to enter credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables
- ECS Role
- Assume role

### Static credentials

Static credentials can be provided by adding `access_key`, `secret_key` and `region` in-line in the
alicloud provider block:

Usage:

```hcl
provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```


### Environment variables

You can provide your credentials via `ALICLOUD_ACCESS_KEY` and `ALICLOUD_SECRET_KEY`
environment variables, representing your Alicloud access key and secret key respectively.
`ALICLOUD_REGION` is also used, if applicable:

```hcl
provider "alicloud" {}
```

Usage:

```shell
$ export ALICLOUD_ACCESS_KEY="anaccesskey"
$ export ALICLOUD_SECRET_KEY="asecretkey"
$ export ALICLOUD_REGION="cn-beijing"
$ terraform plan
```

### ECS Role

If you're running Terraform from an ECS instance with RAM Instance using RAM Role,
Terraform will just access
the metadata URL: `http://100.100.100.200/latest/meta-data/ram/security-credentials/<ecs_role_name>`
to obtain the STS credential.
Refer to details [Access other Cloud Product APIs by the Instance RAM Role](https://www.alibabacloud.com/help/doc-detail/54579.htm).

This is a preferred approach over any other when running in ECS as you can avoid
hard coding credentials. Instead these are leased on-the-fly by Terraform
which reduces the chance of leakage.


Usage:

```hcl
provider "alicloud" {
  ecs_role_name = "terraform-provider-alicloud"
  region        = "${var.region}"
}
```

-> **NOTE:** At present, the [MNS Resources](https://www.terraform.io/docs/providers/alicloud/r/mns_queue.html) does not support ECS Role Credential.

### Assume role

If provided with a role ARN, Terraform will attempt to assume this role using the supplied credentials.

Usage:

```hcl
provider "alicloud" {
  assume_role {
    role_arn           = "acs:ram::ACCOUNT_ID:role/ROLE_NAME"
    policy             = "POLICY"
    session_name       = "SESSION_NAME"
    session_expiration = 999
  }
}
```


## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Alibaba Cloud
 `provider` block:

* `access_key` - This is the Alicloud access key. It must be provided, but
  it can also be sourced from the `ALICLOUD_ACCESS_KEY` environment variable, or via
  a dynamic access key if `ecs_role_name` is specified.

* `secret_key` - This is the Alicloud secret key. It must be provided, but
  it can also be sourced from the `ALICLOUD_SECRET_KEY` environment variable, or via
  a dynamic secret key if `ecs_role_name` is specified.

* `security_token` - Alicloud [Security Token Service](https://www.alibabacloud.com/help/doc-detail/66222.html).
  It can be sourced from the `ALICLOUD_SECURITY_TOKEN` environment variable,  or via
  a dynamic security token if `ecs_role_name` is specified.

* `ecs_role_name` - "The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the 'Access Control' section of the Alibaba Cloud console.",

* `region` - This is the Alicloud region. It must be provided, but
  it can also be sourced from the `ALICLOUD_REGION` environment variables.

* `account_id` - (Optional) Alibaba Cloud Account ID. It is used by the Function Compute service and to connect router interfaces.
  If not provided, the provider will attempt to retrieve it automatically with [STS GetCallerIdentity](https://www.alibabacloud.com/help/doc-detail/43767.htm).
  It can be sourced from the `ALICLOUD_ACCOUNT_ID` environment variable.

* `shared_credentials_file` - (Optional, Available in 1.49.0+) This is the path to the shared credentials file. If this is not set and a profile is specified, ~/.aliyun/config.json will be used.

* `profile` - (Optional, Available in 1.49.0+) This is the Alicloud profile name as set in the shared credentials file. It can also be sourced from the `ALICLOUD_PROFILE` environment variable.

* `assume_role` - (Optional) An `assume_role` block (documented below). Only one `assume_role` block may be in the configuration.

* `endpoints` - (Optional) An `endpoints` block (documented below) to support custom endpoints.

* `skip_region_validation` - (Optional, Available in 1.52.0+) Skip static validation of region ID. Used by users of alternative AlibabaCloud-like APIs or users w/ access to regions that are not public (yet).

* `source_name` - (Optional, Available in 1.56.0+) Use a name to mark a template. It can be a source or a specifial usage scenario, like `terraform-alicloud-modules/ram/alicloud` or `examples/vpc`.

The nested `assume_role` block supports the following:

* `role_arn` - (Required) The ARN of the role to assume. If ARN is set to an empty string, it does not perform role switching. It supports environment variable `ALICLOUD_ASSUME_ROLE_ARN`.
  Terraform executes configuration on account with provided credentials.

* `policy` - (Optional) A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict the permissions for the resulting temporary
  security credentials. You cannot use the passed policy to grant permissions that are in excess of those allowed by the access policy of the role that is being assumed.

* `session_name` - (Optional) The session name to use when assuming the role. If omitted, 'terraform' is passed to the AssumeRole call as session name. It supports environment variable `ALICLOUD_ASSUME_ROLE_SESSION_NAME`.

* `session_expiration` - (Optional) The time after which the established session for assuming role expires. Valid value range: [900-3600] seconds. Default to 3600 (in this case Alicloud use own default value). It supports environment variable `ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION`.

Nested `endpoints` block supports the following:

* `ecs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.

* `slb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.

* `cen` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CEN endpoints.

* `ess` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.

* `oss` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.

* `dns` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DNS endpoints.

* `ram` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RAM endpoints.

* `cs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Service endpoints.

* `cr` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Registry endpoints.

* `cdn` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CDN endpoints.

* `kms` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom KMS endpoints.

* `ots` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Table Store endpoints.

* `cms` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Cloud Monitor endpoints.

* `pvtz` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Private Zone endpoints.

* `sts` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom STS endpoints.

* `log` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Log Service endpoints.

* `drds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DRDS endpoints.

* `dds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MongoDB endpoints.

* `gpdb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom GPDB endpoints.

* `kvstore` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom R-KVStore endpoints.

* `fc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Function Computing endpoints.

* `apigateway` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Api Gateway endpoints.

* `datahub` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Datahub endpoints.

* `mns` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MNS endpoints.

* `location` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Location Service endpoints.",

* `nas` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom nas Service endpoints.",

* `actiontrail` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom actiontrail Service endpoints.",

* `cas` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CAS endpoints.

* `bssopenapi` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BSSOPENAPI endpoints.

* `ddoscoo` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BGP-Line Anti-DDoS Pro endpoints.

## Testing

Credentials must be provided via the `ALICLOUD_ACCESS_KEY`, `ALICLOUD_SECRET_KEY` and `ALICLOUD_REGION` environment variables in order to run acceptance tests.
