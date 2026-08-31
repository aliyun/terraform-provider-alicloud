---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_domain_config"
sidebar_current: "docs-alicloud-resource-cdn-domain-config"
description: |-
  Provides a Alicloud Cdn Domain Config resource.
---

# alicloud_cdn_domain_config

Provides a Cdn Domain Config resource.

For information about Cdn Domain Config and how to use it, see [What is Domain Config](https://www.alibabacloud.com/help/en/doc-detail/90915.htm)

-> **NOTE:** Available since v1.34.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cdn_domain_config&exampleId=640bd670-7a13-196b-09ce-ddb572f8f29c9013d7c2&activeTab=example&spm=docs.r.cdn_domain_config.0.640bd6707a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

# Create a new Domain config.
resource "alicloud_cdn_domain_new" "domain" {
  domain_name = "mycdndomain-${random_integer.default.result}.alicloud-provider.cn"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = "20"
    port     = 80
    weight   = "15"
  }
}

resource "alicloud_cdn_domain_config" "config" {
  domain_name   = alicloud_cdn_domain_new.domain.domain_name
  function_name = "ip_allow_list_set"
  function_args {
    arg_name  = "ip_list"
    arg_value = "110.110.110.110"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cdn_domain_config&spm=docs.r.cdn_domain_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the accelerated domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `function_name` - (Required, ForceNew) The name of the domain config.
* `parent_id` - (Optional, Available since v1.220.0) By configuring the function condition (rule engine) in the domain name configuration function parameters, Rule conditions can be created (Rule conditions can match and filter user requests by identifying various parameters carried in user requests). After each rule condition is created, a corresponding ConfigId will be generated, and the ConfigId can be referenced by other functions as a ParentId parameter, in this way, the rule conditions can be combined with the functional configuration to form a more flexible configuration.
* `function_args` - (Required, Set) The args of the domain config. See [`function_args`](#function_args) below.

### `function_args`

The `function_args` block supports the following:

* `arg_name` - (Required) The name of arg.
* `arg_value` - (Required) The value of arg.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain Config. It formats as `<domain_name>:<function_name>:<config_id>`.
-> **NOTE:** Before provider version 1.132.0, it formats as `<domain_name>:<function_name>`
* `config_id` - (Available since v1.132.0) The ID of the domain config function.
* `status` - (Available since v1.132.0) The Status of the function.

## Import

CDN domain config can be imported using the id, e.g.

```shell
terraform import alicloud_cdn_domain_config.example <domain_name>:<function_name>:<config_id>
```

**NOTE:** Before provider version 1.132.0, CDN domain config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_domain_config.example <domain_name>:<function_name>
```
