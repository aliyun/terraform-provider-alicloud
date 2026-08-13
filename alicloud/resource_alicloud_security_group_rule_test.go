package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudECSSecurityGroupRuleBasic(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rac := resourceAttrCheckInit(
		resourceCheckInit(resourceId, &v, serviceFunc),
		resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleCidrIp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"description":              "abcd",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleAll(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abcd",
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

func TestAccAliCloudECSSecurityGroupRuleEgress(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":        "egress",
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"cidr_ip":     "182.254.11.243/32",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRule(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: hclSecurityGroupEgressRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7512",
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

func TestAccAliCloudECSSecurityGroupRuleMulti(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test.2"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_ip":                  "45.20.250.240/32",
						"source_security_group_id": REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAliCloudECSSecurityGroupRulePrefixList(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRulePrefixList)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRulePrefix(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "abc",
						"prefix_list_id": CHECKSET,
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

func TestAccAliCloudECSSecurityGroupRuleEgressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2408:4004:cc:400::/56",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupIngressRuleIpv6Map)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupIngressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2408:4004:cc:400::/56",
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

func TestAccAliCloudECSSecurityGroupRuleEgressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
						"description":  "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupIngressRuleIpv6Map)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupIngressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
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

func TestAccAliCloudECSSecurityGroupRuleEgressICMPv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "ingress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "-1/-1",
		"priority":     "1",
		"ipv6_cidr_ip": "::/0",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupEgressRuleICMPv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

// TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerIDCrossAccount verifies a real
// cross-account security group rule: Terraform manages the account 1 network and
// rule while an SDK-created account 2 peer supplies the literal security group ID.
// Both account IDs are verified with STS before any resource is created.
func TestAccAliCloudECSSecurityGroupRuleSourceGroupOwnerIDCrossAccount(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	resourceId := "alicloud_security_group_rule.test"
	egressResourceId := "alicloud_security_group_rule.egress"
	name := acctest.RandString(4)
	peer := testAccPrepareSecurityGroupRuleCrossAccountPeer(t, name)
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleSourceGroupOwnerIDCrossAccount(name, peer.securityGroupID, peer.accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                     "ingress",
						"ip_protocol":              "tcp",
						"nic_type":                 "intranet",
						"policy":                   "accept",
						"port_range":               "22/22",
						"priority":                 "1",
						"security_group_id":        CHECKSET,
						"source_security_group_id": CHECKSET,
						"source_group_owner_id":    peer.accountID,
						"cidr_ip":                  "",
					}),
					resource.TestCheckResourceAttr(egressResourceId, "type", "egress"),
					resource.TestCheckResourceAttr(egressResourceId, "source_group_owner_account", peer.accountID),
					resource.TestCheckResourceAttr(egressResourceId, "source_group_owner_id", peer.accountID),
				),
			},
			{
				Config:             hclSecurityGroupRuleSourceGroupOwnerIDCrossAccount(name, peer.securityGroupID, peer.accountID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccSecurityGroupRuleCrossAccountPeer struct {
	accountID       string
	vpcID           string
	securityGroupID string
}

func testAccPrepareSecurityGroupRuleCrossAccountPeer(t *testing.T, name string) *testAccSecurityGroupRuleCrossAccountPeer {
	t.Helper()
	ak1 := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_1"))
	sk1 := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_1"))
	ak2 := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2"))
	sk2 := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_2"))
	if ak1 == "" || sk1 == "" || ak2 == "" || sk2 == "" {
		t.Fatal("set ALICLOUD_ACCESS_KEY_1, ALICLOUD_SECRET_KEY_1, ALICLOUD_ACCESS_KEY_2, and ALICLOUD_SECRET_KEY_2 for the cross-account test")
	}

	t.Setenv("ALICLOUD_ACCESS_KEY", ak1)
	t.Setenv("ALICLOUD_SECRET_KEY", sk1)

	region := strings.TrimSpace(os.Getenv("ALICLOUD_REGION"))
	if region == "" {
		region = "cn-beijing"
		t.Setenv("ALICLOUD_REGION", region)
	}
	callerAccountID := func(label, ak, sk string) string {
		client, err := sdk.NewClientWithAccessKey(region, ak, sk)
		if err != nil {
			t.Fatalf("failed to build STS client for %s: %s", label, err)
		}
		request := requests.NewCommonRequest()
		request.Method = requests.POST
		request.Scheme = "https"
		request.Domain = "sts.aliyuncs.com"
		request.Version = "2015-04-01"
		request.Product = "Sts"
		request.ApiName = "GetCallerIdentity"
		request.TransToAcsRequest()
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			t.Fatalf("failed to get caller identity for %s: %s", label, err)
		}
		result := struct {
			AccountID string `json:"AccountId"`
		}{}
		if err := json.Unmarshal([]byte(response.GetHttpContentString()), &result); err != nil {
			t.Fatalf("failed to decode caller identity for %s: %s", label, err)
		}
		if strings.TrimSpace(result.AccountID) == "" {
			t.Fatalf("STS GetCallerIdentity returned an empty AccountId for %s", label)
		}
		return strings.TrimSpace(result.AccountID)
	}
	account1ID := callerAccountID("account 1", ak1, sk1)
	account2ID := callerAccountID("account 2", ak2, sk2)
	if account1ID == account2ID {
		t.Fatal("cross-account test requires the numbered credentials to belong to different accounts")
	}

	vpcClient, err := vpc.NewClientWithAccessKey(region, ak2, sk2)
	if err != nil {
		t.Fatalf("failed to build account 2 VPC client: %s", err)
	}
	ecsClient, err := ecs.NewClientWithAccessKey(region, ak2, sk2)
	if err != nil {
		t.Fatalf("failed to build account 2 ECS client: %s", err)
	}

	peer := &testAccSecurityGroupRuleCrossAccountPeer{accountID: account2ID}
	vpcRequest := vpc.CreateCreateVpcRequest()
	vpcRequest.ClientToken = buildClientToken("CreateVpc")
	vpcRequest.VpcName = fmt.Sprintf("tf-testAccSGRSrcOwnerCross%s-peer", name)
	vpcRequest.CidrBlock = "192.168.0.0/16"
	vpcResponse, err := vpcClient.CreateVpc(vpcRequest)
	if err != nil {
		t.Fatalf("failed to create account 2 peer VPC: %s", err)
	}
	peer.vpcID = vpcResponse.VpcId
	if peer.vpcID == "" {
		t.Fatal("account 2 CreateVpc returned an empty VpcId")
	}
	t.Cleanup(func() {
		if err := testAccCleanupSecurityGroupRulePeerVpc(vpcClient, peer.vpcID); err != nil {
			t.Errorf("failed to delete account 2 peer VPC: %s", err)
		}
	})
	if err := testAccWaitSecurityGroupRulePeerVpcAvailable(vpcClient, peer.vpcID); err != nil {
		t.Fatalf("failed waiting for account 2 peer VPC to become available: %s", err)
	}

	securityGroupRequest := ecs.CreateCreateSecurityGroupRequest()
	securityGroupRequest.ClientToken = buildClientToken("CreateSecurityGroup")
	securityGroupRequest.SecurityGroupName = fmt.Sprintf("tf-testAccSGRSrcOwnerCross%s-peer", name)
	securityGroupRequest.VpcId = peer.vpcID
	securityGroupResponse, err := ecsClient.CreateSecurityGroup(securityGroupRequest)
	if err != nil {
		t.Fatalf("failed to create account 2 peer security group: %s", err)
	}
	peer.securityGroupID = securityGroupResponse.SecurityGroupId
	if peer.securityGroupID == "" {
		t.Fatal("account 2 CreateSecurityGroup returned an empty SecurityGroupId")
	}
	t.Cleanup(func() {
		if err := testAccCleanupSecurityGroupRulePeerSecurityGroup(ecsClient, peer.securityGroupID); err != nil {
			t.Errorf("failed to delete account 2 peer security group: %s", err)
		}
	})

	return peer
}

func testAccWaitSecurityGroupRulePeerVpcAvailable(client *vpc.Client, vpcID string) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := vpc.CreateDescribeVpcsRequest()
		request.VpcId = vpcID
		response, err := client.DescribeVpcs(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(response.Vpcs.Vpc) == 0 {
			return resource.NonRetryableError(fmt.Errorf("VPC %s was not found", vpcID))
		}
		if response.Vpcs.Vpc[0].Status == "Available" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("VPC %s is in status %s", vpcID, response.Vpcs.Vpc[0].Status))
	})
}

func testAccCleanupSecurityGroupRulePeerSecurityGroup(client *ecs.Client, securityGroupID string) error {
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = securityGroupID
	if _, err := client.DeleteSecurityGroup(request); err != nil {
		if !NotFoundError(err) && !IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			return err
		}
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := ecs.CreateDescribeSecurityGroupsRequest()
		request.SecurityGroupId = securityGroupID
		response, err := client.DescribeSecurityGroups(request)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if len(response.SecurityGroups.SecurityGroup) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("security group %s still exists", securityGroupID))
	})
}

func testAccCleanupSecurityGroupRulePeerVpc(client *vpc.Client, vpcID string) error {
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := vpc.CreateDeleteVpcRequest()
		request.VpcId = vpcID
		_, err := client.DeleteVpc(request)
		if err == nil || NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidResource.NotFound", "InvalidVpcID.NotFound"}) {
			return nil
		}
		if IsExpectedErrors(err, []string{"IncorrectVpcStatus"}) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := vpc.CreateDescribeVpcsRequest()
		request.VpcId = vpcID
		response, err := client.DescribeVpcs(request)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidResource.NotFound", "InvalidVpcID.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if len(response.Vpcs.Vpc) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("VPC %s still exists", vpcID))
	})
}

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"security_group_rule_id":   CHECKSET,
	"cidr_ip":                  "",
}

var testAccCheckSecurityGroupIngressRuleIpv6Map = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "drop",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
	"cidr_ip":           "",
}

var testAccCheckSecurityGroupRulePrefixList = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "accept",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_rule" {
			continue
		}
		_, err := ecsService.DescribeSecurityGroupRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
	}

	return nil
}

func hclSecurityGroupRuleBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRBase%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "drop"
  port_range               = "22/22"
  priority                 = 100
  security_group_id        = alicloud_security_group.test.0.id
  source_security_group_id = alicloud_security_group.test.1.id
  description              = "abc"
}
`, name)
}

func hclSecurityGroupRulePrefix(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRPrefix%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.vpcs.0.id
  security_group_name = var.name
}

resource "alicloud_ecs_prefix_list" "test" {
  address_family   = "IPv4"
  max_entries      = 2
  prefix_list_name = "tftest"
  description      = "description"
  entry {
    cidr        = "192.168.0.0/24"
    description = "description"
  }
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  prefix_list_id    = alicloud_ecs_prefix_list.test.id
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.id
  description       = "abc"
}
`, name)
}

func hclSecurityGroupRuleCidrIp(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-CIRDIP%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-desc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "description"
}
`, name)
}

func hclSecurityGroupRuleAll(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR_all%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRMulti%s"
}

variable "cidr_ip_list" {
  type    = list(string)
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  count             = length(compact(var.cidr_ip_list))
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  cidr_ip           = element(var.cidr_ip_list, count.index)
}
`, name)
}

func hclSecurityGroupIngressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
}
`, name)
}

func hclSecurityGroupIngressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
}
`, name)
}

func hclSecurityGroupEgressRule(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgress%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressDesc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
  inner_access_policy = "Accept"
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7512"
}
`, name)
}

func hclSecurityGroupEgressRuleICMPv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "icmpv6"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "::/0"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupRuleSourceGroupOwnerIDCrossAccount(name string, peerSecurityGroupID string, peerAccountID string) string {
	resourceName := fmt.Sprintf("tf-testAccSGRSrcOwnerCross%s", name)
	return fmt.Sprintf(`
# Terraform owns only account 1 resources. The account 2 peer is created by the
# Go test with the numbered account 2 SDK client and passed here by literal ID.
resource "alicloud_vpc" "own" {
  vpc_name   = "%s-own"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_security_group" "test" {
  vpc_id              = alicloud_vpc.own.id
  security_group_name = "%s"
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "accept"
  port_range               = "22/22"
  priority                 = 1
  security_group_id        = alicloud_security_group.test.id
  source_security_group_id = "%s"
  source_group_owner_id    = "%s"
  description              = "tf-testAccSGRSrcOwnerCross"
}

# Exercise independent legacy account and owner ID passthrough in the egress direction.
resource "alicloud_security_group_rule" "egress" {
  type                       = "egress"
  ip_protocol                = "tcp"
  nic_type                   = "intranet"
  policy                     = "accept"
  port_range                 = "443/443"
  priority                   = 1
  security_group_id          = alicloud_security_group.test.id
  source_security_group_id   = "%s"
  source_group_owner_account = "%s"
  source_group_owner_id      = "%s"
  description                = "tf-testAccSGRSrcOwnerCrossEgress"
}
`, resourceName, resourceName, peerSecurityGroupID, peerAccountID, peerSecurityGroupID, peerAccountID, peerAccountID)
}
