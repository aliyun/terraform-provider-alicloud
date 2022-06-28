package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_instance_attachment", &resource.Sweeper{
		Name: "alicloud_cen_instance_attachment",
		F:    testSweepCenInstanceAttachment,
		Dependencies: []string{
			"alicloud_cen_route_service",
		},
	})
}

func testSweepCenInstanceAttachment(region string) error {
	log.Printf("[INFO] Delete cen instance attachment.")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := cbn.CreateDescribeCensRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var cenIds []string
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve cen instance in service list: %s", err)
		}

		response, _ := raw.(*cbn.DescribeCensResponse)

		for _, v := range response.Cens.Cen {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(v.Name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping cen instance: %s ", v.Name)
			} else {
				cenIds = append(cenIds, v.CenId)
			}
		}
		if len(response.Cens.Cen) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, cenId := range cenIds {
		request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)

		for {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenAttachedChildInstances(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete cen instance attachment (%s): %s", cenId, err)
			}
			response, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)

			for _, item := range response.ChildInstances.ChildInstance {
				id := fmt.Sprintf("%v:%v:%v:%v", item.CenId, item.ChildInstanceId, item.ChildInstanceType, item.ChildInstanceRegionId)
				if err := deleteCenInstancAttachmet(id, client); err != nil {
					log.Printf("[ERROR] Failed to delete cen instance attachment (%s): %s", cenId, err)
				} else {
					log.Printf("[INFO] Deleted cen instance attachment success: %s ", id)
				}
			}
			if len(response.ChildInstances.ChildInstance) < PageSizeLarge {
				break
			}
			page, err := getNextpageNumber(request.PageNumber)
			if err != nil {
				return WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

func TestAccAlicloudCenInstanceAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance_attachment.default"
	ra := resourceAttrInit(resourceId, CenInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCenInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-CenInstanceAttachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenInstanceAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":              "${alicloud_cen_instance.default.id}",
					"child_instance_id":        "${alicloud_vpc.default.id}",
					"child_instance_type":      "VPC",
					"child_instance_region_id": defaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              CHECKSET,
						"child_instance_id":        CHECKSET,
						"child_instance_type":      "VPC",
						"child_instance_region_id": defaultRegionToTest,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAlicloudCenInstanceAttachment_multi_same_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiSameRegion(rand, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstanceAttachment_multi_different_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiDifferentRegion(rand, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

var cenInstanceAttachmentMap = map[string]string{
	"instance_id":              CHECKSET,
	"child_instance_id":        CHECKSET,
	"child_instance_region_id": CHECKSET,
}

func testAccCenInstanceAttachmentBasic(rand int, region string) string {
	return fmt.Sprintf(`
	variable "name"{
	    default = "tf-testAcc%sCenInstanceAttachmentBasic-%d"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	    description = "tf-testAccCenInstanceAttachmentBasicDescription"
	}

	resource "alicloud_vpc" "default" {
	    vpc_name = "${var.name}"
	    cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	}
	`, region, rand, region)
}
func testAccCenInstanceAttachmentMultiSameRegion(rand int, region string) string {
	return fmt.Sprintf(`
	provider "alicloud" {
		region = "%[1]s"
	}
	variable "name"{
	    default = "tf-testAcc%[1]sCenInstanceAttachmentBasic-%d"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	    description = "tf-testAccCenInstanceAttachmentBasicDescription"
	}

	data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
	}

	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = data.alicloud_vpcs.default.ids.0
	    child_instance_type = "VPC"
	    child_instance_region_id = "%[1]s"
	}

	resource "alicloud_cen_instance_attachment" "default1" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = data.alicloud_vpcs.default.ids.0
	    child_instance_type = "VPC"
	    child_instance_region_id = "%[1]s"
	}
	`, region, rand)
}

func testAccCenInstanceAttachmentMultiDifferentRegion(rand int, region string) string {
	return fmt.Sprintf(
		`
variable "name"{
    default = "tf-testAccCen%sInstanceAttachmentMultiDifferentRegions-%d"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

data "alicloud_vpcs" "default" {
    provider = "alicloud.fra"
    name_regex = "default-NODELETING"
}

data "alicloud_vpcs" "default1" {
    provider = "alicloud.sh"
    name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance" "default" {
    name = "${var.name}"
    description = "tf-testAccCenInstanceAttachmentMultiDifferentRegionsDescription"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = data.alicloud_vpcs.default.ids.0
	child_instance_type = "VPC"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = data.alicloud_vpcs.default1.ids.0
	child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}
`, region, rand)

}

func testAccCheckCenInstanceAttachmentExistsWithProviders(n string, instance *cbn.DescribeCenAttachedChildInstanceAttributeResponse, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cen Child Instance ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}
			childInstance, err := cbnService.DescribeCenInstanceAttachment(rs.Primary.ID)
			if err != nil {
				return err
			}

			if childInstance.Status != "Attached" {
				return fmt.Errorf("CEN id %s instance id %s status error", childInstance.CenId, childInstance.ChildInstanceId)
			}

			instance = &childInstance
			return nil
		}
		return fmt.Errorf("Cen Child Instance not found")
	}
}

func testAccCheckCenInstanceAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInstanceAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInstanceAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cbnService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}

func deleteCenInstancAttachmet(id string, client *connectivity.AliyunClient) error {
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	request := cbn.CreateDetachCenChildInstanceRequest()
	request.ChildInstanceId = parts[1]
	request.ChildInstanceRegionId = parts[3]
	request.ChildInstanceType = parts[2]
	request.CenId = parts[0]
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DetachCenChildInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return err
}
func resourceCenInstanceAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name"{
	    default = "%s"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	    description = var.name
	}

	resource "alicloud_vpc" "default" {
	    vpc_name = "${var.name}"
	    cidr_block = "192.168.0.0/16"
	}
`, name)
}

var CenInstanceAttachmentBasicMap = map[string]string{
	"child_instance_owner_id": CHECKSET,
	"status":                  CHECKSET,
}
