package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
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
	for {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
			_, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, false)
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

func TestAccAliCloudHbrReplicationVault_basic0(t *testing.T) {
	resourceId := "alicloud_hbr_replication_vault.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacchbrrepvalt%d", rand)
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
		"alicloudshanghai": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckHBRReplicationVaultDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccHBRReplicationVaultConfig(name, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"rep"),
					resource.TestCheckResourceAttr(resourceId, "vault_storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceId, "replication_source_region_id", "cn-shanghai"),
					resource.TestCheckResourceAttrSet(resourceId, "replication_source_vault_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				Config: testAccHBRReplicationVaultConfig(name, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"rep"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
				),
			},
			{
				Config: testAccHBRReplicationVaultConfig(name, name+"update", name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"updaterep"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
				),
			},
			{
				Config: testAccHBRReplicationVaultConfig(name, name+"update", name+"update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"updaterep"),
					resource.TestCheckResourceAttr(resourceId, "description", name+"update"),
				),
			},
			{
				Config: testAccHBRReplicationVaultConfig(name, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"rep"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
				),
			},
		},
	})
}

func TestAccAliCloudHbrReplicationVault_basic1(t *testing.T) {
	resourceId := "alicloud_hbr_replication_vault.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacchbrrepvalt%d", rand)
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
		"alicloudshanghai": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckHBRReplicationVaultDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccHBRReplicationVaultConfig(name, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"rep"),
					resource.TestCheckResourceAttr(resourceId, "vault_storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceId, "replication_source_region_id", "cn-shanghai"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
					resource.TestCheckResourceAttrSet(resourceId, "replication_source_vault_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
		},
	})
}

// KMS-encrypted replication vault (source vault in cn-shanghai, replication in cn-hangzhou).
func TestAccAliCloudHbrReplicationVault_basic10880(t *testing.T) {
	resourceId := "alicloud_hbr_replication_vault.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacchbrrepvalt%d", rand)
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
		"alicloudshanghai": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p)
			return p, nil
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckHBRReplicationVaultDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccHBRReplicationVaultKMSConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vault_name", name+"rep"),
					resource.TestCheckResourceAttr(resourceId, "vault_storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceId, "replication_source_region_id", "cn-shanghai"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
					resource.TestCheckResourceAttr(resourceId, "encrypt_type", "KMS"),
					resource.TestCheckResourceAttr(resourceId, "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet(resourceId, "replication_source_vault_id"),
					resource.TestCheckResourceAttrSet(resourceId, "kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
		},
	})
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
	hbrService := HbrServiceV2{provider.Meta().(*connectivity.AliyunClient)}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_hbr_replication_vault" {
			continue
		}

		_, err := hbrService.DescribeHbrReplicationVault(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("HBR Replication Vault %s still exists", rs.Primary.ID))
	}

	return nil
}

// lintignore: R001
func TestUnitAlicloudHBRReplicationVault(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbr_replication_vault"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbr_replication_vault"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":                  "CreateReplicationVaultValue",
		"replication_source_region_id": "CreateReplicationVaultValue",
		"replication_source_vault_id":  "CreateReplicationVaultValue",
		"vault_name":                   "CreateReplicationVaultValue",
		"vault_storage_class":          "CreateReplicationVaultValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeVaults
		"Vaults": map[string]interface{}{
			"Vault": []interface{}{
				map[string]interface{}{
					"Description":               "CreateReplicationVaultValue",
					"ReplicationSourceRegionId": "CreateReplicationVaultValue",
					"ReplicationSourceVaultId":  "CreateReplicationVaultValue",
					"VaultName":                 "CreateReplicationVaultValue",
					"VaultStorageClass":         "CreateReplicationVaultValue",
					"VaultId":                   "CreateReplicationVaultValue",
					"VaultRegionId":             "CreateReplicationVaultValue",
					"Status":                    "CreateReplicationVaultValue",
				},
			},
		},
		"Success": "true",
		"VaultId": "CreateReplicationVaultValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateReplicationVault
		"VaultId": "CreateReplicationVaultValue",
		"Success": "true",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_hbr_replication_vault", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudHbrReplicationVaultCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeVaults Response
		"VaultId": "CreateReplicationVaultValue",
		"Success": "true",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateReplicationVault" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudHbrReplicationVaultCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_replication_vault"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudHbrReplicationVaultUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateVault
	attributesDiff := map[string]interface{}{
		"description": "UpdateVaultValue",
		"vault_name":  "UpdateVaultValue",
	}
	diff, err := newInstanceDiff("alicloud_hbr_replication_vault", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_replication_vault"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeVaults Response
		"Vaults": map[string]interface{}{
			"Vault": []interface{}{
				map[string]interface{}{
					"Description": "UpdateVaultValue",
					"VaultName":   "UpdateVaultValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateVault" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudHbrReplicationVaultUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_replication_vault"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeVaults" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudHbrReplicationVaultRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudHbrReplicationVaultDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVault" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudHbrReplicationVaultDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

func testAccHBRReplicationVaultConfig(name, vaultName, description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  alias  = "replication"
  region = "cn-hangzhou"
}

provider "alicloudshanghai" {
  alias  = "source"
  region = "cn-shanghai"
}

resource "alicloud_hbr_vault" "default" {
  provider   = alicloudshanghai.source
  vault_name = var.name
}

resource "alicloud_hbr_replication_vault" "default" {
  provider                     = alicloud.replication
  vault_storage_class          = "STANDARD"
  replication_source_vault_id  = alicloud_hbr_vault.default.id
  replication_source_region_id = "cn-shanghai"
  vault_name                   = "%srep"
  description                  = "%s"
}
`, name, vaultName, description)
}

func testAccHBRReplicationVaultKMSConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  alias  = "replication"
  region = "cn-hangzhou"
}

provider "alicloudshanghai" {
  alias  = "source"
  region = "cn-shanghai"
}

resource "alicloud_kms_key" "source" {
  provider               = alicloudshanghai.source
  description            = var.name
  pending_window_in_days = 7
  key_state              = "Enabled"
}

resource "alicloud_kms_key" "replication" {
  provider               = alicloud.replication
  description            = var.name
  pending_window_in_days = 7
  key_state              = "Enabled"
}

resource "alicloud_hbr_vault" "default" {
  provider     = alicloudshanghai.source
  vault_type   = "STANDARD"
  encrypt_type = "KMS"
  vault_name   = var.name
  kms_key_id   = alicloud_kms_key.source.id
}

resource "alicloud_hbr_replication_vault" "default" {
  provider                     = alicloud.replication
  vault_storage_class          = "STANDARD"
  replication_source_vault_id  = alicloud_hbr_vault.default.id
  replication_source_region_id = "cn-shanghai"
  vault_name                   = "%srep"
  description                  = var.name
  encrypt_type                 = "KMS"
  kms_key_id                   = alicloud_kms_key.replication.id
}
`, name, name)
}

