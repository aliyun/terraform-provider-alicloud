---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_audit"
sidebar_current: "docs-alicloud-resource-log-audit"
description: |-
  Provides a Alicloud log audit resource.
---

# alicloud_log_audit

SLS log audit exists in the form of log service app.

In addition to inheriting all SLS functions, it also enhances the real-time automatic centralized collection of audit related logs across multi cloud products under multi accounts, and provides support for storage, query and information summary required by audit. It covers actiontrail, OSS, NAS, SLB, API gateway, RDS, WAF, cloud firewall, cloud security center and other products.

-> **NOTE:** Available since v1.81.0

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_audit&exampleId=771f99a8-11b2-0457-b1eb-811fa59251c3441793f0&activeTab=example&spm=docs.r.log_audit.0.771f99a811&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}
resource "alicloud_log_audit" "example" {
  display_name = "tf-audit-example"
  aliuid       = data.alicloud_account.default.id
  variable_map = {
    "actiontrail_enabled"             = "true",
    "actiontrail_ttl"                 = "180",
    "oss_access_enabled"              = "true",
    "oss_access_ttl"                  = "7",
    "oss_sync_enabled"                = "true",
    "oss_sync_ttl"                    = "180",
    "oss_metering_enabled"            = "true",
    "oss_metering_ttl"                = "180",
    "rds_enabled"                     = "true",
    "rds_audit_collection_policy"     = "",
    "rds_ttl"                         = "180",
    "rds_slow_enabled"                = "false",
    "rds_slow_collection_policy"      = "",
    "rds_slow_ttl"                    = "180",
    "rds_perf_enabled"                = "false",
    "rds_perf_collection_policy"      = "",
    "rds_perf_ttl"                    = "180",
    "vpc_flow_enabled"                = "false",
    "vpc_flow_ttl"                    = "7",
    "vpc_flow_collection_policy"      = "",
    "vpc_sync_enabled"                = "true",
    "vpc_sync_ttl"                    = "180",
    "polardb_enabled"                 = "true",
    "polardb_audit_collection_policy" = "",
    "polardb_ttl"                     = "180",
    "polardb_slow_enabled"            = "false",
    "polardb_slow_collection_policy"  = "",
    "polardb_slow_ttl"                = "180",
    "polardb_perf_enabled"            = "false",
    "polardb_perf_collection_policy"  = "",
    "polardb_perf_ttl"                = "180",
    "drds_audit_enabled"              = "true",
    "drds_audit_collection_policy"    = "",
    "drds_audit_ttl"                  = "7",
    "drds_sync_enabled"               = "true",
    "drds_sync_ttl"                   = "180",
    "slb_access_enabled"              = "true",
    "slb_access_collection_policy"    = "",
    "slb_access_ttl"                  = "7",
    "slb_sync_enabled"                = "true",
    "slb_sync_ttl"                    = "180",
    "bastion_enabled"                 = "true",
    "bastion_ttl"                     = "180",
    "waf_enabled"                     = "true",
    "waf_ttl"                         = "180",
    "cloudfirewall_enabled"           = "true",
    "cloudfirewall_ttl"               = "180",
    "ddos_coo_access_enabled"         = "false",
    "ddos_coo_access_ttl"             = "180",
    "ddos_bgp_access_enabled"         = "false",
    "ddos_bgp_access_ttl"             = "180",
    "ddos_dip_access_enabled"         = "false",
    "ddos_dip_access_ttl"             = "180",
    "sas_crack_enabled"               = "true",
    "sas_dns_enabled"                 = "true",
    "sas_http_enabled"                = "true",
    "sas_local_dns_enabled"           = "true",
    "sas_login_enabled"               = "true",
    "sas_network_enabled"             = "true",
    "sas_process_enabled"             = "true",
    "sas_security_alert_enabled"      = "true",
    "sas_security_hc_enabled"         = "true",
    "sas_security_vul_enabled"        = "true",
    "sas_session_enabled"             = "true",
    "sas_snapshot_account_enabled"    = "true",
    "sas_snapshot_port_enabled"       = "true",
    "sas_snapshot_process_enabled"    = "true",
    "sas_ttl"                         = "180",
    "apigateway_enabled"              = "true",
    "apigateway_ttl"                  = "180",
    "nas_enabled"                     = "true",
    "nas_ttl"                         = "180",
    "appconnect_enabled"              = "false",
    "appconnect_ttl"                  = "180",
    "cps_enabled"                     = "true",
    "cps_ttl"                         = "180",
    "k8s_audit_enabled"               = "true",
    "k8s_audit_collection_policy"     = "",
    "k8s_audit_ttl"                   = "180",
    "k8s_event_enabled"               = "true",
    "k8s_event_collection_policy"     = "",
    "k8s_event_ttl"                   = "180",
    "k8s_ingress_enabled"             = "true",
    "k8s_ingress_collection_policy"   = "",
    "k8s_ingress_ttl"                 = "180"
  }
}
```
Multiple accounts Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_audit&exampleId=4ef78322-06f0-bf07-151a-03b3bba82aa2c9eb6dfb&activeTab=example&spm=docs.r.log_audit.1.4ef7832206&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}

resource "alicloud_log_audit" "example" {
  display_name = "tf-audit-example"
  aliuid       = data.alicloud_account.default.id
  variable_map = {
    "actiontrail_enabled" = "true",
    "actiontrail_ttl"     = "180",
    "oss_access_enabled"  = "true",
    "oss_access_ttl"      = "180",
  }
  multi_account = ["123456789123", "12345678912300123"]
}
```
Resource Directory Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_audit&exampleId=52d8f305-0eaf-1a81-2b89-9c96539bff9f43d4c2ed&activeTab=example&spm=docs.r.log_audit.2.52d8f3050e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}

resource "alicloud_log_audit" "example" {
  display_name = "tf-audit-example"
  aliuid       = data.alicloud_account.default.id
  variable_map = {
    "actiontrail_enabled" = "true",
    "actiontrail_ttl"     = "180",
    "oss_access_enabled"  = "true",
    "oss_access_ttl"      = "180",
  }
  resource_directory_type = "all"
}
```
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_audit&exampleId=9844a14c-5cfc-9781-e1fa-03defa44c0cfba3dbae6&activeTab=example&spm=docs.r.log_audit.3.9844a14c5c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}

resource "alicloud_log_audit" "example" {
  display_name = "tf-audit-example"
  aliuid       = data.alicloud_account.default.id
  variable_map = {
    "actiontrail_enabled" = "true",
    "actiontrail_ttl"     = "180",
    "oss_access_enabled"  = "true",
    "oss_access_ttl"      = "180",
  }
  multi_account           = [] //Put your member accounts here, separated by ","
  resource_directory_type = "custom"
}
```
## Argument Reference

The following arguments are supported:

* `display_name` - (Required, ForceNew) Name of SLS log audit.
* `aliuid` - (Required, ForceNew) Aliuid value of your account.
* `variable_map` - (Optional) Log audit detailed configuration.
    * `actiontrail_enabled` - (Optional) Actiontrail action log switch. Default false.
    * `actiontrail_ttl` - (Optional) Actiontril action log TTL. Default 180.
    
    - `oss_access_enabled` - (Optional) Access log switch of OSS. Default false.
    
    - `oss_access_ttl` - (Optional) Regional Access log TTL of OSS. Default 7.
    - `oss_sync_enabled` - (Optional) OSS synchronization to central configuration switch. Default true.
    
    - `oss_sync_ttl` - (Optional) OSS synchronization to central TTL. Default 180.
    
    - `oss_metering_enabled` - (Optional) OSS metering log switch.Default false.
    - `oss_metering_ttl` - (Optional) OSS measurement log TTL. Default 180.
    
    - `rds_enabled` - (Optional) RDS audit log switch. Default false.
    - `rds_audit_collection_policy` - (Optional) RDS audit log collection policy script. Default empty.
    
    - `rds_ttl` - (Optional) RDS audit log ttl. Default 180.
    
    - `rds_slow_enabled` -  (Optional) RDS slow log switch. Default false.
    - `rds_slow_collection_policy` - (Optional) RDS slow log collection policy script. Default empty.
    
    - `rds_slow_ttl` - (Optional) RDS slow log TTL. Default 180.
    - `rds_perf_enabled` -  (Optional) RDS performance log switch. Default false.
    
    - `rds_perf_collection_policy` - (Optional) RDS performance log collection policy script. Default empty.
    - `rds_perf_ttl` - (Optional) RDS performance log TTL. Default 180.
    - `vpc_flow_enabled` - (Optional) Flow log of VPC. Default false.
    - `vpc_flow_ttl` - (Optional) Regional flow log TTL of VPC. Default 7.
    - `vpc_flow_collection_policy` - (Optional) VPC flow log collection policy script. Default empty.
    - `vpc_sync_enabled` - (Optional) VPC synchronization to central configuration switch. Default true.
    - `vpc_sync_ttl` - (Optional) VPC synchronization to central TTL. Default 180.
    - `dns_intranet_enabled` - (Optional) Specifies whether to collect intranet Alibaba Cloud DNS (DNS) logs. Default false.
    - `dns_intranet_ttl` - (Optional) The retention period of intranet DNS logs in the regional Logstore. Default 7.
    - `dns_intranet_collection_policy` - (Optional) DNS intranet log collection policy script. Default empty.
    - `dns_sync_enabled` - (Optional) Specifies whether to synchronize DNS intranet logs to the central project. Default true.
    - `dns_sync_ttl` - (Optional) The retention period of intranet DNS logs in the central Logstore. Default 180.
    - `polardb_enabled` - (Optional) PolarDB audit log switch. Default false.
    - `polardb_audit_collection_policy` - (Optional) PolarDB audit log collection policy script. Default empty.
    
    - `polardb_ttl` - (Optional) PolarDB audit log ttl. Default 180.
    
    - `polardb_slow_enabled` -  (Optional) PolarDB slow log switch. Default false.
    - `polardb_slow_collection_policy` - (Optional) PolarDB slow log collection policy script. Default empty.
    
    - `polardb_slow_ttl` - (Optional) PolarDB slow log TTL. Default 180.
    - `polardb_perf_enabled` -  (Optional) PolarDB performance log switch. Default false.
    
    - `polardb_perf_collection_policy` - (Optional) PolarDB performance log collection policy script. Default empty.
    - `polardb_perf_ttl` - (Optional) PolarDB performance log TTL. Default 180.
    
    - `drds_audit_enabled` - (Optional) PolarDB-X audit log switch. Default false.
    - `rds_audit_collection_policy` - (Optional) PolarDB-X  audit log collection policy script. Default empty.
    
    - `drds_audit_ttl` - (Optional) Regional PolarDB-X  audit log ttl. Default 7.
    - `drds_sync_enabled` - (Optional) PolarDB-X synchronization to central configuration switch. Default true.
    
    - `drds_sync_ttl` - (Optional) PolarDB-X synchronization to central TTL. Default 180.
    
    - `slb_access_enabled` - (Optional) Slb log switch. Default false.
    - `slb_access_collection_policy` - (Optional) Slb log collection policy script. Default empty.
    
    - `slb_access_ttl` - (Optional) Regional Slb access log ttl. Default 7.
    - `slb_sync_enabled` - (Optional) Slb sync to center switch. Default true.
    
    - `slb_sync_ttl` - (Optional) Slb sync to center switch. Default 180.
    
    - `bastion_enabled` - (Optional) Fortress machine operation log switch.Default false.
    - `bastion_ttl` - (Optional) Fortress machine centralized ttl. Default 180.
    
    - `waf_enabled` - (Optional) Waf log switch. Default false.
    
    - `waf_ttl` - (Optional) Waf centralized ttl. Default 180.
    
    - `cloudfirewall_enabled` - (Optional) Cloudfirewall log switch. Default false.
    - `cloudfirewall_ttl` - (Optional) Cloudfirewall log ttl.Default 180.
    
    - `cloudfirewall_vpc_enabled` - (Optional) Specifies whether to collect VPC firewall traffic logs from Cloud Firewall. Default false.
    - `cloudfirewall_vpc_ttl` - (Optional)The retention period of Cloud Firewall VPC firewall traffic logs in the central Logstore. Default 180.
    
    - `ddos_coo_access_enabled` - (Optional) Anti-DDoS Pro(New BGP) access log switch. Default false.
    
    - `ddos_coo_access_ttl` - (Optional) Anti-DDoS Pro (New BGP) access log ttl. Default 180.
    - `ddos_bgp_access_enabled` - (Optional) Anti-DDoS (Origin) access log switch. Default false.
    - `ddos_bgp_access_ttl` - (Optional) Anti-DDoS (Origin) access log ttl. Default 180.
    - `ddos_dip_access_enabled` - (Optional) Anti-DDoS Premium access log switch. Default false.
    
    - `ddos_dip_access_ttl` - (Optional) Anti-DDoS Premium access log ttl. Default 180.
    - `sas_ttl` - (Optional) Cloud Security Center centralized ttl. Default 180.
    - `sas_process_enabled` - (Optional) Cloud Security Center process startup log switch. Default false.
    
    - `sas_network_enabled` - (Optional) Cloud security center network connection log switch. Default false.
    - `sas_login_enabled` - (Optional) Cloud security center login flow log switch. Default false.
    
    - `sas_crack_enabled` - (Optional) Cloud Security Center Brute Force Log Switch. Default false.
    - `sas_snapshot_process_enabled` - (Optional) Cloud Security Center process snapshot switch. Default false.
    
    - `sas_snapshot_account_enabled` - (Optional) Cloud Security Center account snapshot switch. Default false.
    - `sas_snapshot_port_enabled` - (Optional) Cloud Security Center Port Snapshot Switch. Default false.
    
    - `sas_dns_enabled` - (Optional) Cloud Security Center DNS resolution log switch. Default false.
    - `sas_local_dns_enabled` - (Optional) Cloud security center local DNS log switch. Default false.
    
    - `sas_session_enabled` - (Optional) Cloud security center network session log switch.Default false.
    - `sas_http_enabled` - (Optional). Cloud Security Center WEB access log switch. Default false.
    
    - `sas_security_vul_enabled` - (Optional) Cloud Security Center Vulnerability Log Switch.Default false.
    - `sas_security_hc_enabled` - (Optional) Cloud Security Center Baseline Log Switch. Default false.
    
    - `sas_security_alert_enabled` - (Optional) Cloud Security Center Security Alarm Log Switch .Default false.
    
    - `apigateway_enabled` - (Optional) API Gateway Log Switch. Default false.
    - `apigateway_ttl` - (Optional) API Gateway ttl. Default 180.
    
    - `nas_enabled` - (Optional) Nas log switch. Default false.
    
    - `nas_ttl` - (Optional) Nas centralized ttl. Default 180.
    
    - `appconnect_enabled` - (Optional) App Connect operation log switch. Default false.
    - `appconnect_ttl` - (Optional) App Connect operation log ttl. Default 180.
    
    - `cps_enabled` - (Optional) Mobile push log switch. Default false.
    
    - `cps_ttl` - (Optional) Mobile push ttl. Default 180.
    
    - `k8s_audit_enabled` - (Optional) K8s audit log switch. Default false.
    - `k8s_audit_collection_policy` - (Optional) K8s audit log collection policy script. Default empty.
    
    - `k8s_audit_ttl` - (Optional) K8s audit log ttl. Default 180.
    - `k8s_event_enabled` - (Optional) K8s event log switch. Default false.
    
    - `k8s_event_collection_policy` - (Optional) K8s event log collection policy script. Default empty.
    - `k8s_event_ttl` - (Optional) K8s event log ttl. Default 180.
    
    - `k8s_ingress_enabled` - (Optional) K8s ingress log switch. Default false.
    - `k8s_ingress_collection_policy` - (Optional) K8s ingress log collection policy script. Default empty.
    
    - `k8s_ingress_ttl` - (Optional) K8s ingress log ttl. Default 180.

    - `idaas_mng_enabled` - (Optional) Specifies whether to collect IDaaS management log. Default false.
    - `idaas_mng_ttl` - (Optional)   IDaaS management log TTL. Default 180.
    - `idaas_mng_collection_policy` - (Optional) IDaaS management log collection policy script. Default empty.

    - `idaas_user_enabled` - (Optional) Specifies whether to collect IDaaS user behavior log. Default false.
    - `idaas_user_ttl` - (Optional) IDaaS user behavior log. Default 180.
    - `idaas_user_collection_policy` - (Optional) IDaaS user behavior log collection policy script. Default empty.
  
* `multi_account` - (Optional) Multi-account configuration, please fill in multiple aliuid.
* `resource_directory_type` - (Optional, Available in 1.135.0+) Resource Directory type. Optional values are all or custom. If the value is custom, argument multi_account should be provided.
                

## Attributes Reference

The following attributes are exported:

*  `id` - The ID of the log audit. It formats of `display_name`.

## Import

-> **NOTE:** The UI settings of collection policy scripts for related products (such as rds, slb and etc.) will be cleared when imported using terraform. So you need to modify collection policy scripts directly if you want to edit collection policy in terraform.

Log audit can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_audit.example tf-audit-example
```
