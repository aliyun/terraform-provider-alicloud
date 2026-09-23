---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance_attachment"
sidebar_current: "docs-alicloud-resource-ots-instance-attachment"
description: |-
  Provides an OTS (Open Table Service) resource to attach VPC to instance.
---

# alicloud_ots_instance_attachment

This resource will help you to bind a VPC to an OTS instance.

-> **NOTE:** Available since v1.10.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ots_instance_attachment&exampleId=74477c1c-a6a2-5067-51ff-ce85bf235b767fc0d5d4&activeTab=example&spm=docs.r.ots_instance_attachment.0.74477c1ca6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Vpc"
  tags = {
    Created = "TF",
    For     = "example",
  }
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ots_instance_attachment" "default" {
  instance_name = alicloud_ots_instance.default.name
  vpc_name      = "examplename"
  vswitch_id    = alicloud_vswitch.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ots_instance_attachment&spm=docs.r.ots_instance_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance.
* `vpc_name` - (Required, ForceNew) The name of attaching VPC to instance. It can only contain letters and numbers, must start with a letter, and is limited to 3-16 characters in length.
* `vswitch_id` - (Required, ForceNew) The ID of attaching VSwitch to instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as "instance_name".
* `vpc_id` - The ID of attaching VPC to instance.


