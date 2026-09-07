// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ehpc User. >>> Resource test cases, automatically generated.
// Case User_password_test 13038
func TestAccAliCloudEhpcUser_basic13038(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_user.default"
	ra := resourceAttrInit(resourceId, AliCloudEhpcUserMap13038)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEhpcUserBasicDependence13038)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group":      "users",
					"user_name":  name,
					"cluster_id": "${alicloud_ehpc_cluster_v2.user_cluster_test.id}",
					"password":   "YourPassword1234!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group":      "users",
						"user_name":  name,
						"cluster_id": CHECKSET,
						"password":   "YourPassword1234!",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group": "wheel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group": "wheel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group": "users",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group": "users",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword1234!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "YourPassword1234!update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_pair_name", "password"},
			},
		},
	})
}

func TestAccAliCloudEhpcUser_basic13038_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_user.default"
	ra := resourceAttrInit(resourceId, AliCloudEhpcUserMap13038)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEhpcUserBasicDependence13038)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group":      "users",
					"user_name":  name,
					"cluster_id": "${alicloud_ehpc_cluster_v2.user_cluster_test.id}",
					"password":   "YourPassword1234!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group":      "users",
						"user_name":  name,
						"cluster_id": CHECKSET,
						"password":   "YourPassword1234!",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_pair_name", "password"},
			},
		},
	})
}

var AliCloudEhpcUserMap13038 = map[string]string{
	"user_id":  CHECKSET,
	"group_id": CHECKSET,
}

func AliCloudEhpcUserBasicDependence13038(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = var.name
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
  description  = var.name
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

resource "alicloud_ehpc_cluster_v2" "user_cluster_test_key_pair_name" {
  depends_on = [alicloud_nas_access_rule.user_test_access_rule]
  cluster_credentials {
    key_pair_name = alicloud_ecs_key_pair.default.id
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
`, name)
}

// Case User_key_pair_name_test 13039
func TestAccAliCloudEhpcUser_basic13039(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_user.default"
	ra := resourceAttrInit(resourceId, AliCloudEhpcUserMap13038)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEhpcUserBasicDependence13038)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group":         "users",
					"user_name":     name,
					"cluster_id":    "${alicloud_ehpc_cluster_v2.user_cluster_test_key_pair_name.id}",
					"key_pair_name": "${alicloud_ecs_key_pair.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group":         "users",
						"user_name":     name,
						"cluster_id":    CHECKSET,
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group": "wheel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group": "wheel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group": "users",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group": "users",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_pair_name", "password"},
			},
		},
	})
}

func TestAccAliCloudEhpcUser_basic13039_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_user.default"
	ra := resourceAttrInit(resourceId, AliCloudEhpcUserMap13038)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEhpcUserBasicDependence13038)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group":         "users",
					"user_name":     name,
					"cluster_id":    "${alicloud_ehpc_cluster_v2.user_cluster_test_key_pair_name.id}",
					"key_pair_name": "${alicloud_ecs_key_pair.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group":         "users",
						"user_name":     name,
						"cluster_id":    CHECKSET,
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_pair_name", "password"},
			},
		},
	})
}
