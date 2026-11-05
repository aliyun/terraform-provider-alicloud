package sls_test

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	slssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogIngestion_basic(t *testing.T) {
	var inges *sls.Ingestion
	resourceId := "alicloud_log_ingestion.default"
	ra := acctest.ResourceAttrInit(resourceId, logIngestionMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &inges, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogingestion-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogIngestionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":         "${alicloud_log_project.default.name}",
					"logstore":        "${alicloud_log_store.default.name}",
					"ingestion_name":  "test_ingestion",
					"display_name":    "test_display",
					"interval":        "5m",
					"run_immediately": "false",
					"source":          `{\"bucket\":\"bucket_name\",\"compressionCodec\":\"none\",\"encoding\":\"UTF-8\",\"endpoint\":\"oss-cn-hangzhou-internal.aliyuncs.com\",\"format\":{\"escapeChar\":\"\\\\\",\"fieldDelimiter\":\",\",\"fieldNames\":[],\"firstRowAsHeader\":true,\"maxLines\":1,\"quoteChar\":\"\\\"\",\"skipLeadingRows\":0,\"timeField\":\"\",\"type\":\"DelimitedText\"},\"pattern\":\"\",\"prefix\":\"test_prefix/\",\"restoreObjectEnabled\":false,\"roleARN\":\"acs:ram::1049446484210612:role/aliyunlogimportossrole\",\"type\":\"AliyunOSS\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":         name,
						"logstore":        name,
						"ingestion_name":  "test_ingestion",
						"display_name":    "test_display",
						"interval":        "5m",
						"run_immediately": "false",
						"source":          `{"bucket":"bucket_name","compressionCodec":"none","encoding":"UTF-8","endpoint":"oss-cn-hangzhou-internal.aliyuncs.com","format":{"escapeChar":"\\","fieldDelimiter":",","fieldNames":[],"firstRowAsHeader":true,"maxLines":1,"quoteChar":"\"","skipLeadingRows":0,"timeField":"","type":"DelimitedText"},"pattern":"","prefix":"test_prefix/","restoreObjectEnabled":false,"roleARN":"acs:ram::1049446484210612:role/aliyunlogimportossrole","type":"AliyunOSS"}`,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "test_display_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "test_display_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"interval": "30m",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"interval": "30m",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"run_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_immediately": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_zone": "+0800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_zone": "+0800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": `{\"bucket\":\"bucket_name\",\"compressionCodec\":\"none\",\"encoding\":\"UTF-8\",\"endpoint\":\"oss-cn-hangzhou-internal.aliyuncs.com\",\"format\":{\"escapeChar\":\"\\\\\",\"fieldDelimiter\":\",\",\"fieldNames\":[],\"firstRowAsHeader\":false,\"maxLines\":1,\"quoteChar\":\"\\\"\",\"skipLeadingRows\":0,\"timeField\":\"\",\"type\":\"DelimitedText\"},\"pattern\":\"\",\"prefix\":\"test_prefix/\",\"restoreObjectEnabled\":false,\"roleARN\":\"acs:ram::1049446484210612:role/aliyunlogimportossrole\",\"type\":\"AliyunOSS\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": `{"bucket":"bucket_name","compressionCodec":"none","encoding":"UTF-8","endpoint":"oss-cn-hangzhou-internal.aliyuncs.com","format":{"escapeChar":"\\","fieldDelimiter":",","fieldNames":[],"firstRowAsHeader":false,"maxLines":1,"quoteChar":"\"","skipLeadingRows":0,"timeField":"","type":"DelimitedText"},"pattern":"","prefix":"test_prefix/","restoreObjectEnabled":false,"roleARN":"acs:ram::1049446484210612:role/aliyunlogimportossrole","type":"AliyunOSS"}`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name":    "test_display",
					"interval":        "5m",
					"run_immediately": "false",
					"source":          `{\"bucket\":\"bucket_name\",\"compressionCodec\":\"none\",\"encoding\":\"UTF-8\",\"endpoint\":\"oss-cn-hangzhou-internal.aliyuncs.com\",\"format\":{\"escapeChar\":\"\\\\\",\"fieldDelimiter\":\",\",\"fieldNames\":[],\"firstRowAsHeader\":true,\"maxLines\":1,\"quoteChar\":\"\\\"\",\"skipLeadingRows\":0,\"timeField\":\"\",\"type\":\"DelimitedText\"},\"pattern\":\"\",\"prefix\":\"test_prefix/\",\"restoreObjectEnabled\":false,\"roleARN\":\"acs:ram::1049446484210612:role/aliyunlogimportossrole\",\"type\":\"AliyunOSS\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":    "test_display",
						"interval":        "5m",
						"run_immediately": "false",
						"source":          `{"bucket":"bucket_name","compressionCodec":"none","encoding":"UTF-8","endpoint":"oss-cn-hangzhou-internal.aliyuncs.com","format":{"escapeChar":"\\","fieldDelimiter":",","fieldNames":[],"firstRowAsHeader":true,"maxLines":1,"quoteChar":"\"","skipLeadingRows":0,"timeField":"","type":"DelimitedText"},"pattern":"","prefix":"test_prefix/","restoreObjectEnabled":false,"roleARN":"acs:ram::1049446484210612:role/aliyunlogimportossrole","type":"AliyunOSS"}`,
					}),
				),
			},
		},
	})
}

func resourceLogIngestionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
`, name)
}

var logIngestionMap = map[string]string{
	"project":        acctest.CHECKSET,
	"logstore":       acctest.CHECKSET,
	"ingestion_name": acctest.CHECKSET,
}
