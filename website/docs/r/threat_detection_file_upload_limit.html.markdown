---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_file_upload_limit"
description: |-
  Provides a Alicloud Threat Detection File Upload Limit resource.
---

# alicloud_threat_detection_file_upload_limit

Provides a Threat Detection File Upload Limit resource. User-defined file upload limit.

For information about Threat Detection File Upload Limit and how to use it, see [What is File Upload Limit](https://next.api.alibabacloud.com/document/Sas/2018-12-03/GetFileUploadLimit).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_file_upload_limit&exampleId=28e17fc4-c3be-dc39-4564-8972b39eb12b1438f8d3&activeTab=example&spm=docs.r.threat_detection_file_upload_limit.0.28e17fc4c3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_file_upload_limit" "default" {
  limit = "100"
}
```

### Deleting `alicloud_threat_detection_file_upload_limit` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_file_upload_limit`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_file_upload_limit&spm=docs.r.threat_detection_file_upload_limit.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `limit` - (Required) File Upload Threshold.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the File Upload Limit.
* `update` - (Defaults to 5 mins) Used when update the File Upload Limit.

## Import

Threat Detection File Upload Limit can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_file_upload_limit.example 
```