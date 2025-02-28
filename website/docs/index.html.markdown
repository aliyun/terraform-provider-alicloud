---
layout: "alicloud"
page_title: "Provider: alicloud"
sidebar_current: "docs-alicloud-index"
description: |-
  The Alicloud provider is used to interact with many resources supported by Alibaba Cloud. The provider needs to be configured with the proper credentials before it can be used.
---

# Alibaba Cloud Provider

The Alibaba Cloud provider is used to interact with the
many resources supported by [Alibaba Cloud](https://www.alibabacloud.com). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.50.0, the provider start to support Terraform 0.12.x.


## Example Usage

```terraform
# Configure the AliCloud Provider

provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  # If not set, cn-beijing will be used.
  region = var.region
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

# Create a new ECS instance for VPC
resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

# Create a new Security in a VPC
resource "alicloud_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.vpc.id
}
# Create a kms to encrypt the disk
resource "alicloud_kms_key" "key" {
  description            = "Hello KMS"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_instance" "instance" {
  # cn-beijing
  availability_zone = data.alicloud_zones.default.zones.0.id
  security_groups   = alicloud_security_group.group.*.id

  # series III
  instance_type              = "ecs.n4.large"
  system_disk_category       = "cloud_efficiency"
  system_disk_name           = var.name
  system_disk_description    = "system_disk_description"
  image_id                   = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_name              = var.name
  vswitch_id                 = alicloud_vswitch.vswitch.id
  internet_max_bandwidth_out = 10
  data_disks {
    name        = "data-disk"
    size        = 20
    category    = "cloud_efficiency"
    description = "disk-description"
    encrypted   = true
    kms_key_id  = alicloud_kms_key.key.id
  }
}
```

## Authentication

The Alicloud provider accepts several ways to enter credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables
- Shared credentials/configuration file  
- ECS Instance Role
- Assuming A RAM Role
- Assuming A RAM Role With OIDC
- Sidecar Credentials

### Static credentials

Static credentials can be provided by adding `access_key`, `secret_key` and `region` in-line in the
alicloud provider block:

Usage:

```terraform
provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}
```

### Environment variables

You can provide your credentials via `ALIBABA_CLOUD_ACCESS_KEY_ID`, `ALIBABA_CLOUD_ACCESS_KEY_SECRET` and optionally
`ALIBABA_CLOUD_SECURITY_TOKEN` environment variables. The Region can be set using the `ALIBABA_CLOUD_REGION` environment variables.

Usage:
```terraform
provider "alicloud" {}
```
```shell
$ export ALIBABA_CLOUD_ACCESS_KEY_ID="<Your-Access-Key-ID>"
$ export ALIBABA_CLOUD_ACCESS_KEY_SECRET="<Your-Access-Key-Secret>"
$ export ALIBABA_CLOUD_REGION="cn-beijing"
$ terraform plan
```

### Shared Credentials File

You can use an [Alibaba Cloud credentials or configuration file](https://www.alibabacloud.com/help/doc-detail/110341.htm) to specify your credentials. 
The default location is `$HOME/.aliyun/config.json` on Linux and macOS, or `"%USERPROFILE%\.aliyun/config.json"` on Windows. 
You can optionally specify a different location in the Terraform configuration by providing the `shared_credentials_file` argument 
or using the `ALIBABA_CLOUD_CREDENTIALS_FILE` environment variable. 
This method also supports a `profile` configuration and matching `ALIBABA_CLOUD_PROFILE` environment variable:

Usage:

```terraform
provider "alicloud" {
  region                  = "cn-hangzhou"
  shared_credentials_file = "/Users/tf_user/.aliyun/creds"
  profile                 = "customprofile"
}
```

### ECS Instance Role

If you're running Terraform from an ECS instance with RAM Instance using RAM Role,
Terraform will just access
the metadata URL: `http://100.100.100.200/latest/meta-data/ram/security-credentials/<ecs_role_name>`
to obtain the STS credential.
Refer to details [Access other Cloud Product APIs by the Instance RAM Role](https://www.alibabacloud.com/help/doc-detail/54579.htm).

This is a preferred approach over any other when running in ECS as you can avoid
hard coding credentials. Instead, these are leased on-the-fly by Terraform
which reduces the chance of leakage.

The ECS Instance Role can be set using the `ALIBABA_CLOUD_ECS_METADATA` environment variables.

Usage:

```terraform
provider "alicloud" {
  ecs_role_name = "terraform-provider-alicloud"
  region        = var.region
}
```

### Assuming A RAM Role

If provided with a role ARN, Terraform will attempt to assume this role using the supplied credentials. 
The role arn can be set using the `ALIBABA_CLOUD_ROLE_ARN` environment variables, 
and the role session name using the `ALIBABA_CLOUD_ROLE_SESSION_NAME` environment variables.

Usage:

```terraform
provider "alicloud" {
  access_key = "<One-AccessKeyId-With-AssumeRole-Policy>"
  secret_key = "<One-AccessKeySecret-With-AssumeRole-Policy>"
  assume_role {
    role_arn           = "acs:ram::ACCOUNT_ID:role/ROLE_NAME"
    policy             = "Policy Content"
    session_name       = "A Role Session Name"
    session_expiration = 999
  }
}
```

### Assuming A RAM Role With OIDC

If provided with a role ARN and a token from a service account OpenID Connect (OIDC),
the Alibaba CLoud Provider will attempt to assume this role using the supplied credentials.

**NOTE:** Assuming-Role-With-OIDC is a no-AK auth type, and there is no need setting access_key and secret_key while using it.

Usage:

```terraform
provider "alicloud" {
  assume_role_with_oidc {
    oidc_provider_arn = "acs:ram::ACCOUNT_ID:oidc-provider/ROLE_NAME"
    role_arn          = "acs:ram::ACCOUNT_ID:role/ROLE_NAME"
    oidc_token_file   = "/Users/tf_user/secrets/rrsa-tokens/token"
    role_session_name = "A Role Session Name"
  }
}
```

### Sidecar Credentials

You can deploy a sidecar to storage alibaba cloud credentials. 
Then, you can optionally specify a credentials URI in the Terraform configuration by providing the `credentials_uri` argument 
or using the `ALIBABA_CLOUD_CREDENTIALS_URI` environment variable to get the credentials automatically. 
The Sidecar Credentials is available since v1.141.0.

Usage:

```terraform
provider "alicloud" {
  region          = "cn-hangzhou"
  credentials_uri = "<Your-Credential-URI>"
}
```

### Custom User-Agent Information

By default, the underlying AlibabaCloud client used by the Terraform AliCloud Provider creates requests with User-Agent headers including information about Terraform and AlibabaCloud Go SDK versions. 
To provide additional information in the User-Agent headers, the provider variable `configuration_source` or `TF_APPEND_USER_AGENT` environment variable can be set and its value will be directly added to HTTP requests.

Usage:

```terraform
provider "alicloud" {
  region               = "cn-hangzhou"
  configuration_source = "ArgoAgent/argo-12345678 NodeID/1234"
}
```

or

```shell
$ export TF_APPEND_USER_AGENT="ArgoAgent/argo-12345678 NodeID/1234 (Optional Extra Information)"
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Alibaba Cloud
 `provider` block:

* `access_key` - Alibaba Cloud access key. It is required for the provider. 
  Can also be set with the `ALIBABA_CLOUD_ACCESS_KEY_ID` environment variable since v1.228.0, 
  or via a shared credentials file if profile is specified. See also `secret_key`. 
  Environment variable `ALICLOUD_ACCESS_KEY` and `ALIBABACLOUD_ACCESS_KEY_ID` have been deprecated since v1.228.0.

* `secret_key` - Alibaba Cloud secret key. It is required for the provider. 
  Can also be set with the `ALIBABA_CLOUD_ACCESS_KEY_SECRET` environment variable since v1.228.0,
  or via a shared credentials file if profile is specified. See also `access_key`.
  Environment variable `ALICLOUD_SECRET_KEY` and `ALIBABACLOUD_ACCESS_KEY_SECRET` have been deprecated since v1.228.0.

* `security_token` - Alibaba Cloud [Security Token Service](https://www.alibabacloud.com/help/en/ram/product-overview/what-is-sts).
  Can also be set with the `ALIBABA_CLOUD_SECURITY_TOKEN` environment variable since v1.228.0,
  or via a shared credentials file if profile is specified. See also `access_key`.
  Environment variable `ALICLOUD_SECURITY_TOKEN` and `ALIBABACLOUD_SECURITY_TOKEN` have been deprecated since v1.228.0.

* `ecs_role_name` - The RAM Role Name attached on a ECS instance for API operations.
  Can also be set with the `ALIBABA_CLOUD_ECS_METADATA` environment variable since v1.228.0.
  Environment variable `ALICLOUD_ECS_ROLE_NAME` has been deprecated since v1.228.0.

* `region` - Alibaba Cloud region. Default to `cn-beijing`. 
  Can also be set with the `ALIBABA_CLOUD_REGION` environment variable since v1.228.0.
  Environment variable `ALICLOUD_REGION` has been deprecated since v1.228.0.

* `account_id` - (Optional) Alibaba Cloud Account ID. It is used by the Function Compute service and to connect router interfaces.
  If not provided, the provider will attempt to retrieve it automatically with [STS GetCallerIdentity](https://www.alibabacloud.com/help/doc-detail/43767.htm).
  Can also be set with the `ALIBABA_CLOUD_ACCOUNT_ID` environment variable since v1.228.0.
  Environment variable `ALICLOUD_ACCOUNT_ID` has been deprecated since v1.228.0.

* `account_type` - (Optional, Available since 1.240.0) Alibaba Cloud [Account Type](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/getting-account). 
  It used to indicate caller identity's account type. Can also be set with the `ALIBABA_CLOUD_ACCOUNT_TYPE` environment variable. Valid values:
  - `Domestic`(Default): China-Site Account.
  - `International`: International-Site Account.
  
* `shared_credentials_file` - (Optional, Available since 1.49.0) This is the path to the shared credentials file.
  Can also be set with the `ALIBABA_CLOUD_CREDENTIALS_FILE` environment variable since v1.228.0.
  Environment variable `ALICLOUD_SHARED_CREDENTIALS_FILE` has been deprecated since v1.228.0.
  If this is not set and `profile` is specified, "~/.aliyun/config.json" will be used.

* `profile` - (Optional, Available since 1.49.0) This is the Alibaba Cloud profile name as set in the shared credentials file.
  Can also be set with the `ALIBABA_CLOUD_PROFILE` environment variable since v1.228.0.
  Environment variable `ALICLOUD_PROFILE` has been deprecated since v1.228.0.

* `assume_role` - (Optional) An [`assume_role` Configuration Block](#assume_role-configuration-block) block. Only one `assume_role` block may be in the configuration.

* `assume_role_with_oidc` - (Optional, Available since v1.220.0) Configuration block for assuming an RAM role using an OIDC. See the [`assume_role_with_oidc` Configuration Block](#assume_role_with_oidc-configuration-block) section below. Only one `assume_role_with_oidc` block may be in the configuration.

* `credentials_uri` - (Optional, Available since 1.141.0) The URI of sidecar credentials service. 
  Can also be set with the `ALIBABA_CLOUD_CREDENTIALS_URI` environment variable since v1.228.0.
  Environment variable `ALICLOUD_CREDENTIALS_URI` has been deprecated since v1.228.0.

* `endpoints` - (Optional) An [`endpoints`](#endpoints) block to support custom endpoints.

* `skip_region_validation` - (Optional, Available since 1.52.0) Skip static validation of region ID. Used by users of alternative AlibabaCloud-like APIs or users w/ access to regions that are not public (yet).

* `configuration_source` - (Optional, Available since 1.56.0) Use a string to mark a configuration file source, like `terraform-alicloud-modules/terraform-alicloud-ecs-instance` or `terraform-provider-alicloud/examples/vpc`.
The length should not more than 128(Before 1.207.2, it should not more than 64). Since the version 1.145.0, it supports to be set by environment variable `TF_APPEND_USER_AGENT`. See `Custom User-Agent Information`.

* `protocol` - (Optional, Available since 1.72.0) The Protocol of used by API request. Valid values: `HTTP` and `HTTPS`. Default to `HTTPS`. 

* `client_read_timeout` - (Optional, Available since 1.125.0) The maximum timeout in millisecond second of the client read request. Default to 60000.

* `client_connect_timeout` - (Optional, Available since 1.125.0) The maximum timeout in millisecond second of the client connection server. Default to 60000.

* `max_retry_timeout` - (Optional, Available since 1.183.0) The maximum retry timeout in second of the request. Default to `0`.

### `assume_role` Configuration Block

* `role_arn` - (Required) The ARN of the role to assume. If ARN is set to an empty string, it does not perform role switching. 
  Can also be set with the `ALIBABA_CLOUD_ROLE_ARN` environment variable since v1.228.0.
  Environment variable `ALICLOUD_ASSUME_ROLE_ARN` has been deprecated since v1.228.0.
  Terraform executes configuration on account with provided credentials.

* `policy` - (Optional) A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict the permissions for the resulting temporary
  security credentials. You cannot use the passed policy to grant permissions that are in excess of those allowed by the access policy of the role that is being assumed.

* `session_name` - (Optional) The session name to use when assuming the role. If omitted, 'terraform' is passed to the AssumeRole call as session name. 
  Can also be set with the `ALIBABA_CLOUD_ROLE_SESSION_NAME` environment variable since v1.228.0.
  Environment variable `ALICLOUD_ASSUME_ROLE_SESSION_NAME` has been deprecated since v1.228.0.

* `session_expiration` - (Optional) The time after which the established session for assuming role expires. Valid value range: [900-43200] seconds. Default to 3600 (in this case Alicloud use own default value). It supports environment variable `ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION`.

* `external_id` - (Optional, Available since 1.207.1) The external ID of the RAM role. 
  This parameter is provided by an external party and is used to prevent the confused deputy problem. 
  The value must be 2 to 1,224 characters in length and can contain letters, digits, and the following special characters:`= , . @ : / - _`.

### assume_role_with_oidc Configuration Block

The `assume_role_with_oidc` configuration block supports the following arguments:

* `oidc_provider_arn` - (Required) ARN of the OIDC IdP. Can also be set with the `ALIBABA_CLOUD_OIDC_PROVIDER_ARN` environment variable.
* `role_arn` - (Required) ARN of the RAM Role to assume. Can also be set with the `ALIBABA_CLOUD_ROLE_ARN` environment variable.
* `oidc_token` - (Optional) Value of a RRSA security token from an OIDC Idp. One of `oidc_token` or `oidc_token_file` is required.
  Can also be set with the `ALIBABA_CLOUD_OIDC_TOKEN` environment variable.
* `oidc_token_file` - (Optional) File containing a RRSA security token from an OIDC. One of `oidc_token_file` or `oidc_token` is required.
  Can also be set with the `ALIBABA_CLOUD_OIDC_TOKEN_FILE` environment variable.
* `role_session_name` - (Optional) The session name to use when assuming the role. If omitted, 'terraform' is passed to the AssumeRoleWithOIDC call as session name. 
  Can also be set with the `ALIBABA_CLOUD_ROLE_SESSION_NAME` environment variable.
* `session_expiration` - (Optional) The validity period of the STS token. Unit: seconds. Default value: 3600. Minimum value: 900. Maximum value: the value of the MaxSessionDuration parameter when creating a ram role.
* `policy` - (Optional) The policy that specifies the permissions of the returned STS token. You can use this parameter to grant the STS token fewer permissions than the permissions granted to the RAM role.
 
### `endpoints`

**NOTE:** Due to certain API restrictions, the endpoints pointing to the area should be consistent with the `region_id`.

* `ecs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.

* `slb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.

* `cbn` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CEN endpoints.

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

* `market` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Market endpoints.

* `cddc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ApsaraDB for MyBase endpoints.

* `ehpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Elastic High Performance Computing endpoints.

* `mscsub` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Message Center endpoints.

* `hitsdb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Lindorm endpoints.

* `sddp` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Data Security Center endpoints.

* `sas` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Security Center endpoints.

* `dts` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Data Transmission endpoints.

* `ens` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ens endpoints.

* `alidfs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Apsara File Storage for HDFS endpoints.

* `arms` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Application Real-Time Monitoring Service endpoints.

* `bastionhost` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Bastion Host endpoints.

* `waf` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Web Application Firewall endpoints.

* `alb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Application Load Balancer endpoints.

* `hbr` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Hybrid Backup Recovery endpoints.

* `dataworkspublic` - - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Data Works endpoints.

* `cloudfw` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Cloud Firewall endpoints.

* `dm` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Direct Mail endpoints.

* `eais` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Elastic Accelerated Computing Instances endpoints.

* `dg` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Database Gateway endpoints.

* `imm` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Intelligent Media Management endpoints.

* `iot` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Internet of Things endpoints.

* `vod` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VOD endpoints.

* `gds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Graph Database endpoints.

* `swas` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Simple Application Server endpoints.

* `opensearch` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Open Search endpoints.

* `clickhouse` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Click House endpoints.

* `vs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Video Surveillance System endpoints.

* `quickbi` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Quick BI endpoints.

* `cloudsso` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Cloud SSO endpoints.

* `edas` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom EDAS endpoints.
* `dmsenterprise` - - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DMS Enterprise endpoints.

## Testing

Credentials must be provided via the `ALIBABA_CLOUD_ACCESS_KEY_ID`, `ALIBABA_CLOUD_ACCESS_KEY_SECRET` and `ALIBABA_CLOUD_REGION` environment variables in order to run acceptance tests.
