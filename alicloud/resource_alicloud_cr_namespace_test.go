package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_namespace", &resource.Sweeper{
		Name: "alicloud_cr_namespace",
		F:    testSweepCRNamespace,
	})
}

func testSweepCRNamespace(region string) error {
	// skip not supported region
	for _, r := range connectivity.CRNoSupportedRegions {
		if region == string(r) {
			log.Printf("[INFO] testSweepCRNamespace skipped not supported region: %s", region)
			return nil
		}
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := cr.CreateGetNamespaceListRequest()

	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.GetNamespaceList(req)
	})

	if err != nil {
		log.Printf("[ERROR] %s ", WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_namespace", req.GetActionName(), AlibabaCloudSdkGoERROR))
		return nil
	}

	var resp crDescribeNamespaceListResponse
	err = json.Unmarshal(raw.(*cr.GetNamespaceListResponse).GetHttpContentBytes(), &resp)
	if err != nil {
		log.Printf("[ERROR] %s", WrapError(err))
		return nil
	}

	var ns []string
	for _, n := range resp.Data.Namespace {
		for _, p := range prefixes {
			if strings.HasPrefix(n.Namespace, strings.ToLower(p)) {
				ns = append(ns, n.Namespace)
			}
		}
	}

	for _, n := range ns {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			req := cr.CreateDeleteNamespaceRequest()
			req.Namespace = n

			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.DeleteNamespace(req)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
					return nil
				}
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, n, req.GetActionName(), AlibabaCloudSdkGoERROR))
			}

			crService := CrService{client}

			_, err = crService.DescribeCrNamespace(n)
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, n, req.GetActionName(), AlibabaCloudSdkGoERROR))
			}

			time.Sleep(15 * time.Second)
			return resource.RetryableError(WrapError(Error("DeleteNamespace timeout")))
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Namespace: %s", n)
		}
	}
	return nil
}

func TestAccAlicloudCRNamespace_Basic(t *testing.T) {
	var v *cr.GetNamespaceResponse
	resourceId := "alicloud_cr_namespace.default"
	ra := resourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRNamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
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
					"auto_create": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_visibility": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_visibility": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func resourceCRNamespaceConfigDependence(name string) string {
	return ""
}

var crNamespaceMap = map[string]string{
	"name":               CHECKSET,
	"auto_create":        CHECKSET,
	"default_visibility": CHECKSET,
}
