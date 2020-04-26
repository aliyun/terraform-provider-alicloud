---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_audit"
sidebar_current: "docs-alicloud-resource-log-audit"
description: |-
Provides a Alicloud log audit resource.
---

# alicloud\_log\_audit

SLS log audit exists in the form of log service app.

In addition to inheriting all SLS functions, it also enhances the real-time automatic centralized collection of audit related logs across multi cloud products under multi accounts, and provides support for storage, query and information summary required by audit. It covers actiontrail, OSS, NAS, SLB, API gateway, RDS, WAF, cloud firewall, cloud security center and other products.

-> **NOTE:** Available in 1.81.0

## Example Usage

Basic Usage

```
resource "alicloud_log_audit" "example" {
    display_name = "tf-audit-test"
    aliuid = "12345678"
    variable_map = {
        "actiontrail_enabled" = "true",
        "actiontrail_ttl" = "180",
        "oss_access_enabled" = "true",
        "oss_access_ttl" = "180",
    }
}
```
Multiple accounts Usage

```
resource "alicloud_log_audit" "example" {
    display_name = "tf-audit-test"
    aliuid = "12345678"
    variable_map = {
        "actiontrail_enabled" = "true",
        "actiontrail_ttl" = "180",
        "oss_access_enabled" = "true",
        "oss_access_ttl" = "180",
    }
    multi_account = ["123456789123", "12345678912300123"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required, ForceNew) Name of SLS log audit.
* `aliuid` - (Required, ForceNew) Aliuid value of your account.
* `variable_map` - (Required)  Log audit detailed configuration.
    * `actiontrail_enabled` - (Optional) Notification type. support Email, SMS, DingTalk. Default true.
    * `actiontrail_ttl` - (Optional) Actiontril action log TTL. Default 180.
    * `oss_access_enabled` - (Optional). Access log switch of OSS. Default true.
    * `oss_access_ttl` - (Optional). Access log TTL of OSS. Default 180.
    * `oss_sync_enabled` - (Optional).OSS synchronization to central configuration switch. Default true.
    * `oss_sync_ttl` - (Optional).OSS synchronization to central TTL. Default 180.
    * `oss_metering_enabled` - (Optional). OSS metering log switch.Default true.
    * `oss_metering_ttl` - (Optional). OSS measurement log TTL. Default 180.
    * `rds_enabled` - (Optional). RDS audit log switch. Default true.
    * `rds_ttl` - (Optional). Dds log centralization ttl.Default 180.
    * `slb_access_enabled` - (Optional). Slb log switch. Default true.
    * `slb_access_ttl` - (Optional). Slb centralized ttl. Default 180.
    * `slb_sync_enabled` - (Optional). Slb sync to center switch. Default true.
    * `slb_sync_ttl` - (Optional). Slb sync to center switch。Default 180.
    * `bastion_enabled` - (Optional). Fortress machine operation log switch.Default true.
    * `bastion_ttl` - (Optional). Fort machine centralized ttl. Default 180.
    * `waf_enabled` - (Optional). Waf log switch .Default true.
    * `waf_ttl` - (Optional). Waf centralized ttl.Default true.
    * `cloudfirewall_ttl` - (Optional). Cloud firewall switch.Default true.
    * `sas_ttl` - (Optional).Cloud Security Center centralized ttl. Default 180.
    * `sas_process_enabled` - (Optional).Cloud Security Center process startup log switch. Default false.
    * `sas_network_enabled` - (Optional).Cloud security center network connection log switch. Default false.
    * `sas_login_enabled` - (Optional).Cloud security center login flow log switch. Default false.
    * `sas_crack_enabled` - (Optional).Cloud Security Center Brute Force Log Switch. Default false.
    * `sas_snapshot_process_enabled` - (Optional). Cloud Security Center process snapshot switch. Default false.
    * `sas_snapshot_account_enabled` - (Optional). Cloud Security Center account snapshot switch.Default false.
    * `sas_snapshot_port_enabled` - (Optional).Cloud Security Center Port Snapshot Switch. Default false.
    * `sas_dns_enabled` - (Optional).Cloud Security Center DNS resolution log switch. Default false.
    * `sas_local_dns_enabled` - (Optional).Cloud security center local DNS log switch. Default false.
    * `sas_session_enabled` - (Optional). Cloud security center network session log switch.Default false.
    * `sas_http_enabled` - (Optional). Cloud Security Center WEB access log switch.Default false.
    * `sas_security_vul_enabled` - (Optional). Cloud Security Center Vulnerability Log Switch.Default false.
    * `sas_security_hc_enabled` - (Optional).Cloud Security Center Baseline Log Switch. Default false.
    * `sas_security_alert_enabled` - (Optional). Cloud Security Center Security Alarm Log Switch .Default false.
    * `apigateway_enabled` - (Optional). API Gateway Log Switch.Default true.
    * `apigateway_ttl` - (Optional). API Gateway ttl.Default 180.
    * `nas_enabled` - (Optional). Nas log switch.Default true.
    * `nas_ttl` - (Optional). Nas centralized ttl. Default 180.
    * `cps_enabled` - (Optional). Mobile push log switch.Default true.
    * `cps_ttl` - (Optional).Mobile push ttl. Default 180.
* `multi_account` - (Optional).Multi-account configuration, please fill in multiple aliuid.
                

## Attributes Reference

The following attributes are exported:

*  `id` - The ID of the log audit. It formats of `display_name`.

## Import

Log alert can be imported using the id, e.g.

```
$ terraform import alicloud_log_audit.example tf-audit-test
```