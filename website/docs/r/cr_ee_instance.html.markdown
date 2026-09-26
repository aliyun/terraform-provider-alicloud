---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instance"
description: |-
  Provides a Alicloud CR Instance resource.
---

# alicloud_cr_ee_instance

Provides a CR Instance resource.


For information about Container Registry Instance and how to use it, see [What is Container Registry](https://www.alibabacloud.com/help/en/acr/product-overview/what-is-container-registry).

For information about CR Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/doc-detail/208144.htm).

-> **NOTE:** Available since v1.124.0.

-> **NOTE:** International site only support `Basic` and `Advance` instance types.

## Example Usage

Basic Usage

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
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `custom_oss_bucket` - (Optional) Custom OSS Bucket name

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `default_oss_bucket` - (Optional, Available since v1.235.0) Whether to use the default OSS Bucket. Value:
true: Use the default OSS Bucket.
false: Use a custom OSS Bucket.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `image_scanner` - (Optional, Available since v1.235.0) The security scan engine used by the Enterprise Edition of Container Image Service. Value:
  - `ACR`: Uses the Trivy scan engine provided by default.
  - `SAS`: uses the enhanced cloud security scan engine.
  - `DISABLE`: Disables the image security scan engine.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `instance_name` - (Required, ForceNew) InstanceName
* `instance_type` - (Required) The Value configuration of the Group 1 attribute of Container Mirror Service Enterprise Edition. Valid values:
Economy: Economy instance
Basic: Basic instance
Standard: Standard instance
Advanced: Advanced Edition Instance

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `kms_encrypted_password` - (Optional, Available since v1.132.0) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.132.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `logistics` - (Optional, ForceNew, Available since v1.295.0) Logistics information. A JSON string describing the delivery address of the instance order, for example: `{"cityCode":"330100","cityName":"Hangzhou","contactName":"Terraform","countryCode":"CN","detailAddress":"No. 1 Example Road","districtCode":"330106","districtName":"Xihu District","email":"tf@example.com","mobilePhone":"15300000000","phone":"1234567","provCode":"330000","provName":"Zhejiang","streetCode":"33010610","streetName":"Zhuantang","zipCode":"310000"}`.

-> **NOTE:** This parameter only takes effect when the instance is created. Changing it forces Terraform to destroy the existing instance and create a new one.

* `namespace_quota` - (Optional, Int, Available since v1.268.0) Additional namespace quota.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `password` - (Optional) Login password, 8-32 digits, must contain at least two letters, symbols, or numbers
* `payment_type` - (Required, ForceNew) Payment type, value:
  - Subscription: Prepaid.
* `period` - (Optional, Int) Prepaid cycle. The unit is Monthly, please enter an integer multiple of 12 for annual paid products.

-> **NOTE:**  must be set when creating a prepaid instance.


-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `renew_period` - (Optional, ForceNew, Int) Automatic renewal cycle, in months.

-> **NOTE:**  When `RenewalStatus` is set to `AutoRenewal`, it must be set.

* `renewal_status` - (Optional, ForceNew, Computed) Automatic renewal status, value:
  - AutoRenewal: automatic renewal.
  - ManualRenewal: manual renewal.

Default ManualRenewal.
* `repo_quota` - (Optional, Int, Available since v1.268.0) Additional repository quota.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `resource_group_id` - (Optional, Computed, Available since v1.235.0) The ID of the resource group
* `tags` - (Optional, Map, Available since v1.288.0) A mapping of tags to assign to the resource.
* `vpc_quota` - (Optional, Int, Available since v1.268.0) VPC access control quota.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.


The following arguments will be discarded. Please use new fields as soon as possible:
* `created_time` - (Deprecated since v1.235.0). Field 'created_time' has been deprecated from provider version 1.235.0. New field 'create_time' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `end_time` - Expiration Time.
* `instance_endpoints` - (Available since v1.240.0) Instance Network Access Endpoint List.
  * `domains` - Domain List.
    * `domain` - Domain.
    * `type` - Domain Type.
  * `enable` - enable.
  * `endpoint_type` - Network Access Endpoint Type.
* `instance_issue` - (Available since v1.295.0) Instance issue.
* `modified_time` - (Available since v1.295.0) Last modification time, as a Unix timestamp in milliseconds.
* `region_id` - RegionId.
* `status` - Instance Status.
* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 6 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

CR Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_instance.example <instance_id>
```