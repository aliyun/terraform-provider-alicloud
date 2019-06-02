package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"

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
		Dependencies: []string{
			"alicloud_cen_bandwidth_package",
		},
	})
}

func testSweepCenInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testAcc%s", defaultRegionToTest),
		fmt.Sprintf("tf_testAcc%s", defaultRegionToTest),
	}

	var insts []cbn.Cen
	describeCensRequest := cbn.CreateDescribeCensRequest()
	describeCensRequest.RegionId = client.RegionId
	describeCensRequest.PageSize = requests.NewInteger(PageSizeLarge)
	describeCensRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCens(describeCensRequest)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CEN Instances: %s", err)
		}
		describeCensResponse, _ := raw.(*cbn.DescribeCensResponse)
		if len(describeCensResponse.Cens.Cen) < 1 {
			break
		}
		insts = append(insts, describeCensResponse.Cens.Cen...)

		if len(describeCensResponse.Cens.Cen) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(describeCensRequest.PageNumber); err != nil {
			return err
		} else {
			describeCensRequest.PageNumber = page
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
		describeCenAttachedChildInstancesRequest := cbn.CreateDescribeCenAttachedChildInstancesRequest()
		describeCenAttachedChildInstancesRequest.CenId = id
		raw, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DescribeCenAttachedChildInstances(describeCenAttachedChildInstancesRequest)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to Describe CEN Attached Instance (%s (%s)): %s", name, id, err)
		}
		describeCenAttachedChildInstancesResponse, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)
		for _, childInstance := range describeCenAttachedChildInstancesResponse.ChildInstances.ChildInstance {
			instanceId := childInstance.ChildInstanceId
			log.Printf("[INFO] Detaching CEN Child Instance: %s (%s %s)", name, id, instanceId)
			detachCenChildInstanceRequest := cbn.CreateDetachCenChildInstanceRequest()
			detachCenChildInstanceRequest.CenId = id
			detachCenChildInstanceRequest.ChildInstanceId = instanceId
			detachCenChildInstanceRequest.ChildInstanceType = childInstance.ChildInstanceType
			detachCenChildInstanceRequest.ChildInstanceRegionId = childInstance.ChildInstanceRegionId
			_, err := client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
				return cenClient.DetachCenChildInstance(detachCenChildInstanceRequest)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to Detach CEN Attached Instance (%s (%s %s)): %s", name, id, instanceId, err)
			}
			cenService := CenService{client}
			err = cenService.WaitForCenChildInstanceDetached(instanceId, id, DefaultCenTimeoutLong)
			if err != nil {
				log.Printf("[ERROR] Failed to WaitFor CEN Attached Instance Detached (%s (%s %s)): %s", name, id, instanceId, err)
			}
		}
		log.Printf("[INFO] Deleting CEN Instance: %s (%s)", name, id)
		deleteCenRequest := cbn.CreateDeleteCenRequest()
		deleteCenRequest.CenId = id
		_, err = client.WithCenClient(func(cenClient *cbn.Client) (interface{}, error) {
			return cenClient.DeleteCen(deleteCenRequest)
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
	resourceId := "alicloud_cen_instance.default"
	ra := resourceAttrInit(resourceId, cenInstanceMap)
	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cen, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccCenInstanceNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sCenConfig-%d-N", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccCenInstanceDescriptionConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "tf-testAccCenConfigDescription-N"}),
				),
			},
			{
				Config: testAccCenInstanceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCenConfig-%d", defaultRegionToTest, rand),
						"description": "tf-testAccCenConfigDescription",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudCenInstance_multi(t *testing.T) {
	var cen cbn.Cen
	resourceId := "alicloud_cen_instance.default.9"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cen, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCenConfig-%d-9", defaultRegionToTest, rand),
						"description": "tf-testAccCenConfigDescription",
					}),
				),
			},
		},
	})
}

var cenInstanceMap = map[string]string{
	"description": "tf-testAccCenConfigDescription",
}

func testAccCenInstanceConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
		name = "tf-testAcc%sCenConfig-%d"
		description = "tf-testAccCenConfigDescription"
}
`, defaultRegionToTest, rand)
}
func testAccCenInstanceNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
		name = "tf-testAcc%sCenConfig-%d-N"
		description = "tf-testAccCenConfigDescription"
}
`, defaultRegionToTest, rand)
}
func testAccCenInstanceDescriptionConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
		name = "tf-testAcc%sCenConfig-%d-N"
		description = "tf-testAccCenConfigDescription-N"
}
`, defaultRegionToTest, rand)
}
func testAccCenInstanceMultiConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
		name = "tf-testAcc%sCenConfig-%d-${count.index}"
		description = "tf-testAccCenConfigDescription"
		count = 10
}
`, defaultRegionToTest, rand)
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
