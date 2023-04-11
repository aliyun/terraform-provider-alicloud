---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheis"
sidebar_current: "docs-alicloud-datasource-arms-prometheis"
description: |-
  Provides a list of Arms Prometheis to the user.
---

# alicloud\_arms\_prometheis

This data source provides the Arms Prometheis of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.203.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_prometheis" "ids" {
  ids = ["example_id"]
}

output "arms_prometheis_id_1" {
  value = data.alicloud_arms_prometheis.ids.prometheis.0.id
}

data "alicloud_arms_prometheis" "nameRegex" {
  name_regex = "tf-example"
}

output "arms_prometheis_id_2" {
  value = data.alicloud_arms_prometheis.nameRegex.prometheis.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Prometheus IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Prometheus name.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prometheus names.
* `prometheis` - A list of Prometheis. Each element contains the following attributes:
  * `id` - The ID of the Prometheus.
  * `cluster_id` - The ID of the cluster.
  * `cluster_type` - The type of the cluster.
  * `cluster_name` - The name of the cluster.
  * `vpc_id` - The ID of the VPC.
  * `vswitch_id` - The ID of the VSwitch.
  * `security_group_id` - The ID of the security group.
  * `sub_clusters_json` - The child instance json string of the globalView instance.
  * `grafana_instance_id` - The ID of the Grafana workspace.
  * `resource_group_id` - The ID of the resource group.
  * `tags` - The tag of the Prometheus.
  