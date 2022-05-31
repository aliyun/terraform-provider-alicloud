data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

data "alicloud_alidns_gtm_instances" "default" {}

resource "alicloud_alidns_gtm_instance" "default" {
  count                   = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? 0 : 1
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

locals {
  gtm_instance_id = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? data.alicloud_alidns_gtm_instances.default.ids[0] : concat(alicloud_alidns_gtm_instance.default.*.id, [""])[0]
}

resource "alicloud_alidns_address_pool" "ipv4" {
  count             = 2
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
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

resource "alicloud_alidns_access_strategy" "default" {
  strategy_name                  = var.name
  strategy_mode                  = "GEO"
  instance_id                    = local.gtm_instance_id
  default_addr_pool_type         = "IPV4"
  default_lba_strategy           = "RATIO"
  default_min_available_addr_num = 1
  default_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.0.id
  }
  failover_addr_pool_type         = "IPV4"
  failover_lba_strategy           = "RATIO"
  failover_min_available_addr_num = 1
  failover_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.1.id
  }
  lines {
    line_code = "default"
  }
}

