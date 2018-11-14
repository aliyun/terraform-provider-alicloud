---
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interface_attachment
sidebar_current: "docs-alicloud-resource-network-interface-attachment
description: |-
  Provides an Alicloud ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.
---

# alicloud\_network\_interface\_attachment

Provides an Alicloud ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.

For information about Elastic Network Interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html).

## Example Usage

Bacis Usage

```
...
resource "alicloud_eni_attachment" "at" {
    instance_id = "${alicloud_instance.instance.id}"
    network_interface_id = "${alicloud_eni.eni.id}"
}
...
```

## Argument Reference

The following argument are supported:

* `instance_id` - (Required, ForceNew) The instance ID to attach.
* `network_interface_id` - (Required, ForceNew) The ENI ID to attach.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource, formatted as `<network_interface_id>:<instance_id>`.

## Import

Network Interfaces Attachment resource can be imported using the id, e.g.

```
$ terraform import alicloud_network_interface.eni eni-abc123456789000:i-abc123456789000
```
