---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_applications"
sidebar_current: "docs-alicloud-datasource-edas-applications"
description: |-
    Provides a list of EDAS applications available to the user.
---

# alicloud\_edas\_applications

This data source provides a list of EDAS application in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.82.0+

## Example Usage

```
data "alicloud_edas_applications" "applications" {
  ids = ["xxx"]
  output_file = "application.txt"
}

output "first_application_name" {
  value = data.alicloud_edas_applications.applications.applications.0.app_name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) An ids string to filter results by the application id. 
* `name_regex` - (Optional) A regex string to filter results by the application name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of application IDs.
* `names` - A list of applications names.
* `applications` - A list of applications.
  * `app_name` - The name of your EDAS application. Only letters '-' '_' and numbers are allowed. The length cannot exceed 36 characters.
  * `app_id` - The ID of the application that you want to deploy.
  * `application_type` - The type of the package for the deployment of the application that you want to create. The valid values are: WAR and JAR. We strongly recommend you to set this parameter when creating the application.
  * `build_package_id` - The package ID of Enterprise Distributed Application Service (EDAS) Container.
  * `cluster_id` - The ID of the cluster that you want to create the application.
  * `cluster_type` -  The type of the cluster that you want to create. Valid values: 1: Swarm cluster. 2: ECS cluster. 3: Kubernates cluster.
  * `region_id` - The ID of the namespace the application belongs to.

