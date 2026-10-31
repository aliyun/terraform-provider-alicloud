---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_namespace"
description: |-
  Provides a Alicloud Realtime Compute Namespace resource.
---

# alicloud_realtime_compute_namespace

Provides a Realtime Compute Namespace resource.

Namespace of Realtime Compute for Apache Flink, used to reference reusable values in deployments.

For information about Realtime Compute Namespace and how to use it, see [What is Namespace](https://api.aliyun.com/api/foasconsole/2021-10-28/CreateNamespace).

-> **NOTE:** Available since 2.3.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_realtime_compute_namespace&exampleId=780af765-a79c-16b4-8c0b-d416be66fd32082ae05c&activeTab=example&spm=docs.r.realtime_compute_variable.0.780af765a7&intl_lang=EN_US" target="_blank">
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

data "alicloud_oss_buckets" "default" {
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = "code-test-tf"
  storage {
    oss {
      bucket = data.alicloud_oss_buckets.default.buckets.0.name
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "1000"
    memory_gb = "4000"
  }
  default_namespace_resource_spec {
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = var.name
  resource_spec {
    cpu       = "2"
    memory_gb = "8"
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_realtime_compute_namespace&spm=docs.r.realtime_compute_variable.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The instance ID.
* `namespace` - (Required, ForceNew) The namespace name.
* `ha` - (Optional, Bool) Specifies whether to enable zone-disaster recovery high availability.
* `resource_spec` - (Required, Set) The resource specifications. See [`resource_spec`](#resource_spec) below.

### `resource_spec`

The resource_spec supports the following:

* `cpu` - (Required, Int) The number of CPUs.
* `memory_gb` - (Required, Int) The memory size.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. It formats as `<instance_id>:<namespace>`.
* `status` - The namespace status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Namespace.
* `update` - (Defaults to 5 mins) Used when update the Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Namespace.

## Import

Realtime Compute Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_namespace.example <instance_id>:<namespace>
```
