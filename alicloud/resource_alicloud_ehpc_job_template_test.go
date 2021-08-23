package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEhpcJobTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_job_template.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcJobTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcJobTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpcjobtemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcJobTemplateBasicDependence0)
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
					"job_template_name": "JobTemplateNameT",
					"command_line":      "./LammpsTest/lammps.pbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_template_name": "JobTemplateNameT",
						"command_line":      "./LammpsTest/lammps.pbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stderr_redirect_path": "./LammpsTest1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stderr_redirect_path": "./LammpsTest1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"clock_time": "12:00:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"clock_time": "12:00:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runas_user": "user1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runas_user": "user1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"thread": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"thread": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_template_name": "JobTemplateNameH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_template_name": "JobTemplateNameH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"variables": "[{Name:,Value:},{Name:,Value:}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"variables": "[{Name:,Value:},{Name:,Value:}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"re_runable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"re_runable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command_line": "./LammpsTestOne/lammps.pbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_line": "./LammpsTestOne/lammps.pbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem": "1GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem": "1GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stdout_redirect_path": "./LammpsTest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stdout_redirect_path": "./LammpsTest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"array_request": "1-10:2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"array_request": "1-10:2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue": "workq",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue": "workq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_path": "./jobfolderOne",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_path": "./jobfolderOne",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task":                 "4",
					"priority":             "2",
					"stderr_redirect_path": "./LammpsTestT",
					"node":                 "4",
					"clock_time":           "14:00:00",
					"gpu":                  "3",
					"runas_user":           "user3",
					"thread":               "3",
					"job_template_name":    "JobTemplateNameY",
					"variables":            "[{Demo:,Test:},{Test:,Demo:}]",
					"re_runable":           "true",
					"command_line":         "./LammpsTestT/lammps.pbs",
					"mem":                  "3GB",
					"stdout_redirect_path": "./LammpsTestH",
					"array_request":        "1-12:2",
					"queue":                "workq",
					"package_path":         "./jobfolderT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task":                 "4",
						"priority":             "2",
						"stderr_redirect_path": "./LammpsTestT",
						"node":                 "4",
						"clock_time":           "14:00:00",
						"gpu":                  "3",
						"runas_user":           "user3",
						"thread":               "3",
						"job_template_name":    "JobTemplateNameY",
						"variables":            "[{Demo:,Test:},{Test:,Demo:}]",
						"re_runable":           "true",
						"command_line":         "./LammpsTestT/lammps.pbs",
						"mem":                  "3GB",
						"stdout_redirect_path": "./LammpsTestH",
						"array_request":        "1-12:2",
						"queue":                "workq",
						"package_path":         "./jobfolderT",
					}),
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

var AlicloudEhpcJobTemplateMap0 = map[string]string{
	"command_line":      "./LammpsTest/lammps.pbs",
	"job_template_name": "JobTemplateName",
}

func AlicloudEhpcJobTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
