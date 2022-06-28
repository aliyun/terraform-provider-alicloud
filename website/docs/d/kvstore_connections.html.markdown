---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_connections"
sidebar_current: "docs-alicloud-datasource-kvstore-connections"
description: |-
    Query the public IP of the specified KVStore DBInstance.
---

# alicloud\_kvstore\_connections

This data source can query the public IP of the specified KVStore DBInstance.
 
-> **NOTE:** Available in v1.101.0+.

## Example Usage

```terraform
# Declare the data source

data "alicloud_kvstore_connections" "example" {
  ids = ["r-wer123456"]
}

output "connection_string" {
  value = data.alicloud_kvstore_connections.example.connections.0.connection_string
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Required) A list of KVStore DBInstance ids, only support one item.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of KVStore DBInstance ids.
* `connections` - Public network details of the specified resource. contains the following attributes:
  * `connection_string` - The connection string of the instance.
  * `db_instance_net_type` - The network type of the instance.
  * `expired_time` - The expiration time of the classic network address.
  * `ip_address` - The IP address of the instance.
  * `port` - The port number of the instance.
  * `upgradeable` - The remaining validity period of the endpoint of the classic network.
  * `vpc_id` - The ID of the VPC where the instance is deployed.
  * `vpc_instance_id` - The ID of the instance. It is returned only when the value of the DBInstanceNetType parameter is 2 (indicating VPC).
  * `vswitch_id` - The ID of the VSwitch.
