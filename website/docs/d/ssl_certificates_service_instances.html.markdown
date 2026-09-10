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
* `brand` - (ForceNew, Optional) The certificate brand. Valid values: WoSign, CFCA, DigiCert, GeoTrust, GlobalSign, vTrus, and Alibaba.
* `certificate_status` - (ForceNew, Optional) The status of the certificate.
  - `issued`: Issued
  - `revoked`: Revoked
  - `willExpire`: Expiring Soon
  - `expired`: Expired
* `certificate_type` - (ForceNew, Optional) The type of the certificate. Valid values: DV, OV, and EV.
* `instance_type` - (ForceNew, Optional) The instance type. Valid values: BUY: official certificate; TEST: test certificate.
* `keyword` - (ForceNew, Optional) Fuzzy search that matches domain names, instance names, or corresponding resource IDs.
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group.
* `status` - (ForceNew, Optional) The instance status.
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
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `auto_reissue` - **NOTE:** This field is only available when `enable_details` is `true`. Specifies whether to enable managed renewal.
  * `average_waiting_time` - **NOTE:** This field is only available when `enable_details` is `true`. Average waiting time for issuing a certificate of this specification, in seconds.
  * `brand` - **NOTE:** This field is only available when `enable_details` is `true`. The certificate brand.
  * `certificate_type` - **NOTE:** This field is only available when `enable_details` is `true`. The type of the certificate.
  * `city` - **NOTE:** This field is only available when `enable_details` is `true`. The city where the company or organization to which the certificate purchaser belongs is located.
  * `company_id` - **NOTE:** This field is only available when `enable_details` is `true`. The company information ID.
  * `contact_id_list` - **NOTE:** This field is only available when `enable_details` is `true`. The list of contact IDs.
  * `country_code` - **NOTE:** This field is only available when `enable_details` is `true`. The code of the country or region where the certificate organization is located.
  * `csr` - **NOTE:** This field is only available when `enable_details` is `true`. The CSR content.
  * `domain` - **NOTE:** This field is only available when `enable_details` is `true`. The domain names to be bound to the certificate.
  * `full_domain_count` - **NOTE:** This field is only available when `enable_details` is `true`. The number of single domain names included in the instance.
  * `generate_csr_method` - **NOTE:** This field is only available when `enable_details` is `true`. The method used to generate the Certificate Signing Request (CSR).
  * `instance_end_time` - **NOTE:** This field is only available when `enable_details` is `true`. Instance expiration time, a UNIX timestamp in seconds.
  * `instance_id` - The instance ID.
  * `instance_name` - **NOTE:** This field is only available when `enable_details` is `true`. The name of the instance.
  * `instance_start_time` - **NOTE:** This field is only available when `enable_details` is `true`. Instance start time, a UNIX timestamp in seconds.
  * `instance_type` - **NOTE:** This field is only available when `enable_details` is `true`. The instance type.
  * `key_algorithm` - **NOTE:** This field is only available when `enable_details` is `true`. The certificate algorithm.
  * `order_end_time` - **NOTE:** This field is only available when `enable_details` is `true`. Order end time, a UNIX timestamp in seconds.
  * `order_start_time` - **NOTE:** This field is only available when `enable_details` is `true`. Order start time, a UNIX timestamp in seconds.
  * `province` - **NOTE:** This field is only available when `enable_details` is `true`. The province or region where the company is located.
  * `resource_group_id` - **NOTE:** This field is only available when `enable_details` is `true`. The ID of the resource group.
  * `spec` - **NOTE:** This field is only available when `enable_details` is `true`. The specification of the purchased instance.
  * `status` - **NOTE:** This field is only available when `enable_details` is `true`. The instance status.
  * `upgrade_status` - **NOTE:** This field is only available when `enable_details` is `true`. The upgrade status of the instance.
  * `validation_method` - **NOTE:** This field is only available when `enable_details` is `true`. The verification method for the certificate application.
  * `wildcard_domain_count` - **NOTE:** This field is only available when `enable_details` is `true`. The number of wildcard domain names included in the instance.
  * `tags` - **NOTE:** This field is only available when `enable_details` is `true`. The tags of the instance.
  * `id` - The ID of the resource supplied above.
