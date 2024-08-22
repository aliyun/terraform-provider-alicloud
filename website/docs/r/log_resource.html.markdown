---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_resource"
sidebar_current: "docs-alicloud-resource-log-resource"
description: |-
  Provides a Alicloud log resource.
---

# alicloud_log_resource

Log resource is a meta store service provided by log service, resource can be used to define meta store's table structure. 

For information about SLS Resource and how to use it, see [Resource management](https://www.alibabacloud.com/help/en/doc-detail/207732.html)

-> **NOTE:** Available since v1.162.0. log resource region should be set a main region: cn-heyuan.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_log_resource&exampleId=79d92b9d-00f8-cdf3-3e28-5e3bbf189fe8f866022f&activeTab=example&spm=docs.r.log_resource.0.79d92b9d00" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The meta store's name, can be used as table name.
* `type` - (Required) The meta store's type, userdefine e.g.
* `description` - (Optional) The meta store's description.
* `schema` - (Required) The meta store's schema info, which is json string format, used to define table's fields.
* `ext_info` - (Optional) The ext info of meta store.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It formats of `<name>`.

## Import

Log resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_resource_record.example <id>
```
