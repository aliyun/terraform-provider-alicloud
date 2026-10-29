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

func init() {
	resource.AddTestSweepers("alicloud_ssl_certificates_service_certificate", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_certificate",
		F:    testSweepSslCertificatesServiceCertificate,
	})
}

func testSweepSslCertificatesServiceCertificate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"cert-tf-testacc",
		"cert-tf_testacc",
	}

	action := "DescribeUserCertificateList"
	request := make(map[string]interface{})
	request["ShowSize"] = PageSizeXLarge
	request["CurrentPage"] = 1
	ids := make([]string, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_certificates_service_certificates", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.CertificateList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CertificateList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["name"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping ssl certificate: %s ", item["name"])
					continue
				}
			}
			ids = append(ids, fmt.Sprint(item["id"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, sslId := range ids {
		log.Printf("[INFO] Deleting ssl centrficate: %s ", sslId)
		action = "DeleteUserCertificate"
		request = map[string]interface{}{
			"CertId": sslId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ssl centrficate (%s): %s", sslId, err)
		}
	}
	return nil
}
