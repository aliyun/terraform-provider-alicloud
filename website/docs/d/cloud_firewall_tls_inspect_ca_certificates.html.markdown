---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_tls_inspect_ca_certificates"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-tls-inspect-ca-certificates"
description: |-
  Provides a list of Cloud Firewall Tls Inspect Ca Certificate owned by an Alibaba Cloud account.
---

# alicloud_cloud_firewall_tls_inspect_ca_certificates

This data source provides Cloud Firewall Tls Inspect Ca Certificate available to the user.[What is Tls Inspect Ca Certificate](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/GetTlsInspectCertificateDownloadUrl)

-> **NOTE:** Available since v1.262.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_firewall_tls_inspect_ca_certificate" "default" {
}

data "alicloud_cloud_firewall_tls_inspect_ca_certificates" "default" {
  ids = ["${alicloud_cloud_firewall_tls_inspect_ca_certificate.default.id}"]
}

output "alicloud_cloud_firewall_tls_inspect_ca_certificate_example_id" {
  value = data.alicloud_cloud_firewall_tls_inspect_ca_certificates.default.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ca_cert_id` - (ForceNew, Optional) CA certificate ID
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `ids` - (Optional, ForceNew, Computed) A list of Tls Inspect Ca Certificate IDs. 
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Tls Inspect Ca Certificate IDs.
* `certificates` - A list of Tls Inspect Ca Certificate Entries. Each element contains the following attributes:
  * `ca_cert_id` - CA certificate ID
  * `id` - The ID of the resource supplied above.
