---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application"
sidebar_current: "docs-alicloud-resource-sae-application"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Application resource.
---

# alicloud_sae_application

Provides a Serverless App Engine (SAE) Application resource.

For information about Serverless App Engine (SAE) Application and how to use it, see [What is Application](https://www.alibabacloud.com/help/en/sae/latest/createapplication).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sae_application&exampleId=e567998f-7b80-0723-0c08-0c427c7d00e4a3e17339&activeTab=example&spm=docs.r.sae_application.0.e567998f7b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

data "alicloud_regions" "default" {
  current = true
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
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_application" "default" {
  app_description   = var.name
  app_name          = "${var.name}-${random_integer.default.result}"
  namespace_id      = alicloud_sae_namespace.default.id
  image_url         = "registry-vpc.${data.alicloud_regions.default.regions.0.id}.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type      = "Image"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vswitch_id        = alicloud_vswitch.default.id
  timezone          = "Asia/Beijing"
  replicas          = "5"
  cpu               = "500"
  memory            = "2048"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, ForceNew) Application Name. Combinations of numbers, letters, and dashes (-) are allowed. It must start with a letter and the maximum length is 36 characters.
* `package_type` - (Required, ForceNew) Application package type. Valid values: `FatJar`, `War`, `Image`, `PhpZip`, `IMAGE_PHP_5_4`, `IMAGE_PHP_5_4_ALPINE`, `IMAGE_PHP_5_5`, `IMAGE_PHP_5_5_ALPINE`, `IMAGE_PHP_5_6`, `IMAGE_PHP_5_6_ALPINE`, `IMAGE_PHP_7_0`, `IMAGE_PHP_7_0_ALPINE`, `IMAGE_PHP_7_1`, `IMAGE_PHP_7_1_ALPINE`, `IMAGE_PHP_7_2`, `IMAGE_PHP_7_2_ALPINE`, `IMAGE_PHP_7_3`, `IMAGE_PHP_7_3_ALPINE`, `PythonZip`.
* `replicas` - (Required, Int) Initial number of instances.
* `namespace_id` - (Optional, ForceNew) SAE namespace ID. Only namespaces whose names are lowercase letters and dashes (-) are supported, and must start with a letter. The namespace can be obtained by calling the DescribeNamespaceList interface.
* `vpc_id` - (Optional, ForceNew) The vpc id.
* `vswitch_id` - (Optional) The vswitch id. **NOTE:** From version 1.211.0, `vswitch_id` can be modified.
* `package_version` - (Optional) The version number of the deployment package. Required when the Package Type is War and FatJar.
* `package_url` - (Optional) Deployment package address. Only FatJar or War type applications can configure the deployment package address.
* `image_url` - (Optional) Mirror address. Only Image type applications can configure the mirror address.
* `cpu` - (Optional, Int) The CPU required for each instance, in millicores, cannot be 0. Valid values: `500`, `1000`, `2000`, `4000`, `8000`, `16000`, `32000`.
* `memory` - (Optional, Int) The memory required for each instance, in MB, cannot be 0. One-to-one correspondence with CPU. Valid values: `1024`, `2048`, `4096`, `8192`, `12288`, `16384`, `24576`, `32768`, `65536`, `131072`.
* `command` - (Optional) Mirror start command. The command must be an executable object in the container. For example: sleep. Setting this command will cause the original startup command of the mirror to become invalid.
* `web_container` - (Optional) The version of tomcat that the deployment package depends on. Image type applications are not supported.
* `jdk` - (Optional) The JDK version that the deployment package depends on. Image type applications are not supported.
* `jar_start_options` - (Optional) The JAR package starts the application option. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
* `jar_start_args` - (Optional) The JAR package starts application parameters. Application default startup command: $JAVA_HOME/bin/java $JarStartOptions -jar $CATALINA_OPTS "$package_path" $JarStartArgs.
* `app_description` - (Optional) Application description information. No more than 1024 characters. **NOTE:** From version 1.211.0, `app_description` can be modified.
* `auto_config` - (Optional, Bool) The auto config. Valid values: `true`, `false`.
* `auto_enable_application_scaling_rule` - (Optional, Bool) The auto enable application scaling rule. Valid values: `true`, `false`.
* `batch_wait_time` - (Optional, Int) The batch wait time.
* `change_order_desc` - (Optional) The change order desc.
* `deploy` - (Optional, Bool) The deploy. Valid values: `true`, `false`.
* `edas_container_version` - (Optional) The operating environment used by the Pandora application.
* `enable_ahas` - (Optional) The enable ahas. Valid values: `true`, `false`.
* `enable_grey_tag_route` - (Optional, Bool) The enable grey tag route. Default value: `false`. Valid values:
  - `true`: The canary release rule is enabled.
  - `false`: The canary release rule is disabled.
  **NOTE:** Currently, `enable_grey_tag_route` can only be set to `false`, and if you want to set it to `true`, you must operate on the web console.
* `min_ready_instances` - (Optional, Int) The Minimum Available Instance. On the Change Had Promised during the Available Number of Instances to Be.
* `min_ready_instance_ratio` - (Optional, Int) Minimum Survival Instance Percentage. **NOTE:** When `min_ready_instances` and `min_ready_instance_ratio` are passed at the same time, and the value of `min_ready_instance_ratio` is not -1, the `min_ready_instance_ratio` parameter shall prevail. Assuming that `min_ready_instances` is 5 and `min_ready_instance_ratio` is 50, 50 is used to calculate the minimum number of surviving instances.The value description is as follows:
  * `-1`: Initialization value, indicating that percentages are not used.
  * `0~100`: The unit is percentage, rounded up. For example, if it is set to 50%, if there are currently 5 instances, the minimum number of surviving instances is 3.
* `oss_ak_id` - (Optional, Sensitive) OSS AccessKey ID.
* `oss_ak_secret` - (Optional, Sensitive) OSS  AccessKey Secret.
* `php_arms_config_location` - (Optional) The PHP application monitors the mount path, and you need to ensure that the PHP server will load the configuration file of this path. You don't need to pay attention to the configuration content, SAE will automatically render the correct configuration file.
* `php_config` - (Optional) PHP configuration file content.
* `php_config_location` - (Optional) PHP application startup configuration mount path, you need to ensure that the PHP server will start using this configuration file.
* `security_group_id` - (Optional) Security group ID.
* `termination_grace_period_seconds` - (Optional, Int) Graceful offline timeout, the default is 30, the unit is seconds. The value range is 1~60. Valid values: [1,60].
* `timezone` - (Optional) Time zone. Default value: `Asia/Shanghai`.
* `war_start_options` - (Optional) WAR package launch application option. Application default startup command: java $JAVA_OPTS $CATALINA_OPTS [-Options] org.apache.catalina.startup.Bootstrap "$@" start.
* `acr_instance_id` - (Optional, Available since v1.189.0) The ID of the ACR EE instance. Only necessary if the image_url is pointing to an ACR EE instance.
* `acr_assume_role_arn` - (Optional, Available since v1.189.0) The ARN of the RAM role required when pulling images across accounts. Only necessary if the image_url is pointing to an ACR EE instance.
* `micro_registration` - (Optional, Available since v1.198.0) Select the Nacos registry. Valid values: `0`, `1`, `2`.
* `envs` - (Optional) Container environment variable parameters. For example,`	[{"name":"envtmp","value":"0"}]`. The value description is as follows:
  * `name` - environment variable name.
  * `value` - Environment variable value or environment variable reference.
* `sls_configs` - (Optional) SLS  configuration.
* `php` - (Optional, Available since v1.211.0) The Php environment.
* `image_pull_secrets` - (Optional, Available since v1.211.0) The ID of the corresponding Secret.
* `programming_language` - (Optional, ForceNew, Available since v1.211.0) The programming language that is used to create the application. Valid values: `java`, `php`, `other`.
* `command_args_v2` - (Optional, List, Available since v1.211.0) The parameters of the image startup command.
* `custom_host_alias_v2` - (Optional, Set, Available since v1.211.0) The custom mapping between the hostname and IP address in the container. See [`custom_host_alias_v2`](#custom_host_alias_v2) below.
* `oss_mount_descs_v2` - (Optional, Set, Available since v1.211.0) The description of the mounted Object Storage Service (OSS) bucket. See [`oss_mount_descs_v2`](#oss_mount_descs_v2) below.
* `config_map_mount_desc_v2` - (Optional, Set, Available since v1.211.0) The description of the ConfigMap that is mounted to the application. A ConfigMap that is created on the ConfigMaps page of a namespace is used to inject configurations into containers. See [`config_map_mount_desc_v2`](#config_map_mount_desc_v2) below.
* `liveness_v2` - (Optional, Set, Available since v1.211.0) The liveness check settings of the container. See [`liveness_v2`](#liveness_v2) below.
* `readiness_v2` - (Optional, Set, Available since v1.211.0) The readiness check settings of the container. If a container fails this health check multiple times, the container is stopped and then restarted. See [`readiness_v2`](#readiness_v2) below.
* `post_start_v2` - (Optional, Set, Available since v1.211.0) The script that is run immediately after the container is started. See [`post_start_v2`](#post_start_v2) below.
* `pre_stop_v2` - (Optional, Set, Available since v1.211.0) The script that is run before the container is stopped. See [`pre_stop_v2`](#pre_stop_v2) below.
* `tomcat_config_v2` - (Optional, Set, Available since v1.211.0) The Tomcat configuration. See [`tomcat_config_v2`](#tomcat_config_v2) below.
* `update_strategy_v2` - (Optional, Set, Available since v1.211.0) The release policy. See [`update_strategy_v2`](#update_strategy_v2) below.
* `nas_configs` - (Optional, Set, Available since v1.211.0) The configurations for mounting the NAS file system. See [`nas_configs`](#nas_configs) below.
* `kafka_configs` - (Optional, Set, Available since v1.211.0) The logging configurations of ApsaraMQ for Kafka. See [`kafka_configs`](#kafka_configs) below.
* `pvtz_discovery_svc` - (Optional, Set, Available since v1.211.0) The configurations of Kubernetes Service-based service registration and discovery. See [`pvtz_discovery_svc`](#pvtz_discovery_svc) below.
* `tags` - (Optional, Available since v1.167.0) A mapping of tags to assign to the resource.
* `status` - (Optional) The status of the resource. Valid values: `RUNNING`, `STOPPED`, `UNKNOWN`.
* `command_args` - (Deprecated since v1.211.0) Mirror startup command parameters. The parameters required for the above start command. For example: 1d. **NOTE:** Field `command_args` has been deprecated from provider version 1.211.0. New field `command_args_v2` instead.
* `custom_host_alias` - (Deprecated since v1.211.0) Custom host mapping in the container. For example: [{`hostName`:`samplehost`,`ip`:`127.0.0.1`}]. **NOTE:** Field `custom_host_alias` has been deprecated from provider version 1.211.0. New field `custom_host_alias_v2` instead.
* `oss_mount_descs` - (Deprecated since v1.211.0) OSS mount description information. **NOTE:** Field `oss_mount_descs` has been deprecated from provider version 1.211.0. New field `oss_mount_descs_v2` instead.
* `config_map_mount_desc` - (Deprecated since v1.211.0) ConfigMap mount description. **NOTE:** Field `config_map_mount_desc` has been deprecated from provider version 1.211.0. New field `config_map_mount_desc_v2` instead.
* `liveness` - (Deprecated since v1.211.0) Container health check. Containers that fail the health check will be shut down and restored. Currently, only the method of issuing commands in the container is supported.
  **NOTE:** Field `liveness` has been deprecated from provider version 1.211.0. New field `liveness_v2` instead.
* `readiness` - (Deprecated since v1.211.0) Application startup status checks, containers that fail multiple health checks will be shut down and restarted. Containers that do not pass the health check will not receive SLB traffic. For example: {`exec`:{`command`:[`sh`,"-c","cat /home/admin/start.sh"]},`initialDelaySeconds`:30,`periodSeconds`:30,"timeoutSeconds ":2}. Valid values: `command`, `initialDelaySeconds`, `periodSeconds`, `timeoutSeconds`.
  **NOTE:** Field `readiness` has been deprecated from provider version 1.211.0. New field `readiness_v2` instead.
* `post_start` - (Deprecated since v1.211.0) Execute the script after startup, the format is like: {`exec`:{`command`:[`cat`,"/etc/group"]}}. **NOTE:** Field `post_start` has been deprecated from provider version 1.211.0. New field `post_start_v2` instead.
* `pre_stop` - (Deprecated since v1.211.0) Execute the script before stopping, the format is like: {`exec`:{`command`:[`cat`,"/etc/group"]}}. **NOTE:** Field `pre_stop` has been deprecated from provider version 1.211.0. New field `pre_stop_v2` instead.
* `tomcat_config` - (Deprecated since v1.211.0) Tomcat file configuration, set to "{}" means to delete the configuration:  useDefaultConfig: Whether to use a custom configuration, if it is true, it means that the custom configuration is not used; if it is false, it means that the custom configuration is used. If you do not use custom configuration, the following parameter configuration will not take effect.  contextInputType: Select the access path of the application.  war: No need to fill in the custom path, the access path of the application is the WAR package name. root: No need to fill in the custom path, the access path of the application is /. custom: You need to fill in the custom path in the custom path below. contextPath: custom path, this parameter only needs to be configured when the contextInputType type is custom.  httpPort: The port range is 1024~65535. Ports less than 1024 need Root permission to operate. Because the container is configured with Admin permissions, please fill in a port greater than 1024. If not configured, the default is 8080. maxThreads: Configure the number of connections in the connection pool, the default size is 400. uriEncoding: Tomcat encoding format, including UTF-8, ISO-8859-1, GBK and GB2312. If not set, the default is ISO-8859-1. useBodyEncoding: Whether to use BodyEncoding for URL. Valid values: `contextInputType`, `contextPath`, `httpPort`, `maxThreads`, `uriEncoding`, `useBodyEncoding`, `useDefaultConfig`.
  **NOTE:** Field `tomcat_config` has been deprecated from provider version 1.211.0. New field `tomcat_config_v2` instead.
* `update_strategy` - (Deprecated since v1.211.0) The update strategy. **NOTE:** Field `update_strategy` has been deprecated from provider version 1.211.0. New field `update_strategy_v2` instead.
* `nas_id` - (Removed since v1.211.0) ID of the mounted NAS, Must be in the same region as the cluster. It must have an available mount point creation quota, or its mount point must be on a switch in the VPC. If it is not filled in and the mountDescs field is present, a NAS will be automatically purchased and mounted on the switch in the VPC by default.
  **NOTE:** Field `nas_id` has been removed from provider version 1.211.0.
* `mount_host` - (Removed since v1.211.0) Mount point of NAS in application VPC. **NOTE:** Field `mount_host` has been removed from provider version 1.211.0.
* `mount_desc` - (Removed since v1.211.0) Mount description. **NOTE:** Field `mount_desc` has been removed from provider version 1.211.0.
* `version_id` - (Removed since v1.211.0) Application version id. **NOTE:** Field `version_id` has been removed from provider version 1.211.0.

### `custom_host_alias_v2`

The custom_host_alias_v2 supports the following:

* `host_name` (Optional) The domain name or hostname.
* `ip` (Optional) The IP address.

### `oss_mount_descs_v2`

The oss_mount_descs_v2 supports the following:

* `bucket_name` (Optional) The name of the OSS bucket.
* `bucket_path` (Optional) The directory or object in OSS.
* `mount_path` (Optional) The path of the container in SAE.
* `read_only` (Optional, Bool) Specifies whether the application can use the container path to read data from or write data to resources in the directory of the OSS bucket. Valid values:
  - `true`: The application has the read-only permissions.
  - `false`: The application has the read and write permissions.

### `config_map_mount_desc_v2`

The config_map_mount_desc_v2 supports the following:

* `config_map_id` (Optional) The ID of the ConfigMap.
* `mount_path` (Optional) The mount path.
* `key` (Optional) The key.

### `liveness_v2`

The liveness_v2 supports the following:

* `initial_delay_seconds` (Optional, Int) The delay of the health check.
* `period_seconds` (Optional, Int) The interval at which the health check is performed.
* `timeout_seconds` (Optional, Int) The timeout period of the health check.
* `exec` - (Optional, Set) Execute. See [`exec`](#liveness_v2-exec) below.
* `tcp_socket` - (Optional, Set) The liveness check settings of the container. See [`tcp_socket`](#liveness_v2-tcp_socket) below.
* `http_get` - (Optional, Set) The liveness check settings of the container. See [`http_get`](#liveness_v2-http_get) below.

### `liveness_v2-exec`

The exec supports the following:

* `command` - (Optional, List) The health check command.

### `liveness_v2-tcp_socket`

The tcp_socket supports the following:

* `port` (Optional, Int) The port that is used to check the status of TCP connections.

### `liveness_v2-http_get`

The http_get supports the following:

* `path` (Optional) The request path.
* `port` (Optional, Int) The port.
* `scheme` (Optional) The protocol that is used to perform the health check. Valid values: `HTTP` and `HTTPS`.
* `key_word` (Optional) The custom keywords.
* `is_contain_key_word` (Optional, Bool) Specifies whether the response contains keywords. Valid values: `true` and `false`. If you do not set it, the advanced settings are not used.

### `readiness_v2`

The readiness_v2 supports the following:

* `initial_delay_seconds` (Optional, Int) The delay of the health check.
* `period_seconds` (Optional, Int) The interval at which the health check is performed.
* `timeout_seconds` (Optional, Int) The timeout period of the health check.
* `exec` - (Optional, Set) Execute. See [`exec`](#readiness_v2-exec) below.
* `tcp_socket` - (Optional, Set) The liveness check settings of the container. See [`tcp_socket`](#readiness_v2-tcp_socket) below.
* `http_get` - (Optional, Set) The liveness check settings of the container. See [`http_get`](#readiness_v2-http_get) below.

### `readiness_v2-exec`

The exec supports the following:

* `command` - (Optional, List) The health check command.

### `readiness_v2-tcp_socket`

The tcp_socket supports the following:

* `port` (Optional, Int) The port that is used to check the status of TCP connections.

### `readiness_v2-http_get`

The http_get supports the following:

* `path` (Optional) The request path.
* `port` (Optional, Int) The port.
* `scheme` (Optional) The protocol that is used to perform the health check. Valid values: `HTTP` and `HTTPS`.
* `key_word` (Optional) The custom keywords.
* `is_contain_key_word` (Optional, Bool) Specifies whether the response contains keywords. Valid values: `true` and `false`. If you do not set it, the advanced settings are not used.

### `post_start_v2`

The post_start_v2 supports the following:

* `exec` - (Optional, Set) Execute. See [`exec`](#post_start_v2-exec) below.

### `post_start_v2-exec`

The exec supports the following:

* `command` - (Optional, List) The command.

### `pre_stop_v2`

The pre_stop_v2 supports the following:

* `exec` - (Optional, Set) Execute. See [`exec`](#pre_stop_v2-exec) below.

### `pre_stop_v2-exec`

The exec supports the following:

* `command` - (Optional, List) The command.

### `tomcat_config_v2`

The tomcat_config_v2 supports the following:

* `port` (Optional, Int) The port.
* `max_threads` (Optional, Int) The maximum number of connections in the connection pool.
* `context_path` (Optional) The path.
* `uri_encoding` (Optional) The URI encoding scheme in the Tomcat container.
* `use_body_encoding_for_uri` (Optional) Specifies whether to use the encoding scheme that is specified by BodyEncoding for URL.

### `update_strategy_v2`

The update_strategy_v2 supports the following:

* `type` (Optional) The type of the release policy. Valid values: `GrayBatchUpdate` and `BatchUpdate`.
* `batch_update` - (Optional, Set) The phased release policy. See [`batch_update`](#update_strategy_v2-batch_update) below.

### `update_strategy_v2-batch_update`

The batch_update supports the following:

* `release_type` - (Optional) The processing method for the batches. Valid values: `auto` and `manual`.
* `batch` (Optional, Int) The number of batches in which you want to release the instances.
* `batch_wait_time` (Optional, Int) The time interval at which the instances in a batch are deployed. Unit: seconds.

### `nas_configs`

The nas_configs supports the following:

* `nas_id` (Optional) The ID of the NAS file system.
* `nas_path` (Optional, Int) The directory in the NAS file system.
* `mount_path` (Optional) The mount path of the container.
* `mount_domain` (Optional) The domain name of the mount target.
* `read_only` (Optional, Bool) Specifies whether the application can read data from or write data to resources in the directory of the NAS. Valid values: `true` and `false`. If you set `read_only` to `false`, the application has the read and write permissions.

### `kafka_configs`

The kafka_configs supports the following:

* `kafka_instance_id` (Optional) The  ID of the ApsaraMQ for Kafka instance.
* `kafka_endpoint` (Optional) The endpoint of the ApsaraMQ for Kafka API.
* `kafka_configs` - (Optional, Set) One or more logging configurations of ApsaraMQ for Kafka. See [`kafka_configs`](#kafka_configs-kafka_configs) below.

### `kafka_configs-kafka_configs`

The kafka_configs supports the following:

* `log_type` - (Optional) The type of the log.
* `log_dir` (Optional) The path in which logs are stored.
* `kafka_topic` (Optional) The topic of the Kafka.

### `pvtz_discovery_svc`

The pvtz_discovery_svc supports the following:

* `service_name` - (Optional, ForceNew) The name of the Service.
* `namespace_id` (Optional, ForceNew) The ID of the namespace.
* `enable` (Optional, Bool) Enables the Kubernetes Service-based registration and discovery feature.
* `port_protocols` - (Optional, Set) The port number and protocol. See [`port_protocols`](#pvtz_discovery_svc-port_protocols) below.

### `pvtz_discovery_svc-port_protocols`

The port_protocols supports the following:

* `port` - (Optional, Int) The port.
* `protocol` (Optional) The protocol. Valid values: `TCP` and `UDP`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application.

## Import

Serverless App Engine (SAE) Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_application.example <id>
```
