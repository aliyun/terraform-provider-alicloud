---
subcategory: "Time Series Database (TSDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_tsdb_instance"
sidebar_current: "docs-alicloud-resource-tsdb-instance"
description: |-
  Provides a Alicloud Time Series Database (TSDB) Instance resource.
---

# alicloud\_tsdb\_instance

Provides a Time Series Database (TSDB) Instance resource.

For information about Time Series Database (TSDB) Instance and how to use it, see [What is Time Series Database (TSDB)](https://www.alibabacloud.com/help/en/doc-detail/55652.htm).

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_tsdb_zones" "example" {}

resource "alicloud_vpc" "example" {
  cidr_block = "192.168.0.0/16"
  name       = "tf-testaccTsdbInstance"
}

resource "alicloud_vswitch" "example" {
  availability_zone = data.alicloud_tsdb_zones.example.ids.0
  cidr_block        = "192.168.1.0/24"
  vpc_id            = alicloud_vpc.example.id
}

resource "alicloud_tsdb_instance" "example" {
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.example.id
  instance_storage = "50"
  instance_class   = "tsdb.1x.basic"
  engine_type      = "tsdb_tsdb"
  instance_alias   = "tf-testaccTsdbInstance"
}

```

## Argument Reference

The following arguments are supported:

* `app_key` - (Optional) The app key.
* `disk_category` - (Optional, ForceNew) The disk type of instance. Valid when the engine type is `tsdb_influxdb`. `cloud_ssd` refers to SSD disk, `cloud_efficiency` refers to efficiency disk, `cloud_essd` refers to ESSD PL1 disk. Valid values: `cloud_efficiency`, `cloud_essd`, `cloud_ssd`.
* `duration` - (Optional, ForceNew) The duration.
* `engine_type` - (Optional, ForceNew) The engine type of instance Enumerative: `tsdb_tsdb` refers to TSDB, `tsdb_influxdb` refers to TSDB for InfluxDB️.
* `instance_alias` - (Optional) The alias of the instance.
* `instance_class` - (Required) The specification of the instance. 
    - Following enumerative value for TSDB for InfluxDB️ standart edition:
    - `influxdata.n1.mxlarge` refers to 2 CPU 8GB memory;
    - `influxdata.n1.xlarge` refers to 4 CPU 16GB memory;
    - `influxdata.n1.2xlarge` refers to 8 CPU 32 GB memory;
    - `influxdata.n1.4xlarge` refers to 16 CPU 64 GB memory;
    - `influxdata.n1.8xlarge` refers to 32 CPU 128 GB memory;
    - `influxdata.n1.16xlarge` refers to 64 CPU 256 GB memory. 
    - Following enumerative value for TSDB for InfluxDB High-availability edition:
    - `influxdata.n1.mxlarge_ha` refers to 2 CPU 8GB memory;
    - `influxdata.n1.xlarge_ha` refers to 4 CPU 16GB memory;
    - `influxdata.n1.2xlarge_ha` refers to 8 CPU 32 GB memory;
    - `influxdata.n1.4xlarge_ha` refers to 16 CPU 64 GB memory;
    - `influxdata.n1.8xlarge_ha` refers to 32 CPU 128 GB memory;
    - `influxdata.n1.16xlarge_ha` refers to 64 CPU 256 GB memory. 
    - Following enumerative value for TSDB:
    - `tsdb.1x.basic` refers to basic edition I;
    - `tsdb.3x.basic` refers to basic edition II; 
    - `tsdb.4x.basic` refers to basic edtion III;
    - `tsdb.12x.standard` refers to standard edition I; 
    - `tsdb.24x.standard` refers to standard edition II; 
    - `tsdb.48x.large` refers to ultimate edition I;
    - `tsdb.96x.large` refers to ultimate edition II.
* `instance_storage` - (Required) The storage capacity of the instance. Unit: GB. For example, the value 50 indicates 50 GB. Does not support shrink storage.
* `payment_type` - (Required, ForceNew) The billing method. Valid values: `PayAsYouGo` and `Subscription`. The `PayAsYouGo` value indicates the pay-as-you-go method, and the `Subscription` value indicates the subscription method.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - Instance status, enumerative: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 31 mins) Used when create the Instance.
* `update` - (Defaults to 31 mins) Used when update the Instance.

## Import

TSDB Instance can be imported using the id, e.g.

```
$ terraform import alicloud_tsdb_instance.example <id>
```
