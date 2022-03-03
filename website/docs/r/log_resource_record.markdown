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

-> **NOTE:** Available in 1.160.0

## Example Usage

Basic Usage

```
resource "alicloud_log_resource_record" "example_6" {
  resource_name         = "user.tf.test_resource"
  record_id             = "user.tf.test_resource-1"
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
* `create_time` - (Optional) The create time of record, unixtimestamp.
* `last_modify_time` - (Optional) The last_modify time of record, unixtimestamp.
