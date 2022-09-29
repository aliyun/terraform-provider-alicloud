---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_cluster_credential"
sidebar_current: "docs-alicloud-datasource-cs-cluster-credential"
description: |-
  Provides a list of Container Service Cluster's credential to be used by the alicloud_cs_cluster_credential resource.
---

# alicloud\_cs\_cluster\_credential

This data source provides Container Service cluster credential on Alibaba Cloud.

-> **NOTE:** Available in v1.187.0+

-> **NOTE:** This datasource can be used on all kinds of ACK clusters, including managed clusters, imported kubernetes clusters, serverless clusters and edge clusters. Please make sure that the target cluster is not in the failed state before using this datasource, since the api server of clusters in the failed state cannot be accessed.

## Example Usage

```
# Declare the data source
data "alicloud_cs_cluster_credential" "auth" {
  cluster_id                 = "cluster-id"
  temporary_duration_minutes = 60
}
```

```
# Declare the data source
data "alicloud_cs_managed_kubernetes_clusters" "k8s"{
  name_regex     = "my-cluster"
  enable_details = false
}

data "alicloud_cs_cluster_credential" "auth" {
  for_each                   = toset(data.alicloud_cs_managed_kubernetes_clusters.k8s.ids)
  cluster_id                 = each.key
  temporary_duration_minutes = 60
  output_file                = "my-auth-json"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (**Required**) The id of target cluster.
* `temporary_duration_minutes` - (Optional) Automatic expiration time of the returned credential. The valid value between `15` and `4320`, in minutes. When this field is omitted, the expiration time will be determined by the system automatically and the result will be in the attributed field `expiration`.
* `output_file` - (Optional) File name where to save the returned KubeConfig (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `cluster_id` - The id of target cluster.
* `cluster_name` - The name of target cluster.
* `kube_config` - (Sensitive) The kube config to use to authenticate with the cluster.
* `certificate_authority` - (Available in 1.105.0+) Nested attribute containing certificate authority data for your cluster.
  * `cluster_cert` - The base64 encoded cluster certificate data required to communicate with your cluster. Add this to the certificate-authority-data section of the kube config file for your cluster.
  * `client_cert` - The base64 encoded client certificate data required to communicate with your cluster. Add this to the client-certificate-data section of the kube config file for your cluster.
  * `client_key` - The base64 encoded client key data required to communicate with your cluster. Add this to the client-key-data section of the kube config file for your cluster.
* `expiration` - Expiration time of kube config. Format: UTC time in rfc3339.
