---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_message"
description: |-
  Provides a Alicloud PAIPlugin Message resource.
---

# alicloud_pai_plugin_message

Provides a PAIPlugin Message resource. SendMessage creates and sends a message;
the message status is later read back via ListMessages using the captured send
date.

For information about PAIPlugin Message and how to use it, see
[What is SendMessage](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-sendmessage).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
resource "alicloud_pai_plugin_message" "default" {
  phone_number  = "13800000000"
  signature     = "tf-sig"
  template_code = "SMS_1"
}
```

## Argument Reference

The following arguments are supported:

* `phone_number` - (Required, ForceNew) The phone number.
* `signature` - (Optional, ForceNew) The signature name. Either `signature` or `signature_id` may be set.
* `signature_id` - (Optional, ForceNew) The signature id. Either `signature` or `signature_id` may be set.
* `template_code` - (Optional, ForceNew) The template code. Either `template_code` or `template_id` may be set.
* `template_id` - (Optional, ForceNew) The template id. Either `template_code` or `template_id` may be set.
* `template_params` - (Optional, ForceNew) The template params, JSON strings keyed by variable name.
* `payment_type` - (Optional, ForceNew) The payment type.
* `schedule_id_create` - (Optional, ForceNew) The schedule id used to associate a reach plan at send time.
* `sms_up_extend_codes` - (Optional, ForceNew) The sms up extend codes.
* `out_ids` - (Optional, ForceNew) The out ids.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `message_id`.
* `message_id` - The message id, returned by SendMessage.
* `send_date` - The send date (yyyymmdd) captured at creation.
* `group_id` - The associated group id, empty when unassociated.
* `out_id` - The out id.
* `schedule_id` - The associated schedule id, empty when unassociated.
* `error_code` - The error code of the message.
* `status` - The send status. 0: sending, 1: success, 2: failure.
* `template_type` - The template type. 0: verification code, 1: notification, 2: marketing.

-> **NOTE:** A sent message cannot be unsent; PAIPlugin exposes no delete API, so destroying the resource only removes it from Terraform state. Every input is ForceNew, so changing any argument recreates the message.

## Import

PAIPlugin Message can be imported using the id, e.g.

```shell
terraform import alicloud_pai_plugin_message.default <message_id>
```

-> **NOTE:** Import without a captured send date cannot re-read the message status, since ListMessages requires a Datetime filter. The resource is best managed via create.
