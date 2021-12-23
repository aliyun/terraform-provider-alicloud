package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_direct_mail_domain",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_domain",
			F:    testSweepDirectMailDomain,
		})
}

func testSweepDirectMailDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "QueryDomainByParam"
	request := map[string]interface{}{
		"PageNo":   requests.NewInteger(1),
		"PageSize": requests.NewInteger(PageSizeLarge),
	}

	var dmDomains []interface{}

	conn, err := client.NewDmClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	for {
		var response map[string]interface{}

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
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.data.domain", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.data.domain", action, err)
			return nil
		}
		result, _ := resp.([]interface{})

		if len(result) < 1 {
			break
		}
		dmDomains = append(dmDomains, result...)
		if len(result) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request["PageNo"].(requests.Integer)); err != nil {
			log.Printf("[ERROR] %s get an error: %#v", "getNextpageNumber", err)
			break
		} else {
			request["PageNo"] = page
		}
	}

	sweeped := false
	for _, v := range dmDomains {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["DomainName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Direct Mail Domain : %s", item["DomainName"].(string))
			continue
		}

		sweeped = true
		action := "DeleteDomain"
		request := map[string]interface{}{
			"DomainId": item["DomainId"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Direct Mail Domain (%s): %s", item["DomainName"].(string), err)
		}

		if sweeped {
			// Waiting 5 seconds to ensure Direct Mail Domain have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete Direct Mail Domain success: %s ", item["DomainName"].(string))
	}
	return nil
}

func TestAccAlicloudDirectMailDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d.pop.onaliyun.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailDomainBasicDependence0)
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
					"domain_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name,
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

var AlicloudDirectMailDomainMap0 = map[string]string{
	"status":      CHECKSET,
	"domain_name": CHECKSET,
}

func AlicloudDirectMailDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
