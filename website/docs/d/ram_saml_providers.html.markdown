---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_saml_providers"
sidebar_current: "docs-alicloud-datasource-ram-saml-providers"
description: |-
  Provides a list of Ram Saml Providers to the user.
---

# alicloud\_ram\_saml\_providers

This data source provides the Ram Saml Providers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.114.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ram_saml_providers" "example" {
  ids        = ["samlProviderName"]
  name_regex = "tf-testAcc"
}

output "first_ram_saml_provider_id" {
  value = data.alicloud_ram_saml_providers.example.providers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of SAML Provider IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by SAML Provider name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of SAML Provider names.
* `providers` - A list of Ram Saml Providers. Each element contains the following attributes:
	* `arn` - The Alibaba Cloud Resource Name (ARN) of the IdP.
	* `description` - The description of SAML Provider.
	* `encodedsaml_metadata_document` - The encodedsaml metadata document.
	* `id` - The ID of the SAML Provider.
	* `saml_provider_name` - The saml provider name.
	* `update_date` - The update time.
