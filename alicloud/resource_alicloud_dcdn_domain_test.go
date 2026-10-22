package alicloud

import (
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_domain",
		&resource.Sweeper{
			Name: "alicloud_dcdn_domain",
			F:    testSweepDcdnDomain,
		})
}

func testSweepDcdnDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeDcdnUserDomains"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var allDomains []string
	var response map[string]interface{}
	for {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if strings.HasPrefix(item["DomainName"].(string), "tf-testacc") {
				allDomains = append(allDomains, item["DomainName"].(string))
			} else {
				log.Printf("Skip %#v", item["DomainName"].(string))
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	action = "DeleteDcdnDomain"
	for _, domain := range allDomains {
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		query["DomainName"] = domain

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, query, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, domain, action, AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}
