---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_firewall_rule"
sidebar_current: "docs-alicloud-resource-simple-application-server-firewall-rule"
description: |-
  Provides a Alicloud Simple Application Server Firewall Rule resource.
---

# alicloud_simple_application_server_firewall_rule

Provides a Simple Application Server Firewall Rule resource.

For information about Simple Application Server Firewall Rule and how to use it, see [What is Firewall Rule](https://www.alibabacloud.com/help/doc-detail/190449.htm).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_simple_application_server_firewall_rule&exampleId=0e591dba-527a-0e06-b288-536cd17677570be13048&activeTab=example&spm=docs.r.simple_application_server_firewall_rule.0.0e591dba52&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}

data "alicloud_simple_application_server_images" "default" {}
data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

resource "alicloud_simple_application_server_firewall_rule" "default" {
  instance_id   = alicloud_simple_application_server_instance.default.id
  rule_protocol = "Tcp"
  port          = "9999"
  remark        = var.name
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Alibaba Cloud simple application server instance ID.
* `port` - (Required, ForceNew) The port range. Valid values of port numbers: `1` to `65535`. Specify a port range in the format of `<start port number>/<end port number>`. Example: `1024/1055`, which indicates the port range of `1024` through `1055`.
* `remark` - (Optional, ForceNew) The remarks of the firewall rule.
* `rule_protocol` - (Required, ForceNew) The transport layer protocol. Valid values: `Tcp`, `Udp`, `TcpAndUdp`.

## Attributes Reference

The following attributes are exported:

* `firewall_rule_id` - The ID of the firewall rule.
* `id` - The resource ID of Firewall Rule. The value formats as `<instance_id>:<firewall_rule_id>`.

## Import

Simple Application Server Firewall Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_firewall_rule.example <instance_id>:<firewall_rule_id>
```