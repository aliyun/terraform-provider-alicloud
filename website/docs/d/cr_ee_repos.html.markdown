---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_repos"
description: |-
  Provides a list of Container Registry Enterprise Edition Repositories to the user.
---

# alicloud_cr_ee_repos

This data source provides the Container Registry Enterprise Edition Repositories of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.87.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_cr_ee_instances" "default" {
  name_regex = "default-nodeleting"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  name               = var.name
  auto_create        = true
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_namespace.default.instance_id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = var.name
}

data "alicloud_cr_ee_repos" "ids" {
  ids         = [alicloud_cr_ee_repo.default.repo_id]
  instance_id = alicloud_cr_ee_repo.default.instance_id
}

output "cr_ee_repos_id_0" {
  value = data.alicloud_cr_ee_repos.ids.repos.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Repository IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Repository name.
* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `namespace` - (Optional, ForceNew) The name of the namespace to which the Repository belongs.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Repository names.
* `repos` -  A list of Repositories. Each element contains the following attributes:
  * `id` - The ID of the Repository.
  * `instance_id` - The ID of the Container Registry instance to which the Repository belongs.
  * `namespace` - The name of the namespace to which the Repository belongs.
  * `name` - The name of the Repository.
  * `summary` - The summary of the Repository.
  * `repo_type` - The type of the Repository.
  * `tags` - A list of image tags belong to this Repository. **Note:** `tags` takes effect only if `enable_details` is set to `true`.
    * `tag` - The tag of the image.
    * `image_id` - The ID of the image.
    * `image_size` - The size of the image.
    * `digest` - The digest of the image.
    * `status` - The status of the image.
    * `image_create` - The time when the image was created.  
    * `image_update` - The time when the image was last updated.
