---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ack_inspect_config"
description: |-
  Provides a Alicloud Container Service for Kubernetes (ACK) Inspect Config resource.
---

# alicloud_ack_inspect_config

Provides a Container Service for Kubernetes (ACK) Inspect Config resource.

cluster inspect config.

For information about Container Service for Kubernetes (ACK) Inspect Config and how to use it, see [What is Inspect Config](https://next.api.alibabacloud.com/document/CS/2015-12-15/CreateClusterInspectConfig).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cs_managed_kubernetes" "创建Cluster" {
  addons {
    name   = "terway-eniip"
    config = "{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"true\"}"
  }
  addons {
    name   = "terway-controlplane"
    config = "{\"ENITrunking\":\"true\"}"
  }
  addons {
    name = "csi-plugin"
  }
  addons {
    name = "managed-csiprovisioner"
  }
  addons {
    name = "nginx-ingress-controller"
  }
  addons {
    name = "metrics-server"
  }
  addons {
    name = "coredns"
  }
  ip_stack                     = "ipv4"
  is_enterprise_security_group = true
  service_cidr                 = var.service_cidr
  proxy_mode                   = "ipvs"
  deletion_protection          = false
  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable           = true
    maintenance_time = "2025-11-03T00:00:00.000+08:00"
    duration         = "3h"
    weekly_period    = "Monday"
  }
  zone_ids = [data.alicloud_zones.default.zones.0.id]
}


resource "alicloud_ack_inspect_config" "default" {
  recurrence           = "FREQ=DAILY;BYHOUR=10;BYMINUTE=15"
  disabled_check_items = ["APIServerCLBListenerAbnormal"]
  enabled              = true
  inspect_config_id    = alicloud_cs_managed_kubernetes.创建Cluster.id
}
```

## Argument Reference

The following arguments are supported:
* `disabled_check_items` - (Optional, List) List of disabled inspection items.
* `enabled` - (Required) Whether to enable patrol inspection.
* `inspect_config_id` - (Required, ForceNew) The first ID of the resource
* `recurrence` - (Required) Use the RFC5545 Recurrence Rule syntax to define the inspection cycle. BYHOUR and BYMINUTE must be specified. Only FREQ = DAILY is supported, and COUNT or UNTIL is not supported.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Inspect Config.
* `delete` - (Defaults to 5 mins) Used when delete the Inspect Config.
* `update` - (Defaults to 5 mins) Used when update the Inspect Config.

## Import

Container Service for Kubernetes (ACK) Inspect Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_ack_inspect_config.example <inspect_config_id>
```