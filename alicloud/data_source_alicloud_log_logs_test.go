package alicloud

import (
	"fmt"
	"os"
	"testing"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogLogsDataSource_basic(t *testing.T) {
	if val := os.Getenv("TEST_LOG_LOGS"); val != "true" {
		return
	}
	config := fmt.Sprintf(testAlicloudLogLogsDefault, time.Now().Unix()-1000, time.Now().Unix()+1000)

	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	accessValue := os.Getenv("ALICLOUD_SECRET_KEY")
	region := os.Getenv("ALICLOUD_REGION")

	client := sls.CreateNormalInterface(region+".log.aliyuncs.com", accessKey, accessValue, "")

	client.CreateProject("tf-testacc-log-logs", "test for TestAccAlicloudLogLogsDataSource_basic")
	client.CreateLogStore("tf-testacc-log-logs", "tf-testacc-log-logs", 2, 2, false, 0)

	defer func() {
		client.DeleteLogStore("tf-testacc-log-logs", "tf-testacc-log-logs")
		client.DeleteProject("tf-testacc-log-logs")
	}()

	client.CreateIndex("tf-testacc-log-logs", "tf-testacc-log-logs", *sls.CreateDefaultIndex())

	time.Sleep(time.Second * 60)

	c1 := &sls.LogContent{
		Key:   proto.String("key-1"),
		Value: proto.String("error"),
	}
	c2 := &sls.LogContent{
		Key:   proto.String("key-2"),
		Value: proto.String("InternalServerError"),
	}
	c3 := &sls.LogContent{
		Key:   proto.String("content"),
		Value: proto.String("internal server have some errors"),
	}
	l := &sls.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
		Contents: []*sls.LogContent{
			c1,
			c2,
			c3,
		},
	}
	lg := &sls.LogGroup{
		Topic:  proto.String("demo topic"),
		Source: proto.String("10.230.201.117"),
		Logs:   []*sls.Log{},
	}
	logCount := 50
	for i := 0; i < logCount; i++ {
		lg.Logs = append(lg.Logs, l)
	}

	err := client.PostLogStoreLogs("tf-testacc-log-logs", "tf-testacc-log-logs", lg, nil)
	if err != nil {
		t.Error("post logs error : ", err)
	}

	time.Sleep(time.Second * 3)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_log_logs.all", "logs.#", "50"),
				),
			},
		},
	})
}

const testAlicloudLogLogsDefault = `
variable "name" {
    default = "tf-testacc-log-logs"
}

variable "key" {
    default = "key"
}

data "alicloud_log_logs" "all" {
    project = "${var.name}"
	logstore = "${var.name}"
	from = %d
	to = %d
	query = "${var.key}-1 error"
	output_file = "./logs.json"
}
`
