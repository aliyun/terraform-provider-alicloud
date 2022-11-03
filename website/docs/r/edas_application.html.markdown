---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_application"
sidebar_current: "docs-alicloud-resource-edas-application"
description: |-
  Creates an EDAS ecs application on EDAS.
---

# alicloud\_edas\_application

Creates an EDAS ecs application on EDAS. The application will be deployed when `group_id` and `war_url` are given.

-> **NOTE:** Available in 1.82.0+

## Example Usage

Basic Usage

```terraform
resource "alicloud_edas_application" "default" {
  application_name  = "xxx"
  cluster_id        = "xxx"
  package_type      = "JAR"
  build_pack_id     = xxx
  descriotion       = "xxx"
  health_check_url  = "xxx"
  logical_region_id = "cn-xxxx:xxx"
  component_ids     = xxx
  ecu_info          = ["xxx"]
  group_id          = "xxx"
  package_version   = "xxx"
  war_url           = "http://xxx"
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) Name of your EDAS application. Only letters '-' '_' and numbers are allowed. The length cannot exceed 36 characters.
* `package_type` - (Required) The type of the package for the deployment of the application that you want to create. The valid values are: WAR and JAR. We strongly recommend you to set this parameter when creating the application.
* `cluster_id` - (Required) The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter. 
* `build_pack_id` - (Optional) The package ID of Enterprise Distributed Application Service (EDAS) Container, which can be retrieved by calling container version list interface ListBuildPack or the "Pack ID" column in container version list. When creating High-speed Service Framework (HSF) application, this parameter is required.
* `descriotion` - (Optional) The description of the application that you want to create.
* `health_check_url` - (Optional) The URL for health checking of the application.
* `logical_region_id` - (Optional) The ID of the namespace where you want to create the application. You can call the ListUserDefineRegion operation to query the namespace ID.
* `component_ids` - (Optional) The ID of the component in the container where the application is going to be deployed. If the runtime environment is not specified when the application is created and the application is not deployed, you can set the parameter as fellow: when deploying a native Dubbo or Spring Cloud application using a WAR package for the first time, you must specify the version of the Apache Tomcat component based on the deployed application. You can call the ListClusterOperation interface to query the components. When deploying a non-native Dubbo or Spring Cloud application using a WAR package for the first time, you can leave this parameter empty. 
* `ecu_info` - (Optional) The ID of the Elastic Compute Unit (ECU) where you want to deploy the application. Type: List.
* `group_id` - (Optional) The ID of the instance group where the application is going to be deployed. Set this parameter to all if you want to deploy the application to all groups.
* `package_version` - (Optional) The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommended you to use a timestamp.
* `war_url` - (Optional) The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the deployType parameter is set as url.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `app_Id`.

## Import

EDAS application can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_application.app app_Id
```
