---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_nlb"
sidebar_current: "docs-alicloud-datasource-sae-load-balancer-nlb"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Load Balancer NLB data source.
---

# alicloud_sae_load_balancer_nlb

Provides a Serverless App Engine (SAE) Load Balancer NLB data source.

For information about SAE NLB and how to use it, see [Describe application NLB](https://www.alibabacloud.com/help/en/sae/developer-reference/api-sae-2019-05-06-describeapplicationnlbs).

-> **NOTE:** Available since v1.293.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_sae_load_balancer_nlb" "default" {
  app_id = "example-app-id"
}

output "nlb_dns_name" {
  value = data.alicloud_sae_load_balancer_nlb.default.dns_name
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the SAE application to query.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source, which is the app ID.
* `dns_name` - The DNS name of the NLB instance bound to the application.
* `created_by_sae` - Whether the NLB instance was created by SAE.
* `listeners` - A list of NLB listener configurations.
  * `port` - The NLB listening port.
  * `target_port` - The container port that the traffic is forwarded to.
  * `protocol` - The protocol of the listener.
  * `cert_ids` - The certificate IDs for the TCPSSL protocol.
