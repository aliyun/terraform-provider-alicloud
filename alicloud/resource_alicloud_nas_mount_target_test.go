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
	resource.AddTestSweepers("alicloud_nas_mounttarget", &resource.Sweeper{
		Name: "alicloud_nas_mounttarget",
		F:    testSweepNasMT,
	})
}

func testSweepNasMT(region string) error {
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

	var mt []nas.MountTarget
	req := nas.CreateDescribeMountTargetsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving filesystem: %s", err)
		}
		resp, _ := raw.(*nas.DescribeMountTargetsResponse)
		if resp == nil || len(resp.MountTargets.MountTarget) < 1 {
			break
		}
		mt = append(mt, resp.MountTargets.MountTarget...)

		if len(resp.MountTargets.MountTarget) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, fs := range mt {

		id := fs.MountTargetDomain
		AccessGroupName := fs.AccessGroup
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(AccessGroupName), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping MountTarget: %s (%s)", AccessGroupName, id)
			continue
		}
		log.Printf("[INFO] Deleting MountTarget: %s (%s)", AccessGroupName, id)
		req := nas.CreateDeleteMountTargetRequest()
		req.FileSystemId = id
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteMountTarget(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete MountTarget (%s (%s)): %s", AccessGroupName, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNas_MountTarget_basic(t *testing.T) {
	var mt nas.MountTarget
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_nas_mounttarget.foo:",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMtDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasMtConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMtExists("alicloud_nas_mounttarget.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mounttarget.foo", "networktype", "Classic"),
				),
			},
		},
	})

}

func TestAccAlicloudNas_MountTarget_update(t *testing.T) {
	var mt nas.MountTarget

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMtDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasMtConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMtExists("alicloud_nas_mounttarget.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mounttarget.foo", "accessgroup_name", "test_wang"),
				),
			},
			resource.TestStep{
				Config: testAccNasMtConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMtExists("alicloud_nas_mounttarget.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mounttarget.foo", "accessgroup_name", "wang_test"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mounttarget.foo", "status", "Active"),
				),
			},
		},
	})
}

func testAccCheckMtExists(n string, nas *nas.MountTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NAS ID is set")
		}
		split := strings.Split(rs.Primary.ID, ":")
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		instance, err := nasService.DescribeMountTargets(split[0])

		if err != nil {
			return err
		}

		*nas = instance
		return nil
	}
}

func testAccCheckMtDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_filesystem" {
			continue
		}

		// Try to find the NAS
		split := strings.Split(rs.Primary.ID, ":")
		instance, err := nasService.DescribeMountTargets(split[0])

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.MountTargetDomain != "" {
			return fmt.Errorf("NAS %s still exist", instance.MountTargetDomain)
		}
	}

	return nil
}

const testAccNasMtConfig = `
resource "alicloud_nas_filesystem" "foo" {
		protocol_type = "NFS"
		storage_type = "Performance"
}
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_mounttarget" "foo" {
		filesystem_id = "${alicloud_nas_filesystem.foo.id}"
		accessgroup_name = "${alicloud_nas_accessgroup.foo.id}"
		networktype = "Classic"
}
`

const testAccNasMtConfigUpdate = `
resource "alicloud_nas_filesystem" "foo" {
		protocol_type = "NFS"
		storage_type = "Performance"
}
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_mounttarget" "foo" {
		filesystem_id = "${alicloud_nas_filesystem.foo.id}"
		accessgroup_name = "{alicloud_nas_accessgroup.foo.id}"
		networktype = "Classic"
		status = "Active"
}
`
