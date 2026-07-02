---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chain"
description: |-
  Provides a Alicloud CR Chain resource.
---

# alicloud_cr_chain

Provides a CR Chain resource.



For information about CR Chain and how to use it, see [What is Chain](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchain).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 100000
  max = 999999
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = "${var.name}-${random_integer.default.result}"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}

resource "alicloud_cr_chain" "default" {
  chain_name          = "${var.name}-${random_integer.default.result}"
  description         = var.name
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
* `chain_config` - (Optional, Set) Delivery chain configuration description See [`chain_config`](#chain_config) below.
* `chain_name` - (Required) Delivery chain name
* `description` - (Optional) Delivery chain description
* `instance_id` - (Required, ForceNew) Instance ID
* `repo_name` - (Required, ForceNew) Warehouse name
* `repo_namespace_name` - (Required, ForceNew) Namespace name

### `chain_config`

The chain_config supports the following:
* `nodes` - (Optional, List) Each node in the delivery chain See [`nodes`](#chain_config-nodes) below.
* `routers` - (Optional, List) Execution sequence relationship between delivery chain nodes See [`routers`](#chain_config-routers) below.

### `chain_config-nodes`

The chain_config-nodes supports the following:
* `enable` - (Optional) Whether to enable the delivery chain node. Valid values:
  -'true': Enable delivery chain nodes
  -'false': do not enable the delivery chain node
* `node_config` - (Optional, ForceNew, Set) Delivery chain node configuration See [`node_config`](#chain_config-nodes-node_config) below.
* `node_name` - (Optional) Delivery chain node name

### `chain_config-routers`

The chain_config-routers supports the following:
* `from` - (Optional, ForceNew, Set) Source node See [`from`](#chain_config-routers-from) below.
* `to` - (Optional, ForceNew, Set) Destination node See [`to`](#chain_config-routers-to) below.

### `chain_config-routers-from`

The chain_config-routers-from supports the following:
* `node_name` - (Optional) Source node name

### `chain_config-routers-to`

The chain_config-routers-to supports the following:
* `node_name` - (Optional) Destination node name

### `chain_config-nodes-node_config`

The chain_config-nodes-node_config supports the following:
* `scan_engine` - (Optional, Available since v1.283.0) Delivery chain scan node engine

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<chain_id>`.
* `chain_config` - Delivery chain configuration description.
  * `chain_config_id` - Delivery chain configuration ID.
  * `is_active` - Whether the delivery chain configuration takes effect.
  * `nodes` - Each node in the delivery chain.
    * `node_config` - Delivery chain node configuration.
        * `deny_policy` - Blocking Rules for Scanning Nodes in Delivery Chain Nodes.
            * `action` - Blocking action, value:.
            * `issue_count` - Trigger blocking when the number of scanning vulnerabilities reaches.
            * `issue_level` - Trigger blocking when scanning vulnerability level reaches.
            * `logic` - Scan logic that triggers blocking.
        * `retry` - Number of retries.
        * `timeout` - Timeout.
  * `version` - Delivery chain version.
* `chain_id` - Delivery chain ID.
* `code` - Return code.
* `create_time` - Delivery chain creation time.
* `is_success` - Success.
* `modified_time` - Delivery chain description modification time.
* `scope_id` - Delivery chain scope ID.
* `scope_type` - Delivery chain scope type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Chain.
* `delete` - (Defaults to 5 mins) Used when delete the Chain.
* `update` - (Defaults to 5 mins) Used when update the Chain.

## Import

CR Chain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_chain.example <instance_id>:<chain_id>
```