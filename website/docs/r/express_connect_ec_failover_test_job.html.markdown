---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_ec_failover_test_job"
description: |-
  Provides a Alicloud Express Connect Ec Failover Test Job resource.
---

# alicloud_express_connect_ec_failover_test_job

Provides a Express Connect Ec Failover Test Job resource. Express Connect Failover Test Job.

For information about Express Connect Ec Failover Test Job and how to use it, see [What is Ec Failover Test Job](https://www.alibabacloud.com/help/zh/express-connect/developer-reference/api-vpc-2016-04-28-createfailovertestjob-efficiency-channels).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_ec_failover_test_job&exampleId=b53a3a45-f617-164b-24e8-6c17d9997a2b45116598&activeTab=example&spm=docs.r.express_connect_ec_failover_test_job.0.b53a3a45f6&intl_lang=EN_US" target="_blank">
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

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_ec_failover_test_job" "default" {
  description = var.name
  job_type    = "StartNow"
  resource_id = [
    "${data.alicloud_express_connect_physical_connections.default.ids.0}",
    "${data.alicloud_express_connect_physical_connections.default.ids.1}"
  ]
  job_duration              = "1"
  resource_type             = "PHYSICALCONNECTION"
  ec_failover_test_job_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_ec_failover_test_job&spm=docs.r.express_connect_ec_failover_test_job.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Job description.
* `ec_failover_test_job_name` - (Optional) Job name.
* `job_duration` - (Required) Job duration.
* `job_type` - (Required, ForceNew) Job type.
* `resource_id` - (Required) Resource id list.
* `resource_type` - (Required, ForceNew) Resource type.
* `status` - (Optional, Computed) The status of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ec Failover Test Job.
* `delete` - (Defaults to 5 mins) Used when delete the Ec Failover Test Job.
* `update` - (Defaults to 5 mins) Used when update the Ec Failover Test Job.

## Import

Express Connect Ec Failover Test Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_ec_failover_test_job.example <id>
```