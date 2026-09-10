---
subcategory: "SCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_scdn_doamin_config"
sidebar_current: "docs-alicloud-resource-scdn-domain-config"
description: |-
  Provides a Alicloud Scdn domain config Resource.
---

# alicloud_scdn_domain_config

Provides a SCDN Accelerated Domain resource.

For information about domain config and how to use it, see [Batch set config](https://help.aliyun.com/document_detail/92912.html)

-> **NOTE:** Available in v1.131.0+.

-> **NOTE:** Alibaba Cloud SCDN has stopped new customer purchases from January 26, 2023, and you can choose to buy Alibaba Cloud DCDN products with more comprehensive acceleration and protection capabilities. If you are already a SCDN customer, you can submit a work order at any time to apply for a smooth migration to Alibaba Cloud DCDN products. In the future, we will provide better acceleration and security protection services in Alibaba Cloud DCDN, thank you for your understanding and cooperation.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.aliyun.com/product/scdn) from version `1.219.0`.

## Example Usage

Basic Usage

```terraform
# Create a new Domain config.
resource "alicloud_scdn_domain" "domain" {
  domain_name = "mydomain.alicloud-provider.cn"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = "20"
    port     = 80
  }
}

resource "alicloud_scdn_domain_config" "config" {
  domain_name   = alicloud_scdn_domain.domain.domain_name
  function_name = "ip_allow_list_set"
  function_args {
    arg_name  = "ip_list"
    arg_value = "110.110.110.110"
  }
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_scdn_domain_config&spm=docs.r.scdn_domain_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `function_name` - (Required, ForceNew) The name of the domain config.
* `function_args` - (Required, Type: list) The args of the domain config.

### Block function_args

The `function_args` block supports the following:

* `arg_name` - (Required) The name of arg.
* `arg_value` - (Required) The value of arg.

## Attributes Reference

The following attributes are exported:

* `config_id` - The SCDN domain config id.
* `id` - The ID of this resource. The value is formate as `<domain_name>:<function_name>:<config_id>`.
* `status` -  The status of this resource.

## Import

SCDN domain config can be imported using the id, e.g.

```
terraform import alicloud_scdn_domain_config.example <domain_name>:<function_name>:<config_id>
```
