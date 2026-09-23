package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ehpc Queue. >>> Resource test cases, automatically generated.
// Case Queue_minimal_test 12111
func TestAccAliCloudEhpcQueue_basic12111(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcQueueMap12111)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcQueueBasicDependence12111)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_ehpc_cluster_v2.queue_minimal_cluster_test.id}",
					"queue_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": CHECKSET,
						"queue_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEhpcQueueMap12111 = map[string]string{}

func AlicloudEhpcQueueBasicDependence12111(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "queue_minimal_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_vswitch" "queue_minimal_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_minimal_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_file_system" "queue_minimal_test_nas" {
  description  = "test-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_nas_access_group" "queue_minimal_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_security_group" "queue_minimal_test_security_group" {
  vpc_id              = alicloud_vpc.queue_minimal_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_mount_target" "queue_minimal_test_mount_domain" {
  vpc_id            = alicloud_vpc.queue_minimal_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.queue_minimal_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.queue_minimal_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.queue_minimal_test_nas.id
}

resource "alicloud_nas_access_rule" "queue_minimal_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.queue_minimal_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.queue_minimal_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}

resource "alicloud_ehpc_cluster_v2" "queue_minimal_cluster_test" {
  depends_on = [alicloud_nas_access_rule.queue_minimal_test_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id    = alicloud_vpc.queue_minimal_test_vpc.id
  cluster_category  = "Standard"
  cluster_mode      = "Integrated"
  security_group_id = alicloud_security_group.queue_minimal_test_security_group.id
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
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.queue_minimal_test_vswitch.id
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

// Case Queue_autoscale_test 12112
func TestAccAliCloudEhpcQueue_basic12112(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcQueueMap12112)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcQueueBasicDependence12112)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "true",
					"cluster_id":      "${alicloud_ehpc_cluster_v2.queue_autoscale_cluster_test.id}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_autoscale_test_vswitch.id}"},
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud_essd",
									"size":     "40",
									"level":    "PL0",
								},
							},
							"instance_charge_type": "PostPaid",
							"image_id":             "centos_7_6_x64_20G_alibase_20211130.vhd",
							"spot_price_limit":     "0",
							"duration":             "0",
							"instance_type":        "ecs.c6.xlarge",
							"spot_strategy":        "NoSpot",
						},
					},
					"inter_connect":    "vpc",
					"enable_scale_out": "true",
					"max_count":        "10",
					"queue_name":       name,
					"initial_count":    "0",
					"hostname_prefix":  "compute",
					"min_count":        "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "true",
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"compute_nodes.#":  "1",
						"inter_connect":    "vpc",
						"enable_scale_out": "true",
						"max_count":        "10",
						"queue_name":       name,
						"initial_count":    "0",
						"hostname_prefix":  "compute",
						"min_count":        "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "false",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_autoscale_test_vswitch_modified.id}"},
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud",
									"size":     "50",
									"level":    "PL1",
								},
							},
							"enable_ht":            "true",
							"instance_charge_type": "PostPaid",
							"image_id":             "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
							"spot_price_limit":     "3",
							"duration":             "1",
							"instance_type":        "ecs.hpc8ae.32xlarge",
							"spot_strategy":        "SpotWithPriceLimit",
						},
					},
					"inter_connect":    "eRDMA",
					"enable_scale_out": "false",
					"max_count":        "8",
					"hostname_prefix":  "test",
					"min_count":        "1",
					"hostname_suffix":  "hpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "false",
						"vswitch_ids.#":    "1",
						"compute_nodes.#":  "1",
						"inter_connect":    "eRDMA",
						"enable_scale_out": "false",
						"max_count":        "8",
						"hostname_prefix":  "test",
						"min_count":        "1",
						"hostname_suffix":  "hpc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEhpcQueueMap12112 = map[string]string{}

func AlicloudEhpcQueueBasicDependence12112(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "queue_autoscale_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_access_group" "queue_autoscale_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "queue_autoscale_test_nas" {
  description  = "test-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_vswitch" "queue_autoscale_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_autoscale_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_mount_target" "queue_autoscale_test_mount_domain" {
  vpc_id            = alicloud_vpc.queue_autoscale_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.queue_autoscale_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.queue_autoscale_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.queue_autoscale_test_nas.id
}

resource "alicloud_security_group" "queue_autoscale_test_security_group" {
  vpc_id              = alicloud_vpc.queue_autoscale_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_vswitch" "queue_autoscale_test_vswitch_modified" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_autoscale_test_vpc.id
  zone_id      = "cn-hangzhou-h"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_ehpc_cluster_v2" "queue_autoscale_cluster_test" {
  depends_on = [alicloud_nas_access_rule.queue_autoscale_test_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id      = alicloud_vpc.queue_autoscale_test_vpc.id
  cluster_category    = "Standard"
  cluster_mode        = "Integrated"
  security_group_id   = alicloud_security_group.queue_autoscale_test_security_group.id
  cluster_name        = "autoscale-test-cluster"
  deletion_protection = false
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_autoscale_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_autoscale_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_autoscale_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_autoscale_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_autoscale_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_autoscale_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.queue_autoscale_test_vswitch.id
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

resource "alicloud_nas_access_rule" "queue_autoscale_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.queue_autoscale_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.queue_autoscale_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/8"
}


`, name)
}

// Case queue_prepaid_test 12113
func TestAccAliCloudEhpcQueue_basic12113(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcQueueMap12113)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcQueueBasicDependence12113)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "true",
					"cluster_id":      "${alicloud_ehpc_cluster_v2.queue_prepaid_cluster_test.id}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_prepaid_test_vswitch.id}"},
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud_essd",
									"size":     "40",
									"level":    "PL0",
								},
							},
							"auto_renew_period":    "1",
							"instance_charge_type": "PrePaid",
							"auto_renew":           "false",
							"image_id":             "centos_7_6_x64_20G_alibase_20211130.vhd",
							"period":               "1",
							"instance_type":        "ecs.c6.xlarge",
							"period_unit":          "Week",
						},
					},
					"inter_connect":    "vpc",
					"enable_scale_out": "true",
					"max_count":        "10",
					"queue_name":       name,
					"initial_count":    "0",
					"hostname_prefix":  "compute",
					"min_count":        "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "true",
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "1",
						"compute_nodes.#":  "1",
						"inter_connect":    "vpc",
						"enable_scale_out": "true",
						"max_count":        "10",
						"queue_name":       name,
						"initial_count":    "0",
						"hostname_prefix":  "compute",
						"min_count":        "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "false",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_prepaid_test_vswitch_modified.id}"},
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud_essd",
									"size":     "40",
									"level":    "PL0",
								},
							},
							"enable_ht":            "true",
							"auto_renew_period":    "2",
							"instance_charge_type": "PrePaid",
							"auto_renew":           "true",
							"image_id":             "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
							"period":               "2",
							"instance_type":        "ecs.hpc8ae.32xlarge",
							"period_unit":          "Month",
						},
					},
					"inter_connect":    "eRDMA",
					"enable_scale_out": "false",
					"max_count":        "8",
					"hostname_prefix":  "test",
					"min_count":        "1",
					"hostname_suffix":  "hpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "false",
						"vswitch_ids.#":    "1",
						"compute_nodes.#":  "1",
						"inter_connect":    "eRDMA",
						"enable_scale_out": "false",
						"max_count":        "8",
						"hostname_prefix":  "test",
						"min_count":        "1",
						"hostname_suffix":  "hpc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "false",
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud",
									"size":     "50",
									"level":    "PL1",
								},
							},
							"enable_ht":            "true",
							"instance_charge_type": "PostPaid",
							"image_id":             "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
							"spot_price_limit":     "3",
							"duration":             "1",
							"instance_type":        "ecs.hpc8ae.32xlarge",
							"spot_strategy":        "SpotWithPriceLimit",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEhpcQueueMap12113 = map[string]string{}

func AlicloudEhpcQueueBasicDependence12113(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "queue_prepaid_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_access_group" "queue_prepaid_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "queue_prepaid_test_nas" {
  description  = "test-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_vswitch" "queue_prepaid_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_prepaid_test_vpc.id
  zone_id      = "cn-hangzhou-h"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_security_group" "queue_prepaid_test_security_group" {
  vpc_id              = alicloud_vpc.queue_prepaid_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_mount_target" "queue_prepaid_test_mount_domain" {
  vpc_id            = alicloud_vpc.queue_prepaid_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.queue_prepaid_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.queue_prepaid_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.queue_prepaid_test_nas.id
}

resource "alicloud_vswitch" "queue_prepaid_test_vswitch_modified" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_prepaid_test_vpc.id
  zone_id      = "cn-hangzhou-h"
  cidr_block   = "10.0.2.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_ehpc_cluster_v2" "queue_prepaid_cluster_test" {
  depends_on = [alicloud_nas_access_rule.queue_prepaid_test_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id      = alicloud_vpc.queue_prepaid_test_vpc.id
  cluster_category    = "Standard"
  cluster_mode        = "Integrated"
  security_group_id   = alicloud_security_group.queue_prepaid_test_security_group.id
  cluster_name        = "autoscale-test-cluster"
  deletion_protection = false
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_prepaid_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_prepaid_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_prepaid_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_prepaid_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_prepaid_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_prepaid_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.queue_prepaid_test_vswitch.id
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

resource "alicloud_nas_access_rule" "queue_prepaid_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.queue_prepaid_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.queue_prepaid_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/8"
}


`, name)
}

// Case queue_vswitchIds_test 12114
func TestAccAliCloudEhpcQueue_basic12114(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcQueueMap12114)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcQueueBasicDependence12114)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "true",
					"cluster_id":      "${alicloud_ehpc_cluster_v2.queue_vswitchids_test_cluster.id}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_vswitchIds_test_vswitch1.id}", "${alicloud_vswitch.queue_vswitchIds_test_vswitch2.id}", "${alicloud_vswitch.queue_vswitchIds_test_vswitch3.id}"},
					"hostname_suffix": "node",
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud_essd",
									"size":     "40",
									"level":    "PL0",
								},
							},
							"enable_ht":            "true",
							"instance_charge_type": "PostPaid",
							"image_id":             "centos_7_6_x64_20G_alibase_20211130.vhd",
							"duration":             "1",
							"instance_type":        "ecs.c6.xlarge",
							"spot_strategy":        "NoSpot",
						},
					},
					"inter_connect":    "vpc",
					"enable_scale_out": "true",
					"max_count":        "10",
					"queue_name":       name,
					"initial_count":    "0",
					"min_count":        "0",
					"hostname_prefix":  "compute",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "true",
						"cluster_id":       CHECKSET,
						"vswitch_ids.#":    "3",
						"hostname_suffix":  "node",
						"compute_nodes.#":  "1",
						"inter_connect":    "vpc",
						"enable_scale_out": "true",
						"max_count":        "10",
						"queue_name":       name,
						"initial_count":    "0",
						"min_count":        "0",
						"hostname_prefix":  "compute",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_scale_in": "false",
					"vswitch_ids": []string{
						"${alicloud_vswitch.queue_vswitchIds_test_vswitch1.id}", "${alicloud_vswitch.queue_vswitchIds_test_vswitch2.id}"},
					"hostname_suffix": "hpc",
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud",
									"size":     "50",
									"level":    "PL1",
								},
							},
							"enable_ht":            "true",
							"instance_charge_type": "PostPaid",
							"image_id":             "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
							"spot_price_limit":     "3",
							"duration":             "1",
							"instance_type":        "ecs.hpc8ae.32xlarge",
							"spot_strategy":        "SpotWithPriceLimit",
						},
					},
					"inter_connect":    "eRDMA",
					"enable_scale_out": "false",
					"max_count":        "8",
					"min_count":        "1",
					"hostname_prefix":  "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_scale_in":  "false",
						"vswitch_ids.#":    "2",
						"hostname_suffix":  "hpc",
						"compute_nodes.#":  "1",
						"inter_connect":    "eRDMA",
						"enable_scale_out": "false",
						"max_count":        "8",
						"min_count":        "1",
						"hostname_prefix":  "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{},
					"compute_nodes": []map[string]interface{}{
						{
							"system_disk": []map[string]interface{}{
								{
									"category": "cloud",
									"size":     "50",
									"level":    "PL1",
								},
							},
							"enable_ht":            "true",
							"instance_charge_type": "PostPaid",
							"image_id":             "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
							"spot_price_limit":     "3",
							"duration":             "1",
							"instance_type":        "ecs.hpc8ae.32xlarge",
							"spot_strategy":        "SpotWithPriceLimit",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#":   "0",
						"compute_nodes.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEhpcQueueMap12114 = map[string]string{}

func AlicloudEhpcQueueBasicDependence12114(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "queue_vswitchIds_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_file_system" "queue_vswitchIds_test_nas" {
  description  = "test-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_vswitch" "queue_vswitchIds_test_vswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_vswitchIds_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = "test-cluster-vsw1"
}

resource "alicloud_nas_access_group" "queue_vswitchIds_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_security_group" "queue_vswitchIds_test_security_group" {
  vpc_id              = alicloud_vpc.queue_vswitchIds_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_mount_target" "queue_vswitchIds_test_mount_domain" {
  vpc_id            = alicloud_vpc.queue_vswitchIds_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.queue_vswitchIds_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.queue_vswitchIds_test_vswitch1.id
  file_system_id    = alicloud_nas_file_system.queue_vswitchIds_test_nas.id
}

resource "alicloud_nas_access_rule" "queue_vswitchIds_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.queue_vswitchIds_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.queue_vswitchIds_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/8"
}

resource "alicloud_ehpc_cluster_v2" "queue_vswitchids_test_cluster" {
  depends_on = [alicloud_nas_access_rule.queue_vswitchIds_test_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id    = alicloud_vpc.queue_vswitchIds_test_vpc.id
  cluster_category  = "Standard"
  cluster_mode      = "Integrated"
  security_group_id = alicloud_security_group.queue_vswitchIds_test_security_group.id
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
    mount_target_domain = alicloud_nas_mount_target.queue_vswitchIds_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_vswitchIds_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_vswitchIds_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_vswitchIds_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_vswitchIds_test_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_vswitchIds_test_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.queue_vswitchIds_test_vswitch1.id
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

resource "alicloud_vswitch" "queue_vswitchIds_test_vswitch2" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_vswitchIds_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.3.0/24"
  vswitch_name = "test-cluster-vsw2"
}

resource "alicloud_vswitch" "queue_vswitchIds_test_vswitch3" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_vswitchIds_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.2.0/24"
  vswitch_name = "test-cluster-vsw3"
}


`, name)
}

// Test Ehpc Queue. <<< Resource test cases, automatically generated.
