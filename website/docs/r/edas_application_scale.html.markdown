---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_application_scale"
sidebar_current: "docs-alicloud-resource-edas-application-scale"
description: |-
  This operation is provided to scale out an EDAS application.
---

# alicloud_edas_application_scale

This operation is provided to scale out an EDAS application, see [What is EDAS Application Scale](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-scaleoutapplication).


-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

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
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
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
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = var.name
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  internet_max_bandwidth_out = "10"
  internet_charge_type       = "PayByTraffic"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
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

resource "alicloud_edas_application" "default" {
  application_name = var.name
  cluster_id       = alicloud_edas_cluster.default.id
  package_type     = "WAR"
}

resource "alicloud_edas_deploy_group" "default" {
  app_id     = alicloud_edas_application.default.id
  group_name = var.name
}

data "alicloud_edas_deploy_groups" "default" {
  app_id = alicloud_edas_deploy_group.default.app_id
}

resource "alicloud_edas_application_scale" "default" {
  app_id       = alicloud_edas_application.default.id
  deploy_group = data.alicloud_edas_deploy_groups.default.groups.0.group_id
  ecu_info     = [alicloud_edas_instance_cluster_attachment.default.ecu_map[alicloud_instance.default.id]]
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy.
* `deploy_group` - (Required, ForceNew) The ID of the instance group to which you want to add ECS instances to scale out the application.
* `ecu_info` - (Required, ForceNew) The IDs of the Elastic Compute Unit (ECU) where you want to deploy the application. Type: List.
* `force_status` - (Optional) This parameter specifies whether to forcibly remove an ECS instance where the application is deployed. It is set as true only after the ECS instance expires. In normal cases, this parameter do not need to be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<ecu1,ecu2>`.
* `ecc_info` - The ecc information of the resource supplied above. The value is formulated as `<ecc1,ecc2>`.

