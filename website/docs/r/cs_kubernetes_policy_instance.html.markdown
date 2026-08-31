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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_kubernetes_policy_instance&exampleId=f7586c37-a483-bb38-9ed8-f40c7cf24848943a61cb&activeTab=example&spm=docs.r.cs_kubernetes_policy_instance.0.f7586c37a4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "vpc_cidr" {
  default = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "cluster_name" {
  default = "terraform-example-"
}

variable "pod_cidr" {
  default = "172.16.0.0/16"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

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
  name_prefix                  = var.cluster_name
  cluster_spec                 = "ack.standard"
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

resource "alicloud_cs_kubernetes_policy_instance" "base" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = "ACKPSPReadOnlyRootFilesystem"
}

resource "alicloud_cs_kubernetes_policy_instance" "string" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = "ACKPVSizeConstraint"
  action      = "deny"
  parameters = {
    maxSize = "60Gi"
  }
}

resource "alicloud_cs_kubernetes_policy_instance" "int_bool" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = "ACKPSPHostNetworkingPorts"
  namespaces = [
    "test"
  ]
  parameters = {
    hostNetwork = true
    min         = 20
    max         = 200
  }
}

resource "alicloud_cs_kubernetes_policy_instance" "array" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = "ACKAllowedRepos"
  parameters = {
    repos = jsonencode([
      "docker.io/library/nginx",
      "docker.io/library/redis"
    ])
  }
}

resource "alicloud_cs_kubernetes_policy_instance" "object" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = "ACKRequiredLabels"
  action      = "warn"
  namespaces = [
    "test1",
    "test2",
    "test3"
  ]
  parameters = {
    labels = jsonencode([
      {
        key          = "test"
        allowedRegex = "^test.*$"
      },
      {
        key          = "env"
        allowedRegex = "^(dev|prod)$"
        optional     = false
      }
    ])
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cs_kubernetes_policy_instance&spm=docs.r.cs_kubernetes_policy_instance.example&intl_lang=EN_US)

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
