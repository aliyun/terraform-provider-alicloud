---
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_instance"
sidebar_current: "docs-alicloud-resource-elasticsearch-instance"
description: |-
  Provides a Alicloud Elasticsearch instance resource.
---

# alicloud\_elasticsearch\_instance

Provides a Elasticsearch instance resource. It contains data nodes, dedicated master node(optional) and etc. It can be associated with private IP whitelists and kibana IP whitelist.

~> **NOTE:** Only one operation is supported in a request. So if `data_node_spec` and `data_node_disk_size` are both changed, system will respond error.

## Example Usage

Basic Usage

```
resource "alicloud_elasticsearch" "default" {
  instance_charge_type = "PostPaid"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  vswitch_id           = "some vswitch id"
  password             = "Your password"
  version              = "5.5.3_with_X-Pack"
  description          = "description"
}
```
## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of instance. It a string of 0 to 30 characters..
* `instance_charge_type` - (Required) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `period` - (Optional) The duration that you will buy Elasticsearch instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `data_node_amount` - (Required) The Elasticsearch cluster's data node quantity, between 2 and 50.
* `data_node_spec` - (Required) The data node specifications of the Elasticsearch instance.
* `data_node_disk_size` - (Required) The single data node storage space. An SSD (`cloud_ssd`) supports a maximum of 2048 GiB (2 TB). An ultra disk (`cloud_efficiency`) supports a maximum of 5120 GiB (5 TB). If the data to be stored is larger than 2048 GiB, an ultra disk can only support the following data sizes: 2560 GiB, 3072 GiB, 3584 GiB, 4096 GiB, 4608 GiB, or 5120 GiB.
* `data_node_disk_type` - (Required) The data node disk type. Supported values: cloud_ssd, cloud_efficiency.
* `vswitch_id` - (Required) The ID of VSwitch.
* `password` - (Required) The password of the instance. The password can be 8 to 32 characters in length and must contain three of the following conditions: uppercase letters, lowercase letters, numbers, and special characters (!@#$%^&*()_+-=).
* `version` - (Required) Elasticsearch version. Supported values: 5.5.3_with_X-Pack and 6.3_with_X-Pack.
* `private_whitelist` - (Optional) Set the instance's IP whitelist in VPC network.
* `kibana_whitelist` - (Optional) Set the Kibana's IP whitelist in internet network.
* `master_node_spec` - (Optional) The dedicated master node spec. If specified, dedicated master node will be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Elasticsearch instance.
* `domain` - Instance connection domain (only VPC network access supported).
* `port` - Instance connection port.
* `kibana_domain` - Kibana console domain (Internet access supported).
* `kibana_port` - Kibana console port.
* `status` - The Elasticsearch instance status. Includes `active`, `activating`, `inactive`. Some operations are denied when status is not `active`.

## Import

Elasticsearch can be imported using the id, e.g.

```
$ terraform import alicloud_elasticsearch.example es-cn-abcde123456
```

