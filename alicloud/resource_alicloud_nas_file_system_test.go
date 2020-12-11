package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_file_system",
		&resource.Sweeper{
			Name: "alicloud_nas_file_syster",
			F:    testSweepNasFileSystem,
			// When implemented, these should be removed firstly
			Dependencies: []string{
				"alicloud_nas_mount_target",
			},
		})
}

func testSweepNasFileSystem(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var filesystems []nas.FileSystem
	req := nas.CreateDescribeFileSystemsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(req)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving filesystem: %s", err)
		}
		resp, _ := raw.(*nas.DescribeFileSystemsResponse)
		if resp == nil || len(resp.FileSystems.FileSystem) < 1 {
			break
		}
		filesystems = append(filesystems, resp.FileSystems.FileSystem...)

		if len(resp.FileSystems.FileSystem) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, fs := range filesystems {

		id := fs.FileSystemId
		destription := fs.Description
		domain := fs.MountTargets.MountTarget
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(destription), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping FileSystem: %s (%s)", destription, id)
			continue
		}
		if len(domain) > 0 {
			for _, mount_target := range domain {
				request := nas.CreateDeleteMountTargetRequest()
				request.FileSystemId = id
				request.MountTargetDomain = mount_target.MountTargetDomain
				_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
					return nasClient.DeleteMountTarget(request)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete MountTarget (%s (%s)): %s", destription, id, err)
				}
			}
		}
		log.Printf("[INFO] Deleting FileSystem: %s (%s)", destription, id)
		req := nas.CreateDeleteFileSystemRequest()
		req.FileSystemId = id
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteFileSystem(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete FileSystem (%s (%s)): %s", destription, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNas_FileSystem_basic(t *testing.T) {
	var v nas.FileSystem
	resourceID := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"description":   "tf-testAccNasConfig",
						"storage_type":  "Performance",
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccNasConfigUpdateName",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudNas_FileSystem_basicT(t *testing.T) {
	var v nas.FileSystem
	resourceID := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasConfigT(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"description":   "tf-testAccNasConfig",
						"storage_type":  "Capacity",
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasConfigUpdateT(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccNasConfigUpdateName",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudNas_FileSystem_multi(t *testing.T) {
	var v nas.FileSystem
	resourceID := "alicloud_nas_file_system.default.2"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasConfigMulti(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"description":   "tf-testAccNasConfig_multil-1",
						"storage_type":  "Performance",
					}),
				),
			},
			{
				Config: testAccNasConfigMultiT(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"description":   "tf-testAccNasConfig_multil-1",
						"storage_type":  "Capacity",
					}),
				),
			},
		},
	})
}

func testAccCheckNasExists(n string, nas *nas.FileSystem) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No NAS ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		instance, err := nasService.DescribeNasFileSystem(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*nas = instance
		return nil
	}
}

func testAccCheckNasDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_file_system" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeNasFileSystem(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return WrapError(fmt.Errorf("NAS %s still exist", instance.FileSystemId))
	}
	return nil
}

func testAccNasConfig() string {
	return fmt.Sprintf(`
variable "storage_type" {
  default = "Performance"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
	protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
	storage_type = "${var.storage_type}"
	description = "tf-testAccNasConfig"
}`)
}

func testAccNasConfigT() string {
	return fmt.Sprintf(`
variable "storage_type" {
  default = "Capacity"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        storage_type = "${var.storage_type}"
        description = "tf-testAccNasConfig"
}`)
}

func testAccNasConfigUpdate() string {
	return fmt.Sprintf(`
variable "storage_type" {
  default = "Performance"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        storage_type = "${var.storage_type}"
        description = "tf-testAccNasConfigUpdateName"
}`)
}

func testAccNasConfigUpdateT() string {
	return fmt.Sprintf(`
variable "storage_type" {
  default = "Capacity"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        storage_type = "${var.storage_type}"
        description = "tf-testAccNasConfigUpdateName"
}`)
}

func testAccNasConfigMulti() string {
	return fmt.Sprintf(`
variable "description" {
        default = "tf-testAccNasConfig_multil"
}

data "alicloud_nas_protocols" "bar_1" {
        type = "Performance"
}

resource "alicloud_nas_file_system" "default" {
        protocol_type = "${data.alicloud_nas_protocols.bar_1.protocols.0}"
        storage_type = "Performance"
        description = "${var.description}-1"
	count = 3
}`)
}

func testAccNasConfigMultiT() string {
	return fmt.Sprintf(`
variable "description" {
        default = "tf-testAccNasConfig_multil"
}

data "alicloud_nas_protocols" "bar_1" {
        type = "Capacity"
}

resource "alicloud_nas_file_system" "default" {
        protocol_type = "${data.alicloud_nas_protocols.bar_1.protocols.0}"
        storage_type = "Capacity"
        description = "${var.description}-1"
        count = 3
}`)
}
