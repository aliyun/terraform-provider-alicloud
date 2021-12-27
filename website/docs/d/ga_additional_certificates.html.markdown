---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_additional_certificates"
sidebar_current: "docs-alicloud-datasource-ga-additional-certificates"
description: |-
  Provides a list of Ga Additional Certificates to the user.
---

# alicloud\_ga\_additional\_certificates

This data source provides the Ga Additional Certificates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.150.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_additional_certificates" "ids" {
  accelerator_id = "example_value"
  listener_id    = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "ga_additional_certificate_id_1" {
  value = data.alicloud_ga_additional_certificates.ids.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Additional Certificate IDs.
* `listener_id` - (Required, ForceNew) The ID of the listener. Only HTTPS listeners support this parameter.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `certificates` - A list of Ga Additional Certificates. Each element contains the following attributes:
	* `accelerator_id` - The ID of the GA instance.
	* `certificate_id` - The Certificate ID.
	* `domain` - The domain name specified by the certificate.
	* `id` - The ID of the Additional Certificate. The value formats as `<accelerator_id>:<listener_id>:<domain>`.
	* `listener_id` - The ID of the listener. Only HTTPS listeners support this parameter.