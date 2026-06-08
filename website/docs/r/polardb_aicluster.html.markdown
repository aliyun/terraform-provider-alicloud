---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_aicluster"
sidebar_current: "docs-alicloud-resource-polardb-aicluster"
description: |-
  Provides a PolarDB AI Cluster resource.
---

# alicloud_polardb_aicluster

Provides a PolarDB AI Cluster resource. A PolarDB AI Cluster is a managed inference cluster for deploying AI models on PolarDB infrastructure.

-> **NOTE:** Available since v1.279.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_aicluster&exampleId=fe7cf825-a403-e4b1-0392-cb8d9f8831fbd973c279&activeTab=example&spm=docs.r.polardb_aicluster.0.fe7cf825a4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_aicluster" "default" {
  region_id              = "cn-beijing"
  zone_id                = "cn-beijing-k"
  db_node_class          = "polar.mysql.g8.4xlarge.gu50"
  db_cluster_description = "tf-aicluster-example"
  pay_type               = "Postpaid"
  vpc_id                 = "vpc-xxx"
  vswitch_id             = "vsw-xxx"
  kube_type              = "ainode"
  model_name             = "Qwen3.5-9B"
  extension              = "maas"
  inference_engine       = "sglang"
  db_cluster_id          = "pc-xxx"
  security_group_id      = "sg-xxx"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_aicluster&spm=docs.r.polardb_aicluster.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `region_id` - (Required, ForceNew) The region ID of the AI cluster.
* `zone_id` - (Optional, ForceNew) The zone ID of the AI cluster.
* `db_node_class` - (Required, ForceNew) The DB node class of the AI cluster.
* `db_cluster_description` - (Optional, ForceNew) The description of the AI DB cluster.
* `pay_type` - (Required, ForceNew) The billing method. Valid values: `Postpaid`, `Prepaid`.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal.
* `period` - (Optional, ForceNew) The subscription period.
* `used_time` - (Optional, ForceNew) The subscription duration.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID.
* `security_group_id` - (Optional, ForceNew) The security group ID.
* `kube_type` - (Optional, ForceNew) The type of the Kubernetes cluster. Valid values: `ainode`.
* `model_name` - (Optional, ForceNew) The model name. Example: `Qwen3.5-9B`.
* `extension` - (Optional, ForceNew) The extension type. Valid values: `maas`, `custom`.
* `inference_engine` - (Optional, ForceNew) The inference engine. Valid values: `sglang`, `vllm`.
* `db_cluster_id` - (Optional, ForceNew) The ID of the associated DB cluster.
* `auto_use_coupon` - (Optional, ForceNew) Whether to use coupons automatically. Default value: `true`.
* `promotion_code` - (Optional, ForceNew) The promotion code.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID (same as the AI cluster ID).
* `status` - The status of the AI cluster.
* `model_type` - The model type. Example: `public`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the PolarDB AI Cluster.
* `delete` - (Defaults to 10 mins) Used when deleting the PolarDB AI Cluster.

## Import

PolarDB AI Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_aicluster.example pm-abc12345678
```
