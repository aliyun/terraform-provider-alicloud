---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_resource_record"
sidebar_current: "docs-alicloud-resource-log-resource-record"
description: |-
  Provides a Alicloud log resource record.
---

# alicloud_log_resource_record

Log resource is a meta store service provided by log service, resource can be used to define meta store's table structure, record can be used for table's row data. 

For information about SLS Resource and how to use it, see [Resource management](https://www.alibabacloud.com/help/en/doc-detail/207732.html)

-> **NOTE:** Available since v1.162.0. log resource region should be set a main region: cn-heyuan.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-heyuan"
}

resource "alicloud_log_resource" "example" {
  type        = "userdefine"
  name        = "user.tf.resource"
  description = "user tf resource desc"
  ext_info    = "{}"
  schema      = <<EOF
    {
      "schema": [
        {
          "column": "col1",
          "desc": "col1   desc",
          "ext_info": {
          },
          "required": true,
          "type": "string"
        },
        {
          "column": "col2",
          "desc": "col2   desc",
          "ext_info": "optional",
          "required": true,
          "type": "string"
        }
      ]
    }
  EOF
}

resource "alicloud_log_resource_record" "example" {
  resource_name = alicloud_log_resource.example.id
  record_id     = "user_tf_resource_1"
  tag           = "resource tag"
  value         = <<EOF
    {
      "col1": "this is col1 value",
      "col2": "col2   value"
    }
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `resource_name` - (Required) The name defined in log_resource, log service have some internal resource, like sls.common.user, sls.common.user_group.
* `record_id` - (Required, ForceNew) The record's id, should be unique.
* `tag` - (Required) The record's tag, can be used for search.
* `value` - (Required) The json value of record.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource record. It formats of `<resource_name>:<record_id>`.

## Import

Log resource record can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_resource_record.example <resource_name>:<record_id>
```
