package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogMachineGroup_basic(t *testing.T) {
	var v *sls.MachineGroup
	resourceId := "alicloud_log_machine_group.default"
	ra := resourceAttrInit(resourceId, logMachineGroupMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogmachinegroupip-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogMachineGroupConfigDependence)

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
					"name":          name,
					"project":       "${alicloud_log_project.default.name}",
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name,
						"project": name,
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
					"identify_type": "userdefined",
					"identify_list": []string{"terraform", "abc1234"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"identify_type":   "userdefined",
						"identify_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"topic": "terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic": "terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"identify_type": REMOVEKEY,
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
					"topic":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"identify_type":   "ip",
						"identify_list.#": "3",
						"topic":           REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogMachineGroup_multi(t *testing.T) {
	var v *sls.MachineGroup
	resourceId := "alicloud_log_machine_group.default.4"
	ra := resourceAttrInit(resourceId, logMachineGroupMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogmachinegroupip-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogMachineGroupConfigDependence)

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
					"name":          name + "${count.index}",
					"project":       "${alicloud_log_project.default.name}",
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
					"count":         "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogMachineGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

var logMachineGroupMap = map[string]string{
	"name":            CHECKSET,
	"project":         CHECKSET,
	"identify_list.#": "3",
}

// TestAccAlicloudLogMachineGroup_signVersionV4 is the path-B representative
// for v4 signature regression testing. It drives a resource whose CRUD is
// entirely on the legacy WithLogClient path (aliyun-log-go-sdk/*sls.Client)
// under a provider configured with sign_version.sls = "v4", exercising the
// applyLogClientSignVersion wiring added in alicloud/connectivity/client.go.
//
// Coverage map — two acceptance tests validate v4 signing across all 22
// alicloud_sls_* / alicloud_log_* resources:
//
//	Path A — SDK v2 client.Do("Sls", ...) → applyOpenapiSignVersion:
//	  alicloud_sls_alert, alicloud_sls_collection_policy, alicloud_sls_etl,
//	  alicloud_sls_index, alicloud_sls_logtail_config,
//	  alicloud_sls_logtail_pipeline_config, alicloud_sls_machine_group,
//	  alicloud_sls_oss_export_sink, alicloud_sls_scheduled_sql,
//	  alicloud_log_project (Create / main Read / Update / Delete)
//	  → represented by TestAccAliCloudSlsProject_signVersionV4.
//
//	Path B — WithLogClient (aliyun-log-go-sdk) → applyLogClientSignVersion:
//	  alicloud_log_alert, alicloud_log_alert_resource, alicloud_log_audit,
//	  alicloud_log_dashboard, alicloud_log_etl, alicloud_log_ingestion,
//	  alicloud_log_machine_group, alicloud_log_oss_export,
//	  alicloud_log_oss_shipper, alicloud_log_resource,
//	  alicloud_log_resource_record, alicloud_log_store,
//	  alicloud_log_store_index, alicloud_log_project (policy sub-path)
//	  → represented by this test.
//
//	Path C — WithLogPopClient (alibaba-cloud-sdk-go/services/sls, POP API):
//	  alicloud_log_audit, alicloud_log_alert_resource (POP calls only).
//	  Not covered: the pinned alibaba-cloud-sdk-go v1.62.590 has no v4
//	  signer. Intentional gap; tracked separately.
//
// Helper-level correctness is further locked down by
// TestApplyOpenapiSignVersion and TestApplyLogClientSignVersion in
// alicloud/connectivity/sign_version_test.go. Together the three layers
// (helper unit test + one acceptance per path) catch any regression that
// would drop v4 signing for any resource on either path.
//
// The endpoint is hard-coded to cn-hangzhou-acdr-ut-3.log.aliyuncs.com
// because only a v4-only endpoint is a meaningful environment for this
// test — see the sibling TestAccAliCloudSlsProject_signVersionV4 doc
// comment for why public endpoints are unsuitable. Running this test
// requires ALICLOUD_ACCESS_KEY whose account has access to
// cn-hangzhou-acdr-ut-3.
func TestAccAlicloudLogMachineGroup_signVersionV4(t *testing.T) {
	var v *sls.MachineGroup
	resourceId := "alicloud_log_machine_group.default"
	ra := resourceAttrInit(resourceId, logMachineGroupMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogmgv4-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogMachineGroupSignVersionV4Dependence)

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
					"name":          name,
					"project":       "${alicloud_log_project.default.name}",
					"identify_list": []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name,
						"project": name,
					}),
				),
			},
			// ImportState intentionally omitted: terraform-plugin-sdk v1's
			// EvalImportStateVerify builds a fresh eval context for the
			// post-import refresh and inline `provider "alicloud" {...}`
			// blocks do not reliably propagate custom endpoints /
			// sign_version into it. See the matching note on
			// TestAccAliCloudSlsProject_signVersionV4.
		},
	})
}

// resourceLogMachineGroupSignVersionV4Dependence inlines a provider block
// targeting a v4-only SLS endpoint (cn-hangzhou-acdr-ut-3) so the test
// exercises the full schema → config.SignVersion → WithLogClient →
// applyLogClientSignVersion chain.
//
// The endpoint, region and sign_version are all hard-coded and must stay
// consistent with each other — see AlicloudSlsProjectSignVersionV4Dependence
// for the rationale on why the credential-scope region is derived from the
// endpoint and why the provider.region pins the account-level region the
// acdr-ut endpoint belongs to.
func resourceLogMachineGroupSignVersionV4Dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  skip_region_validation = true
  region = "cn-hangzhou-acdr-ut-3"
  endpoints {
    log = "cn-hangzhou-acdr-ut-3.log.aliyuncs.com"
  }
  sign_version {
    sls = "v4"
  }
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "tf acc test for log_machine_group sign_version v4"
}
`, name)
}
