package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_access_group", &resource.Sweeper{
		Name: "alicloud_nas_access_group",
		F:    testSweepNasAccessGroup,
		Dependencies: []string{
			"alicloud_nas_file_system",
		},
	})
}

func testSweepNasAccessGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeAccessGroups"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	ids := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			log.Println("[ERROR] List nas access groups failed. err:", err)
		}
		resp, err := jsonpath.Get("$.AccessGroups.AccessGroup", response)
		if err != nil {
			log.Println("Get $.AccessGroups.AccessGroup failed. err:", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["AccessGroupName"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(name, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Nas Access Group: %s ", name)
				continue
			}
			ids = append(ids, fmt.Sprint(item["AccessGroupName"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, id := range ids {
		log.Printf("[Info] Delete Nas Access Group: %s", id)
		action := "DeleteAccessGroup"
		request := map[string]interface{}{
			"AccessGroupName": id,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete  Nas Access Group (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAlicloudNasAccessGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_group.default"
	ra := resourceAttrInit(resourceId, AlicloudNasAccessGroup0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasAccessGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasAccessGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasAccessGroupBasicDependence0)
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
					"access_group_name": "${var.name}",
					"access_group_type": "Vpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"access_group_type": "Vpc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
					}),
				),
			},
		},
	})
}

var AlicloudNasAccessGroup0 = map[string]string{
	"file_system_type": CHECKSET,
}

func AlicloudNasAccessGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
