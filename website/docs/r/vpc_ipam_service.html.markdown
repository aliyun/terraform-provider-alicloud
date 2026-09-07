---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_service"
description: |-
  Provides a Alicloud Vpc Ipam Service resource.
---

# alicloud_vpc_ipam_service

Provides a Vpc Ipam Service resource.

Ipam service, used to support automatic provisioning of Terraform.

For information about Vpc Ipam Service and how to use it, see [What is Service](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/OpenVpcIpamService).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipam_service&exampleId=282715d4-1ddb-181e-3553-4c53ba2518b892707963&activeTab=example&spm=docs.r.vpc_ipam_service.0.282715d41d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_vpc_ipam_service" "default" {
}
```

### Deleting `alicloud_vpc_ipam_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_vpc_ipam_service`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_ipam_service&spm=docs.r.vpc_ipam_service.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.
* `enabled` - Whether the IPAM service has been activated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.

## Import

Vpc Ipam Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipam_service.example 
```