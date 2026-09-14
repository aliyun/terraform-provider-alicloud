package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccCmsDigitalEmployeesDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ram_role" "default" {
  name     = "tf-acc-cms-de-role-%[1]s"
  document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "RAM": ["acs:ram::1632449562862807:root"]
      }
    }
  ]
}
EOF
}

resource "alicloud_cms_digital_employee" "default" {
  digital_employee_name = "%[1]s"
  role_arn             = alicloud_ram_role.default.arn
  description          = "description-ds"
  display_name         = "display-ds"
  default_rule         = "rule-ds"
  attributes = {
    k1 = "v1"
  }
  sandbox_network_policy {
    allow_fqdns = ["api.example.com"]
    allow_cidrs = ["10.0.0.0/8"]
    enable_acl  = false
  }
  tags {
    key   = "env"
    value = "tf-acc-ds"
  }
}
`, name)
}

func TestAccAliCloudCmsDigitalEmployees_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	resourceId := "alicloud_cms_digital_employee.default"
	dataSourceId := "data.alicloud_cms_digital_employees.default"
	dataSourceConfig := dataSourceTestAccConfigFunc(dataSourceId, name, testAccCmsDigitalEmployeesDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsDigitalEmployeeDestroy(resourceId),
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig(map[string]interface{}{
					"digital_employee_name": "${alicloud_cms_digital_employee.default.digital_employee_name}",
					"ids":                   []string{"${alicloud_cms_digital_employee.default.id}"},
					"output_file":           "digital_employees.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceId, "ids.#"),
					resource.TestCheckResourceAttrSet(dataSourceId, "digital_employees.#"),
					resource.TestCheckResourceAttr(dataSourceId, "digital_employees.0.digital_employee_name", name),
					resource.TestCheckResourceAttr(dataSourceId, "digital_employees.0.description", "description-ds"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsDigitalEmployees_empty(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms-nonexistent-%d", rand)
	dataSourceId := "data.alicloud_cms_digital_employees.default"
	dataSourceConfig := dataSourceTestAccConfigFunc(dataSourceId, name, func(name string) string { return "" })
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig(map[string]interface{}{
					"digital_employee_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "digital_employees.#", "0"),
				),
			},
		},
	})
}
