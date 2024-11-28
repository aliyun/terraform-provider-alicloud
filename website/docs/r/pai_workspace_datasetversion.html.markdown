---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_datasetversion"
description: |-
  Provides a Alicloud PAI Workspace Datasetversion resource.
---

# alicloud_pai_workspace_datasetversion

Provides a PAI Workspace Datasetversion resource.



For information about PAI Workspace Datasetversion and how to use it, see [What is Datasetversion](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_workspace" "defaultAiWorkspace" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "defaultDataset" {
  accessibility    = "PRIVATE"
  source_type      = "USER"
  data_type        = "PIC"
  workspace_id     = alicloud_pai_workspace_workspace.defaultAiWorkspace.id
  options          = jsonencode({ "mountPath" : "/mnt/data/" })
  description      = var.name
  source_id        = "d-xxxxx_v1"
  uri              = "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/"
  dataset_name     = format("%s1", var.name)
  user_id          = "1511928242963727"
  data_source_type = "OSS"
  property         = "DIRECTORY"
}


resource "alicloud_pai_workspace_datasetversion" "default" {
  options          = jsonencode({ "mountPath" : "/mnt/data/verion/" })
  description      = var.name
  data_source_type = "OSS"
  source_type      = "USER"
  source_id        = "d-xxxxx_v1"
  data_size        = "2068"
  data_count       = "1000"
  labels {
    key   = "key1"
    value = "example1"
  }
  uri        = "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/"
  property   = "DIRECTORY"
  dataset_id = alicloud_pai_workspace_dataset.defaultDataset.id
}
```

## Argument Reference

The following arguments are supported:
* `data_count` - (Optional, Int) Data count.
* `data_size` - (Optional, Int) Data size.
* `data_source_type` - (Required, ForceNew) The data source type. The following values are supported:
  - OSS: Alibaba Cloud Object Storage (OSS).
  - NAS: Alibaba cloud file storage (NAS).
* `dataset_id` - (Required, ForceNew) The first ID of the resource
* `description` - (Optional) Description of dataset version.
* `labels` - (Optional, ForceNew, List) The tag of the resource See [`labels`](#labels) below.
* `options` - (Optional) The extended field, which is of the JsonString type.

  When DLC uses a dataset, you can specify the default Mount path for the dataset by configuring the mountPath field.
* `property` - (Required, ForceNew) The properties of the dataset. The following values are supported:
  - FILE: FILE.
  - DIRECTORY: folder.
* `source_id` - (Optional, ForceNew) The data source ID.
* `source_type` - (Optional, ForceNew) The data source type. The default value is USER. 
* `uri` - (Required, ForceNew) The Uri configuration sample is as follows:
  - The data source type is OSS:'oss:// bucket.endpoint/object'
  - The data source type is NAS:

  The general NAS format is: 'nas://.region/subpath/to/dir/';

  CPFS1.0:'nas://.region/subpath/to/dir /';

  CPFS2.0:'nas://.region//'.

  CPFS1.0 and CPFS2.0 are distinguished by the format of fsid: CPFS1.0 is cpfs-;CPFS2.0 is cpfs-.

### `labels`

The labels supports the following:
* `key` - (Optional, ForceNew) The key of the tags
* `value` - (Optional, ForceNew) The value of the tags

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<dataset_id>:<version_name>`.
* `create_time` - Update time.
* `version_name` - The name of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Datasetversion.
* `delete` - (Defaults to 5 mins) Used when delete the Datasetversion.
* `update` - (Defaults to 5 mins) Used when update the Datasetversion.

## Import

PAI Workspace Datasetversion can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_datasetversion.example <dataset_id>:<version_name>
```