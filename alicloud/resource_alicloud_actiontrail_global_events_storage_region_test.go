package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudActiontrailGlobalEventsStorageRegion_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ActiontrailGlobalEventsStorageRegionSupportRegions)
	resourceId := "alicloud_actiontrail_global_events_storage_region.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailGlobalEventsStorageRegionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailGlobalEventsStorageRegion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailGlobalEventsStorageRegionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			// Step1: create/update to cn-hangzhou (valid enum). If the singleton
			// residual differs (e.g. ap-southeast-1 from a prior run), this toggles
			// the value and exercises the UpdateGlobalEventsStorageRegion API plus
			// the Read eventual-consistency polling path.
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": "cn-hangzhou",
					}),
				),
			},
			// Step2: update to ap-southeast-1 (the other valid enum). The toggle
			// from cn-hangzhou forces HasChange -> Update API -> Read polling
			// convergence; state settles on ap-southeast-1.
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": "ap-southeast-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": "ap-southeast-1",
					}),
				),
			},
			// Step3: an invalid storage_region must be rejected by ValidateFunc at
			// plan/validate time. state is unchanged (still ap-southeast-1).
			// NB: this ExpectError step MUST NOT be the last step, because the
			// SDK (terraform-plugin-sdk v1.17.2) derives the auto-destroy config
			// from the last step's Config, which would re-trigger ValidateFunc and
			// fail destroy with "config is invalid".
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": "us-east-1",
				}),
				ExpectError: regexp.MustCompile(`expected storage_region to be one of`),
			},
			// Step4: a legal config as the last step. state is already
			// ap-southeast-1, so apply is a no-op, but being the last step means
			// the SDK auto-destroy reuses this (valid) config and passes
			// ValidateFunc. The Check also confirms Step3's failed plan left
			// state untouched.
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": "ap-southeast-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": "ap-southeast-1",
					}),
				),
			},
		},
	})
}

var AlicloudActiontrailGlobalEventsStorageRegionMap0 = map[string]string{}

func AlicloudActiontrailGlobalEventsStorageRegionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

`, name)
}

func TestUnitAlicloudActiontrailGlobalEventsStorageRegionStorageRegionValidation(t *testing.T) {
	schemaMap := resourceAlicloudActiontrailGlobalEventsStorageRegion().Schema
	s, ok := schemaMap["storage_region"]
	if !ok {
		t.Fatal("storage_region schema is missing")
	}
	if s.ValidateFunc == nil {
		t.Fatal("storage_region ValidateFunc is not configured")
	}
	for _, region := range []string{"ap-southeast-1", "cn-hangzhou"} {
		if _, errs := s.ValidateFunc(region, "storage_region"); len(errs) != 0 {
			t.Fatalf("expected storage_region %q to pass validation, got: %v", region, errs)
		}
	}
	for _, region := range []string{"us-east-1", "cn-beijing", "eu-central-1"} {
		if _, errs := s.ValidateFunc(region, "storage_region"); len(errs) == 0 {
			t.Fatalf("expected storage_region %q to be rejected by validation", region)
		}
	}
}
