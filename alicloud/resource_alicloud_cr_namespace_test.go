package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_namespace", &resource.Sweeper{
		Name: "alicloud_cr_namespace",
		F:    testSweepCRNamespace,
	})
}

func testSweepCRNamespace(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
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
				if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
					return nil
				}
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, n, req.GetActionName(), AlibabaCloudSdkGoERROR))
			}

			crService := CrService{client}

			_, err = crService.DescribeNamespace(n)
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_namespace.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRNamespace_Basic(acctest.RandIntRange(100000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-testacc-cr-ns-basic*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "false"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PUBLIC"),
				),
			},
		},
	})
}

func TestAccAlicloudCRNamespace_Update(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "alicloud_cr_namespace.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRNamespace_UpdateBefore(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-testacc-cr-ns*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "false"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PUBLIC"),
				),
			},
			{
				Config: testAccCRNamespace_UpdateAfter(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cr_namespace.default", "id"),
					resource.TestMatchResourceAttr("alicloud_cr_namespace.default", "name", regexp.MustCompile("tf-testacc-cr-ns*")),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "auto_create", "true"),
					resource.TestCheckResourceAttr("alicloud_cr_namespace.default", "default_visibility", "PRIVATE"),
				),
			},
		},
	})
}

func testAccCheckCRNamespaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_namespace" {
			continue
		}

		crService := CrService{client}
		raw, err := crService.DescribeNamespace(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		var resp crDescribeNamespaceResponse
		err = json.Unmarshal(raw.GetHttpContentBytes(), &resp)
		if err != nil {
			return WrapError(err)
		}

		if resp.Data.Namespace.Namespace != "" {
			return fmt.Errorf("error namespace %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCRNamespace_Basic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-ns-basic-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}
`, rand)
}

func testAccCRNamespace_UpdateBefore(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-ns-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PUBLIC"
}
`, rand)
}

func testAccCRNamespace_UpdateAfter(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-cr-ns-%d"
}

resource "alicloud_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}
`, rand)
}
