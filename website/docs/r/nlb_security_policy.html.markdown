---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_security_policy"
description: |-
  Provides a Alicloud NLB Security Policy resource.
---

# alicloud_nlb_security_policy

Provides a NLB Security Policy resource.



For information about NLB Security Policy and how to use it, see [What is Security Policy](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createsecuritypolicy-nlb).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_security_policy&exampleId=25030d7e-01a3-a404-96ee-61a95508d8bd2c8d0138&activeTab=example&spm=docs.r.nlb_security_policy.0.25030d7e01&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_nlb_security_policy" "default" {
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_policy_name = var.name
  ciphers              = ["ECDHE-RSA-AES128-SHA", "ECDHE-ECDSA-AES128-SHA"]
  tls_versions         = ["TLSv1.0", "TLSv1.1", "TLSv1.2"]
  tags = {
    Created = "TF"
    For     = "example"
  }
}
```

## Argument Reference

The following arguments are supported:
* `ciphers` - (Required, List) The supported cipher suites, which are determined by the TLS protocol version. You can specify at most 32 cipher suites.
  - TLS 1.0 and TLS 1.1 support the following cipher suites: `ECDHE-ECDSA-AES128-SHA`, `ECDHE-ECDSA-AES256-SHA`, `ECDHE-RSA-AES128-SHA`, `ECDHE-RSA-AES256-SHA`, `AES128-SHA`, `AES256-SHA`, `DES-CBC3-SHA`
  - TLS 1.2 supports the following cipher suites: `ECDHE-ECDSA-AES128-SHA`, `ECDHE-ECDSA-AES256-SHA`, `ECDHE-RSA-AES128-SHA`, `ECDHE-RSA-AES256-SHA`, `AES128-SHA`, `AES256-SHA, DES-CBC3-SHA`, `ECDHE-ECDSA-AES128-GCM-SHA256`, `ECDHE-ECDSA-AES256-GCM-SHA384`, `ECDHE-ECDSA-AES128-SHA256`, `ECDHE-ECDSA-AES256-SHA384`, `ECDHE-RSA-AES128-GCM-SHA256`, `ECDHE-RSA-AES256-GCM-SHA384`, `ECDHE-RSA-AES128-SHA256`, `ECDHE-RSA-AES256-SHA384`, `AES128-GCM-SHA256`, `AES256-GCM-SHA384`, `AES128-SHA256`, `AES256-SHA256`
  - TLS 1.3 supports the following cipher suites: `TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`, `TLS_AES_128_CCM_SHA256`, `TLS_AES_128_CCM_8_SHA256`.
* `resource_group_id` - (Optional, Computed) The ID of the new resource group.

  You can log on to the [Resource Management console](https://resourcemanager.console.aliyun.com/resource-groups) to view resource group IDs.
* `security_policy_name` - (Optional) The name of the security policy.

  The name must be 1 to 200 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-).
* `tags` - (Optional, Map) The tag of the resource
* `tls_versions` - (Required, List) The supported versions of the Transport Layer Security (TLS) protocol. Valid values: `TLSv1.0`, `TLSv1.1`, `TLSv1.2`, and `TLSv1.3`. You can specify at most four TLS versions.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Security Policy.
* `update` - (Defaults to 5 mins) Used when update the Security Policy.

## Import

NLB Security Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_security_policy.example <id>
```