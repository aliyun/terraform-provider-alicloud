package kms_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test Kms ValueAddedService. >>> Resource test cases, automatically generated.
// Case 默认密钥增值服务 11636
func TestAccAliCloudKmsValueAddedService_basic11636(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_value_added_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudKmsValueAddedServiceMap11636)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsValueAddedService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudKmsValueAddedServiceBasicDependence11636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"value_added_service": "2",
					"payment_type":        "Subscription",
					"period":              "1",
					"renew_period":        "1",
					"renew_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value_added_service": acctest.CHECKSET,
						"payment_type":        "Subscription",
						"period":              "1",
						"renew_period":        "1",
						"renew_status":        "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_status": "ManualRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "value_added_service"},
			},
		},
	})
}

var AlicloudKmsValueAddedServiceMap11636 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudKmsValueAddedServiceBasicDependence11636(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}`, name)
}

// Test Kms ValueAddedService. <<< Resource test cases, automatically generated.

// TestAccAliCloudKmsValueAddedService_duplicateOrder reproduces the recreate-on-every-plan report. The
// loop needs the account to already hold the service, so it cannot be reached from a single resource on
// a clean account: it takes a second resource with the same config, which is what a duplicate order
// actually looks like.
//
// Both resources are declared in one config with no dependency between them, so Terraform creates them
// concurrently and either guard may be the one that fires - the pre-flight, if the first order is
// already effective by the time the second one checks, or the create wait, if both checks run before
// either order lands and the loser's order is refunded. Whichever it is, the create has to fail with
// the shared already-held message instead of writing the refunded instance into state, which is what
// used to start the loop. Asserting on the message the two guards share is what keeps that stable.
func TestAccAliCloudKmsValueAddedService_duplicateOrder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckKmsValueAddedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccKmsValueAddedServiceDuplicateConfig(),
				ExpectError: regexp.MustCompile("This account may already hold this value added service"),
			},
		},
	})
}

func testAccKmsValueAddedServiceConfig() string {
	return `
resource "alicloud_kms_value_added_service" "default" {
  value_added_service = "2"
  payment_type        = "Subscription"
  period              = "1"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}
`
}

func testAccKmsValueAddedServiceDuplicateConfig() string {
	return testAccKmsValueAddedServiceConfig() + `
resource "alicloud_kms_value_added_service" "duplicate" {
  value_added_service = "2"
  payment_type        = "Subscription"
  period              = "1"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}
`
}

// TestUnitKmsValueAddedServiceInEffect pins which listed instances count as holding the service. The
// cases that matter are the ones that must NOT count: refusing a create over any of them would leave
// the user with no way to proceed, while letting a duplicate order through only costs an order that
// BssOpenApi refunds by itself and that the create wait already catches.
// TestAccAliCloudKmsValueAddedService_alreadyHeld covers the already-held pre-flight. The second step
// adds a second resource for the same value added service in the same region, and because steps run
// sequentially the first instance is already effective when the second create starts, so the prefix
// pre-flight has to recognise it and refuse instead of placing an order that would only be refunded.
func TestAccAliCloudKmsValueAddedService_alreadyHeld(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckKmsValueAddedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsValueAddedServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsValueAddedServiceExists("alicloud_kms_value_added_service.default"),
					resource.TestCheckResourceAttr("alicloud_kms_value_added_service.default", "status", "Normal"),
					resource.TestCheckResourceAttrSet("alicloud_kms_value_added_service.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_kms_value_added_service.default", "region_id"),
				),
			},
			{
				Config:      testAccKmsValueAddedServiceAlreadyHeldConfig(),
				ExpectError: regexp.MustCompile("already holds this value added service"),
			},
		},
	})
}

func testAccKmsValueAddedServiceAlreadyHeldConfig() string {
	return testAccKmsValueAddedServiceConfig() + `
resource "alicloud_kms_value_added_service" "duplicate" {
  depends_on          = [alicloud_kms_value_added_service.default]
  value_added_service = "2"
  payment_type        = "Subscription"
  period              = "1"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}
`
}

// TestAccAliCloudKmsValueAddedService_twoRegions covers the confirmed per-region uniqueness of the
// service: two providers pinned to different regions each order the same value added service (code 2)
// and both succeed independently, because the service is scoped per region. That also pins the
// already-held pre-flight to the right scope - it filters QueryAvailableInstances on the response
// Region, so the cn-hangzhou instance must not make the cn-shanghai create refuse.
//
// Two same-region resources in one apply are deliberately not exercised: Terraform creates them
// concurrently, so both pre-flights can run before either order lands and the race is not fixable
// from inside the resource.
func TestAccAliCloudKmsValueAddedService_twoRegions(t *testing.T) {
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
		"alicloudshanghai": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou", "cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckKmsValueAddedServiceDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsValueAddedServiceTwoRegionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsValueAddedServiceExistsWithProviders("alicloud_kms_value_added_service.hz", &providers),
					testAccCheckKmsValueAddedServiceExistsWithProviders("alicloud_kms_value_added_service.sh", &providers),
					resource.TestCheckResourceAttr("alicloud_kms_value_added_service.hz", "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_kms_value_added_service.sh", "region_id", "cn-shanghai"),
				),
			},
		},
	})
}

func testAccKmsValueAddedServiceTwoRegionsConfig() string {
	return `
provider "alicloud" {
  region = "cn-hangzhou"
}

provider "alicloudshanghai" {
  alias  = "sh"
  region = "cn-shanghai"
}

resource "alicloud_kms_value_added_service" "hz" {
  value_added_service = "2"
  payment_type        = "Subscription"
  period              = "1"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}

resource "alicloud_kms_value_added_service" "sh" {
  provider            = alicloudshanghai.sh
  value_added_service = "2"
  payment_type        = "Subscription"
  period              = "1"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}
`
}

func testAccCheckKmsValueAddedServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS value added service ID is set")
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		kmsServiceV2 := kms.NewKmsServiceV2(client)
		if _, err := kmsServiceV2.DescribeKmsValueAddedService(rs.Primary.ID); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckKmsValueAddedServiceDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	kmsServiceV2 := kms.NewKmsServiceV2(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kms_value_added_service" {
			continue
		}
		_, err := kmsServiceV2.DescribeKmsValueAddedService(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("KMS value added service %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCheckKmsValueAddedServiceExistsWithProviders(n string, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS value added service ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers.
			if provider.Meta() == nil {
				continue
			}
			client := provider.Meta().(*connectivity.AliyunClient)
			kmsServiceV2 := kms.NewKmsServiceV2(client)
			if _, err := kmsServiceV2.DescribeKmsValueAddedService(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return nil
		}
		return fmt.Errorf("KMS value added service %s not found in any provider", rs.Primary.ID)
	}
}

func testAccCheckKmsValueAddedServiceDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		checked := false
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			checked = true
			client := provider.Meta().(*connectivity.AliyunClient)
			kmsServiceV2 := kms.NewKmsServiceV2(client)
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "alicloud_kms_value_added_service" {
					continue
				}
				_, err := kmsServiceV2.DescribeKmsValueAddedService(rs.Primary.ID)
				if err != nil {
					if errs.NotFoundError(err) {
						continue
					}
					return err
				}
				return fmt.Errorf("KMS value added service %s still exists", rs.Primary.ID)
			}
		}
		if !checked {
			return fmt.Errorf("no configured provider to check KMS value added service destroy with")
		}
		return nil
	}
}
