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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_mesh_service_mesh&exampleId=7a493e32-9016-a2c2-381a-2a54cf1aad5348f76c48&activeTab=example&spm=docs.r.service_mesh_service_mesh.0.7a493e3290&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

Basic Usage with mesh config
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_mesh_service_mesh&exampleId=0892d3e5-fcbd-c733-99cc-f4abe500c5db1471b0f9&activeTab=example&spm=docs.r.service_mesh_service_mesh.1.0892d3e5fc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf-example"
}
data "alicloud_service_mesh_versions" "default" {
  edition = "Default"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.3.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = "tf-test-service_mesh"
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = var.name
}

data "alicloud_instance_types" "default" {
  availability_zone    = alicloud_vswitch.default.zone_id
  kubernetes_node_role = "Worker"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_ecs_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "created by terraform"
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
  mesh_config {
    customized_zipkin  = false
    enable_locality_lb = false
    telemetry          = true
    kiali {
      enabled = true
    }

    tracing = true
    pilot {
      http10_enabled = true
      trace_sampling = 100
    }
    opa {
      enabled        = true
      log_level      = "info"
      request_cpu    = "1"
      request_memory = "512Mi"
      limit_cpu      = "2"
      limit_memory   = "1024Mi"
    }
    audit {
      enabled = true
      project = alicloud_log_project.default.project_name
    }
    proxy {
      request_memory = "128Mi"
      limit_memory   = "1024Mi"
      request_cpu    = "100m"
      limit_cpu      = "2000m"
    }
    sidecar_injector {
      enable_namespaces_by_default  = false
      request_memory                = "128Mi"
      limit_memory                  = "1024Mi"
      request_cpu                   = "100m"
      auto_injection_policy_enabled = true
      limit_cpu                     = "2000m"
    }
    outbound_traffic_policy = "ALLOW_ANY"
    access_log {
      enabled         = true
      gateway_enabled = true
      sidecar_enabled = true
    }
  }
  cluster_ids = [alicloud_cs_kubernetes_node_pool.default.cluster_id]
  extra_configuration {
    cr_aggregation_enabled = true
  }
}
```

## Argument Reference

The following arguments are supported:
* `cluster_ids` - (Optional, List, Available since v1.166.0) List of clusters.
* `cluster_spec` - (Optional, Computed, Available since v1.166.0) Cluster specification
* `customized_prometheus` - (Optional, Available since v1.211.2) Whether to customize Prometheus. Value:
  -'true': custom Prometheus.
  -'false': Do not customize Prometheus.

  Default value: 'false '.
* `edition` - (Optional, ForceNew) Grid instance version type (for example: the standard, the Pro version, etc.)
* `extra_configuration` - (Optional, Computed, List) Data plane KubeAPI access capability See [`extra_configuration`](#extra_configuration) below.
* `force` - (Optional) Whether to forcibly delete the ASM instance. Value:
  -'true': force deletion of ASM instance
  -'false': no forced deletion of ASM instance

  Default value: false
* `load_balancer` - (Optional, Computed, List) Load balancing information See [`load_balancer`](#load_balancer) below.
* `mesh_config` - (Optional, Computed, List) Service grid configuration information See [`mesh_config`](#mesh_config) below.
* `network` - (Required, ForceNew, List) Service grid network configuration information See [`network`](#network) below.
* `prometheus_url` - (Optional, Available since v1.211.2) The Prometheus service address (in non-custom cases, use the ARMS address format).
* `service_mesh_name` - (Optional) ServiceMeshName
* `tags` - (Optional, Map, Available since v1.211.2) The tag of the resource
* `version` - (Optional) Service grid version number

### `extra_configuration`

The extra_configuration supports the following:
* `cr_aggregation_enabled` - (Optional) Whether the data plane KubeAPI access capability is enabled.

### `load_balancer`

The load_balancer supports the following:
* `api_server_public_eip` - (Optional, ForceNew) Indicates whether to use the IP address of a public network exposed API Server
* `pilot_public_eip` - (Optional, ForceNew) Indicates whether to use the IP address of a public network exposure Istio Pilot. **Note**: This field has been deprecated and is readonly as of 1.232.0. Use pilot_public_eip_id instead.
* `pilot_public_eip_id` - (Optional, Available since v1.232.0) the EIP instance id of Pilot load balancer.

### `mesh_config`

The mesh_config supports the following:
* `access_log` - (Optional, List) The access logging configuration See [`access_log`](#mesh_config-access_log) below.
* `audit` - (Optional, Computed, List) Audit information See [`audit`](#mesh_config-audit) below.
* `control_plane_log` - (Optional, List) Control plane log collection configuration. See [`control_plane_log`](#mesh_config-control_plane_log) below.
* `customized_zipkin` - (Optional) Whether or not to enable the use of a custom zipkin
* `enable_locality_lb` - (Optional, ForceNew) Whether to enable service can access the service through the nearest node access
* `include_ip_ranges` - (Optional, Computed) The IP ADDRESS range
* `kiali` - (Optional, List) Kiali configuration See [`kiali`](#mesh_config-kiali) below.
* `opa` - (Optional, List) The open-door policy of agent (OPA) plug-in information See [`opa`](#mesh_config-opa) below.
* `outbound_traffic_policy` - (Optional) Out to the traffic policy
* `pilot` - (Optional, List) Link trace sampling information See [`pilot`](#mesh_config-pilot) below.
* `proxy` - (Optional, List) Proxy configuration, the fields under this structure have service segment default values, if not explicitly specified, you need to manually add them based on the return value of the server after the instance is created. See [`proxy`](#mesh_config-proxy) below.
* `sidecar_injector` - (Optional, List) Sidecar injector configuration See [`sidecar_injector`](#mesh_config-sidecar_injector) below.
* `telemetry` - (Optional) Whether to enable acquisition Prometheus metrics (it is recommended that you use [Alibaba Cloud Prometheus monitoring](https://arms.console.aliyun.com/)
* `tracing` - (Optional) Whether to enable link trace (you need to have [Alibaba Cloud link tracking service](https://tracing-analysis.console.aliyun.com/)

### `mesh_config-access_log`

The mesh_config-access_log supports the following:
* `enabled` - (Optional) Whether to enable access log
* `gateway_enabled` - (Optional, Available since v1.223.1) Whether collect AccessLog of ASM Gateway to Alibaba Cloud SLS
* `gateway_lifecycle` - (Optional, Computed, Int, Available since v1.223.1) Lifecycle of AccessLog of ASM Gateways which have been collected to Alibaba Cloud SLS
* `project` - (Optional) Access the SLS Project of log collection.
* `sidecar_enabled` - (Optional, Available since v1.223.1) Whether collect AccessLog of ASM Gateway to Alibaba Cloud SLS
* `sidecar_lifecycle` - (Optional, Computed, Int, Available since v1.223.1) Lifecycle of AccessLog of ASM Sidecars which have been collected to Alibaba Cloud SLS

### `mesh_config-audit`

The mesh_config-audit supports the following:
* `enabled` - (Optional, Computed) Enable Audit
* `project` - (Optional, Computed) Audit Log Items

### `mesh_config-control_plane_log`

The mesh_config-control_plane_log supports the following:
* `enabled` - (Required) Whether to enable control plane log collection. Value:
  -'true': enables control plane log collection.
  -'false': does not enable control plane log collection.
* `log_ttl_in_day` - (Optional, Computed, Int, Available since v1.223.1) Lifecycle of logs has been collected to Alibaba Cloud SLS
* `project` - (Optional) The name of the SLS Project to which the control plane logs are collected.

### `mesh_config-kiali`

The mesh_config-kiali supports the following:
* `auth_strategy` - (Optional, Computed, Available since v1.232.0) The authentication strategy used when logging into the mesh topology. In data plane deployment mode, the mesh topology can use token, openid, or ramoauth authentication strategies; in managed mode, the mesh topology can use openid or ramoauth authentication strategies.
* `custom_prometheus_url` - (Optional, Available since v1.232.0) When the mesh topology cannot automatically use the integrated ARMS Prometheus, you need to use this property to specify a custom Prometheus HTTP API Url. The corresponding Prometheus instance needs to have been configured to collect Istio metrics in the cluster within the service mesh.
* `enabled` - (Optional) Whether to enable kiali, you must first open the collection Prometheus, when the configuration update is false, the system automatically set this value to false)
* `integrate_clb` - (Optional, Available since v1.232.0) Whether to integrate CLB for mesh topology services to provide external access.
* `kiali_arms_auth_tokens` - (Optional, Available since v1.232.0) When the mesh topology automatically uses the integrated ARMS Prometheus, if the ARMS Prometheus instance in the cluster has token authentication enabled, you need to use this property to provide the corresponding authentication token for the mesh topology. The key of the property is the Kubernetes cluster id, and the value is the authentication token of the ARMS Prometheus instance corresponding to the cluster. (Service mesh instance version 1.15.3.113 or above is required)
* `kiali_service_annotations` - (Optional, Available since v1.232.0) Annotations for the Service corresponding to the mesh topology service. When the mesh topology service integrates CLB, annotations can be used to control the CLB specifications. The attribute type is map, the key is the Kubernetes cluster id, and the value is the mesh topology service annotation map under the corresponding Kubernetes cluster. When using the managed mode mesh topology, the key is the service mesh instance id. For annotation content, refer to [Configuring traditional load balancing CLB through Annotation](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/user-guide/use-annotations-to-configure-load-balancing).(Service mesh instance version 1.17.2.19 or above is required)
* `open_id_config` - (Optional, List, Available since v1.232.0) When the mesh topology's authentication policy is openid, the configuration used when the mesh topology and OIDC application are connected. If the authentication policy is openid, this configuration must be provided. See [`open_id_config`](#mesh_config-kiali-open_id_config) below.
* `ram_oauth_config` - (Optional, List, Available since v1.232.0) When the authentication strategy of the mesh topology is ramoauth, the mesh topology will be connected to the RAM OAuth application to log in with the Alibaba Cloud account. In this case, this attribute must be provided to configure the connection with the RAM OAuth application. See [`ram_oauth_config`](#mesh_config-kiali-ram_oauth_config) below.
* `server_config` - (Optional, List, Available since v1.232.0) When you need to configure external access to the mesh topology through ASM gateway or other means, and access the mesh topology through a custom domain name or address, you need to specify this property. (The service mesh instance version must be 1.16.4.5 or above) See [`server_config`](#mesh_config-kiali-server_config) below.

### `mesh_config-opa`

The mesh_config-opa supports the following:
* `enabled` - (Optional) Whether integration-open policy agent (OPA) plug-in
* `limit_cpu` - (Optional) OPA proxy container of CPU resource limits
* `limit_memory` - (Optional) OPA proxy container of the memory resource limit
* `log_level` - (Optional) OPA proxy container log level
* `request_cpu` - (Optional) OPA proxy container of CPU resource request
* `request_memory` - (Optional) OPA proxy container of the memory resource request

### `mesh_config-pilot`

The mesh_config-pilot supports the following:
* `http10_enabled` - (Optional) Whether to support the HTTP1.0
* `trace_sampling` - (Optional, Float) Link trace sampling percentage

### `mesh_config-proxy`

The mesh_config-proxy supports the following:
* `cluster_domain` - (Optional, ForceNew, Computed) Cluster domain name
* `limit_cpu` - (Optional) CPU resources
* `limit_memory` - (Optional) Memory limit resource
* `request_cpu` - (Optional) CPU requests resources
* `request_memory` - (Optional) A memory request resources

### `mesh_config-sidecar_injector`

The mesh_config-sidecar_injector supports the following:
* `auto_injection_policy_enabled` - (Optional) Whether to enable by Pod Annotations automatic injection Sidecar
* `enable_namespaces_by_default` - (Optional) Whether it is the all namespaces you turn on the auto injection capabilities
* `init_cni_configuration` - (Optional, Computed, List) CNI configuration See [`init_cni_configuration`](#mesh_config-sidecar_injector-init_cni_configuration) below.
* `limit_cpu` - (Optional) Sidecar injector Pods on the throttle
* `limit_memory` - (Optional) Sidecar injector Pods on the throttle
* `request_cpu` - (Optional) Sidecar injector Pods on the requested resource
* `request_memory` - (Optional) Sidecar injector Pods on the requested resource

### `mesh_config-sidecar_injector-init_cni_configuration`

The mesh_config-sidecar_injector-init_cni_configuration supports the following:
* `enabled` - (Optional) Enable CNI
* `exclude_namespaces` - (Optional) The excluded namespace

### `mesh_config-kiali-open_id_config`

The mesh_config-kiali-open_id_config supports the following:
* `client_id` - (Optional, Available since v1.232.0) The client id provided by the OIDC application
* `client_secret` - (Optional, Available since v1.232.0) The client secret provided by the OIDC application
* `issuer_uri` - (Optional, Available since v1.232.0) OIDC应用的Issuer URI
* `scopes` - (Optional, List, Available since v1.232.0) The scope of the mesh topology request to the OIDC application

### `mesh_config-kiali-ram_oauth_config`

The mesh_config-kiali-ram_oauth_config supports the following:
* `redirect_uris` - (Optional, Available since v1.232.0) The redirect Uri provided to the RAM OAuth application. This needs to be the access address of the mesh topology service. When not provided, the redirect Uri will be automatically inferred based on the ServerConfig or the CLB address of the mesh topology integration.

### `mesh_config-kiali-server_config`

The mesh_config-kiali-server_config supports the following:
* `web_fqdn` - (Optional, Available since v1.232.0) The domain name or address used when accessing the mesh topology in a custom way
* `web_port` - (Optional, Int, Available since v1.232.0) The port used when accessing the mesh topology in a custom way
* `web_root` - (Optional, Available since v1.232.0) The root path of the service when accessing the mesh topology in a custom way
* `web_schema` - (Optional, Available since v1.232.0) The protocol used when accessing the mesh topology in a custom way. Can only be http or https

### `network`

The network supports the following:
* `vswitche_list` - (Required, ForceNew, List) Virtual Switch ID
* `vpc_id` - (Required, ForceNew) VPC ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Service grid creation time
* `kubeconfig` - Kubeconfig configuration content
* `load_balancer` - Load balancing information
  * `api_server_loadbalancer_id` - The Instance ID of APIServer Load Balancer
  * `pilot_public_loadbalancer_id` - The Instance ID of Pilot Load Balancer
* `mesh_config` - Service grid configuration information
  * `kiali` - Kiali configuration
    * `aggregated_kiali_address` - When the mesh topology is deployed in managed mode and integrated with CLB to provide external access, the external access address is automatically generated.
    * `distributed_kiali_access_tokens` - The login token provided when the mesh topology is deployed in data plane deployment mode. When the mesh topology authentication strategy is token, this token can be used to log in to the mesh topology service. The key of the property is the Kubernetes cluster id, and the value of the property is the login token of the mesh topology service in the cluster.
    * `distributed_kiali_addresses` - When the mesh topology is deployed in data plane deployment mode and integrated with CLB to provide external access, the external access address is automatically generated. The key of the attribute is the Kubernetes cluster id, and the value is the external access address of the mesh topology service in the cluster.
    * `url` - Kiali service address
    * `use_populated_arms_prometheus` - Whether the mesh topology automatically uses the integrated ARMS Prometheus. When the integrated ARMS Prometheus is automatically used, there is no need to specify the dependent Prometheus HTTP API Url.
  * `prometheus` - Prometheus configuration
    * `external_url` - Prometheus service addresses (enabled external Prometheus when the system automatically populates)
    * `use_external` - Whether to enable external Prometheus
  * `sidecar_injector` - Sidecar injector configuration
    * `sidecar_injector_webhook_as_yaml` - Other automatic injection Sidecar configuration (in YAML format)
* `network` - Service grid network configuration information
  * `security_group_id` - Security group ID
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service Mesh.
* `delete` - (Defaults to 20 mins) Used when delete the Service Mesh.
* `update` - (Defaults to 10 mins) Used when update the Service Mesh.

## Import

Service Mesh Service Mesh can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_mesh_service_mesh.example <id>
```
