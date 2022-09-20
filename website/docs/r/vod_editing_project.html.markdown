---
subcategory: "ApsaraVideo VoD"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_editing_project"
sidebar_current: "docs-alicloud-resource-vod-editing-project"
description: |-
  Provides a Alicloud VOD Editing Project resource.
---

# alicloud\_vod\_editing\_project

Provides a VOD Editing Project resource.

For information about VOD Editing Project and how to use it, see [What is Editing Project](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/addeditingproject#doc-api-vod-AddEditingProject).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vod_editing_project" "example" {
  editing_project_name = "example_value"
  title                = "example_value"
  timeline             = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `cover_url` - (Optional) The thumbnail URL of the online editing project. If you do not specify this parameter and the video track in the timeline has mezzanine files, the thumbnail of the first mezzanine file in the timeline is used.
* `division` - (Optional) The region where you want to create the online editing project.
* `editing_project_name` - (Optional, Computed) The description of the online editing project.
* `timeline` - (Optional, Computed) The timeline of the online editing project, in JSON format. For more information about the structure, see [Timeline](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/basic-structures). If you do not specify this parameter, an empty timeline is created and the duration of the online editing project is zero.
* `title` - (Required) The title of the online editing project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Editing Project.
* `status` - The Status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Editing Project.
* `update` - (Defaults to 1 mins) Used when update the Editing Project.
* `delete` - (Defaults to 1 mins) Used when delete the Editing Project.


## Import

VOD Editing Project can be imported using the id, e.g.

```
$ terraform import alicloud_vod_editing_project.example <id>
```