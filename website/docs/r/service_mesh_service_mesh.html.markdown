---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_mesh"
sidebar_current: "docs-alicloud-resource-service-mesh-service-mesh"
description: |-
  Provides a Alicloud Service Mesh Service Mesh resource.
---

# alicloud_service_mesh_service_mesh

Provides a Service Mesh Service Mesh resource.

For information about Service Mesh Service Mesh and how to use it, see [What is Service Mesh](https://www.alibabacloud.com/help/en/alibaba-cloud-service-mesh/latest/api-servicemesh-2020-01-11-createservicemesh).

-> **NOTE:** Available since v1.138.0.

## Example Usage

creating a standard cluster
```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf_example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_service_mesh_service_mesh" "default" {
  service_mesh_name = var.name
  edition           = "Default"
  version           = data.alicloud_service_mesh_versions.default.versions.0.version
  cluster_spec      = "standard"
  network {
    vpc_id        = alicloud_vpc.default.id
    vswitche_list = [alicloud_vswitch.default.id]
  }
  load_balancer {
    pilot_public_eip      = false
    api_server_public_eip = false
  }
}
```

creating an enterprise cluster
```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_service_mesh_service_mesh" "default" {
  service_mesh_name = var.name
  edition           = "Pro"
  version           = data.alicloud_service_mesh_versions.default.versions.0.version
  cluster_spec      = "enterprise"
  network {
    vpc_id        = alicloud_vpc.default.id
    vswitche_list = [alicloud_vswitch.default.id]
  }
  load_balancer {
    pilot_public_eip      = false
    api_server_public_eip = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer` - (Optional, ForceNew) The configuration of the Load Balancer. See [`load_balancer`](#load_balancer) below.
* `mesh_config` - (Optional) The configuration of the Service grid. See [`mesh_config`](#mesh_config) below.
* `network` - (Required, ForceNew) The network configuration of the Service grid. See [`network`](#network) below.
* `service_mesh_name` - (Optional, ForceNew) The name of the resource.
* `version` - (Optional) The version of the resource. you can look up the version using `alicloud_service_mesh_versions`. **Note:** The `version` supports updating from v1.170.0, the relevant version can be obtained via `istio_operator_version` in `alicloud_service_mesh_service_meshes`.
* `edition` - (Optional, ForceNew) The type  of the resource. Valid values: `Default` and `Pro`. `Default`: the standard. `Pro`: the Pro version.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.
* `cluster_spec` - (Optional, Available since v1.166.0) The service mesh instance specification. 
  Valid values: `standard`,`enterprise`,`ultimate`. Default to `standard`.
* `cluster_ids` - (Optional, Available since v1.166.0) The array of the cluster ids.
* `extra_configuration` - (Optional, Available since v1.169.0) The configurations of additional features for the ASM instance. See [`extra_configuration`](#extra_configuration) below.


### `network`

The network supports the following: 

* `vpc_id` - (Required) The ID of the VPC.
* `vswitche_list` - (Required) The list of Virtual Switch.

### `extra_configuration`

The extra_configuration supports the following:

* `cr_aggregation_enabled` - (Optional, Available since v1.169.0) Indicates whether the Kubernetes API of clusters on the data plane is used to access Istio resources. A value of `true` indicates that the Kubernetes API is used.


### `mesh_config`

The mesh_config supports the following: 

* `access_log` - (Optional) The configuration of the access logging. See [`access_log`](#mesh_config-access_log) below.
* `control_plane_log` - (Optional, ForceNew, Available since v1.174.0) The configuration of the control plane logging. See [`control_plane_log`](#mesh_config-control_plane_log) below.
* `audit` - (Optional) The configuration of the audit. See [`audit`](#mesh_config-audit) below.
* `customized_zipkin` - (Optional) Whether to enable the use of a custom zipkin.
* `enable_locality_lb` - (Optional) The enable locality lb.
* `kiali` - (Optional) The configuration of the Kiali. See [`kiali`](#mesh_config-kiali) below.
* `opa` - (Optional) The open-door policy of agent (OPA) plug-in information. See [`opa`](#mesh_config-opa) below.
* `outbound_traffic_policy` - (Optional) The policy of the Out to the traffic. Valid values: `ALLOW_ANY` and `REGISTRY_ONLY`.
* `pilot` - (Optional, ForceNew) The configuration of the Link trace sampling. See [`pilot`](#mesh_config-pilot) below.
* `proxy` - (Optional) The configuration of the Proxy. See [`proxy`](#mesh_config-proxy) below.
* `sidecar_injector` - (Optional)The configuration of the Sidecar injector. See [`sidecar_injector`](#mesh_config-sidecar_injector) below.
* `telemetry` - (Optional) Whether to enable acquisition Prometheus metrics it is recommended that you use [Alibaba Cloud Prometheus monitoring](https://arms.console.aliyun.com/).
* `tracing` - (Optional) Whether to enable link trace you need to have [Alibaba Cloud link tracking service](https://tracing-analysis.console.aliyun.com/).

### `mesh_config-access_log`

The access_log supports the following:

* `enabled` - (Optional, Available since v1.174.0) Whether to enable of the access logging. Valid values: `true` and `false`.
* `project` - (Optional, Available since v1.174.0) The SLS Project of the access logging.

### `mesh_config-control_plane_log`

The control_plane_log supports the following:

* `enabled` - (Optional, Available since v1.174.0) Whether to enable of the control plane logging. Valid values: `true` and `false`.
* `project` - (Optional, Available since v1.174.0) The SLS Project of the control plane logging.

### `mesh_config-sidecar_injector`

The sidecar_injector supports the following: 

* `auto_injection_policy_enabled` - (Optional) Whether to enable by Pod Annotations automatic injection Sidecar.
* `enable_namespaces_by_default` - (Optional) Whether it is the all namespaces you turn on the auto injection capabilities.
* `limit_cpu` - (Optional) The limit cpu of the Sidecar injector Pods.
* `limit_memory` - (Optional) Sidecar injector Pods on the throttle.
* `request_cpu` - (Optional) The requested cpu the Sidecar injector Pods.
* `request_memory` - (Optional) The requested memory the Sidecar injector Pods.

### `mesh_config-proxy`

The proxy supports the following: 

* `limit_cpu` - (Optional) The limit cpu of the resource.
* `limit_memory` - (Optional) The memory limit of the resource.
* `request_cpu` - (Optional) The request cpu of the resource.
* `request_memory` - (Optional) The request memory of the resource.

### `mesh_config-pilot`

The pilot supports the following: 

* `http10_enabled` - (Optional) Whether to support the HTTP1.0.
* `trace_sampling` - (Optional) The  percentage of the Link trace sampling.

### `mesh_config-kiali`

The kiali supports the following: 

* `enabled` - (Optional) Whether to enable kiali, you must first open the collection Prometheus, when the configuration update is false, the system automatically set this value to false.

### `mesh_config-opa`

The opa supports the following:

* `enabled` - (Optional) Whether integration-open policy agent (OPA) plug-in.
* `limit_cpu` - (Optional) The CPU resource  of the limitsOPA proxy container.
* `limit_memory` - (Optional) The memory resource limit of the OPA proxy container.
* `log_level` - (Optional) The log level of the OPA proxy container .
* `request_cpu` - (Optional) The CPU resource request of the OPA proxy container.
* `request_memory` - (Optional) The memory resource request of the OPA proxy container.

### `mesh_config-audit`

The audit supports the following: 

* `enabled` - (Optional) Whether to enable Service grid audit.
* `project` - (Optional) The Service grid audit that to the project.

### `load_balancer`

The load_balancer supports the following: 

* `api_server_public_eip` - (Optional)  Whether to use the IP address of a public network exposed the API Server.
* `pilot_public_eip` - (Optional) Whether to use the IP address of a public network exposure the Istio Pilot.
* `pilot_public_loadbalancer_id` - (Optional) The ID of the Server Load Balancer (SLB) instance that is used when Istio Pilot is exposed to the Internet.
* `api_server_loadbalancer_id` - (Optional)  The ID of the SLB instance that is used when the API server is exposed to the Internet.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Service Mesh.
* `status` - The status of the resource. Valid values: `running` or `initial`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Service Mesh.
* `update` - (Defaults to 5 mins) Used when update the Service Mesh.
* `delete` - (Defaults to 10 mins) Used when delete the Service Mesh.

## Import

Service Mesh Service Mesh can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_mesh_service_mesh.example <id>
```
