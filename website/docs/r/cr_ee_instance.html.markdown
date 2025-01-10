---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instance"
description: |-
  Provides a Alicloud CR Instance resource.
---

# alicloud_cr_ee_instance

Provides a CR Instance resource.

For information about Container Registry Enterprise Edition instances and how to use it, see [Create a Instance](https://www.alibabacloud.com/help/en/doc-detail/208144.htm)

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_ee_instance&exampleId=22381149-4704-d4bc-dbb5-cafc5c4bdc8e6d63f70f&activeTab=example&spm=docs.r.cr_ee_instance.0.2238114947&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000000
  max = 99999999
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `custom_oss_bucket` - (Optional) Custom OSS Bucket name
* `default_oss_bucket` - (Optional, Available since v1.235.0) Whether to use the default OSS Bucket. Value:
  - `true`: Use the default OSS Bucket.
  - `false`: Use a custom OSS Bucket.
* `image_scanner` - (Optional, Available since v1.235.0) The security scan engine used by the Enterprise Edition of Container Image Service. Value:
  - `ACR`: Uses the Trivy scan engine provided by default.
  - `SAS`: uses the enhanced cloud security scan engine.
* `instance_name` - (Required, ForceNew) InstanceName
* `instance_type` - (Required) The Value configuration of the Group 1 attribute of Container Mirror Service Enterprise Edition. Valid values:
  - `Basic`: Basic instance
  - `Standard`: Standard instance
  - `Advanced`: Advanced Edition Instance
* `password` - (Optional) Login password, 8-32 digits, must contain at least two letters, symbols, or numbers
* `payment_type` - (Required, ForceNew) Payment type, value:
  - Subscription: Prepaid.
* `period` - (Optional, Int) Prepaid cycle. The unit is Monthly, please enter an integer multiple of 12 for annual paid products.

-> **NOTE:**  must be set when creating a prepaid instance.

* `renew_period` - (Optional, ForceNew, Int) Automatic renewal cycle, in months.

-> **NOTE:**  When `RenewalStatus` is set to `AutoRenewal`, it must be set.

* `renewal_status` - (Optional, ForceNew, Computed) Automatic renewal status, value:
  - AutoRenewal: automatic renewal.
  - ManualRenewal: manual renewal.

  Default ManualRenewal.
* `resource_group_id` - (Optional, Computed, Available since v1.235.0) The ID of the resource group

The following arguments will be discarded. Please use new fields as soon as possible:
* `created_time` - (Deprecated since v1.235.0). Field 'created_time' has been deprecated from provider version 1.235.0. New field 'create_time' instead.
* `kms_encrypted_password` - (Optional, Available since v1.132.0) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.132.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `end_time` - Expiration Time
* `instance_endpoints` - (Available since v1.240.0) Instance Network Access Endpoint List
  * `domains` - Domain List
    * `domain` - Domain
    * `type` - Domain Type
  * `enable` - enable
  * `endpoint_type` - Network Access Endpoint Type
* `region_id` - RegionId
* `status` - Instance Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

CR Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_instance.example <id>
```