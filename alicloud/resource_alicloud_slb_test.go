package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_slb", &resource.Sweeper{
		Name: "alicloud_slb",
		F:    testSweepSLBs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_cs_cluster",
		},
	})
}

func testSweepSLBs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var slbs []slb.LoadBalancer
	req := slb.CreateDescribeLoadBalancersRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving SLBs: %s", err)
		}
		resp, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if resp == nil || len(resp.LoadBalancers.LoadBalancer) < 1 {
			break
		}
		slbs = append(slbs, resp.LoadBalancers.LoadBalancer...)

		if len(resp.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	service := SlbService{client}
	vpcService := VpcService{client}
	for _, loadBalancer := range slbs {
		name := loadBalancer.LoadBalancerName
		id := loadBalancer.LoadBalancerId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(loadBalancer.VpcId, loadBalancer.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping SLB: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting SLB: %s (%s)", name, id)
		if err := service.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlb_classictest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbClassicInstanceConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbClassicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"internet": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"address_ip_version":   "ipv4",
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
					"specification": "slb.s2.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "slb.s2.medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccSlbClassicInstanceConfigSpot%d_change", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSlbClassicInstanceConfigSpot%d_change", rand),
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tag_a": "1",
						"tag_b": "2",
						"tag_c": "3",
						"tag_d": "4",
						"tag_e": "5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"specification": "slb.s2.small",
					"internet":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"tags.%":               REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_ip_version": "ipv6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version": "ipv6",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpctest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbVpcInstanceConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbVpcConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"delete_protection": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"delete_protection":    "on",
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
					"delete_protection": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_protection": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "slb.s2.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "slb.s2.medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tag_a": "1",
						"tag_b": "2",
						"tag_c": "3",
						"tag_d": "4",
						"tag_e": "5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccSlbVpcInstanceConfigSpot%d_change", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccSlbVpcInstanceConfigSpot%d_change", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"specification":     "slb.s2.small",
					"delete_protection": "off",
					"address":           "172.16.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"tags.%":               REMOVEKEY,
						"address":              "172.16.0.1",
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"delete_protection":    "off",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpcmulti(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbVpcInstancemultiConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbVpcConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":         "10",
					"name":          name,
					"specification": "slb.s2.small",
					"vswitch_id":    "${alicloud_vswitch.default.id}",
					"tags": map[string]string{
						"tag_a": "1",
						"tag_b": "2",
						"tag_c": "3",
						"tag_d": "4",
						"tag_e": "5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"tags.%":               "5",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceSlbVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, SlbVpcCommonTestCase, name)
}

func resourceSlbClassicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	`, name)
}
