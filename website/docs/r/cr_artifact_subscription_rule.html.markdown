---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_artifact_subscription_rule"
description: |-
  Provides a Alicloud cr Artifact Subscription Rule resource.
---

# alicloud_cr_artifact_subscription_rule

Provides a cr Artifact Subscription Rule resource.

Artifact subscription rule.

For information about cr Artifact Subscription Rule and how to use it, see [What is Artifact Subscription Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateArtifactSubscriptionRule).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_artifact_subscription_rule&exampleId=942f8b26-e53e-f384-d8ef-71f43119a90c311af05f&activeTab=example&spm=docs.r.cr_artifact_subscription_rule.0.942f8b26e5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cr_ee_instance" "example" {
  default_oss_bucket = "true"
  instance_name      = "tf-demo-cr-instance"
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

resource "alicloud_cr_ee_namespace" "example" {
  instance_id        = alicloud_cr_ee_instance.example.id
  name               = "tf-demo-cr-namespace"
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "example" {
  instance_id = alicloud_cr_ee_instance.example.id
  namespace   = alicloud_cr_ee_namespace.example.name
  name        = "tf-demo-cr-repo"
  repo_type   = "PRIVATE"
  summary     = "demo repository"
}

resource "alicloud_cr_artifact_subscription_rule" "example" {
  instance_id           = alicloud_cr_ee_instance.example.id
  source_provider       = "DOCKER_HUB"
  source_namespace_name = "library"
  source_repo_name      = "alpine"
  namespace_name        = alicloud_cr_ee_namespace.example.name
  repo_name             = alicloud_cr_ee_repo.example.name
  tag_regexp            = ".*"
  tag_count             = 10
  override              = true
  accelerate            = false
  platform              = ["linux/amd64", "linux/arm64"]
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_artifact_subscription_rule&spm=docs.r.cr_artifact_subscription_rule.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `accelerate` - (Optional) Whether to enable acceleration
* `instance_id` - (Required, ForceNew) Instance ID
* `namespace_name` - (Required) Namespace name
* `override` - (Optional) Whether to override existing tags
* `platform` - (Required, List) Subscription platform list
* `repo_name` - (Required) Repository name
* `source_namespace_name` - (Optional) Source namespace name
* `source_provider` - (Required) Source image registry provider, e.g. DOCKER_HUB, GCR, QUAY
* `source_repo_name` - (Required) Source repository name
* `tag_count` - (Required, Int) Number of tags to subscribe
* `tag_regexp` - (Required) Regular expression for subscribing tags

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<artifact_subscription_rule_id>`.
* `artifact_subscription_rule_id` - The first ID of the resource.
* `create_time` - Creation time.
* `modified_time` - Modification time.
* `source_domain` - Source domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Artifact Subscription Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Artifact Subscription Rule.
* `update` - (Defaults to 5 mins) Used when update the Artifact Subscription Rule.

## Import

cr Artifact Subscription Rule can be imported using the id, which consists of instance_id and artifact_subscription_rule_id, e.g.

```shell
$ terraform import alicloud_cr_artifact_subscription_rule.example <instance_id>:<artifact_subscription_rule_id>
```