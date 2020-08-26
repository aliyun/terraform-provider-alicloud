---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_k8s_application_package_attachment"
sidebar_current: "docs-alicloud-resource-edas-k8s-application-package-attachment"
description: |-
  Deploy an k8s application 
---

# alicloud\_edas\_k8s\_application\_package\_attachment

Deploy an EDAS k8s application.For information about EDAS K8s Application Attachment and how to use it, see [What is EDAS K8s Application Application](https://www.alibabacloud.com/help/doc-detail/149420.htm).

-> **NOTE:** Available in 1.94.0+

## Example Usage

Basic Usage

```
resource "alicloud_edas_k8s_application_package_attachment" "default" {
  app_id                = var.app_id

  // set 'image_url' and 'repo_id' when package_type is 'image'
  image_url             = var.image_url
  repo_id               = var.repo_id

  image                 = var.image

  // set 'package_url','package_version' and 'jdk' when package_type is not 'image'
  package_url           = var.package_url
  package_version       = var.package_version
  jdk = var.jdk

  // set 'web_container' and 'edas_container' when package_type is 'war'
  web_container         = var.web_container
  edas_container_version= var.edas_container_version

  intranet_target_port  = var.intranet_target_port
  intranet_slb_port     = var.intranet_slb_port
  intranet_slb_protocol = var.intranet_slb_protocol
  intranet_slb_id       = var.intranet_slb_id
  
  pre_stop              = var.pre_stop
  envs                  = var.envs
  image_tag             = var.image_tag
  batch_wait_time       = var.batch_wait_time
  command               = var.command
  post_start            = var.post_start
  liveness              = var.liveness
  readiness             = var.readiness
  args                  = var.args
  replicas              = var.replicas
  
  internet_target_port  = var.internet_target_port
  internet_slb_port     = var.internet_slb_port
  internet_slb_protocol = var.internet_slb_protocol
  internet_slb_id       = var.internet_slb_id
  limit_cpu             = var.limit_cpu
  limit_mem             = var.limit_mem
  requests_cpu          = var.requests_cpu
  requests_mem          = var.requests_mem
  requests_m_cpu        = var.requests_m_cpu
  limit_m_cpu           = var.limit_m_cpu
  command               = var.command
  command_args          = var.command_args
  envs                  = var.envs
  pre_stop              = var.pre_stop
  post_start            = var.post_start
  liveness              = var.liveness
  readiness             = var.readiness
  nas_id                = var.nas_id
  mount_descs           = var.mount_descs
  local_volume          = var.local_volume
  namespace             = var.namespace
  logical_region_id     = var.logical_region_id
  uri_encoding          = var.uri_encoding
  use_body_encoding     = var.use_body_encoding
  volumes_str           = var.volumes_str
  package_version_id    = var.package_version_id
  change_order_desc     = var.change_order_desc
  update_strategy       = var.update_strategy
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, ForceNew) The name of the application you want to create. Must start with character,supports numbers, letters and dashes (-), supports up to 36 characters
* `cluster_id` - (Required, ForceNew) The ID of the alicloud container service kubernetes cluster that you want to import to. You can call the ListCluster operation to query.
* `app_id` - (Required, ForceNew) Application ID
* `pre_stop` - (Optional, ForceNew) Execute script before stopping
* `envs` - (Optional, ForceNew)  Deployment environment variables, the format must conform to the JSON object array, such as: `{"name":"x","value":"y"},{"name":"x2","value":"y2"}`, If you want to cancel the configuration, you need to set an empty JSON array "" to indicate no configuration.
* `image_tag` - (Optional, ForceNew) Tag of image.
* `batch_wait_time` - (Optional, ForceNew) The minimum interval between Pod releases in batches, [Click to view more information](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/?spm=a2c4g.11186623.2.20.5663573dc7A1Y9#min-ready-seconds).
* `command` - (Optional, ForceNew) The startup command for container.
* `post_start` - (Optional, ForceNew) Execute script after startup
* `liveness` - (Optional, ForceNew) Container survival status monitoring, format such as: {"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1,"tcpSocket":{"host":"", "port":8080} }.
* `readiness` - (Optional, ForceNew) Container service status check. If the check fails, the traffic passing through K8s Service will not be transferred to the container. The format is: {"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1, "httpGet": {"path": "/consumer","port": 8080,"scheme": "HTTP","httpHeaders": [{"name": "test","value": "testvalue"} ]}}.
* `args` - (Optional, ForceNew) The container start Args command parameter, the format must conform to the JSON array ["parameter 1", "parameter 2"], where the key is a string, and the empty JSON array must be set to an empty JSON array "[]".
* `replicas` - (Optional, ForceNew) Number of application instances.
* `package_type` - (Required, ForceNew) Application package type. Optional parameter values include: FatJar, WAR and Image.
* `image_url` - (Optional, ForceNew) Mirror address. When the package_type is set to 'Image', this parameter item is required.
* `repo_id` - (Optional, ForceNew) The repository id of Image
* `application_descriotion` - (Optional, ForceNew) The description of the application
* `package_url` - (Optional, ForceNew) The url of the package to deploy.Applications deployed through FatJar or WAR packages need to configure it.
* `package_version` - (Optional, ForceNew) The version number of the deployment package. WAR and FatJar types are required. Please customize its meaning.
* `jdk` - (Optional, ForceNew) The JDK version that the deployed package depends on. The optional parameter values are Open JDK 7 and Open JDK 8. Image does not support this parameter.
* `web_container` - (Optional, ForceNew) The Tomcat version that the deployment package depends on. Applicable to Spring Cloud and Dubbo applications deployed through WAR packages. Image does not support this parameter.
* `edas_container_version` - (Optional, ForceNew) EDAS-Container version that the deployed package depends on. Image does not support this parameter.

* `intranet_target_port` - (Optional, ForceNew) The internal SLB back-end port, is also the service port of the application, ranging from 1 to 65535.
* `intranet_slb_port` - (Optional, ForceNew) Intranet SLB front-end port, range 1~65535.
* `intranet_slb_protocol` - (Optional, ForceNew) The private network SLB protocol, supports TCP, HTTP and HTTPS protocols.
* `intranet_slb_id` - (Optional, ForceNew) Private network SLB ID. If not configured, EDAS will automatically purchase a new SLB for user.

* `internet_target_port` - (Optional, ForceNew) The private SLB back-end port, is also the service port of the application, ranging from 1 to 65535.
* `internet_slb_port` - (Optional, ForceNew) The public network SLB front-end port, range 1~65535.
* `internet_slb_protocol` - (Optional, ForceNew) The public network SLB protocol supports TCP, HTTP and HTTPS protocols.
* `internet_slb_id` - (Optional, ForceNew) Public network SLB ID. If not configured, EDAS will automatically purchase a new SLB for the user.

* `image` - (Optional, ForceNew)  The complete URL of mirror, overwrite the mirror Tag parameter.
* `limit_cpu` - (Optional, ForceNew) During application operation, the CPU quota of the application instance, unit: number of cores.
* `limit_mem` - (Optional, ForceNew) The memory limit of the application instance during application operation, unit: M.
* `requests_cpu` - (Optional, ForceNew) When the application is created, the CPU quota of the application instance, unit: number of cores. When set to 0, it means unlimited.
* `requests_mem` - (Optional, ForceNew) When the application is created, the memory limit of the application instance, unit: M. When set to 0, it means unlimited. 
* `requests_mem` - (Optional, ForceNew) When the application is created, the memory limit of the application instance, unit: M. When set to 0, it means unlimited.
* `requests_m_cpu` - (Optional, ForceNew) When the application is created, the CPU quota of the application instance, unit: number of millcores, similar to request_cpu
* `limit_m_cpu` - (Optional, ForceNew)  The CPU quota of the application instance during application operation. Unit: Number of millcores, set to 0 means unlimited, similar to request_cpu.
* `command` - (Optional, ForceNew) The set command, if set, will replace the startup command in the mirror when the mirror is started.
* `command_args` - (Optional, ForceNew) Used in combination with the command, the parameter of the command is a JsonArray string in the format: `[{"argument":"-c"},{"argument":"test"}]`. Among them, -c and test are two parameters that need to be set. 
* `envs` - (Optional, ForceNew)  Deployment environment variables, the format must conform to the JSON object array, such as: `{"name":"x","value":"y"},{"name":"x2","value":"y2"}`, If you want to cancel the configuration, you need to set an empty JSON array "" to indicate no configuration.
* `pre_stop` - (Optional, ForceNew) Execute script before stopping
* `post_start` - (Optional, ForceNew) Execute script after startup
* `liveness` - (Optional, ForceNew) Container survival status monitoring, format such as: {"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1,"tcpSocket":{"host":"", "port":8080} }.
* `readiness` - (Optional, ForceNew) Container service status check. If the check fails, the traffic passing through K8s Service will not be transferred to the container. The format is: {"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1, "httpGet": {"path": "/consumer","port": 8080,"scheme": "HTTP","httpHeaders": [{"name": "test","value": "testvalue"} ]}}.
* `limit_m_cpu` - (Optional, ForceNew)  The CPU quota of the application instance during application operation. Unit: Number of millcores, set to 0 means unlimited, similar to request_cpu. 
* `nas_id` - (Optional, ForceNew) The ID of the mounted NAS must be in the same region as the cluster. It must have an available mount point creation quota, or its mount point must be on a switch in the VPC. If it is not filled in and the mountDescs field exists, a NAS will be automatically purchased and mounted on the switch in the VPC by default.
* `mount_descs` - (Optional, ForceNew) Mount configuration description, as a serialized JSON. For example: `[{"nasPath": "/k8s","mountPath": "/mnt"},{"nasPath": "/files","mountPath": "/app/files"}]`. Among them, nasPath refers to the file storage path; mountPath refers to the path mounted in the container.
* `local_volume` - (Optional, ForceNew) The configuration of the host file mounted to the container. For example: `[{"type":"","nodePath":"/localfiles","mountPath":"/app/files"},{"type":"Directory","nodePath":"/mnt", "mountPath":"/app/storage"}]`. Among them, nodePath is the host path; mountPath is the path in the container; type is the mount type.
* `namespace` - (Optional, ForceNew) The namespace of the K8s cluster, it will determine which K8s namespace your application is deployed in. The default is 'default'.
* `logical_region_id` - (Optional, ForceNew) The ID corresponding to the EDAS namespace, the non-default namespace must be filled in.
* `package_url` - (Optional, ForceNew) The url of the package to deploy.Applications deployed through FatJar or WAR packages need to configure it.
* `package_version` - (Optional, ForceNew) The version number of the deployment package. WAR and FatJar types are required. Please customize its meaning.
* `jdk` - (Optional, ForceNew) The JDK version that the deployed package depends on. The optional parameter values are Open JDK 7 and Open JDK 8. Image does not support this parameter.
* `web_container` - (Optional, ForceNew) The Tomcat version that the deployment package depends on. Applicable to Spring Cloud and Dubbo applications deployed through WAR packages. Image does not support this parameter.
* `edas_container_version` - (Optional, ForceNew) EDAS-Container version that the deployed package depends on. Image does not support this parameter.
* `uri_encoding` - (Optional, ForceNew) URI encoding method supports ISO-8859-1, GBK, GB2312 and UTF-8. The application configuration does not set this parameter and uses the Tomcat default value.
* `use_body_encoding` - (Optional, ForceNew) Whether useBodyEncodingForURI is enabled 
* `volumes_str` - (Optional, ForceNew) Data volume. 
* `package_version_id` - (Optional, ForceNew) Deployment package version ID. 
* `change_order_desc` - (Optional, ForceNew) Change record description. 
* `update_strategy` - (Optional, ForceNew) Batch release strategy. Example 1: 1 set of gray scale + 2 subsequent batches + automatic batching + batching interval 1 minute. {"type":"GrayBatchUpdate","batchUpdate":{"batch":2,"releaseType":"auto","batchWaitTime":1},"grayUpdate":{"gray":1}}; Example 2: 1 set of gray scale + 2 subsequent batches + manual batching. {"type":"GrayBatchUpdate","batchUpdate":{"batch":2,"releaseType":"manual"},"grayUpdate":{"gray":1}}; Example 3: Divide into 2 batches + automatic batch + batch interval 0 minutes. {"type":"BatchUpdate","batchUpdate":{"batch":2,"releaseType":"auto","batchWaitTime":0}} 

## Attributes Reference

The following attributes are exported:

* `application_name` - The name of the application you want to create. Must start with character,supports numbers, letters and dashes (-), supports up to 36 characters
* `cluster_id` - The ID of the alicloud container service kubernetes cluster that you want to import to. You can call the ListCluster operation to query.
* `replicas` - Number of application instances.
* `package_type` -  Application package type. Optional parameter values include: FatJar, WAR and Image.
* `image_url` - Mirror address. When the package_type is set to 'Image', this parameter item is available.
* `id` - The latest deployed package id

## Import

EDAS k8s application deployment can be imported as below, e.g.

```
$ terraform import alicloud_edas_k8s_application_package_attachment.deploy application_id
```