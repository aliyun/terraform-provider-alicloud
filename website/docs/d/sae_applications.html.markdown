---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_applications"
sidebar_current: "docs-alicloud-datasource-sae-applications"
description: |-
  Provides a list of Sae Applications to the user.
---

# alicloud\_sae\_applications

This data source provides the Sae Applications of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.161.0+.

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
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = "cn-hangzhou:tfacctest"
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = "tf-testaccDescription"
  app_name        = "tf-testaccAppName131"
  namespace_id    = alicloud_sae_namespace.default.id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  vswitch_id      = alicloud_vswitch.vsw.id
  timezone        = "Asia/Beijing"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}
data "alicloud_sae_applications" "default" {
  ids = [alicloud_sae_application.default.id]
}
output "sae_application_id" {
  value = data.alicloud_sae_applications.default.applications.0.id
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Optional, ForceNew) Application Name. Combinations of numbers, letters, and dashes (-) are allowed. It must start with a letter and the maximum length is 36 characters.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `field_type` - (Optional, ForceNew) The field type. Valid values:`appName`, `appIds`, `slbIps`, `instanceIps`
* `field_value` - (Optional, ForceNew) The field value.
* `ids` - (Optional, ForceNew, Computed)  A list of Application IDs.
* `namespace_id` - (Optional, ForceNew) SAE namespace ID. Only namespaces whose names are lowercase letters and dashes (-) are supported, and must start with a letter. The namespace can be obtained by calling the DescribeNamespaceList interface.
* `order_by` - (Optional, ForceNew) The order by.Valid values:`running`,`instances`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `reverse` - (Optional, ForceNew) The reverse.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `RUNNING`, `STOPPED`,`UNKNOWN`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `applications` - A list of Sae Applications. Each element contains the following attributes:
    * `acr_assume_role_arn` - The ARN of the RAM role required when pulling images across accounts.
    * `app_description` - Application description information. No more than 1024 characters.
    * `app_name` - Application Name. Combinations of numbers, letters, and dashes (-) are allowed. It must start with a letter and the maximum length is 36 characters.
    * `application_id` - The first ID of the resource.
    * `command` - Mirror start command. The command must be an executable object in the container. For example: sleep. Setting this command will cause the original startup command of the mirror to become invalid.
    * `command_args` - Mirror startup command parameters. The parameters required for the above start command. For example: 1d.
    * `config_map_mount_desc` - ConfigMap mount description.
    * `cpu` - The CPU required for each instance, in millicores, cannot be 0.
    * `create_time` - Indicates That the Application of the Creation Time.
    * `custom_host_alias` - Custom host mapping in the container. For example: [{"hostName":"samplehost","ip":"127.0.0.1"}].
    * `edas_container_version` - The operating environment used by the Pandora application.
    * `envs` - The virtual switch where the elastic network card of the application instance is located. The switch must be located in the aforementioned VPC. The switch also has a binding relationship with the SAE namespace. If it is left blank, the default is the vSwitch ID bound to the namespace.
    * `id` - The ID of the Application.
    * `image_url` - Mirror address. Only Image type applications can configure the mirror address.
    * `jar_start_args` - The JAR package starts application parameters. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
    * `jar_start_options` - The JAR package starts the application option. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
    * `jdk` - The JDK version that the deployment package depends on. Image type applications are not supported.
    * `liveness` - Container health check. Containers that fail the health check will be shut down and restored. Currently, only the method of issuing commands in the container is supported.
    * `memory` - The memory required for each instance, in MB, cannot be 0. One-to-one correspondence with CPU.
    * `min_ready_instances` - The Minimum Available Instance. On the Change Had Promised during the Available Number of Instances to Be.
    * `mount_desc` - Mount description information.
      * `mount_path` - Container mount path.
      * `nas_path` - NAS relative file directory.
    * `mount_host` - Mount point of NAS in application VPC.
    * `namespace_id` - SAE namespace ID. Only namespaces whose names are lowercase letters and dashes (-) are supported, and must start with a letter. The namespace can be obtained by calling the DescribeNamespaceList interface.
    * `nas_id` - ID of the mounted NAS, Must be in the same region as the cluster. It must have an available mount point creation quota, or its mount point must be on a switch in the VPC. If it is not filled in and the mountDescs field is present, a NAS will be automatically purchased and mounted on the switch in the VPC by default.
    * `oss_ak_id` - OSS AccessKey ID.
    * `oss_ak_secret` - OSS  AccessKey Secret.
    * `oss_mount_descs` - OSS mount description information.
    * `oss_mount_details` - The OSS mount detail. 
      * `bucket_name` - The name of the bucket.
      * `bucket_path` - The path of the bucket.
      * `mount_path` - The Container mount path.
      * `read_only` - Whether the container path has readable permission to mount directory resources.
    * `package_type` - Application package type. Support FatJar, War and Image.
    * `package_url` - Deployment package address. Only FatJar or War type applications can configure the deployment package address.
    * `package_version` - The version number of the deployment package. Required when the Package Type is War and FatJar.
    * `php_arms_config_location` - The PHP application monitors the mount path, and you need to ensure that the PHP server will load the configuration file of this path. You don't need to pay attention to the configuration content, SAE will automatically render the correct configuration file.
    * `php_config` - PHP configuration file content.
    * `php_config_location` - PHP application startup configuration mount path, you need to ensure that the PHP server will start using this configuration file.
    * `post_start` - Execute the script after startup, the format is like: {"exec":{"command":["cat","/etc/group"]}}.
    * `pre_stop` - Execute the script before stopping, the format is like: {"exec":{"command":["cat","/etc/group"]}}.
    * `readiness` - Application startup status checks, containers that fail multiple health checks will be shut down and restarted. Containers that do not pass the health check will not receive SLB traffic. For example: {"exec":{"command":["sh","-c","cat /home/admin/start.sh"]},"initialDelaySeconds":30,"periodSeconds":30,"timeoutSeconds ":2}.
    * `replicas` - Initial number of instances.
    * `security_group_id` - Security group ID.
    * `sls_configs` - SLS  configuration.
    * `status` - The status of the resource.
    * `termination_grace_period_seconds` - Graceful offline timeout, the default is 30, the unit is seconds. The value range is 1~60.
    * `timezone` - Time zone, the default value is Asia/Shanghai.
    * `tomcat_config` - Tomcat file configuration, set to "" or "{}" means to delete the configuration:  useDefaultConfig: Whether to use a custom configuration, if it is true, it means that the custom configuration is not used; if it is false, it means that the custom configuration is used. If you do not use custom configuration, the following parameter configuration will not take effect.  contextInputType: Select the access path of the application.  war: No need to fill in the custom path, the access path of the application is the WAR package name. root: No need to fill in the custom path, the access path of the application is /. custom: You need to fill in the custom path in the custom path below. contextPath: custom path, this parameter only needs to be configured when the contextInputType type is custom.  httpPort: The port range is 1024~65535. Ports less than 1024 need Root permission to operate. Because the container is configured with Admin permissions, please fill in a port greater than 1024. If not configured, the default is 8080. maxThreads: Configure the number of connections in the connection pool, the default size is 400. uriEncoding: Tomcat encoding format, including UTF-8, ISO-8859-1, GBK and GB2312. If not set, the default is ISO-8859-1. useBodyEncoding: Whether to use BodyEncoding for URL.
    * `vpc_id` - The VPC corresponding to the SAE namespace. In SAE, a namespace can only correspond to one VPC and cannot be modified. Creating a SAE application in the namespace for the first time will form a binding relationship. Multiple namespaces can correspond to a VPC. If you leave it blank, it will default to the VPC ID bound to the namespace.
    * `vswitch_id` - The vswitch id.
    * `war_start_options` - WAR package launch application option. Application default startup command: java $JAVA_OPTS $CATALINA_OPTS [-Options] org.apache.catalina.startup.Bootstrap "$@" start.
    * `web_container` - The version of tomcat that the deployment package depends on. Image type applications are not supported.
    * `tags` - A mapping of tags to assign to the resource.
    