package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_vpc_traffic_mirror_filter",
		&resource.Sweeper{
			Name: "alicloud_vpc_traffic_mirror_filter",
			F:    testSweepVPCTrafficMirrorFilter,
		})
}

func testSweepVPCTrafficMirrorFilter(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListTrafficMirrorFilters"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeLarge
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.TrafficMirrorFilters", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.TrafficMirrorFilters", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TrafficMirrorFilterName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping VPC Traffic Mirror Filter: %s", item["TrafficMirrorFilterName"].(string))
				continue
			}
			action := "DeleteTrafficMirrorFilter"
			request := map[string]interface{}{
				"TrafficMirrorFilterId": item["TrafficMirrorFilterId"],
			}
			request["RegionId"] = client.RegionId
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete VPC Traffic Mirror Filter (%s): %s", item["TrafficMirrorFilterName"].(string), err)
			}
			log.Printf("[INFO] Delete VPC Traffic Mirror Filter success: %s ", item["TrafficMirrorFilterName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudVPCTrafficMirrorFilter_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_filter.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorFilterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorFilter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccvpctrafficmirrorfilter%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorFilterBasicDependence0)
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
					"traffic_mirror_filter_name":        "${var.name}",
					"traffic_mirror_filter_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_name":        name,
						"traffic_mirror_filter_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_filter_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_filter_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mirror_filter_name":        "${var.name}",
					"traffic_mirror_filter_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_name":        name,
						"traffic_mirror_filter_description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudVPCTrafficMirrorFilter_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_filter.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorFilterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorFilter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccvpctrafficmirrorfilter%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorFilterBasicDependence0)
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
					"traffic_mirror_filter_name":        "${var.name}",
					"traffic_mirror_filter_description": "${var.name}",
					"dry_run":                           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_name":        name,
						"traffic_mirror_filter_description": name,
						"dry_run":                           "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVPCTrafficMirrorFilterMap0 = map[string]string{
	"status":                     CHECKSET,
	"dry_run":                    NOSET,
	"traffic_mirror_filter_name": CHECKSET,
}

func AlicloudVPCTrafficMirrorFilterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
