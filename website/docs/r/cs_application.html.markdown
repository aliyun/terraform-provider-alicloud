---
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_application"
sidebar_current: "docs-alicloud-resource-cs-application"
description: |-
  Provides a resource to deploy application in one container cluster.
---

# alicloud\_cs\_application

This resource use an orchestration template to define and deploy a multi-container application. An application is created by using an orchestration template.
Each application can contain one or more services.

-> **NOTE:** Application orchestration template must be a valid Docker Compose YAML template.

-> **NOTE:** At present, this resource only support swarm cluster.

## Example Usage

Basic Usage

```
resource "alicloud_cs_application" "app" {
  cluster_name = "my-first-swarm"
  name = "wordpress"
  version = "1.2"
  template = "${file("wordpress.yml")}"
  latest_image = true
  environment = {
    EXTERNAL_URL = "123.123.123.123:8080"
  }
}
```
## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) The swarm cluster's name.
* `name` - (Required, ForceNew) The application name. It should be 1-64 characters long, and can contain numbers, English letters and hyphens, but cannot start with hyphens.
* `description` - The description of application.
* `version` - The application deploying version. Each updating, it must be different with current. Default to "1.0"
* `template` - (Required) The application deployment template and it must be [Docker Compose format](https://docs.docker.com/compose/).
* `environment` - A key/value map used to replace the variable parameter in the Compose template.
* `latest_image` - Whether to use latest docker image while each updating application. Default to false.
* `blue_green` - Wherther to use "Blue Green" method when release a new version. Default to false.
* `blue_green_confirm` - Whether to confirm a "Blue Green" application. Default to false. It will be ignored when `blue_green` is false.

-> **NOTE:** Each update of `template`, `environment`, `latest_image` and `blue_green`, it requires a new `version`. Otherwise, the update will be ignored.

-> **NOTE:** If you want to rollback a "Blue Green" application, just set `blue_green` as false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container application. It's formate is `<cluster_name>:<name>`.
* `cluster_name` - The name of the container cluster.
* `name` - The application name.
* `description` - The application description.
* `template` - The application deploying template.
* `environment` - The application environment variables.
* `services` - List of services in the application. It contains several attributes to `Block Nodes`.
* `default_domain` - The application default domain and it can be used to configure routing service.


### Block Nodes

* `id` - ID of the service.
* `name` - Service name.
* `status` - The current status of service.
* `version` - The current version of service.


## Import

Swarm application can be imported using the id, e.g.

```
$ terraform import alicloud_cs_application.app my-first-swarm:wordpress
```