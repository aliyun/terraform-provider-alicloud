package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudALBSecurityPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_security_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudALBSecurityPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbSecurityPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbsecuritypolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBSecurityPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "tf-testAcc-test-secrity",
					"tls_versions":         []string{"TLSv1.0"},
					"ciphers":              []string{"ECDHE-ECDSA-AES128-SHA", "AES256-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity",
						"ciphers.#":            "2",
						"tls_versions.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "tf-testAcc-test-secrity-new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity-new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{"TLSv1.1"},
					"ciphers":      []string{"AES128-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#":      "1",
						"tls_versions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc4",
						"For":     "Tftestacc4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc4",
						"tags.For":     "Tftestacc4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc42",
						"For":     "Tftestacc42",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc42",
						"tags.For":     "Tftestacc42",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "tf-testAcc-test-secrity",
					"ciphers":              []string{"TLS_AES_128_GCM_SHA256"},
					"tls_versions":         []string{"TLSv1.3"},
					"tags": map[string]string{
						"Created": "tfTestAcc58",
						"For":     "Tftestacc58",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity",
						"ciphers.#":            "1",
						"tls_versions.#":       "1",
						"tags.%":               "2",
						"tags.Created":         "tfTestAcc58",
						"tags.For":             "Tftestacc58",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBSecurityPolicyMap0 = map[string]string{
	"dry_run":              NOSET,
	"ciphers.#":            CHECKSET,
	"status":               CHECKSET,
	"tls_versions.#":       CHECKSET,
	"resource_group_id":    CHECKSET,
	"tags.%":               NOSET,
	"security_policy_name": "tf-testAcc-test-secrity",
}

func AlicloudALBSecurityPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
