---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_slb_attachment"
sidebar_current: "docs-alicloud-resource-edas-slb-attachment"
description: |-
  Binds SLB to an EDAS application.
---

# alicloud_edas_slb_attachment

Binds SLB to an EDAS application.

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_slb_attachment&exampleId=5ce98a60-4c59-31dc-3860-d5584340f8a2ec82cac4&activeTab=example&spm=docs.r.edas_slb_attachment.0.5ce98a604c&intl_lang=EN_US" target="_blank">
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
  package_type     = "JAR"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_edas_slb_attachment" "default" {
  app_id = alicloud_edas_application.default.id
  slb_id = alicloud_slb_load_balancer.default.id
  slb_ip = alicloud_slb_load_balancer.default.address
  type   = alicloud_slb_load_balancer.default.address_type
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_slb_attachment&spm=docs.r.edas_slb_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application to which you want to bind an SLB instance.
* `slb_id` - (Required, ForceNew) The ID of the SLB instance that is going to be bound.
* `slb_ip` - (Required, ForceNew) The IP address that is allocated to the bound SLB instance.
* `type` - (Required, ForceNew) The type of the bound SLB instance.
* `listener_port` - (Optional, ForceNew) The listening port for the bound SLB instance.
* `vserver_group_id` - (Optional, ForceNew) The ID of the virtual server (VServer) group associated with the intranet SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<slb_id>`.
* `slb_status` - Running Status of SLB instance. Inactive：The instance is stopped, and listener will not monitor and forward traffic. Active：The instance is running. After the instance is created, the default state is active. Locked：The instance is locked, the instance has been owed or locked by Alibaba Cloud. Expired: The instance has expired.
* `vswitch_id` - VPC related vswitch ID.


