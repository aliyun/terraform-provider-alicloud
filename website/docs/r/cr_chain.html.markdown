---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chain"
sidebar_current: "docs-alicloud-resource-cr-chain"
description: |-
  Provides a Alicloud CR Chain resource.
---

# alicloud\_cr\_chain

Provides a CR Chain resource.

For information about CR Chain and how to use it, see [What is Chain](https://www.alibabacloud.com/help/en/doc-detail/357808.html).

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "default-name"
}

data "alicloud_cr_ee_instances" "default" {
  name_regex = "default-instance"
}

resource "alicloud_cr_ee_instance" "default" {
  count          = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = "default-instance"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_namespace.default.instance_id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = "my-repo"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}

resource "alicloud_cr_chain" "default" {
  chain_name          = var.name
  description         = "description"
  instance_id         = alicloud_cr_ee_namespace.default.instance_id
  repo_name           = alicloud_cr_ee_repo.default.name
  repo_namespace_name = alicloud_cr_ee_namespace.default.name
  chain_config {
    routers {
      from {
        node_name = "DOCKER_IMAGE_BUILD"
      }
      to {
        node_name = "DOCKER_IMAGE_PUSH"
      }
    }
    routers {
      from {
        node_name = "DOCKER_IMAGE_PUSH"
      }
      to {
        node_name = "VULNERABILITY_SCANNING"
      }
    }
    routers {
      from {
        node_name = "VULNERABILITY_SCANNING"
      }
      to {
        node_name = "ACTIVATE_REPLICATION"
      }
    }
    routers {
      from {
        node_name = "ACTIVATE_REPLICATION"
      }
      to {
        node_name = "TRIGGER"
      }
    }
    routers {
      from {
        node_name = "VULNERABILITY_SCANNING"
      }
      to {
        node_name = "SNAPSHOT"
      }
    }
    routers {
      from {
        node_name = "SNAPSHOT"
      }
      to {
        node_name = "TRIGGER_SNAPSHOT"
      }
    }

    nodes {
      enable    = true
      node_name = "DOCKER_IMAGE_BUILD"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "DOCKER_IMAGE_PUSH"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "VULNERABILITY_SCANNING"
      node_config {
        deny_policy {
          issue_level = "MEDIUM"
          issue_count = 1
          action      = "BLOCK_DELETE_TAG"
          logic       = "AND"
        }
      }
    }
    nodes {
      enable    = true
      node_name = "ACTIVATE_REPLICATION"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "TRIGGER"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = false
      node_name = "SNAPSHOT"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = false
      node_name = "TRIGGER_SNAPSHOT"
      node_config {
        deny_policy {}
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `chain_name` - (Required) The name of delivery chain. The length of the name is 1-64 characters, lowercase English letters and numbers, and the separators "_", "-", "." can be used, noted that the separator cannot be at the first or last position.
* `description` - (Optional) The description delivery chain.
* `repo_name` - (Optional, ForceNew) The name of CR Enterprise Edition repository. **NOTE:** This parameter must specify a correct value, otherwise the created resource will be incorrect.
* `repo_namespace_name` - (Optional, ForceNew) The name of CR Enterprise Edition namespace. **NOTE:** This parameter must specify the correct value, otherwise the created resource will be incorrect.
* `instance_id` - (Required, ForceNew) The ID of CR Enterprise Edition instance.
* `chain_config` - (Optional) The configuration of delivery chain. **NOTE:** This parameter must specify the correct value, otherwise the created resource will be incorrect.

#### Block chain_config

The `chain_config` block supports the following:
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

-> **NOTE:** The `from` and `to` fields are all fixed, and their structure and the value of `node_name` are fixed. You can refer to the template given in the example for configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chain. The value formats as `<instance_id>:<chain_id>`.
* `chain_id` - Delivery chain ID.

## Import

CR Chain can be imported using the id, e.g.

```
$ terraform import alicloud_cr_chain.example <instance_id>:<chain_id>
```