---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_honey_pots"
sidebar_current: "docs-alicloud-datasource-threat-detection-honey-pots"
description: |-
  Provides a list of Threat Detection Honey Pot owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_honey_pots

This data source provides Threat Detection Honey Pot available to the user.[What is Honey Pot](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-createhoneypot)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_threat_detection_honey_pots" "default" {
  ids           = ["xxxx"]
  honeypot_name = "tf-test"
  node_id       = "a44e1ab3-6945-444c-889d-5bacee7056e8"
}

output "alicloud_threat_detection_honey_pot_example_id" {
  value = data.alicloud_threat_detection_honey_pots.default.pots.0.id
}
```

## Argument Reference

The following arguments are supported:
* `honeypot_id` - (ForceNew,Optional) Honeypot ID.
* `honeypot_name` - (ForceNew,Optional) Honeypot custom name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by honey pot name.
* `node_id` - (ForceNew,Optional) The ID of the honeypot management node.
* `ids` - (Optional, ForceNew, Computed) A list of Honey Pot IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Honey Pot IDs.
* `pots` - A list of Honey Pot Entries. Each element contains the following attributes:
    * `id` - Honeypot ID. The value is the same as `honeypot_id`.
    * `honeypot_id` - Honeypot ID.
    * `honeypot_image_id` - The image ID of the honeypot.
    * `honeypot_image_name` - Honeypot mirror name.
    * `honeypot_name` - Honeypot custom name.
    * `node_id` - The ID of the honeypot management node.
    * `page_total` - Total pages.
    * `preset_id` - The custom parameter ID of honeypot.
    * `state` - Honeypot status.
    * `status` - The status of the resource
