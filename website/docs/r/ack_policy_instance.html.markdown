---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ack_policy_instance"
description: |-
  Provides a Alicloud Container Service for Kubernetes (ACK) Policy Instance resource.
---

# alicloud_ack_policy_instance

Provides a Container Service for Kubernetes (ACK) Policy Instance resource.



For information about Container Service for Kubernetes (ACK) Policy Instance and how to use it, see [What is Policy Instance](https://next.api.alibabacloud.com/document/CS/2015-12-15/DeployPolicyInstance).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

variable "policy-template-controller_version" {
  default = "v0.4.0.0-gddee19d-aliyun"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "cluster_name" {
  default = "example-create-cluster"
}

variable "policy_name" {
  default = "ACKAllowedRepos"
}

variable "policy_scope" {
  default = "default"
}

variable "policy_scope-kube-public" {
  default = "kube-public"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "policy_scope_update" {
  default = "kube-system"
}

variable "cidr_block" {
  default = "172.18.0.0/21"
}

variable "gatekeeper_version" {
  default = "3.18.2-release"
}

variable "loongcollector_version" {
  default = "3.1.6"
}

resource "alicloud_vpc" "创建VPC" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "创建VSwitch" {
  vpc_id     = alicloud_vpc.创建VPC.id
  zone_id    = var.zone_id
  cidr_block = var.cidr_block
}

resource "alicloud_cs_managed_kubernetes" "创建Cluster" {
  addons {
    name     = "gatekeeper"
    disabled = false
  }
  addons {
    name     = "loongcollector"
    disabled = false
  }
  addons {
    name     = "policy-template-controller"
    disabled = false
  }
  addons {
    name     = "terway-eniip"
    config   = "{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"true\"}"
    disabled = false
  }
  addons {
    name     = "terway-controlplane"
    config   = "{\"ENITrunking\":\"true\"}"
    disabled = false
  }
  addons {
    name     = "coredns"
    disabled = false
  }
  addons {
    name     = "metrics-server"
    disabled = false
  }
  addons {
    name     = "nginx-ingress-controller"
    disabled = false
  }
  addons {
    name     = "managed-csiprovisioner"
    disabled = false
  }
  addons {
    name     = "csi-plugin"
    disabled = false
  }
  addons {
    name     = "storage-operator"
    disabled = false
  }
  is_enterprise_security_group = true
  vswitch_ids                  = ["${alicloud_vswitch.创建VSwitch.id}"]
  service_cidr                 = var.service_cidr
  pod_vswitch_ids              = ["${alicloud_vswitch.创建VSwitch.id}"]
  ip_stack                     = "ipv4"
  proxy_mode                   = "ipvs"
  deletion_protection          = false
  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable = false
  }
  profile      = "Default"
  cluster_spec = "ack.pro.small"
}


resource "alicloud_ack_policy_instance" "default" {
  policy_action = "deny"
  namespaces    = []
  parameters {
  }
  cluster_id  = alicloud_cs_managed_kubernetes.创建Cluster.id
  policy_name = var.policy_name
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required, ForceNew) Target cluster ID
* `namespaces` - (Optional, Computed, List) Limits the namespace of the policy implementation. Empty indicates all namespaces.
* `parameters` - (Optional, Map) The parameter configuration of the current rule instance. For more information about the parameters supported by each policy rule, see [Container Security Policy Rule Base Description](https://www.alibabacloud.com/help/doc-detail/359819.html).
* `policy_action` - (Optional) Policy Governance Implementation Actions
* `policy_name` - (Required, ForceNew) Policy Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<policy_name>:<instance_name>`.
* `instance_name` - Rule Instance Name

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Policy Instance.
* `update` - (Defaults to 5 mins) Used when update the Policy Instance.

## Import

Container Service for Kubernetes (ACK) Policy Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ack_policy_instance.example <cluster_id>:<policy_name>:<instance_name>
```
