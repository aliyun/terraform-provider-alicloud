---
subcategory: "ApsaraVideo VoD (VOD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_editing_project"
sidebar_current: "docs-alicloud-resource-vod-editing-project"
description: |-
  Provides a Alicloud VOD Editing Project resource.
---

# alicloud_vod_editing_project

Provides a VOD Editing Project resource.

For information about VOD Editing Project and how to use it, see [What is Editing Project](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/addeditingproject#doc-api-vod-AddEditingProject).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vod_editing_project&exampleId=f0e8df38-bb3c-d8b4-c822-14eced3b4bfda1bfabf3&activeTab=example&spm=docs.r.vod_editing_project.0.f0e8df38bb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_vod_editing_project" "example" {
  editing_project_name = var.name
  title                = var.name
  timeline             = <<EOF
  {
    "VideoTracks":[
      {
        "VideoTrackClips":[
          {
          "MediaId":"0c60e6f02dae71edbfaa472190a90102",
          "In":2811
          }
        ]
      }
    ]
  }
  EOF
  cover_url            = "https://demo.aliyundoc.com/6AB4D0E1E1C74468883516C2349D1FC2-6-2.png"
  division             = data.alicloud_regions.default.regions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cover_url` - (Optional) The thumbnail URL of the online editing project. If you do not specify this parameter and the video track in the timeline has mezzanine files, the thumbnail of the first mezzanine file in the timeline is used.
* `division` - (Optional) The region where you want to create the online editing project.
* `editing_project_name` - (Optional, Computed) The description of the online editing project.
* `timeline` - (Optional) The timeline of the online editing project, in JSON format. For more information about the structure, see [Timeline](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/basic-structures). If you do not specify this parameter, an empty timeline is created and the duration of the online editing project is zero.
* `title` - (Required) The title of the online editing project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Editing Project.
* `status` - The Status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Editing Project.
* `update` - (Defaults to 1 mins) Used when update the Editing Project.
* `delete` - (Defaults to 1 mins) Used when delete the Editing Project.


## Import

VOD Editing Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_vod_editing_project.example <id>
```