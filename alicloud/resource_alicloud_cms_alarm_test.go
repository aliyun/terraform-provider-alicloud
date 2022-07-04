package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_alarm", &resource.Sweeper{
		Name: "alicloud_cms_alarm",
		F:    testSweepCMSAlarms,
	})
}

func testSweepCMSAlarms(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var alarms []cms.AlarmInDescribeMetricRuleList
	req := cms.CreateDescribeMetricRuleListRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.Page = requests.NewInteger(1)
	for {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(req)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving CMS Alarm: %s", err)
		}
		resp, _ := raw.(*cms.DescribeMetricRuleListResponse)
		if resp == nil || len(resp.Alarms.Alarm) < 1 {
			break
		}
		alarms = append(alarms, resp.Alarms.Alarm...)

		if len(resp.Alarms.Alarm) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(req.Page); err != nil {
			return WrapError(err)
		} else {
			req.Page = page
		}
	}

	for _, v := range alarms {
		name := v.RuleName
		id := v.RuleId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip && v.AlertState == "INSUFFICIENT_DATA" {
			skip = false
		}
		if skip {
			log.Printf("[INFO] Skipping CMS Alarm: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting CMS Alarm: %s (%s). Status: %s", name, id, v.AlertState)
		req := cms.CreateDeleteMetricRulesRequest()
		req.Id = &[]string{id}
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CMS Alarm (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

// At present, the provider does not support creating contact group resource, so you should create manually a contact group
// by web console and set it by environment variable ALICLOUD_CMS_CONTACT_GROUP before running the following test case.
func TestAccAlicloudCmsAlarm_basic(t *testing.T) {
	var alarm map[string]interface{}
	resourceName := "alicloud_cms_alarm.basic"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactGroup%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceName,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.basic", alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "name", "tf-testAccCmsAlarm_basic"),
					resource.TestCheckResourceAttrSet("alicloud_cms_alarm.basic", "metric_dimensions"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_critical.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_warn.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_info.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dimensions", "start_time", "end_time"},
			},
		},
	})
}

func TestAccAlicloudCmsAlarm_basic1(t *testing.T) {
	var alarm map[string]interface{}
	resourceName := "alicloud_cms_alarm.basic"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactGroup%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceName,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_basic1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.basic", alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "name", "tf-testAccCmsAlarm_basic"),
					resource.TestCheckResourceAttrSet("alicloud_cms_alarm.basic", "metric_dimensions"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_critical.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_warn.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "escalations_info.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
			},
		},
	})
}

func TestAccAlicloudCmsAlarm_update(t *testing.T) {
	var alarm map[string]interface{}
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactGroup%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_alarm.update",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "name", "tf-testAccCmsAlarm_update"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_critical.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_warn.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_info.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_cms_alarm.update", "metric_dimensions"),
					resource.TestMatchResourceAttr("alicloud_cms_alarm.update", "webhook", regexp.MustCompile("^https://[0-9]+.eu-central-1.fc.aliyuncs.com/[0-9-]+/proxy/Terraform/AlarmEndpointMock/$")),
				),
			},

			{
				Config: testAccCmsAlarm_updateAfter(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_critical.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_warn.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "escalations_info.#", "1"),
					resource.TestCheckResourceAttrSet("alicloud_cms_alarm.update", "metric_dimensions"),
					resource.TestMatchResourceAttr("alicloud_cms_alarm.update", "webhook", regexp.MustCompile("^https://[0-9]+.eu-central-1.fc.aliyuncs.com/[0-9-]+/proxy/Terraform/AlarmEndpointMock/updated$")),
				),
			},
		},
	})
}

func TestAccAlicloudCmsAlarm_disable(t *testing.T) {
	var alarm map[string]interface{}
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCmsAlarmContactGroup%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_alarm.disable",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_disable(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.disable", alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.disable", "name", "tf-testAccCmsAlarm_disable"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.disable", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckCmsAlarmExists(n string, d map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		alarm, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}

		if alarm.Primary.ID == "" {
			return fmt.Errorf("No Cloud monitor alarm ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsService := CmsService{client}
		attr, err := cmsService.DescribeAlarm(alarm.Primary.ID)
		log.Printf("[DEBUG] check alarm %s attribute %#v", alarm.Primary.ID, attr)

		if err != nil {
			return err
		}

		if attr["RuleId"] == "" {
			return fmt.Errorf("Alarm rule not found")
		}

		d = attr
		return nil
	}
}

func testAccCheckCmsAlarmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_alarm" {
			continue
		}

		alarm, err := cmsService.DescribeAlarm(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if alarm["RuleId"] != "" {
			return fmt.Errorf("Error alarm rule %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCmsAlarm_basic(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	
	resource "alicloud_cms_alarm_contact_group" "default" {
	  alarm_contact_group_name = "${var.name}"
	  describe                 = "Test For Alarm."
	}
	
	resource "alicloud_cms_alarm" "basic" {
	  name               = "tf-testAccCmsAlarm_basic"
	  project            = "acs_ecs_dashboard"
	  metric             = "disk_writebytes"
	  metric_dimensions  = "[{\"instanceId\":\"i-bp1247jeep0y53nu3bnk\",\"device\":\"/dev/vda1\"},{\"instanceId\":\"i-bp11gdcik8z6dl5jm84p\",\"device\":\"/dev/vdb1\"}]"
	  period             = 900
	  escalations_critical {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_warn {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_info {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  contact_groups     = [
      alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
	  effective_interval = "06:00-20:00"
	}
	`, name)
}

func testAccCmsAlarm_basic1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	
	resource "alicloud_cms_alarm_contact_group" "default" {
	  alarm_contact_group_name = "${var.name}"
	  describe                 = "Test For Alarm."
	}
	
	resource "alicloud_cms_alarm" "basic" {
	  name               = "tf-testAccCmsAlarm_basic"
	  project            = "acs_ecs_dashboard"
	  metric             = "disk_writebytes"
	  metric_dimensions  = "[{\"instanceId\":\"i-bp1247jeep0y53nu3bnk\",\"device\":\"/dev/vda1\"},{\"instanceId\":\"i-bp11gdcik8z6dl5jm84p\",\"device\":\"/dev/vdb1\"}]"
	  period             = 900
	  escalations_critical {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_warn {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_info {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  contact_groups     = [
      alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
	  effective_interval = "06:00-20:00"
	}
	`, name)
}

func testAccCmsAlarm_update(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	
	resource "alicloud_cms_alarm_contact_group" "default" {
	  alarm_contact_group_name = "${var.name}"
	  describe                 = "Test For Alarm."
	}
	
	data "alicloud_account" "current" {
	}
	
	resource "alicloud_cms_alarm" "update" {
	  name               = "tf-testAccCmsAlarm_update"
	  project            = "acs_ecs_dashboard"
	  metric             = "disk_writebytes"
	  metric_dimensions  = "[{\"instanceId\":\"i-bp1247jeep0y53nu3bnk\",\"device\":\"/dev/vda1\"},{\"instanceId\":\"i-bp11gdcik8z6dl5jm84p\",\"device\":\"/dev/vdb1\"}]"
	  period             = 900
	  escalations_critical {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_warn {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  escalations_info {
		statistics          = "Average"
		comparison_operator = "<="
		threshold           = 35
		times               = 2
	  }
	  contact_groups     = [
	  alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
	  effective_interval = "06:00-20:00"
	  webhook            = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
	}
`, name)
}

func testAccCmsAlarm_updateAfter(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	
	resource "alicloud_cms_alarm_contact_group" "default" {
	  alarm_contact_group_name = "${var.name}"
	  describe                 = "Test For Alarm."
	}
	
	data "alicloud_account" "current" {
	}
	
	resource "alicloud_cms_alarm" "update" {
	  name               = "tf-testAccCmsAlarm_update"
	  project            = "acs_ecs_dashboard"
	  metric             = "disk_writebytes"
	  metric_dimensions  = "[{\"instanceId\":\"i-bp1247jeep0y53nu3bnk\",\"device\":\"/dev/vda1\"},{\"instanceId\":\"i-bp11gdcik8z6dl5jm84p\",\"device\":\"/dev/vdb1\"}]"
	  period             = 900
	  escalations_critical {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 3
	  }
	  escalations_warn {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 3
	  }
	  escalations_info {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 2
	  }
	  contact_groups     = [
	  alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
	  effective_interval = "06:00-20:00"
	  webhook            = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/updated"
	}
	`, name)
}

func testAccCmsAlarm_disable(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	
	resource "alicloud_cms_alarm_contact_group" "default" {
	  alarm_contact_group_name = "${var.name}"
	  describe                 = "Test For Alarm."
	}
	
	data "alicloud_account" "current" {
	}
	
	resource "alicloud_cms_alarm" "disable" {
	  name               = "tf-testAccCmsAlarm_disable"
	  project            = "acs_ecs_dashboard"
	  metric             = "disk_writebytes"
	  metric_dimensions  = "[{\"instanceId\":\"i-bp1247jeep0y53nu3bnk\",\"device\":\"/dev/vda1\"},{\"instanceId\":\"i-bp11gdcik8z6dl5jm84p\",\"device\":\"/dev/vdb1\"}]"
	  period             = 900
	  escalations_critical {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 3
	  }
	  escalations_warn {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 3
	  }
	  escalations_info {
		statistics          = "Average"
		comparison_operator = "<"
		threshold           = 35
		times               = 2
	  }
	  contact_groups     = [
	  alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
	  effective_interval = "06:00-20:00"
	  enabled            = false
	  webhook            = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
	}
	`, name)
}
