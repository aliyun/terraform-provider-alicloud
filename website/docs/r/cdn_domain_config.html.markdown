---
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_doamin_config"
sidebar_current: "docs-alicloud-resource-cdn-domain-config"
description: |-
  Provides a Alicloud Cdn domain config  Resource.
---

# alicloud_cdn_domain_config

Provides a CDN Accelerated Domain resource.

For information about domain config and how to use it, see [Batch set config](https://www.alibabacloud.com/help/doc-detail/91179.html)

-> **NOTE:** Available in v1.34.0+.

## Example Usage

Basic Usage

```
# Create a new Domain config.
resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
      scope="overseas"
      sources {
         content = "1.1.1.1"
         type = "ipaddr"
         priority = "20"
         port = 80
         weight = "15"
      }
}
resource "alicloud_cdn_domain_config" "config" {
	  domain_name = "${alicloud_cdn_domain_new.domain.domain_name}"
	  function_name = "ip_allow_list_set"
      function_args {
            arg_name = "ip_list"
            arg_value = "110.110.110.110"
      }
}

```
## Argument Reference

The following arguments are supported:
* `domain_name` - (Required, ForceNew) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `function_name` - (Required, ForceNew) The name of the domain config.
* `function_args` - (Required, ForceNew, Type: list) The args of the domain config.

### Block function_args

The `function_args` block supports the following:

* `arg_name` - (Required) The name of arg.
* `arg_value` - (Required) The value of arg.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the domain config. The value is formate as `<domain_name>:<function_name>`.

## Import

CDN domain config can be imported using the id, e.g.

```
terraform import alicloud_cdn_domain_config.example cdn:config-abc123456
```
