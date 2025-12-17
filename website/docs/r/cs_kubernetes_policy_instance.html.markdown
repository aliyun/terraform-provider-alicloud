---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_policy_instance"
description: |-
  Provides a Alicloud Container Service for Kubernetes (ACK) Policy Instance resource.
---

# alicloud_cs_kubernetes_policy_instance

Provides a Container Service for Kubernetes (ACK) Policy Instance resource.

For information about Container Service for Kubernetes (ACK) Policy Instance and how to use it, see [What is Policy Instance](https://next.api.alibabacloud.com/document/CS/2015-12-15/DeployPolicyInstance).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "vpc_cidr" {
  default = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "cluster_name" {
  default = "example-create-cluster"
}

variable "pod_cidr" {
  default = "172.16.0.0/16"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

variable "policy_name" {
  default = "ACKPSPHostNetworkingPorts"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_vpc" "CreateVPC" {
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "CreateVSwitch" {
  count      = length(var.vswitch_cidrs)
  vpc_id     = alicloud_vpc.CreateVPC.id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "CreateCluster" {
  name                         = "${var.cluster_name}-${random_integer.default.result}"
  cluster_spec                 = "ack.pro.small"
  profile                      = "Default"
  vswitch_ids                  = split(",", join(",", alicloud_vswitch.CreateVSwitch.*.id))
  pod_cidr                     = var.pod_cidr
  service_cidr                 = var.service_cidr
  is_enterprise_security_group = true
  ip_stack                     = "ipv4"
  proxy_mode                   = "ipvs"
  deletion_protection          = false

  addons {
    name = "gatekeeper"
  }
  addons {
    name = "loongcollector"
  }
  addons {
    name = "policy-template-controller"
  }

  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable = false
  }
}

resource "alicloud_cs_kubernetes_policy_instance" "default" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = var.policy_name
  action      = "deny"
  namespaces = [
    "test"
  ]
  parameters = {
    hostNetwork = true
    min         = 20
    max         = 200
  }
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required, ForceNew) Target cluster ID
* `namespaces` - (Optional, Computed, List) Limits the namespace of the policy implementation. Empty indicates all namespaces.
* `parameters` - (Optional, Map) The parameter configuration of the current rule instance. For more information about the parameters supported by each policy rule, see [Container Security Policy Rule Base Description](https://www.alibabacloud.com/help/doc-detail/359819.html).
* `action` - (Optional) Policy Governance Implementation Actions
* `policy_name` - (Required, ForceNew) Policy Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<cluster_id>:<policy_name>:<instance_name>`.
* `instance_name` - Rule Instance Name

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Policy Instance.
* `update` - (Defaults to 5 mins) Used when update the Policy Instance.

## Import

Container Service for Kubernetes (ACK) Policy Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cs_kubernetes_policy_instance.example <cluster_id>:<policy_name>:<instance_name>
```
