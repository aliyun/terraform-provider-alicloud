---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_deploy_group"
sidebar_current: "docs-alicloud-resource-edas-deploy-group"
description: |-
  Provides an EDAS deploy group resource.
---

# alicloud_edas_deploy_group

Provides an EDAS deploy group resource, see [What is EDAS Deploy Group](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-insertdeploygroup).

-> **NOTE:** Available since v1.82.0.


## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_edas_cluster" "default" {
  cluster_name      = var.name
  cluster_type      = "2"
  network_mode      = "2"
  logical_region_id = data.alicloud_regions.default.regions.0.id
  vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_edas_application" "default" {
  application_name = var.name
  cluster_id       = alicloud_edas_cluster.default.id
  package_type     = "JAR"
}

resource "alicloud_edas_deploy_group" "default" {
  app_id     = alicloud_edas_application.default.id
  group_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy.
* `group_name` - (Required, ForceNew) The name of the instance group that you want to create. 
* `group_type` - (ForceNew) The type of the instance group that you want to create. Valid values: 0: Default group. 1: Phased release is disabled for traffic management. 2: Phased release is enabled for traffic management.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<group_name>:<group_id>`.

## Import

EDAS deploy group can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_deploy_group.group app_id:group_name:group_id
```
