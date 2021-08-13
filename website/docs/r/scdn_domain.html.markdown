---
subcategory: "SCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_scdn_domain"
sidebar_current: "docs-alicloud-resource-scdn-domain"
description: |-
  Provides a Alicloud SCDN Domain resource.
---

# alicloud\_scdn\_domain

Provides a SCDN Domain resource.

For information about SCDN Domain and how to use it, see [What is Domain](https://help.aliyun.com/document_detail/63672.html).

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_scdn_domain" "example" {
  domain_name = "my-Domain"
  sources {
    content  = "xxx.aliyuncs.com"
    enabled  = "online"
    port     = 80
    priority = "20"
    type     = "oss"
  }
  domain_configs {
    function_name = "referer_white_list_set"
    function_args {
      arg_name = "refer_domain_allow_list"
      arg_value = "110.110.110.110"
    }
  }
  domain_configs {
    function_name = "filetype_based_ttl_set"
    function_args {
      arg_name = "ttl"
      arg_value = "330"
    }
    function_args {
      arg_name = "file_type"
      arg_value = "jpg"
    }
    function_args {
      arg_name = "weight"
      arg_value = "1"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `biz_name` - (Optional) from the Business Type Drop-down List. Valid values: `download`, `image`, `scdn`, `video`.
* `cert_infos` - (Optional) Certificate Information. See the following `Block cert_infos`.
* `check_url` - (Optional) The health check url.
* `domain_configs` - (Optional) Domain name configuration data list. See the following `Block domain_configs`.
* `domain_name` - (Required, ForceNew) The name of domain.
* `force_set` - (Optional) Whether to set certificate forcibly.
* `resource_group_id` - (Optional, Computed) The resource group id.
* `sources` - (Required) the Origin Server Information. See the following `Block sources`.
* `status` - (Optional, Computed) The status of the resource. Valid values: `offline`, `online`.

#### Block sources

The sources supports the following: 

* `content` - (Required) The Back-to-Source Address.
* `enabled` - (Optional) The source status. Valid values: online, offline.
* `port` - (Required) Port.
* `priority` - (Required) Priority.
* `type` - (Required) The Origin Server Type. Valid Values: 
  * ipaddr: IP Source Station 
  * domain: the Domain Name 
  * oss: OSS Bucket as a Source Station.

#### Block domain_configs

The domain_configs supports the following: 

* `config_id` - (Optional) Configuration ID.
* `function_args` - (Optional) Each Function. See the following `Block function_args`.
* `function_name` - (Optional) Function Name. See [funtion name](https://help.aliyun.com/document_detail/92912.html).
* `status` - (Optional) Configure State, Including success, testing, failed, configuring.

#### Block function_args

The function_args supports the following: 

* `arg_name` - (Optional) The Configuration Name.
* `arg_value` - (Optional) and Leave the Configuration Values.

#### Block cert_infos

The cert_infos supports the following: 

* `cert_name` - (Optional) If You Enable HTTPS Here Certificate Name.
* `cert_type` - (Optional) Certificate Type. Value Range: 
  * upload: Certificate
  * cas: Certificate Authority Certificate. 
  * free: Free Certificate.
* `ssl_pri` - (Optional) Private Key. Do Not Enable Certificate without Entering a User Name and Configure Certificates Enter Private Key.
* `ssl_protocol` - (Optional) Whether to Enable SSL Certificate. Valid Values: on, off. Valid values: `on`, `off`.
* `ssl_pub` - (Optional) If You Enable HTTPS Here Key.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain. Its value is same as `domain_name`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Domain.
* `delete` - (Defaults to 1 mins) Used when delete the Domain.
* `update` - (Defaults to 11 mins) Used when update the Domain.

## Import

SCDN Domain can be imported using the id, e.g.

```
$ terraform import alicloud_scdn_domain.example <domain_name>
```