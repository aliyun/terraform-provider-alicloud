---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_loadbalancer_common_bandwidth_package_attachment"
description: |-
  Provides a Alicloud NLB Loadbalancer Common Bandwidth Package Attachment resource.
---

# alicloud_nlb_loadbalancer_common_bandwidth_package_attachment

Provides a NLB Loadbalancer Common Bandwidth Package Attachment resource. Bandwidth Package Operation.

For information about NLB Loadbalancer Common Bandwidth Package Attachment and how to use it, see [What is Loadbalancer Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/nlb-instances-change).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment" "default" {
  bandwidth_package_id = "cbwp-2zexv44uov1m4b7xnh60j"
  load_balancer_id     = "nlb-f6gdwdsnt02uzx002l"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_package_id` - (Required, ForceNew) The ID of the bound shared bandwidth package.
* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<bandwidth_package_id>`.
* `status` - Network-based load balancing instance status. Value:, indicating that the instance listener will no longer forward traffic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Loadbalancer Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Loadbalancer Common Bandwidth Package Attachment.

## Import

NLB Loadbalancer Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_loadbalancer_common_bandwidth_package_attachment.example <load_balancer_id>:<bandwidth_package_id>
```