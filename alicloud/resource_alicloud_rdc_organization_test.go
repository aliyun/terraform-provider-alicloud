package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_rdc_organization", &resource.Sweeper{
		Name: "alicloud_rdc_organization",
		F:    testSweepRdcOrganization,
	})
}

func testSweepRdcOrganization(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	conn, err := client.NewDevopsrdcClient()

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListUserOrganization"

	request := map[string]interface{}{}

	var response map[string]interface{}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-03"), StringPointer("AK"), nil, request, &runtime)

	if err != nil {
		log.Printf("[ERROR] %s got an error: %s", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.Object", response)
	if err != nil {
		log.Printf("[ERROR] %s parsing response got an error: %s", action, err)
		return nil
	}

	result, _ := resp.([]interface{})

	for _, v := range result {
		item := v.(map[string]interface{})

		name := item["Name"]
		id := item["Id"]
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping organization: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting organization: %s (%s)", name, id)
		action = "DeleteDevopsOrganization"
		request = map[string]interface{}{
			"OrgId": item["Id"],
		}

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			log.Printf("[ERROR] Failed to delete organization(%s (%s)): %s", name, id, err)
		}
	}

	return nil
}

func TestAccAlicloudRDCOrganization_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rdc_organization.default"
	ra := resourceAttrInit(resourceId, AlicloudRDCOrganizationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DevopsRdcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdcOrganization")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdcorganization%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRDCOrganizationBasicDependence0)
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
					"organization_name": name,
					"source":            "tf-testaccsource",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"organization_name": name,
						"source":            "tf-testaccsource",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"real_pk", "source", "desired_member_count"},
			},
		},
	})
}

var AlicloudRDCOrganizationMap0 = map[string]string{
	"real_pk":              NOSET,
	"source":               NOSET,
	"desired_member_count": NOSET,
}

func AlicloudRDCOrganizationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
