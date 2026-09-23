---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_bindings"
sidebar_current: "docs-alicloud-datasource-amqp-bindings"
description: |-
  Provides a list of Amqp Bindings to the user.
---

# alicloud\_amqp\_bindings

This data source provides the Amqp Bindings of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_amqp_bindings" "examples" {
  instance_id       = "amqp-cn-xxxxx"
  virtual_host_name = "my-vh"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Instance Id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `virtual_host_name` - (Required, ForceNew) Virtualhost Name.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Binding IDs. Each item value formats as `<instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>`.
* `bindings` - A list of Amqp Bindings. Each element contains the following attributes:
	* `argument` - X-match Attributes. Valid Values: All: Default Value, All the Message Header of Key-Value Pairs Stored in the Must Match. Any: at Least One Pair of the Message Header of Key-Value Pairs Stored in the Must Match. This Parameter Applies Only to Headers Exchange Other Types of Exchange Is Invalid. Other Types of Exchange Here Can Either Be an Arbitrary Value.
	* `binding_key` - The Binding Key. The Source of the Binding Exchange Non-Topic Type: Can Only Contain Letters, Lowercase Letters, Numbers, and the Dash (-), the Underscore Character (_), English Periods (.) and the at Sign (@). Length from 1 to 255 Characters. The Source of the Binding Exchange Topic Type: Can Contain Letters, Lowercase Letters, Numbers, and the Dash (-), the Underscore Character (_), English Periods (.) and the at Sign (@). If You Include the Hash (.
	* `binding_type` - The Target Binding Types.
	* `destination_name` - The Target Queue Or Exchange of the Name.
	* `id` - The ID of the Binding. The value formats as `<instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>`.
	* `instance_id` - Instance Id.
	* `source_exchange` - The Source Exchange Name.
	* `virtual_host_name` - Virtualhost Name.
