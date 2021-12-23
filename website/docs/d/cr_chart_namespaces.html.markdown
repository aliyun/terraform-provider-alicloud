---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_namespaces"
sidebar_current: "docs-alicloud-datasource-cr-chart-namespaces"
description: |-
  Provides a list of Cr Chart Namespaces to the user.
---

# alicloud\_cr\_chart\_namespaces

This data source provides the Cr Chart Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_ee_instance" "default" {
  payment_type  = "Subscription"
  period        = 1
  instance_type = "Advanced"
  instance_name = "name"
}

resource "alicloud_cr_chart_namespace" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  namespace_name = "name"
}
data "alicloud_cr_chart_namespaces" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
}
output "output" {
  default = data.alicloud_cr_chart_namespaces.default.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Chart Namespace IDs.
* `instance_id`- (Reu, ForceNew, Computed)  A list of Chart Namespace IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by name space name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of matched Container Registry Enterprise Edition namespaces.
* `names` - A list of namespace names.
* `namespaces` - A list of Cr Chart Namespaces. Each element contains the following attributes:
	* `auto_create_repo` - AutoCreateRepo.
	* `chart_namespace_id` - The first ID of the resource.
	* `default_repo_type` - DefaultRepoType.
	* `id` - The ID of the Chart Namespace.
	* `instance_id` - InstanceId.
	* `namespace_name` - NamespaceName.