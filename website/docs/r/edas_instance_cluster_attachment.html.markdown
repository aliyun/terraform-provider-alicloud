---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_instance_cluster_attachment"
sidebar_current: "docs-alicloud-resource-edas-instance-cluster-attachment"
description: |-
  Provides an EDAS instance cluster attachment resource.
---

# alicloud_edas_instance_cluster_attachment

Provides an EDAS instance cluster attachment resource, see [What is EDAS Instance Cluster Attachment](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-installagent).

-> **NOTE:** Available since v1.82.0.


## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_instance_cluster_attachment&exampleId=2f892ae6-8656-f488-10d5-806421c51d104f786c6e&activeTab=example&spm=docs.r.edas_instance_cluster_attachment.0.2f892ae686&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
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

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}


resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_edas_cluster" "default" {
  cluster_name      = var.name
  cluster_type      = "2"
  network_mode      = "2"
  logical_region_id = data.alicloud_regions.default.regions.0.id
  vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_edas_instance_cluster_attachment" "default" {
  cluster_id   = alicloud_edas_cluster.default.id
  instance_ids = [alicloud_instance.default.id]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_instance_cluster_attachment&spm=docs.r.edas_instance_cluster_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the cluster that you want to create the application.
* `instance_ids` - (Required, ForceNew) The ID of instance. Type: list.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<cluster_id>:<instance_id1,instance_id2>`.
* `status_map` - The status map of the resource supplied above. The key is instance_id and the values are 1(running) 0(converting) -1(failed) and -2(offline).
* `ecu_map` - The ecu map of the resource supplied above. The key is instance_id and the value is ecu_id.
* `cluster_member_ids` - The cluster members map of the resource supplied above. The key is instance_id and the value is cluster_member_id.


