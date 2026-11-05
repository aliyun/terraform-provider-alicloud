package cr_test

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	crsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := crsdk.CreateGetNamespaceListRequest()

	raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
		return crClient.GetNamespaceList(req)
	})

	if err != nil {
		log.Printf("[ERROR] %s ", errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cr_namespace", req.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		return nil
	}

	var resp crsvc.CrDescribeNamespaceListResponse
	err = json.Unmarshal(raw.(*crsdk.GetNamespaceListResponse).GetHttpContentBytes(), &resp)
	if err != nil {
		log.Printf("[ERROR] %s", errs.WrapError(err))
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
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			req := crsdk.CreateDeleteNamespaceRequest()
			req.Namespace = n

			_, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
				return crClient.DeleteNamespace(req)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
					return nil
				}
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, n, req.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}

			crService := crsvc.NewCrService(client)

			_, err = crService.DescribeCrNamespace(n)
			if err != nil {
				if errs.NotFoundError(err) {
					return nil
				}
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, n, req.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}

			time.Sleep(15 * time.Second)
			return retry.RetryableError(errs.WrapError(errs.Error("DeleteNamespace timeout")))
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Namespace: %s", n)
		}
	}
	return nil
}

func TestAccAliCloudCRNamespace_Basic(t *testing.T) {
	var v *crsdk.GetNamespaceResponse
	resourceId := "alicloud_cr_namespace.default"
	ra := acctest.ResourceAttrInit(resourceId, crNamespaceMap)
	serviceFunc := func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCRNamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"name":               acctest.CHECKSET,
	"auto_create":        acctest.CHECKSET,
	"default_visibility": acctest.CHECKSET,
}
