---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_instance"
description: |-
  Provides a Alicloud Anti-DDoS Pro (DdosBgp) Instance resource.
---

# alicloud_ddosbgp_instance

Provides a Anti-DDoS Pro (DdosBgp) Instance resource.



For information about Anti-DDoS Pro (DdosBgp) Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddosbgp_instance&exampleId=73013b8f-f4af-d04d-5bc2-974159e015748ffa2e4f&activeTab=example&spm=docs.r.ddosbgp_instance.0.73013b8ff4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_ddosbgp_instance" "instance" {
  name             = var.name
  base_bandwidth   = 20
  bandwidth        = -1
  ip_count         = 100
  ip_type          = "IPv4"
  normal_bandwidth = 100
  type             = "Enterprise"
}
```

### Deleting `alicloud_ddosbgp_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ddosbgp_instance`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ddosbgp_instance&spm=docs.r.ddosbgp_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Required, ForceNew, Int) The bandwidth of the package configuration.
* `base_bandwidth` - (Optional, ForceNew) The basic protection bandwidth of the Anti-DDoS Origin Enterprise instance. Default value: `20`. Valid values: `20`.
* `instance_name` - (Optional, Available since v1.259.0) The name of the instance.
* `ip_count` - (Required, ForceNew, Int) The number of IP addresses that can be protected by the Anti-DDoS Origin Enterprise instance.
* `ip_type` - (Required, ForceNew) The protection IP address type of the protection package. Valid values:
  - `IPv4`
  - `IPv6`
* `normal_bandwidth` - (Required, ForceNew, Int) The normal clean bandwidth. Unit: Mbit/s.
* `period` - (Optional) The duration that you will buy Ddosbgp instance (in month). Valid values: [1~9], 12, 24, 36. Default to 12. At present, the provider does not support modify "period".
* `resource_group_id` - (Optional, Available since v1.259.0) Resource Group ID
* `tags` - (Optional, Map, Available since v1.259.0) The key of the tag that is added to the Anti-DDoS Origin instance.
* `type` - (Optional, ForceNew) The protection package type of the DDoS native protection instance. Default value: `Enterprise`. Valid values: `Enterprise`, `Professional`.
* `name` - (Optional, Deprecated since v1.259.0) Field `name` has been deprecated from provider version 1.259.0. New field `instance_name` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - (Available since v1.259.0) The status of the Instance.

## Timeouts

-> **NOTE:** Available since v1.259.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 26 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Anti-DDoS Pro (DdosBgp) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddosbgp_instance.example <id>
```
