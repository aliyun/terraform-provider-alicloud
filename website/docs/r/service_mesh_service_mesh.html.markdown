---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_mesh"
description: |-
  Provides a Alicloud Service Mesh Service Mesh resource.
---

# alicloud_service_mesh_service_mesh

Provides a Service Mesh Service Mesh resource. 

For information about Service Mesh Service Mesh and how to use it, see [What is Service Mesh](https://www.alibabacloud.com/help/en/asm/developer-reference/api-servicemesh-2020-01-11-createservicemesh).

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

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
  version           = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
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
  version           = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
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
* `cluster_ids` - (Optional, Available since v1.166.0+.) List of clusters.
* `cluster_spec` - (Optional, ForceNew, Computed, Available since v1.166.0+.) Cluster specification. The service mesh instance specification. Valid values: `standard`,`enterprise`,`ultimate`. Default to `standard`.
* `customized_prometheus` - (Optional, Available since v1.211.2) Whether to customize Prometheus. Value:
  -'true': custom Prometheus.
  -'false': Do not customize Prometheus. Default value: 'false '.
* `edition` - (Optional, ForceNew, Computed) Grid instance version type. Valid values: `Default` and `Pro`. Default: the standard. Pro: the Pro version.
* `extra_configuration` - (Optional, Computed) Data plane KubeAPI access capability. See [`extra_configuration`](#extra_configuration) below.
* `force` - (Optional) Whether to forcibly delete the ASM instance. Value:
  -'true': force deletion of ASM instance
  -'false': no forced deletion of ASM instance. Default value: false.
* `load_balancer` - (Optional, ForceNew, Computed) Load balancing information. See [`load_balancer`](#load_balancer) below.
* `mesh_config` - (Optional, ForceNew, Computed) Service grid configuration information. See [`mesh_config`](#mesh_config) below.
* `network` - (Required, ForceNew) Service grid network configuration information. See [`network`](#network) below.
* `prometheus_url` - (Optional, Available since v1.211.2) The Prometheus service address (in non-custom cases, use the ARMS address format).
* `service_mesh_name` - (Optional, ForceNew) ServiceMeshName.
* `tags` - (Optional, Map, Available since v1.211.2) The tag of the resource.
* `version` - (Optional, Computed) Service grid version number. The version of the resource. you can look up the version using alicloud_service_mesh_versions. Note: The version supports updating from v1.170.0, the relevant version can be obtained via istio_operator_version in `alicloud_service_mesh_service_meshes`.

### `extra_configuration`

The extra_configuration supports the following:
* `cr_aggregation_enabled` - (Optional) Whether the data plane KubeAPI access capability is enabled. Indicates whether the Kubernetes API of clusters on the data plane is used to access Istio resources. A value of true indicates that the Kubernetes API is used.

### `load_balancer`

The load_balancer supports the following:
* `api_server_public_eip` - (Optional, ForceNew) Indicates whether to use the IP address of a public network exposed API Server.
* `pilot_public_eip` - (Optional, ForceNew) Indicates whether to use the IP address of a public network exposure Istio Pilot.

### `mesh_config`

The mesh_config supports the following:
* `access_log` - (Optional, Computed) The access logging configuration. See [`access_log`](#mesh_config-access_log) below.
* `audit` - (Optional, ForceNew, Computed) Audit information. See [`audit`](#mesh_config-audit) below.
* `control_plane_log` - (Optional, ForceNew, Computed) Control plane log collection configuration. See [`control_plane_log`](#mesh_config-control_plane_log) below.
* `customized_zipkin` - (Optional, Computed) Whether or not to enable the use of a custom zipkin.
* `enable_locality_lb` - (Optional, ForceNew) Whether to enable service can access the service through the nearest node access.
* `include_ip_ranges` - (Optional, Computed) The IP ADDRESS range.
* `kiali` - (Optional, Computed) Kiali configuration. See [`kiali`](#mesh_config-kiali) below.
* `opa` - (Optional) The open-door policy of agent (OPA) plug-in information. See [`opa`](#mesh_config-opa) below.
* `outbound_traffic_policy` - (Optional) Out to the traffic policy.
* `pilot` - (Optional, ForceNew) Link trace sampling information. See [`pilot`](#mesh_config-pilot) below.
* `proxy` - (Optional) Proxy configuration. See [`proxy`](#mesh_config-proxy) below.
* `sidecar_injector` - (Optional) Sidecar injector configuration. See [`sidecar_injector`](#mesh_config-sidecar_injector) below.
* `telemetry` - (Optional) Whether to enable acquisition Prometheus metrics (it is recommended that you use [Alibaba Cloud Prometheus monitoring](https://arms.console.aliyun.com/).
* `tracing` - (Optional) Whether to enable link trace (you need to have [Alibaba Cloud link tracking service](https://tracing-analysis.console.aliyun.com/).

### `mesh_config-access_log`

The access_log supports the following:
* `enabled` - (Optional, Computed) Whether to enable access log.
* `project` - (Optional, Computed) Access the SLS Project of log collection.

### `mesh_config-audit`

The audit supports the following:
* `enabled` - (Optional, ForceNew, Computed) Enable Audit.
* `project` - (Optional, ForceNew, Computed) Audit Log Items.

### `mesh_config-control_plane_log`

The control_plane_log supports the following:
* `enabled` - (Optional, ForceNew, Computed) Whether to enable control plane log collection. Value:
  -'true': enables control plane log collection.
  -'false': does not enable control plane log collection.
* `project` - (Optional, ForceNew, Computed) The name of the SLS Project to which the control plane logs are collected.

### `mesh_config-kiali`

The kiali supports the following:
* `enabled` - (Optional) Whether to enable kiali, you must first open the collection Prometheus, when the configuration update is false, the system automatically set this value to false).

### `mesh_config-opa`

The opa supports the following:
* `enabled` - (Optional) Whether integration-open policy agent (OPA) plug-in.
* `limit_cpu` - (Optional) OPA proxy container of CPU resource limits.
* `limit_memory` - (Optional) OPA proxy container of the memory resource limit.
* `log_level` - (Optional) OPA proxy container log level.
* `request_cpu` - (Optional) OPA proxy container of CPU resource request.
* `request_memory` - (Optional) OPA proxy container of the memory resource request.

### `mesh_config-pilot`

The pilot supports the following:
* `http10_enabled` - (Optional, Computed) Whether to support the HTTP1.0.
* `trace_sampling` - (Optional) Link trace sampling percentage.

### `mesh_config-proxy`

The proxy supports the following:
* `limit_cpu` - (Optional) CPU resources.
* `limit_memory` - (Optional) Memory limit resource.
* `request_cpu` - (Optional) CPU requests resources.
* `request_memory` - (Optional) A memory request resources.

### `mesh_config-sidecar_injector`

The sidecar_injector supports the following:
* `auto_injection_policy_enabled` - (Optional, Computed) Whether to enable by Pod Annotations automatic injection Sidecar.
* `enable_namespaces_by_default` - (Optional, Computed) Whether it is the all namespaces you turn on the auto injection capabilities.
* `init_cni_configuration` - (Optional, Computed) CNI configuration. See [`init_cni_configuration`](#mesh_config-sidecar_injector-init_cni_configuration) below.
* `limit_cpu` - (Optional, Computed) Sidecar injector Pods on the throttle.
* `limit_memory` - (Optional, Computed) Sidecar injector Pods on the throttle.
* `request_cpu` - (Optional, Computed) Sidecar injector Pods on the requested resource.
* `request_memory` - (Optional, Computed) Sidecar injector Pods on the requested resource.

### `mesh_config-sidecar_injector-init_cni_configuration`

The init_cni_configuration supports the following:
* `enabled` - (Optional) Enable CNI.
* `exclude_namespaces` - (Optional) The excluded namespace.

### `network`

The network supports the following:
* `vswitche_list` - (Required, ForceNew) Virtual Switch ID.
* `vpc_id` - (Required, ForceNew) VPC ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Service grid creation time.
* `load_balancer` - Load balancing information.
  * `api_server_loadbalancer_id` - The IP address of a public network exposed API Server corresponding to the load balancing ID.
  * `pilot_public_loadbalancer_id` - The IP address of a public network exposure Istio Pilot corresponds to the load balancing ID.
* `mesh_config` - Service grid configuration information.
  * `prometheus` - Prometheus configuration.
    * `external_url` - Prometheus service addresses (enabled external Prometheus when the system automatically populates).
    * `use_external` - Whether to enable external Prometheus.
  * `sidecar_injector` - The configuration of the Sidecar injector.
    * `sidecar_injector_webhook_as_yaml` - Other configurations of automatically injected sidecar (in YAML format).
  * `kiali` - Kiali configuration.
    * `url` - Grid topology service address.
  * `proxy` - Proxy configuration.
    * `cluster_domain` - Trust cluster domain.
* `network` - Service grid network configuration information.
  * `security_group_id` - Security group ID.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service Mesh.
* `delete` - (Defaults to 5 mins) Used when delete the Service Mesh.
* `update` - (Defaults to 5 mins) Used when update the Service Mesh.

## Import

Service Mesh Service Mesh can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_mesh_service_mesh.example <id>
```