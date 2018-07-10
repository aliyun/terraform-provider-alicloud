package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// At present, the provider does not support creating contact group resource, so you should add a contact group called "tf-acc-test-group"
// by web console manually before running the following test case.
func SkipTestAccAlicloudCmsAlarm_basic(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_alarm.basic",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCmsAlarm_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.basic", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "name", "testAccCmsAlarm_basic"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "dimensions.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "dimensions.device", "/dev/vda1,/dev/vdb1"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCmsAlarm_update(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_alarm.update",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCmsAlarm_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "name", "testAccCmsAlarm_update"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "operator", "<="),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "triggered_count", "2"),
				),
			},

			resource.TestStep{
				Config: testAccCmsAlarm_updateAfter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "operator", "=="),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "triggered_count", "3"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCmsAlarm_disable(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_alarm.disable",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCmsAlarm_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.disable", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.disable", "name", "testAccCmsAlarm_disable"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.disable", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckCmsAlarmExists(n string, d *cms.AlarmInListAlarm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		alarm, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}

		if alarm.Primary.ID == "" {
			return fmt.Errorf("No Cloud monitor alarm ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeAlarm(alarm.Primary.ID)
		log.Printf("[DEBUG] check alarm %s attribute %#v", alarm.Primary.ID, attr)

		if err != nil {
			return err
		}

		if attr.Id == "" {
			return fmt.Errorf("Alarm rule not found")
		}

		*d = attr
		return nil
	}
}

func testAccCheckCmsAlarmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_alarm" {
			continue
		}

		alarm, err := client.DescribeAlarm(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if alarm.Id != "" {
			return fmt.Errorf("Error alarm rule %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCmsAlarm_basic = `
resource "alicloud_cms_alarm" "basic" {
  name = "testAccCmsAlarm_basic"
  project = "acs_ecs_dashboard"
  metric = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247jeep0y53nu3bnk,i-bp11gdcik8z6dl5jm84p"
    device = "/dev/vda1,/dev/vdb1"
  }
  statistics ="Average"
  period = 900
  operator = "<="
  threshold = 35
  triggered_count = 2
  contact_groups = ["tf-acc-test-group"]
  end_time = 20
  start_time = 6
  notify_type = 1
}
`

const testAccCmsAlarm_update = `
resource "alicloud_cms_alarm" "update" {
  name = "testAccCmsAlarm_update"
  project = "acs_ecs_dashboard"
  metric = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247jeep0y53nu3bnk,i-bp11gdcik8z6dl5jm84p"
    device = "/dev/vda1,/dev/vdb1"
  }
  statistics ="Average"
  period = 900
  operator = "<="
  threshold = 35
  triggered_count = 2
  contact_groups = ["tf-acc-test-group"]
  end_time = 20
  start_time = 6
  notify_type = 1
}
`

const testAccCmsAlarm_updateAfter = `
resource "alicloud_cms_alarm" "update" {
  name = "testAccCmsAlarm_update"
  project = "acs_ecs_dashboard"
  metric = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247jeep0y53nu3bnk,i-bp11gdcik8z6dl5jm84p"
    device = "/dev/vda1,/dev/vdb1"
  }
  statistics ="Average"
  period = 900
  operator = "=="
  threshold = 35
  triggered_count = 3
  contact_groups = ["tf-acc-test-group"]
  end_time = 20
  start_time = 6
  notify_type = 1
}
`

const testAccCmsAlarm_disable = `
resource "alicloud_cms_alarm" "disable" {
  name = "testAccCmsAlarm_disable"
  project = "acs_ecs_dashboard"
  metric = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247jeep0y53nu3bnk,i-bp11gdcik8z6dl5jm84p"
    device = "/dev/vda1,/dev/vdb1"
  }
  statistics ="Average"
  period = 900
  operator = "=="
  threshold = 35
  triggered_count = 3
  contact_groups = ["tf-acc-test-group"]
  end_time = 20
  start_time = 6
  notify_type = 1
  enabled = false
}
`
