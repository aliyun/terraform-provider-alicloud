package alicloud

import (
	"fmt"
	"log"
	"os"
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
	resource.AddTestSweepers(
		"alicloud_scdn_domain",
		&resource.Sweeper{
			Name: "alicloud_scdn_domain",
			F:    testSweepScdnDomain,
		})
}

func testSweepScdnDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "DescribeScdnUserDomains"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeXLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewScdnClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DomainName"])
			if !strings.HasPrefix(strings.ToLower(name), "tf-testacc") {
				continue
			}
			log.Println("[ERROR] Deleting the domain ", name)
			request := map[string]interface{}{
				"DomainName": name,
			}

			_, err = conn.DoRequest(StringPointer("DeleteScdnDomain"), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Println("[ERROR] Deleting the domain ", name, " failed! Error: ", err)
			}
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

func TestAccAlicloudScdnDomain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_scdn_domain.default"
	ra := resourceAttrInit(resourceId, ScdnDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ScdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeScdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ScdnDomainBasicdependence)
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
					"domain_name": name,
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name,
						"sources.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_infos": []map[string]interface{}{
						{
							"cert_name":    name,
							"cert_type":    "upload",
							"ssl_pub":      testScdnPublicKey,
							"ssl_pri":      testScdnPrivateKey,
							"ssl_protocol": "off",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_infos.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_url": "www.yourdomain.com/test.html",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_url": "www.yourdomain.com/test.html",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "xiaozhu.aliyuncs.com",
							"port":     "80",
							"priority": "20",
							"type":     "oss",
						},
						{
							"content":  "xiaozhu.aliyuncs.com",
							"type":     "domain",
							"port":     "90",
							"priority": "21",
							"enabled":  "offline",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
					"cert_infos": []map[string]interface{}{
						{
							"cert_name":    name + "update",
							"cert_type":    "upload",
							"ssl_pub":      testScdnPublicKey,
							"ssl_pri":      testScdnPrivateKey,
							"ssl_protocol": "off",
						},
					},
					"check_url": "www.yourdomainupdate.com/test.html",
					// There is an OpenAPI bug
					//"sources": []map[string]interface{}{
					//	{
					//		"content":  "1.1.1.1",
					//		"port":     "80",
					//		"priority": "20",
					//		"type":     "ipaddr",
					//	},
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":       "online",
						"cert_infos.#": "1",
						"check_url":    "www.yourdomainupdate.com/test.html",
						//"sources.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_set", "check_url", "force_set", "cert_infos"},
			},
		},
	})
}

var ScdnDomainMap = map[string]string{
	"resource_group_id": CHECKSET,
	"status":            "online",
}

func ScdnDomainBasicdependence(name string) string {
	return ""
}

const testScdnPrivateKey = `-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n`
const testScdnPublicKey = `-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n`
