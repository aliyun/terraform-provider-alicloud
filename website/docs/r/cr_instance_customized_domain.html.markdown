---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_instance_customized_domain"
description: |-
  Provides a Alicloud CR Instance Customized Domain resource.
---

# alicloud_cr_instance_customized_domain

Provides a CR Instance Customized Domain resource.


For information about CR Instance Customized Domain and how to use it, see [What is Instance Customized Domain](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateInstanceCustomizedDomain).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_ssl_certificates_service_certificates" "default" {
  keyword = "alicloud-provider.cn"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}

resource "alicloud_cr_instance_customized_domain" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  module_name    = "Registry"
  domain         = "alicloud-provider.cn"
  cert_id        = data.alicloud_ssl_certificates_service_certificates.default.certificates.0.id
  cert_region_id = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The instance ID.
* `module_name` - (Required, ForceNew) The custom module name.
* `domain` - (Required, ForceNew) The custom domain name.
* `cert_id` - (Required) The ID of the custom domain name certificate.
* `cert_region_id` - (Optional) The region to which the certificate belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance Customized Domain. It formats as `<instance_id>:<module_name>:<domain>`.
* `region_id` - The region ID.
* `create_time` - The creation time.
* `modified_time` - The modification time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance Customized Domain.
* `update` - (Defaults to 5 mins) Used when update the Instance Customized Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Instance Customized Domain.

## Import

CR Instance Customized Domain can be imported using the id, which consists of instance_id, module_name and domain, e.g.

```shell
$ terraform import alicloud_cr_instance_customized_domain.example <instance_id>:<module_name>:<domain>
```
