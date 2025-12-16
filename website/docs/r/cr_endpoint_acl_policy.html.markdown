---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_endpoint_acl_policy"
sidebar_current: "docs-alicloud-resource-cr-endpoint-acl-policy"
description: |-
  Provides a Alicloud CR Endpoint Acl Policy resource.
---

# alicloud_cr_endpoint_acl_policy

Provides a CR Endpoint Acl Policy resource.

For information about CR Endpoint Acl Policy and how to use it, see [What is Endpoint Acl Policy](https://www.alibabacloud.com/help/doc-detail/145275.htm).

-> **NOTE:** Available since v1.139.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_endpoint_acl_policy&exampleId=94de1ced-0cee-db8d-ae0d-6a2c39086e0ec0b0a975&activeTab=example&spm=docs.r.cr_endpoint_acl_policy.0.94de1ced0c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000000
  max = 99999999
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = alicloud_cr_ee_instance.default.id
  module_name   = "Registry"
}

resource "alicloud_cr_endpoint_acl_policy" "default" {
  instance_id   = data.alicloud_cr_endpoint_acl_service.default.instance_id
  entry         = "192.168.1.0/24"
  description   = var.name
  module_name   = "Registry"
  endpoint_type = "internet"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_endpoint_acl_policy&spm=docs.r.cr_endpoint_acl_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the entry.
* `endpoint_type` - (Required, ForceNew) The type of endpoint. Valid values: `internet`.
* `entry` - (Required, ForceNew) The IP segment that allowed to access.
* `instance_id` - (Required, ForceNew) The ID of the CR Instance.
* `module_name` - (Optional, ForceNew) The module that needs to set the access policy. Valid values: `Registry`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Endpoint Acl Policy. The value formats as `<instance_id>:<endpoint_type>:<entry>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Acl Policy.
* `delete` - (Defaults to 10 mins) Used when delete the Acl Policy.

## Import

CR Endpoint Acl Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_endpoint_acl_policy.example <instance_id>:<endpoint_type>:<entry>
```
