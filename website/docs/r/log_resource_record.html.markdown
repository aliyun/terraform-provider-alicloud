---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_resource_record"
sidebar_current: "docs-alicloud-resource-log-resource-record"
description: |-
  Provides a Alicloud log resource record.
---

# alicloud\_log\_resource\_record

Log resource is a meta store service provided by log service, resource can be used to define meta store's table structure, record can be used for table's row data. 

For information about SLS Resource and how to use it, see [Resource management](https://www.alibabacloud.com/help/en/doc-detail/207732.html)

-> **NOTE:** Available in 1.162.0+, log resource region should be set a main region: cn-heyuan

## Example Usage

Basic Usage

```
resource "alicloud_log_resource_record" "example" {
  resource_name         = "user.tf.test_resource"
  record_id             = "user_tf_test_resource_1"
  tag                   = "test resource tag"
  value                 = "{\"col1\": \"this is col1 value\", \"col2\": \"col2 value\"}"
}
```
## Argument Reference

The following arguments are supported:

* `resource_name` - (Required) The name defined in log_resource, log service have some internal resource, like sls.common.user, sls.common.user_group.
* `record_id` - (Required) The record's id, should be unique.
* `tag` - (Required) The record's tag, can be used for search.
* `value` - (Required) The json value of record.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource record. It formats of `<resource_name>:<record_id>`.

## Import

Log resource record can be imported using the id, e.g.

```
$ terraform import alicloud_log_resource_record.example user.tf.test_resource:user_tf_test_resource_1
```
