---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_folder"
sidebar_current: "docs-alicloud-resource-data-works-folder"
description: |-
  Provides a Alicloud Data Works Folder resource.
---

# alicloud_data_works_folder

Provides a Data Works Folder resource.

For information about Data Works Folder and how to use it, see [What is Folder](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-createfolder).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_folder&exampleId=85c5de6b-3c82-dea2-eec1-bb74fea767d42934f637&activeTab=example&spm=docs.r.data_works_folder.0.85c5de6b3c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_data_works_folder" "example" {
  project_id  = "320687"
  folder_path = "Business Flow/tfTestAcc/folderDi/tftest1"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_data_works_folder&spm=docs.r.data_works_folder.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `folder_path` - (Required) Folder Path. The folder path composed with for part: `Business Flow/{Business Flow Name}/[folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined]/{Directory Name}`. The first segment of path must be `Business Flow`, and sencond segment of path must be a Business Flow Name within the project. The third part of path must be one of those keywords:`folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined`. Then the finial part of folder path can be specified in yourself.
* `project_id` - (Optional, ForceNew, Available in v1.131.0+) The ID of the DataWorks workspace. You can obtain the workspace ID from the workspace configuration page in the DataWorks console. Either this parameter or `project_identifier` must be specified to determine which DataWorks workspace this folder belongs to.
* `project_identifier` - (Optional) The name of the DataWorks workspace. You can obtain the workspace name from the workspace configuration page in the DataWorks console. Either this parameter or `project_id` must be specified to determine which DataWorks workspace this folder belongs to.

## Attributes Reference

The following attributes are exported:

* `folder_id` - The ID of the DataWorks folder.
* `id` - The resource ID of Folder. The value formats as `<folder_id>:<$.ProjectId>`.

## Import

Data Works Folder can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_folder.example <folder_id>:<$.ProjectId>
```
