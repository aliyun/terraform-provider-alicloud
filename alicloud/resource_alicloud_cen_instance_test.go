package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_instance", &resource.Sweeper{
		Name: "alicloud_cen_instance",
		F:    testSweepCenInstances,
	})
}

func testSweepCenInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
	}

	var insts []cbn.Cen
	req := cbn.CreateDescribeCensRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCens(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CEN Instances: %s", err)
		}
		resp, _ := raw.(*cbn.DescribeCensResponse)
		if resp == nil || len(resp.Cens.Cen) < 1 {
			break
		}
		insts = append(insts, resp.Cens.Cen...)

		if len(resp.Cens.Cen) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.Name
		id := v.CenId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CEN Instance: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting CEN Instance: %s (%s)", name, id)
		req := cbn.CreateDeleteCenRequest()
		req.CenId = id
		_, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DeleteCen(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 5 seconds to eusure these instances have been deleted.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudCenInstance_basic(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceExists("alicloud_cen_instance.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "name", "tf-testAccCenConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "description", "tf-testAccCenConfigDescription"),
				),
			},
		},
	})

}

func TestAccAlicloudCenInstance_update(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceExists("alicloud_cen_instance.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "name", "tf-testAccCenConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "description", "tf-testAccCenConfigDescription"),
				),
			},
			resource.TestStep{
				Config: testAccCenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceExists("alicloud_cen_instance.foo", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "name", "tf-testAccCenConfigUpdate"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.foo", "description", "tf-testAccCenConfigDescriptionUpdate"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstance_multi(t *testing.T) {
	var cen cbn.Cen

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceExists("alicloud_cen_instance.bar_1", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_1", "name", "tf-testAccCenConfig-1"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_1", "description", "tf-testAccCenConfigDescription-1"),
					testAccCheckCenInstanceExists("alicloud_cen_instance.bar_2", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_2", "name", "tf-testAccCenConfig-2"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_2", "description", "tf-testAccCenConfigDescription-2"),
					testAccCheckCenInstanceExists("alicloud_cen_instance.bar_3", &cen),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_3", "name", "tf-testAccCenConfig-3"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_instance.bar_3", "description", "tf-testAccCenConfigDescription-3"),
				),
			},
		},
	})
}

func testAccCheckCenInstanceExists(n string, cen *cbn.Cen) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CEN ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cenService := CenService{client}
		instance, err := cenService.DescribeCenInstance(rs.Primary.ID)

		if err != nil {
			return err
		}

		*cen = instance
		return nil
	}
}

func testAccCheckCenInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_instance" {
			continue
		}

		// Try to find the CEN
		cenService := CenService{client}
		instance, err := cenService.DescribeCenInstance(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.CenId != "" {
			return fmt.Errorf("CEN %s still exist", instance.CenId)
		}
	}

	return nil
}

const testAccCenConfig = `
resource "alicloud_cen_instance" "foo" {
	name = "tf-testAccCenConfig"
	description = "tf-testAccCenConfigDescription"
}
`

const testAccCenConfigUpdate = `
resource "alicloud_cen_instance" "foo" {
	name = "tf-testAccCenConfigUpdate"
	description = "tf-testAccCenConfigDescriptionUpdate"
}
`

const testAccCenConfigMulti = `
resource "alicloud_cen_instance" "bar_1" {
	name = "tf-testAccCenConfig-1"
	description = "tf-testAccCenConfigDescription-1"
}
resource "alicloud_cen_instance" "bar_2" {
	name = "tf-testAccCenConfig-2"
	description = "tf-testAccCenConfigDescription-2"
}
resource "alicloud_cen_instance" "bar_3" {
	name = "tf-testAccCenConfig-3"
	description = "tf-testAccCenConfigDescription-3"
}
`
