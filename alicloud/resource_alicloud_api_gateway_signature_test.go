package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_signature", &resource.Sweeper{
		Name: "alicloud_api_gateway_signature",
		F:    testSweepApiGatewaySignature,
	})
}

func testSweepApiGatewaySignature(region string) error {
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

	req := cloudapi.CreateDescribeSignaturesRequest()
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	sweeped := false

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeSignatures(req)
		})
		if err != nil {
			log.Printf("[ERROR] Describe Signatures: %s", err)
			return nil
		}
		response, _ := raw.(*cloudapi.DescribeSignaturesResponse)

		for _, v := range response.SignatureInfos.SignatureInfo {
			name := v.SignatureName
			id := v.SignatureId
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping api gateway signature: %s", name)
				continue
			}
			sweeped = true

			log.Printf("[INFO] Deleting Api Gateway Signature: %s", name)

			req := cloudapi.CreateDeleteSignatureRequest()
			req.SignatureId = id
			_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
				return cloudApiClient.DeleteSignature(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Api Gateway Signature (%s): %s", name, err)
			}
		}

		if len(response.SignatureInfos.SignatureInfo) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		req.PageNumber = page
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAliCloudApiGatewaySignature(t *testing.T) {
	var v *cloudapi.SignatureInfo

	resourceId := "alicloud_api_gateway_signature.default"
	ra := resourceAttrInit(resourceId, apigatewaySignatureBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccSignature_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewaySignatureConfigDependence)

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
					"signature_name":   "${var.name}",
					"signature_key":    "${var.key}",
					"signature_secret": "${var.secret}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"signature_name":   name,
						"signature_key":    "tf_testAccKey",
						"signature_secret": "tf_testAccSecret",
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
					"signature_name":   "${var.name}_u",
					"signature_key":    "${var.key}",
					"signature_secret": "${var.secret}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"signature_name":   name + "_u",
						"signature_key":    "tf_testAccKey",
						"signature_secret": "tf_testAccSecret",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"signature_name":   "${var.name}",
					"signature_key":    "${var.key}_u",
					"signature_secret": "${var.secret}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"signature_name":   name,
						"signature_key":    "tf_testAccKey_u",
						"signature_secret": "tf_testAccSecret_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"signature_name":   "${var.name}",
					"signature_key":    "${var.key}",
					"signature_secret": "${var.secret}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"signature_name":   name,
						"signature_key":    "tf_testAccKey",
						"signature_secret": "tf_testAccSecret",
					}),
				),
			},
		},
	})
}

var apigatewaySignatureBasicMap = map[string]string{
	"signature_name":   CHECKSET,
	"signature_key":    CHECKSET,
	"signature_secret": CHECKSET,
}

func resourceApigatewaySignatureConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "key" {
  default = "tf_testAccKey"
}
variable "secret" {
  default = "tf_testAccSecret"
}
`, name)
}
