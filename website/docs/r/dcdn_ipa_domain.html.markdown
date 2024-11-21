---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_ipa_domain"
sidebar_current: "docs-alicloud-resource-dcdn-ipa-domain"
description: |-
  Provides a Alicloud DCDN Ipa Domain resource.
---

# alicloud_dcdn_ipa_domain

Provides a DCDN Ipa Domain resource.

For information about DCDN Ipa Domain and how to use it, see [What is Ipa Domain](https://www.alibabacloud.com/help/en/doc-detail/130634.html).

-> **NOTE:** Available since v1.158.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_ipa_domain&exampleId=068383b5-4eac-3449-d83b-89fd6655a82f6da386a6&activeTab=example&spm=docs.r.dcdn_ipa_domain.0.068383b54e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_dcdn_ipa_domain" "example" {
  domain_name       = "example-${random_integer.default.result}.com"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  scope             = "overseas"
  status            = "online"
  sources {
    content  = "www.alicloud-provider.cn"
    port     = 8898
    priority = "20"
    type     = "domain"
    weight   = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The domain name to be added to IPA. Wildcard domain names are supported. A wildcard domain name must start with a period (.).
* `resource_group_id` - (Optional) The ID of the resource group. If you do not set this parameter, the system automatically assigns the ID of the default resource group.
* `scope` - (Optional, ForceNew) The accelerated region. Valid values: `domestic`, `global`, `overseas`.
* `sources` - (Required) Sources. See [`sources`](#sources) below.
* `status` - (Optional) The status of DCDN Ipa Domain. Valid values: `online`, `offline`. Default to `online`.

### `sources`

The sources supports the following: 

* `content` - (Required) The address of the origin server. You can specify an IP address or a domain name.
* `port` - (Required) The custom port number. Valid values: `0` to `65535`.
* `priority` - (Required) The priority of the origin server. Valid values: `20` and `30`. Default value: `20`. A value of 20 specifies that the origin is a primary origin. A value of 30 specifies that the origin is a secondary origin.
* `type` - (Required) The type of the origin server. Valid values: `ipaddr`, `domain`, `oss`.
* `weight` - (Required) The weight of the origin server. You must specify a value that is less than `100`. Default value: `10`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipa Domain. Its value is same as `domain_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when Creating DCDN Ipa domain instance.
* `update` - (Defaults to 5 mins) Used when Creating DCDN Ipa domain instance.
* `delete` - (Defaults to 10 mins) Used when terminating the DCDN Ipa domain instance.


## Import

DCDN Ipa Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_ipa_domain.example <domain_name>
```