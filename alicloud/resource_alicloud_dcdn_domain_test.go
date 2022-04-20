package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	queryRequest := dcdn.CreateDescribeDcdnUserDomainsRequest()
	var allDomains []dcdn.PageData
	queryRequest.PageSize = requests.NewInteger(PageSizeLarge)
	queryRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.DescribeDcdnUserDomains(queryRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %#v", queryRequest.GetActionName(), err)
		}
		addDebug(queryRequest.GetActionName(), raw)
		response, _ := raw.(*dcdn.DescribeDcdnUserDomainsResponse)
		domains := response.Domains.PageData

		for _, domain := range domains {
			if strings.HasPrefix(domain.DomainName, "tf-testacc") {
				allDomains = append(allDomains, domain)
			} else {
				log.Printf("Skip %#v", domain)
			}
		}

		if len(domains) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(queryRequest.PageNumber); err != nil {
			return WrapError(err)
		} else {
			queryRequest.PageNumber = page
		}
	}
	removeRequest := dcdn.CreateDeleteDcdnDomainRequest()
	removeRequest.DomainName = ""
	for _, domain := range allDomains {
		removeRequest.DomainName = domain.DomainName
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.DeleteDcdnDomain(removeRequest)
		})

		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		addDebug(removeRequest.GetActionName(), raw)
	}

	return nil
}

func TestAccAlicloudDCDNDomain_basic0(t *testing.T) {
	var v dcdn.DomainDetail
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, DcdnDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DcdnDomainBasicdependence)
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
					"scope":       "overseas",
					"status":      "online",
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "security_token", "cert_type", "check_url", "force_set", "ssl_pri"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_protocol": "on",
					"ssl_pub":      testDcdnPublicKey,
					"ssl_pri":      testDcdnPrivateKey,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_protocol": "on",
						"ssl_pub":      strings.Replace(testDcdnPublicKey, `\n`, "\n", -1),
						"ssl_pri":      strings.Replace(testDcdnPrivateKey, `\n`, "\n", -1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.update.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
						{
							"content": "2.2.2.2",
							"type":    "ipaddr",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDCDNDomain_basic1(t *testing.T) {
	var v dcdn.DomainDetail
	resourceId := "alicloud_dcdn_domain.default"
	ra := resourceAttrInit(resourceId, DcdnDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DcdnDomainBasicdependence)
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
					"domain_name":       name,
					"scope":             "overseas",
					"status":            "online",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"port":     "80",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":       name,
						"sources.#":         "1",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "security_token", "cert_type", "check_url", "force_set", "ssl_pri"},
			},
		},
	})
}

var DcdnDomainMap = map[string]string{
	"resource_group_id": CHECKSET,
	"ssl_protocol":      "off",
	"scope":             "overseas",
	"status":            "online",
}

func DcdnDomainBasicdependence(name string) string {
	return fmt.Sprintf(`
	
variable "name" {
	default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_resource_manager_resource_groups" "update" {
  name_regex = "terraformci"
}
`, name)
}

const testDcdnPrivateKey = `-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n`
const testDcdnPublicKey = `-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n`
