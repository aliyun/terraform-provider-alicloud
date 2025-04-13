package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_cloud_sso_directory",
		&resource.Sweeper{
			Name: "alicloud_cloud_sso_directory",
			F:    testSweepCloudSsoDirectory,
		})
}

func testSweepCloudSsoDirectory(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListDirectories"
	request := map[string]interface{}{}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.Directories", response)
	if formatInt(response["TotalCounts"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Directories", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["DirectoryName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Cloud Sso Directory: %s", item["DirectoryName"].(string))
			continue
		}
		action := "DeleteDirectory"
		request := map[string]interface{}{
			"DirectoryId": item["DirectoryId"],
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso Directory (%s): %s", item["DirectoryName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud Sso Directory success: %s ", item["DirectoryName"].(string))
	}
	return nil
}

func TestAccAliCloudCloudSSODirectory_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_directory.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudSSODirectoryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoDirectory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudssodirectory%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSODirectoryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_name": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_authentication_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_authentication_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_authentication_status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_authentication_status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scim_synchronization_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scim_synchronization_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scim_synchronization_status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scim_synchronization_status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status":                "Enabled",
							"encoded_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG1kOkVudGl0eURlc2NyaXB0b3IgZW50aXR5SUQ9Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9hc3NlcnRpb24vdGVzdCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlES3pDQ0FoT2dBd0lCQWdJQkFUQU5CZ2txaGtpRzl3MEJBUXNGQURCWU1Sa3dGd1lEVlFRREV4QmhiR2xpWVdKaFkyeHZkV1F1ClkyOXRNUnd3R2dZRFZRUUxFeE5KWkdWdWRHbDBlU0JOWVc1aFoyVnRaVzUwTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVFzd0NRWUQKVlFRR0V3SkRUakFnRncweU1UQTFNRFl3T0RNMk1EQmFHQTh5TVRJeE1EUXhNakE0TXpZd01Gb3dXREVaTUJjR0ExVUVBeE1RWVd4cApZbUZpWVdOc2IzVmtMbU52YlRFY01Cb0dBMVVFQ3hNVFNXUmxiblJwZEhrZ1RXRnVZV2RsYldWdWRERVFNQTRHQTFVRUNoTUhRV3hwClltRmlZVEVMTUFrR0ExVUVCaE1DUTA0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURabHVOUlc2N1gKQlhVdFJMaG9BcTJ5d3VGRDd5WVpsOEliTUNJZ0NBdnJzekxhRmVQM0VvRWg4RjJ0M2FBWlJyZHJGMDdSbmNjTTZQRC9Ud2d4dzNoVgpDSmh2QnloNjh4UXVENXRja2pLei90Wnl3MVB5NXJPb05Va001N09ZL1AvZWc2MHg1SlJxaXdmZXZFcnJPK0xKWnVxMTNqTmpzRWFNCnRsRnhuTUkrY3FvUW9QbnZSUnVWNFVQdXJ0TEpJemNwQ2QvRnJxQjVvZTg5TjgySFpCcVE4eTdscVNsdHNuOTF0S05xRm44U3lnN3QKWitJQzRweWVsVDJVYit4RkFXRStjTWd4RVloR0R5UzVWSGNaWkc5OC90d0RWUHlRSjV5NWhJM0U2MWE3OW1mTkxHa0IvTGVXODF3MQpTZFM4SldQL3hsNWt6dDR0aXNkd0kyejJmMm8zQWdNQkFBRXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBTDFseUZIcnRET1R6bUZsCkpScmFsNlVLekRSSlU4djVOcGNPeTdlbVZBR2RMamNqRll5ZzR3cG1FZ3pOUWxscFZZd0VBRi9LbG5nOGxNbC8zekI2cWJEVFNSWnEKbFhpdDU4eHJsM0RiWnJ1OUdVOVMrdW9XYngrdVkyWStVNkpOOWhJSFRkcHc5SmFkTWY3akJzL2pFWWFuRDNObUU0SjltSThMeDNXYwpHL0lFLzNOOTY1NXM1UHJPTkc0MTgySXJFTG1TdzRrRkxDS0gzdi9IUnN0Ym9PSkttTWRuTHc4VGNaYWFNZHp2N1VEK1NoRlRPbFdmCmJrRmhISVZVSGNVSGN3S1kxSitjSmF0VVdPakNmWGZkNVl4Uy9aM3dhbURtbmRxTXJaT3kxUEh2dkhxNHFGQ3dsQTY1V2pzTDF6dkYKcE51MjQ2WFIxSHNRcVo4UTlvVk1kNU09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpTaW5nbGVMb2dvdXRTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2N5LnRlc3QuY29tL3NhbWwvbG9nb3V0L3Rlc3QiLz48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9sb2dvdXQvdGVzdCIvPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6cGVyc2lzdGVudDwvbWQ6TmFtZUlERm9ybWF0PjxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2N5LnRlc3QuY29tL3NhbWwvYXNzZXJ0aW9uL3Rlc3QiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9hc3NlcnRpb24vdGVzdCIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_identity_provider_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_name":              "${var.name}",
					"mfa_authentication_status":   "Enabled",
					"scim_synchronization_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_name":              name,
						"mfa_authentication_status":   "Enabled",
						"scim_synchronization_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status":                "Disabled",
							"encoded_metadata_document": "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_identity_provider_configuration.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCloudSSODirectoryMap0 = map[string]string{}

func AlicloudCloudSSODirectoryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// Test CloudSSO Directory. >>> Resource test cases, automatically generated.
// Case cloudsso目录 9760
func TestAccAliCloudCloudSSODirectory_basic9760(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_directory.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudSSODirectoryMap9760)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSODirectory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSODirectoryBasicDependence9760)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scim_synchronization_status": "Enabled",
					"directory_name":              name,
					"mfa_authentication_status":   "Disabled",
					"password_policy": []map[string]interface{}{
						{
							"min_password_length":           "8",
							"min_password_different_chars":  "8",
							"max_password_age":              "90",
							"password_reuse_prevention":     "1",
							"max_login_attempts":            "5",
							"password_not_contain_username": "true",
						},
					},
					"mfa_authentication_setting_info": []map[string]interface{}{
						{
							"mfa_authentication_advance_settings": "OnlyRiskyLogin",
							"operation_for_risk_login":            "EnforceVerify",
						},
					},
					"directory_global_access_status": "Disabled",
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status":                "Enabled",
							"encoded_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG1kOkVudGl0eURlc2NyaXB0b3IgZW50aXR5SUQ9Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9hc3NlcnRpb24vdGVzdCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlES3pDQ0FoT2dBd0lCQWdJQkFUQU5CZ2txaGtpRzl3MEJBUXNGQURCWU1Sa3dGd1lEVlFRREV4QmhiR2xpWVdKaFkyeHZkV1F1ClkyOXRNUnd3R2dZRFZRUUxFeE5KWkdWdWRHbDBlU0JOWVc1aFoyVnRaVzUwTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVFzd0NRWUQKVlFRR0V3SkRUakFnRncweU1UQTFNRFl3T0RNMk1EQmFHQTh5TVRJeE1EUXhNakE0TXpZd01Gb3dXREVaTUJjR0ExVUVBeE1RWVd4cApZbUZpWVdOc2IzVmtMbU52YlRFY01Cb0dBMVVFQ3hNVFNXUmxiblJwZEhrZ1RXRnVZV2RsYldWdWRERVFNQTRHQTFVRUNoTUhRV3hwClltRmlZVEVMTUFrR0ExVUVCaE1DUTA0d2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURabHVOUlc2N1gKQlhVdFJMaG9BcTJ5d3VGRDd5WVpsOEliTUNJZ0NBdnJzekxhRmVQM0VvRWg4RjJ0M2FBWlJyZHJGMDdSbmNjTTZQRC9Ud2d4dzNoVgpDSmh2QnloNjh4UXVENXRja2pLei90Wnl3MVB5NXJPb05Va001N09ZL1AvZWc2MHg1SlJxaXdmZXZFcnJPK0xKWnVxMTNqTmpzRWFNCnRsRnhuTUkrY3FvUW9QbnZSUnVWNFVQdXJ0TEpJemNwQ2QvRnJxQjVvZTg5TjgySFpCcVE4eTdscVNsdHNuOTF0S05xRm44U3lnN3QKWitJQzRweWVsVDJVYit4RkFXRStjTWd4RVloR0R5UzVWSGNaWkc5OC90d0RWUHlRSjV5NWhJM0U2MWE3OW1mTkxHa0IvTGVXODF3MQpTZFM4SldQL3hsNWt6dDR0aXNkd0kyejJmMm8zQWdNQkFBRXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBTDFseUZIcnRET1R6bUZsCkpScmFsNlVLekRSSlU4djVOcGNPeTdlbVZBR2RMamNqRll5ZzR3cG1FZ3pOUWxscFZZd0VBRi9LbG5nOGxNbC8zekI2cWJEVFNSWnEKbFhpdDU4eHJsM0RiWnJ1OUdVOVMrdW9XYngrdVkyWStVNkpOOWhJSFRkcHc5SmFkTWY3akJzL2pFWWFuRDNObUU0SjltSThMeDNXYwpHL0lFLzNOOTY1NXM1UHJPTkc0MTgySXJFTG1TdzRrRkxDS0gzdi9IUnN0Ym9PSkttTWRuTHc4VGNaYWFNZHp2N1VEK1NoRlRPbFdmCmJrRmhISVZVSGNVSGN3S1kxSitjSmF0VVdPakNmWGZkNVl4Uy9aM3dhbURtbmRxTXJaT3kxUEh2dkhxNHFGQ3dsQTY1V2pzTDF6dkYKcE51MjQ2WFIxSHNRcVo4UTlvVk1kNU09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpTaW5nbGVMb2dvdXRTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2N5LnRlc3QuY29tL3NhbWwvbG9nb3V0L3Rlc3QiLz48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9sb2dvdXQvdGVzdCIvPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6cGVyc2lzdGVudDwvbWQ6TmFtZUlERm9ybWF0PjxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2N5LnRlc3QuY29tL3NhbWwvYXNzZXJ0aW9uL3Rlc3QiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vY3kudGVzdC5jb20vc2FtbC9hc3NlcnRpb24vdGVzdCIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+",
							"binding_type":              "Post",
						},
					},
					"login_preference": []map[string]interface{}{
						{
							"allow_user_to_get_credentials": "true",
							"login_network_masks":           "192.163.2.0/24",
						},
					},
					"user_provisioning_configuration": []map[string]interface{}{
						{
							"default_landing_page": "https://home.console.aliyun.com",
							"session_duration":     "12",
						},
					},
					"saml_service_provider": []map[string]interface{}{
						{
							"authn_sign_algo":             "rsa-sha256",
							"certificate_type":            "self-signed",
							"support_encrypted_assertion": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scim_synchronization_status":    "Enabled",
						"directory_name":                 name,
						"mfa_authentication_status":      "Disabled",
						"directory_global_access_status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scim_synchronization_status": "Disabled",
					"mfa_authentication_status":   "Enabled",
					"password_policy": []map[string]interface{}{
						{
							"min_password_length":           "9",
							"min_password_different_chars":  "9",
							"password_not_contain_username": "false",
							"max_password_age":              "80",
							"password_reuse_prevention":     "2",
							"max_login_attempts":            "6",
						},
					},
					"mfa_authentication_setting_info": []map[string]interface{}{
						{
							"mfa_authentication_advance_settings": "Enabled",
							"operation_for_risk_login":            "Autonomous",
						},
					},
					"directory_global_access_status": "Enabled",
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status":                "Enabled",
							"encoded_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cHM6Ly9wb3J0YWwuc3NvLmFwLXNvdXRoZWFzdC0xLmFtYXpvbmF3cy5jb20vc2FtbC9hc3NlcnRpb24vTmpjeE1qRTNPRE0wTlRZeFgybHVjeTAxWVRRM1l6STFNMlpsWXpjNU1qZ3oiIHhtbG5zOm1kPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bWV0YWRhdGEiPjxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+PG1kOktleURlc2NyaXB0b3IgdXNlPSJzaWduaW5nIj48ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+PGRzOlg1MDlEYXRhPjxkczpYNTA5Q2VydGlmaWNhdGU+TUlJREF6Q0NBZXVnQXdJQkFnSUJBVEFOQmdrcWhraUc5dzBCQVFzRkFEQkZNUll3RkFZRFZRUUREQTFoYldGNmIyNWhkM011WTI5dE1RMHdDd1lEVlFRTERBUkpSRUZUTVE4d0RRWURWUVFLREFaQmJXRjZiMjR4Q3pBSkJnTlZCQVlUQWxWVE1CNFhEVEl3TURNd09UQXlOREkxTjFvWERUSTFNRE13T1RBeU5ESTFOMW93UlRFV01CUUdBMVVFQXd3TllXMWhlbTl1WVhkekxtTnZiVEVOTUFzR0ExVUVDd3dFU1VSQlV6RVBNQTBHQTFVRUNnd0dRVzFoZW05dU1Rc3dDUVlEVlFRR0V3SlZVekNDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFKbXVmZE9PTWt5Z2dZQ05mN0twWjhwWnpDY1pCNmhFMy9JbW5zOEJzd2ZYTVZRNWVDVWkxaS8ydzVQclVTaVJzMzdRTU04QURuZHVDQ21tY1FhMHhzZ0tFQVZSU2pDV3JTTXEvQWNVZ0Jsamt0dk5JN2dmaVpCNTRoSnV4RzY5T2o5dlp3S1FoaTVYdGZDcDZFbXVWVnNiRjFXdy9XMDl0Wm54VkJZV1dOUS9vS0xBQ3BjTGZmWnFZVndmN0RLazZNUkV0eE5wUThiN21TSVdBcEMzS2d2RjBZUldpeVBkYkZqT2ZKWjRuZXJtSk1NN0F4cWc5TXdmWngyN3MrbnkyTERzODlKdER0VGpvNng2NzZaZXZNYXVDMVVYZDNneW0vVnFOVDNUb1NHYzMwTW9NTmJheVkxZ0FnUVRIWHFvZndCUVVjR1EvQXhEeHFUdmtkdWVRNDBDQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFSOW1FcG5FWXV5MGFwa0prM2Y1Q1g2VnRxSVlVL3FEUUtBY1UyMzZtSzFyTmpuRmJYM2Yvek5pRUpkOThzL01GcUE5VDJoOVVsaVZkaGNHamRyWVhwQ3VlTGV4eVJyNnp4UGJFWld3ZXcyOTVCYVdKemVtS0ZlTU5rVWtUSk1GSW9yRDU4RGk2QzdrYlhZdDdUa1pidWJMdUV4KzMzNGFNTkFyU1VNRnRYc0tBVk5tdkt2dW5TOHhSa0hZN01hNmdBMVhyOUFKNytFaDdzZE5talh4aDVSOUdoVlJhUmdhMnlraFRUaWo2R1B2WEt6NFFhenVUTHNyazh0QlBNcXhuTWZEdnppQ1pEeUxySjNadVVMNEFkdlE2bVBRWHFvZ1R6ajhTOVhrajNuQWcvL0w0bnlrK2xqaEZnU0RuQW5TOHF2bnR6UnRMaGxxSDdtTG1ibENjSFE9PTwvZHM6WDUwOUNlcnRpZmljYXRlPjwvZHM6WDUwOURhdGE+PC9kczpLZXlJbmZvPjwvbWQ6S2V5RGVzY3JpcHRvcj48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1QT1NUIiBMb2NhdGlvbj0iaHR0cHM6Ly9wb3J0YWwuc3NvLmFwLXNvdXRoZWFzdC0xLmFtYXpvbmF3cy5jb20vc2FtbC9sb2dvdXQvTmpjeE1qRTNPRE0wTlRZeFgybHVjeTAxWVRRM1l6STFNMlpsWXpjNU1qZ3oiLz48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vcG9ydGFsLnNzby5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3NhbWwvbG9nb3V0L05qY3hNakUzT0RNME5UWXhYMmx1Y3kwMVlUUTNZekkxTTJabFl6YzVNamd6Ii8+PG1kOk5hbWVJREZvcm1hdC8+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHBzOi8vcG9ydGFsLnNzby5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3NhbWwvYXNzZXJ0aW9uL05qY3hNakUzT0RNME5UWXhYMmx1Y3kwMVlUUTNZekkxTTJabFl6YzVNamd6Ii8+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUmVkaXJlY3QiIExvY2F0aW9uPSJodHRwczovL3BvcnRhbC5zc28uYXAtc291dGhlYXN0LTEuYW1hem9uYXdzLmNvbS9zYW1sL2Fzc2VydGlvbi9OamN4TWpFM09ETTBOVFl4WDJsdWN5MDFZVFEzWXpJMU0yWmxZemM1TWpneiIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+",
							"binding_type":              "Post",
						},
					},
					"login_preference": []map[string]interface{}{
						{
							"login_network_masks":           "127.0.0.1/24",
							"allow_user_to_get_credentials": "false",
						},
					},
					"user_provisioning_configuration": []map[string]interface{}{
						{
							"default_landing_page": "https://www.aliyun.com",
							"session_duration":     "8",
						},
					},
					"saml_service_provider": []map[string]interface{}{
						{
							"authn_sign_algo":             "rsa-sha256",
							"certificate_type":            "self-signed",
							"support_encrypted_assertion": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scim_synchronization_status":    "Disabled",
						"mfa_authentication_status":      "Enabled",
						"directory_global_access_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_authentication_status": "Disabled",
					"password_policy": []map[string]interface{}{
						{
							"min_password_length":          "8",
							"min_password_different_chars": "8",
							"max_password_age":             "90",
							"password_reuse_prevention":    "1",
							"max_login_attempts":           "5",
						},
					},
					"mfa_authentication_setting_info": []map[string]interface{}{
						{
							"mfa_authentication_advance_settings": "OnlyRiskyLogin",
							"operation_for_risk_login":            "EnforceVerify",
						},
					},
					"directory_global_access_status": "Disabled",
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status":                "Enabled",
							"entity_id":                 "https://portal.sso.ap-southeast-1.amazonaws.com/saml/assertion/NjcxMjE3ODM0NTYxX2lucy01YTQ3YzI1M2ZlYzc5Mjgz",
							"login_url":                 "https://portal.sso.ap-southeast-1.amazonaws.com/saml/assertion/NjcxMjE3ODM0NTYxX2lucy01YTQ3YzI1M2ZlYzc5Mjgz",
							"want_request_signed":       "false",
							"binding_type":              "Post",
							"encoded_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cHM6Ly9wb3J0YWwuc3NvLmFwLXNvdXRoZWFzdC0xLmFtYXpvbmF3cy5jb20vc2FtbC9hc3NlcnRpb24vTmpjeE1qRTNPRE0wTlRZeFgybHVjeTAxWVRRM1l6STFNMlpsWXpjNU1qZ3oiIHhtbG5zOm1kPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bWV0YWRhdGEiPjxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+PG1kOktleURlc2NyaXB0b3IgdXNlPSJzaWduaW5nIj48ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+PGRzOlg1MDlEYXRhPjxkczpYNTA5Q2VydGlmaWNhdGU+TUlJREF6Q0NBZXVnQXdJQkFnSUJBVEFOQmdrcWhraUc5dzBCQVFzRkFEQkZNUll3RkFZRFZRUUREQTFoYldGNmIyNWhkM011WTI5dE1RMHdDd1lEVlFRTERBUkpSRUZUTVE4d0RRWURWUVFLREFaQmJXRjZiMjR4Q3pBSkJnTlZCQVlUQWxWVE1CNFhEVEl3TURNd09UQXlOREkxTjFvWERUSTFNRE13T1RBeU5ESTFOMW93UlRFV01CUUdBMVVFQXd3TllXMWhlbTl1WVhkekxtTnZiVEVOTUFzR0ExVUVDd3dFU1VSQlV6RVBNQTBHQTFVRUNnd0dRVzFoZW05dU1Rc3dDUVlEVlFRR0V3SlZVekNDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFKbXVmZE9PTWt5Z2dZQ05mN0twWjhwWnpDY1pCNmhFMy9JbW5zOEJzd2ZYTVZRNWVDVWkxaS8ydzVQclVTaVJzMzdRTU04QURuZHVDQ21tY1FhMHhzZ0tFQVZSU2pDV3JTTXEvQWNVZ0Jsamt0dk5JN2dmaVpCNTRoSnV4RzY5T2o5dlp3S1FoaTVYdGZDcDZFbXVWVnNiRjFXdy9XMDl0Wm54VkJZV1dOUS9vS0xBQ3BjTGZmWnFZVndmN0RLazZNUkV0eE5wUThiN21TSVdBcEMzS2d2RjBZUldpeVBkYkZqT2ZKWjRuZXJtSk1NN0F4cWc5TXdmWngyN3MrbnkyTERzODlKdER0VGpvNng2NzZaZXZNYXVDMVVYZDNneW0vVnFOVDNUb1NHYzMwTW9NTmJheVkxZ0FnUVRIWHFvZndCUVVjR1EvQXhEeHFUdmtkdWVRNDBDQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFSOW1FcG5FWXV5MGFwa0prM2Y1Q1g2VnRxSVlVL3FEUUtBY1UyMzZtSzFyTmpuRmJYM2Yvek5pRUpkOThzL01GcUE5VDJoOVVsaVZkaGNHamRyWVhwQ3VlTGV4eVJyNnp4UGJFWld3ZXcyOTVCYVdKemVtS0ZlTU5rVWtUSk1GSW9yRDU4RGk2QzdrYlhZdDdUa1pidWJMdUV4KzMzNGFNTkFyU1VNRnRYc0tBVk5tdkt2dW5TOHhSa0hZN01hNmdBMVhyOUFKNytFaDdzZE5talh4aDVSOUdoVlJhUmdhMnlraFRUaWo2R1B2WEt6NFFhenVUTHNyazh0QlBNcXhuTWZEdnppQ1pEeUxySjNadVVMNEFkdlE2bVBRWHFvZ1R6ajhTOVhrajNuQWcvL0w0bnlrK2xqaEZnU0RuQW5TOHF2bnR6UnRMaGxxSDdtTG1ibENjSFE9PTwvZHM6WDUwOUNlcnRpZmljYXRlPjwvZHM6WDUwOURhdGE+PC9kczpLZXlJbmZvPjwvbWQ6S2V5RGVzY3JpcHRvcj48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1QT1NUIiBMb2NhdGlvbj0iaHR0cHM6Ly9wb3J0YWwuc3NvLmFwLXNvdXRoZWFzdC0xLmFtYXpvbmF3cy5jb20vc2FtbC9sb2dvdXQvTmpjeE1qRTNPRE0wTlRZeFgybHVjeTAxWVRRM1l6STFNMlpsWXpjNU1qZ3oiLz48bWQ6U2luZ2xlTG9nb3V0U2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vcG9ydGFsLnNzby5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3NhbWwvbG9nb3V0L05qY3hNakUzT0RNME5UWXhYMmx1Y3kwMVlUUTNZekkxTTJabFl6YzVNamd6Ii8+PG1kOk5hbWVJREZvcm1hdC8+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHBzOi8vcG9ydGFsLnNzby5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tL3NhbWwvYXNzZXJ0aW9uL05qY3hNakUzT0RNME5UWXhYMmx1Y3kwMVlUUTNZekkxTTJabFl6YzVNamd6Ii8+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUmVkaXJlY3QiIExvY2F0aW9uPSJodHRwczovL3BvcnRhbC5zc28uYXAtc291dGhlYXN0LTEuYW1hem9uYXdzLmNvbS9zYW1sL2Fzc2VydGlvbi9OamN4TWpFM09ETTBOVFl4WDJsdWN5MDFZVFEzWXpJMU0yWmxZemM1TWpneiIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+",
						},
					},
					"login_preference": []map[string]interface{}{
						{
							"login_network_masks":           "255.255.255.255/32",
							"allow_user_to_get_credentials": "false",
						},
					},
					"user_provisioning_configuration": []map[string]interface{}{
						{
							"default_landing_page": "https://www.alibabacloud.com",
							"session_duration":     "3",
						},
					},
					"saml_service_provider": []map[string]interface{}{
						{
							"authn_sign_algo":             "rsa-sha1",
							"certificate_type":            "public",
							"support_encrypted_assertion": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_authentication_status":      "Disabled",
						"directory_global_access_status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_policy": []map[string]interface{}{
						{
							"min_password_length":          "8",
							"min_password_different_chars": "8",
							"max_password_age":             "90",
							"password_reuse_prevention":    "1",
							"max_login_attempts":           "5",
						},
					},
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status": "Disabled",
						},
					},
					"login_preference": []map[string]interface{}{
						{
							"login_network_masks": "255.255.255.254/32",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_policy": []map[string]interface{}{
						{
							"min_password_length":          "8",
							"min_password_different_chars": "8",
							"max_password_age":             "90",
							"password_reuse_prevention":    "1",
							"max_login_attempts":           "5",
						},
					},
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"encoded_metadata_document": "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"directory_name", "saml_identity_provider_configuration"},
			},
		},
	})
}

var AlicloudCloudSSODirectoryMap9760 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudCloudSSODirectoryBasicDependence9760(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "directoryname2" {
  default = "apispecdnameprod02"
}

variable "directoryname1" {
  default = "apispecdnameprod01"
}


`, name)
}

// Test CloudSSO Directory. <<< Resource test cases, automatically generated.
