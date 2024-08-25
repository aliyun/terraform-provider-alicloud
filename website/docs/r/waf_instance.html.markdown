---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_instance"
sidebar_current: "docs-alicloud-resource-waf-instance"
description: |-
  Provides a Web Application Firewall Instance resource.
---

# alicloud\_waf\_instance

-> **DEPRECATED:**  This resource has been deprecated and using [alicloud_wafv3_instance](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/wafv3_instance) instead.

Provides a WAF Instance resource to create instance in the Web Application Firewall.

For information about WAF and how to use it, see [What is Alibaba Cloud WAF](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available in 1.83.0+ .

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_waf_instance&exampleId=bb172d48-8404-035c-b38d-686ce2b74420b94ffc1b&activeTab=example&spm=docs.r.waf_instance.0.bb172d4884&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_waf_instances" "default" {}
resource "alicloud_waf_instance" "default" {
  count                = length(data.alicloud_waf_instances.default.instances) > 0 ? 0 : 1
  big_screen           = "0"
  exclusive_ip_package = "1"
  ext_bandwidth        = "50"
  ext_domain_package   = "1"
  package_code         = "version_3"
  prefessional_service = "false"
  subscription_type    = "Subscription"
  period               = 1
  waf_log              = "false"
  log_storage          = "3"
  log_time             = "180"
  resource_group_id    = "rs-abc12345"
}
```
## Argument Reference

The following arguments are supported:

* `big_screen` - (Required, String) Specify whether big screen is supported. Valid values: ["0", "1"]. "0" for false and "1" for true.
* `exclusive_ip_package` - (Required, String) Specify the number of exclusive WAF IP addresses.
* `ext_bandwidth` - (Required, String) The extra bandwidth. Unit: Mbit/s.
* `ext_domain_package` - (Required, String) The number of extra domains.
* `log_storage` - (Required, String) Log storage size. Unit: T. Valid values: [3, 5, 10, 20, 50].
* `log_time` - (Required, String) Log storage period. Unit: day. Valid values: [180, 360].
* `modify_type` - (Optional) Type of configuration change. Valid value: Upgrade.
* `package_code` - (Required, String) Subscription plan:
    * China site customers can purchase the following versions of China Mainland region, valid values: ["version_3", "version_4", "version_5"].
    * China site customers can purchase the following versions of International region, valid values: ["version_pro_asia", "version_business_asia", "version_enterprise_asia"]
    * International site customers can purchase the following versions of China Mainland region: ["version_pro_china", "version_business_china", "version_enterprise_china"]
    * International site customers can purchase the following versions of International region: ["version_pro", "version_business", "version_enterprise"].

* `period` - (ForceNew) Service time of Web Application Firewall.
* `prefessional_service` - (Required, String) Specify whether professional service is supported. Valid values: ["true", "false"]
* `renew_period` - (ForceNew) Renewal period of WAF service. Unit: month
* `renewal_status` - (ForceNew) Renewal status of WAF service. Valid values: 
    * AutoRenewal: The service time of WAF is renewed automatically.
    * ManualRenewal (default): The service time of WAF is renewed manually.Specifies whether to configure a Layer-7 proxy, such as Anti-DDoS Pro or CDN, to filter the inbound traffic before it is forwarded to WAF. Valid values: "On" and "Off". Default to "Off".
* `resource_group_id` - (Optional) The resource group ID.
* `region` - (Optional, Available in 1.139.0+) The instance region ID.
* `subscription_type` - (Required, String) Subscription of WAF service. Valid values: ["Subscription", "PayAsYouGo"].
* `waf_log` - (Required, String) Specify whether Log service is supported. Valid values: ["true", "false"]                                           
			
## Attributes Reference

The following attributes are exported:

* `id` - This resource instance id.
* `status` - The status of the instance.

## Import

WAF instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_waf_instance.default waf-cn-132435
```
