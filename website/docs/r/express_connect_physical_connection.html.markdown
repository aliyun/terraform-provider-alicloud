---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_physical_connection"
sidebar_current: "docs-alicloud-resource-express-connect-physical-connection"
description: |-
  Provides a Alicloud Express Connect Physical Connection resource.
---

# alicloud_express_connect_physical_connection

Provides a Express Connect Physical Connection resource.

For information about Express Connect Physical Connection and how to use it, see [What is Physical Connection](https://www.alibabacloud.com/help/en/express-connect/developer-reference/api-vpc-2016-04-28-createphysicalconnection-efficiency-channels).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_physical_connection&exampleId=14b33bb7-2a9a-34e6-d4ce-3c91fc5111636f98acb5&activeTab=example&spm=docs.r.express_connect_physical_connection.0.14b33bb72a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
  alias  = "hz"
}

provider "alicloud" {
  region = "ap-southeast-1"
  alias  = "sgp"
}

resource "alicloud_express_connect_physical_connection" "domestic" {
  access_point_id          = "ap-cn-hangzhou-yh-B"
  line_operator            = "CT"
  peer_location            = "example_value"
  physical_connection_name = "example_value"
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
  provider                 = alicloud.hz
}

resource "alicloud_express_connect_physical_connection" "international" {
  access_point_id          = "ap-sg-singpore-A"
  line_operator            = "Other"
  peer_location            = "example_value"
  physical_connection_name = "example_value"
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
  provider                 = alicloud.sgp
}
```

## Argument Reference

The following arguments are supported:

* `access_point_id` - (Required, ForceNew) The access point ID of the Express Connect circuit.
* `line_operator` - (Required) The connectivity provider of the Express Connect circuit. Valid values:
  - `CT`: China Telecom.
  - `CU`: China Unicom.
  - `CM`: China Mobile.
  - `CO`: Other connectivity providers in the Chinese mainland.
  - `Equinix`: Equinix.
  - `Other`: Other connectivity providers outside the Chinese mainland.
* `type` - (Optional, ForceNew) The type of Express Connect circuit. Default value: `VPC`. Valid values: `VPC`.
* `port_type` - (Optional, ForceNew) The port type of the Express Connect circuit. Valid values:
  - `100Base-T`: 100 Mbit/s copper Ethernet port.
  - `1000Base-T`: 1000 Mbit/s copper Ethernet port.
  - `1000Base-LX`: 1000 Mbit/s single-mode optical port (10 km).
  - `10GBase-T`: 10000 Mbit/s copper Ethernet port.
  - `10GBase-LR`: 10000 Mbit/s single-mode optical port (10 km).
  - `40GBase-LR`: 40000 Mbit/s single-mode optical port.
  - `100GBase-LR`: 100000 Mbit/s single-mode optical port.
-> **NOTE:** From version 1.185.0, `port_type` can be set to `40GBase-LR`, `100GBase-LR`. From version 1.230.1, `port_type` cannot be modified.
* `bandwidth` - (Optional) The maximum bandwidth of the hosted connection.
* `circuit_code` - (Optional) The circuit code of the Express Connect circuit.
* `peer_location` - (Optional) The geographical location of the data center.
* `redundant_physical_connection_id` - (Optional, ForceNew) The ID of the redundant Express Connect circuit. **NOTE:** From version 1.230.1, `redundant_physical_connection_id` cannot be modified.
* `physical_connection_name` - (Optional) The name of the Express Connect circuit.
* `description` - (Optional) The description of the Express Connect circuit.
* `status` - (Optional) The status of the Express Connect circuit. Valid values: `Confirmed`, `Enabled`, `Canceled`, `Terminated`. **NOTE:** From version 1.230.1, `status` can be set to `Confirmed`. If you want to set `status` to `Enabled`, `period` must be set.
* `period` - (Optional, Int, Available since v1.230.1) The subscription duration. Valid values:
  - If `pricing_cycle` is set to `Month`. Valid values: `1` to `9`.
  - If `pricing_cycle` is set to `Year`. Valid values: `1` to `5`.
* `pricing_cycle` - (Optional, Available since v1.230.1) The billing cycle of the subscription. Default value: `Month`. Valid values: `Month`, `Year`.
-> **NOTE:** `period` and `pricing_cycle` are valid only when `status` is set to `Enabled`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Physical Connection.
* `order_id` - The ID of the order that is placed. **Note:** `order_id` takes effect only if `status` is set to `Enabled`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Physical Connection.

## Import

Express Connect Physical Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_physical_connection.example <id>
```
