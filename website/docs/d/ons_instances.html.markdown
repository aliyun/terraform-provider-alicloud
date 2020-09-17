---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_instances"
sidebar_current: "docs-alicloud-datasource-ons-instances"
description: |-
    Provides a list of ons instances available to the user.
---

# alicloud\_ons\_instances

This data source provides a list of ONS Instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.52.0+

## Example Usage

```terraform
variable "name" {
  default = "onsInstanceDatasourceName"
}

resource "alicloud_ons_instance" "default" {
  name   = "${var.name}"
  remark = "default_ons_instance_remark"
}

data "alicloud_ons_instances" "instances_ds" {
  ids         = ["${alicloud_ons_instance.default.id}"]
  name_regex  = "${alicloud_ons_instance.default.name}"
  output_file = "instances.txt"
}

output "first_instance_id" {
  value = "${data.alicloud_ons_instances.instances_ds.instances.0.instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by the instance name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, Available in v1.97.0+) The status of Ons instance. Valid values: `0` deploying, `2` arrears, `5` running, `7` upgrading.
* `enable_details` - (Optional, Available in v1.97.0+) Default to `false`. Set it to true can output more details.
* `tags` - (Optional, Available in v1.97.0+) A map of tags assigned to the Ons instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance.
  * `instance_type` - The type of the instance. Read [Fields in InstanceVO](https://www.alibabacloud.com/help/doc-detail/106351.html) for further details.
  * `instance_status` - The status of the instance. Read [Fields in InstanceVO](https://www.alibabacloud.com/help/doc-detail/106351.html) for further details.
  * `status` - The status of the instance. Read [Fields in InstanceVO](https://www.alibabacloud.com/help/doc-detail/106351.html) for further details.
  * `release_time` - The automatic release time of an Enterprise Platinum Edition instance.
  * `tags` - A map of tags assigned to the Ons instance.
  * `remark` - This attribute is a concise description of instance.
  * `http_internet_endpoint` - The public HTTP endpoint for the Message Queue for Apache RocketMQ instance.
  * `http_internal_endpoint` - The internal HTTP endpoint for the Message Queue for Apache RocketMQ instance.
  * `http_internet_secure_endpoint` - The public HTTPS endpoint for the Message Queue for Apache RocketMQ instance.
  * `tcp_endpoint` - The TCP endpoint for the Message Queue for Apache RocketMQ instance.
  * `independent_naming` - Indicates whether any namespace is configured for the Message Queue for Apache RocketMQ instance.
