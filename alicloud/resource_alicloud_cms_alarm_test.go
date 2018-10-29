package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var alarms []cms.AlarmInListAlarm
	req := cms.CreateListAlarmRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ListAlarm(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CMS Alarm: %s", err)
		}
		resp, _ := raw.(*cms.ListAlarmResponse)
		if resp == nil || len(resp.AlarmList.Alarm) < 1 {
			break
		}
		alarms = append(alarms, resp.AlarmList.Alarm...)

		if len(resp.AlarmList.Alarm) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range alarms {
		name := v.Name
		id := v.Id
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CMS Alarm: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting CMS Alarm: %s (%s)", name, id)
		req := cms.CreateDeleteAlarmRequest()
		req.Id = id
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteAlarm(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CMS Alarm (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

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
					resource.TestCheckResourceAttr("alicloud_cms_alarm.basic", "name", "tf-testAccCmsAlarm_basic"),
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
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "name", "tf-testAccCmsAlarm_update"),
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
					resource.TestCheckResourceAttr("alicloud_cms_alarm.disable", "name", "tf-testAccCmsAlarm_disable"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsService := CmsService{client}
		attr, err := cmsService.DescribeAlarm(alarm.Primary.ID)
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

		if alarm.Id != "" {
			return fmt.Errorf("Error alarm rule %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCmsAlarm_basic = `
resource "alicloud_cms_alarm" "basic" {
  name = "tf-testAccCmsAlarm_basic"
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
  name = "tf-testAccCmsAlarm_update"
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
  name = "tf-testAccCmsAlarm_update"
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
  name = "tf-testAccCmsAlarm_disable"
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
