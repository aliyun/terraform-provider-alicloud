---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network_grant"
sidebar_current: "docs-alicloud-resource-cloud-connect-network-grant"
description: |-
  Provides a Alicloud Cloud Connect Network Grant resource.
---

# alicloud_cloud_connect_network_grant

Provides a Cloud Connect Network Grant resource. If the CEN instance to be attached belongs to another account, authorization by the CEN instance is required.

For information about Cloud Connect Network Grant and how to use it, see [What is Cloud Connect Network Grant](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/grantinstancetocbn).

-> **NOTE:** Available since v1.63.0.

-> **NOTE:** Only the following regions support create Cloud Connect Network Grant. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_connect_network_grant&exampleId=931a835b-f9e6-c6b1-0acd-c83a1f1c0a193a3d3a9a&activeTab=example&spm=docs.r.cloud_connect_network_grant.0.931a835bf9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "another_uid" {
  default = 123456789
}

provider "alicloud" {
  region = "cn-shanghai"
  alias  = "default"
}

# Method 1: Use assume_role to operate resources in the target cen account, detail see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#assume-role
provider "alicloud" {
  region = "cn-hangzhou"
  alias  = "cen_account"
  assume_role {
    role_arn = "acs:ram::${var.another_uid}:role/terraform-example-assume-role"
  }
}


# Method 2: Use the target cen account's access_key, secret_key
# provider "alicloud" {
#   region     = "cn-hangzhou"
#   access_key = "access_key"
#   secret_key = "secret_key"
#   alias      = "cen_account"
# }

resource "alicloud_cloud_connect_network" "default" {
  provider    = alicloud.default
  name        = var.name
  description = var.name
  cidr_block  = "192.168.0.0/24"
  is_default  = true
}

resource "alicloud_cen_instance" "cen" {
  provider          = alicloud.cen_account
  cen_instance_name = var.name
}

resource "alicloud_cloud_connect_network_grant" "default" {
  provider = alicloud.default
  ccn_id   = alicloud_cloud_connect_network.default.id
  cen_id   = alicloud_cen_instance.cen.id
  cen_uid  = var.another_uid
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_connect_network_grant&spm=docs.r.cloud_connect_network_grant.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, ForceNew) The ID of the CCN instance.
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `cen_uid` - (Required, ForceNew) The ID of the account to which the CEN instance belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cloud Connect Network grant Id and formates as `<ccn_id>:<cen_id>`.

## Import

The Cloud Connect Network Grant can be imported using the instance_id, e.g.

```shell
$ terraform import alicloud_cloud_connect_network_grant.example ccn-abc123456:cen-abc123456
```

