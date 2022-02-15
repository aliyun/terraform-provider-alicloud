package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_direct_mail_receivers", &resource.Sweeper{
		Name: "alicloud_direct_mail_receivers",
		F:    testSweepDirectMailReceivers,
	})
}

func testSweepDirectMailReceivers(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	conn, err := client.NewDmClient()
	if err != nil {
		log.Println(WrapError(err))
		return nil
	}
	action := "QueryReceiverByParam"
	request := map[string]interface{}{
		"PageSize": PageSizeLarge,
		"PageNo":   1,
	}
	Ids := make([]string, 0)
	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List Direct Mail Receivers Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.data.receiver", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.receiver", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ReceiversName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Direct Mail Receivers: %v (%v)", item["ReceiversName"], item["ReceiverId"])
				continue
			}
			Ids = append(Ids, fmt.Sprint(item["ReceiverId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}
	for _, Id := range Ids {
		log.Printf("[INFO] Deleting Direct Mail Receivers: (%s)", Id)
		action := "DeleteReceiver"
		conn, err := client.NewDmClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"ReceiverId": Id,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*9, func() *resource.RetryError {
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed To Delete Direct Mail Receivers : %s", err)
			continue
		}
		log.Printf("[INFO] Delete Direct Mail Receivers Success : %s", Id)
	}
	return nil
}

func TestAccAlicloudDirectMailReceivers_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_receivers.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailReceiversMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailReceivers")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailReceiversBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"receivers_name":  "${var.name}",
					"receivers_alias": name + "@onaliyun.com",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"receivers_name":  name,
						"receivers_alias": name + "@onaliyun.com",
						"description":     name,
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

var AlicloudDirectMailReceiversMap0 = map[string]string{}

func AlicloudDirectMailReceiversBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
