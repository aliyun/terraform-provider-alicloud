package ecd_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestMain wires the terraform-plugin-testing sweep machinery so the
// registered sweeper for alicloud_ecd_ad_connector_directory runs under
// `go test -sweep ./alicloud/service/ecd/`. The flat root package keeps its
// own SDK TestMain for the resources that still live there.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers(
		"alicloud_ecd_ad_connector_directory",
		&resource.Sweeper{
			Name: "alicloud_ecd_ad_connector_directory",
			F:    testSweepEcdAdConnectorDirectory,
		})
}

// testSweepEcdClient builds an AliyunClient from the environment for a given
// region, mirroring the flat suite's sharedClientForRegion so the sweeper and
// CheckDestroy can call the package's Service methods from outside the root
// package.
func testSweepEcdClient(region string) (*connectivity.AliyunClient, error) {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	endpoints := &sync.Map{}
	signVersion := &sync.Map{}
	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   endpoints,
		SignVersion: signVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	if accountId := os.Getenv("ALICLOUD_ACCOUNT_ID"); accountId == "" {
		if v := os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID"); v != "" {
			conf.AccountId = v
		}
	} else {
		conf.AccountId = accountId
	}
	if accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE"); accountType == "" {
		if v := os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE"); v != "" {
			conf.AccountType = v
		}
	} else {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential
	return conf.Client()
}

func testSweepEcdAdConnectorDirectory(region string) error {
	if !ecdRegionSupported(region) {
		log.Printf("[INFO] Skipping Ecd Ad Connector Directory unsupported region: %s", region)
		return nil
	}

	rawClient, err := testSweepEcdClient(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDirectories"
	request := map[string]interface{}{
		"RegionId":      rawClient.RegionId,
		"MaxResults":    helper.PageSizeLarge,
		"DirectoryType": "AD_CONNECTOR",
	}
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = rawClient.RpcPost("ecd", "2020-09-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Directories", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Directories", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecd Ad Connector Directory: %s", item["Name"].(string))
				continue
			}
			deleteAction := "DeleteDirectories"
			deleteRequest := map[string]interface{}{
				"DirectoryId": []string{fmt.Sprint(item["DirectoryId"])},
				"RegionId":    rawClient.RegionId,
			}
			_, err = rawClient.RpcPost("ecd", "2020-09-30", deleteAction, nil, deleteRequest, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecd Ad Connector Directory (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Ecd Ad Connector Directory success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// ecdRegionSupported reports whether the region is one ECD lists as supported.
// ECD is available in a fixed set of regions; the sweeper skips the others so it
// never lists resources in a region where the ECD control plane is absent.
func ecdRegionSupported(region string) bool {
	for _, r := range connectivity.EcdSupportRegions {
		if region == string(r) {
			return true
		}
	}
	return false
}

// testAccEcdAdConnectorDirectoryPreCheck pins the test region to cn-hangzhou.
// acctest.PreCheck defaults ALICLOUD_REGION to cn-beijing when the runner has not
// set it; cn-beijing is in EcdSupportRegions yet its ECD control plane has no
// cluster for the AD connector's VPC CIDR, so CreateADConnectorDirectory fails
// with chooseClusterError ("Cannot choose cluster, please try other CIDR block").
// The flat suite's checkoutSupportedRegions backs off to Hangzhou, and the
// original test's hardcoded cn-hangzhou-j zone targeted cn-hangzhou, so pin here.
func testAccEcdAdConnectorDirectoryPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func TestAccAliCloudECDAdConnectorDirectory_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testaccecdadconnectordirectory%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEcdAdConnectorDirectoryPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcdAdConnectorDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcdAdConnectorDirectoryConfig_basic0(name, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcdAdConnectorDirectoryExists("alicloud_ecd_ad_connector_directory.default"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "directory_name", name),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "domain_user_name", "sAMAccountName"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "vswitch_ids.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "dns_address.#", "1"),
				),
			},
			{
				ResourceName:            "alicloud_ecd_ad_connector_directory.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification", "domain_password"},
			},
		},
	})
}

func TestAccAliCloudECDAdConnectorDirectory_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testaccecdadconnectordirectory%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEcdAdConnectorDirectoryPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcdAdConnectorDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcdAdConnectorDirectoryConfig_basic1(name, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcdAdConnectorDirectoryExists("alicloud_ecd_ad_connector_directory.default"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "directory_name", name),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "desktop_access_type", "INTERNET"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "sub_domain_name", "child.example.com"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "enable_admin_access", "true"),
					resource.TestCheckResourceAttr("alicloud_ecd_ad_connector_directory.default", "mfa_enabled", "false"),
				),
			},
			{
				ResourceName:            "alicloud_ecd_ad_connector_directory.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification", "domain_password"},
			},
		},
	})
}

// testAccCheckEcdAdConnectorDirectoryExists verifies the directory was created
// and settled into REGISTERED (the lifecycle wait the resource now enforces).
func testAccCheckEcdAdConnectorDirectoryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		region := os.Getenv("ALICLOUD_REGION")
		client, err := testSweepEcdClient(region)
		if err != nil {
			return fmt.Errorf("error getting AliCloud client: %s", err)
		}
		ecdService := ecd.NewService(client)
		object, err := ecdService.DescribeEcdAdConnectorDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Resource alicloud_ecd_ad_connector_directory DescribeEcdAdConnectorDirectory Failed!!! %s", err)
		}
		if status, ok := object["Status"]; !ok || fmt.Sprint(status) != "REGISTERED" {
			return fmt.Errorf("expected directory status REGISTERED, got %v", status)
		}
		return nil
	}
}

// testAccCheckEcdAdConnectorDirectoryDestroy verifies the directory is gone
// after destroy, the coverage the flat suite's checkResourceDestroy provided.
func testAccCheckEcdAdConnectorDirectoryDestroy(s *terraform.State) error {
	region := os.Getenv("ALICLOUD_REGION")
	client, err := testSweepEcdClient(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	ecdService := ecd.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecd_ad_connector_directory" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		_, err := ecdService.DescribeEcdAdConnectorDirectory(rs.Primary.ID)
		if err != nil {
			if ecdServiceNotFound(err) {
				continue
			}
			return fmt.Errorf("Error checking directory destroyed: %s", err)
		}
		return fmt.Errorf("ECD AD Connector Directory %s still exists", rs.Primary.ID)
	}
	return nil
}

// ecdServiceNotFound reports whether the error surfaced by the ecd Service is a
// not-found condition (errs.NotFoundError is package-private; the wrapped error
// string carries the ResourceNotFound marker used across the suite).
func ecdServiceNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ResourceNotfound")
}

func testAccEcdAdConnectorDirectoryConfig_basic0(name string, rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ecd_ad_connector_directory" "default" {
  directory_name  = var.name
  domain_name     = "corp.%[2]d.example.com"
  vswitch_ids     = [alicloud_vswitch.default.id]
  dns_address     = ["127.0.0.2"]
  domain_password = "YourPassword1234"
  domain_user_name = "sAMAccountName"
}
`, name, rand)
}

func testAccEcdAdConnectorDirectoryConfig_basic1(name string, rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ecd_ad_connector_directory" "default" {
  directory_name         = var.name
  domain_name            = "corp.%[2]d.example.com"
  vswitch_ids            = [alicloud_vswitch.default.id]
  dns_address            = ["127.0.0.2"]
  desktop_access_type    = "INTERNET"
  domain_password        = "YourPassword1234"
  domain_user_name       = "sAMAccountName"
  specification          = "1"
  sub_domain_dns_address = ["127.0.0.3"]
  sub_domain_name        = "child.example.com"
  enable_admin_access    = true
  mfa_enabled            = false
}
`, name, rand)
}
