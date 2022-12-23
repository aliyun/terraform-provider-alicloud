---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honeypot_images"
sidebar_current: "docs-alicloud-datasource-threat-detection-honeypot-images"
description: |-
  Provides a list of Threat Detection Honeypot Image owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_honeypot_images

This data source provides Threat Detection Honeypot Image available to the user.[What is Honeypot Image](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-listavailablehoneypot)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_threat_detection_honeypot_images" "default" {
  ids        = ["sha256:02882320c9a55303410127c5dc4ae2dc470150f9d7f2483102d994f5e5f4d9df"]
  name_regex = "^meta"
}

output "alicloud_threat_detection_honeypot_image_example_id" {
  value = data.alicloud_threat_detection_honeypot_images.default.images.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Honeypot Image IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Honeypot mirror nam.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Honeypot Image IDs.
* `names` - A list of name of Honeypot Images.
* `images` - A list of Honeypot Image Entries. Each element contains the following attributes:
    * `id` - The image ID of the honeypot.The value is the same as `honeypot_image_id`.
    * `honeypot_image_display_name` - The name of the honeypot image display.
    * `honeypot_image_id` - The image ID of the honeypot.
    * `honeypot_image_name` - Honeypot mirror name.
    * `honeypot_image_type` - Honeypot mirror type.
    * `honeypot_image_version` - Honeypot Mirror version.
    * `multiports` - Ports supported by honeypots. In JSON format. Contains the following fields:-**log_type**: log type-**proto**: Support Protocol-**description**: description-**ports**: supports Port collection-**port_str**: supports port strings-**type**: type
    * `proto` - Honeypot-supported protocols.
    * `service_port` - Honeypot service port.
    * `template` - Honeypot configuration parameter template.
