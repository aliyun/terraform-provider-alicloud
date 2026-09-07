---
subcategory: "Elastic High Performance Computing (Ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_user"
description: |-
  Provides a Alicloud Ehpc User resource.
---

# alicloud_ehpc_user

Provides a Ehpc User resource.

A user in an E-HPC cluster.

For information about Ehpc User and how to use it, see [What is User](https://next.api.alibabacloud.com/document/EHPC/2024-07-30/CreateUsers).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "user_password_example_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "user-password-example-vpc"
}

resource "alicloud_vswitch" "user_password_example_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.user_password_example_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "user-password-example-vsw"
}

resource "alicloud_nas_file_system" "user_password_example_nas" {
  description  = "user-password-example-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
  recycle_bin {
    status        = "Disable"
    reserved_days = "7"
  }
}

resource "alicloud_nas_access_group" "user_password_example_access_group" {
  access_group_type = "Vpc"
  description       = "user password example access group"
  access_group_name = "UserPasswordTest"
  file_system_type  = "standard"
}

resource "alicloud_nas_mount_target" "user_password_example_mount_domain" {
  vpc_id            = alicloud_vpc.user_password_example_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.user_password_example_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.user_password_example_vswitch.id
  file_system_id    = alicloud_nas_file_system.user_password_example_nas.id
}

resource "alicloud_security_group" "user_password_example_security_group" {
  vpc_id              = alicloud_vpc.user_password_example_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_access_rule" "user_password_example_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.user_password_example_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.user_password_example_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/8"
}

resource "alicloud_ehpc_cluster_v2" "user_password_example_cluster" {
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id      = alicloud_vpc.user_password_example_vpc.id
  cluster_category    = "Standard"
  cluster_mode        = "Integrated"
  security_group_id   = alicloud_security_group.user_password_example_security_group.id
  cluster_name        = "user-password-example-cluster"
  deletion_protection = false
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_password_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_password_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_password_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_password_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_password_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_password_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.user_password_example_vswitch.id
  manager {
    manager_node {
      system_disk {
        category = "cloud_essd"
        size     = "40"
        level    = "PL0"
      }
      enable_ht            = true
      instance_charge_type = "PostPaid"
      image_id             = "centos_7_6_x64_20G_alibase_20211130.vhd"
      instance_type        = "ecs.c6.xlarge"
      spot_strategy        = "NoSpot"
    }
    scheduler {
      type    = "SLURM"
      version = "22.05.8"
    }
    dns {
      type    = "nis"
      version = "1.0"
    }
    directory_service {
      type    = "nis"
      version = "1.0"
    }
  }
}


resource "alicloud_ehpc_user" "default" {
  group      = "users"
  user_name  = "passworduser1"
  cluster_id = alicloud_ehpc_cluster_v2.user_password_example_cluster.id
  password   = "UserPass123"
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Optional) The SSH public-key content. This field is accepted only for create and is never returned. This parameter is mutually exclusive with the `password` parameter. This parameter takes effect when the cluster authentication method is set to key (not recommended).

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `cluster_id` - (Required, ForceNew) The cluster ID.
* `group` - (Required) The user group. Valid values: `users`, `wheel`.
* `password` - (Optional) The user password. This parameter is mutually exclusive with the `key_pair_name` parameter. This parameter takes effect when the cluster authentication method is set to password (recommended).

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `user_name` - (Required, ForceNew) The username.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User. It formats as `<cluster_id>:<user_name>`.
* `group_id` - The user group ID.
* `user_id` - The user ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the User.
* `update` - (Defaults to 5 mins) Used when update the User.
* `delete` - (Defaults to 5 mins) Used when delete the User.

## Import

Ehpc User can be imported using the id, e.g.

```shell
$ terraform import alicloud_ehpc_user.example <cluster_id>:<user_name>
```
