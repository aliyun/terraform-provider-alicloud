package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudEhpcUsersDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ehpc_users.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEhpcUsersConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"ids":        []string{"${alicloud_ehpc_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"ids":        []string{"${alicloud_ehpc_user.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"name_regex": "${alicloud_ehpc_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"name_regex": "${alicloud_ehpc_user.default.user_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"ids":        []string{"${alicloud_ehpc_user.default.id}"},
			"name_regex": "${alicloud_ehpc_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_ehpc_user.default.cluster_id}",
			"ids":        []string{"${alicloud_ehpc_user.default.id}_fake"},
			"name_regex": "${alicloud_ehpc_user.default.user_name}_fake",
		}),
	}

	var existAliCloudEhpcUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"names.#":            "1",
			"users.#":            "1",
			"users.0.id":         CHECKSET,
			"users.0.cluster_id": CHECKSET,
			"users.0.user_name":  CHECKSET,
			"users.0.user_id":    CHECKSET,
			"users.0.group":      CHECKSET,
			"users.0.group_id":   CHECKSET,
		}
	}

	var fakeAliCloudEhpcUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"users.#": "0",
		}
	}

	var aliCloudEhpcUsersInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ehpc_users.default",
		existMapFunc: existAliCloudEhpcUsersMapFunc,
		fakeMapFunc:  fakeAliCloudEhpcUsersMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEhpcUsersInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceEhpcUsersConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "user_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "user_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.user_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_nas_file_system" "user_test_nas" {
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_nas_access_group" "user_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = var.name
  file_system_type  = "standard"
}

resource "alicloud_security_group" "user_test_security_group" {
  vpc_id              = alicloud_vpc.user_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_mount_target" "user_test_mount_domain" {
  vpc_id            = alicloud_vpc.user_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.user_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.user_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.user_test_nas.id
}

resource "alicloud_nas_access_rule" "user_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.user_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.user_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}

resource "alicloud_ehpc_cluster_v2" "user_cluster_test" {
  depends_on = [alicloud_nas_access_rule.user_test_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id    = alicloud_vpc.user_test_vpc.id
  cluster_category  = "Standard"
  cluster_mode      = "Integrated"
  security_group_id = alicloud_security_group.user_test_security_group.id
  addons {
    version        = "1.0"
    services_spec  = <<EOF
      [
        {
          "ServiceName": "SSH",
          "NetworkACL": [
            {
              "Port": 22,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        },
        {
          "ServiceName": "VNC",
          "NetworkACL": [
            {
              "Port": 12016,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        },
        {
          "ServiceName": "CLIENT",
          "ServiceAccessType": "URL",
          "ServiceAccessUrl": "https://ehpc-app.oss-cn-hangzhou.aliyuncs.com/ClientRelease/E-HPC-Client-Mac-zh-cn.zip",
          "NetworkACL": [
            {
              "Port": 12011,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        }
      ]
  EOF
    resources_spec = <<EOF
      {
        "EipResource": {
          "AutoCreate": true
        },
        "EcsResources": [
          {
            "ImageId": "centos_7_6_x64_20G_alibase_20211130.vhd",
            "EnableHT": true,
            "InstanceChargeType": "PostPaid",
            "InstanceType": "ecs.c7.xlarge",
            "SpotStrategy": "NoSpot",
            "SystemDisk": {
              "Category": "cloud_essd",
              "Size": 40,
              "Level": "PL0"
            },
            "DataDisks": [
              {
                "Category": "cloud_essd",
                "Size": 40,
                "Level": "PL0"
              }
            ]
          }
        ]
      }
  EOF
    name           = "Login"
  }
  cluster_name        = "minimal-test-cluster"
  deletion_protection = false
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.user_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.user_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.user_test_vswitch.id
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
  user_name  = var.name
  cluster_id = alicloud_ehpc_cluster_v2.user_cluster_test.id
  password   = "YourPassword1234!"
}
`, name)
}
