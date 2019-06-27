package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_app", &resource.Sweeper{
		Name: "alicloud_api_gateway_app",
		F:    testSweepApiGatewayApp,
	})
}

func testSweepApiGatewayApp(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeAppAttributesRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeAppAttributes(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Apps: %s", err)
	}
	apps, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

	swept := false

	for _, v := range apps.Apps.AppAttribute {
		name := v.AppName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping app: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting App: %s", name)

		req := cloudapi.CreateDeleteAppRequest()
		req.AppId = requests.NewInteger(v.AppId)
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApp(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete App (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func SkipTestAccAlicloudApigatewayApp_basic(t *testing.T) {
	var v *cloudapi.DescribeAppResponse

	resourceId := "alicloud_api_gateway_app.default"
	ra := resourceAttrInit(resourceId, apigatewayAppBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApp_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayAppConfigDependence)

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
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"description": "${var.description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf_testAcc api gateway description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.description}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf_testAcc api gateway description_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "${var.description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tf_testAcc api gateway description",
					}),
				),
			},
		},
	})
}

func SkipTestAccAlicloudApigatewayApp_multi(t *testing.T) {
	var v *cloudapi.DescribeAppResponse
	resourceId := "alicloud_api_gateway_app.default.9"
	ra := resourceAttrInit(resourceId, apigatewayAppBasicMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApp_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayAppConfigDependence)

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
					"name":        "${var.name}",
					"description": "${var.description}",
					"count":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}

func resourceApigatewayAppConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

`, name)
}

var apigatewayAppBasicMap = map[string]string{
	"name": CHECKSET,
}
