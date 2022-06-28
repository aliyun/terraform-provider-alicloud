---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_k8s_application"
sidebar_current: "docs-alicloud-resource-edas-k8s-application"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alicloud\_edas\_k8s\_application

Create an EDAS k8s application.For information about EDAS K8s Application and how to use it, see [What is EDAS K8s Application](https://www.alibabacloud.com/help/doc-detail/85029.htm). 

-> **NOTE:** Available in 1.105.0+

## Example Usage

Basic Usage

```terraform
resource "alicloud_edas_k8s_application" "default" {
  // package type is Image / FatJar / War
  package_type            = "Image"
  application_name        = "DemoApplication"
  application_descriotion = "This is description of application"
  cluster_id              = var.cluster_id
  replicas                = 2

  // set 'image_url' and 'repo_id' when package_type is 'image'
  image_url = "registry-vpc.cn-beijing.aliyuncs.com/edas-demo-image/consumer:1.0"

  // set 'package_url','package_version' and 'jdk' when package_type is not 'image'
  package_url     = var.package_url
  package_version = var.package_version
  jdk             = var.jdk

  // set 'web_container' and 'edas_container' when package_type is 'war'
  web_container          = var.web_container
  edas_container_version = var.edas_container_version

  internet_target_port  = var.internet_target_port
  internet_slb_port     = var.internet_slb_port
  internet_slb_protocol = var.internet_slb_protocol
  internet_slb_id       = var.internet_slb_id
  limit_cpu             = 4
  limit_mem             = 2048
  requests_cpu          = 0
  requests_mem          = 0
  requests_m_cpu        = 0
  limit_m_cpu           = 4000
  command               = var.command
  command_args          = var.command_args
  envs                  = var.envs
  pre_stop              = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  post_start            = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  liveness              = var.liveness
  readiness             = var.readiness
  nas_id                = var.nas_id
  mount_descs           = var.mount_descs
  local_volume          = var.local_volume
  namespace             = "default"
  logical_region_id     = cn-beijing
}
```

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
* `internet_target_port` - (Optional, ForceNew) The private SLB back-end port, is also the service port of the application, ranging from 1 to 65535.
* `internet_slb_port` - (Optional, ForceNew) The public network SLB front-end port, range 1~65535.
* `internet_slb_protocol` - (Optional, ForceNew) The public network SLB protocol supports TCP, HTTP and HTTPS protocols.
* `internet_slb_id` - (Optional, ForceNew) Public network SLB ID. If not configured, EDAS will automatically purchase a new SLB for the user.
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
* `application_name` - The name of the application you want to create. Must start with character,supports numbers, letters and dashes (-), supports up to 36 characters
* `cluster_id` - The ID of the alicloud container service kubernetes cluster that you want to import to. You can call the ListCluster operation to query.

## Import

EDAS k8s application can be imported as below, e.g.

```
$ terraform import alicloud_edas_k8s_application.new_k8s_application application_id
```
