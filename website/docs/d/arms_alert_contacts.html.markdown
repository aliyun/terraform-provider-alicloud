---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_alert_contacts"
sidebar_current: "docs-alicloud-datasource-arms-alert-contacts"
description: |-
  Provides a list of Arms Alert Contacts to the user.
---

# alicloud\_arms\_alert\_contacts

This data source provides the Arms Alert Contacts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_alert_contacts" "ids" {}
output "arms_alert_contact_id_1" {
  value = data.alicloud_arms_alert_contacts.ids.contacts.0.id
}

data "alicloud_arms_alert_contacts" "nameRegex" {
  name_regex = "^my-AlertContact"
}
output "arms_alert_contact_id_2" {
  value = data.alicloud_arms_alert_contacts.nameRegex.contacts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `alert_contact_name` - (Optional, ForceNew) The name of the alert contact.
* `email` - (Optional, ForceNew) The email address of the alert contact.
* `ids` - (Optional, ForceNew, Computed)  A list of Alert Contact IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Alert Contact name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `phone_num` - (Optional, ForceNew) The mobile number of the alert contact. 

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Alert Contact names.
* `contacts` - A list of Arms Alert Contacts. Each element contains the following attributes:
	* `alert_contact_id` - Contact ID.
	* `alert_contact_name` - The name of the alert contact.
	* `create_time` - The Creation Time Timestamp.
	* `ding_robot_webhook_url` - The webhook URL of the DingTalk chatbot. 
	* `email` - The email address of the alert contact. 
	* `id` - The ID of the Alert Contact.
	* `phone_num` - The mobile number of the alert contact. 
	* `system_noc` - Specifies whether the alert contact receives system notifications. 
	* `webhook` - Webhook Information.
