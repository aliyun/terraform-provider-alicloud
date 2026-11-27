---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_tr_association"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router Tr Association resource.
---

# alicloud_express_connect_router_tr_association

Provides a Express Connect Router Express Connect Router Tr Association resource. Leased line gateway and TR binding relationship object.

For information about Express Connect Router Express Connect Router Tr Association and how to use it, see [What is Express Connect Router Tr Association](https://next.api.alibabacloud.com/api/ExpressConnectRouter/2023-09-01/CreateExpressConnectRouterAssociation).

-> **NOTE:** Available since v1.224.0.

-> **NOTE:** When an association is created with `create_attachment = true`, destroying it also deletes the TR-side attachment created by that request. When `create_attachment` is `false` or unavailable after import, the provider does not send `DeleteAttachment` and does not inspect or delete independently managed attachments. If deletion returns `OperationDenied.EcrAttachmentExisted`, the attachment manager must delete the TR-side attachment before the association destroy is retried; the provider does not retry that error with `DeleteAttachment = true`. No pre-delete `cen:ListTransitRouterEcrAttachments` permission is required.

-> **NOTE:** The TR-side ECR attachment can also be managed by [`alicloud_cen_transit_router_ecr_attachment`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_ecr_attachment). `CreateTransitRouterEcrAttachment` requires the association to exist first. When the two resources are used together, preserve the association-to-attachment reference shown in the CEN attachment example so Terraform creates the association first and destroys the independently managed attachment before the association. The association resource does not automatically delete an independently managed or unknown-ownership attachment.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_tr_association&exampleId=dd472839-6d9b-59e6-6edb-30df6adee7ca28a7b887&activeTab=example&spm=docs.r.express_connect_router_tr_association.0.dd4728396d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "alowprefix1" {
  default = "10.0.0.0/24"
}

variable "allowprefix2" {
  default = "10.0.1.0/24"
}

variable "allowprefix3" {
  default = "10.0.2.0/24"
}

variable "allowprefix4" {
  default = "10.0.3.0/24"
}

variable "asn" {
  default = "4200001003"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultpX0KlC" {
  alibaba_side_asn = var.asn
}

resource "alicloud_cen_instance" "default418DC9" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultRYcjsc" {
  cen_id = alicloud_cen_instance.default418DC9.id
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_tr_association" "default" {
  ecr_id                  = alicloud_express_connect_router_express_connect_router.defaultpX0KlC.id
  cen_id                  = alicloud_cen_instance.default418DC9.id
  transit_router_owner_id = data.alicloud_account.current.id
  allowed_prefixes = [
    "${var.alowprefix1}",
    "${var.allowprefix3}",
    "${var.allowprefix2}"
  ]
  transit_router_id     = alicloud_cen_transit_router.defaultRYcjsc.transit_router_id
  association_region_id = "cn-hangzhou"
  create_attachment     = true
  allowed_prefixes_mode = "MatchMode"
  description           = "tf-testacc-description"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_router_tr_association&spm=docs.r.express_connect_router_tr_association.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `allowed_prefixes` - (Optional) List of allowed route prefixes.
* `allowed_prefixes_mode` - (Optional, String, Computed, Available since v2.4.0) The prefix routing mode. Valid values: `MatchMode` and `IncrementalMode`. Configuring this argument requires a non-empty `allowed_prefixes` list; otherwise Terraform rejects the configuration during plan. To clear a non-empty prefix list, remove `allowed_prefixes_mode` together with setting `allowed_prefixes = []` or removing both arguments. `MatchMode` advertises only the configured `allowed_prefixes` to the data center and withdraws the detailed routes that were advertised before; `IncrementalMode` advertises the configured `allowed_prefixes` and withdraws only the detailed routes inside them, while detailed routes outside the configured prefixes stay advertised. When omitted, the value returned by the API is kept in state.
* `association_region_id` - (Required, ForceNew) The region to which the VPC or TR belongs.
* `cen_id` - (Optional, ForceNew) The ID of the CEN instance.
* `create_attachment` - (Optional, Bool, ForceNew, Computed, Available since v2.4.0) Whether the ECR side initiates the association when the ECR binds the TR. `true` means the association does not need to be initiated on the transit router side; `false` means it must be initiated from the transit router side. Changing this argument replaces the association. The value is only sent when the association is created and is not returned by the API, so it is not restored after import.
* `description` - (Optional, String, Available since v2.4.0) The description of the association resource. The length is 0-128 characters. Removing the argument clears the description.
* `ecr_id` - (Required, ForceNew) The ID of the leased line gateway instance.
* `transit_router_id` - (Optional, ForceNew, Computed) The ID of the forwarding router instance.
* `transit_router_owner_id` - (Optional, ForceNew) The ID of the Alibaba Cloud account to which the forwarding router belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<association_id>:<transit_router_id>`.
* `association_id` - The first ID of the resource.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Express Connect Router Tr Association.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router Tr Association.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router Tr Association.

## Import

Express Connect Router Express Connect Router Tr Association can be imported using the id, which consists of ecr_id, association_id and transit_router_id, e.g.

```shell
$ terraform import alicloud_express_connect_router_tr_association.example <ecr_id>:<association_id>:<transit_router_id>
```

`create_attachment` is create-only and is not returned by the API, so it is not restored on import. An imported association has unknown attachment ownership: the provider omits `DeleteAttachment`, does not inspect the TR side, and does not escalate `OperationDenied.EcrAttachmentExisted` into a cascading delete. The attachment manager must remove an independently managed or unknown-ownership attachment before the association can be destroyed. Only an association configured with `create_attachment = true` may use `DeleteAttachment = true`. Other transient errors (`IncorrectStatus.EcrAssociation`, `Conflict.Lock`, `IncorrectStatus.Ecr`, `OperationDenied.ForceDeleting`) are retried without changing the delete scope.
