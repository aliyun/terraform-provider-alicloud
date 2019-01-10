package alicloud

import (
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccAlicloudLogtailToMachineGroupDataSource(t *testing.T) {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	accessValue := os.Getenv("ALICLOUD_SECRET_KEY")
	region := os.Getenv("ALICLOUD_REGION")
	localFileConfigInputDetail := sls.LocalFileConfigInputDetail{LogType: "json_log", LogPath: "/root", TopicFormat: "default"}
	outputDetail := sls.OutputDetail{ProjectName: "tf-logproject1", LogStoreName: "tf-testacc-log-logs"}
	logconfig := sls.LogConfig{
		Name: "evan-terraform-config",
		InputType: "file",
		InputDetail: sls.JSONConfigInputDetail{LocalFileConfigInputDetail: localFileConfigInputDetail},
		OutputType: "LogService",
		OutputDetail: outputDetail,
	}
	m := sls.MachineGroup{
		Name: "evan-machine-group",
		MachineIDType: "ip",
		MachineIDList: []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
		Attribute: sls.MachinGroupAttribute{ExternalName: "testgroup", TopicName: "test"},
	}

	client := sls.CreateNormalInterface(region+".log.aliyuncs.com", accessKey, accessValue, "")
	client.CreateProject("tf-logproject1", "test for TestAccAlicloudLogtailToMachineGroupDataSource")
	client.CreateLogStore("tf-testacc-log-logs", "tf-testacc-log-logs", 2, 2, false, 0)
	client.CreateMachineGroup("tf-logproject1", &m)
	client.CreateConfig("tf-logproject1", &logconfig)
	defer func() {
		client.DeleteConfig("tf-logproject1", "evan-terraform-config")
		client.DeleteMachineGroup("tf-logproject1", "evan-machine-group")
		client.DeleteLogStore("tf-logproject1", "tf-testacc-log-logs")
		client.DeleteProject("tf-logproject1")
	}()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailToMachineGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_logtail_to_machine_group.example"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "project", "tf-logproject1"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "logtail_config.0", "evan-terraform-config"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "machine_group.0", "evan-machine-group"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogtailToMachineGroupDataSource = `
data "alicloud_logtail_to_machine_group" "example" {
   project = "tf-logproject1"
   output_file = "~/newdata/map.json"
}
`
