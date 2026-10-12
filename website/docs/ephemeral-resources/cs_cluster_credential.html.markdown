---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_cluster_credential"
sidebar_current: "docs-alicloud-ephemeral-cs-cluster-credential"
description: |-
    Fetches an ACK cluster kubeconfig without persisting the credential to the state or plan.
---

# Ephemeral: alicloud_cs_cluster_credential

This ephemeral resource fetches an ACK cluster's kubeconfig through the `DescribeClusterUserKubeconfig` API without persisting the credential to the state or plan. The kubeconfig is available for the duration of a single Terraform operation and can be referenced in provider configuration, locals, provisioners, and other ephemeral resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.0.0-beta5. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The ephemeral resource can be used on all kinds of ACK clusters, including managed clusters, imported kubernetes clusters, serverless clusters and edge clusters. Please make sure that the target cluster is not in the failed state, since the API server of clusters in the failed state cannot be accessed.

-> **NOTE:** The ephemeral resource is opened during the plan phase when `cluster_id` is known at plan time, so the cluster must already exist; if the cluster does not exist, the open fails with a not-found error such as `ErrorClusterNotFound`. When `cluster_id` references another resource in the same configuration, its value is unknown during plan and the open is deferred to the apply phase automatically. The credential is re-fetched on every operation that opens it, and a kubeconfig issued with `temporary_duration_minutes` expires at the reported `expiration`.

-> **NOTE:** The
[`alicloud_cs_cluster_credential`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential)
data source persists the kubeconfig into the state and can also write it to disk through `output_file`;
use this ephemeral resource when the credential must never appear in any Terraform artifact.

## Example Usage

Configure the kubernetes provider with an ephemeral kubeconfig:

```terraform
ephemeral "alicloud_cs_cluster_credential" "example" {
  cluster_id = "c-xxxxxxxxxxxxxxxx"
}

provider "kubernetes" {
  host                   = yamldecode(ephemeral.alicloud_cs_cluster_credential.example.kube_config)["clusters"][0]["cluster"]["server"]
  cluster_ca_certificate = base64decode(yamldecode(ephemeral.alicloud_cs_cluster_credential.example.kube_config)["clusters"][0]["cluster"]["certificate-authority-data"])
  client_certificate     = base64decode(yamldecode(ephemeral.alicloud_cs_cluster_credential.example.kube_config)["users"][0]["user"]["client-certificate-data"])
  client_key             = base64decode(yamldecode(ephemeral.alicloud_cs_cluster_credential.example.kube_config)["users"][0]["user"]["client-key-data"])
}
```

Fetch a short-lived credential that expires at the reported `expiration`:

```terraform
ephemeral "alicloud_cs_cluster_credential" "example" {
  cluster_id                 = "c-xxxxxxxxxxxxxxxx"
  temporary_duration_minutes = 60
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the target cluster.
* `temporary_duration_minutes` - (Optional) Automatic expiration time of the returned credential. The valid value between `15` and `4320`, in minutes. When this field is omitted, the expiration time is determined by the system automatically and reported in `expiration`.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `kube_config` - (Sensitive) The kubeconfig to use to authenticate with the cluster.
* `expiration` - Expiration time of the kubeconfig. Format: UTC time in RFC 3339.
