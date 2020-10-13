---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_application_deployment"
sidebar_current: "docs-alicloud-resource-edas-application-deployment"
description: |-
  Deploys applications on EDAS.
---

# alicloud\_edas\_application\_deployment

Deploys applications on EDAS.

-> **NOTE:** Available in 1.82.0+

## Example Usage

Basic Usage

```
resource "alicloud_edas_application_deployment" "default" {
  app_id          = var.app_id
  group_id        = var.group_id
  package_version = var.package_version
  war_url         = var.war_url
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy.
* `group_id` - (Required, ForceNew) The ID of the instance group where the application is going to be deployed. Set this parameter to all if you want to deploy the application to all groups.
* `package_version` - (Optional, ForceNew) The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommended you to use a timestamp.
* `war_url` - (Required, ForceNew) The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the deployType parameter is set as url.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<package_version>`.
* `last_package_version` - Last package version deployed.

