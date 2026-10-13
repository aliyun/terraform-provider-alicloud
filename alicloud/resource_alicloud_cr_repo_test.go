package alicloud

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlicloudCRRepo_Basic(t *testing.T) {
	var v *cr.GetRepoResponse
	resourceId := "alicloud_cr_repo.default"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${alicloud_cr_namespace.default.name}",
					"name":      "${var.name}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace": name,
						"name":      name,
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
					"detail": "detail",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": "summary update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": "summary update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": "detail update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"detail":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary":   "summary",
						"repo_type": "PUBLIC",
						"detail":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCRRepo_Multi(t *testing.T) {
	var v *cr.GetRepoResponse
	resourceId := "alicloud_cr_repo.default.4"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${alicloud_cr_namespace.default.name}",
					"name":      "${var.name}${count.index}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"count":     "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceCRRepoConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_namespace" "default" {
  name               = "${var.name}"
  auto_create        = false
  default_visibility = "PRIVATE"
}
`, name)
}

var crRepoMap = map[string]string{
	"namespace": CHECKSET,
	"name":      CHECKSET,
	"summary":   "summary",
	"repo_type": "PUBLIC",
}

func TestCRRepoStateUpgradeV0(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected []interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id":        "test-ns/test-repo",
				"namespace": "test-ns",
				"name":      "test-repo",
				"domain_list": map[string]interface{}{
					"vpc":      "registry-vpc.cn-hangzhou.cr.aliyuncs.com",
					"public":   "registry.cn-hangzhou.cr.aliyuncs.com",
					"internal": "registry-internal.cn-hangzhou.cr.aliyuncs.com",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"vpc":      "registry-vpc.cn-hangzhou.cr.aliyuncs.com",
					"public":   "registry.cn-hangzhou.cr.aliyuncs.com",
					"internal": "registry-internal.cn-hangzhou.cr.aliyuncs.com",
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":          "test-ns/test-repo",
				"domain_list": map[string]interface{}{},
			},
			expected: []interface{}{},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":          "test-ns/test-repo",
				"domain_list": nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "test-ns/test-repo",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceAlicloudCRRepoStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := result["domain_list"]
			if tc.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			gotList, ok := got.([]interface{})
			if !ok {
				t.Fatalf("expected []interface{}, got %T", got)
			}
			if len(gotList) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(gotList))
			}
			for i, item := range gotList {
				gotMap := item.(map[string]interface{})
				expMap := tc.expected[i].(map[string]interface{})
				for k, v := range expMap {
					if gotMap[k] != v {
						t.Errorf("item[%d].%s: expected %v, got %v", i, k, v, gotMap[k])
					}
				}
			}
		})
	}
}

func TestCRRepoSchemaVersionV0ToV1(t *testing.T) {
	r := resourceAlicloudCRRepo()
	if r.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion 1, got %d", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 StateUpgrader, got %d", len(r.StateUpgraders))
	}
	if r.StateUpgraders[0].Version != 0 {
		t.Errorf("expected StateUpgrader version 0, got %d", r.StateUpgraders[0].Version)
	}
}

func TestAccAlicloudCRRepo_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_cr_repo.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrRepo")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccCRRepoStateMigrationConfigV0(name),
			},
			{
				ProviderFactories: testAccProviderFactory,
				Config:            testAccCRRepoStateMigrationConfigV1(name),
			},
		},
	})
}

func testAccCRRepoStateMigrationConfigV0(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cr_namespace" "default" {
  name               = "%[1]s"
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_repo" "default" {
  namespace = alicloud_cr_namespace.default.name
  name      = "%[1]s"
  summary   = "state migration test"
  repo_type = "PUBLIC"
}

output "domain_list" {
  value = alicloud_cr_repo.default.domain_list["vpc"]
}
`, name)
}

func testAccCRRepoStateMigrationConfigV1(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cr_namespace" "default" {
  name               = "%[1]s"
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_repo" "default" {
  namespace = alicloud_cr_namespace.default.name
  name      = "%[1]s"
  summary   = "state migration test"
  repo_type = "PUBLIC"
}

output "domain_list" {
  value = alicloud_cr_repo.default.domain_list.0.vpc
}
`, name)
}
