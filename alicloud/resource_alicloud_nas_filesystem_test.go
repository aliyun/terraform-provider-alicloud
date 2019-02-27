package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_filesystem", &resource.Sweeper{
		Name: "alicloud_nas_filesystem",
		F:    testSweepNas,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_nas_mounttarget",
		},
	})
}

func testSweepNas(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
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
			return fmt.Errorf("Error retrieving filesystem: %s", err)
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
		destription := fs.Destription
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
	var fs nas.FileSystem
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasExists("alicloud_nas_filesystem.foo", &fs),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.foo", "protocol_type", "NFS"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.foo", "storage_type", "Performance"),
				),
			},
		},
	})

}

func TestAccAlicloudNas_FileSystem_update(t *testing.T) {
	var fs nas.FileSystem

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasExists("alicloud_nas_filesystem.foo", &fs),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.foo", "protocol_type", "NFS"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.foo", "storage_type", "Performance"),
				),
			},
			resource.TestStep{
				Config: testAccNasConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasExists("alicloud_nas_filesystem.foo", &fs),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.foo", "description", "wang_test"),
				),
			},
		},
	})
}

func TestAccAlicloudNas_FileSystem_multi(t *testing.T) {
	var fs nas.FileSystem

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasExists("alicloud_nas_filesystem.bar_1", &fs),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.bar_1", "protocol_type", "NFS"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.bar_1", "storage_type", "Performance"),
					testAccCheckNasExists("alicloud_nas_filesystem.bar_2", &fs),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.bar_2", "protocol_type", "SMB"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_filesystem.bar_2", "storage_type", "Capacity"),
				),
			},
		},
	})
}

func testAccCheckNasExists(n string, nas *nas.FileSystem) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NAS ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		instance, err := nasService.DescribeFileSystems(rs.Primary.ID)

		if err != nil {
			return err
		}

		*nas = instance
		return nil
	}
}

func testAccCheckNasDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_filesystem" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeFileSystems(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.FileSystemId != "" {
			return fmt.Errorf("NAS %s still exist", instance.FileSystemId)
		}
	}
	return nil
}

const testAccNasConfig = `
resource "alicloud_nas_filesystem" "foo" {
		protocol_type = "NFS"
		storage_type = "Performance"
		description = "test_wang"
}
`

const testAccNasConfigUpdate = `
resource "alicloud_nas_filesystem" "foo" {
		protocol_type = "NFS"
		storage_type = "Performance"
		description = "wang_test"
}
`

const testAccNasConfigMulti = `
variable "description" {
  	default = "test_wang"
}
resource "alicloud_nas_filesystem" "bar_1" {
	protocol_type = "NFS"
	storage_type = "Performance"
	description = "${var.description}-1"
}
resource "alicloud_nas_filesystem" "bar_2" {
	protocol_type = "SMB"
	storage_type = "Capacity"
	description = "${var.description}-2"
}
`
