---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_gtm_instances"
sidebar_current: "docs-alicloud-datasource-alidns-gtm-instances"
description: |-
  Provides a list of Alidns Gtm Instances to the user.
---

# alicloud\_alidns\_gtm\_instances

This data source provides the Alidns Gtm Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alidns_gtm_instances" "ids" {}
output "alidns_gtm_instance_id_1" {
  value = data.alicloud_alidns_gtm_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Gtm Instance IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).
* `lang` - (Optional, ForceNew) The lang.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Alidns Gtm Instances. Each element contains the following attributes:
    * `alert_group` - The alert group.
    * `alert_config` - The alert notification methods.
      * `dingtalk_notice` - Whether to configure DingTalk notifications.
      * `email_notice` -  Whether to configure mail notification.
      * `notice_type` - The Alarm Event Type. 
      * `sms_notice` - Whether to configure SMS notification.
    * `cname_type` - The access type of the CNAME domain name.
    * `create_time` - The CreateTime of the Gtm Instance.
    * `expire_time` - The ExpireTime of the Gtm Instance.
    * `id` - The ID of the Gtm Instance.
    * `instance_id` - The ID of the Gtm Instance.
    * `instance_name` - The name of the Gtm Instance.
    * `strategy_mode` - The type of the access policy.
    * `payment_type` - The paymentype of the resource.
    * `public_cname_mode` - The Public Network domain name access method.
    * `public_rr` - The CNAME access domain name.
    * `public_user_domain_name` - The website domain name that the user uses on the Internet.
    * `public_zone_name` - The domain name that is used to access GTM over the Internet.
    * `resource_group_id` - The ID of the resource group.
    * `ttl` - The global time to live.
    * `package_edition` - The version of the instance.