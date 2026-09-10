---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_log_analysis"
description: |-
  Provides an Alicloud File Storage (NAS) Log Analysis resource.
---

# alicloud_nas_log_analysis

Provides a File Storage (NAS) Log Analysis resource.

The log delivery configuration of a NAS file system.

For information about File Storage (NAS) Log Analysis and how to use it, see [What is Log Analysis](https://next.api.alibabacloud.com/document/NAS/2017-06-26/CreateLogAnalysis).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_log_analysis&exampleId=acb22f59-bcb9-4c32-1926-5108f65761cd895e686d&activeTab=example&spm=docs.r.nas_log_analysis.0.acb22f59bc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_log_analysis" "default" {
  file_system_id = alicloud_nas_file_system.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_log_analysis&spm=docs.r.nas_log_analysis.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `file_system_id` - (Required, ForceNew) The ID of the file system for which log delivery is enabled.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `logstore` - The name of the Logstore that receives NAS logs.
* `project` - The name of the project that receives NAS logs.
* `region` - The Simple Log Service region of the log project.
* `role_arn` - The ARN of the service role used by NAS to deliver logs to Simple Log Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when creating the Log Analysis.
* `delete` - (Defaults to 5 mins) Used when deleting the Log Analysis.

## Import

File Storage (NAS) Log Analysis can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_log_analysis.example <file_system_id>
```
