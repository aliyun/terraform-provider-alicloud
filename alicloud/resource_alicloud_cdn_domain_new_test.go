package alicloud

import (
	"fmt"
	"log"
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

func TestAccAliCloudCDNDomainNew_basic(t *testing.T) {
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
	name := fmt.Sprintf("test%s%d.pfytlm.xyz", defaultRegionToTest, rand)
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
					"env":               "online",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"sources": []map[string]interface{}{
						{
							"content": "${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}",
							"type":    "oss",
						},
					},
					"check_url": "http://" + name + "/test.html",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":       CHECKSET,
						"scope":             "domestic",
						"resource_group_id": CHECKSET,
						"sources.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env":               "",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url", "env"},
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
							"server_certificate_status": "on",
							"server_certificate":        testServerCertificate,
							"private_key":               testPrivateKey,
							"cert_name":                 name,
							"cert_type":                 "cas",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.server_certificate_status": "on",
						"certificate_config.0.cert_type":                 "cas",
						"certificate_config.0.cert_name":                 name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate_status": "on",
							"cert_id":                   "${alicloud_ssl_certificates_service_certificate.default.id}",
							"cert_region":               "cn-hangzhou",
							"cert_type":                 "cas",
							"cert_name":                 REMOVEKEY,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_config.0.server_certificate_status": "on",
						"certificate_config.0.cert_id":                   CHECKSET,
						"certificate_config.0.cert_region":               CHECKSET,
						"certificate_config.0.cert_type":                 "cas",
						"certificate_config.0.cert_name":                 REMOVEKEY,
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
					"domain_name": name,
					"cdn_type":    "web",
					"scope":       REMOVEKEY,
					"certificate_config": []map[string]interface{}{
						{
							"server_certificate_status": "off",
						},
					},
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
						"scope":                            REMOVEKEY,
						"sources.#":                        "3",
						"certificate_config.0.cert_id":     "",
						"certificate_config.0.cert_name":   "",
						"certificate_config.0.cert_region": "",
						"certificate_config.0.server_certificate_status": "off",
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

func TestAccAliCloudCDNDomainNew_scope(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
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
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = alicloud_oss_bucket.default.bucket
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIDeDCCAmCgAwIBAgIEN3ZT6zANBgkqhkiG9w0BAQsFADBVMQswCQYDVQQGEwJD
TjEVMBMGA1UEAwwMKi50ZnRlc3QudG9wMRAwDgYDVQQIDAdCZWlKaW5nMRAwDgYD
VQQHDAdCZWlKaW5nMQswCQYDVQQKDAJBQTAeFw0yMzA4MjgwNjQ5NDNaFw0yNTA4
MjcwNjQ5NDNaMFUxCzAJBgNVBAYTAkNOMRUwEwYDVQQDDAwqLnRmdGVzdC50b3Ax
EDAOBgNVBAgMB0JlaUppbmcxEDAOBgNVBAcMB0JlaUppbmcxCzAJBgNVBAoMAkFB
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzkk9NJUH7PLSQK4RRrGQ
Y5dVsftkhnKh9HhI6yrnlowWIDPS1PZHOU/5gQ7xPUPGdKQV5S7x8wROnAaXEimx
N4GdQw25pGhRJvlwme9fzJJiSe6lG49NCxmuBiEdJAzPKaTPpK1cG1f1TqdgCfHR
HAL6Jxb3ylHG2LlTNFLXikubUi5RT6/9C8psr713Zm4HveCI/cx0WdgZ+fmsc9ft
rkIB1DdyV1kQ51m8r2rLi3J7aC5ggGOiex/VlGSd4e6SOQLpdQEdDbodtOJ4LgVk
+arFNCMinUWIOPGFzXhdm6lssPbh4MOBrz8c/M9TcF4hoMn5/o/9johZIZ/DOvXt
ZQIDAQABo1AwTjAdBgNVHQ4EFgQUOnWiddgeZj17IeysatqhE361o5YwHwYDVR0j
BBgwFoAUOnWiddgeZj17IeysatqhE361o5YwDAYDVR0TBAUwAwEB/zANBgkqhkiG
9w0BAQsFAAOCAQEAfh3cnOszHM/5wXjY7BIkmgDOReksS+87ibhBz7T2ddZj+yCF
9GdIBzXCiHpQFXpW8a3kc3I7l3nGfMTkmF6ld3ot/6SXP17QKJwxtvUA4ib8QkWD
S7FT+UcHCUHv42Sh1e5uAlQ5pMSul7iKcR7jwlwZGZ0905HOqrmdyUGJ+Ud2uZWD
AC0dJF6Bv9VhNtci8Imp05PaPH09deXLZu8LRrKRZFy9qLW5R6Swv7nzxckOAqDk
TTc40xwvQROekWUyxeJL7xaHuylUHE0bxsiIfx5bZsBizRjprIwGlj85CSPuTZyP
DPfaiZAN/61h5HNAnxLltOZfqabKYYw7l9LBDg==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDOST00lQfs8tJA
rhFGsZBjl1Wx+2SGcqH0eEjrKueWjBYgM9LU9kc5T/mBDvE9Q8Z0pBXlLvHzBE6c
BpcSKbE3gZ1DDbmkaFEm+XCZ71/MkmJJ7qUbj00LGa4GIR0kDM8ppM+krVwbV/VO
p2AJ8dEcAvonFvfKUcbYuVM0UteKS5tSLlFPr/0LymyvvXdmbge94Ij9zHRZ2Bn5
+axz1+2uQgHUN3JXWRDnWbyvasuLcntoLmCAY6J7H9WUZJ3h7pI5Aul1AR0Nuh20
4nguBWT5qsU0IyKdRYg48YXNeF2bqWyw9uHgw4GvPxz8z1NwXiGgyfn+j/2OiFkh
n8M69e1lAgMBAAECggEAevPgTTT+0lYwx2h416ACJboP09O5KQGuUl5XaAPcoTjB
/1OkOFbKQPjQCAJ1+0QoR2F9w2plv6kziX/MD4FWJXVV3J+TpNCgfhBy8u1gNjiR
6Osa8gBJtXIK7ZBTJCeWWoXnVYoWuh2FEupkLck6D+4eV6oy6x4u3QIo+6jc24n9
dIXQG6/v/Iao34kB0LUdp/4WNaUDvfI6NDhEwchpKE95dtWIDlIN/YhfiYAdjrnl
YmH2VDbAGgsdEiHP4wLZfjgsGPPDGS0+qBHoSiJGH0E6wWEZdAE4TsYGRFsO86n3
LfjEPFGfPlcnZe2cTTe3kmyKq/DTjxtu2rh3I8o2CQKBgQD/5Xe7cenaOBefzPlx
GOEsB+qv49UmzEPOXDNZe9hmAawuuuxPUM+xlE++P+mEgQm1LPT4WWgtFLPVuwmx
ncxt4CJNZh+ZGFyAZ4dm4M4ZhIBXNonyIP+yGyAJUUVF9Iy3TYcJNiGzv2Rx9JRQ
XWJMQnTDILmZbmU+ltTea7/zqwKBgQDOXqCqb17MuLt7OcKWSgthm79OlaOdzZvl
i9qU6VzZKG7Axc5gA9yq6tHp3vWPI4bNdvwqIIa/nzVILjGA5fcYFbRN+7gHwo8s
rNAgi5PAoKWqQRovyJRY9Eq/sn6l1jbJZAOUAMZMWDm8z89OqK7PNQSIAtfFSneo
2QxJkGeTLwKBgGJkafBB8af9b1/7YWISLepPNPbihH/BhMThAMGEdAVs2TaymtA4
g1OFck/1pSVUtFXcbmjbf8ntruQcYbLQuNz6lFXsUXP9QPwCUrbE85ouL2bZSps2
AvsJoPzUKe2nBUAp6CUrkjPaAJYsc6ae8X/fAaRRfeu33ef9+OV4yrq3AoGAYFZo
ZmfrN2Kdkt7Z6dLTEVPlsMfGQ6pyNmxdM9rkzzNC0JcGymfDIb7RE35T3+hTy6La
AMiCXv3xn6qAzY2NFh87tpPlyymWMOLTnf3Kkcfszlfp45idOBGCu46V9NDVbppT
2UmrSIR/H5dbTXsNcAlt/hhlpeInjhkU1VqmH10CgYEA7Kk+QhWq705SczpWjm5J
9kHqfFzJLwAWNBduiia0WypgPhLe/4wT1rYQkBtKMVKrgFo7Cvi4YKlrtlDnXyeU
CIFqfEL5NriQelqrFsvgHsmD+MpvDoSWm5C8IrTubtlNyWUzXSVT4OIwzPobzPqG
LILJ+e7bLw8RrM0HfgFnl8c=
-----END PRIVATE KEY-----
EOF
}

`, name)
}

var cdnDomainBasicMap = map[string]string{
	"domain_name": CHECKSET,
	"scope":       CHECKSET,
	"cdn_type":    "web",
}

const testServerCertificate = `-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n`
const testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n`

// Test Cdn Domain. >>> Resource test cases, automatically generated.
// Case 3176
func TestAccAliCloudCdnDomain_basic3176(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_domain_new.default"
	ra := resourceAttrInit(resourceId, AlicloudCdnDomainMap3176)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CdnServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("test%s%d.pfytlm.xyz", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCdnDomainBasicDependence3176)
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
					"cdn_type":    "web",
					"sources": []map[string]interface{}{
						{
							"type":     "ipaddr",
							"content":  "1.1.1.1",
							"priority": "20",
							"port":     "80",
							"weight":   "15",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": CHECKSET,
						"cdn_type":    "web",
						"sources.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope": "domestic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope": "domestic",
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url"},
			},
		},
	})
}

var AlicloudCdnDomainMap3176 = map[string]string{}

func AlicloudCdnDomainBasicDependence3176(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3176  twin
func TestAccAliCloudCdnDomain_basic3176_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_domain_new.default"
	ra := resourceAttrInit(resourceId, AlicloudCdnDomainMap3176)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CdnServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCdnDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("test%s%d.pfytlm.xyz", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCdnDomainBasicDependence3176)
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
					"scope":       "domestic",
					"domain_name": name,
					"cdn_type":    "web",
					"sources": []map[string]interface{}{
						{
							"type":     "ipaddr",
							"content":  "1.1.1.1",
							"priority": "20",
							"port":     "80",
							"weight":   "16",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":        "domestic",
						"domain_name":  CHECKSET,
						"cdn_type":     "web",
						"sources.#":    "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url"},
			},
		},
	})
}

// Test Cdn Domain. <<< Resource test cases, automatically generated.
