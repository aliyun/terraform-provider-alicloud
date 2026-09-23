---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_physical_connection"
sidebar_current: "docs-alicloud-resource-express-connect-virtual-physical-connection"
description: |-
  Provides a Alicloud Express Connect Virtual Physical Connection resource.
---

# alicloud_express_connect_virtual_physical_connection

Provides a Express Connect Virtual Physical Connection resource.

For information about Express Connect Virtual Physical Connection and how to use it, see [What is Virtual Physical Connection](https://www.alibabacloud.com/help/en/express-connect/latest/createvirtualphysicalconnection#doc-api-Vpc-CreateVirtualPhysicalConnection).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_virtual_physical_connection&exampleId=6c1a81e4-689b-36c7-4d27-e62c70cfa3c9cc63df38&activeTab=example&spm=docs.r.express_connect_virtual_physical_connection.0.6c1a81e468&intl_lang=EN_US" target="_blank">
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
data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
data "alicloud_account" "default" {}
resource "alicloud_express_connect_virtual_physical_connection" "example" {
  virtual_physical_connection_name = var.name
  description                      = var.name
  order_mode                       = "PayByPhysicalConnectionOwner"
  parent_physical_connection_id    = data.alicloud_express_connect_physical_connections.example.ids.0
  spec                             = "50M"
  vlan_id                          = random_integer.vlan_id.id
  vpconn_ali_uid                   = data.alicloud_account.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_virtual_physical_connection&spm=docs.r.express_connect_virtual_physical_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the physical connection.
* `expect_spec` - (Optional) The estimated bandwidth value of the shared line. Valid values: `50M`, `100M`, `200M`, `300M`, `400M`, `500M`, `1G`, `2G`, `5G`, `8G`, and `10G`. **Note**: By default, the values of 2G, 5G, 8G, and 10G are unavailable. If you want to specify these values, contact your customer manager. Unit: **M** indicates Mbps, **G** indicates Gbps.
* `order_mode` - (Required, ForceNew) The payment method of shared dedicated line. Value:
  - **PayByPhysicalConnectionOwner**: indicates that the owner of the physical line associated with the shared line pays.
  - **PayByVirtualPhysicalConnectionOwner**: indicates that the owner of the shared line pays.
* `parent_physical_connection_id` - (Required, ForceNew) The ID of the instance of the physical connection.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `spec` - (Required, ForceNew) The bandwidth value of the shared line. Valid values: `50M`, `100M`, `200M`, `300M`, `400M`, `500M`, `1G`, `2G`, `5G`, `8G`, and `10G`. **Note**: By default, the values of 2G, 5G, 8G, and 10G are unavailable. If you want to specify these values, contact your customer manager. Unit: **M** indicates Mbps, **G** indicates Gbps.
* `virtual_physical_connection_name` - (Optional, ForceNew) The name of the physical connection.
* `vlan_id` - (Required) The VLAN ID of the shared leased line. Valid values: `0` to `2999`.
* `vpconn_ali_uid` - (Required, ForceNew) The ID of the Alibaba Cloud account (primary account) of the owner of the shared line.
* `dry_run` - (Optional) Specifies whether to precheck the API request. Valid values: `true` and `false`.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `access_point_id` - The ID of the access point of the physical connection.
* `ad_location` - The physical location where the physical connection access device is located.
* `bandwidth` - The bandwidth of the physical connection. Unit: Mbps.
* `business_status` - The commercial status of the physical line. Value:-**Normal**: activated.-**Financialized**: Arrears locked.-**SecurityLocked**: locked for security reasons.
* `circuit_code` - The circuit code provided by the operator for the physical connection.
* `create_time` - The creation time of the resource
* `enabled_time` - The opening time of the physical connection.
* `end_time` - The expiration time of the shared line.Time is expressed according to ISO8601 standard and UTC time is used. The format is: YYYY-MM-DDThh:mm:ssZ.
* `line_operator` - Operators that provide access to physical lines. Value:-**CT**: China Telecom.-**CU**: China Unicom.-**CM**: China Mobile.-**CO**: China Other.-**Equinix**:Equinix.-**Other**: Other abroad.
* `loa_status` - The state of LOA. Value:-**Applying**:LOA application.-**Accept**:LOA application passed.-**Available**:LOA is Available.-**Rejected**:LOA application Rejected.-**Completing**: The dedicated line is under construction.-**Complete**: The construction of the dedicated line is completed.-**Deleted**:LOA has been Deleted.
* `parent_physical_connection_ali_uid` - The ID of the Alibaba Cloud account (primary account) to which the physical connection belongs.
* `peer_location` - The geographic location of the local data center.
* `port_number` - The port number of the physical connection device.
* `port_type` - Physical connection port type. Value:-**100Base-T**: 100 megabytes port.-**1000Base-T**: Gigabit port.-**1000Base-LX**: Gigabit single mode optical port (10km).-**10GBase-T**: 10 Gigabit port.-**10GBase-LR**: 10 Gigabit single mode optical port (10km).-**40GBase-LR**: 40 megabytes single-mode optical port.-**100GBase-LR**: 100,000 megabytes single-mode optical port.
* `redundant_physical_connection_id` - The ID of the redundant physical connection.
* `status` - The status of the resource
* `virtual_physical_connection_status` - The business status of the shared line. Value:-**Confirmed**: The shared line has been Confirmed to receive.-**UnConfirmed**: The shared line has not been confirmed to be received.-**Deleted**: The shared line has been Deleted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Virtual Physical Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Virtual Physical Connection.
* `update` - (Defaults to 5 mins) Used when update the Virtual Physical Connection.

## Import

Express Connect Virtual Physical Connection can be imported using the id, e.g.

```shell
$terraform import alicloud_express_connect_virtual_physical_connection.example <id>
```