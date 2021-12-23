package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
	conn, err := client.NewCloudssoClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso Directory (%s): %s", item["DirectoryName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud Sso Directory success: %s ", item["DirectoryName"].(string))
	}
	return nil
}

func TestAccAlicloudCloudSSODirectory_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_directory.default"
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
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
			testAccPreCheckEnterpriseAccountEnabled(t)
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
					"saml_identity_provider_configuration": []map[string]interface{}{
						{
							"sso_status": "Disabled",
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
