---
subcategory: "Ack One"
layout: "alicloud"
page_title: "Alicloud: alicloud_ack_one_cluster"
description: |-
  Provides a Alicloud Ack One Cluster resource.
---

# alicloud_ack_one_cluster

Provides a Ack One Cluster resource. Fleet Manager Cluster.

For information about Ack One Cluster and how to use it, see [What is Cluster](https://www.alibabacloud.com/help/en/ack/distributed-cloud-container-platform-for-kubernetes/developer-reference/api-adcp-2022-01-01-createhubcluster).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ack_one_cluster&exampleId=b4cd95b5-eaa3-1433-1ebc-eaa890f904691c4ef5fa&activeTab=example&spm=docs.r.ack_one_cluster.0.b4cd95b5ea&intl_lang=EN_US" target="_blank">
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "defaultyVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "172.16.2.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name

}


resource "alicloud_ack_one_cluster" "default" {
  network {
    vpc_id    = alicloud_vpc.defaultVpc.id
    vswitches = ["${alicloud_vswitch.defaultyVSwitch.id}"]
  }
  profile = "XFlow"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ack_one_cluster&spm=docs.r.ack_one_cluster.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cluster_name` - (Optional, ForceNew, Computed) Cluster name.
* `network` - (Required, ForceNew) Cluster network information. See [`network`](#network) below.
* `profile` - (Optional, ForceNew, Computed) Cluster attributes. Valid values: 'Default', 'XFlow'.

**Note**: When profile is Default, vswitches might not be deleted when cluster is deleted because there are some remaining resources in the vswitches. We are still fixing this problem.

* `argocd_enabled` - (Optional) (Available since v1.243.0) Whether to enable ArgoCD. Default to true. Only valid when `profile` is 'Default'. It has to be false when cluster is deleted.


### `network`

The network supports the following:
* `vpc_id` - (Required, ForceNew) VpcId to which the cluster belongs.
* `vswitches` - (Required, ForceNew) Switch to which the cluster belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Cluster creation time.
* `network` - Cluster network information.
  * `security_group_ids` - Security group to which the cluster belongs.
* `status` - The status of the resource.
* `argocd_enabled` - Whether to enable ArgoCD.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 25 mins) Used when create the Cluster.
* `delete` - (Defaults to 25 mins) Used when delete the Cluster.

## Import

Ack One Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_ack_one_cluster.example <id>
```