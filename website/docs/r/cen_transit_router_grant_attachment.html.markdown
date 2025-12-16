---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_grant_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit-router-grant-attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Grant Attachment resource.
---

# alicloud_cen_transit_router_grant_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router Grant Attachment resource.

For information about Cloud Enterprise Network (CEN) Transit Router Grant Attachment and how to use it, see [What is Transit Router Grant Attachment](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/grantinstancetotransitrouter).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_grant_attachment&exampleId=d3da962c-17b8-e48b-28ad-cc619d02ca0d7a77d422&activeTab=example&spm=docs.r.cen_transit_router_grant_attachment.0.d3da962c17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}

resource "alicloud_vpc" "example" {
  vpc_name   = "tf_example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_transit_router_grant_attachment" "example" {
  cen_id        = alicloud_cen_instance.example.id
  cen_owner_id  = data.alicloud_account.default.id
  instance_id   = alicloud_vpc.example.id
  instance_type = "VPC"
  order_type    = "PayByCenOwner"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_transit_router_grant_attachment&spm=docs.r.cen_transit_router_grant_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cen_id`- (Required, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance to which the transit router belongs.
* `cen_owner_id` - (Required, ForceNew) The ID of the Alibaba Cloud account to which the CEN instance belongs.
* `instance_id` - (Required, ForceNew) The ID of the network instance.
* `instance_type` - (Required, ForceNew) The type of the network instance. Valid values: `VPC`, `ExpressConnect`, `VPN`.
* `order_type` - (Optional, Computed, ForceNew) The entity that pays the fees of the network instance. Valid values: `PayByResourceOwner`, `PayByCenOwner`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Transit Router Grant Attachment. The value formats as `<instance_type>:<instance_id>:<cen_owner_id>:<cen_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the cen transit router grant attachment.
* `delete` - (Defaults to 1 mins) Used when delete the cen transit router grant attachment.


## Import

Cloud Enterprise Network (CEN) Transit Router Grant Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_grant_attachment.example <instance_type>:<instance_id>:<cen_owner_id>:<cen_id>
```