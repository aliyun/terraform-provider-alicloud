---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_artifact_lifecycle_rule"
description: |-
  Provides a Alicloud CR Artifact Lifecycle Rule resource.
---

# alicloud_cr_artifact_lifecycle_rule

Provides a CR Artifact Lifecycle Rule resource.

Retention policies for versions in the warehouse.

For information about CR Artifact Lifecycle Rule and how to use it, see [What is Artifact Lifecycle Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateArtifactLifecycleRule).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_artifact_lifecycle_rule&exampleId=a7ba9972-6bd1-87fc-bb90-0a0382f9a4138687aae0&activeTab=example&spm=docs.r.cr_artifact_lifecycle_rule.0.a7ba99726b&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cr_ee_instance" "default" {
  default_oss_bucket = "true"
  instance_name      = var.name
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = "example repository for lifecycle rule"
}

resource "alicloud_cr_artifact_lifecycle_rule" "default" {
  auto                = true
  namespace_name      = alicloud_cr_ee_namespace.default.name
  retention_tag_count = "30"
  schedule_time       = "WEEK"
  scope               = "REPO"
  instance_id         = alicloud_cr_ee_instance.default.id
  tag_regexp          = ".*"
  repo_name           = alicloud_cr_ee_repo.default.name
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_artifact_lifecycle_rule&spm=docs.r.cr_artifact_lifecycle_rule.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `auto` - (Required) Whether to execute automatically
* `instance_id` - (Required, ForceNew) Instance ID
* `namespace_name` - (Optional) Namespace name
* `repo_name` - (Optional) Repository Name
* `retention_tag_count` - (Optional, Int) Number of Retention Tags
* `schedule_time` - (Optional) Execution cycle
* `scope` - (Optional) Scope of cleaning
* `tag_regexp` - (Optional) Retain regular expressions for mirrored versions

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<artifact_lifecycle_rule_id>`.
* `artifact_lifecycle_rule_id` - The first ID of the resource.
* `create_time` - Creation time.
* `modified_time` - Change time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Artifact Lifecycle Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Artifact Lifecycle Rule.
* `update` - (Defaults to 5 mins) Used when update the Artifact Lifecycle Rule.

## Import

CR Artifact Lifecycle Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_artifact_lifecycle_rule.example <instance_id>:<artifact_lifecycle_rule_id>
```