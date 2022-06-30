---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_service_meshes"
sidebar_current: "docs-alicloud-datasource-service-mesh-service-meshes"
description: |-
  Provides a list of Service Mesh Service Meshes to the user.
---

# alicloud\_service\_mesh\_service\_meshes

This data source provides the Service Mesh Service Meshes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_service_mesh_service_meshes" "ids" {
  ids = ["example_id"]
}
output "service_mesh_service_mesh_id_1" {
  value = data.alicloud_service_mesh_service_meshes.ids.meshes.0.id
}

data "alicloud_service_mesh_service_meshes" "nameRegex" {
  name_regex = "^my-ServiceMesh"
}
output "service_mesh_service_mesh_id_2" {
  value = data.alicloud_service_mesh_service_meshes.nameRegex.meshes.0.id
}

data "alicloud_service_mesh_service_meshes" "status" {
  ids    = ["example_id"]
  status = "running"
}
output "service_mesh_service_mesh_id_3" {
  value = data.alicloud_service_mesh_service_meshes.status.meshes.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Service Mesh IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Service Mesh name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `running` or `initial`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Service Mesh names.
* `meshes` - A list of Service Mesh Service Meshes. Each element contains the following attributes:
    * `clusters` - Cluster List.
    * `create_time` - The created time of the resource.
    * `istio_operator_version` - The Istio Operator Version. **Note:** the `istio_operator_version` is available from the version v1.170.0.
    * `sidecar_version` - The SideCar Version. **Note:** the `sidecar_version` is available from the version v1.170.0.
    * `endpoints` - The endpoint details.
        * `intranet_api_server_endpoint` - The internal address of the API Server.
        * `intranet_pilot_endpoint` - The internal address of the Istio Pilot.
        * `public_api_server_endpoint` - The public address of the API Server.
        * `public_pilot_endpoint` - The public address of the Istio Pilot.
    * `error_message` - The Cause of the Error.
    * `edition` - Grid instance version type. Valid values: `Default` and `Pro`. `Default`:the standard. `Pro`:the Pro version.
    * `id` - The ID of the Service Mesh.
    * `load_balancer` - The configuration of the Load Balancer.
        * `api_server_loadbalancer_id` - The IP address of a public network exposed API Server corresponding to the Load Balance.
        * `api_server_public_eip` - Whether to use the IP address of a public network exposed the API Server.
        * `pilot_public_eip` - Whether to use the IP address of a public network exposure the Istio Pilot.
        * `pilot_public_loadbalancer_id` - The IP address of a public network exposure Istio Pilot corresponds to the Load Balance.
    * `mesh_config` - The configuration of the Service grid.
        * `access_log` - The configuration of the access logging.
            * `enabled` - Whether to enable of the access logging. Valid values: `true` and `false`.
            * `project` - The SLS Project of the access logging.
        * `control_plane_log` - The configuration of the control plane logging. **NOTE:** Available in 1.174.0+
            * `enabled` - Whether to enable of the control plane logging. Valid values: `true` and `false`.
            * `project` - The SLS Project of the control plane logging.            
        * `customized_zipkin` - Whether or not to enable the use of a custom zipkin.
        * `enable_locality_lb` - Whether to enable service can access the service through the nearest node access.
        * `tracing` - Whether to enable link trace.
        * `outbound_traffic_policy` - The policy of the Out to the traffic. Valid values: `ALLOW_ANY` and `REGISTRY_ONLY`.
        * `proxy` - The configuration of the Proxy.
            * `limit_memory` - The Memory limit of the resource.
            * `request_cpu` - The  CPU requests of the resources.
            * `request_memory` - The  memory request of the resource.
            * `cluster_domain` - The domain name of the Cluster.
            * `limit_cpu` - The CPU limited of the resource for the proxy container.
        * `telemetry` - Whether to enable acquisition Prometheus metrics.
        * `kiali` - The configuration of the Kiali.
            * `enabled` - Whether to enable kiali, you must first open the collection Prometheus, when the configuration update is false, the system automatically set this value to false).
            * `url` - The service address of the Kiali.
        * `opa` - The open-door policy of agent (OPA) plug-in information.
            * `enabled` - Whether integration-open policy agent (OPA) plug-in.
            * `limit_cpu` - The CPU resource  of the limitsOPA proxy container.
            * `limit_memory` - The memory resource limit of the OPA proxy container.
            * `log_level` - The log level of the OPA proxy container .
            * `request_cpu` - The CPU resource request of the OPA proxy container.
            * `request_memory` - The memory resource request of the OPA proxy container.
        * `pilot` - The configuration of the Link trace sampling.
            * `http10_enabled` - Whether to support the HTTP1.0.
            * `trace_sampling` - The  percentage of the Link trace sampling.
        * `prometheus` - the configuration of the Prometheus.
            * `external_url` - The  service addresses of the Prometheus.
            * `use_external` - Whether to enable external Prometheus.
        * `sidecar_injector` - The configuration of the Sidecar injector.
            * `limit_memory` - The memory limit  of the Sidecar injector Pods.
            * `request_cpu` - The requested cpu the Sidecar injector Pods.
            * `request_memory` - The requested memory the Sidecar injector Pods.
            * `sidecar_injector_webhook_as_yaml` - Other automatic injection Sidecar configuration (in YAML format).
            * `auto_injection_policy_enabled` - Whether to enable by Pod Annotations automatic injection Sidecar.
            * `enable_namespaces_by_default` - Whether it is the all namespaces you turn on the auto injection capabilities.
            * `init_cni_configuration` - The configuration of the CNI
                * `enabled` - Whether to enable CNI.
                * `exclude_namespaces` - The excluded namespace of the CNI.
            * `limit_cpu` - Sidecar injector Pods on the throttle.
        * `audit` - The configuration of the Service grid audit.
            * `enabled` - Whether to enable Service grid audit.
            * `project` - The Service grid audit that to the project.
        * `include_ip_ranges` - The IP ADDRESS range.
    * `network` - The configuration of the Service grid network.
        * `vswitche_list` - The list of Virtual Switch.
        * `vpc_id` - The ID of the VPC.
        * `security_group_id` - The ID of the Security group
    * `service_mesh_id` - The first ID of the resource.
    * `service_mesh_name` - The name of the resource.
    * `status` - The status of the resource.
    * `version` - The version of the resource.
