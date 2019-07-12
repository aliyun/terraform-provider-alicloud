package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		"tf-testacc",
		"tf_testacc",
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

func TestAccAlicloudCdnDomainNew_basic(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainDependence)

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
					"scope":       "domestic",
					"sources": []map[string]interface{}{
						{
							"content": "www.aliyuntest.com",
							"type":    "oss",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope": "domestic",
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
							"content":  "www.aliyuntest.com",
							"type":     "oss",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.0.weight": "30",
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
						"sources.0.type": "domain",
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
						"sources.0.weight":  "10",
						"sources.0.content": "1.1.1.1",
						"sources.0.type":    "ipaddr",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "1.1.1.1",
							"type":     "ipaddr",
							"priority": "40",
							"port":     "80",
							"weight":   "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.0.priority": "40",
						"sources.0.weight":   "10",
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
					"domain_name":        name,
					"cdn_type":           "web",
					"scope":              REMOVEKEY,
					"certificate_config": REMOVEKEY,
					"sources": []map[string]interface{}{
						{
							"content": "www.aliyuntest.com",
							"type":    "oss",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":              REMOVEKEY,
						"sources.0.content":  "www.aliyuntest.com",
						"sources.0.type":     "oss",
						"sources.0.priority": "20",
						"sources.0.weight":   "10",

						"certificate_config.0.server_certificate_status": REMOVEKEY,
						"certificate_config.0.force_set":                 REMOVEKEY,
						"certificate_config.0.cert_name":                 REMOVEKEY,
						"certificate_config.0.cert_type":                 REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCdnDomainNew_scope(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCdnDomainDependence)

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
							"content": "www.aliyuntest.com",
							"type":    "oss",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope": "overseas",
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
							"content":  "www.aliyuntest.com",
							"type":     "oss",
							"priority": "20",
							"port":     "80",
							"weight":   "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.0.weight": "30",
					}),
				),
			},
		},
	})
}

func resourceCdnDomainDependence(name string) string {
	return ""
}

var cdnDomainBasicMap = map[string]string{
	"domain_name":        CHECKSET,
	"scope":              CHECKSET,
	"cdn_type":           "web",
	"sources.0.content":  "www.aliyuntest.com",
	"sources.0.type":     "oss",
	"sources.0.priority": "20",
	"sources.0.weight":   "10",
	"sources.0.port":     "80",
}

const testServerCertificate = `-----BEGIN CERTIFICATE-----\nMIICrDCCAZQCCQDApyXUTYDE+DANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDDA0q\nLnhpYW96aHUuY29tMB4XDTE4MTIzMTEwMDY1OVoXDTE5MTIzMTEwMDY1OVowGDEW\nMBQGA1UEAwwNKi54aWFvemh1LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\nAQoCggEBAM51j0hM0dqY/gHed5LhsSbTShXmw2Sqg18avBpHup/C8YzvpVddvRUe\ndE54idcCfdCVhoWEtcNFXIkKpWDpImFtovoupdNdC2XoBHEC42or9ZA+wVsaCih+\nSRuB3r+yNzXHh+rmUa0FLeij/x0gUovmznUsk7UMMzwJLZwmuyi8LCuiTIlQzQ9R\nTdaYo2t4OuGVkdQiJzsYiRRNPOqCKtYvYcEBLalLFOcVn0aG/I9Fn0P3rc8fK9BE\nHaoYRunmusUCdCcKpisHHKYdtmd3Zgz+Z+PBkjARtufO6kOXSov6u7azLQxZgZJm\neagNkqQ4/R+9b4GQ6crj0Pi655QWZVECAwEAATANBgkqhkiG9w0BAQsFAAOCAQEA\nJAWcQb7942jWDWjFqW5C4eyJBuJxK8dsTsYpOa0ccpp/A2cVCTNMzV/sMCGROy2B\nuhrTMG4q3QawFwV+K4nQX8yJUJp57zyUfdit+yOD4hDzPyF5MOlOwzdqeg9Y9ODL\nSc2O6J3tFsx2332Hb7RXFlYodEk1uFezrj4YLVrAgk1QojaBHpuFiA/O3eCmErjW\nc88bcil60qMvlhBhCWaZzivji6oc3y/6qWRZ3apxp9/6sjOvlm5Q+wwoLHXaU+L6\nw7TZptw6upuSrAS+N4Rwqgn5TfPvbGkdRG0X0TDzE1+F+w377kE71nV8lmi8F786\nbqloD1gWJER071WRz1coHQ==\n-----END CERTIFICATE-----\n`
const testPrivateKey = `-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDOdY9ITNHamP4B\n3neS4bEm00oV5sNkqoNfGrwaR7qfwvGM76VXXb0VHnROeInXAn3QlYaFhLXDRVyJ\nCqVg6SJhbaL6LqXTXQtl6ARxAuNqK/WQPsFbGgoofkkbgd6/sjc1x4fq5lGtBS3o\no/8dIFKL5s51LJO1DDM8CS2cJrsovCwrokyJUM0PUU3WmKNreDrhlZHUIic7GIkU\nTTzqgirWL2HBAS2pSxTnFZ9GhvyPRZ9D963PHyvQRB2qGEbp5rrFAnQnCqYrBxym\nHbZnd2YM/mfjwZIwEbbnzupDl0qL+ru2sy0MWYGSZnmoDZKkOP0fvW+BkOnK49D4\nuueUFmVRAgMBAAECggEAGFrP2zSMsN/JXwkSS/Zpwm28WJcPR6nBs49gzyzU/BGw\nEvMWKxc4vewIxlT71axKkTeCVe/QzUc6YkQqPCNkVd/sEN093JAmTxAurfIsR5MF\n9c0hXBDXT+2NzDvmvfBVCPgPtYsT6Xgp8T6fUp1Ef5JrmnD2v62/wX5Hrhr3ixdp\nFoEJ92/FJN+Dt/Duri481+RYqUUbNqiZAuEMsny/hM+aijtpQgZK8FqnVfjIOo0l\np4TwfDoxtTwkj0QeJxelX/sNSD+iQ5xVmnxQ9Q1qmKMk0M1Wl8318ovcrXDTN3RZ\nR/Fp+lNw0aC4OLuOJDc7A13l5unHmex4Ka1kbjOZwQKBgQD8w0hJiVVmiojJ3/YX\n06rG3ufGlKXu4qvRL6eXNqnDp8AnKXpjKiDV5lGz+vK0pcAyF5AhpddDN20Un+nR\n4VU7NslNJdBfjQxfEKJRkGBP22C+pdlxOBDgdgouN9q7TTr9StaQC0KF+HUcF8fO\nsYpEp1R+27iHie+F6uS78ur0mQKBgQDRGnbv8acTDTBW8nll1wTRa0+iVR2yxglS\nndC97HzYHOFb+caVmaM18QzHaaRvAhQ0bZXOjFV0QVv1nytlmHxsMlG/N6LDz5jA\nNLGvbrzDzki+RUk0TWJuYpvdevsnwaCxMwyY6CQ5MKhKAvwgE8AhdAq5MGcE+6U/\noAOZNyOxeQKBgQCiCB2i5mLUpSIjJ2r+wzXK3sH9zvTAOpaiNsZcbTJOto67jB9k\nynDaLhdaJRjJLSgT9H700vc3o6RNgGXHoYedufU5e3AkkKrJlkQ3vTHAf4V5MaA+\nsA5BlenYzv1s7IlQLlV1aYJvl2Kba7MukSlt8UZ9PCUC3i2pz3Zp9cMgoQKBgQCK\nGkR7bMq/1nIausJa9IwGFC3gNP8MV6dInVqEVXCO+2QL7wetPm+A7NdXzPoBJwpZ\nJhdO93ho89Hcg2eSDgf/HazH8eLaGH32U9cW2rhpShDZOcGDfaiI5y+yM8s1Erki\nz2h+hLOH4g8D8ry6ItE+RvneHY2syNb3EqPNyZEVYQKBgGeC4KiyTpW+SqmkDRU/\n/FxIAlZE2vWmtJptAUkmLEVoBi0/GvJoYYJR4Jhj4gS9tJcrI860X0j05NZHvIX8\nUJ/vZZeiFRiatXYUtKbD/ngoOXTF92ew2J3eCwNHH5OKwEO3aejaYwD9m0QM4kTl\nVLWIf79+G3NrxIeqPnqv+l3J\n-----END PRIVATE KEY-----`
