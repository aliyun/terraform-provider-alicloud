// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ehpc ClusterV2. >>> Resource test cases, automatically generated.
// Case minimal_cluster_test_clientVersion 12089
func TestAccAliCloudEhpcClusterV2_basic12089(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterV2Map12089)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcClusterV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterV2BasicDependence12089)
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
					"cluster_credentials": []map[string]interface{}{
						{
							"password": "aliHPC123",
						},
					},
					"cluster_vpc_id":      "${alicloud_vpc.minimal_test_clientVersion_vpc.id}",
					"cluster_category":    "Standard",
					"cluster_mode":        "Integrated",
					"security_group_id":   "${alicloud_security_group.minimal_test_clientVersion_security_group.id}",
					"cluster_name":        "minimal-test-cluster",
					"deletion_protection": "false",
					"client_version":      "2.0.47",
					"shared_storages": []map[string]interface{}{
						{
							"mount_directory":     "/home",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_clientVersion_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_clientVersion_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/opt",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_clientVersion_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_clientVersion_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/ehpcdata",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_clientVersion_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_clientVersion_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
					},
					"cluster_vswitch_id": "${alicloud_vswitch.minimal_test_clientVersion_vswitch.id}",
					"manager": []map[string]interface{}{
						{
							"manager_node": []map[string]interface{}{
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
									"instance_type":        "ecs.c6.xlarge",
									"spot_strategy":        "NoSpot",
								},
							},
							"scheduler": []map[string]interface{}{
								{
									"type":    "SLURM",
									"version": "22.05.8",
								},
							},
							"dns": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
							"directory_service": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_vpc_id":      CHECKSET,
						"cluster_category":    "Standard",
						"cluster_mode":        "Integrated",
						"security_group_id":   CHECKSET,
						"cluster_name":        "minimal-test-cluster",
						"deletion_protection": "false",
						"client_version":      "2.0.47",
						"shared_storages.#":   "3",
						"cluster_vswitch_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_version": "2.0.48",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_version": "2.0.48",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addons", "cluster_credentials"},
			},
		},
	})
}

var AlicloudEhpcClusterV2Map12089 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEhpcClusterV2BasicDependence12089(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "minimal_test_clientVersion_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_access_group" "minimal_test_clientVersion_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_vswitch" "minimal_test_clientVersion_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.minimal_test_clientVersion_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_file_system" "minimal_test_clientVersion_nas" {
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

resource "alicloud_nas_mount_target" "minimal_test_clientVersion_mount_domain" {
  vpc_id            = alicloud_vpc.minimal_test_clientVersion_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.minimal_test_clientVersion_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.minimal_test_clientVersion_vswitch.id
  file_system_id    = alicloud_nas_file_system.minimal_test_clientVersion_nas.id
}

resource "alicloud_security_group" "minimal_test_clientVersion_security_group" {
  vpc_id              = alicloud_vpc.minimal_test_clientVersion_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_access_rule" "minimal_test_clientVersion_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.minimal_test_clientVersion_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.minimal_test_clientVersion_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}


`, name)
}

// Case minimal_cluster_test_prepaid 12090
func TestAccAliCloudEhpcClusterV2_basic12090(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterV2Map12090)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcClusterV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterV2BasicDependence12090)
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
					"cluster_credentials": []map[string]interface{}{
						{
							"password": "aliHPC123",
						},
					},
					"cluster_vpc_id":      "${alicloud_vpc.minimal_test_prepaid_vpc.id}",
					"cluster_category":    "Standard",
					"cluster_mode":        "Integrated",
					"security_group_id":   "${alicloud_security_group.minimal_test_prepaid_security_group.id}",
					"cluster_name":        "minimal-test-cluster",
					"deletion_protection": "true",
					"shared_storages": []map[string]interface{}{
						{
							"mount_directory":     "/home",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_prepaid_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_prepaid_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/opt",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_prepaid_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_prepaid_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/ehpcdata",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_prepaid_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_prepaid_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
					},
					"cluster_vswitch_id": "${alicloud_vswitch.minimal_test_prepaid_vswitch.id}",
					"manager": []map[string]interface{}{
						{
							"manager_node": []map[string]interface{}{
								{
									"system_disk": []map[string]interface{}{
										{
											"category": "cloud_essd",
											"size":     "40",
											"level":    "PL0",
										},
									},
									"auto_renew_period":    "1",
									"enable_ht":            "true",
									"instance_charge_type": "PrePaid",
									"auto_renew":           "true",
									"image_id":             "centos_7_6_x64_20G_alibase_20211130.vhd",
									"period":               "1",
									"instance_type":        "ecs.c6.xlarge",
									"period_unit":          "Month",
								},
							},
							"scheduler": []map[string]interface{}{
								{
									"type":    "SLURM",
									"version": "22.05.8",
								},
							},
							"dns": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
							"directory_service": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_vpc_id":      CHECKSET,
						"cluster_category":    "Standard",
						"cluster_mode":        "Integrated",
						"security_group_id":   CHECKSET,
						"cluster_name":        "minimal-test-cluster",
						"deletion_protection": "true",
						"shared_storages.#":   "3",
						"cluster_vswitch_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addons", "cluster_credentials"},
			},
		},
	})
}

var AlicloudEhpcClusterV2Map12090 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEhpcClusterV2BasicDependence12090(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "minimal_test_prepaid_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_file_system" "minimal_test_prepaid_nas" {
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

resource "alicloud_vswitch" "minimal_test_prepaid_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.minimal_test_prepaid_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_access_group" "minimal_test_prepaid_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_nas_mount_target" "minimal_test_prepaid_mount_domain" {
  vpc_id            = alicloud_vpc.minimal_test_prepaid_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.minimal_test_prepaid_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.minimal_test_prepaid_vswitch.id
  file_system_id    = alicloud_nas_file_system.minimal_test_prepaid_nas.id
}

resource "alicloud_security_group" "minimal_test_prepaid_security_group" {
  vpc_id              = alicloud_vpc.minimal_test_prepaid_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_access_rule" "minimal_test_prepaid_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.minimal_test_prepaid_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.minimal_test_prepaid_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}


`, name)
}

// Case full_cluster_test 12035
func TestAccAliCloudEhpcClusterV2_basic12035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterV2Map12035)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcClusterV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterV2BasicDependence12035)
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
					"cluster_credentials": []map[string]interface{}{
						{
							"password": "aliHPC123",
						},
					},
					"cluster_vpc_id":    "${alicloud_vpc.full_test_vpc.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"cluster_category":  "Standard",
					"cluster_mode":      "Integrated",
					"security_group_id": "${alicloud_security_group.full_test_security_group.id}",
					"addons": []map[string]interface{}{
						{
							"version":        "1.0",
							"services_spec":  "[\\n        {\\n          \\\"ServiceName\\\": \\\"SSH\\\",\\n          \\\"NetworkACL\\\": [\\n            {\\n              \\\"Port\\\": 22,\\n              \\\"SourceCidrIp\\\": \\\"0.0.0.0/0\\\",\\n              \\\"IpProtocol\\\": \\\"TCP\\\"\\n            }\\n          ]\\n        },\\n        {\\n          \\\"ServiceName\\\": \\\"VNC\\\",\\n          \\\"NetworkACL\\\": [\\n            {\\n              \\\"Port\\\": 12016,\\n              \\\"SourceCidrIp\\\": \\\"0.0.0.0/0\\\",\\n              \\\"IpProtocol\\\": \\\"TCP\\\"\\n            }\\n          ]\\n        },\\n        {\\n          \\\"ServiceName\\\": \\\"CLIENT\\\",\\n          \\\"ServiceAccessType\\\": \\\"URL\\\",\\n          \\\"ServiceAccessUrl\\\": \\\"https://ehpc-app.oss-cn-hangzhou.aliyuncs.com/ClientRelease/E-HPC-Client-Mac-zh-cn.zip\\\",\\n          \\\"NetworkACL\\\": [\\n            {\\n              \\\"Port\\\": 12011,\\n              \\\"SourceCidrIp\\\": \\\"0.0.0.0/0\\\",\\n              \\\"IpProtocol\\\": \\\"TCP\\\"\\n            }\\n          ]\\n        }\\n      ]",
							"resources_spec": "{\\n        \\\"EipResource\\\": {\\n          \\\"AutoCreate\\\": true\\n        },\\n        \\\"EcsResources\\\": [\\n          {\\n            \\\"ImageId\\\": \\\"centos_7_6_x64_20G_alibase_20211130.vhd\\\",\\n            \\\"EnableHT\\\": true,\\n            \\\"InstanceChargeType\\\": \\\"PostPaid\\\",\\n            \\\"InstanceType\\\": \\\"ecs.c7.xlarge\\\",\\n            \\\"SpotStrategy\\\": \\\"NoSpot\\\",\\n            \\\"SystemDisk\\\": {\\n              \\\"Category\\\": \\\"cloud_essd\\\",\\n              \\\"Size\\\": 40,\\n              \\\"Level\\\": \\\"PL0\\\"\\n            },\\n            \\\"DataDisks\\\": [\\n              {\\n                \\\"Category\\\": \\\"cloud_essd\\\",\\n                \\\"Size\\\": 40,\\n                \\\"Level\\\": \\\"PL0\\\"\\n              }\\n            ]\\n          }\\n        ]\\n      }",
							"name":           "Login",
						},
					},
					"cluster_name":        "full-test-cluster",
					"deletion_protection": "true",
					"shared_storages": []map[string]interface{}{
						{
							"mount_directory":     "/home",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.full_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.full_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/opt",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.full_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.full_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/ehpcdata",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.full_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.full_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
					},
					"cluster_vswitch_id": "${alicloud_vswitch.full_test_vswitch.id}",
					"manager": []map[string]interface{}{
						{
							"manager_node": []map[string]interface{}{
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
									"duration":             "0",
									"instance_type":        "ecs.c6.xlarge",
									"spot_strategy":        "NoSpot",
								},
							},
							"scheduler": []map[string]interface{}{
								{
									"type":    "SLURM",
									"version": "22.05.8",
								},
							},
							"dns": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
							"directory_service": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_vpc_id":      CHECKSET,
						"resource_group_id":   CHECKSET,
						"cluster_category":    "Standard",
						"cluster_mode":        "Integrated",
						"security_group_id":   CHECKSET,
						"addons.#":            "1",
						"cluster_name":        "full-test-cluster",
						"deletion_protection": "true",
						"shared_storages.#":   "3",
						"cluster_vswitch_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":        "full-test-modified-cluster",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":        "full-test-modified-cluster",
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addons", "cluster_credentials"},
			},
		},
	})
}

var AlicloudEhpcClusterV2Map12035 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEhpcClusterV2BasicDependence12035(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "full_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_access_group" "full_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "full_test_nas" {
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

resource "alicloud_vswitch" "full_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.full_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_access_rule" "full_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.full_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.full_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}

resource "alicloud_nas_mount_target" "full_test_mount_domain" {
  vpc_id            = alicloud_vpc.full_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.full_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.full_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.full_test_nas.id
}

resource "alicloud_security_group" "full_test_security_group" {
  vpc_id              = alicloud_vpc.full_test_vpc.id
  security_group_type = "normal"
}


`, name)
}

// Case minimal_cluster_test 12036
func TestAccAliCloudEhpcClusterV2_basic12036(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterV2Map12036)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcClusterV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccehpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterV2BasicDependence12036)
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
					"cluster_credentials": []map[string]interface{}{
						{
							"password": "aliHPC123",
						},
					},
					"cluster_vpc_id":      "${alicloud_vpc.minimal_test_vpc.id}",
					"cluster_category":    "Standard",
					"cluster_mode":        "Integrated",
					"security_group_id":   "${alicloud_security_group.minimal_test_security_group.id}",
					"cluster_name":        "minimal-test-cluster",
					"deletion_protection": "false",
					"shared_storages": []map[string]interface{}{
						{
							"mount_directory":     "/home",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/opt",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
						{
							"mount_directory":     "/ehpcdata",
							"nas_directory":       "/",
							"mount_target_domain": "${alicloud_nas_mount_target.minimal_test_mount_domain.mount_target_domain}",
							"protocol_type":       "NFS",
							"file_system_id":      "${alicloud_nas_file_system.minimal_test_nas.id}",
							"mount_options":       "-t nfs -o vers=3,nolock,proto=tcp,noresvport",
						},
					},
					"cluster_vswitch_id": "${alicloud_vswitch.minimal_test_vswitch.id}",
					"manager": []map[string]interface{}{
						{
							"manager_node": []map[string]interface{}{
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
									"instance_type":        "ecs.c6.xlarge",
									"spot_strategy":        "SpotWithPriceLimit",
									"spot_price_limit":     "3",
								},
							},
							"scheduler": []map[string]interface{}{
								{
									"type":    "SLURM",
									"version": "22.05.8",
								},
							},
							"dns": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
							"directory_service": []map[string]interface{}{
								{
									"type":    "nis",
									"version": "1.0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_vpc_id":      CHECKSET,
						"cluster_category":    "Standard",
						"cluster_mode":        "Integrated",
						"security_group_id":   CHECKSET,
						"cluster_name":        "minimal-test-cluster",
						"deletion_protection": "false",
						"shared_storages.#":   "3",
						"cluster_vswitch_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addons", "cluster_credentials"},
			},
		},
	})
}

var AlicloudEhpcClusterV2Map12036 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEhpcClusterV2BasicDependence12036(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "minimal_test_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "test-cluster-vpc"
}

resource "alicloud_nas_access_group" "minimal_test_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_vswitch" "minimal_test_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.minimal_test_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "test-cluster-vsw"
}

resource "alicloud_nas_file_system" "minimal_test_nas" {
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

resource "alicloud_nas_mount_target" "minimal_test_mount_domain" {
  vpc_id            = alicloud_vpc.minimal_test_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.minimal_test_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.minimal_test_vswitch.id
  file_system_id    = alicloud_nas_file_system.minimal_test_nas.id
}

resource "alicloud_security_group" "minimal_test_security_group" {
  vpc_id              = alicloud_vpc.minimal_test_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_access_rule" "minimal_test_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.minimal_test_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.minimal_test_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}


`, name)
}

// Test Ehpc ClusterV2. <<< Resource test cases, automatically generated.
