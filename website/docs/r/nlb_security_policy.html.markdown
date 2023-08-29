---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_security_policy"
description: |-
  Provides a Alicloud NLB Security Policy resource.
---

# alicloud_nlb_security_policy

Provides a NLB Security Policy resource.

For information about NLB Security Policy and how to use it, see [What is Security Policy](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-nlb-2022-04-30-createsecuritypolicy).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "sgp1"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "rg2" {
  display_name        = "sgp2"
  resource_group_name = "${var.name}1"
}


resource "alicloud_nlb_security_policy" "default" {
  security_policy_name = var.name
  ciphers              = ["ECDHE-ECDSA-AES128-SHA"]
  tls_versions         = ["TLSv1.0"]
  resource_group_id    = alicloud_resource_manager_resource_group.rg1.id
}
```

## Argument Reference

The following arguments are supported:
* `ciphers` - (Required) The supported cipher suites, which are determined by the TLS protocol version. You can specify at most 32 cipher suites.
  - TLS 1.0 and TLS 1.1 support the following cipher suites: `ECDHE-ECDSA-AES128-SHA`, `ECDHE-ECDSA-AES256-SHA`, `ECDHE-RSA-AES128-SHA`, `ECDHE-RSA-AES256-SHA`, `AES128-SHA`, `AES256-SHA`, `DES-CBC3-SHA`.
  - TLS 1.2 supports the following cipher suites: `ECDHE-ECDSA-AES128-SHA`, `ECDHE-ECDSA-AES256-SHA`, `ECDHE-RSA-AES128-SHA`, `ECDHE-RSA-AES256-SHA`, `AES128-SHA`, `AES256-SHA`, `DES-CBC3-SHA`,` ECDHE-ECDSA-AES128-GCM-SHA256`, `ECDHE-ECDSA-AES256-GCM-SHA384`, `ECDHE-ECDSA-AES128-SHA256`, `ECDHE-ECDSA-AES256-SHA384`, `ECDHE-RSA-AES128-GCM-SHA256`, `ECDHE-RSA-AES256-GCM-SHA384`, `ECDHE-RSA-AES128-SHA256`, `ECDHE-RSA-AES256-SHA384`, `AES128-GCM-SHA256`, `AES256-GCM-SHA384`, `AES128-SHA256`, `AES256-SHA256`.
  - TLS 1.3 supports the following cipher suites: `TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`, `TLS_AES_128_CCM_SHA256`, `TLS_AES_128_CCM_8_SHA256`.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_policy_name` - (Optional) The name of the security policy. The name must be 1 to 200 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-).
* `tags` - (Optional, Map) The tag of the resource.
* `tls_versions` - (Required) The TLS protocol versions that are supported. Valid values: `TLSv1.0`, `TLSv1.1`, `TLSv1.2` and `TLSv1.3`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Security Policy.
* `update` - (Defaults to 5 mins) Used when update the Security Policy.

## Import

NLB Security Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_security_policy.example <id>
```