---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_cluster_scanner_yaml"
description: |-
  Provides a Alicloud Threat Detection Cluster Scanner Yaml resource.
---

# alicloud_threat_detection_cluster_scanner_yaml

Provides a Threat Detection Cluster Scanner Yaml resource.

Generates the k8s cluster scanner webhook yaml configuration for a Security Center (Sas) managed cluster.

For information about Threat Detection Cluster Scanner Yaml and how to use it, see [GenerateClusterScannerWebhookYaml](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-generateclusterscannerwebhookyaml).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_cluster_scanner_yaml" "default" {
  cluster_id   = "cxxxxxxxxxxxxxxxxx"
  webhook_open = 1
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required, ForceNew) The ID of the k8s cluster to which the scan component needs to be deployed.
* `webhook_open` - (Optional, ForceNew, Int) Specifies the webhook enabled state of the scan component. Valid values:
  - `1`: webhook needs to be enabled.
  - `0`: the webhook needs to be closed.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. It is the cluster ID.
* `ca_cert_base64` - The base64-encoded CA certificate of the cluster scan component.
* `tls_key_base64` - The base64-encoded TLS private key of the cluster scan component.
* `tls_cert_base64` - The base64-encoded TLS certificate corresponding to the cluster scan component.
* `cluster_env_info` - The cluster environment information.
* `image` - The image address corresponding to the cluster scan component.
* `region_id` - The region ID of the resource.

-> **NOTE:** This resource only supports create and read. The GenerateClusterScannerWebhookYaml API does not provide a delete operation, so removing the resource from Terraform state does not modify the cloud-side scanner configuration. Changes to `cluster_id` or `webhook_open` force a new resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cluster Scanner Yaml.
* `delete` - (Defaults to 5 mins) Used when delete the Cluster Scanner Yaml.

## Import

Threat Detection Cluster Scanner Yaml can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_cluster_scanner_yaml.example <cluster_id>
```
