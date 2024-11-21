---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_security_policy"
sidebar_current: "docs-alicloud-resource-alb-security-policy"
description: |-
  Provides a Alicloud ALB Security Policy resource.
---

# alicloud_alb_security_policy

Provides a ALB Security Policy resource.

For information about ALB Security Policy and how to use it, see [What is Security Policy](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createsecuritypolicy).

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_security_policy&exampleId=bd72c280-b973-e348-4cf6-f3af2cbf9ea1f8c1dff8&activeTab=example&spm=docs.r.alb_security_policy.0.bd72c280b9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_alb_security_policy" "default" {
  security_policy_name = "tf_example"
  tls_versions         = ["TLSv1.0"]
  ciphers              = ["ECDHE-ECDSA-AES128-SHA", "AES256-SHA"]
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Optional) The ID of the resource group.
* `security_policy_name` - (Required) The name of the resource. The name must be 2 to 128 characters in length and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-).
* `dry_run` - (Optional) The dry run.
* `tls_versions` - (Required) The TLS protocol versions that are supported. Valid values: TLSv1.0, TLSv1.1, TLSv1.2 and TLSv1.3.
* `ciphers` - (Required) The supported cipher suites, which are determined by the TLS protocol version.The specified cipher suites must be supported by at least one TLS protocol version that you select. 
* `tags` - (Optional) A mapping of tags to assign to the resource.
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Security Policy.
* `status` - The status of the resource.

## Import

ALB Security Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_security_policy.example <id>
```
