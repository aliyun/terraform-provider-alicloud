---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_sas_private_link_endpoint"
sidebar_current: "docs-alicloud-resource-threat-detection-sas-private-link-endpoint"
description: |-
  Provides a Alicloud Threat Detection Sas Private Link Endpoint resource.
---

# alicloud_threat_detection_sas_private_link_endpoint

Provides a Threat Detection Sas Private Link Endpoint resource.

For information about Threat Detection Sas Private Link Endpoint and how to use it, see [What is Sas Private Link Endpoint](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createsasprivatelinkendpoint).

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "default" {
  vpc_name   = "tf-sas-pl-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/21"
  zone_id    = "cn-hangzhou-h"
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = "tf-sas-pl-example"
}

resource "alicloud_threat_detection_sas_private_link_endpoint" "default" {
  node_name         = "tf-sas-pl-example"
  vpc_id            = alicloud_vpc.default.id
  security_group_id = alicloud_security_group.default.id
  zones {
    v_switch_id = alicloud_vswitch.default.id
    zone_id     = "cn-hangzhou-h"
  }
}
```

## Argument Reference

The following arguments are supported:
* `node_name` - (Required) The name of the private link endpoint node.
* `security_group_id` - (Optional, ForceNew) The security group ID associated with the endpoint.
* `vpc_id` - (Optional, ForceNew) The VPC ID.
* `region_id` - (Optional, ForceNew) The region ID. If not set, the provider region will be used.
* `zones` - (Required, ForceNew) The availability zone information. See [`zones`](#zones) below. **Note: The parameter is immutable after resource creation.**

### `zones`

The zones supports the following:
* `v_switch_id` - (Required) The VSwitch ID.
* `zone_id` - (Required) The Zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Sas Private Link Endpoint. It is the same as the node ID.
* `status` - The status of the endpoint.
* `update_domain` - The update domain.
* `jsrv_domain` - The jsrv domain.

## Timeouts

The Timeouts block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for actions:
* `create` - (Defaults to 5 mins) Used when create the Sas Private Link Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Sas Private Link Endpoint.
* `delete` - (Defaults to 10 mins) Used when delete the Sas Private Link Endpoint.

## Import

Threat Detection Sas Private Link Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_sas_private_link_endpoint.example <id>
```
