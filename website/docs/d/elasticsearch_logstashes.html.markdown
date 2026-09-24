---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_logstashes"
sidebar_current: "docs-alicloud-datasource-elasticsearch-logstashes"
description: |-
  Provides a list of Elasticsearch Logstash owned by an Alibaba Cloud account.
---

# alicloud_elasticsearch_logstashes

This data source provides Elasticsearch Logstash available to the user.[What is Logstash](https://next.api.alibabacloud.com/document/elasticsearch/2017-06-13/CreateLogstash)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_elasticsearch_logstash" "default" {
  description = "tf-acc-create-test-1"
  version     = "7.4_with_X-Pack"
  node_spec {
    disk_type = "cloud_efficiency"
    spec      = "elasticsearch.sn1ne.large"
    disk      = "20"
  }
  network_config {
    type       = "vpc"
    vpc_id     = "vpc-bp1jy348ibzulk6hn65xf"
    vswitch_id = "vsw-bp13kz5zhn6flmqmh9fyn"
    vs_area    = "cn-hangzhou-i"
  }
  payment_type = "Subscription"
  node_amount  = "1"
}

data "alicloud_elasticsearch_logstashes" "default" {
  ids               = ["${alicloud_elasticsearch_logstash.default.id}"]
  description       = "tf-acc-create-test-1"
  resource_group_id = ""
  version           = "7.4_with_X-Pack"
}

output "alicloud_elasticsearch_logstash_example_id" {
  value = data.alicloud_elasticsearch_logstashes.default.logstashes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (ForceNew, Optional) Description
* `instance_id` - (ForceNew, Optional) Resource id
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group
* `version` - (ForceNew, Optional) Version number
* `ids` - (Optional, Computed) A list of Logstash IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Logstash IDs.
* `logstashes` - A list of Logstash Entries. Each element contains the following attributes:
    * `create_time` - Instance creation time.
    * `description` - Description.
    * `instance_id` - Resource id.
    * `network_config` - VPC configuration.
        * `type` - Network type.
        * `vswitch_id` - Unique identification of switch network.
        * `vpc_id` - VPC unique identifier.
        * `vs_area` - Unique ID of zone.
    * `node_amount` - Number of nodes.
    * `node_spec` - Elasticsearch node disk.
        * `disk` - Disk size.
        * `disk_type` - Disk type.
        * `spec` - Disk specification.
    * `payment_type` - The payment type of the resource.
    * `resource_group_id` - The ID of the resource group.
    * `status` - The status of the resource.
    * `tags` - The tag of the resource.
    * `updated_at` - Instance update time.
    * `version` - Version number.
    * `id` - The ID of the resource supplied above.
