---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vpd_grant_rule"
description: |-
  Provides a Alicloud Eflo Vpd Grant Rule resource.
---

# alicloud_eflo_vpd_grant_rule

Provides a Eflo Vpd Grant Rule resource.

Lingjun Network Segment Cross-Account Authorization Information.

For information about Eflo Vpd Grant Rule and how to use it, see [What is Vpd Grant Rule](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateVpdGrantRule).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_vpd_grant_rule&exampleId=26277719-74f4-2182-e98c-6bb7408d0229b75ac7c0&activeTab=example&spm=docs.r.eflo_vpd_grant_rule.0.2627771974&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-hangzhou-a"
}

resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = var.name
}

resource "alicloud_eflo_vpd_grant_rule" "default" {
  grant_tenant_id = data.alicloud_account.default.id
  er_id           = alicloud_eflo_er.default.id
  instance_id     = alicloud_eflo_vpd.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_vpd_grant_rule&spm=docs.r.eflo_vpd_grant_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `er_id` - (Required, ForceNew) The ID of the ER instance under the cross-account tenant.
* `grant_tenant_id` - (Required, ForceNew) Cross-account authorized tenant ID.
* `instance_id` - (Required, ForceNew) Instance ID of VPD.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The Creation time.
* `region_id` - The Region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpd Grant Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Vpd Grant Rule.

## Import

Eflo Vpd Grant Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vpd_grant_rule.example <id>
```
