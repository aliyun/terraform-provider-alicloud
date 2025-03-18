---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus"
sidebar_current: "docs-alicloud-datasource-arms-prometheus"
description: |-
  Provides a list of Arms Prometheus to the user.
---

# alicloud_arms_prometheus

This data source provides the Arms Prometheus of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "ecs"
  grafana_instance_id = "free"
  vpc_id              = data.alicloud_vpcs.default.ids.0
  vswitch_id          = data.alicloud_vswitches.default.ids.0
  security_group_id   = alicloud_security_group.default.id
  cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  tags = {
    Created = "TF"
    For     = "Prometheus"
  }
}

data "alicloud_arms_prometheus" "nameRegex" {
  name_regex = "${alicloud_arms_prometheus.default.cluster_name}"
}

output "arms_prometheus_id" {
  value = data.alicloud_arms_prometheus.nameRegex.prometheis.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Prometheus IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Prometheus name.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `enable_details` - (Optional, Bool, Available since v1.214.0) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prometheus names.
* `prometheis` - A list of Prometheus. Each element contains the following attributes:
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
  * `remote_read_intra_url` - (Available since v1.214.0) The internal URL for remote read. **Note:** `remote_read_intra_url` takes effect only if `enable_details` is set to `true`.
  * `remote_read_inter_url` - (Available since v1.214.0) The public URL for remote read. **Note:** `remote_read_inter_url` takes effect only if `enable_details` is set to `true`.
  * `remote_write_intra_url` - (Available since v1.214.0) The internal URL for remote write. **Note:** `remote_write_intra_url` takes effect only if `enable_details` is set to `true`.
  * `remote_write_inter_url` - (Available since v1.214.0) The public URL for remote write. **Note:** `remote_write_inter_url` takes effect only if `enable_details` is set to `true`.
  * `push_gate_way_intra_url` - (Available since v1.214.0) The internal URL for Pushgateway. **Note:** `push_gate_way_intra_url` takes effect only if `enable_details` is set to `true`.
  * `push_gate_way_inter_url` - (Available since v1.214.0) The public URL for Pushgateway. **Note:** `push_gate_way_inter_url` takes effect only if `enable_details` is set to `true`.
  * `http_api_intra_url` - (Available since v1.214.0) The internal URL for the HTTP API. **Note:** `http_api_intra_url` takes effect only if `enable_details` is set to `true`.
  * `http_api_inter_url` - (Available since v1.214.0) The public URL for the HTTP API. **Note:** `http_api_inter_url` takes effect only if `enable_details` is set to `true`.
  * `auth_token` - (Available since v1.214.0) The authorization token. **Note:** `auth_token` takes effect only if `enable_details` is set to `true`.
  * `tags` - The tag of the Prometheus.
  