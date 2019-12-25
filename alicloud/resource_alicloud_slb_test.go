package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	k8sPrefix := "kubernetes"

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

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	service := SlbService{client}
	vpcService := VpcService{client}
	csService := CsService{client}
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
		// If a slb tag key has prefix "kubernetes", this is a slb for k8s cluster and it should be deleted if cluster not exist.
		if skip {
			for _, t := range loadBalancer.Tags.Tag {
				if strings.HasPrefix(strings.ToLower(t.TagKey), strings.ToLower(k8sPrefix)) {
					_, err := csService.DescribeCsKubernetes(name)
					if NotFoundError(err) {
						skip = false
					} else {
						skip = true
						break
					}
				}
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
						"address_type":         "internet",
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
				Config: testAccConfig(map[string]interface{}{
					"name":         name,
					"internet":     REMOVEKEY,
					"address_type": "intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"address_type":         "intranet",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"address_ip_version":   "ipv4",
						"resource_group_id":    CHECKSET,
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
						"tag_A1": "value_A1",
						"tag_B2": "value_B2",
						"tag_C3": "value_C3",
						"tag_D4": "value_D4",
						"tag_E5": "value_E5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":      "5",
						"tags.tag_A1": "value_A1",
						"tags.tag_B2": "value_B2",
						"tags.tag_C3": "value_C3",
						"tags.tag_D4": "value_D4",
						"tags.tag_E5": "value_E5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"specification": "slb.s2.small",
					"address_type":  REMOVEKEY,
					"internet":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"address_type":         "internet",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"tags.%":               REMOVEKEY,
						"tags.tag_A1":          REMOVEKEY,
						"tags.tag_B2":          REMOVEKEY,
						"tags.tag_C3":          REMOVEKEY,
						"tags.tag_D4":          REMOVEKEY,
						"tags.tag_E5":          REMOVEKEY,
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
					"vswitch_id":        alicloud_vswitch.default.id,
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
						"resource_group_id":    CHECKSET,
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
						"tag_A1": "value_A1",
						"tag_B2": "value_B2",
						"tag_C3": "value_C3",
						"tag_D4": "value_D4",
						"tag_E5": "value_E5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":      "5",
						"tags.tag_A1": "value_A1",
						"tags.tag_B2": "value_B2",
						"tags.tag_C3": "value_C3",
						"tags.tag_D4": "value_D4",
						"tags.tag_E5": "value_E5",
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
					"vswitch_id":    alicloud_vswitch.default.id,
					"tags": map[string]string{
						"tag_A1": "value_A1",
						"tag_B2": "value_B2",
						"tag_C3": "value_C3",
						"tag_D4": "value_D4",
						"tag_E5": "value_E5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"tags.%":               "5",
						"tags.tag_A1":          "value_A1",
						"tags.tag_B2":          "value_B2",
						"tags.tag_C3":          "value_C3",
						"tags.tag_D4":          "value_D4",
						"tags.tag_E5":          "value_E5",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"resource_group_id":    CHECKSET,
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
