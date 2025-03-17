---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_log_meta"
description: |-
  Provides a Alicloud Threat Detection Log Meta resource.
---

# alicloud_threat_detection_log_meta

Provides a Threat Detection Log Meta resource.

Log analysis shipping status.

For information about Threat Detection Log Meta and how to use it, see [What is Log Meta](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ModifyLogMetaStatus).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_log_meta&exampleId=8a8717ed-1051-cbcf-a183-a30f1e2b0597497b2a89&activeTab=example&spm=docs.r.threat_detection_log_meta.0.8a8717ed10&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}


resource "alicloud_threat_detection_log_meta" "default" {
  status        = "disabled"
  log_meta_name = "aegis-log-client"
}
```

### Deleting `alicloud_threat_detection_log_meta` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_log_meta`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `log_meta_name` - (Required, ForceNew) The name of the exclusive Logstore where logs are stored. Value:
  - aegis-log-client: client event log
  - aegis-log-crack: Brute Force log
  - aegis-log-dns-query:DNS request log
  - aegis-log-login: login log
  - aegis-log-network: network connection log
  - aegis-log-process: process startup log
  - aegis-snapshot-host: account snapshot log
  - aegis-snapshot-port: port snapshot log
  - aegis-snapshot-process: process snapshot log
  - local-dns: local DNS log
  - sas-log-dns:DNS resolution log
  - sas-log-http:WEB access log
  - sas-log-session: Web session log
  - sas-security-log: alarm log
  - sas-vul-log: Vulnerability log
  - sas-cspm-log: Cloud platform configuration check log
  - sas-hc-log: baseline log
  - sas-rasp-log: Application Protection Log
  - sas-filedetect-log: file detection log
  - sas-net-block: Network Defense Log
* `status` - (Required) The status of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Log Meta.
* `update` - (Defaults to 5 mins) Used when update the Log Meta.

## Import

Threat Detection Log Meta can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_log_meta.example <id>
```