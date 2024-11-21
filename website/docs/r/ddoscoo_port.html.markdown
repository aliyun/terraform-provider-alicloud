---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_port"
description: |-
  Provides a Alicloud Ddos Coo Port resource.
---

# alicloud_ddoscoo_port

Provides a Ddos Coo Port resource.


For information about Anti-DDoS Pro Port and how to use it, see [What is Port](https://www.alibabacloud.com/help/en/ddos-protection/latest/api-ddoscoo-2020-01-01-createport).

-> **NOTE:** Available since v1.123.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddoscoo_port&exampleId=682d6caf-f018-3608-ef78-77710086b53b734b4057&activeTab=example&spm=docs.r.ddoscoo_port.0.682d6caff0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
resource "alicloud_ddoscoo_instance" "default" {
  name              = var.name
  bandwidth         = "30"
  base_bandwidth    = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
  period            = "1"
  product_type      = "ddoscoo"
}

resource "alicloud_ddoscoo_port" "default" {
  instance_id       = alicloud_ddoscoo_instance.default.id
  frontend_port     = "7001"
  backend_port      = "7002"
  frontend_protocol = "tcp"
  real_servers      = ["1.1.1.1", "2.2.2.2"]
}
```

## Argument Reference

The following arguments are supported:
* `backend_port` - (Optional, ForceNew) The port of the origin server. Valid values: `0` to `65535`.

* `config` - (Optional, List, Available since v1.230.0) Session persistence settings for port forwarding rules. Use a string representation in JSON format. The specific structure is described as follows.
  - `PersistenceTimeout`: is of Integer type and is required. The timeout period of the session. Value range: `30` to `3600`, in seconds. The default value is `0`, which is closed. See [`config`](#config) below.
* `frontend_port` - (Required, ForceNew, Int) The forwarding port to query. Valid values: `0` to `65535`.

* `frontend_protocol` - (Required, ForceNew) The type of the forwarding protocol to query. Valid values:
  - `tcp`
  - `udp`

* `instance_id` - (Required, ForceNew) The ID of the Anti-DDoS Pro or Anti-DDoS Premium instance to which the port forwarding rule belongs.

-> **NOTE:**  You can call the [DescribeInstanceIds](https://www.alibabacloud.com/help/en/doc-detail/157459.html) operation to query the IDs of all instances.

* `real_servers` - (Required, List) List of source IP addresses

### `config`

The config supports the following:
* `persistence_timeout` - (Optional, Int, Available since v1.230.0) The timeout period for session retention. Value range: 30~3600, unit: second. The default is 0, which means off.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<frontend_port>:<frontend_protocol>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Port.
* `delete` - (Defaults to 5 mins) Used when delete the Port.
* `update` - (Defaults to 5 mins) Used when update the Port.

## Import

Ddos Coo Port can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_port.example <instance_id>:<frontend_port>:<frontend_protocol>
```