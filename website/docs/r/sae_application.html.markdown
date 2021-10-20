---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application"
sidebar_current: "docs-alicloud-resource-sae-application"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Application resource.
---

# alicloud\_sae\_application

Provides a Serverless App Engine (SAE) Application resource.

For information about Serverless App Engine (SAE) Application and how to use it, see [What is Application](https://help.aliyun.com/document_detail/97792.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_testacc"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id = "cn-hangzhou:tfacctest"
  namespace_name = var.name
}

resource "alicloud_sae_application" "default" {
  app_description= "tf-testaccDescription"
  app_name=        "tf-testaccAppName"
  namespace_id=    alicloud_sae_namespace.default.id
  image_url=     "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type=    "Image"
  vswitch_id=      alicloud_vswitch.vsw.id
  timezone =     "Asia/Beijing"
  replicas=        "5"
  cpu=             "500"
  memory =          "2048"
}
```

## Argument Reference

The following arguments are supported:

* `app_description` - (Optional, ForceNew) Application description information. No more than 1024 characters.
* `app_name` - (Required, ForceNew) Application Name. Combinations of numbers, letters, and dashes (-) are allowed. It must start with a letter and the maximum length is 36 characters.
* `auto_config` - (Optional) The auto config. Valid values: `false`, `true`.
* `auto_enable_application_scaling_rule` - (Optional) The auto enable application scaling rule. Valid values: `false`, `true`.
* `batch_wait_time` - (Optional) The batch wait time.
* `change_order_desc` - (Optional) The change order desc.
* `command` - (Optional) Mirror start command. The command must be an executable object in the container. For example: sleep. Setting this command will cause the original startup command of the mirror to become invalid.
* `command_args` - (Optional) Mirror startup command parameters. The parameters required for the above start command. For example: 1d.
* `config_map_mount_desc` - (Optional) ConfigMap mount description.
* `cpu` - (Optional) The CPU required for each instance, in millicores, cannot be 0. Valid values: `1000`, `16000`, `2000`, `32000`, `4000`, `500`, `8000`.
* `custom_host_alias` - (Optional) Custom host mapping in the container. For example: [{`hostName`:`samplehost`,`ip`:`127.0.0.1`}]. 
* `deploy` - (Optional) The deploy. Valid values: `false`, `true`.
* `edas_container_version` - (Optional) The operating environment used by the Pandora application.
* `enable_ahas` - (Optional) The enable ahas.
* `enable_grey_tag_route` - (Optional) The enable grey tag route.
* `envs` - (Optional) The virtual switch where the elastic network card of the application instance is located. The switch must be located in the aforementioned VPC. The switch also has a binding relationship with the SAE namespace. If it is left blank, the default is the vSwitch ID bound to the namespace.
* `image_url` - (Optional) Mirror address. Only Image type applications can configure the mirror address.
* `jar_start_args` - (Optional) The JAR package starts application parameters. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
* `jar_start_options` - (Optional) The JAR package starts the application option. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
* `jdk` - (Optional) The JDK version that the deployment package depends on. Image type applications are not supported.
* `liveness` - (Optional) Container health check. Containers that fail the health check will be shut down and restored. Currently, only the method of issuing commands in the container is supported.
* `memory` - (Optional) The memory required for each instance, in MB, cannot be 0. One-to-one correspondence with CPU. Valid values: `1024`, `131072`, `16384`, `2048`, `32768`, `4096`, `65536`, `8192`.
* `min_ready_instances` - (Optional) The Minimum Available Instance. On the Change Had Promised during the Available Number of Instances to Be.
* `mount_desc` - (Optional) Mount description.
* `mount_host` - (Optional) Mount point of NAS in application VPC.
* `namespace_id` - (Optional, ForceNew) SAE namespace ID. Only namespaces whose names are lowercase letters and dashes (-) are supported, and must start with a letter. The namespace can be obtained by calling the DescribeNamespaceList interface.
* `nas_id` - (Optional) ID of the mounted NAS, Must be in the same region as the cluster. It must have an available mount point creation quota, or its mount point must be on a switch in the VPC. If it is not filled in and the mountDescs field is present, a NAS will be automatically purchased and mounted on the switch in the VPC by default.
* `oss_ak_id` - (Optional, Sensitive) OSS AccessKey ID.
* `oss_ak_secret` - (Optional, Sensitive) OSS  AccessKey Secret.
* `oss_mount_descs` - (Optional) OSS mount description information.
* `package_type` - (Required, ForceNew) Application package type. Support FatJar, War and Image. Valid values: `FatJar`, `Image`, `War`.
* `package_url` - (Optional) Deployment package address. Only FatJar or War type applications can configure the deployment package address.
* `package_version` - (Optional) The version number of the deployment package. Required when the Package Type is War and FatJar.
* `php_arms_config_location` - (Optional) The PHP application monitors the mount path, and you need to ensure that the PHP server will load the configuration file of this path. You don't need to pay attention to the configuration content, SAE will automatically render the correct configuration file.
* `php_config` - (Optional) PHP configuration file content.
* `php_config_location` - (Optional) PHP application startup configuration mount path, you need to ensure that the PHP server will start using this configuration file.
* `post_start` - (Optional) Execute the script after startup, the format is like: {`exec`:{`command`:[`cat`,"/etc/group"]}}.
* `pre_stop` - (Optional) Execute the script before stopping, the format is like: {`exec`:{`command`:[`cat`,"/etc/group"]}}.
* `readiness` - (Optional) Application startup status checks, containers that fail multiple health checks will be shut down and restarted. Containers that do not pass the health check will not receive SLB traffic. For example: {`exec`:{`command`:[`sh`,"-c","cat /home/admin/start.sh"]},`initialDelaySeconds`:30,`periodSeconds`:30,"timeoutSeconds ":2}. Valid values: `command`, `initialDelaySeconds`, `periodSeconds`, `timeoutSeconds`.
* `replicas` - (Required, ForceNew) Initial number of instances.
* `security_group_id` - (Optional) Security group ID.
* `sls_configs` - (Optional) SLS  configuration.
* `status` - (Optional, Computed) The status of the resource. Valid values: `RUNNING`, `STOPPED`.
* `termination_grace_period_seconds` - (Optional, Computed) Graceful offline timeout, the default is 30, the unit is seconds. The value range is 1~60. Valid values: [1,60].
* `timezone` - (Optional) Time zone, the default value is Asia/Shanghai.
* `tomcat_config` - (Optional) Tomcat file configuration, set to "{}" means to delete the configuration:  useDefaultConfig: Whether to use a custom configuration, if it is true, it means that the custom configuration is not used; if it is false, it means that the custom configuration is used. If you do not use custom configuration, the following parameter configuration will not take effect.  contextInputType: Select the access path of the application.  war: No need to fill in the custom path, the access path of the application is the WAR package name. root: No need to fill in the custom path, the access path of the application is /. custom: You need to fill in the custom path in the custom path below. contextPath: custom path, this parameter only needs to be configured when the contextInputType type is custom.  httpPort: The port range is 1024~65535. Ports less than 1024 need Root permission to operate. Because the container is configured with Admin permissions, please fill in a port greater than 1024. If not configured, the default is 8080. maxThreads: Configure the number of connections in the connection pool, the default size is 400. uriEncoding: Tomcat encoding format, including UTF-8, ISO-8859-1, GBK and GB2312. If not set, the default is ISO-8859-1. useBodyEncoding: Whether to use BodyEncoding for URL. Valid values: `contextInputType`, `contextPath`, `httpPort`, `maxThreads`, `uriEncoding`, `useBodyEncoding`, `useDefaultConfig`.
* `update_strategy` - (Optional) The update strategy.
* `version_id` - (Optional, ForceNew) Application version id.
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `vpc_id` - (Optional, ForceNew, Available in v1.139+) The vpc id.
* `war_start_options` - (Optional) WAR package launch application option. Application default startup command: java $JAVA_OPTS $CATALINA_OPTS [-Options] org.apache.catalina.startup.Bootstrap "$@" start.
* `web_container` - (Optional) The version of tomcat that the deployment package depends on. Image type applications are not supported.
* `intranet` - (Optional,Available in v1.139+) Bound public network SLB. The details see Block intranet.
* `internet` - (Optional,Available in v1.139+) Bound private network SLB. The details see Block internet.
* `internet_slb_id` - (Optional,Available in v1.139+) public network SLB ID.
* `intranet_slb_id` - (Optional,Available in v1.139+) private network SLB ID.

#### intranet

The intranet supports the following:
* `https_cert_id` - (Optional) SSL certificate. `https_cert_id` is required when HTTPS is selected
* `protocol` - (Optional) Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `target_port` - (Optional) Container port.
* `port` - (Optional) SLB Port.

#### internet

The internet supports the following:
* `https_cert_id` - (Optional) SSL certificate. `https_cert_id` is required when HTTPS is selected
* `protocol` - (Optional) Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `target_port` - (Optional) Container port.
* `port` - (Optional) SLB Port.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application.
* `internet_ip` - Use designated public network SLBs that have been purchased to support non-shared instances. **NOTE:** Available in v1.139+.
* `intranet_ip` - Use the designated private network SLB that has been purchased to support non-shared instances. **NOTE:** Available in v1.139+.

## Import

Serverless App Engine (SAE) Application can be imported using the id, e.g.

```
$ terraform import alicloud_sae_application.example <id>
```