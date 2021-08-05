---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_config_maps"
sidebar_current: "docs-alicloud-datasource-sae-config-maps"
description: |-
  Provides a list of Sae Config Maps to the user.
---

# alicloud\_sae\_config\_maps

This data source provides the Sae Config Maps of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_sae_config_maps" "nameRegex" {
  namespace_id = "cn-hangzhou:namespaceId"
  name_regex   = "^my-ConfigMap"
}
output "sae_config_map_id_2" {
  value = data.alicloud_sae_config_maps.nameRegex.maps.0.id
}
            
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Config Map IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Config Map name.
* `namespace_id` - (Required, ForceNew) The NamespaceId of Config Maps.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Config Map names.
* `maps` - A list of Sae Config Maps. Each element contains the following attributes:
	* `config_map_id` - The first ID of the resource.
	* `create_time` - The Creation Time of the ConfigMap.
	* `data` - ConfigMap instance data. The value's format is a `json` string
	* `description` - The Description of Config Map.
	* `id` - The ID of the Config Map.
	* `name` - ConfigMap instance name.
	* `namespace_id` - The NamespaceId of Config Maps.