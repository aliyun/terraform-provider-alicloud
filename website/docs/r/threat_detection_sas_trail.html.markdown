---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_sas_trail"
description: |-
  Provides a Alicloud Threat Detection Sas Trail resource.
---

# alicloud_threat_detection_sas_trail

Provides a Threat Detection Sas Trail resource. 

For information about Threat Detection Sas Trail and how to use it, see [What is Sas Trail](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-createservicetrail).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_sas_trail&exampleId=6c289b18-fd60-a14a-cc0a-e7348dd3f12c3382041c&activeTab=example&spm=docs.r.threat_detection_sas_trail.0.6c289b18fd&intl_lang=EN_US" target="_blank">
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


resource "alicloud_threat_detection_sas_trail" "default" {
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_sas_trail&spm=docs.r.threat_detection_sas_trail.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.
* `create_time` - The service trace creation timestamp, in milliseconds.
* `service_trail` - Service trace configuration information.
  * `config` - Service tracking on status. The value is:
  - **on:** Open
  - **off:** off.
  * `update_time` - The timestamp of the last service update. Unit: milliseconds.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Sas Trail.
* `delete` - (Defaults to 5 mins) Used when delete the Sas Trail.

## Import

Threat Detection Sas Trail can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_sas_trail.example 
```