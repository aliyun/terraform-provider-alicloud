---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_k8s_application"
sidebar_current: "docs-alicloud-resource-edas-k8s-application"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alicloud_edas_k8s_application

Create an EDAS k8s application.For information about EDAS K8s Application and how to use it, see [What is EDAS K8s Application](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-insertk8sapplication). 

-> **NOTE:** Available since v1.105.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_k8s_application&exampleId=2e94df1e-5028-caf3-ff5a-270fcf613afcf1425e1c&activeTab=example&spm=docs.r.edas_k8s_application.0.2e94df1e50&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  desired_size         = 2
}
resource "alicloud_edas_k8s_cluster" "default" {
  cs_cluster_id = alicloud_cs_kubernetes_node_pool.default.cluster_id
}

resource "alicloud_edas_k8s_application" "default" {
  application_name        = var.name
  cluster_id              = alicloud_edas_k8s_cluster.default.id
  package_type            = "FatJar"
  package_url             = "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/SPRING_CLOUD_PROVIDER.jar"
  jdk                     = "Open JDK 8"
  replicas                = 2
  readiness               = "{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"port\":18081}}"
  liveness                = "{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"port\":18081}}"
  application_descriotion = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_k8s_application&spm=docs.r.edas_k8s_application.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, ForceNew) The name of the application you want to create. Must start with character,supports numbers, letters and dashes (-), supports up to 36 characters
* `cluster_id` - (Required, ForceNew) The ID of the alicloud container service kubernetes cluster that you want to import to. You can call the ListCluster operation to query.
* `package_type` - (Optional, ForceNew) Application package type. Optional parameter values include: FatJar, WAR and Image.
* `replicas` - (Optional) Number of application instances.
* `image_url` - (Optional) Mirror address. When the package_type is set to 'Image', this parameter item is required.
* `application_descriotion` - (Optional) The description of the application
* `package_url` - (Optional) The url of the package to deploy.Applications deployed through FatJar or WAR packages need to configure it.
* `package_version` - (Optional) The version number of the deployment package. WAR and FatJar types are required. Please customize its meaning.
* `jdk` - (Optional) The JDK version that the deployed package depends on. The optional parameter values are Open JDK 7 and Open JDK 8. Image does not support this parameter.
* `web_container` - (Optional) The Tomcat version that the deployment package depends on. Applicable to Spring Cloud and Dubbo applications deployed through WAR packages. Image does not support this parameter.
* `edas_container_version` - (Optional) EDAS-Container version that the deployed package depends on. Image does not support this parameter.
* `internet_target_port` - (Optional, ForceNew, Deprecated since v1.194.0) The private SLB back-end port, is also the service port of the application, ranging from 1 to 65535.
  It has been deprecated, and new resource 'alicloud_edas_k8s_slb_attachment' replaces it.
* `internet_slb_port` - (Optional, ForceNew, Deprecated since v1.194.0) The public network SLB front-end port, range 1~65535. It has been deprecated and new resource 'alicloud_edas_k8s_slb_attachment' replaces it.
* `internet_slb_protocol` - (Optional, ForceNew, Deprecated since v1.194.0) The public network SLB protocol supports TCP, HTTP and HTTPS protocols. It has been deprecated, and new resource 'alicloud_edas_k8s_slb_attachment' replaces it.
* `internet_slb_id` - (Optional, ForceNew, Deprecated since v1.194.0) Public network SLB ID. If not configured, EDAS will automatically purchase a new SLB for the user.
  It has been deprecated, and new resource 'alicloud_edas_k8s_slb_attachment' replaces it.
* `limit_mem` - (Optional) The memory limit of the application instance during application operation, unit: M.
* `requests_mem` - (Optional) When the application is created, the memory limit of the application instance, unit: M. When set to 0, it means unlimited. 
* `requests_m_cpu` - (Optional) When the application is created, the CPU quota of the application instance, unit: number of millcores, similar to request_cpu
* `limit_m_cpu` - (Optional)  The CPU quota of the application instance during application operation. Unit: Number of millcores, set to 0 means unlimited, similar to request_cpu.
* `command` - (Optional) The set command, if set, will replace the startup command in the mirror when the mirror is started.
* `command_args` - (Optional) Used in combination with the command, the parameter of the command is a JsonArray string in the format: `[{"argument":"-c"},{"argument":"test"}]`. Among them, -c and test are two parameters that need to be set. 
* `envs` - (Optional)  Deployment environment variables, the format must conform to the JSON object array, such as: `{"name":"x","value":"y"},{"name":"x2","value":"y2"}`, If you want to cancel the configuration, you need to set an empty JSON array "" to indicate no configuration.
* `pre_stop` - (Optional) Execute script before stopping
* `post_start` - (Optional) Execute script after startup
* `liveness` - (Optional) Container survival status monitoring, format such as: `{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1,"tcpSocket":{"host":"", "port":8080} }`.
* `readiness` - (Optional) Container service status check. If the check fails, the traffic passing through K8s Service will not be transferred to the container. The format is: `{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1, "httpGet": {"path": "/consumer","port": 8080,"scheme": "HTTP","httpHeaders": [{"name": "test","value": "testvalue"} ]}}`.
* `nas_id` - (Optional) The ID of the mounted NAS must be in the same region as the cluster. It must have an available mount point creation quota, or its mount point must be on a switch in the VPC. If it is not filled in and the mountDescs field exists, a NAS will be automatically purchased and mounted on the switch in the VPC by default.
* `mount_descs` - (Optional) Mount configuration description, as a serialized JSON. For example: `[{"nasPath": "/k8s","mountPath": "/mnt"},{"nasPath": "/files","mountPath": "/app/files"}]`. Among them, nasPath refers to the file storage path; mountPath refers to the path mounted in the container.
* `local_volume` - (Optional) The configuration of the host file mounted to the container. For example: `[{"type":"","nodePath":"/localfiles","mountPath":"/app/files"},{"type":"Directory","nodePath":"/mnt", "mountPath":"/app/storage"}]`. Among them, nodePath is the host path; mountPath is the path in the container; type is the mount type.
* `namespace` - (Optional) The namespace of the K8s cluster, it will determine which K8s namespace your application is deployed in. The default is 'default'.
* `logical_region_id` - (Optional) The ID corresponding to the EDAS namespace, the non-default namespace must be filled in.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the edas k8s application.

## Import

EDAS k8s application can be imported as below, e.g.

```shell
$ terraform import alicloud_edas_k8s_application.new_k8s_application application_id
```
