---
subcategory: ""
layout: "alicloud"
page_title: "Terraform Alicloud Provider Version 2 Upgrade Guide"
description: |-
  Terraform Alicloud Provider Version 2 Upgrade Guide
---

# Terraform AliCloud Provider Version 2 Upgrade Guide

Version 2.0.0 of the AliCloud provider for Terraform is a major release and includes some changes that you need to consider when upgrading. This guide is intended to help with that process and focuses only on changes from version 1.x.x to version 2.0.0.

Upgrade topics:

- [Provider Version Configuration](#provider-version-configuration)
- [Resource: alicloud_api_gateway_instance](#resource-alicloud_api_gateway_instance)
- [Resource: alicloud_cr_repo](#resource-alicloud_cr_repo)
- [Resource: alicloud_cs_edge_kubernetes](#resource-alicloud_cs_edge_kubernetes)
- [Resource: alicloud_cs_kubernetes](#resource-alicloud_cs_kubernetes)
- [Resource: alicloud_cs_managed_kubernetes](#resource-alicloud_cs_managed_kubernetes)
- [Data Source: alicloud_cr_repos](#data-source-alicloud_cr_repos)
- [Data Source: alicloud_cs_cluster_credential](#data-source-alicloud_cs_cluster_credential)
- [Data Source: alicloud_cs_edge_kubernetes_clusters](#data-source-alicloud_cs_edge_kubernetes_clusters)
- [Data Source: alicloud_cs_kubernetes_clusters](#data-source-alicloud_cs_kubernetes_clusters)
- [Data Source: alicloud_cs_managed_kubernetes_clusters](#data-source-alicloud_cs_managed_kubernetes_clusters)
- [Data Source: alicloud_cs_serverless_kubernetes_clusters](#data-source-alicloud_cs_serverless_kubernetes_clusters)
- [Data Source: alicloud_db_instance_classes](#data-source-alicloud_db_instance_classes)
- [Data Source: alicloud_instance_types](#data-source-alicloud_instance_types)

## Provider Version Configuration

-> Before upgrading to `v2.0.0` or later, it is recommended to upgrade to the v1.X of the provider (v1.282.0) and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) without unexpected changes or deprecation notices.

We recommend using [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init --upgrade`](https://www.terraform.io/docs/commands/init.html) to download the new version.

Update to latest 1.X version:

```terraform
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.282.0"
    }
  }
}
```

Update to latest 2.X version:

```terraform
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 2.0.0"
    }
  }
}
```

## Resource: alicloud_api_gateway_instance

### to_connect_vpc_ip_block Argument Type Change

The `to_connect_vpc_ip_block` argument has changed from a `TypeMap` to a `TypeList`. Update the configuration from map syntax to block syntax, and update any attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
resource "alicloud_api_gateway_instance" "example" {
  # ... other configuration ...

  to_connect_vpc_ip_block = {
    cidr_block = "172.16.0.0/12"
    vswitch_id = "vsw-abc123"
    zone_id    = "cn-hangzhou-a"
  }
}

output "cidr_block" {
  value = alicloud_api_gateway_instance.example.to_connect_vpc_ip_block["cidr_block"]
}
```

Updated configuration:

```terraform
resource "alicloud_api_gateway_instance" "example" {
  # ... other configuration ...

  to_connect_vpc_ip_block {
    cidr_block = "172.16.0.0/12"
    vswitch_id = "vsw-abc123"
    zone_id    = "cn-hangzhou-a"
  }
}

output "cidr_block" {
  value = alicloud_api_gateway_instance.example.to_connect_vpc_ip_block[0].cidr_block
}
```

## Resource: alicloud_cr_repo

### domain_list Attribute Type Change

The `domain_list` attribute has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
resource "alicloud_cr_repo" "example" {
  namespace = "my-namespace"
  name      = "my-repo"
  summary   = "Chinese mainland China regions"
  repo_type = "PUBLIC"
}

output "public_domain" {
  value = alicloud_cr_repo.example.domain_list["public"]
}

output "internal_domain" {
  value = alicloud_cr_repo.example.domain_list["internal"]
}

output "vpc_domain" {
  value = alicloud_cr_repo.example.domain_list["vpc"]
}
```

Updated configuration:

```terraform
resource "alicloud_cr_repo" "example" {
  namespace = "my-namespace"
  name      = "my-repo"
  summary   = "Chinese mainland China regions"
  repo_type = "PUBLIC"
}

output "public_domain" {
  value = alicloud_cr_repo.example.domain_list[0].public
}

output "internal_domain" {
  value = alicloud_cr_repo.example.domain_list[0].internal
}

output "vpc_domain" {
  value = alicloud_cr_repo.example.domain_list[0].vpc
}
```

## Resource: alicloud_cs_edge_kubernetes

### runtime, certificate_authority, connections Type Change

The `runtime`, `certificate_authority`, and `connections` attributes have changed from `TypeMap` to `TypeList`. Update the configuration from map syntax to block syntax, and update any attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
resource "alicloud_cs_edge_kubernetes" "example" {
  # ... other configuration ...

  runtime = {
    name    = "containerd"
    version = "1.5.13"
  }
}

output "cluster_cert" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority["cluster_cert"]
}

output "client_cert" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority["client_cert"]
}

output "client_key" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority["client_key"]
}

output "api_server_internet" {
  value = alicloud_cs_edge_kubernetes.example.connections["api_server_internet"]
}

output "api_server_intranet" {
  value = alicloud_cs_edge_kubernetes.example.connections["api_server_intranet"]
}

output "master_public_ip" {
  value = alicloud_cs_edge_kubernetes.example.connections["master_public_ip"]
}

output "service_domain" {
  value = alicloud_cs_edge_kubernetes.example.connections["service_domain"]
}
```

Updated configuration:

```terraform
resource "alicloud_cs_edge_kubernetes" "example" {
  # ... other configuration ...

  runtime {
    name    = "containerd"
    version = "1.5.13"
  }
}

output "cluster_cert" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority[0].cluster_cert
}

output "client_cert" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority[0].client_cert
}

output "client_key" {
  value = alicloud_cs_edge_kubernetes.example.certificate_authority[0].client_key
}

output "api_server_internet" {
  value = alicloud_cs_edge_kubernetes.example.connections[0].api_server_internet
}

output "api_server_intranet" {
  value = alicloud_cs_edge_kubernetes.example.connections[0].api_server_intranet
}

output "master_public_ip" {
  value = alicloud_cs_edge_kubernetes.example.connections[0].master_public_ip
}

output "service_domain" {
  value = alicloud_cs_edge_kubernetes.example.connections[0].service_domain
}
```

## Resource: alicloud_cs_kubernetes

### runtime, certificate_authority, connections Type Change

The `runtime`, `certificate_authority`, and `connections` attributes have changed from `TypeMap` to `TypeList`. Update the configuration from map syntax to block syntax, and update any attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
resource "alicloud_cs_kubernetes" "example" {
  # ... other configuration ...

  runtime = {
    name    = "containerd"
    version = "1.5.13"
  }
}

output "cluster_cert" {
  value = alicloud_cs_kubernetes.example.certificate_authority["cluster_cert"]
}

output "client_cert" {
  value = alicloud_cs_kubernetes.example.certificate_authority["client_cert"]
}

output "client_key" {
  value = alicloud_cs_kubernetes.example.certificate_authority["client_key"]
}

output "api_server_internet" {
  value = alicloud_cs_kubernetes.example.connections["api_server_internet"]
}

output "api_server_intranet" {
  value = alicloud_cs_kubernetes.example.connections["api_server_intranet"]
}

output "master_public_ip" {
  value = alicloud_cs_kubernetes.example.connections["master_public_ip"]
}

output "service_domain" {
  value = alicloud_cs_kubernetes.example.connections["service_domain"]
}
```

Updated configuration:

```terraform
resource "alicloud_cs_kubernetes" "example" {
  # ... other configuration ...

  runtime {
    name    = "containerd"
    version = "1.5.13"
  }
}

output "cluster_cert" {
  value = alicloud_cs_kubernetes.example.certificate_authority[0].cluster_cert
}

output "client_cert" {
  value = alicloud_cs_kubernetes.example.certificate_authority[0].client_cert
}

output "client_key" {
  value = alicloud_cs_kubernetes.example.certificate_authority[0].client_key
}

output "api_server_internet" {
  value = alicloud_cs_kubernetes.example.connections[0].api_server_internet
}

output "api_server_intranet" {
  value = alicloud_cs_kubernetes.example.connections[0].api_server_intranet
}

output "master_public_ip" {
  value = alicloud_cs_kubernetes.example.connections[0].master_public_ip
}

output "service_domain" {
  value = alicloud_cs_kubernetes.example.connections[0].service_domain
}
```

## Resource: alicloud_cs_managed_kubernetes

### certificate_authority, connections Type Change

The `certificate_authority` and `connections` attributes have changed from `TypeMap` to `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
resource "alicloud_cs_managed_kubernetes" "example" {
  # ... other configuration ...
}

output "cluster_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority["cluster_cert"]
}

output "client_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority["client_cert"]
}

output "client_key" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority["client_key"]
}

output "api_server_internet" {
  value = alicloud_cs_managed_kubernetes.example.connections["api_server_internet"]
}

output "api_server_intranet" {
  value = alicloud_cs_managed_kubernetes.example.connections["api_server_intranet"]
}

output "master_public_ip" {
  value = alicloud_cs_managed_kubernetes.example.connections["master_public_ip"]
}

output "service_domain" {
  value = alicloud_cs_managed_kubernetes.example.connections["service_domain"]
}
```

Updated configuration:

```terraform
resource "alicloud_cs_managed_kubernetes" "example" {
  # ... other configuration ...
}

output "cluster_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority[0].cluster_cert
}

output "client_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority[0].client_cert
}

output "client_key" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority[0].client_key
}

output "api_server_internet" {
  value = alicloud_cs_managed_kubernetes.example.connections[0].api_server_internet
}

output "api_server_intranet" {
  value = alicloud_cs_managed_kubernetes.example.connections[0].api_server_intranet
}

output "master_public_ip" {
  value = alicloud_cs_managed_kubernetes.example.connections[0].master_public_ip
}

output "service_domain" {
  value = alicloud_cs_managed_kubernetes.example.connections[0].service_domain
}
```

## Data Source: alicloud_cr_repos

### repos.domain_list Attribute Type Change

The `domain_list` attribute within `repos` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cr_repos" "example" {
  namespace = "my-namespace"
}

output "first_repo_public_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list["public"]
}

output "first_repo_internal_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list["internal"]
}

output "first_repo_vpc_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list["vpc"]
}
```

Updated configuration:

```terraform
data "alicloud_cr_repos" "example" {
  namespace = "my-namespace"
}

output "first_repo_public_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list[0].public
}

output "first_repo_internal_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list[0].internal
}

output "first_repo_vpc_domain" {
  value = data.alicloud_cr_repos.example.repos[0].domain_list[0].vpc
}
```

## Data Source: alicloud_cs_cluster_credential

### certificate_authority Attribute Type Change

The `certificate_authority` attribute has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cs_cluster_credential" "example" {
  cluster_id = "cluster-example-id"
}

output "cluster_cert" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority["cluster_cert"]
}

output "client_cert" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority["client_cert"]
}

output "client_key" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority["client_key"]
}
```

Updated configuration:

```terraform
data "alicloud_cs_cluster_credential" "example" {
  cluster_id = "cluster-example-id"
}

output "cluster_cert" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority[0].cluster_cert
}

output "client_cert" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority[0].client_cert
}

output "client_key" {
  value = data.alicloud_cs_cluster_credential.example.certificate_authority[0].client_key
}
```

## Data Source: alicloud_cs_edge_kubernetes_clusters

### clusters.connections Attribute Type Change

The `connections` attribute within `clusters` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cs_edge_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_edge_kubernetes_clusters.example.clusters[0].connections["api_server_internet"]
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_edge_kubernetes_clusters.example.clusters[0].connections["api_server_intranet"]
}
```

Updated configuration:

```terraform
data "alicloud_cs_edge_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_edge_kubernetes_clusters.example.clusters[0].connections[0].api_server_internet
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_edge_kubernetes_clusters.example.clusters[0].connections[0].api_server_intranet
}
```

## Data Source: alicloud_cs_kubernetes_clusters

### clusters.connections Attribute Type Change

The `connections` attribute within `clusters` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cs_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_kubernetes_clusters.example.clusters[0].connections["api_server_internet"]
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_kubernetes_clusters.example.clusters[0].connections["api_server_intranet"]
}
```

Updated configuration:

```terraform
data "alicloud_cs_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_kubernetes_clusters.example.clusters[0].connections[0].api_server_internet
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_kubernetes_clusters.example.clusters[0].connections[0].api_server_intranet
}
```

## Data Source: alicloud_cs_managed_kubernetes_clusters

### clusters.connections Attribute Type Change

The `connections` attribute within `clusters` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cs_managed_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_managed_kubernetes_clusters.example.clusters[0].connections["api_server_internet"]
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_managed_kubernetes_clusters.example.clusters[0].connections["api_server_intranet"]
}
```

Updated configuration:

```terraform
data "alicloud_cs_managed_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_managed_kubernetes_clusters.example.clusters[0].connections[0].api_server_internet
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_managed_kubernetes_clusters.example.clusters[0].connections[0].api_server_intranet
}
```

## Data Source: alicloud_cs_serverless_kubernetes_clusters

### clusters.connections Attribute Type Change

The `connections` attribute within `clusters` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_cs_serverless_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_serverless_kubernetes_clusters.example.clusters[0].connections["api_server_internet"]
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_serverless_kubernetes_clusters.example.clusters[0].connections["api_server_intranet"]
}
```

Updated configuration:

```terraform
data "alicloud_cs_serverless_kubernetes_clusters" "example" {
  name_regex = "my-cluster"
}

output "first_cluster_api_server_internet" {
  value = data.alicloud_cs_serverless_kubernetes_clusters.example.clusters[0].connections[0].api_server_internet
}

output "first_cluster_api_server_intranet" {
  value = data.alicloud_cs_serverless_kubernetes_clusters.example.clusters[0].connections[0].api_server_intranet
}
```

## Data Source: alicloud_db_instance_classes

### instance_classes.storage_range Attribute Type Change

The `storage_range` attribute within `instance_classes` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_db_instance_classes" "example" {
  engine         = "MySQL"
  engine_version = "8.0"
}

output "first_class_storage_min" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range["min"]
}

output "first_class_storage_max" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range["max"]
}

output "first_class_storage_step" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range["step"]
}
```

Updated configuration:

```terraform
data "alicloud_db_instance_classes" "example" {
  engine         = "MySQL"
  engine_version = "8.0"
}

output "first_class_storage_min" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range[0].min
}

output "first_class_storage_max" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range[0].max
}

output "first_class_storage_step" {
  value = data.alicloud_db_instance_classes.example.instance_classes[0].storage_range[0].step
}
```

## Data Source: alicloud_instance_types

### instance_types.gpu Attribute Type Change

The `gpu` attribute within `instance_types` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_gpu_amount" {
  value = data.alicloud_instance_types.example.instance_types[0].gpu["amount"]
}

output "first_type_gpu_category" {
  value = data.alicloud_instance_types.example.instance_types[0].gpu["category"]
}
```

Updated configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_gpu_amount" {
  value = data.alicloud_instance_types.example.instance_types[0].gpu[0].amount
}

output "first_type_gpu_category" {
  value = data.alicloud_instance_types.example.instance_types[0].gpu[0].category
}
```

### instance_types.burstable_instance Attribute Type Change

The `burstable_instance` attribute within `instance_types` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_initial_credit" {
  value = data.alicloud_instance_types.example.instance_types[0].burstable_instance["initial_credit"]
}

output "first_type_baseline_credit" {
  value = data.alicloud_instance_types.example.instance_types[0].burstable_instance["baseline_credit"]
}
```

Updated configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_initial_credit" {
  value = data.alicloud_instance_types.example.instance_types[0].burstable_instance[0].initial_credit
}

output "first_type_baseline_credit" {
  value = data.alicloud_instance_types.example.instance_types[0].burstable_instance[0].baseline_credit
}
```

### instance_types.local_storage Attribute Type Change

The `local_storage` attribute within `instance_types` has changed from a `TypeMap` to a `TypeList`. Update attribute references from map key syntax to list index syntax.

Previous configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_local_storage_capacity" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage["capacity"]
}

output "first_type_local_storage_amount" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage["amount"]
}

output "first_type_local_storage_category" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage["category"]
}
```

Updated configuration:

```terraform
data "alicloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 4
}

output "first_type_local_storage_capacity" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage[0].capacity
}

output "first_type_local_storage_amount" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage[0].amount
}

output "first_type_local_storage_category" {
  value = data.alicloud_instance_types.example.instance_types[0].local_storage[0].category
}
```
