---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_namespaces"
sidebar_current: "docs-alicloud-datasource-sae-namespaces"
description: |-
  Provides a list of Sae Namespaces to the user.
---

# alicloud\_sae\_namespaces

This data source provides the Sae Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sae_namespaces" "nameRegex" {
  name_regex = "^my-Namespace"
}
output "sae_namespace_id" {
  value = data.alicloud_sae_namespaces.nameRegex.namespaces.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Namespace IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Namespace name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Namespace names.
* `namespaces` - A list of Sae Namespaces. Each element contains the following attributes:
	* `id` - The ID of the Namespace.
	* `namespace_description` - The Description of Namespace.
	* `namespace_id` - The Id of Namespace.It can contain 2 to 32 characters.The value is in format {RegionId}:{namespace}.
	* `namespace_name` - The Name of Namespace.
