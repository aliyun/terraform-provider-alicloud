---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_gateway_route_table_attachment"
sidebar_current: "docs-alicloud-resource-vpc-gateway-route-table-attachment"
description: |-
  Provides a Alicloud VPC Gateway Route Table Attachment resource.
---

# alicloud_vpc_gateway_route_table_attachment

Provides a VPC Gateway Route Table Attachment resource. 

For information about VPC Gateway Route Table Attachment and how to use it, see [What is Gateway Route Table Attachment](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_gateway_route_table_attachment&exampleId=4137772f-0504-841d-3397-66cb3a3d5196eada3fd1&activeTab=example&spm=docs.r.vpc_gateway_route_table_attachment.0.4137772f05&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_vpc" "example" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "terraform-example"
}

resource "alicloud_route_table" "example" {
  vpc_id           = alicloud_vpc.example.id
  route_table_name = "terraform-example"
  description      = "terraform-example"
  associate_type   = "Gateway"
}

resource "alicloud_vpc_ipv4_gateway" "example" {
  ipv4_gateway_name = "terraform-example"
  vpc_id            = alicloud_vpc.example.id
  enabled           = "true"
}

resource "alicloud_vpc_gateway_route_table_attachment" "example" {
  ipv4_gateway_id = alicloud_vpc_ipv4_gateway.example.id
  route_table_id  = alicloud_route_table.example.id
}

```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_gateway_route_table_attachment&spm=docs.r.vpc_gateway_route_table_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) Specifies whether to only precheck this request. Default value: `false`.
* `ipv4_gateway_id` - (Required, ForceNew) The ID of the IPv4 Gateway instance.
* `route_table_id` - (Required, ForceNew) The ID of the Gateway route table to be bound.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<route_table_id>:<ipv4_gateway_id>`.
* `create_time` - The creation time of the resource.
* `status` - The status of the IPv4 Gateway instance. Value:
  - **Creating**: The function is being created.
  - **Created**: Created and available.
  - **Modifying**: is being modified.
  - **Deleting**: Deleting.
  - **Deleted**: Deleted.
  - **Activating**: enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gateway Route Table Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway Route Table Attachment.

## Import

VPC Gateway Route Table Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_gateway_route_table_attachment.example <route_table_id>:<ipv4_gateway_id>
```