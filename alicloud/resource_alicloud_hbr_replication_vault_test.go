package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_hbr_replication_vault",
		&resource.Sweeper{
			Name: "alicloud_hbr_replication_vault",
			F:    testSweepHbrReplicationVault,
		})
}

func testSweepHbrReplicationVault(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeVaults"
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}

	conn, err := client.NewHbrClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	for {
		var response map[string]interface{}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.Vaults.Vault", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Vaults.Vault", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["VaultName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping HBR Replication Vault : %s", item["VaultName"].(string))
				continue
			}

			action := "DeleteVault"
			request := map[string]interface{}{
				"VaultId": item["VaultId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete HBR Replication Vault (%s): %s", item["VaultName"].(string), err)
			}

			log.Printf("[INFO] Delete HBR Replication Vault success: %s ", item["VaultName"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudHBRReplicationVault_basic0(t *testing.T) {
	resourceId := "alicloud_hbr_replication_vault.default"
	checkoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBRReplicationVaultMap0)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrreplicationvault%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRReplicationVaultBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckHBRReplicationVaultDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name":                   name,
					"vault_storage_class":          "STANDARD",
					"replication_source_vault_id":  "${alicloud_hbr_vault.default.id}",
					"replication_source_region_id": "${var.region_source}",
					"provider":                     "alicloud.replication",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name":                   name,
						"vault_storage_class":          "STANDARD",
						"replication_source_vault_id":  CHECKSET,
						"replication_source_region_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name":  name,
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name":  name,
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudHBRReplicationVault_basic1(t *testing.T) {
	resourceId := "alicloud_hbr_replication_vault.default"
	checkoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBRReplicationVaultMap0)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrreplicationvault%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRReplicationVaultBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckHBRReplicationVaultDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name":                   name,
					"vault_storage_class":          "STANDARD",
					"replication_source_vault_id":  "${alicloud_hbr_vault.default.id}",
					"replication_source_region_id": "${var.region_source}",
					"description":                  name,
					"provider":                     "alicloud.replication",
					"depends_on":                   []string{"alicloud_hbr_vault.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name":                   name,
						"vault_storage_class":          "STANDARD",
						"replication_source_vault_id":  CHECKSET,
						"replication_source_region_id": CHECKSET,
						"description":                  name,
					}),
				),
			},
		},
	})
}

var AlicloudHBRReplicationVaultMap0 = map[string]string{}

func AlicloudHBRReplicationVaultBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region_source" {
  default = "%s"
}

provider "alicloud" {
  alias = "source"
  region = var.region_source
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  provider   = alicloud.source
}

data "alicloud_hbr_replication_vault_regions" "default" {}

locals {
	region_replication = data.alicloud_hbr_replication_vault_regions.default.regions.0.replication_region_id
}

provider "alicloud" {
  alias = "replication"
  region = local.region_replication
}
`, name, defaultRegionToTest)
}

func testAccCheckHBRReplicationVaultDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckHBRReplicationVaultDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckHBRReplicationVaultDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	hbrService := HbrService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_hbr_replication_vault" {
			continue
		}

		_, err := hbrService.DescribeHbrReplicationVault(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}
