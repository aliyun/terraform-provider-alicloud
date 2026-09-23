package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_sae_ingress", &resource.Sweeper{
		Name: "alicloud_sae_ingress",
		F:    testSweepSaeIngress,
	})
}

func testSweepSaeIngress(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("%s:tftestacc", region),
	}

	request := make(map[string]*string)
	var response map[string]interface{}

	request["ContainCustom"] = StringPointer(strconv.FormatBool(true))

	action := "/pop/v1/sam/namespace/describeNamespaceList"
	response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		action = "/pop/v1/sam/ingress/IngressList"
		request["RegionId"] = StringPointer(client.RegionId)
		request["NamespaceId"] = StringPointer(item["NamespaceId"].(string))
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.Data.IngressList", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data.IngressList", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				app_name := ""
				if val, exist := item["Description"]; exist {
					app_name = val.(string)
				}
				if strings.Contains(strings.ToLower(app_name), strings.ToLower(prefix)) {
					skip = false
				}
			}

			if skip {
				log.Printf("[INFO] Skipping Ecs SnapShot Policy: %s (%s)", item["Description"], item["Id"])
				continue
			}
			sweeped = true
			action = "/pop/v1/sam/ingress/Ingress"
			request := map[string]*string{
				"IngressId": StringPointer(strconv.FormatFloat(item["Id"].(float64), 'f', 0, 64)),
			}
			response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs SnapShot Policy (%s (%v)): %s", item["Description"].(string), item["Id"].(float64), err)
			}
			if sweeped {
				// Waiting 30 seconds to ensure these Ecs SnapShot Policy have been deleted.
				time.Sleep(30 * time.Second)
			}
			log.Printf("[INFO] Delete Ecs SnapShot Policy success: %v ", item["Id"].(float64))
		}

	}

	return nil
}

func TestAccAlicloudSAEIngress_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_ingress.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEIngressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeIngress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssaeingress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEIngressBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_id":  "${alicloud_sae_namespace.default.id}",
					"slb_id":        "${alicloud_slb.default.id}",
					"listener_port": "80",
					"rules": []map[string]interface{}{
						{
							"app_id":         "${alicloud_sae_application.default.id}",
							"container_port": "443",
							"domain":         "www.alicloud.com",
							"app_name":       "${alicloud_sae_application.default.app_name}",
							"path":           "/",
						},
					},
					"default_rule": []map[string]interface{}{
						{
							"app_id":         "${alicloud_sae_application.default.id}",
							"container_port": "443",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slb_id":         CHECKSET,
						"namespace_id":   CHECKSET,
						"listener_port":  "80",
						"rules.#":        "1",
						"default_rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_rule": []map[string]interface{}{
						{
							"app_id":         "${alicloud_sae_application.default.id}",
							"container_port": "443",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_port": "443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_port": "443",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "ingress-sae-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "ingress-sae-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cert_id": "${alicloud_slb_server_certificate.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"app_name":       "${alicloud_sae_application.default.app_name}",
							"container_port": "443",
							"domain":         "www.sohu.com",
							"app_id":         "${alicloud_sae_application.default.id}",
							"path":           "/",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_rule": []map[string]interface{}{
						{
							"app_id":         "${alicloud_sae_application.default.id}",
							"container_port": "443",
						},
					},
					"listener_port": "443",
					"description":   "ingress-sae-test",
					"cert_id":       "${alicloud_slb_server_certificate.default.id}",
					"rules": []map[string]interface{}{
						{
							"app_name":       "${alicloud_sae_application.default.app_name}",
							"container_port": "443",
							"domain":         "www.sohu.com",
							"app_id":         "${alicloud_sae_application.default.id}",
							"path":           "/",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule.#": "1",
						"listener_port":  "443",
						"description":    "ingress-sae-test",
						"cert_id":        CHECKSET,
						"rules.#":        "1",
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

func TestAccAlicloudSAEIngress_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_ingress.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEIngressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeIngress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssaeingress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEIngressBasicDependence1)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"slb_id":            "${alicloud_alb_load_balancer.default.id}",
					"namespace_id":      "${alicloud_sae_namespace.default.id}",
					"listener_port":     "80",
					"load_balance_type": "alb",
					"listener_protocol": "HTTP",
					"rules": []map[string]interface{}{
						{
							"app_id":           "${alicloud_sae_application.default.id}",
							"container_port":   "8080",
							"domain":           "www.alicloud.com",
							"app_name":         "${alicloud_sae_application.default.app_name}",
							"path":             "/",
							"rewrite_path":     "/${1}",
							"backend_protocol": "http",
						},
					},
					"default_rule": []map[string]interface{}{
						{
							"app_id":         "${alicloud_sae_application.default.id}",
							"container_port": "8080",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slb_id":            CHECKSET,
						"namespace_id":      CHECKSET,
						"listener_port":     "80",
						"load_balance_type": "alb",
						"listener_protocol": "HTTP",
						"rules.#":           "1",
						"default_rule.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_port":     "443",
					"cert_ids":          "${alicloud_slb_server_certificate.default.id}",
					"listener_protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_port":     "443",
						"cert_ids":          CHECKSET,
						"listener_protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"app_name":         "${alicloud_sae_application.default.app_name}",
							"container_port":   "443",
							"domain":           "www.sohu.com",
							"app_id":           "${alicloud_sae_application.default.id}",
							"path":             "/",
							"rewrite_path":     "/${2}",
							"backend_protocol": "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_port":     "443",
					"cert_ids":          "${alicloud_slb_server_certificate.default.id}",
					"listener_protocol": "HTTPS",
					"description":       "ingress-sae-test",
					"rules": []map[string]interface{}{
						{
							"app_name":         "${alicloud_sae_application.default.app_name}",
							"container_port":   "443",
							"domain":           "www.sohu.com",
							"app_id":           "${alicloud_sae_application.default.id}",
							"path":             "/",
							"rewrite_path":     "/${2}",
							"backend_protocol": "https",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_port":     "443",
						"cert_ids":          CHECKSET,
						"listener_protocol": "HTTPS",
						"description":       "ingress-sae-test",
						"rules.#":           "1",
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

var AlicloudSAEIngressMap0 = map[string]string{
	"listener_port": CHECKSET,
	"rules.#":       CHECKSET,
	"namespace_id":  CHECKSET,
	"slb_id":        CHECKSET,
}

func AlicloudSAEIngressBasicDependence0(name string) string {
	config := fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	variable "namespace_id" {
  		default = "%s:tftestacc%d"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_sae_application" "default" {
  		app_description = var.name
  		app_name        = var.name
  		namespace_id    = alicloud_sae_namespace.default.id
  		image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/sae-demo-image/consumer:1.0"
  		package_type    = "Image"
  		vswitch_id      = data.alicloud_vswitches.default.vswitches.0.id
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		timezone        = "Asia/Beijing"
  		replicas        = "5"
  		cpu             = "500"
  		memory          = "2048"
	}

	resource "alicloud_slb" "default" {
  		load_balancer_name = var.name
  		load_balancer_spec = "slb.s2.small"
  		vswitch_id         = data.alicloud_vswitches.default.vswitches.0.id
	}

	resource "alicloud_sae_namespace" "default" {
  		namespace_description = var.name
  		namespace_id          = var.namespace_id
  		namespace_name        = var.name
	}

	resource "alicloud_slb_server_certificate" "default" {
  		name               = "slbservercertificate"
  		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  		private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
	}
`, name, defaultRegionToTest, acctest.RandIntRange(100, 999))
	return config
}

func AlicloudSAEIngressBasicDependence1(name string) string {
	config := fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	variable "namespace_id" {
  		default = "%s:tftestacc%d"
	}

	data "alicloud_alb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default_1" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_alb_zones.default.zones.0.id
	}

	data "alicloud_vswitches" "default_2" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_alb_zones.default.zones.1.id
	}

	resource "alicloud_sae_application" "default" {
  		app_description = var.name
  		app_name        = var.name
  		namespace_id    = alicloud_sae_namespace.default.id
  		image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/sae-demo-image/consumer:1.0"
  		package_type    = "Image"
  		vswitch_id      = data.alicloud_vswitches.default_1.vswitches.0.id
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		timezone        = "Asia/Beijing"
  		replicas        = "5"
  		cpu             = "500"
  		memory          = "2048"
	}

	resource "alicloud_alb_load_balancer" "default" {
  		vpc_id                 = data.alicloud_vpcs.default.ids.0
  		address_type           = "Internet"
  		address_allocated_mode = "Fixed"
  		load_balancer_name     = var.name
  		load_balancer_edition  = "Standard"
  		load_balancer_billing_config {
    		pay_type = "PayAsYouGo"
  		}
  		zone_mappings {
    		vswitch_id = data.alicloud_vswitches.default_1.ids[0]
    		zone_id    = data.alicloud_alb_zones.default.zones.0.id
  		}
  		zone_mappings {
    		vswitch_id = data.alicloud_vswitches.default_2.ids[0]
    		zone_id    = data.alicloud_alb_zones.default.zones.1.id
  		}
	}

	resource "alicloud_sae_namespace" "default" {
  		namespace_description = var.name
  		namespace_id          = var.namespace_id
  		namespace_name        = var.name
	}

	resource "alicloud_slb_server_certificate" "default" {
  		name               = "slbservercertificate"
  		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  		private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
	}
`, name, defaultRegionToTest, acctest.RandIntRange(100, 999))
	return config
}
