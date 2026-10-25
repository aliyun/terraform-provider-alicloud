package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The DTS sweepers have to stay in the root test package: the sweep entry
// (make sweep) runs go test ./alicloud only, so a sweeper registered from
// alicloud/service/dts would never be executed.
func init() {
	resource.AddTestSweepers(
		"alicloud_dts_migration_job",
		&resource.Sweeper{
			Name: "alicloud_dts_migration_job",
			F:    testSweepDTSMigrationJob,
		})
	resource.AddTestSweepers(
		"alicloud_dts_check_job",
		&resource.Sweeper{
			Name: "alicloud_dts_check_job",
			F:    testSweepDTSCheckJob,
		})
}

func testSweepDTSMigrationJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDtsJobs"
	request := map[string]interface{}{}
	request["JobType"] = "MIGRATION"
	request["PageNumber"] = 1
	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DtsJobList", response)
		totalCount, formatErr := formatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DtsJobList", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["DtsJobName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DtsJobName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DTS Migration Job: %s", item["DtsJobName"].(string))
				continue
			}
			action := "DeleteDtsJob"
			request := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			_, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DTS Migration Job (%s): %s", item["DtsJobId"].(string), err)
			}
			log.Printf("[INFO] Delete DTS Migration Job success: %s ", item["DtsJobId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepDTSCheckJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		// The resource documentation example names its job "terraform-example".
		"terraform-example",
	}
	action := "DescribeCheckJobs"
	request := map[string]interface{}{}
	request["PageNumber"] = 1
	request["PageSize"] = PageSizeXLarge

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.CheckJobs", response)
		if err != nil {
			if fmt.Sprint(response["TotalRecordCount"]) == "0" {
				return nil
			}
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.CheckJobs", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			jobName, _ := item["JobName"].(string)
			if jobName == "" {
				continue
			}
			if regionId, ok := item["RegionId"].(string); ok && regionId != "" && regionId != region {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(jobName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DTS Check Job: %s", jobName)
				continue
			}
			action := "DeleteDtsJob"
			deleteRequest := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			if v, ok := item["DtsInstanceID"]; ok {
				deleteRequest["DtsInstanceId"] = v
			}
			_, err = client.RpcPost("Dts", "2020-01-01", action, nil, deleteRequest, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DTS Check Job (%s): %s", item["DtsJobId"], err)
				continue
			}
			log.Printf("[INFO] Delete DTS Check Job success: %s ", item["DtsJobId"])
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
