---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_meshes"
sidebar_current: "docs-alicloud-datasource-service-mesh-service-meshes"
description: |-
  Provides a list of Service Mesh Service Meshes to the user.
---

# alicloud_service_mesh_service_meshes

This data source provides the Service Mesh Service Meshes of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

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

data "alicloud_service_mesh_service_meshes" "status" {
  ids    = [alicloud_service_mesh_service_mesh.default.id]
  status = "running"
}

output "service_mesh_service_mesh_id_3" {
  value = data.alicloud_service_mesh_service_meshes.status.meshes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List)  A list of Service Mesh IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Service Mesh name.
* `status` - (Optional, ForceNew) The status of the Service Mesh. Valid values: `running`, `initial`.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Service Mesh names.
* `meshes` - A list of Service Mesh Service Meshes. Each element contains the following attributes:
  * `id` - The ID of the Service Mesh instance.
  * `service_mesh_id` - The ID of the Service Mesh instance.
  * `service_mesh_name` - The name of the Service Mesh instance.
  * `clusters` - The Cluster List.
  * `error_message` - The returned error message.
  * `version` - The version of the Service Mesh instance.
  * `status` - The status of the Service Mesh instance.
  * `create_time` - The created time of the Service Mesh instance.
  * `istio_operator_version` - (Available since v1.170.0) The Istio Operator Version. **Note:** `istio_operator_version` takes effect only if `enable_details` is set to `true`.
  * `sidecar_version` - (Available since v1.170.0) The SideCar Version. **Note:** `sidecar_version` takes effect only if `enable_details` is set to `true`.
  * `kube_config` - The content of Kube config. **Note:** `kube_config` takes effect only if `enable_details` is set to `true`.
  * `edition` - Grid instance version type. **Note:** `edition` takes effect only if `enable_details` is set to `true`.
  * `endpoints` - The endpoint details. **Note:** `endpoints` takes effect only if `enable_details` is set to `true`.
    * `intranet_api_server_endpoint` - The internal address of the API Server.
    * `intranet_pilot_endpoint` - The internal address of the Istio Pilot.
    * `public_api_server_endpoint` - The public address of the API Server.
    * `public_pilot_endpoint` - The public address of the Istio Pilot.
  * `load_balancer` - The configuration of the Load Balancer. **Note:** `load_balancer` takes effect only if `enable_details` is set to `true`.
    * `api_server_loadbalancer_id` - The IP address of a public network exposed API Server corresponding to the Load Balance.
    * `api_server_public_eip` - Whether to use the IP address of a public network exposed the API Server.
    * `pilot_public_eip` - Whether to use the IP address of a public network exposure the Istio Pilot.
    * `pilot_public_loadbalancer_id` - The IP address of a public network exposure Istio Pilot corresponds to the Load Balance.
  * `network` - The configuration of the Service grid network. **Note:** `network` takes effect only if `enable_details` is set to `true`.
    * `vpc_id` - The ID of the VPC.
    * `vswitche_list` - The list of Virtual Switch.
    * `security_group_id` - The ID of the Security group
  * `mesh_config` - The configuration of the Service grid. **Note:** `mesh_config` takes effect only if `enable_details` is set to `true`.
    * `telemetry` - Whether to enable acquisition Prometheus metrics.
    * `customized_zipkin` - Whether or not to enable the use of a custom zipkin.
    * `enable_locality_lb` - Whether to enable service can access the service through the nearest node access.
    * `tracing` - Whether to enable link trace.
    * `outbound_traffic_policy` - The policy of the Out to the traffic.
    * `include_ip_ranges` - The IP ADDRESS range.
    * `access_log` - The configuration of the access logging.
      * `enabled` - Whether to enable of the access logging.
      * `project` - The SLS Project of the access logging.
    * `control_plane_log` - (Available since v1.174.0) The configuration of the control plane logging.
      * `enabled` - Whether to enable of the control plane logging.
      * `project` - The SLS Project of the control plane logging.
    * `audit` - The configuration of the Service grid audit.
      * `enabled` - Whether to enable Service grid audit.
      * `project` - The Service grid audit that to the project.
    * `kiali` - The configuration of the Kiali.
      * `enabled` - Whether to enable kiali.
      * `url` - The service address of the Kiali.
    * `pilot` - The configuration of the Link trace sampling.
      * `http10_enabled` - Whether to support the HTTP1.0.
      * `trace_sampling` - The  percentage of the Link trace sampling.    
    * `prometheus` - the configuration of the Prometheus.
      * `external_url` - The  service addresses of the Prometheus.
      * `use_external` - Whether to enable external Prometheus.
    * `opa` - The open-door policy of agent (OPA) plug-in information.
      * `enabled` - Whether integration-open policy agent (OPA) plug-in.
      * `limit_cpu` - The CPU resource  of the limitsOPA proxy container.
      * `limit_memory` - The memory resource limit of the OPA proxy container.
      * `log_level` - The log level of the OPA proxy container .
      * `request_cpu` - The CPU resource request of the OPA proxy container.
      * `request_memory` - The memory resource request of the OPA proxy container.
    * `proxy` - The configuration of the Proxy.
      * `limit_memory` - The Memory limit of the resource.
      * `request_cpu` - The  CPU requests of the resources.
      * `request_memory` - The  memory request of the resource.
      * `cluster_domain` - The domain name of the Cluster.
      * `limit_cpu` - The CPU limited of the resource for the proxy container.
    * `sidecar_injector` - The configuration of the Sidecar injector.
      * `limit_memory` - The memory limit  of the Sidecar injector Pods.
      * `request_cpu` - The requested cpu the Sidecar injector Pods.
      * `request_memory` - The requested memory the Sidecar injector Pods.
      * `sidecar_injector_webhook_as_yaml` - Other automatic injection Sidecar configuration (in YAML format).
      * `auto_injection_policy_enabled` - Whether to enable by Pod Annotations automatic injection Sidecar.
      * `enable_namespaces_by_default` - Whether it is the all namespaces you turn on the auto injection capabilities.
      * `limit_cpu` - Sidecar injector Pods on the throttle.
      * `init_cni_configuration` - The configuration of the CNI
        * `enabled` - Whether to enable CNI.
        * `exclude_namespaces` - The excluded namespace of the CNI.