---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_repos"
sidebar_current: "docs-alicloud-datasource-cr-ee-repos"
description: |-
  Provides a list of Container Registry Enterprise Edition repositories.
---

# alicloud\_cr\_ee\_repos

This data source provides a list Container Registry Enterprise Edition repositories on Alibaba Cloud.

-> **NOTE:** Available in v1.87.0+

## Example Usage

```
# Declare the data source
data "alicloud_cr_ee_repos" "my_repos" {
  instance_id = "cri-xx"
  name_regex  = "my-repos"
  output_file = "my-repo-json"
}

output "output" {
  value = "${data.alicloud_cr_ee_repos.my_repos.repos}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of Container Registry Enterprise Edition instance.
* `namespace` - (Optional) Name of Container Registry Enterprise Edition namespace where the repositories are located in.
* `ids` - (Optional) A list of ids to filter results by repository id.
* `name_regex` - (Optional) A regex string to filter results by repository name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Boolean, false by default, only repository attributes are exported. Set to true if tags belong to this repository are needed. See `tags` in attributes.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Enterprise Edition repositories. Its element is a repository id.
* `names` - A list of repository names.
* `repos` - A list of matched Container Registry Enterprise Edition namespaces. Each element contains the following attributes:
  * `instance_id` - ID of Container Registry Enterprise Edition instance.
  * `namespace` - Name of Container Registry Enterprise Edition namespace where repo is located.
  * `id` - ID of Container Registry Enterprise Edition repository.
  * `name` - Name of Container Registry Enterprise Edition repository.
  * `summary` - The repository general information.
  * `repo_type` - `PUBLIC` or `PRIVATE`, repository's visibility.
  * `tags` - A list of image tags belong to this repository. Each contains several attributes, see `Block Tag`.

### Block Tag

* `tag` - Tag of this image.
* `image_id` - Id of this image.
* `digest` - Digest of this image.
* `status` - Status of this image.
* `image_size` - Status of this image, in bytes.
* `image_update` - Last update time of this image, unix time in nanoseconds.
* `image_create` - Create time of this image, unix time in nanoseconds.

