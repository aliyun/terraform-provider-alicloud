---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_groups"
description: |-
  Provides a list of Sae Application Groups to the user.
---

# alicloud_sae_application_groups

This data source provides the Sae Application Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_sae_applications" "default" {
  app_name = "default-NODELETING-image"
}

data "alicloud_sae_application_groups" "default" {
  app_id = data.alicloud_sae_applications.default.applications.0.id
}

output "sae_application_groups_id_0" {
  value = data.alicloud_sae_application_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Group IDs.
* `name_regex` - (Optional) A regex string to filter results by Group name.
* `app_id` - (Required) The application ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Group names.
* `groups` - A list of Groups. Each element contains the following attributes:
  * `id` - The ID of the Group.
  * `group_id` - The ID of the application group.
  * `group_name` - The name of the application group.
  * `group_type` - The type of the application group.
  * `replicas` - The total number of instances.
  * `package_type` - The type of the deployment package.
  * `package_version` - The deployment package version.
  * `package_version_id` - The ID of the deployment package version.
  * `package_url` - The URL of the deployment package.
  * `image_url` - The image URL.
  * `web_container` - The Tomcat container version required by the deployment package.
  * `jdk` - The JDK version required by the deployment package.
  * `running_instances` - The number of running instances.
  * `edas_container_version` - The version of the application runtime environment in the High-Speed Service Framework (HSF), such as an Ali-Tomcat container.
