---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_mesh"
sidebar_current: "docs-alicloud-resource-service-mesh-service-mesh"
description: |-
  Provides a Alicloud Service Mesh Service Mesh resource.
---

# alicloud\_service\_mesh\_service\_mesh

Provides a Service Mesh Service Mesh resource.

For information about Service Mesh Service Mesh and how to use it, see [What is Service Mesh](https://help.aliyun.com/document_detail/171559.html).

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "alicloud_vpc" "default" {
  count    = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  vpc_name = "example_value"
}
data "alicloud_vswitches" "default" {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}
resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = "example_value"
}
resource "alicloud_service_mesh_service_mesh" "example" {
  service_mesh_name = "example_value"
  network {
    vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
    vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
  }
}

```

## Argument Reference

The following arguments are supported:

* `load_balancer` - (Optional, ForceNew) The configuration of the Load Balancer. See the following `Block load_balancer`.
* `mesh_config` - (Optional) The configuration of the Service grid. See the following `Block mesh_config`.
* `network` - (Required, ForceNew) The network configuration of the Service grid. See the following `Block network`.
* `service_mesh_name` - (Optional, ForceNew) The name of the resource.
* `version` - (Optional) The version of the resource. you can look up the version using `alicloud_service_mesh_versions`. **Note:** The `version` supports updating from v1.170.0, the relevant version can be obtained via `istio_operator_version` in `alicloud_service_mesh_service_meshes`.
* `edition` - (Optional, ForceNew) The type  of the resource. Valid values: `Default` and `Pro`. `Default`:the standard. `Pro`:the Pro version.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.
* `cluster_spec` - (Optional,Available in v1.166.0+.) The service mesh instance specification. Valid values: `standard`,`enterprise`,`ultimate`.
* `cluster_ids` - (Optional,Available in v1.166.0+.) The array of the cluster ids.
* `extra_configuration` - (Optional, Available in v1.169.0+.) The configurations of additional features for the ASM instance. See the following `Block extra_configuration`.


#### Block network

The network supports the following: 

* `vpc_id` - (Required) The ID of the VPC.
* `vswitche_list` - (Required) The list of Virtual Switch.

#### Block extra_configuration

The extra_configuration supports the following:

* `cr_aggregation_enabled` - (Optional, Available in v1.169.0+.) Indicates whether the Kubernetes API of clusters on the data plane is used to access Istio resources. A value of `true` indicates that the Kubernetes API is used.


#### Block mesh_config

The mesh_config supports the following: 

* `access_log` - (Optional) The configuration of the access logging.
* `control_plane_log` - (Optional, ForceNew, Available in 1.174.0+) The configuration of the control plane logging. 
* `audit` - (Optional) The configuration of the audit. See the following `Block audit`.
* `customized_zipkin` - (Optional) Whether to enable the use of a custom zipkin.
* `enable_locality_lb` - (Optional) The enable locality lb.
* `kiali` - (Optional) The configuration of the Kiali. See the following `Block kiali`.
* `opa` - (Optional) The open-door policy of agent (OPA) plug-in information. See the following `Block opa`.
* `outbound_traffic_policy` - (Optional) The policy of the Out to the traffic. Valid values: `ALLOW_ANY` and `REGISTRY_ONLY`.
* `pilot` - (Optional, ForceNew) The configuration of the Link trace sampling. See the following `Block pilot`.
* `proxy` - (Optional) The configuration of the Proxy. See the following `Block proxy`.
* `sidecar_injector` - (Optional)The configuration of the Sidecar injector. See the following `Block sidecar_injector`.
* `telemetry` - (Optional) Whether to enable acquisition Prometheus metrics (it is recommended that you use [Alibaba Cloud Prometheus monitoring](https://arms.console.aliyun.com/).
* `tracing` - (Optional) Whether to enable link trace (you need to have [Alibaba Cloud link tracking service](https://tracing-analysis.console.aliyun.com/).

#### Block access_log

The access_log supports the following:

* `enabled` - (Optional, Available in 1.174.0+) Whether to enable of the access logging. Valid values: `true` and `false`.
* `project` - (Optional, Available in 1.174.0+) The SLS Project of the access logging.

#### Block control_plane_log

The control_plane_log supports the following:

* `enabled` - (Optional, Available in 1.174.0+) Whether to enable of the control plane logging. Valid values: `true` and `false`.
* `project` - (Optional, Available in 1.174.0+) The SLS Project of the control plane logging.

#### Block sidecar_injector

The sidecar_injector supports the following: 

* `auto_injection_policy_enabled` - (Optional) Whether to enable by Pod Annotations automatic injection Sidecar.
* `enable_namespaces_by_default` - (Optional) Whether it is the all namespaces you turn on the auto injection capabilities.
* `limit_cpu` - (Optional) The limit cpu of the Sidecar injector Pods.
* `limit_memory` - (Optional) Sidecar injector Pods on the throttle.
* `request_cpu` - (Optional) The requested cpu the Sidecar injector Pods.
* `request_memory` - (Optional) The requested memory the Sidecar injector Pods.

#### Block proxy

The proxy supports the following: 

* `limit_cpu` - (Optional) The limit cpu of the resource.
* `limit_memory` - (Optional) The memory limit of the resource.
* `request_cpu` - (Optional) The request cpu of the resource.
* `request_memory` - (Optional) The request memory of the resource.

#### Block pilot

The pilot supports the following: 

* `http10_enabled` - (Optional) Whether to support the HTTP1.0.
* `trace_sampling` - (Optional) The  percentage of the Link trace sampling.

#### Block kiali

The kiali supports the following: 

* `enabled` - (Optional) Whether to enable kiali, you must first open the collection Prometheus, when the configuration update is false, the system automatically set this value to false.

#### Block opa

The opa supports the following:

* `enabled` - (Optional) Whether integration-open policy agent (OPA) plug-in.
* `limit_cpu` - (Optional) The CPU resource  of the limitsOPA proxy container.
* `limit_memory` - (Optional) The memory resource limit of the OPA proxy container.
* `log_level` - (Optional) The log level of the OPA proxy container .
* `request_cpu` - (Optional) The CPU resource request of the OPA proxy container.
* `request_memory` - (Optional) The memory resource request of the OPA proxy container.

#### Block audit

The audit supports the following: 

* `enabled` - (Optional) Whether to enable Service grid audit.
* `project` - (Optional) The Service grid audit that to the project.

#### Block load_balancer

The load_balancer supports the following: 

* `api_server_public_eip` - (Optional)  Whether to use the IP address of a public network exposed the API Server.
* `pilot_public_eip` - (Optional) Whether to use the IP address of a public network exposure the Istio Pilot.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Service Mesh.
* `status` - The status of the resource. Valid values: `running` or `initial`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Service Mesh.
* `delete` - (Defaults to 5 mins) Used when delete the Service Mesh.
* `update` - (Defaults to 5 mins) Used when update the Service Mesh.

## Import

Service Mesh Service Mesh can be imported using the id, e.g.

```
$ terraform import alicloud_service_mesh_service_mesh.example <id>
```
