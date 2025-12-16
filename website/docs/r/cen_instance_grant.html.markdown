---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_grant"
sidebar_current: "docs-alicloud-resource-cen-instance-grant"
description: |-
  Provides a Alicloud CEN child instance grant resource.
---

# alicloud_cen_instance_grant

Provides a CEN child instance grant resource, which allow you to authorize a VPC or VBR to a CEN of a different account.

For more information about how to use it, see [Attach a network in a different account](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-attachcenchildinstance). 

-> **NOTE:** Deprecated since v1.241.0. The resource have been deprecated and new resource type [alicloud_cen_transit_router_grant_attachment](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_grant_attachment) is recommended.

-> **NOTE:** Available since v1.37.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_instance_grant&exampleId=53269692-a358-5a74-e55d-be8b1b6aefcfea760c20&activeTab=example&spm=docs.r.cen_instance_grant.0.53269692a3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "another_uid" {
  default = "xxxx"
}
# Method 1: Use assume_role to operate resources in the target cen account, detail see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#assume-role
provider "alicloud" {
  region = "cn-hangzhou"
  alias  = "child_account"
  assume_role {
    role_arn = "acs:ram::${var.another_uid}:role/terraform-example-assume-role"
  }
}

# Method 2: Use the target cen account's access_key, secret_key
# provider "alicloud" {
#   region     = "cn-hangzhou"
#   access_key = "access_key"
#   secret_key = "secret_key"
#   alias      = "child_account"
# }

provider "alicloud" {
  alias = "your_account"
}
data "alicloud_account" "your_account" {
  provider = alicloud.your_account
}
data "alicloud_account" "child_account" {
  provider = alicloud.child_account
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_cen_instance" "example" {
  provider          = alicloud.your_account
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_vpc" "child_account" {
  provider   = alicloud.child_account
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}
resource "alicloud_cen_instance_grant" "child_account" {
  provider          = alicloud.child_account
  cen_id            = alicloud_cen_instance.example.id
  child_instance_id = alicloud_vpc.child_account.id
  cen_owner_id      = data.alicloud_account.your_account.id
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_instance_grant&spm=docs.r.cen_instance_grant.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `child_instance_id` - (Required, ForceNew) The ID of the child instance to grant.
* `cen_owner_id` - (Required, ForceNew) The owner UID of the  CEN which the child instance granted to.

## Attributes Reference

The following attributes are exported:

- `id` - ID of the resource, formatted as `<cen_id>:<child_instance_id>:<cen_owner_id>`.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance_grant.example cen-abc123456:vpc-abc123456:uid123456
```
