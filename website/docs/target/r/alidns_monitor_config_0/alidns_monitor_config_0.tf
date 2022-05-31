variable "name" {
  default = "tf-testacc"
}

variable "domain_name" {
  default = "your_domain_name"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}

resource "alicloud_alidns_address_pool" "default" {
  address_pool_name = var.name
  instance_id       = alicloud_alidns_gtm_instance.default.id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}

resource "alicloud_alidns_monitor_config" "default" {
  addr_pool_id        = alicloud_alidns_address_pool.default.id
  evaluation_count    = "1"
  interval            = "60"
  timeout             = "5000"
  protocol_type       = "TCP"
  monitor_extend_info = "{\"failureRate\"=50,\"port\"=80}"
  isp_city_node = {
    city_code = "503"
    isp_code  = "465"
  }
}
