---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance"
sidebar_current: "docs-alicloud-resource-ots-instance"
description: |-
  Provides an OTS (Open Table Service) instance resource.
---

# alicloud_ots_instance

This resource will help you to manager a [Table Store](https://www.alibabacloud.com/help/doc-detail/27280.htm) Instance.
It is foundation of creating data table.

-> **NOTE:** Available since v1.10.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ots_instance&exampleId=8a075a3f-f162-bf06-1a9f-1df4ce7c5e5cad33bafd&activeTab=example&spm=docs.r.ots_instance.0.8a075a3ff1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Vpc"
  tags = {
    Created = "TF"
    For     = "Building table"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ots_instance&spm=docs.r.ots_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the instance.
* `network_type_acl` - (Optional, Available since v1.221.0) The set of network types that are allowed access. Valid optional values:
  * `CLASSIC` - Classic network.
  * `VPC` - VPC network.
  * `INTERNET` - Public internet.

  Default to ["VPC", "CLASSIC", "INTERNET"].
* `network_source_acl` - (Optional, Available since v1.221.0) The set of request sources that are allowed access. Valid optional values:
  * `TRUST_PROXY` - Trusted proxy, usually the Alibaba Cloud console.

  Default to ["TRUST_PROXY"].
* `accessed_by` - (Optional, Deprecated since v1.221.0) The network limitation of accessing instance. Valid values:
  * `Any` - Allow all network to access the instance.
  * `Vpc` - Only can the attached VPC allow to access the instance.
  * `ConsoleOrVpc` - Allow web console or the attached VPC to access the instance.

  Default to "Any".
* `resource_group_id` - (Optional, Available since v1.221.0) The resource group the instance belongs to.
  Default to Alibaba Cloud default resource group.
* `instance_type` - (Optional, ForceNew) The type of instance. Valid values are "Capacity" and "HighPerformance". Default to "HighPerformance".
* `description` - (Optional, ForceNew) The description of the instance. Currently, it does not support modifying.
* `tags` - (Optional) A mapping of tags to assign to the instance.
* `vcu` - (Optional, ForceNew) The number of VCUs for a reserved-mode (VCU) instance. When this value is greater than 0, the instance is created as a VCU instance via the CreateVCUInstance API; otherwise a free instance is created.
* `period_in_month` - (Optional, ForceNew) The subscription duration in months. Required for VCU instances.
* `enable_auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal for the VCU instance.
* `enable_elastic_vcu` - (Optional, ForceNew) Whether to enable instance elasticity. When enabled, the peak VCU usage can exceed the purchased VCU amount.
* `auto_renew_period_in_month` - (Optional, ForceNew) The auto-renewal cycle in months. Required when `enable_auto_renew` is true.
* `alias_name` - (Optional) The alias of the instance.
* `elastic_vcu_upper_limit` - (Optional) The upper limit for the VCUs of the instance. Valid values: number of reserved VCUs + 0.1 to 2000.
* `policy` - (Optional) The instance access control policy in JSON format.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as the "name".
* `payment_type` - The payment type of the instance.
* `policy_version` - The version of the instance policy.
* `vcu` - The number of VCUs of the instance.
* `alias_name` - The alias of the instance.
* `elastic_vcu_upper_limit` - The upper limit for the VCUs of the instance.
* `policy` - The instance access control policy in JSON format.

## Import

OTS instance can be imported using instance id or name, e.g.

```shell
$ terraform import alicloud_ots_instance.foo "my-ots-instance"
```

