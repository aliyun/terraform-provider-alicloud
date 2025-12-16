---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_hpc_cluster"
sidebar_current: "docs-alicloud-resource-ecs-hpc-cluster"
description: |-
  Provides a Alicloud ECS Hpc Cluster resource.
---

# alicloud\_ecs\_hpc\_cluster

Provides a ECS Hpc Cluster resource.

For information about ECS Hpc Cluster and how to use it, see [What is Hpc Cluster](https://www.alibabacloud.com/help/en/doc-detail/109138.htm).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_hpc_cluster&exampleId=f89770ff-c9df-6eb1-a7fc-4962e10ffb0da51368ca&activeTab=example&spm=docs.r.ecs_hpc_cluster.0.f89770ffc9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecs_hpc_cluster" "example" {
  name        = "tf-testAcc"
  description = "For Terraform Test"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_hpc_cluster&spm=docs.r.ecs_hpc_cluster.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of ECS Hpc Cluster.
* `name` - (Required) The name of ECS Hpc Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hpc Cluster.

## Import

ECS Hpc Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_hpc_cluster.example <id>
```
