---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_instances"
sidebar_current: "docs-alicloud-datasource-alikafka-instances"
description: |-
    Provides a list of alikafka instances available to the user.
---

# alicloud\_alikakfa\_instances

This data source provides a list of ALIKAFKA Instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.59.0+

## Example Usage

```
variable "instance_name" {
 default = "alikafkaInstanceName"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}
resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.instance_name}"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "4"
  io_max = "20"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

data "alicloud_alikafka_instances" "instances_ds" {
  name_regex = "alikafkaInstanceName"
  output_file = "instances.txt"
}

output "first_instance_name" {
  value = "${data.alicloud_alikafka_instances.instances_ds.instances.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by the instance name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `name` - Name of the instance.
  * `create_time` - The create time of the instance.
  * `service_status` - The current status of the instance. -1: unknown status, 0: wait deploy, 1: initializing, 2: preparing, 3 starting, 5: in service, 7: wait upgrade, 8: upgrading, 10: released, 15: freeze, 101: deploy error, 102: upgrade error. 
  * `deploy_type` - The deploy type of the instance. 0: sharing instance, 1: vpc instance, 2: vpc instance(support ip mapping), 3: eip instance, 4: eip/vpc instance, 5: vpc instance.
  * `vpc_id` - The ID of attaching VPC to instance.
  * `vswitch_id` - The ID of attaching vswitch to instance.
  * `io_max` - The peak value of io of the instance.
  * `eip_max` - The peak bandwidth of the instance.
  * `disk_type` - The disk type of the instance. 0: efficient cloud disk , 1: SSD.
  * `disk_size` - The disk size of the instance.
  * `topic_quota` - The max num of topic can be create of the instance.
  * `zone_id` - The ID of attaching zone to instance.
  * `paid_type` - The paid type of the instance.
  * `spec_type` - The spec type of the instance.

