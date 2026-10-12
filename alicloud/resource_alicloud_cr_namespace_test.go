package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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

// Test Cr Namespace. >>> Resource test cases, automatically generated.
// Case Namespace 生命周期 5730
func TestAccAliCloudCrNamespace_basic5730(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudCrNamespaceMap5730)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrNamespaceBasicDependence5730)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":    name,
					"instance_id":       "${alicloud_cr_ee_instance.defaultTwqdy8.id}",
					"auto_create_repo":  "false",
					"default_repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":    name,
						"instance_id":       CHECKSET,
						"auto_create_repo":  "false",
						"default_repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_create_repo":  "true",
					"default_repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create_repo":  "true",
						"default_repo_type": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCrNamespaceMap5730 = map[string]string{
	"status":       CHECKSET,
	"namespace_id": CHECKSET,
}

func AlicloudCrNamespaceBasicDependence5730(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "defaultTwqdy8" {
  instance_name  = "auto-test-731"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}


`, name)
}

// Case Namespace 生命周期_succ 5850
func TestAccAliCloudCrNamespace_basic5850(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudCrNamespaceMap5850)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrNamespaceBasicDependence5850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":    name,
					"instance_id":       "${alicloud_cr_ee_instance.defaultTwqdy8.id}",
					"auto_create_repo":  "false",
					"default_repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":    name,
						"instance_id":       CHECKSET,
						"auto_create_repo":  "false",
						"default_repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_create_repo":  "true",
					"default_repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create_repo":  "true",
						"default_repo_type": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCrNamespaceMap5850 = map[string]string{
	"status":       CHECKSET,
	"namespace_id": CHECKSET,
}

func AlicloudCrNamespaceBasicDependence5850(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "defaultTwqdy8" {
  instance_name  = "auto-test-instance"
  renewal_status = "ManualRenewal"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}


`, name)
}

// Case Namespace 生命周期_副本1706083189606 5840
func TestAccAliCloudCrNamespace_basic5840(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudCrNamespaceMap5840)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrNamespaceBasicDependence5840)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":    name,
					"instance_id":       "${alicloud_cr_ee_instance.defaultTwqdy8.id}",
					"auto_create_repo":  "false",
					"default_repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":    name,
						"instance_id":       CHECKSET,
						"auto_create_repo":  "false",
						"default_repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_create_repo":  "true",
					"default_repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create_repo":  "true",
						"default_repo_type": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCrNamespaceMap5840 = map[string]string{
	"status":       CHECKSET,
	"namespace_id": CHECKSET,
}

func AlicloudCrNamespaceBasicDependence5840(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "defaultTwqdy8" {
  default_oss_bucket = "true"
  renewal_status     = "ManualRenewal"
  period             = "1"
  instance_name      = "api-test-345"
  image_scanner      = "ACR"
  payment_type       = "Subscription"
  instance_type      = "Basic"
}


`, name)
}

// Test Cr Namespace. <<< Resource test cases, automatically generated.
