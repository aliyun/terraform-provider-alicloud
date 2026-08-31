---
subcategory: "Graph Database"
layout: "alicloud"
page_title: "Alicloud: alicloud_graph_database_db_instance"
sidebar_current: "docs-alicloud-resource-graph-database-db-instance"
description: |-
  Provides a Alicloud Graph Database Db Instance resource.
---

# alicloud_graph_database_db_instance

Provides a Graph Database Db Instance resource.

For information about Graph Database Db Instance and how to use it, see [What is Db Instance](https://www.alibabacloud.com/help/en/graph-compute/latest/placeholder).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_graph_database_db_instance&exampleId=703d51a5-0232-e3ef-2323-33f25db64b34f9b8dd08&activeTab=example&spm=docs.r.graph_database_db_instance.0.703d51a502&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_graph_database_db_instance" "example" {
  db_node_class            = "gdb.r.2xlarge"
  db_instance_network_type = "vpc"
  db_version               = "1.0"
  db_instance_category     = "HA"
  db_instance_storage_type = "cloud_ssd"
  db_node_storage          = "50"
  payment_type             = "PayAsYouGo"
  db_instance_description  = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_graph_database_db_instance&spm=docs.r.graph_database_db_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `db_instance_category` - (Required, ForceNew) The category of the db instance. Valid values: `HA`, `SINGLE`(Available in 1.173.0+).
* `db_instance_description` - (Optional) According to the practical example or notes.
* `db_instance_network_type` - (Required, ForceNew) The network type of the db instance. Valid values: `vpc`.
* `db_instance_storage_type` - (Required, ForceNew) Disk storage type. Valid values: `cloud_essd`, `cloud_ssd`. Modification is not supported.
* `db_node_class` - (Required) The class of the db node. Valid values: `gdb.r.xlarge`, `gdb.r.2xlarge`, `gdb.r.4xlarge`, `gdb.r.8xlarge`, `gdb.r.16xlarge`, `gdb.r.xlarge_basic`, `gdb.r.2xlarge_basic`, `gdb.r.4xlarge_basic`, `gdb.r.8xlarge_basic`, `gdb.r.16xlarge_basic`.
* `db_node_storage` - (Required) Instance storage space, which is measured in GB.
* `db_version` - (Required, ForceNew) Kernel Version. Valid values: `1.0` or `1.0-OpenCypher`. `1.0`: represented as gremlin, `1.0-OpenCypher`: said opencypher.
* `payment_type` - (Required, ForceNew) The paymen type of the resource. Valid values: `PayAsYouGo`.
* `db_instance_ip_array` - (Optional) IP ADDRESS whitelist for the instance group list. See [`db_instance_ip_array`](#db_instance_ip_array) below.
* `vswitch_id` - (Optional, ForceNew, Available since v1.171.0) The ID of attaching vswitch to instance.
* `vpc_id` - (Optional, ForceNew, Available since v1.171.0) ID of the VPC.
* `zone_id` - (Optional, ForceNew, Available since v1.171.0) The zone ID of the resource.

### `db_instance_ip_array`

The db_instance_ip_array supports the following:

* `db_instance_ip_array_attribute` - (Optional) The default is empty. To distinguish between the different property console does not display a `hidden` label grouping.
* `db_instance_ip_array_name` - (Optional) IP ADDRESS whitelist group name.
* `security_ips` - (Optional) IP ADDRESS whitelist addresses in the IP ADDRESS list, and a maximum of 1000 comma-separated format is as follows: `0.0.0.0/0` and `10.23.12.24`(IP) or `10.23.12.24/24`(CIDR mode, CIDR (Classless Inter-Domain Routing)/24 represents the address prefixes in the length of the range [1,32]).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Db Instance.
* `status` - Instance status. Value range: `Creating`, `Running`, `Deleting`, `Rebooting`, `DBInstanceClassChanging`, `NetAddressCreating` and `NetAddressDeleting`.
* `connection_string` - (Available in 1.196.0+)  The connection string of the instance.
* `port` - (Available in 1.196.0+) The connection port of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when create the Db Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Db Instance.
* `update` - (Defaults to 60 mins) Used when update the Db Instance.

## Import

Graph Database Db Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_graph_database_db_instance.example <id>
```
