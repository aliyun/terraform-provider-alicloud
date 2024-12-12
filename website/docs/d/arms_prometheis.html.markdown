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

-> **NOTE:** Available since v1.203.0.

-> **DEPRECATED:** This data source has been renamed to [alicloud_arms_prometheus](https://www.terraform.io/docs/providers/alicloud/d/arms_prometheus) from version 1.214.0.

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

data "alicloud_resource_manager_resource_groups" "default" {
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
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  tags = {
    Created = "TF"
    For     = "Prometheus"
  }
}

data "alicloud_arms_prometheis" "nameRegex" {
  name_regex = "${alicloud_arms_prometheus.default.cluster_name}"
}

output "arms_prometheis_id" {
  value = data.alicloud_arms_prometheis.nameRegex.prometheis.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Prometheus IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Prometheus name.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `enable_details` - (Optional) Whether to query details about the instance.
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
  * `vswitch_id` - The ID of the vSwitch.
  * `security_group_id` - The ID of the security group.
  * `sub_clusters_json` - The child instance json string of the globalView instance.
  * `grafana_instance_id` - The ID of the Grafana workspace.
  * `resource_group_id` - The ID of the resource group.
  * `remote_read_intra_url` - RemoteRead intranet Url.
  * `remote_read_inter_url` - Public Url of remoteRead.
  * `remote_write_intra_url` - RemoteWrite Intranet Url.
  * `remote_write_inter_url` - RemoteWrite public Url.
  * `push_gate_way_intra_url` - PushGateway intranet Url.
  * `push_gate_way_inter_url` - PushGateway public network Url.
  * `http_api_intra_url` - Http api intranet address.
  * `http_api_inter_url` - Http api public network address.
  * `auth_token` - The token used to access the data source.
  * `tags` - The tag of the Prometheus.
  