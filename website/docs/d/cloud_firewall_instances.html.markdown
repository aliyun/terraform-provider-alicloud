---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instances"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-instances"
description: |-
  Provides a list of Cloud Firewall Instances to the user.
---

# alicloud\_cloud\_firewall\_instances

This data source provides the Cloud Firewall Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_firewall_instances" "ids" {}
output "cloud_firewall_instance_id_1" {
  value = data.alicloud_cloud_firewall_instances.ids.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Cloud Firewall Instances. Each element contains the following attributes:
    * `status` - The Status of Instance.
    * `renewal_status` - Automatic renewal status. Valid values: `AutoRenewal`,`ManualRenewal`. Default Value: `ManualRenewal`.
    * `id` - The ID of the Instance.
    * `instance_id` - The first ID of the resource.
    * `end_time` - The end time of the resource..
    * `create_time` - The Creation time of the resource.
    * `renewal_duration_unit` - Automatic renewal period unit. Valid values: `Month`,`Year`.
    * `payment_type` - The payment type of the resource. Valid values: `Subscription`.
    
