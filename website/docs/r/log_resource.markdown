---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_resource"
sidebar_current: "docs-alicloud-resource-log-resource"
description: |-
  Provides a Alicloud log resource.
---

# alicloud\_log\_resource

Log resource is a meta store service provided by log service, resource can be used to define meta store's table structure. 

For information about SLS Resource and how to use it, see [Resource management](https://www.alibabacloud.com/help/en/doc-detail/207732.html)

-> **NOTE:** Available in 1.160.0

## Example Usage

Basic Usage

```
resource "alicloud_log_resource" "example" {
  name                  = "user.tf.test_resource"
  description           = "user tf test resource desc"
  schema                = "{\"schema\":[{\"column\":\"col1\",\"desc\":\"col1 desc\",\"ext_info\":{},\"required\":true,\"type\":\"string\"},{\"column\":\"col2\",\"desc\":\"col2 desc\",\"ext_info\":\"optional\",\"required\":true,\"type\":\"string\"}]}"
}


```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The meta store's name, can be used as table name.
* `description` - (Optional) The meta store's description.
* `schema` - (Required) The meta store's schema info, which is json string format, used to define table's fields.
* `ext_info` - (Optional) The ext info of meta store.
* `create_time` - (Optional) The create time of meta store, unixtimestamp.
* `last_modify_time` - (Optional) The last_modify time of meta store, unixtimestamp.
