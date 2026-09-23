---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network"
description: |-
  Provides a Alicloud ENS Network resource.
---

# alicloud_ens_network

Provides a ENS Network resource. 

For information about ENS Network and how to use it, see [What is Network](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createnetwork-1).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_network&exampleId=9317eafe-9bd2-b5b1-c6db-79edc7036c857c72e1cd&activeTab=example&spm=docs.r.ens_network.0.9317eafe9b&intl_lang=EN_US" target="_blank">
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


resource "alicloud_ens_network" "default" {
  network_name = var.name

  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ens_network&spm=docs.r.ens_network.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Required, ForceNew) The network segment of the network. You can use the following network segments or a subset of them as the network segment: `10.0.0.0/8` (default), `172.16.0.0/12`, `192.168.0.0/16`.
* `description` - (Optional) Description information.Rules:It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`. Example value: this is my first network.
* `ens_region_id` - (Required, ForceNew) Ens node IDExample value: cn-beijing-telecom.
* `network_name` - (Optional) Name of the network instanceThe naming rules are as follows: 1. Length is 2~128 English or Chinese characters; 2. It must start with a large or small letter or Chinese, not with `http://` and `https://`; 3. Can contain numbers, colons (:), underscores (_), or dashes (-).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time, timestamp (MS).
* `status` - The status of the network instance. Pending: Configuring, Available: Available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network.
* `delete` - (Defaults to 5 mins) Used when delete the Network.
* `update` - (Defaults to 5 mins) Used when update the Network.

## Import

ENS Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network.example <id>
```