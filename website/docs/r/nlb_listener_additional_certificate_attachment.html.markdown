---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_listener_additional_certificate_attachment"
description: |-
  Provides a Alicloud NLB Listener Additional Certificate Attachment resource.
---

# alicloud_nlb_listener_additional_certificate_attachment

Provides a NLB Listener Additional Certificate Attachment resource. 

For information about NLB Listener Additional Certificate Attachment and how to use it, see [What is Listener Additional Certificate Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/nlb-instances-change).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_nlb_listener_additional_certificate_attachment" "default" {
  certificate_id = "10513353-cn-hangzhou"
  listener_id    = "lsn-gmpiutpjy6e58qequ1@14500"
}
```

## Argument Reference

The following arguments are supported:
* `certificate_id` - (Required, ForceNew) Certificate ID. Currently, only server certificates are supported.
* `dry_run` - (Optional) Whether to PreCheck only this request, value: - **true**: sends a check request and does not create a resource. Check items include whether required parameters, request format, and business restrictions have been filled in. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '. - **false** (default): Sends a normal request, returns the HTTP 2xx status code after the check, and directly performs the operation.
* `listener_id` - (Required, ForceNew) The ID of the tcpssl listener.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<listener_id>:<certificate_id>`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Listener Additional Certificate Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Listener Additional Certificate Attachment.

## Import

NLB Listener Additional Certificate Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_listener_additional_certificate_attachment.example <listener_id>:<certificate_id>
```