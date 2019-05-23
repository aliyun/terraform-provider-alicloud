package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_instance", &resource.Sweeper{
		Name: "alicloud_instance",
		F:    testSweepInstances,
		// When implemented, these should be removed firstly
		// Now, the resource alicloud_havip_attachment has been published.
		//Dependencies: []string{
		//	"alicloud_havip_attachment",
		//},
	})
}

func testSweepInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []ecs.Instance
	req := ecs.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Instances: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Instance: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Instance: %s (%s)", name, id)
		if v.DeletionProtection {
			request := ecs.CreateModifyInstanceAttributeRequest()
			request.InstanceId = id
			request.DeletionProtection = requests.NewBoolean(false)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				fmt.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
				continue
			}
		}
		if v.InstanceChargeType == string(PrePaid) {
			request := ecs.CreateModifyInstanceChargeTypeRequest()
			request.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), id))
			request.InstanceChargeType = string(PostPaid)
			request.IncludeDataDisks = requests.NewBoolean(true)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				fmt.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
				continue
			}
			time.Sleep(3 * time.Second)
		}

		req := ecs.CreateDeleteInstanceRequest()
		req.InstanceId = id
		req.Force = requests.NewBoolean(true)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 20 seconds to eusure these instances have been deleted.
		time.Sleep(20 * time.Second)
	}
	return nil
}

func TestAccAlicloudInstance_basic(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(*terraform.State) error {
		if instance.ZoneId == "" {
			return fmt.Errorf("bad availability zone")
		}
		if len(instance.SecurityGroupIds.SecurityGroupId) == 0 {
			return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"tf-testAccInstanceConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"security_enhancement_strategy",
						"Active"),
					testAccCheckSystemDiskSize("alicloud_instance.foo", 80),
				),
			},

			// test for multi steps
			{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"tf-testAccInstanceConfig"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_vpc(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigVPC(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"system_disk_category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_userData(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigUserData(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"system_disk_category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"internet_charge_type",
						"PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"user_data",
						"echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_multipleRegions(t *testing.T) {
	var instance ecs.Instance

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckInstanceDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigMultipleRegions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExistsWithProviders(
						"alicloud_instance.foo", &instance, &providers),
					testAccCheckInstanceExistsWithProviders(
						"alicloud_instance.bar", &instance, &providers),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_multiSecurityGroup(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(sgCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(instance.SecurityGroupIds.SecurityGroupId) < 0 {
				return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			if len(instance.SecurityGroupIds.SecurityGroupId) < sgCount {
				return fmt.Errorf("less security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_multiSecurityGroup(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(2),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"2"),
				),
			},
			{
				Config: testAccInstanceConfig_multiSecurityGroup_add(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(3),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"3"),
				),
			},
			{
				Config: testAccInstanceConfig_multiSecurityGroup_remove(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(1),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"1"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_multiSecurityGroupByCount(t *testing.T) {
	var instance ecs.Instance

	testCheck := func(sgCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(instance.SecurityGroupIds.SecurityGroupId) < 0 {
				return fmt.Errorf("no security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			if len(instance.SecurityGroupIds.SecurityGroupId) < sgCount {
				return fmt.Errorf("less security group: %#v", instance.SecurityGroupIds.SecurityGroupId)
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_multiSecurityGroupByCount(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					testCheck(3),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"security_groups.#",
						"3"),
				),
			},
		},
	})

}

func TestAccAlicloudInstance_NetworkInstanceSecurityGroups(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceNetworkInstanceSecurityGroups(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_tags(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceConfigTags(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo", "tags.foo", "bar"),
				),
			},

			{
				Config: testAccCheckInstanceConfigTagsUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo", "tags.%", "6"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo", "tags.bar5", "zzz"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceConfigOrigin(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"tf-testAccCheckInstanceConfigOrigin-foo"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"host_name",
						"host-foo"),
				),
			},

			{
				Config: testAccCheckInstanceConfigOriginUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"tf-testAccCheckInstanceConfigOrigin-bar"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"host_name",
						"host-bar"),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceImage_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceImageOrigin(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.update_image", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.update_image",
						"system_disk_size",
						"50"),
				),
			},

			{
				Config: testAccCheckInstanceImageUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.update_image", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.update_image",
						"system_disk_size",
						"60"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.update_image", "key_name", "tf-testAccCheckInstanceImageOrigin"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_associatePublicIP(t *testing.T) {
	var instance ecs.Instance

	testCheckPrivateIP := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			var privateIP string
			if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
				privateIP = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
			}
			if privateIP != "" {
				return nil
			}

			return fmt.Errorf("can't get private IP")
		}
	}

	testCheckPublicIP := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			publicIP := instance.PublicIpAddress.IpAddress[0]
			if publicIP == "" {
				return fmt.Errorf("can't get public IP")
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigAssociatePublicIP(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.foo", &instance),
					testCheckPrivateIP(),
					testCheckPublicIP(),
				),
			},
		},
	})
}

func TestAccAlicloudInstancePrivateIp_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstancePrivateIp(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"private_ip",
						"172.16.0.10"),
				),
			},

			{
				Config: testAccCheckInstancePrivateIpUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"private_ip",
						"172.16.1.10"),
				),
			},
		},
	})
}
func TestAccAlicloudInstanceSecurityEnhancementStrategy_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceSecurityEnhancementStrategy(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"security_enhancement_strategy",
						"Active"),
				),
			},

			{
				Config: testAccCheckInstanceSecurityEnhancementStrategyUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.private_ip", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.private_ip",
						"security_enhancement_strategy",
						"Deactive"),
				),
			},
		},
	})
}

// At present, One account only support at most 16 cpu core modify in one month.
func SkipTestAccAlicloudInstanceChargeType_post2Pre(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceChargeTypePostPaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PostPaid"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "period"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "period_unit"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "renewal_status"), // string(RenewNormal)),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "auto_renew_period"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "include_data_disks"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "dry_run"),
					resource.TestCheckNoResourceAttr("alicloud_instance.charge_type", "force_delete"),
				),
			},

			{
				Config: testAccCheckInstanceChargeTypePrePaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period_unit", string(Week)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "renewal_status", string(RenewNormal)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "auto_renew_period", "0"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "include_data_disks", "true"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "dry_run", "false"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "force_delete", "true"),
				),
			},

			{
				Config: testAccCheckInstanceChargeTypePostPaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period_unit", string(Week)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "renewal_status", string(RenewNormal)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "auto_renew_period", "0"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "include_data_disks", "true"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "dry_run", "false"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "force_delete", "true"),
				),
			},
		},
	})
}

// At present, One account only support at most 16 cpu core modify in one month.
func SkipTestAccAlicloudInstanceChargeType_pre2Post(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceChargeTypePrePaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period_unit", string(Week)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "renewal_status", string(RenewNormal)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "auto_renew_period", "0"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "include_data_disks", "true"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "dry_run", "false"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "force_delete", "true"),
				),
			},

			{
				Config: testAccCheckInstanceChargeTypePostPaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period_unit", string(Week)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "renewal_status", string(RenewNormal)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "auto_renew_period", "0"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "include_data_disks", "true"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "dry_run", "false"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "force_delete", "true"),
				),
			},

			{
				Config: testAccCheckInstanceChargeTypePrePaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.charge_type", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "period_unit", string(Week)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "renewal_status", string(RenewNormal)),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "auto_renew_period", "0"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "include_data_disks", "true"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "dry_run", "false"),
					resource.TestCheckResourceAttr("alicloud_instance.charge_type", "force_delete", "true"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_spot(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.EcsSpotNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: "alicloud_instance.spot",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSpotInstance(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.spot", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.spot",
						"spot_strategy", "SpotWithPriceLimit"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.spot",
						"spot_price_limit", "1.002"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudInstanceType_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceType(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
					resource.TestMatchResourceAttr("alicloud_instance.type", "instance_type", regexp.MustCompile("^ecs.t5-[a-z0-9]{1,}.nano")),
				),
			},

			{
				Config: testAccCheckInstanceTypeUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
					resource.TestMatchResourceAttr("alicloud_instance.type", "instance_type", regexp.MustCompile("^ecs.t5-[a-z0-9]{1,}.small")),
				),
			},

			{
				Config: testAccCheckInstanceTypePrepaid(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
					resource.TestMatchResourceAttr("alicloud_instance.type", "instance_type", regexp.MustCompile("^ecs.t5-[a-z0-9]{1,}.small")),
				),
			},

			{
				Config: testAccCheckInstanceTypePrepaidUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.type", &instance),
					resource.TestMatchResourceAttr("alicloud_instance.type", "instance_type", regexp.MustCompile("^ecs.t5-[a-z0-9]{1,}.large")),
				),
			},
		},
	})
}

func TestAccAlicloudInstanceNetworkSpec_update(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceNetworkSpec(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.network", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_in", "-1"),
				),
			},

			{
				Config: testAccCheckInstanceNetworkSpecUpdate(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.network", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_out", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.network",
						"internet_max_bandwidth_in", "50"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_ramrole(t *testing.T) {
	var instance ecs.Instance
	rand := acctest.RandIntRange(100000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.role",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceRamRole(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("alicloud_instance.role", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.role", "role_name",
						fmt.Sprintf("tf-testAccCheckInstanceRamRole-%d", rand)),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_dataDisk(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigDataDisk(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"system_disk_category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.#",
						"2"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.name",
						"disk1"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.size",
						"20"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.description",
						"disk1"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.encrypted",
						"false"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.snapshot_id",
						""),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.0.delete_with_instance",
						"true"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.name",
						"disk2"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.size",
						"20"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.description",
						"disk2"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.encrypted",
						"false"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.snapshot_id",
						""),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"data_disks.1.delete_with_instance",
						"true"),
				),
			},
		},
	})
}

func TestAccAlicloudInstance_deletionProtection(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceEnableDeletionProtection(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"deletion_protection",
						"true")),
			},
			{
				Config: testAccCheckInstanceDisableDeletionProtection(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"deletion_protection",
						"false")),
			},
		},
	})
}

//testAccCheckInstanceVolumeTags
func TestAccAlicloudInstance_volumeTags(t *testing.T) {
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceVolumeTags(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.%",
						"1"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.tag1",
						"test")),
			},
			{
				Config: testAccCheckInstanceVolumeTags_update(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.%",
						"2"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.tag1",
						"test1"),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.tag2",
						"test2")),
			},
			{
				Config: testAccCheckInstanceVolumeTags_delete(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_instance.foo",
						"volume_tags.%",
						"0")),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, i *ecs.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckInstanceExistsWithProviders(n, i, &providers)
}

func testAccCheckInstanceExistsWithProviders(n string, i *ecs.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			ecsService := EcsService{client}
			instance, err := ecsService.DescribeInstance(rs.Primary.ID)
			log.Printf("[WARN]get ecs instance %#v", instance)
			// Verify the error is what we want
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err

			}

			*i = instance
			return nil
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckInstanceDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance" {
			continue
		}

		// Try to find the resource
		instance, err := ecsService.DescribeInstance(rs.Primary.ID)
		if err == nil {
			if instance.Status != "" && instance.Status != string(Stopped) {
				return fmt.Errorf("Found unstopped instance: %s", instance.InstanceId)
			}
		}

		// Verify the error is what we want
		if NotFoundError(err) {
			continue
		}

		return err
	}

	return nil
}

func testAccCheckSystemDiskSize(n string, size int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providers := []*schema.Provider{testAccProvider}
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		for _, provider := range providers {
			if provider.Meta() == nil {
				continue
			}
			client := provider.Meta().(*connectivity.AliyunClient)
			ecsService := EcsService{client}
			systemDisk, err := ecsService.QueryInstanceSystemDisk(rs.Primary.ID)
			if err != nil {
				log.Printf("[ERROR]get system disk size error: %#v", err)
				return err
			}

			if systemDisk.Size != size {
				return fmt.Errorf("system disk size not equal %d, the instance system size is %d",
					size, systemDisk.Size)
			}
		}

		return nil
	}
}

const testAccInstanceConfig = `
data "alicloud_zones" "default" {
	 available_disk_category = "cloud_ssd"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
	name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "tf-testAccInstanceConfig"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}"
	description = "bar"
}

resource "alicloud_instance" "foo" {
	image_id = "${data.alicloud_images.default.images.0.id}"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

	system_disk_category = "cloud_ssd"
	system_disk_size = 80

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
	security_enhancement_strategy = "Active"
    deletion_protection = false

	tags {
		foo = "bar"
		work = "test"
	}
}
`

func testAccInstanceConfigVPC(common string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "tf-testAccEcsInstanceConfigVPC"
	}
	resource "alicloud_instance" "foo" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"

		# series III
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"

		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
	}`, common)
}

func testAccInstanceConfigUserData(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEcsInstanceConfigUserData"
	}
	resource "alicloud_instance" "foo" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		# series III
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
	}
	`, common)
}

const testAccInstanceConfigMultipleRegions = `
provider "alicloud" {
	alias = "beijing"
	region = "cn-beijing"
}
provider "alicloud" {
	alias = "shanghai"
	region = "cn-shanghai"
}
data "alicloud_zones" "default" {
	provider = "alicloud.beijing"
	available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "default" {
	provider = "alicloud.beijing"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
	provider = "alicloud.beijing"
	name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
data "alicloud_zones" "sh" {
	provider = "alicloud.shanghai"
	 available_disk_category = "cloud_efficiency"
}
data "alicloud_instance_types" "sh" {
	provider = "alicloud.shanghai"
 	availability_zone = "${data.alicloud_zones.sh.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "sh" {
	provider = "alicloud.shanghai"
	name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "tf-testAccInstanceConfigMultipleRegions"
}

resource "alicloud_vpc" "vpc_foo" {
  name = "${var.name}"
  provider = "alicloud.beijing"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc_bar" {
  name = "${var.name}"
  provider = "alicloud.shanghai"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vsw_foo" {
  provider = "alicloud.beijing"
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc_foo.id}"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "vsw_bar" {
  provider = "alicloud.shanghai"
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc_bar.id}"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.sh.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	provider = "alicloud.beijing"
	description = "foo"
	vpc_id = "${alicloud_vpc.vpc_foo.id}"
}

resource "alicloud_security_group" "tf_test_bar" {
	name = "${var.name}"
	provider = "alicloud.shanghai"
	description = "bar"
	vpc_id = "${alicloud_vpc.vpc_bar.id}"
}

resource "alicloud_instance" "foo" {
  	# cn-beijing
  	provider = "alicloud.beijing"
  	image_id = "${data.alicloud_images.default.images.0.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"

  	internet_charge_type = "PayByBandwidth"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	system_disk_category = "cloud_efficiency"
  	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.vsw_foo.id}"
}

resource "alicloud_instance" "bar" {
	# cn-shanghai
	provider = "alicloud.shanghai"
  	image_id = "${data.alicloud_images.sh.images.0.id}"
	availability_zone = "${data.alicloud_zones.sh.zones.0.id}"

	internet_charge_type = "PayByBandwidth"

	instance_type = "${data.alicloud_instance_types.sh.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.tf_test_bar.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.vsw_bar.id}"
}
`

func testAccInstanceConfig_multiSecurityGroup(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf_testAccInstanceConfig_multiSecurityGroup"
	}

	resource "alicloud_security_group" "tf_test_foo" {
		name = "${var.name}-foo"
		description = "foo"
		vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"

		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		security_groups = ["${alicloud_security_group.tf_test_foo.id}", "${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, common)
}

func testAccInstanceConfig_multiSecurityGroup_add(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf_testAccInstanceConfig_multiSecurityGroup"
	}
	resource "alicloud_security_group" "tf_test_foo" {
		name = "${var.name}-foo"
		description = "foo"
		vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_security_group" "tf_test_bar" {
		name = "${var.name}-bar"
		description = "bar"
	        vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"

		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		security_groups = ["${alicloud_security_group.tf_test_foo.id}", "${alicloud_security_group.default.id}",
			"${alicloud_security_group.tf_test_bar.id}"]
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccInstanceConfig_multiSecurityGroup_remove(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf_testAccInstanceConfig_multiSecurityGroup"
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"

		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccInstanceConfig_multiSecurityGroupByCount(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf_testAccInstanceConfig_multiSecurityGroupByCount"
	}
	resource "alicloud_security_group" "tf_test_foo" {
		name = "${var.name}"
		count = 2
		description = "foo"
  		vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		security_groups = ["${alicloud_security_group.default.id}", "${alicloud_security_group.tf_test_foo.*.id}"]
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccInstanceNetworkInstanceSecurityGroups(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccInstanceNetworkInstanceSecurityGroups"
	}

	resource "alicloud_instance" "foo" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"

		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"

		internet_max_bandwidth_out = 5
		internet_charge_type = "PayByTraffic"
	}
	`, common)
}
func testAccCheckInstanceConfigTags(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceConfigTags"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"

		tags {
			foo = "bar"
			bar = "foo"
		}
	}
	`, common)
}

func testAccCheckInstanceConfigTagsUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceConfigTags"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"

		tags {
			bar1 = "zzz"
			bar2 = "bar"
			bar3 = "bar"
			bar4 = "bar"
			bar5 = "zzz"
			bar6 = "bar"
		}
	}
	`, common)
}
func testAccCheckInstanceConfigOrigin(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceConfigOrigin"
	}

	resource "alicloud_security_group_rule" "http-in" {
		  type = "ingress"
		  ip_protocol = "tcp"
		  nic_type = "intranet"
		  policy = "accept"
		  port_range = "80/80"
		  priority = 1
		  security_group_id = "${alicloud_security_group.default.id}"
		  cidr_ip = "0.0.0.0/0"
	}

	resource "alicloud_security_group_rule" "ssh-in" {
		  type = "ingress"
		  ip_protocol = "tcp"
		  nic_type = "intranet"
		  policy = "accept"
		  port_range = "22/22"
		  priority = 1
		  security_group_id = "${alicloud_security_group.default.id}"
		  cidr_ip = "0.0.0.0/0"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_groups = ["${alicloud_security_group.default.id}"]

		instance_name = "${var.name}-foo"
		host_name = "host-foo"
	}
	`, common)
}

func testAccCheckInstanceConfigOriginUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceConfigOrigin"
	}

	resource "alicloud_security_group_rule" "http-in" {
		  type = "ingress"
		  ip_protocol = "tcp"
		  nic_type = "intranet"
		  policy = "accept"
		  port_range = "80/80"
		  priority = 1
		  security_group_id = "${alicloud_security_group.default.id}"
		  cidr_ip = "0.0.0.0/0"
	}

	resource "alicloud_security_group_rule" "ssh-in" {
		  type = "ingress"
		  ip_protocol = "tcp"
		  nic_type = "intranet"
		  policy = "accept"
		  port_range = "22/22"
		  priority = 1
		  security_group_id = "${alicloud_security_group.default.id}"
		  cidr_ip = "0.0.0.0/0"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		internet_charge_type = "PayByTraffic"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_groups = ["${alicloud_security_group.default.id}"]

		instance_name = "${var.name}-bar"
		host_name = "host-bar"
	}
	`, common)
}
func testAccInstanceConfigAssociatePublicIP(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccInstanceConfigAssociatePublicIP"
	}

	resource "alicloud_instance" "foo" {
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		internet_max_bandwidth_out = 5
		internet_charge_type = "PayByTraffic"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		instance_name = "${var.name}"
	}
	`, common)
}
func testAccCheckInstanceImageOrigin(common string) string {
	return fmt.Sprintf(`
	%s
	data "alicloud_images" "centos" {
		most_recent = true
		owners = "system"
		name_regex = "^centos_6\\w{1,5}[64]{1}.*"
	}
	variable "name" {
		default = "tf-testAccCheckInstanceImageOrigin"
	}

	resource "alicloud_instance" "update_image" {
		image_id = "${data.alicloud_images.centos.images.0.id}"
		  system_disk_category = "cloud_efficiency"
		  system_disk_size = 50
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		  instance_name = "${var.name}"
		  password = "Yourpassword1234"
		  security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_key_pair" "key" {
		  key_name = "${var.name}"
	}

	resource "alicloud_key_pair_attachment" "atta" {
		  key_name = "${alicloud_key_pair.key.key_name}"
		  instance_ids = ["${alicloud_instance.update_image.id}"]
	}
	`, common)
}
func testAccCheckInstanceImageUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	data "alicloud_images" "ubuntu" {
		most_recent = true
		owners = "system"
		name_regex = "^ubuntu_14\\w{1,5}[64]{1}.*"
	}
	variable "name" {
		default = "tf-testAccCheckInstanceImageOrigin"
	}

	resource "alicloud_instance" "update_image" {
		image_id = "${data.alicloud_images.ubuntu.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 60
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		password = "Yourpassword1234"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_key_pair" "key" {
		  key_name = "${var.name}"
	}

	resource "alicloud_key_pair_attachment" "atta" {
		  key_name = "${alicloud_key_pair.key.key_name}"
		  instance_ids = ["${alicloud_instance.update_image.id}"]
	}
	`, common)
}

func testAccCheckInstancePrivateIp(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstancePrivateIp"
	}

	resource "alicloud_vswitch" "foo" {
		vpc_id = "${alicloud_vpc.default.id}"
		cidr_block = "172.16.1.0/24"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		name = "${var.name}-foo"
	}
	resource "alicloud_instance" "private_ip" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
	}
	`, common)
}

func testAccCheckInstancePrivateIpUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstancePrivateIp"
	}
	resource "alicloud_vswitch" "foo" {
		vpc_id = "${alicloud_vpc.default.id}"
		cidr_block = "172.16.1.0/24"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		name = "${var.name}-foo"
	}

	resource "alicloud_instance" "private_ip" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.foo.id}"
		private_ip = "172.16.1.10"
	}
	`, common)
}

func testAccCheckInstanceSecurityEnhancementStrategy(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceSecurityEnhancementStrategy"
	}
	resource "alicloud_instance" "private_ip" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		security_enhancement_strategy = "Active"
	}
	`, common)
}

func testAccCheckInstanceSecurityEnhancementStrategyUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceSecurityEnhancementStrategy"
	}
	resource "alicloud_instance" "private_ip" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		security_enhancement_strategy = "Deactive"
	}
	`, common)
}

func testAccCheckInstanceChargeTypePostPaid(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceChargeType"
	}

	data "alicloud_zones" "special" {
	  	available_disk_category     = "cloud_efficiency"
	  	available_resource_creation = "VSwitch"
	  	available_instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
	}

	data "alicloud_instance_types" "special" {
	  	cpu_core_count    = 2
	  	memory_size       = 4
	  	instance_type_family = "ecs.t5"
	}

	resource "alicloud_vswitch" "special" {
		vpc_id            = "${alicloud_vpc.default.id}"
		cidr_block        = "172.16.1.0/24"
		availability_zone = "${data.alicloud_zones.special.zones.0.id}"
		name              = "${var.name}"
	}

	resource "alicloud_instance" "charge_type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.special.id}"
		instance_charge_type = "PostPaid"
	}
	`, common)
}

func testAccCheckInstanceChargeTypePrePaid(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceChargeType"
	}

	data "alicloud_zones" "special" {
	  	available_disk_category     = "cloud_efficiency"
	  	available_resource_creation = "VSwitch"
	  	available_instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
	}

	data "alicloud_instance_types" "special" {
	  	cpu_core_count    = 2
	  	memory_size       = 4
	  	instance_type_family = "ecs.t5"
	}

	resource "alicloud_vswitch" "special" {
	  	vpc_id            = "${alicloud_vpc.default.id}"
	  	cidr_block        = "172.16.1.0/24"
	  	availability_zone = "${data.alicloud_zones.special.zones.0.id}"
	  	name              = "${var.name}"
	}

	resource "alicloud_instance" "charge_type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.special.id}"
		instance_charge_type = "PrePaid"
		period_unit = "Week"
		force_delete = "true"
	}
	`, common)
}

func testAccCheckSpotInstance(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckSpotInstance"
	}
	data "alicloud_instance_types" "special" {
	  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  	cpu_core_count    = 2
	  	memory_size       = 4
	  	spot_strategy = "SpotWithPriceLimit"
	}
	resource "alicloud_instance" "spot" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		spot_strategy = "SpotWithPriceLimit"
		spot_price_limit = "1.002"
	}
	`, common)
}

func testAccCheckInstanceType(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${alicloud_vswitch.new.availability_zone}"
		cpu_core_count = 1
		memory_size = 0.5
		instance_type_family = "ecs.t5"
	}
	resource "alicloud_vswitch" "new" {
	  vpc_id            = "${alicloud_vpc.default.id}"
	  cidr_block        = "172.16.1.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
	  name              = "${var.name}"
	}
	resource "alicloud_instance" "type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.new.id}"
	}
	`, common)
}

func testAccCheckInstanceTypeUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${alicloud_vswitch.new.availability_zone}"
		cpu_core_count = 1
		memory_size = 1
		instance_type_family = "ecs.t5"
	}
	resource "alicloud_vswitch" "new" {
	  vpc_id            = "${alicloud_vpc.default.id}"
	  cidr_block        = "172.16.1.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
	  name              = "${var.name}"
	}
	resource "alicloud_instance" "type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.new.id}"
	}
	`, common)
}

func testAccCheckInstanceTypePrepaid(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${alicloud_vswitch.new.availability_zone}"
		cpu_core_count = 1
		memory_size = 2
		instance_type_family = "ecs.t5"
	}
	resource "alicloud_vswitch" "new" {
	  vpc_id            = "${alicloud_vpc.default.id}"
	  cidr_block        = "172.16.1.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
	  name              = "${var.name}"
	}
	resource "alicloud_instance" "type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
		instance_name = "${var.name}"
		instance_charge_type = "PrePaid"
		period_unit = "Week"
		force_delete = true
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.new.id}"
	}
	`, common)
}

func testAccCheckInstanceTypePrepaidUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${alicloud_vswitch.new.availability_zone}"
		cpu_core_count = 2
		memory_size = 4
		instance_type_family = "ecs.t5"
	}
	resource "alicloud_vswitch" "new" {
	  vpc_id            = "${alicloud_vpc.default.id}"
	  cidr_block        = "172.16.1.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
	  name              = "${var.name}"
	}
	resource "alicloud_instance" "type" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
		instance_name = "${var.name}"
		instance_charge_type = "PrePaid"
		period_unit = "Week"
		force_delete = true
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.new.id}"
	}
	`, common)
}

func testAccCheckInstanceNetworkSpec(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceNetworkSpec"
	}

	resource "alicloud_instance" "network" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccCheckInstanceNetworkSpecUpdate(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceNetworkSpec"
	}
	resource "alicloud_instance" "network" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		internet_max_bandwidth_in = 50
	}
	`, common)
}
func testAccCheckInstanceRamRole(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceRamRole-%d"
	}

	resource "alicloud_instance" "role" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 60
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		password = "Yourpassword1234"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		role_name = "${alicloud_ram_role.role.name}"
	}

	resource "alicloud_ram_role" "role" {
		  name = "${var.name}"
		  services = ["ecs.aliyuncs.com"]
		  force = "true"
	}

	resource "alicloud_ram_policy" "policy" {
	  name = "${var.name}"
	  statement = [
		{
		  effect = "Allow"
		  action = ["CreateInstance"]
		  resource = ["*"]
		}
	  ]
	  force = "true"
	}

	resource "alicloud_ram_role_policy_attachment" "role-policy" {
	  policy_name = "${alicloud_ram_policy.policy.name}"
	  role_name = "${alicloud_ram_role.role.name}"
	  policy_type = "${alicloud_ram_policy.policy.type}"
	}
	`, common, rand)
}

func testAccInstanceConfigDataDisk(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceRamRole"
	}

	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 60
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		password = "Yourpassword1234"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"

		data_disks = [
		{
			name = "disk1"
			size = "20"
			category = "cloud_efficiency"
			description = "disk1"
		},
		{
			name = "disk2"
			size = "20"
			category = "cloud_efficiency"
			description = "disk2"
		}
		]

		tags {
			foo = "bar"
			work = "test"
		}
	}
	`, common)
}

func testAccCheckInstanceEnableDeletionProtection(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		deletion_protection = true
	}
	`, common)
}

func testAccCheckInstanceDisableDeletionProtection(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		deletion_protection = false
	}
	`, common)
}

func testAccCheckInstanceVolumeTags(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		volume_tags {
			tag1 = "test"
		}
	}
	`, common)
}

func testAccCheckInstanceVolumeTags_update(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		volume_tags {
			tag1 = "test1"
			tag2 = "test2"
		}
	}
	`, common)
}

func testAccCheckInstanceVolumeTags_delete(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckInstanceType"
	}
	data "alicloud_instance_types" "new" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
	}
	resource "alicloud_instance" "foo" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_efficiency"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		security_groups = ["${alicloud_security_group.default.id}"]
		vswitch_id = "${alicloud_vswitch.default.id}"
		volume_tags {
		}
	}
	`, common)
}
