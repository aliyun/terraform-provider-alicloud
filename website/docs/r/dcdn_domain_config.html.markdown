---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_domain_config"
sidebar_current: "docs-alicloud-resource-dcdn-domain-config"
description: |-
  Provides a Alicloud Dcdn domain config Resource.
---

# alicloud_dcdn_domain_config

Provides a DCDN Accelerated Domain resource.

For information about domain config and how to use it, see [Batch set config](https://www.alibabacloud.com/help/en/doc-detail/130632.htm).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_domain_config&exampleId=d0ea9ec9-32b7-28de-5d96-89f9fbb1a5e0755a3552&activeTab=example&spm=docs.r.dcdn_domain_config.0.d0ea9ec932&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "domain_name" {
  default = "alibaba-example.com"
}

resource "alicloud_dcdn_domain" "example" {
  domain_name = var.domain_name
  scope       = "overseas"
  status      = "online"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = 20
    port     = 80
    weight   = 10
  }
}

resource "alicloud_dcdn_domain_config" "ip_allow_list_set" {
  domain_name   = alicloud_dcdn_domain.example.domain_name
  function_name = "ip_allow_list_set"
  function_args {
    arg_name  = "ip_list"
    arg_value = "192.168.0.1"
  }
}

resource "alicloud_dcdn_domain_config" "referer_white_list_set" {
  domain_name   = alicloud_dcdn_domain.example.domain_name
  function_name = "referer_white_list_set"
  function_args {
    arg_name  = "refer_domain_allow_list"
    arg_value = "110.110.110.110"
  }
}

resource "alicloud_dcdn_domain_config" "filetype_based_ttl_set" {
  domain_name   = alicloud_dcdn_domain.example.domain_name
  function_name = "filetype_based_ttl_set"
  function_args {
    arg_name  = "ttl"
    arg_value = "300"
  }
  function_args {
    arg_name  = "file_type"
    arg_value = "jpg"
  }
  function_args {
    arg_name  = "weight"
    arg_value = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `function_name` - (Required, ForceNew) The name of the domain config.
* `parent_id` - (Optional, Available since v1.221.0) By configuring the function condition (rule engine) in the domain name configuration function parameters, Rule conditions can be created (Rule conditions can match and filter user requests by identifying various parameters carried in user requests). After each rule condition is created, a corresponding ConfigId will be generated, and the ConfigId can be referenced by other functions as a ParentId parameter, in this way, the rule conditions can be combined with the functional configuration to form a more flexible configuration.
* `function_args` - (Required, Set) The args of the domain config. See [`function_args`](#function_args) below.

### `function_args`

The function_args block supports the following:

* `arg_name` - (Required) The name of arg.
* `arg_value` - (Required) The value of arg.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Config. It formats as `<domain_name>:<function_name>:<config_id>`.
* `config_id` - The ID of the configuration.
* `status` -  The status of the Config.

## Import

DCDN domain config can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_domain_config.example <domain_name>:<function_name>:<config_id>
```
