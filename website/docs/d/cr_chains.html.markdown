---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chains"
sidebar_current: "docs-alicloud-datasource-cr-chains"
description: |-
  Provides a list of Cr Chains to the user.
---

# alicloud\_cr\_chains

This data source provides the Cr Chains of the current Alibaba Cloud user.

For information about CR Chains and how to use it, see [What is Chain](https://www.alibabacloud.com/help/en/doc-detail/357821.html).

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cr_chains" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}

output "cr_chain_id_1" {
  value = data.alicloud_cr_chains.ids.chains.0.id
}

data "alicloud_cr_chains" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-Chain"
}

output "cr_chain_id_2" {
  value = data.alicloud_cr_chains.nameRegex.chains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of CR Enterprise Edition instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Chain name.
* `repo_name` - (Optional, ForceNew) The name of CR Enterprise Edition repository.
* `repo_namespace_name` - (Optional, ForceNew) The name of CR Enterprise Edition namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Chain IDs.
* `names` - A list of Chain names.
* `chains` - A list of Cr Chains. Each element contains the following attributes:
  * `id` - The resource ID of the delivery chain. The value formats as `<instance_id>:<chain_id>`.
  * `chain_id` - The ID of delivery chain.
  * `instance_id` - The ID of CR Enterprise Edition instance.
  * `chain_name` - The name of delivery chain.
  * `description` - The description of delivery chain.
  * `chain_config` - The configuration of delivery chain.
    * `routers` - Execution sequence relationship between delivery chain nodes.
      * `from` - Source node.
        * `node_name` - The name of node. Valid values: `DOCKER_IMAGE_BUILD`, `DOCKER_IMAGE_PUSH`, `VULNERABILITY_SCANNING`, `ACTIVATE_REPLICATION`, `TRIGGER`, `SNAPSHOT`, `TRIGGER_SNAPSHOT`.
      * `to` - Destination node.
        * `node_name` - The name of node. Valid values: `DOCKER_IMAGE_BUILD`, `DOCKER_IMAGE_PUSH`, `VULNERABILITY_SCANNING`, `ACTIVATE_REPLICATION`, `TRIGGER`, `SNAPSHOT`, `TRIGGER_SNAPSHOT`.
    * `nodes` - Each node in the delivery chain.
      * `node_name` - The name of delivery chain node.
      * `enable` - Whether to enable the delivery chain node. Valid values: `true`, `false`.
      * `node_config` - The configuration of delivery chain node.
        * `deny_policy` - Blocking rules for scanning nodes in delivery chain nodes. **Note:** When `node_name` is `VULNERABILITY_SCANNING`, the parameters in `deny_policy` need to be filled in.
          * `issue_count` - The count of scanning vulnerabilities that triggers blocking.
          * `issue_level` - The level of scanning vulnerability that triggers blocking. Valid values: `LOW`, `MEDIUM`, `HIGH`, `UNKNOWN`.
          * `logic` - The logic of trigger blocking. Valid values: `AND`, `OR`.
          * `action` - The action of trigger blocking. Valid values: `BLOCK`, `BLOCK_RETAG`, `BLOCK_DELETE_TAG`. While `Block` means block the delivery chain from continuing to execute, `BLOCK_RETAG` means block overwriting push image tag, `BLOCK_DELETE_TAG` means block deletion of mirror tags.
  * `create_time` - The creation time of delivery chain.
  * `modified_time` - The modification time of delivery chain description.
  * `scope_id` - Delivery chain scope ID.
  * `scope_type` - Delivery chain scope type.