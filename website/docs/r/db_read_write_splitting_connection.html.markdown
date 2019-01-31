---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_read_write_splitting_connection"
sidebar_current: "docs-alicloud-resource-db-read-write-splitting-connection"
description: |-
  Provides an RDS instance read write splitting connection resource.
---

# alicloud\_db\_read\_write\_splitting\_connection

Provides an RDS read write splitting connection resource to allocate an Intranet connection string for RDS instance.

## Example Usage

```
resource "alicloud_db_read_write_splitting_connection" "foo" {
    instance_id = "rm-2eps..."
    connection_prefix = "test-connection"
	distribution_type = "Custom"
	max_delay_time = 300
	weight = "{\"rm-2eps...\":\"500\"}"
	
	depends_on = ["alicloud_db_readonly_instance.foo"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance that can run database.
* `distribution_type` - (Required) Read weight distribution mode. Values are as follows: `Standard` indicates automatic weight distribution based on types, `Custom` indicates custom weight distribution. 
* `connection_prefix` - (Optional) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'rw'.
* `port` - (Optional) Intranet connection port. Valid value: [3001-3999]. Default to 3306.
* `max_delay_time` - (Optional) Delay threshold, in seconds. The value range is 0 to 7200. Default to 30. Read requests are not routed to the read-only instances with a delay greater than the threshold.  
* `weight` - (Optional) Read weight distribution. Read weights increase at a step of 100 up to 10,000. Enter weights in the following format: {"Instanceid":"Weight","Instanceid":"Weight"}. This parameter must be set when distribution_type is set to Custom. 

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_prefix` - Prefix of a connection string.
* `port` - Connection instance port.
* `connection_string` - Connection instance string.

## Import

RDS read write splitting connection can be imported using the id, e.g.

```
$ terraform import alicloud_db_read_write_splitting_connection.example abc12345678
```