---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_cluster"
sidebar_current: "docs-alicloud-resource-edas-cluster"
description: |-
  Provides an EDAS cluster resource.
---

# alicloud_edas_cluster

Provides an EDAS cluster resource, see [What is EDAS Cluster](https://www.alibabacloud.com/help/en/edas/developer-reference/api-edas-2017-08-01-insertcluster).

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_cluster&exampleId=6454467d-e77d-4dfc-0858-498320a8da115e6a50a0&activeTab=example&spm=docs.r.edas_cluster.0.6454467de7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_edas_cluster" "default" {
  cluster_name      = var.name
  cluster_type      = "2"
  network_mode      = "2"
  logical_region_id = data.alicloud_regions.default.regions.0.id
  vpc_id            = alicloud_vpc.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_cluster&spm=docs.r.edas_cluster.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) The name of the cluster that you want to create.
* `cluster_type` - (Required, ForceNew) The type of the cluster that you want to create. Valid values only: 2: ECS cluster.
* `network_mode` - (Required, ForceNew) The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC.
* `logical_region_id` - (Optional, ForceNew) The ID of the namespace where you want to create the application. You can call the ListUserDefineRegion operation to query the namespace ID.
* `vpc_id` - (Optional, ForceNew) The ID of the Virtual Private Cloud (VPC) for the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<cluster_id>`.

## Import

EDAS cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_cluster.cluster cluster_id
```
