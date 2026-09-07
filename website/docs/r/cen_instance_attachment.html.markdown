---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_attachment"
sidebar_current: "docs-alicloud-resource-cen-instance-attachment"
description: |-
  Provides a Alicloud CEN child instance attachment resource.
---

# alicloud_cen_instance_attachment

Provides a CEN child instance attachment resource that associate the network(VPC, CCN, VBR) with the CEN instance.

-> **NOTE:** Available since v1.42.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_instance_attachment&exampleId=15793616-df4e-06df-591f-82ab8c62d69ed2f00bf2&activeTab=example&spm=docs.r.cen_instance_attachment.0.15793616df&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "example" {
  vpc_name   = "tf_example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_instance_attachment" "example" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.example.id
  child_instance_type      = "VPC"
  child_instance_region_id = data.alicloud_regions.default.regions.0.id
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_instance_attachment&spm=docs.r.cen_instance_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `child_instance_id` - (Required, ForceNew) The ID of the child instance to attach.
* `child_instance_region_id` - (Required, ForceNew) The region ID of the child instance to attach.
* `child_instance_owner_id` - (Optional, Available since v1.42.0) The uid of the child instance. Only used when attach a child instance of other account.
* `child_instance_type` - (Required, ForceNew, Available since v1.97.0) The type of the associated network. Valid values: `VPC`, `VBR` and `CCN`.
* `cen_owner_id` - (Optional, Available since v1.97.0) The account ID to which the CEN instance belongs.

->**NOTE:** Ensure that the child instance is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<instance_id>:<child_instance_id>:<child_instance_type>:<child_instance_region_id>`. Before version 1.97.0, the value is `<instance_id>:<child_instance_id>`.
* `status` - The associating status of the network.

## Timeouts

-> **NOTE:** Available since v1.199.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the child instance attachment.
* `delete` - (Defaults to 10 mins) Used when delete the child instance attachment.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance_attachment.example cen-m7i7pjmkon********:vpc-2ze2w07mcy9nz********:VPC:cn-beijing
```
