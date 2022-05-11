package alicloud

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cdn_domain_new", &resource.Sweeper{
		Name: "alicloud_cdn_domain_new",
		F:    testSweepCdnDomains_new,
	})
}

func testSweepCdnDomains_new(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testacc%s", region),
		fmt.Sprintf("tf_testacc%s", region),
	}

	var domains []cdn.PageData
	args := cdn.CreateDescribeUserDomainsRequest()
	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	for {

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DescribeUserDomains(args)
		})
		if err != nil {
			log.Printf("Error retrieving CDN Domain new: %s", err)
		}
		addDebug(args.GetActionName(), raw)
		resp, _ := raw.(*cdn.DescribeUserDomainsResponse)
		if resp == nil || len(resp.Domains.PageData) < 1 {
			break
		}
		domains = append(domains, resp.Domains.PageData...)

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			log.Printf("Error get Next page Number: %s", err)
			break
		} else {
			args.PageNumber = page
		}
	}

	for _, v := range domains {
		name := v.DomainName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CDN domain: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting CDN domain: %s", name)
		request := cdn.CreateDeleteCdnDomainRequest()
		request.DomainName = name
		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN domain (%s): %s", name, err)
		}
		addDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAlicloudCDNDomainNew_basic(t *testing.T) {
	var v *cdn.GetDomainDetailModel

	resourceId := "alicloud_cdn_domain_new.domain"
	ra := resourceAttrInit(resourceId, cdnDomainBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, strconv.Itoa(rand), resourceCdnDomainDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":       name,
					"cdn_type":          "web",
					"scope":             "domestic",
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					"sources": []map[string]interface{}{
						{
							"content": "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":    "oss",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":             "domestic",
						"resource_group_id": CHECKSET,
						"sources.#":         "1",
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
					"sources": []map[string]interface{}{
						{
							"content":  "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":     "oss",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "www.aliyuntest.com",
							"type":     "domain",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"scope": "domestic",
					"sources": []map[string]interface{}{
						{
							"content": "1.1.1.1",
							"type":    "ipaddr",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "www.aliyuntest.com",
							"type":     "domain",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
						{
							"content":  "1.1.1.1",
							"type":     "ipaddr",
							"priority": "40",
							"port":     "80",
							"weight":   "10",
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
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate": testServerCertificate,
							"private_key":        testPrivateKey,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.server_certificate_status": "on",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate": testServerCertificate,
							"private_key":        testPrivateKey,
							"force_set":          "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.force_set": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate": testServerCertificate,
							"private_key":        testPrivateKey,
							"force_set":          "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.force_set": "0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate": testServerCertificate,
							"private_key":        testPrivateKey,
							"force_set":          "1",
							"cert_name":          "tf-test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.force_set": "1",
						"certificate_config.0.cert_name": "tf-test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate": testServerCertificate,
							"private_key":        testPrivateKey,
							"force_set":          "1",
							"cert_name":          "tf-test",
							"cert_type":          "cas",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.cert_type": "cas",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":        name,
					"cdn_type":           "web",
					"scope":              REMOVEKEY,
					"certificate_config": REMOVEKEY,
					"sources": []map[string]interface{}{
						{
							"content": "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":    "oss",
						},
						{
							"content":  "www.aliyuntest.com",
							"type":     "domain",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
						{
							"content":  "1.1.1.1",
							"type":     "ipaddr",
							"priority": "40",
							"port":     "80",
							"weight":   "10",
						},
					},
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":     REMOVEKEY,
						"sources.#": "3",

						"certificate_config.0.server_certificate_status": REMOVEKEY,
						"certificate_config.0.force_set":                 REMOVEKEY,
						"certificate_config.0.cert_name":                 REMOVEKEY,
						"certificate_config.0.cert_type":                 REMOVEKEY,

						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
						"tags.Updated": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCDNDomainNew_scope(t *testing.T) {
	var v *cdn.GetDomainDetailModel

	resourceId := "alicloud_cdn_domain_new.domain"
	ra := resourceAttrInit(resourceId, cdnDomainBasicMap)

	serviceFunc := func() interface{} {
		return &CdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.xiaozhu.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, strconv.Itoa(rand), resourceCdnDomainDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": name,
					"cdn_type":    "web",
					"scope":       "overseas",
					"sources": []map[string]interface{}{
						{
							"content": "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":    "oss",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":     "overseas",
						"sources.#": "1",
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
					"sources": []map[string]interface{}{
						{
							"content":  "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":     "oss",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "1",
					}),
				),
			},
		},
	})
}

func resourceCdnDomainDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_oss_bucket" "default" {
	  bucket = "tf-testacc-cdn-%s"
	}`, name)
}

var cdnDomainBasicMap = map[string]string{
	"domain_name": CHECKSET,
	"scope":       CHECKSET,
	"cdn_type":    "web",
}

const testServerCertificate = `-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n`
const testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n`
