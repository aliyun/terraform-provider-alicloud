package fcv2_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcv2"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testFunctionOnDemandConfigCodeZip is a base64 zip containing a minimal
// index.py with `def handler(event, context)`, enough for a python3.9 FCv2
// function that the on-demand config attaches to.
const testFunctionOnDemandConfigCodeZip = "UEsDBBQAAAAIABAfZl27wpNPTgAAAFMAAAAIABwAaW5kZXgucHlVVAkAA2/f7Gpv3+xqdXgLAAEE6AMAAAToAwAAHcpRCoAgDADQf08x/FKIDtBtIhcWc5M5K4juXvS+31aqqMHehF3CFfLMiVADHsg2wCJseFmcHHwUrSv/eUy91BZun5FI/AT+FKXkn+heUEsBAh4DFAAAAAgAEB9mXbvCk09OAAAAUwAAAAgAGAAAAAAAAQAAAKSBAAAAAGluZGV4LnB5VVQFAANv3+xqdXgLAAEE6AMAAAToAwAAUEsFBgAAAAABAAEATgAAAJAAAAAAAA=="

func TestAccAliCloudFcv2FunctionOnDemandConfig_basic(t *testing.T) {
	resourceId := "alicloud_fcv2_function_on_demand_config.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2ondemand%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckFcv2FunctionOnDemandConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudFCV2FunctionOnDemandConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "resource"),
					resource.TestCheckResourceAttr(resourceId, "service_name", name),
					resource.TestCheckResourceAttr(resourceId, "function_name", name),
					resource.TestCheckResourceAttr(resourceId, "qualifier", "LATEST"),
					resource.TestCheckResourceAttr(resourceId, "maximum_instance_count", "10"),
				),
			},
			{
				Config: testAccAliCloudFCV2FunctionOnDemandConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "maximum_instance_count", "20"),
					resource.TestCheckResourceAttr(resourceId, "qualifier", "LATEST"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFcv2FunctionOnDemandConfigDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_fcv2_function_on_demand_config" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := fcv2.NewService(client).DescribeFunctionOnDemandConfig(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing function on-demand config %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("function on-demand config %s still exists", rs.Primary.ID)
	}
	return nil
}

// The on-demand config attaches to a published or LATEST qualifier of an
// existing FCv2 function, so the dependency is a service + function. The
// function_name is shared by the function and the config, keeping the two in
// sync without depending on attributes the function may not export.
func testAccAliCloudFCV2FunctionOnDemandConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf-testacc for fcv2 function on-demand config"
}

resource "alicloud_fcv2_function" "default" {
  service_name  = alicloud_fc_service.default.name
  function_name = var.name
  runtime       = "python3.9"
  handler       = "index.handler"
  code {
    zip_file = "%[2]s"
  }
}

resource "alicloud_fcv2_function_on_demand_config" "default" {
  service_name           = alicloud_fc_service.default.name
  function_name          = var.name
  qualifier              = "LATEST"
  maximum_instance_count = 10
}
`, name, testFunctionOnDemandConfigCodeZip)
}

func testAccAliCloudFCV2FunctionOnDemandConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf-testacc for fcv2 function on-demand config"
}

resource "alicloud_fcv2_function" "default" {
  service_name  = alicloud_fc_service.default.name
  function_name = var.name
  runtime       = "python3.9"
  handler       = "index.handler"
  code {
    zip_file = "%[2]s"
  }
}

resource "alicloud_fcv2_function_on_demand_config" "default" {
  service_name           = alicloud_fc_service.default.name
  function_name          = var.name
  qualifier              = "LATEST"
  maximum_instance_count = 20
}
`, name, testFunctionOnDemandConfigCodeZip)
}
