---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_instances"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-instances"
description: |-
  Provides a list of Ssl Certificates Service Instance owned by an Alibaba Cloud account.
---

# alicloud_ssl_certificates_service_instances

This data source provides Ssl Certificates Service Instance available to the user.[What is Instance](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_ssl_certificates_service_instance" "default" {
  product_type  = "cas"
  period        = 12
  pricing_cycle = 2
  instance_name = var.name

  parameter {
    code  = "fullSpec"
    value = "ws.dv.f"
  }
  parameter {
    code  = "fullDomainCount"
    value = "1"
  }
}

data "alicloud_ssl_certificates_service_instances" "default" {
  ids = ["${alicloud_ssl_certificates_service_instance.default.id}"]
}

output "alicloud_ssl_certificates_service_instance_example_id" {
  value = data.alicloud_ssl_certificates_service_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `brand` - (Optional) The certificate brand. Valid values: WoSign, CFCA, DigiCert, GeoTrust, GlobalSign, vTrus, and Alibaba.
* `certificate_status` - (Optional) The status of the certificate.
  - `issued`: Issued
  - `revoked`: Revoked
  - `willExpire`: Expiring Soon
  - `expired`: Expired
* `certificate_type` - (Optional) The type of the certificate. Valid values: DV, OV, and EV.
* `instance_type` - (Optional) The instance type. Valid values: BUY: official certificate; TEST: test certificate.
* `keyword` - (Optional) Fuzzy search that matches domain names, instance names, or corresponding resource IDs.
* `resource_group_id` - (Optional) The ID of the resource group.
* `status` - (Optional) The instance status.
  - `inactive`: Pending use
  - `pending`: Under review. The latest certificate has been submitted for review.
  - `willExpire`: The instance is about to expire.
  - `expired`: The instance has expired.
  - `refund`: Refunded
  - `normal`: Normal
  - `closed`: Disabled and unavailable
* `ids` - (Optional, Computed) A list of Instance IDs.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Instance IDs.
* `instances` - A list of Instance Entries.

---

### `instances`

The instances supports the following attributes:

* `auto_reissue` - Specifies whether to enable managed renewal.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `average_waiting_time` - Average waiting time for issuing a certificate of this specification, in seconds.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `brand` - The certificate brand.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `certificate_type` - The type of the certificate.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `city` - The city where the company or organization to which the certificate purchaser belongs is located.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `company_id` - The company information ID.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `contact_id_list` - The list of contact IDs.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `country_code` - The code of the country or region where the certificate organization is located.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `csr` - The CSR content.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `domain` - The domain names to be bound to the certificate.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `full_domain_count` - The number of single domain names included in the instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `generate_csr_method` - The method used to generate the Certificate Signing Request (CSR).

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `instance_end_time` - Instance expiration time, a UNIX timestamp in seconds.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `instance_id` - The instance ID.
* `instance_name` - The name of the instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `instance_start_time` - Instance start time, a UNIX timestamp in seconds.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `instance_type` - The instance type.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `key_algorithm` - The certificate algorithm.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `order_end_time` - Order end time, a UNIX timestamp in seconds.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `order_start_time` - Order start time, a UNIX timestamp in seconds.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `province` - The province or region where the company is located.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `resource_group_id` - The ID of the resource group.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `spec` - The specification of the purchased instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `status` - The instance status.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `upgrade_status` - The upgrade status of the instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `validation_method` - The verification method for the certificate application.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `wildcard_domain_count` - The number of wildcard domain names included in the instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `tags` - The tags of the instance.

  -> **NOTE:** This field is only available when `enable_details` is `true`.
* `id` - The ID of the resource supplied above.
