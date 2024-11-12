---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chain"
sidebar_current: "docs-alicloud-resource-cr-chain"
description: |-
  Provides a Alicloud CR Chain resource.
---

# alicloud_cr_chain

Provides a CR Chain resource.

For information about CR Chain and how to use it, see [What is Chain](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchain).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cr_chain&exampleId=045b97ca-e86a-dc23-3218-7c62847ab87ccfd100d3&activeTab=example&spm=docs.r.cr_chain.0.045b97cae8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

* `chain_name` - (Required) The name of delivery chain. The length of the name is 1-64 characters, lowercase English letters and numbers, and the separators "_", "-", "." can be used, noted that the separator cannot be at the first or last position.
* `description` - (Optional) The description delivery chain.
* `repo_name` - (Optional, ForceNew) The name of CR Enterprise Edition repository. **NOTE:** This parameter must specify a correct value, otherwise the created resource will be incorrect.
* `repo_namespace_name` - (Optional, ForceNew) The name of CR Enterprise Edition namespace. **NOTE:** This parameter must specify the correct value, otherwise the created resource will be incorrect.
* `instance_id` - (Required, ForceNew) The ID of CR Enterprise Edition instance.
* `chain_config` - (Optional) The configuration of delivery chain. See [`chain_config`](#chain_config) below. **NOTE:** This parameter must specify the correct value, otherwise the created resource will be incorrect.

### `chain_config`

The `chain_config` block supports the following:

* `routers` - (Optional) Execution sequence relationship between delivery chain nodes. See [`routers`](#chain_config-routers) below. 
* `nodes` - (Optional) Each node in the delivery chain. See [`nodes`](#chain_config-nodes) below.

-> **NOTE:** The `from` and `to` fields are all fixed, and their structure and the value of `node_name` are fixed. You can refer to the template given in the example for configuration.

### `chain_config-routers`

The `routers` block supports the following:
* `from` - (Optional) Source node. See [`from`](#chain_config-routers-from) below.
* `to` - (Optional) Destination node. See [`to`](#chain_config-routers-to) below.

### `chain_config-routers-from`

The `from` block supports the following:
* `node_name` - (Optional) The name of node. Valid values: `DOCKER_IMAGE_BUILD`, `DOCKER_IMAGE_PUSH`, `VULNERABILITY_SCANNING`, `ACTIVATE_REPLICATION`, `TRIGGER`, `SNAPSHOT`, `TRIGGER_SNAPSHOT`.

### `chain_config-routers-to`

The `to` block supports the following:
* `node_name` - (Optional) The name of node. Valid values: `DOCKER_IMAGE_BUILD`, `DOCKER_IMAGE_PUSH`, `VULNERABILITY_SCANNING`, `ACTIVATE_REPLICATION`, `TRIGGER`, `SNAPSHOT`, `TRIGGER_SNAPSHOT`.

### `chain_config-nodes`

The `nodes` block supports the following:
* `node_name` - (Optional) The name of delivery chain node.
* `enable` - (Optional) Whether to enable the delivery chain node. Valid values: `true`, `false`.
* `node_config` - (Optional) The configuration of delivery chain node. See [`node_config`](#chain_config-nodes-node_config) below.

### `chain_config-nodes-node_config`

The `node_config` block supports the following:
* `deny_policy` - (Optional) Blocking rules for scanning nodes in delivery chain nodes. See [`deny_policy`](#chain_config-nodes-node_config-deny_policy) below. **Note:** When `node_name` is `VULNERABILITY_SCANNING`, the parameters in `deny_policy` need to be filled in.

### `chain_config-nodes-node_config-deny_policy`

The `deny_policy` block supports the following:
* `issue_count` - (Optional) The count of scanning vulnerabilities that triggers blocking.
* `issue_level` - (Optional) The level of scanning vulnerability that triggers blocking. Valid values: `LOW`, `MEDIUM`, `HIGH`, `UNKNOWN`.
* `logic` - (Optional) The logic of trigger blocking. Valid values: `AND`, `OR`.
* `action` - (Optional) The action of trigger blocking. Valid values: `BLOCK`, `BLOCK_RETAG`, `BLOCK_DELETE_TAG`. While `Block` means block the delivery chain from continuing to execute, `BLOCK_RETAG` means block overwriting push image tag, `BLOCK_DELETE_TAG` means block deletion of mirror tags.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chain. The value formats as `<instance_id>:<chain_id>`.
* `chain_id` - Delivery chain ID.

## Import

CR Chain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_chain.example <instance_id>:<chain_id>
```