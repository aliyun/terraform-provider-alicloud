---
subcategory: "Ack One"
layout: "alicloud"
page_title: "Alicloud: alicloud_ack_one_membership_attachment"
description: |-
  Provides a Alicloud Ack One Membership Attachment resource.
---

# alicloud_ack_one_membership_attachment

Provides an Ack One Membership Attachment resource. Fleet Manager Membership Attachment.

For information about Ack One Membership Attachment and how to use it, see [How to attach cluster tp hub](https://www.alibabacloud.com/help/en/ack/distributed-cloud-container-platform-for-kubernetes/developer-reference/api-adcp-2022-01-01-attachclustertohub).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ack_one_membership_attachment&exampleId=19b8d166-e942-7d41-5ffc-efc9ee5856172c52c8fb&activeTab=example&spm=docs.r.ack_one_membership_attachment.0.19b8d166e9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

variable "key_name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true

  is_enterprise_security_group = true
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.key_name
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 1
}

resource "alicloud_ack_one_cluster" "default" {
  depends_on = [alicloud_cs_managed_kubernetes.default]
  network {
    vpc_id    = alicloud_vpc.default.id
    vswitches = [alicloud_vswitch.default.id]
  }
  argocd_enabled = false
}

resource "alicloud_ack_one_membership_attachment" "default" {
  cluster_id     = alicloud_ack_one_cluster.default.id
  sub_cluster_id = alicloud_cs_managed_kubernetes.default.id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the cluster to which the membership is being attached.
* `sub_cluster_id` - (Required, ForceNew) The ID of the member being attached to the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Membership Attachment. It formats as < cluster_id >:< sub_cluster_id >.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 25 mins) Used when creating the Membership Attachment.
* `delete` - (Defaults to 25 mins) Used when deleting the Membership Attachment.

## Import

Ack One Membership Attachment can be imported using the id, which consists of cluster_id and sub_cluster_id, e.g.

```shell
terraform import alicloud_ack_one_membership_attachment.example <cluster_id>:<sub_cluster_id>
```
