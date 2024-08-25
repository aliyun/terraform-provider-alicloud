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

For information about Express Connect Physical Connection and how to use it, see [What is Physical Connection](https://www.alibabacloud.com/help/doc-detail/44852.htm).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_express_connect_physical_connection&exampleId=0a4b0339-0535-aea4-006f-4805696f33ebcf70ddd7&activeTab=example&spm=docs.r.express_connect_physical_connection.0.0a4b033905&intl_lang=EN_US" target="_blank">
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

* `access_point_id` - (Required, ForceNew) The Physical Leased Line Access Point ID.
* `bandwidth` - (Optional) On the Bandwidth of the ECC Service and Physical Connection.
* `circuit_code` - (Optional) Operators for Physical Connection Circuit Provided Coding.
* `description` - (Optional) The Physical Connection to Which the Description.
* `line_operator` - (Required) Provides Access to the Physical Line Operator. Valid values:
  * CT: China Telecom
  * CU: China Unicom
  * CM: china Mobile
  * CO: Other Chinese
  * Equinix: Equinix
  * Other: Other Overseas.

* `peer_location` - (Optional) and an on-Premises Data Center Location.
* `physical_connection_name` - (Optional) on Behalf of the Resource Name of the Resources-Attribute Field.
* `port_type` - (Optional) The Physical Leased Line Access Port Type. Valid value:
  * 100Base-T: Fast Electrical Ports
  * 1000Base-T: gigabit Electrical Ports
  * 1000Base-LX: Gigabit Singlemode Optical Ports (10Km)
  * 10GBase-T: Gigabit Electrical Port
  * 10GBase-LR: Gigabit Singlemode Optical Ports (10Km).
  * 40GBase-LR: 40 Gigabit Singlemode Optical Ports.
  * 100GBase-LR: One hundred thousand Gigabit Singlemode Optical Ports.

**NOTE:** From in v1.185.0+, The `40GBase-LR` and `100GBase-LR` is valid. and Set these values based on the water levels of background ports. For details about the water levels, contact the business manager.

* `redundant_physical_connection_id` - (Optional) Redundant Physical Connection to Which the ID.
* `status` - (Optional) Resources on Behalf of a State of the Resource Attribute Field. Valid values: `Canceled`, `Enabled`, `Terminated`.
* `type` - (Optional, ForceNew) Physical Private Line of Type. Default Value: VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Physical Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Physical Connection.

## Import

Express Connect Physical Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_physical_connection.example <id>
```
