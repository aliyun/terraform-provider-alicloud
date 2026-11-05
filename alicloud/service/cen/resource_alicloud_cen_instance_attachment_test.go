package cen_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_example",
		"tf-example",
		"terraform-example",
	}

	request := cbn.CreateDescribeCensRequest()
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
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
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(v.Name), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping cen instance: %s ", v.Name)
				}
			}
			cenIds = append(cenIds, v.CenId)
		}
		if len(response.Cens.Cen) < helper.PageSizeLarge {
			break
		}
		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}

	for _, cenId := range cenIds {
		request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(helper.PageSizeLarge)
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
			if len(response.ChildInstances.ChildInstance) < helper.PageSizeLarge {
				break
			}
			page, err := helper.GetNextpageNumber(request.PageNumber)
			if err != nil {
				return errs.WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

func TestAccAliCloudCenInstanceAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_instance_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, CenInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCenInstanceAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-CenInstanceAttachment%v", rand)
	// cen_owner_id / child_instance_owner_id are TypeInt account-owner fields; for
	// same-account attach they equal the calling account, read from the env so no
	// fabricated ID is used. cen_owner_id is write-only (never d.Set in Read), so
	// it is also listed in ImportStateVerifyIgnore below.
	accountId, _ := strconv.Atoi(os.Getenv("ALICLOUD_ACCOUNT_ID"))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenInstanceAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":              "${alicloud_cen_instance.default.id}",
					"child_instance_id":        "${alicloud_vpc.default.id}",
					"child_instance_type":      "VPC",
					"child_instance_region_id": acctest.DefaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              acctest.CHECKSET,
						"child_instance_id":        acctest.CHECKSET,
						"child_instance_type":      "VPC",
						"child_instance_region_id": acctest.DefaultRegionToTest,
					}),
				),
			},

			{
				// Exercise the same-account owner-id fields: cen_owner_id is
				// write-only (never d.Set in Read); child_instance_owner_id is
				// Computed and round-tripped via Read. Both equal the calling
				// account, so the real account id is used rather than a fabricated
				// value.
				Config: testAccConfig(map[string]interface{}{
					"instance_id":              "${alicloud_cen_instance.default.id}",
					"child_instance_id":        "${alicloud_vpc.default.id}",
					"child_instance_type":      "VPC",
					"child_instance_region_id": acctest.DefaultRegionToTest,
					"cen_owner_id":             accountId,
					"child_instance_owner_id":  accountId,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              acctest.CHECKSET,
						"child_instance_id":        acctest.CHECKSET,
						"child_instance_type":      "VPC",
						"child_instance_region_id": acctest.DefaultRegionToTest,
						"child_instance_owner_id":  acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_owner_id"},
			},
		},
	})
}
func TestAccAliCloudCenInstanceAttachment_multi_same_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}
	ra := acctest.ResourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			//acctest.PreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiSameRegion(rand, acctest.DefaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudCenInstanceAttachment_multi_different_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}
	ra := acctest.ResourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiDifferentRegion(rand, acctest.DefaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

var cenInstanceAttachmentMap = map[string]string{
	"instance_id":              acctest.CHECKSET,
	"child_instance_id":        acctest.CHECKSET,
	"child_instance_region_id": acctest.CHECKSET,
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
			cbnService := cen.NewCbnService(client)
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
	cbnService := cen.NewCbnService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cbnService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
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
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}
	request := cbn.CreateDetachCenChildInstanceRequest()
	request.ChildInstanceId = parts[1]
	request.ChildInstanceRegionId = parts[3]
	request.ChildInstanceType = parts[2]
	request.CenId = parts[0]
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DetachCenChildInstance(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
	"child_instance_owner_id": acctest.CHECKSET,
	"status":                  acctest.CHECKSET,
}
