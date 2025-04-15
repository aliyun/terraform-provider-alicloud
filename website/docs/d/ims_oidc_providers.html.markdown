---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_oidc_providers"
sidebar_current: "docs-alicloud-datasource-ims-oidc-providers"
description: |-
  Provides a list of Ims Oidc Provider owned by an Alibaba Cloud account.
---

# alicloud_ims_oidc_providers

This data source provides Ims Oidc Provider available to the user.[What is Oidc Provider](https://next.api.alibabacloud.com/document/Ims/2019-08-15/CreateOIDCProvider)

-> **NOTE:** Available since v1.248.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "oidc_provider_name" {
  default = "amp-resource-example-oidc-provider"
}

resource "alicloud_ims_oidc_provider" "default" {
  description = var.oidc_provider_name
  issuer_url  = "https://oauth.aliyun.com"
  fingerprints = [
    "0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71"
  ]
  issuance_limit_time = "12"
  oidc_provider_name  = var.oidc_provider_name
  client_ids = [
    "123",
    "456"
  ]
}

data "alicloud_ims_oidc_providers" "default" {
  ids = ["${alicloud_ims_oidc_provider.default.id}"]
}

output "alicloud_ims_oidc_provider_example_id" {
  value = data.alicloud_ims_oidc_providers.default.providers.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Oidc Provider IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Oidc Provider IDs.
* `names` - A list of name of Oidc Providers.
* `providers` - A list of Oidc Provider Entries. Each element contains the following attributes:
  * `arn` - ARN of OIDC identity provider.
  * `client_ids` - Client ID. 
  * `create_time` - Creation Time (UTC time).
  * `description` - Description of OIDC identity provider.
  * `fingerprints` - The authentication fingerprint of the HTTPS CA certificate.
  * `issuance_limit_time` - The earliest time when an external IdP is allowed to issue an ID Token. If the iat field in the ID Token is greater than the current time, the request is rejected.Unit: hours. Value range: 1~168.
  * `issuer_url` - The issuer URL of the OIDC identity provider.
  * `oidc_provider_name` - The name of the OIDC identity provider.
  * `update_time` - Modification Time (UTC time).
  * `id` - The ID of the resource supplied above.
