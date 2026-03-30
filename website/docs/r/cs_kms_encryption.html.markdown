---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kms_encryption"
description: |-
  Provides a Alicloud resource to use KMS to encrypt Kubernetes Secrets in ACK cluster.
---

# alicloud_cs_kms_encryption

This resource will help you to use KMS to encrypt Kubernetes Secrets in ACK cluster.

For information about how to use it, see [Use KMS to encrypt Kubernetes Secrets](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/security-and-compliance/use-kms-to-encrypt-kubernetes-secrets-2)

-> **NOTE:** Available since v1.274.0.

-> **NOTE:** During the process of enabling or disabling encryption at rest, and after successfully enabling this feature, do not use the KMS console or OpenAPI to disable or delete the KMS key used by this feature. Otherwise, the cluster API Server will become unavailable, and you will be unable to retrieve Secrets and ServiceAccounts, which will affect the normal operation of your business applications. For more information, see [Use Alibaba Cloud KMS to encrypt Kubernetes Secrets](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/security-and-compliance/use-kms-to-encrypt-kubernetes-secrets-2).

-> **NOTE:** Users or roles who use the encryption at rest feature need to be granted additional cluster RBAC permissions (operations or administrator permissions are required). Otherwise, the "ForbiddenUpdateKMSState" error code will be returned.

-> **NOTE:** After successfully configuring encryption at rest, the cluster status will change to "Updating". Once the change is complete, the cluster status will return to "Running". After the change for the same cluster is complete, you must wait at least 5 minutes before calling this API again. Otherwise, a "409" status code will be returned.

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
  default = "terraform-example-"
}

variable "pod_cidr" {
  default = "172.16.0.0/16"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

data "alicloud_kms_keys" "kms_keys_ds" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]},{\"Key\":\"KeySpec\",\"Values\":[\"Aliyun_AES_256\"]},{\"Key\":\"KeyUsage\",\"Values\":[\"ENCRYPT/DECRYPT\"]},{\"Key\":\"CreatorType\",\"Values\":[\"User\"]}]"
}

resource "alicloud_vpc" "CreateVPC" {
  cidr_block = var.vpc_cidr
}

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

resource "alicloud_cs_kms_encryption" "default" {
  cluster_id         = alicloud_cs_managed_kubernetes.CreateCluster.id
  disable_encryption = false
  kms_key_id         = data.alicloud_kms_keys.kms_keys_ds.keys.0.id

  provisioner "local-exec" {
    command = "echo 'wait for task conflict...' && sleep 310"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `disable_encryption` - (Required) Whether to disable KMS encryption. Valid values: `true`, `false`.
* `kms_key_id` - (Optional) The KMS key ID used for encryption. Required when `disable_encryption` is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, which is the same as the cluster ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when enabling KMS encryption for the kubernetes cluster.
* `update` - (Defaults to 10 mins) Used when updating KMS encryption configuration.
* `delete` - (Defaults to 10 mins) Used when disabling KMS encryption.

## Import

KMS encryption for ACK cluster can be imported using the cluster ID.

```shell
$ terraform import alicloud_cs_kms_encryption.my_encryption <cluster_id>
```
