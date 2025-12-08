---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_storage_domain_routing_rule"
description: |-
  Provides a Alicloud CR Storage Domain Routing Rule resource.
---

# alicloud_cr_storage_domain_routing_rule

Provides a CR Storage Domain Routing Rule resource.

Instance Storage Domain Routing Rule.

For information about CR Storage Domain Routing Rule and how to use it, see [What is Storage Domain Routing Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateStorageDomainRoutingRule).

-> **NOTE:** Available since v1.265.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_storage_domain_routing_rule&exampleId=c3f442f1-b7d2-e324-d664-1b610d027b130e11347a&activeTab=example&spm=docs.r.cr_storage_domain_routing_rule.0.c3f442f1b7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}

resource "alicloud_cr_storage_domain_routing_rule" "default" {
  routes {
    instance_domain = "${alicloud_cr_ee_instance.default.instance_name}-registry-vpc.cn-hangzhou.cr.aliyuncs.com"
    storage_domain  = "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou-internal.aliyuncs.com"
    endpoint_type   = "Internet"
  }
  instance_id = alicloud_cr_ee_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the Container Registry Instance.
* `routes` - (Required, List) Domain name routing entry See [`routes`](#routes) below.

### `routes`

The routes supports the following:
* `endpoint_type` - (Required) Endpoint Type.
* `instance_domain` - (Required) Instance domain name.
* `storage_domain` - (Required) Storage domain name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<rule_id>`.
* `create_time` - The creation time of the resource.
* `rule_id` - The ID of the Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Storage Domain Routing Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Storage Domain Routing Rule.
* `update` - (Defaults to 5 mins) Used when update the Storage Domain Routing Rule.

## Import

CR Storage Domain Routing Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_storage_domain_routing_rule.example <instance_id>:<rule_id>
```
