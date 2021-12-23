---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_saml_provider"
sidebar_current: "docs-alicloud-resource-ram-saml-provider"
description: |-
  Provides a Alicloud RAM SAML Provider resource.
---

# alicloud\_ram\_saml\_provider

Provides a RAM SAML Provider resource.

For information about RAM SAML Provider and how to use it, see [What is SAML Provider](https://www.alibabacloud.com/help/doc-detail/186846.htm).

-> **NOTE:** Available in v1.114.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_saml_provider" "example" {
  saml_provider_name            = "tf-testAcc"
  encodedsaml_metadata_document = "your encodedsaml metadata document"
  description                   = "For Terraform Test"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of SAML Provider.
* `encodedsaml_metadata_document` - (Optional) The metadata file, which is Base64 encoded. The file is provided by an IdP that supports SAML 2.0.
* `saml_provider_name` - (Required, ForceNew) The name of SAML Provider.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of SAML Provider. Value as `saml_provider_name`.
* `arn` - The Alibaba Cloud Resource Name (ARN) of the IdP.
* `update_date` - The update time.

## Import

RAM SAML Provider can be imported using the id, e.g.

```
$ terraform import alicloud_ram_saml_provider.example <saml_provider_name>
```
