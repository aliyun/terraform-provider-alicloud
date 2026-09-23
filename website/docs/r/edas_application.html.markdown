---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_application"
sidebar_current: "docs-alicloud-resource-edas-application"
description: |-
  Creates an EDAS ecs application on EDAS.
---

# alicloud_edas_application

Creates an EDAS ecs application on EDAS, see [What is EDAS Application](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-insertapplication). The application will be deployed when `group_id` and `war_url` are given.

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_application&exampleId=0d33502b-a28a-930b-1056-9f7ba7a8bbcbe2c86382&activeTab=example&spm=docs.r.edas_application.0.0d33502ba2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_edas_cluster" "default" {
  cluster_name      = "${var.name}-${random_integer.default.result}"
  cluster_type      = "2"
  network_mode      = "2"
  logical_region_id = data.alicloud_regions.default.regions.0.id
  vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_edas_application" "default" {
  application_name = "${var.name}-${random_integer.default.result}"
  cluster_id       = alicloud_edas_cluster.default.id
  package_type     = "JAR"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_application&spm=docs.r.edas_application.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) Name of your EDAS application. Only letters '-' '_' and numbers are allowed. The length cannot exceed 36 characters.
* `package_type` - (Required, ForceNew) The type of the package for the deployment of the application that you want to create. The valid values are: WAR and JAR. We strongly recommend you to set this parameter when creating the application.
* `cluster_id` - (Required, ForceNew) The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter. 
* `build_pack_id` - (Optional) The package ID of Enterprise Distributed Application Service (EDAS) Container, which can be retrieved by calling container version list interface ListBuildPack or the "Pack ID" column in container version list. When creating High-speed Service Framework (HSF) application, this parameter is required.
* `descriotion` - (Optional) The description of the application that you want to create.
* `health_check_url` - (Optional) The URL for health checking of the application.
* `logical_region_id` - (Optional) The ID of the namespace where you want to create the application. You can call the ListUserDefineRegion operation to query the namespace ID.
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
