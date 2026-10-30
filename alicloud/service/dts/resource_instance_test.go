package dts_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dtsInstanceResourceId = "alicloud_dts_instance.default"

// dtsRenderAttrs mirrors the legacy SDKv2 resourceTestAccConfigFunc builder: each
// entry in attrs is rendered as a `key = value` HCL line. Strings are quoted
// (preserving ${...} interpolation and JSON payloads), bools and ints are bare,
// and maps (used for tags) render as HCL blocks. Map key order is sorted for
// deterministic output; Terraform plan/apply is order-insensitive, so this
// does not change behavior.
func dtsRenderAttrs(attrs map[string]interface{}) string {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "  %s = %s\n", k, dtsRenderValue(1, attrs[k]))
	}
	return b.String()
}

// dtsRenderValue renders a single attribute value as HCL. indent is the
// indentation level of the enclosing attribute (1 for top-level resource
// attributes), used to indent nested map keys and the closing brace.
func dtsRenderValue(indent int, v interface{}) string {
	if v == nil {
		return "null"
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		// Mirror the legacy valueConvert: wrap the raw string in double quotes
		// without escaping. Test maps carry pre-escaped HCL payloads (db_list,
		// job_parameters, reserve JSON) where \" already denotes an embedded
		// quote; %q would double the backslashes and break the HCL.
		return fmt.Sprintf("\"%s\"", rv.String())
	case reflect.Bool:
		return fmt.Sprintf("%t", rv.Bool())
	case reflect.Int, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Map:
		mk := make([]string, 0, rv.Len())
		for _, k := range rv.MapKeys() {
			mk = append(mk, fmt.Sprintf("%v", k.Interface()))
		}
		sort.Strings(mk)
		var lines []string
		for _, k := range mk {
			lines = append(lines, fmt.Sprintf("%s%s = %s",
				strings.Repeat("  ", indent+1), k,
				dtsRenderValue(indent+1, rv.MapIndex(reflect.ValueOf(k)).Interface())))
		}
		return fmt.Sprintf("{\n%s\n%s}", strings.Join(lines, "\n"), strings.Repeat("  ", indent))
	default:
		return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
	}
}

func TestAccAliCloudDTSInstance_basic1170(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDtsInstance%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"type":                             "migration",
		"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
		"payment_type":                     "PayAsYouGo",
		"instance_class":                   "large",
		"source_endpoint_engine_name":      "MySQL",
		"source_region":                    "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_engine_name": "MySQL",
		"compute_unit":                     "2",
		"database_count":                   "1",
		"destination_region":               "${data.alicloud_regions.default.regions.0.id}",
		"sync_architecture":                "oneway",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite gates (checkoutSupportedRegions + testAccPreCheckWithRegions)
			// converged to cn-hangzhou for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkDtsInstanceDestroy(dtsInstanceResourceId),
		Steps: func() []resource.TestStep {
			steps := []resource.TestStep{
				{
					Config: dtsInstanceConfig(name, attrs),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "type", "migration"),
						resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "resource_group_id"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "payment_type", "PayAsYouGo"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "instance_class", "large"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "source_region"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "compute_unit", "2"),
						resource.TestCheckResourceAttr(dtsInstanceResourceId, "database_count", "1"),
						resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "destination_region"),
					),
				},
			}
			attrs["resource_group_id"] = "${data.alicloud_resource_manager_resource_groups.default.ids.1}"
			steps = append(steps, resource.TestStep{
				Config: dtsInstanceConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "resource_group_id"),
				),
			})
			attrs["tags"] = map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			}
			steps = append(steps, resource.TestStep{
				Config: dtsInstanceConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "tags.For", "acceptance test"),
				),
			})
			steps = append(steps, resource.TestStep{
				ResourceName:            dtsInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"compute_unit", "database_count", "used_time", "sync_architecture"},
			})
			return steps
		}(),
	})
}

func TestAccAliCloudDTSInstance_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDtsInstance%d", dtsTestRegion(), rand)

	attrs := map[string]interface{}{
		"type":                             "sync",
		"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
		"payment_type":                     "Subscription",
		"instance_class":                   "large",
		"source_endpoint_engine_name":      "MySQL",
		"source_region":                    "${data.alicloud_regions.default.regions.0.id}",
		"destination_endpoint_engine_name": "MySQL",
		"compute_unit":                     "2",
		"database_count":                   "1",
		"destination_region":               "${data.alicloud_regions.default.regions.0.id}",
		"auto_start":                       "true",
		"du":                               "30",
		"period":                           "Month",
		"sync_architecture":                "bidirectional",
		"used_time":                        "1",
		"fee_type":                         "ONLY_CONFIGURATION_FEE",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite testAccPreCheckWithTime helper has no plugin-testing
			// equivalent and is not used by the migrated templates; drop it.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy stays nil to match the legacy basic1 test, which relied on
		// DescribeDtsInstance returning the resource after deletion (an API bug
		// documented in the original). The migration_instance template's
		// DescribeMigrationInstance-based destroy check is therefore not adopted.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: dtsInstanceConfig(name, attrs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "type", "sync"),
					resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "resource_group_id"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "payment_type", "Subscription"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "instance_class", "large"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "source_region"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "compute_unit", "2"),
					resource.TestCheckResourceAttr(dtsInstanceResourceId, "database_count", "1"),
					resource.TestCheckResourceAttrSet(dtsInstanceResourceId, "destination_region"),
				),
			},
			{
				ResourceName:            dtsInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"compute_unit", "database_count", "used_time", "auto_start", "sync_architecture", "period", "du", "fee_type"},
			},
		},
	})
}

func checkDtsInstanceDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_instance" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeDtsInstance(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("the resource %s %s was not destroyed", resourceId, rs.Primary.ID)
		}
		return nil
	}
}

func dtsInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

func dtsInstanceConfig(name string, attrs map[string]interface{}) string {
	return dtsInstanceDependence(name) + fmt.Sprintf(`
resource "alicloud_dts_instance" "default" {
%s}
`, dtsRenderAttrs(attrs))
}
