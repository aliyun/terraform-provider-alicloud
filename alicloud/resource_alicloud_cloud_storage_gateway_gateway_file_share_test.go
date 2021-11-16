package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayGatewayFileShare_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_file_share.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayFileShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayFileShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewayfileshare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayFileShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"local_path":              "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"gateway_file_share_name": "${var.name}",
					"oss_bucket_name":         "${alicloud_oss_bucket.default.bucket}",
					"oss_endpoint":            "${alicloud_oss_bucket.default.extranet_endpoint}",
					"protocol":                "NFS",
					"remote_sync":             "true",
					"polling_interval":        "4500",
					"ignore_delete":           "false",
					"fe_limit":                "0",
					"backend_limit":           "0",
					"in_place":                "false",
					"cache_mode":              "Cache",
					"squash":                  "none",
					"ro_client_list":          "12.12.12.12",
					"rw_client_list":          "12.12.12.12",
					"oss_bucket_ssl":          "false",
					"lag_period":              "5",
					"direct_io":               "false",
					"path_prefix":             "/home",
					"fast_reclaim":            "false",
					"support_archive":         "false",
					"nfs_v4_optimization":     "false",
					"transfer_acceleration":   "false",
					"bypass_cache_read":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":              CHECKSET,
						"local_path":              CHECKSET,
						"gateway_file_share_name": name,
						"oss_bucket_name":         CHECKSET,
						"oss_endpoint":            CHECKSET,
						"protocol":                "NFS",
						"remote_sync":             "true",
						"polling_interval":        "4500",
						"ignore_delete":           "false",
						"fe_limit":                "0",
						"backend_limit":           "0",
						"in_place":                "false",
						"cache_mode":              "Cache",
						"squash":                  "none",
						"ro_client_list":          "12.12.12.12",
						"rw_client_list":          "12.12.12.12",
						"oss_bucket_ssl":          "false",
						"lag_period":              "5",
						"direct_io":               "false",
						"path_prefix":             "/home",
						"fast_reclaim":            "false",
						"support_archive":         "false",
						"nfs_v4_optimization":     "false",
						"transfer_acceleration":   "false",
						"bypass_cache_read":       "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_limit": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_limit": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fe_limit": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fe_limit": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rw_client_list": "13.13.13.13",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_client_list": "13.13.13.13",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"squash": "root_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"squash": "root_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"squash": "all_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"squash": "all_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"squash": "all_anonymous",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"squash": "all_anonymous",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bypass_cache_read": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bypass_cache_read": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transfer_acceleration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transfer_acceleration": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nfs_v4_optimization": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nfs_v4_optimization": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"polling_interval": "5000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"polling_interval": "5000",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"lag_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lag_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ro_client_list": "13.13.13.13",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ro_client_list": "13.13.13.13",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_reclaim": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fast_reclaim": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_sync":      "false",
					"polling_interval": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_sync":      "false",
						"polling_interval": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ignore_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ignore_delete": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_limit":         "0",
					"fe_limit":              "0",
					"rw_client_list":        "12.12.12.12",
					"remote_sync":           "true",
					"squash":                "none",
					"transfer_acceleration": "false",
					"nfs_v4_optimization":   "false",
					"ignore_delete":         "false",
					"polling_interval":      "4500",
					"lag_period":            "5",
					"ro_client_list":        "12.12.12.12",
					"fast_reclaim":          "false",
					"bypass_cache_read":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_limit":         "0",
						"fe_limit":              "0",
						"rw_client_list":        "12.12.12.12",
						"remote_sync":           "true",
						"squash":                "none",
						"transfer_acceleration": "false",
						"nfs_v4_optimization":   "false",
						"ignore_delete":         "false",
						"polling_interval":      "4500",
						"lag_period":            "5",
						"ro_client_list":        "12.12.12.12",
						"fast_reclaim":          "false",
						"bypass_cache_read":     "false",
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

func TestAccAlicloudCloudStorageGatewayGatewayFileShare_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_file_share.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayFileShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayFileShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewayfileshare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayFileShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":               "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"local_path":               "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"gateway_file_share_name":  "${var.name}",
					"oss_bucket_name":          "${alicloud_oss_bucket.default.bucket}",
					"oss_endpoint":             "${alicloud_oss_bucket.default.intranet_endpoint}",
					"protocol":                 "SMB",
					"remote_sync":              "true",
					"polling_interval":         "4500",
					"ignore_delete":            "false",
					"fe_limit":                 "0",
					"backend_limit":            "0",
					"in_place":                 "true",
					"cache_mode":               "Sync",
					"browsable":                "false",
					"oss_bucket_ssl":           "true",
					"lag_period":               "5",
					"direct_io":                "true",
					"ro_user_list":             "user1",
					"rw_user_list":             "user1",
					"path_prefix":              "",
					"fast_reclaim":             "false",
					"support_archive":          "false",
					"windows_acl":              "false",
					"access_based_enumeration": "false",
					"transfer_acceleration":    "false",
					"remote_sync_download":     "false",
					"download_limit":           "0",
					"partial_sync_paths":       "/root/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":               CHECKSET,
						"local_path":               CHECKSET,
						"gateway_file_share_name":  name,
						"oss_bucket_name":          CHECKSET,
						"oss_endpoint":             CHECKSET,
						"protocol":                 "SMB",
						"remote_sync":              "true",
						"polling_interval":         "4500",
						"ignore_delete":            "false",
						"fe_limit":                 "0",
						"backend_limit":            "0",
						"in_place":                 "true",
						"cache_mode":               "Sync",
						"browsable":                "false",
						"oss_bucket_ssl":           "true",
						"lag_period":               "5",
						"direct_io":                "true",
						"ro_user_list":             "user1",
						"rw_user_list":             "user1",
						"path_prefix":              "",
						"fast_reclaim":             "false",
						"support_archive":          "false",
						"windows_acl":              "false",
						"access_based_enumeration": "false",
						"transfer_acceleration":    "false",
						"remote_sync_download":     "false",
						"download_limit":           "0",
						"partial_sync_paths":       "/root/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"windows_acl":              "true",
					"access_based_enumeration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"windows_acl":              "true",
						"access_based_enumeration": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_limit": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_limit": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"browsable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"browsable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fe_limit": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fe_limit": "100",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"remote_sync_download": "true",
					"download_limit":       "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_sync_download": "true",
						"download_limit":       "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ro_user_list": "user1,user2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ro_user_list": "user1,user2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"polling_interval": "5000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"polling_interval": "5000",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"lag_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lag_period": "10",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"rw_user_list": "user1,user2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_user_list": "user1,user2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_reclaim": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fast_reclaim": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_sync":          "false",
					"polling_interval":     "0",
					"remote_sync_download": "false",
					"download_limit":       "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_sync":          "false",
						"polling_interval":     "0",
						"remote_sync_download": "false",
						"download_limit":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_based_enumeration": "false",
					"backend_limit":            "0",
					"browsable":                "false",
					"fe_limit":                 "0",
					"remote_sync":              "true",
					"windows_acl":              "false",
					"ro_user_list":             "user1",
					"polling_interval":         "4500",
					"lag_period":               "5",
					"rw_user_list":             "user1",
					"fast_reclaim":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_based_enumeration": "false",
						"backend_limit":            "0",
						"browsable":                "false",
						"fe_limit":                 "0",
						"remote_sync":              "true",
						"windows_acl":              "false",
						"ro_user_list":             "user1",
						"polling_interval":         "4500",
						"lag_period":               "5",
						"rw_user_list":             "user1",
						"fast_reclaim":             "false",
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

func TestAccAlicloudCloudStorageGatewayGatewayFileShare_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_file_share.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayFileShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayFileShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cloudstoragegatewaygatewayfileshare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayFileShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"local_path":              "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"gateway_file_share_name": "${var.name}",
					"oss_bucket_name":         "${alicloud_oss_bucket.default.bucket}",
					"oss_endpoint":            "${alicloud_oss_bucket.default.intranet_endpoint}",
					"protocol":                "SMB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":              CHECKSET,
						"local_path":              CHECKSET,
						"gateway_file_share_name": name,
						"oss_bucket_name":         CHECKSET,
						"oss_endpoint":            CHECKSET,
						"protocol":                "SMB",
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

var AlicloudCloudStorageGatewayGatewayFileShareMap0 = map[string]string{
	"index_id":                 CHECKSET,
	"remote_sync_download":     CHECKSET,
	"fast_reclaim":             CHECKSET,
	"access_based_enumeration": CHECKSET,
	"windows_acl":              CHECKSET,
	"ignore_delete":            CHECKSET,
	"direct_io":                CHECKSET,
	"browsable":                CHECKSET,
	"gateway_id":               CHECKSET,
	"fe_limit":                 CHECKSET,
	"backend_limit":            CHECKSET,
	"download_limit":           CHECKSET,
	"nfs_v4_optimization":      CHECKSET,
}

func AlicloudCloudStorageGatewayGatewayFileShareBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cloud_storage_gateway_stocks" "default" {
gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
vpc_name   = var.name
cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
vpc_id       = alicloud_vpc.vpc.id
cidr_block   = "172.16.0.0/21"
zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
description              = "tf-acctestDesalone"
gateway_class            = "Standard"
type                     = "File"
payment_type             = "PayAsYouGo"
vswitch_id               = alicloud_vswitch.default.id
release_after_expiration = true
public_network_bandwidth = 10
storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
location                 = "Cloud"
gateway_name             = var.name
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
cache_disk_category   = "cloud_efficiency"
gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}
`, name)
}
