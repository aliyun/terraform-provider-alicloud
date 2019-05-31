package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"strings"

	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var cmsContactGroup = os.Getenv("ALICLOUD_CMS_CONTACT_GROUP")

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
			log.Printf("[ERROR] Error retrieving CMS Alarm: %s", err)
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
		if skip && v.State == "INSUFFICIENT_DATA" {
			skip = false
		}
		if skip {
			log.Printf("[INFO] Skipping CMS Alarm: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting CMS Alarm: %s (%s). Status: %s", name, id, v.State)
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

// At present, the provider does not support creating contact group resource, so you should create manually a contact group
// by web console and set it by environment variable ALICLOUD_CMS_CONTACT_GROUP before running the following test case.
func TestAccAlicloudCmsAlarm_basic(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCmsContactGroupSetting(t)
		},

		IDRefreshName: "alicloud_cms_alarm.basic",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_basic(cmsContactGroup),
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

func TestAccAlicloudCmsAlarm_update(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCmsContactGroupSetting(t)
		},

		IDRefreshName: "alicloud_cms_alarm.update",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_update(cmsContactGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "name", "tf-testAccCmsAlarm_update"),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "operator", "<="),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "triggered_count", "2"),
					resource.TestMatchResourceAttr("alicloud_cms_alarm.update", "webhook", regexp.MustCompile("^https://[0-9]+.eu-central-1.fc.aliyuncs.com/[0-9-]+/proxy/Terraform/AlarmEndpointMock/$")),
				),
			},

			{
				Config: testAccCmsAlarm_updateAfter(cmsContactGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCmsAlarmExists("alicloud_cms_alarm.update", &alarm),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "operator", "=="),
					resource.TestCheckResourceAttr("alicloud_cms_alarm.update", "triggered_count", "3"),
					resource.TestMatchResourceAttr("alicloud_cms_alarm.update", "webhook", regexp.MustCompile("^https://[0-9]+.eu-central-1.fc.aliyuncs.com/[0-9-]+/proxy/Terraform/AlarmEndpointMock/updated$")),
				),
			},
		},
	})
}

func TestAccAlicloudCmsAlarm_disable(t *testing.T) {
	var alarm cms.AlarmInListAlarm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCmsContactGroupSetting(t)
		},

		IDRefreshName: "alicloud_cms_alarm.disable",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_disable(cmsContactGroup),
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

func testAccCmsAlarm_basic(group string) string {
	return fmt.Sprintf(`
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
	  contact_groups = ["%s"]
	  end_time = 20
	  start_time = 6
	  notify_type = 1
	}
	`, group)
}

func testAccCmsAlarm_update(group string) string {
	return fmt.Sprintf(`
data "alicloud_account" "current"{
}

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
  contact_groups = ["%s"]
  end_time = 20
  start_time = 6
  notify_type = 1
  webhook = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
}
`, group)
}

func testAccCmsAlarm_updateAfter(group string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "current"{
	}
	
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
	  contact_groups = ["%s"]
	  end_time = 20
	  start_time = 6
	  notify_type = 1
  	  webhook = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/updated"
	}
	`, group)
}

func testAccCmsAlarm_disable(group string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "current"{
	}
	
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
	  contact_groups = ["%s"]
	  end_time = 20
	  start_time = 6
	  notify_type = 1
	  enabled = false
	  webhook = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
	}
	`, group)
}
