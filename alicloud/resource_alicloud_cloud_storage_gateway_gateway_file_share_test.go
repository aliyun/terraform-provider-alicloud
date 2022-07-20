package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
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
					"cache_mode":              "Cache",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":              CHECKSET,
						"local_path":              CHECKSET,
						"gateway_file_share_name": name,
						"oss_bucket_name":         CHECKSET,
						"oss_endpoint":            CHECKSET,
						"protocol":                "NFS",
						"cache_mode":              "Cache",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_sync":      "true",
					"polling_interval": "4500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_sync":      "true",
						"polling_interval": "4500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ro_client_list": "12.12.12.12",
					"rw_client_list": "12.12.12.12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ro_client_list": "12.12.12.12",
						"rw_client_list": "12.12.12.12",
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
					"gateway_id":              "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"local_path":              "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"gateway_file_share_name": "${var.name}",
					"oss_bucket_name":         "${alicloud_oss_bucket.default.bucket}",
					"oss_endpoint":            "${alicloud_oss_bucket.default.intranet_endpoint}",
					"protocol":                "SMB",
					"cache_mode":              "Sync",
					"partial_sync_paths":      "/root/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":              CHECKSET,
						"local_path":              CHECKSET,
						"gateway_file_share_name": name,
						"oss_bucket_name":         CHECKSET,
						"oss_endpoint":            CHECKSET,
						"protocol":                "SMB",
						"cache_mode":              "Sync",
						"partial_sync_paths":      "/root/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_sync":      "true",
					"polling_interval": "4500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_sync":      "true",
						"polling_interval": "4500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ro_user_list": "user1",
					"rw_user_list": "user1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ro_user_list": "user1",
						"rw_user_list": "user1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"browsable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"browsable": "false",
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
					"windows_acl":              "true",
					"access_based_enumeration": "true",
					"transfer_acceleration":    "false",
					"remote_sync_download":     "true",
					"download_limit":           "1000",
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
						"windows_acl":              "true",
						"access_based_enumeration": "true",
						"transfer_acceleration":    "false",
						"remote_sync_download":     "true",
						"download_limit":           "1000",
						"partial_sync_paths":       "/root/",
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

func TestAccAlicloudCloudStorageGatewayGatewayFileShare_basic3(t *testing.T) {
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
					"support_archive":         "true",
					"nfs_v4_optimization":     "true",
					"transfer_acceleration":   "true",
					"bypass_cache_read":       "true",
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
						"support_archive":         "true",
						"nfs_v4_optimization":     "true",
						"transfer_acceleration":   "true",
						"bypass_cache_read":       "true",
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

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
description              = "tf-acctestDesalone"
gateway_class            = "Standard"
type                     = "File"
payment_type             = "PayAsYouGo"
vswitch_id               = local.vswitch_id
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

func TestUnitAlicloudCloudStorageGatewayGatewayFileShare(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"access_based_enumeration": true,
		"backend_limit":            1,
		"browsable":                true,
		"cache_mode":               "Cache",
		"direct_io":                true,
		"download_limit":           1,
		"fast_reclaim":             true,
		"fe_limit":                 1,
		"bypass_cache_read":        true,
		"gateway_file_share_name":  "CreateGatewayFileShareValue",
		"ignore_delete":            true,
		"in_place":                 true,
		"lag_period":               5,
		"local_path":               "CreateGatewayFileShareValue",
		"nfs_v4_optimization":      true,
		"oss_bucket_name":          "CreateGatewayFileShareValue",
		"oss_bucket_ssl":           true,
		"oss_endpoint":             "CreateGatewayFileShareValue",
		"partial_sync_paths":       "CreateGatewayFileShareValue",
		"path_prefix":              "CreateGatewayFileShareValue",
		"polling_interval":         1,
		"protocol":                 "CreateGatewayFileShareValue",
		"remote_sync":              true,
		"remote_sync_download":     true,
		"ro_client_list":           "CreateGatewayFileShareValue",
		"ro_user_list":             "CreateGatewayFileShareValue",
		"rw_client_list":           "CreateGatewayFileShareValue",
		"rw_user_list":             "CreateGatewayFileShareValue",
		"squash":                   "none",
		"support_archive":          true,
		"transfer_acceleration":    true,
		"windows_acl":              true,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeGatewayFileShares
		"FileShares": map[string]interface{}{
			"FileShare": []interface{}{
				map[string]interface{}{
					"GatewayId":              "CreateGatewayFileShareValue",
					"IndexId":                "CreateGatewayFileShareValue",
					"AccessBasedEnumeration": true,
					"BeLimit":                1,
					"Browsable":              true,
					"CacheMode":              "Cache",
					"DirectIO":               true,
					"DownloadLimit":          1,
					"FastReclaim":            true,
					"FeLimit":                1,
					"Name":                   "CreateGatewayFileShareValue",
					"IgnoreDelete":           true,
					"InPlace":                true,
					"LagPeriod":              5,
					"LocalPath":              "CreateGatewayFileShareValue",
					"NfsV4Optimization":      true,
					"OssBucketName":          "CreateGatewayFileShareValue",
					"OssBucketSsl":           true,
					"OssEndpoint":            "CreateGatewayFileShareValue",
					"PartialSyncPaths":       "CreateGatewayFileShareValue",
					"PathPrefix":             "CreateGatewayFileShareValue",
					"PollingInterval":        1,
					"Protocol":               "CreateGatewayFileShareValue",
					"BypassCacheRead":        true,
					"RemoteSync":             true,
					"RemoteSyncDownload":     true,
					"RoClientList":           "CreateGatewayFileShareValue",
					"RoUserList":             "CreateGatewayFileShareValue",
					"RwClientList":           "CreateGatewayFileShareValue",
					"RwUserList":             "CreateGatewayFileShareValue",
					"Squash":                 "none",
					"SupportArchive":         true,
					"TransferAcceleration":   true,
					"WindowsAcl":             true,
				},
			},
		},
		"Tasks": map[string]interface{}{
			"SimpleTask": []interface{}{
				map[string]interface{}{
					"TaskId":            "CreateGatewayFileShareValue",
					"StateCode":         "task.state.completed",
					"RelatedResourceId": "CreateGatewayFileShareValue",
				},
			},
		},
		"TaskId":  "CreateGatewayFileShareValue",
		"Success": true,
	}
	CreateMockResponse := map[string]interface{}{
		// CreateGatewayFileShare
		"FileShares": map[string]interface{}{
			"FileShare": []interface{}{
				map[string]interface{}{
					"IndexId": "CreateGatewayFileShareValue",
				},
			},
		},
		"Tasks": map[string]interface{}{
			"SimpleTask": []interface{}{
				map[string]interface{}{
					"TaskId":            "CreateGatewayFileShareValue",
					"StateCode":         "task.state.completed",
					"RelatedResourceId": "CreateGatewayFileShareValue",
				},
			},
		},
		"TaskId":  "CreateGatewayFileShareValue",
		"Success": true,
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_storage_gateway_gateway_file_share", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudStorageGatewayGatewayFileShareCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeGatewayFileShares Response
		"FileShares": map[string]interface{}{
			"FileShare": []interface{}{
				map[string]interface{}{
					"IndexId": "CreateGatewayFileShareValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateGatewayFileShare" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCloudStorageGatewayGatewayFileShareCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudStorageGatewayGatewayFileShareUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateGatewayFileShare
	attributesDiff := map[string]interface{}{
		"gateway_file_share_name":  "UpdateGatewayFileShareValue",
		"access_based_enumeration": false,
		"backend_limit":            2,
		"browsable":                false,
		"bypass_cache_read":        false,
		"cache_mode":               "Sync",
		"download_limit":           2,
		"fast_reclaim":             false,
		"fe_limit":                 2,
		"ignore_delete":            false,
		"in_place":                 false,
		"lag_period":               10,
		"nfs_v4_optimization":      false,
		"polling_interval":         2,
		"remote_sync":              false,
		"remote_sync_download":     false,
		"ro_client_list":           "UpdateGatewayFileShareValue",
		"ro_user_list":             "UpdateGatewayFileShareValue",
		"rw_client_list":           "UpdateGatewayFileShareValue",
		"rw_user_list":             "UpdateGatewayFileShareValue",
		"squash":                   "root_squash",
		"transfer_acceleration":    false,
		"windows_acl":              false,
	}
	diff, err := newInstanceDiff("alicloud_cloud_storage_gateway_gateway_file_share", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeGatewayFileShares Response
		"FileShares": map[string]interface{}{
			"FileShare": []interface{}{
				map[string]interface{}{
					"AccessBasedEnumeration": false,
					"Name":                   "UpdateGatewayFileShareValue",
					"BeLimit":                2,
					"Browsable":              false,
					"BypassCacheRead":        false,
					"CacheMode":              "Sync",
					"DownloadLimit":          2,
					"FastReclaim":            false,
					"FeLimit":                2,
					"IgnoreDelete":           false,
					"InPlace":                false,
					"KmsRotatePeriod":        2,
					"LagPeriod":              10,
					"NfsV4Optimization":      false,
					"PollingInterval":        2,
					"RemoteSync":             false,
					"RemoteSyncDownload":     false,
					"RoClientList":           "UpdateGatewayFileShareValue",
					"RoUserList":             "UpdateGatewayFileShareValue",
					"RwClientList":           "UpdateGatewayFileShareValue",
					"RwUserList":             "UpdateGatewayFileShareValue",
					"Squash":                 "root_squash",
					"TransferAcceleration":   false,
					"WindowsAcl":             false,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateGatewayFileShare" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCloudStorageGatewayGatewayFileShareUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeGatewayFileShares" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCloudStorageGatewayGatewayFileShareRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudCloudStorageGatewayGatewayFileShareDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cloud_storage_gateway_gateway_file_share", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_file_share"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "GatewayDeletionError", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteGatewayFileShares" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Tasks": map[string]interface{}{
								"SimpleTask": []interface{}{
									map[string]interface{}{
										"TaskId":    "CreateGatewayFileShareValue",
										"StateCode": "task.state.completed",
									},
								},
							},
							"TaskId":  "CreateGatewayFileShareValue",
							"Success": true,
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCloudStorageGatewayGatewayFileShareDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		}
	}

}
