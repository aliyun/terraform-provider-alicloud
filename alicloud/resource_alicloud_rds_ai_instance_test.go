// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RdsAi Instance. >>> Resource test cases, automatically generated.
// Case AI实例测试用例 12309
func TestAccAliCloudRdsAiInstance_basic12309(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_ai_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsAiInstanceMap12309)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsAiServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAiInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrdsai%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsAiInstanceBasicDependence12309)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":         name,
					"app_type":         "supabase",
					"db_instance_name": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":         name,
						"app_type":         "supabase",
						"db_instance_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_config_list": []map[string]interface{}{
						{
							"name":  "API_EXTERNAL_URL",
							"value": "http://localhost:3001",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_config_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_key":  "-----BEGIN PRIVATE KEY-----\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKR3GSjlFvPyka\\ntE96Zj15bsU1n+TtIpgGwVQvdO3nK8S0XpXhP+9JQlVccTn9oDpJtKoWMRi4DuHO\\nFdqzSNrvCXqcPUaH+7qI1dVoYJQRyI3JdofjAdnKpA4ngZTNlkAx3yI2QjT4vay6\\nbuOuB+KIOvPHJiB/s/JZkhLnEm7mnUrZgrtCqK56qP888UVAFpVEPsanZU2IEHRR\\nG7ETRwbcvDiImjIkd15g3b6w1xHHiSru3a0PIL+BthAhGKkD384zZWmlZCIQb9NH\\nMVeZPIG1SIIPSQfDakLsY8KneXsQnCYvwqvO6bdTqL5NUHlPyNVk+p48tltFeece\\n7aKioOjzAgMBAAECggEAJtGZIRgA1smXONG7ovC2AXTZkdXyl/OYm0tEvarB8Sg8\\nIqU4PDcJ09RQD2KHT34NUZHDRmj7pm7stKELDHcB1PfLuOole/k6LgJjZxmJsPP9\\nCdmecFktk67yLHC4vs+D2E5LAYCpK8cyu8CGHyLSPXSazfAqMne1Ha1jxUaLU+Qq\\nPsmGX2O808R9DTrdcd7fzT9MwBiAbO3zX0fDQeAXF/NFeuty2YU3hMPxrgm5xirL\\nZnbShZJ7uZPXNslOSK8qJTcGjYbBAPydAKWTS8b1fww78dqvUIEUPiqi9CfbH8T5\\nvwq+pROIbBORhMBxgFfeNZg0PW3r0kRNYi7SGNNzcQKBgQDwBnXQUt9AG6rCaF4c\\njiwNr3l4r0wA3ltrNSsxvgw8JJ4+mn0J6NfYpVFfXxteXInGAW1kWOT2U2rUIRLm\\nk4nt48OVBiYMKxk1aqqWRjP+E51Kv3ZJjUGiBuzFLaAlELgr7TeH9M1s9EiujSgx\\nmO44vNtZxnvhx2799wdM5XkT2QKBgQDXvd2xFzDuEM37Al/aAmmCJNvl8vGojiOt\\nWIhKilke8h7NVEum3A6OV80htMp2Nrxa3fbuGv0g2NA3H2Ddh/ALqJ56DZlYRwxc\\n/GkJ4aHr3gvMgmRBNcOSrw8zEutntzj3y/KRdxofbqloTpmHLy7T/dO25ycNGaax\\nFuEOsCh/qwKBgH+C7xO88t0b1Ztx1o1U+hJLJjz242mStv49QLUsQVOyIF8hs0uQ\\nZxqwuInx/JgkkQyftX2ZvAkgR6Bp8aCMwLmgRkbk/VF5k+rMv9MVeImB4g3TqQNq\\nB3QMObyGgI0wVKcBXn7bjkZTgEk6tB+lHukFa4JF74oCaPSCR16SicABAoGAJ9yX\\n8pmTW9lVBbTpmvFpJzfCPZmG1xr0MpoyDHvFfbdEw7F4jOsJ8Xj9mOW7pt30LWHn\\ntxTbgk8tIZI1ppjwXGcaaPjMRYhTG3czvTSoY1lSmsXY2kehzB92UnyDbFVpPDe9\\nqOz1sasTuAcVzOmF4Ht8u8W37G654uyURs97nCMCgYEA0lIoqnup9WeAqxt7kw3S\\n9V0eGt/xbb5GpWRUdTWU2iJ3BgqZyUfYj4l/Gx4/sUFLX6jYKij0Io/jztFiN6Or\\nuWvaR688Y/t6tftzJIqt0Pr9hRrvsye1rzfFx6Ih9S5tx85yRjjw20BtT+KnwvEn\\n5HQc0Dtn9UJTFGYrDgSpJeE=\\n-----END PRIVATE KEY-----",
					"server_cert": "-----BEGIN CERTIFICATE-----\\nMIIESTCCArGgAwIBAgIRAMYShh0RZPfqjGoYAcYkcGQwDQYJKoZIhvcNAQELBQAw\\nczEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSQwIgYDVQQLDBttdXhp\\nbkBCLVFUUFhIMzBOLTIzMTIubG9jYWwxKzApBgNVBAMMIm1rY2VydCBtdXhpbkBC\\nLVFUUFhIMzBOLTIzMTIubG9jYWwwHhcNMjUwOTA4MTQ0ODA2WhcNMjcxMjA4MTQ0\\nODA2WjBPMScwJQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1lbnQgY2VydGlmaWNhdGUx\\nJDAiBgNVBAsMG211eGluQEItUVRQWEgzME4tMjMxMi5sb2NhbDCCASIwDQYJKoZI\\nhvcNAQEBBQADggEPADCCAQoCggEBAMpHcZKOUW8/KRq0T3pmPXluxTWf5O0imAbB\\nVC907ecrxLReleE/70lCVVxxOf2gOkm0qhYxGLgO4c4V2rNI2u8Jepw9Rof7uojV\\n1WhglBHIjcl2h+MB2cqkDieBlM2WQDHfIjZCNPi9rLpu464H4og688cmIH+z8lmS\\nEucSbuadStmCu0Kornqo/zzxRUAWlUQ+xqdlTYgQdFEbsRNHBty8OIiaMiR3XmDd\\nvrDXEceJKu7drQ8gv4G2ECEYqQPfzjNlaaVkIhBv00cxV5k8gbVIgg9JB8NqQuxj\\nwqd5exCcJi/Cq87pt1Oovk1QeU/I1WT6njy2W0V55x7toqKg6PMCAwEAAaN8MHow\\nDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8GA1UdIwQYMBaA\\nFFNGK+jyz7hgjwBTROub+zLMn1CgMDIGA1UdEQQrMCmCCWxvY2FsaG9zdIcECJOK\\nKYcEfwAAAYcQAAAAAAAAAAAAAAAAAAAAATANBgkqhkiG9w0BAQsFAAOCAYEAdsG6\\n7NnFu+fYI9RQmfZG9imiBkkDWW1CY1ZYxBstYd974Lqig0IKU+waVQBAlPOut5Wy\\n6J2SLDvl165mlzvoRiz743Djq03CWaGk2mvrB5vJBTJN+Pe5bV4FgKC3xynuPR8T\\n1NQ8WmBxGhwBpNnpU5q8hxc6hh7oHjJH60g886P+79aZnACSec8hYrg6IxbdtT1c\\nge+Qs7kMBCBFVHuosDQgg2XSgUFf0V73G5IJXcHVa8yxV3/M69tA/ECSfq4qxhgz\\nbakrSdmGNzkqYYEX6GSXxGBDjnq79E6VQDB9oCFD1i9Db2yTImr0wLlvHDQWx21D\\nOEADBGPc8nYHkTSvLERYoz6PHd/OTrC6k6hFtSWINtRWlHanUZfmPun+5UKEibMI\\nhIX0Pjz/WRBfYJBqdxVENa5Tu0goVLxFD4/MtvQPqj9Y863l1vTqEVp5ymMoxE76\\nfeZbson9UP+vzcnY1QAsdC77/T3eqXkM8xrf1b+y8PyONo7dtbbmT1e+Vjb0\\n-----END CERTIFICATE-----",
					"ssl_enabled": "1",
					"ca_type":     "custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_key":  "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKR3GSjlFvPyka\ntE96Zj15bsU1n+TtIpgGwVQvdO3nK8S0XpXhP+9JQlVccTn9oDpJtKoWMRi4DuHO\nFdqzSNrvCXqcPUaH+7qI1dVoYJQRyI3JdofjAdnKpA4ngZTNlkAx3yI2QjT4vay6\nbuOuB+KIOvPHJiB/s/JZkhLnEm7mnUrZgrtCqK56qP888UVAFpVEPsanZU2IEHRR\nG7ETRwbcvDiImjIkd15g3b6w1xHHiSru3a0PIL+BthAhGKkD384zZWmlZCIQb9NH\nMVeZPIG1SIIPSQfDakLsY8KneXsQnCYvwqvO6bdTqL5NUHlPyNVk+p48tltFeece\n7aKioOjzAgMBAAECggEAJtGZIRgA1smXONG7ovC2AXTZkdXyl/OYm0tEvarB8Sg8\nIqU4PDcJ09RQD2KHT34NUZHDRmj7pm7stKELDHcB1PfLuOole/k6LgJjZxmJsPP9\nCdmecFktk67yLHC4vs+D2E5LAYCpK8cyu8CGHyLSPXSazfAqMne1Ha1jxUaLU+Qq\nPsmGX2O808R9DTrdcd7fzT9MwBiAbO3zX0fDQeAXF/NFeuty2YU3hMPxrgm5xirL\nZnbShZJ7uZPXNslOSK8qJTcGjYbBAPydAKWTS8b1fww78dqvUIEUPiqi9CfbH8T5\nvwq+pROIbBORhMBxgFfeNZg0PW3r0kRNYi7SGNNzcQKBgQDwBnXQUt9AG6rCaF4c\njiwNr3l4r0wA3ltrNSsxvgw8JJ4+mn0J6NfYpVFfXxteXInGAW1kWOT2U2rUIRLm\nk4nt48OVBiYMKxk1aqqWRjP+E51Kv3ZJjUGiBuzFLaAlELgr7TeH9M1s9EiujSgx\nmO44vNtZxnvhx2799wdM5XkT2QKBgQDXvd2xFzDuEM37Al/aAmmCJNvl8vGojiOt\nWIhKilke8h7NVEum3A6OV80htMp2Nrxa3fbuGv0g2NA3H2Ddh/ALqJ56DZlYRwxc\n/GkJ4aHr3gvMgmRBNcOSrw8zEutntzj3y/KRdxofbqloTpmHLy7T/dO25ycNGaax\nFuEOsCh/qwKBgH+C7xO88t0b1Ztx1o1U+hJLJjz242mStv49QLUsQVOyIF8hs0uQ\nZxqwuInx/JgkkQyftX2ZvAkgR6Bp8aCMwLmgRkbk/VF5k+rMv9MVeImB4g3TqQNq\nB3QMObyGgI0wVKcBXn7bjkZTgEk6tB+lHukFa4JF74oCaPSCR16SicABAoGAJ9yX\n8pmTW9lVBbTpmvFpJzfCPZmG1xr0MpoyDHvFfbdEw7F4jOsJ8Xj9mOW7pt30LWHn\ntxTbgk8tIZI1ppjwXGcaaPjMRYhTG3czvTSoY1lSmsXY2kehzB92UnyDbFVpPDe9\nqOz1sasTuAcVzOmF4Ht8u8W37G654uyURs97nCMCgYEA0lIoqnup9WeAqxt7kw3S\n9V0eGt/xbb5GpWRUdTWU2iJ3BgqZyUfYj4l/Gx4/sUFLX6jYKij0Io/jztFiN6Or\nuWvaR688Y/t6tftzJIqt0Pr9hRrvsye1rzfFx6Ih9S5tx85yRjjw20BtT+KnwvEn\n5HQc0Dtn9UJTFGYrDgSpJeE=\n-----END PRIVATE KEY-----",
						"server_cert": "-----BEGIN CERTIFICATE-----\nMIIESTCCArGgAwIBAgIRAMYShh0RZPfqjGoYAcYkcGQwDQYJKoZIhvcNAQELBQAw\nczEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSQwIgYDVQQLDBttdXhp\nbkBCLVFUUFhIMzBOLTIzMTIubG9jYWwxKzApBgNVBAMMIm1rY2VydCBtdXhpbkBC\nLVFUUFhIMzBOLTIzMTIubG9jYWwwHhcNMjUwOTA4MTQ0ODA2WhcNMjcxMjA4MTQ0\nODA2WjBPMScwJQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1lbnQgY2VydGlmaWNhdGUx\nJDAiBgNVBAsMG211eGluQEItUVRQWEgzME4tMjMxMi5sb2NhbDCCASIwDQYJKoZI\nhvcNAQEBBQADggEPADCCAQoCggEBAMpHcZKOUW8/KRq0T3pmPXluxTWf5O0imAbB\nVC907ecrxLReleE/70lCVVxxOf2gOkm0qhYxGLgO4c4V2rNI2u8Jepw9Rof7uojV\n1WhglBHIjcl2h+MB2cqkDieBlM2WQDHfIjZCNPi9rLpu464H4og688cmIH+z8lmS\nEucSbuadStmCu0Kornqo/zzxRUAWlUQ+xqdlTYgQdFEbsRNHBty8OIiaMiR3XmDd\nvrDXEceJKu7drQ8gv4G2ECEYqQPfzjNlaaVkIhBv00cxV5k8gbVIgg9JB8NqQuxj\nwqd5exCcJi/Cq87pt1Oovk1QeU/I1WT6njy2W0V55x7toqKg6PMCAwEAAaN8MHow\nDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8GA1UdIwQYMBaA\nFFNGK+jyz7hgjwBTROub+zLMn1CgMDIGA1UdEQQrMCmCCWxvY2FsaG9zdIcECJOK\nKYcEfwAAAYcQAAAAAAAAAAAAAAAAAAAAATANBgkqhkiG9w0BAQsFAAOCAYEAdsG6\n7NnFu+fYI9RQmfZG9imiBkkDWW1CY1ZYxBstYd974Lqig0IKU+waVQBAlPOut5Wy\n6J2SLDvl165mlzvoRiz743Djq03CWaGk2mvrB5vJBTJN+Pe5bV4FgKC3xynuPR8T\n1NQ8WmBxGhwBpNnpU5q8hxc6hh7oHjJH60g886P+79aZnACSec8hYrg6IxbdtT1c\nge+Qs7kMBCBFVHuosDQgg2XSgUFf0V73G5IJXcHVa8yxV3/M69tA/ECSfq4qxhgz\nbakrSdmGNzkqYYEX6GSXxGBDjnq79E6VQDB9oCFD1i9Db2yTImr0wLlvHDQWx21D\nOEADBGPc8nYHkTSvLERYoz6PHd/OTrC6k6hFtSWINtRWlHanUZfmPun+5UKEibMI\nhIX0Pjz/WRBfYJBqdxVENa5Tu0goVLxFD4/MtvQPqj9Y863l1vTqEVp5ymMoxE76\nfeZbson9UP+vzcnY1QAsdC77/T3eqXkM8xrf1b+y8PyONo7dtbbmT1e+Vjb0\n-----END CERTIFICATE-----",
						"ssl_enabled": "1",
						"ca_type":     "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dashboard_password": "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"database_password": "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_config_list": []map[string]interface{}{
						{
							"name":  "FILE_SIZE_LIMIT",
							"value": "52428801",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_config_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dashboard_password", "database_password", "initialize_with_existing_data", "public_endpoint_enabled", "public_network_access_enabled", "auth_config_list", "storage_config_list"},
			},
		},
	})
}

func TestAccAliCloudRdsAiInstance_basic12309_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_ai_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsAiInstanceMap12309)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsAiServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAiInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrdsai%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsAiInstanceBasicDependence12309Twin)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":                      name,
					"app_type":                      "supabase",
					"db_instance_name":              "${data.alicloud_db_instances.default.instances.0.id}",
					"server_key":                    "-----BEGIN PRIVATE KEY-----\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKR3GSjlFvPyka\\ntE96Zj15bsU1n+TtIpgGwVQvdO3nK8S0XpXhP+9JQlVccTn9oDpJtKoWMRi4DuHO\\nFdqzSNrvCXqcPUaH+7qI1dVoYJQRyI3JdofjAdnKpA4ngZTNlkAx3yI2QjT4vay6\\nbuOuB+KIOvPHJiB/s/JZkhLnEm7mnUrZgrtCqK56qP888UVAFpVEPsanZU2IEHRR\\nG7ETRwbcvDiImjIkd15g3b6w1xHHiSru3a0PIL+BthAhGKkD384zZWmlZCIQb9NH\\nMVeZPIG1SIIPSQfDakLsY8KneXsQnCYvwqvO6bdTqL5NUHlPyNVk+p48tltFeece\\n7aKioOjzAgMBAAECggEAJtGZIRgA1smXONG7ovC2AXTZkdXyl/OYm0tEvarB8Sg8\\nIqU4PDcJ09RQD2KHT34NUZHDRmj7pm7stKELDHcB1PfLuOole/k6LgJjZxmJsPP9\\nCdmecFktk67yLHC4vs+D2E5LAYCpK8cyu8CGHyLSPXSazfAqMne1Ha1jxUaLU+Qq\\nPsmGX2O808R9DTrdcd7fzT9MwBiAbO3zX0fDQeAXF/NFeuty2YU3hMPxrgm5xirL\\nZnbShZJ7uZPXNslOSK8qJTcGjYbBAPydAKWTS8b1fww78dqvUIEUPiqi9CfbH8T5\\nvwq+pROIbBORhMBxgFfeNZg0PW3r0kRNYi7SGNNzcQKBgQDwBnXQUt9AG6rCaF4c\\njiwNr3l4r0wA3ltrNSsxvgw8JJ4+mn0J6NfYpVFfXxteXInGAW1kWOT2U2rUIRLm\\nk4nt48OVBiYMKxk1aqqWRjP+E51Kv3ZJjUGiBuzFLaAlELgr7TeH9M1s9EiujSgx\\nmO44vNtZxnvhx2799wdM5XkT2QKBgQDXvd2xFzDuEM37Al/aAmmCJNvl8vGojiOt\\nWIhKilke8h7NVEum3A6OV80htMp2Nrxa3fbuGv0g2NA3H2Ddh/ALqJ56DZlYRwxc\\n/GkJ4aHr3gvMgmRBNcOSrw8zEutntzj3y/KRdxofbqloTpmHLy7T/dO25ycNGaax\\nFuEOsCh/qwKBgH+C7xO88t0b1Ztx1o1U+hJLJjz242mStv49QLUsQVOyIF8hs0uQ\\nZxqwuInx/JgkkQyftX2ZvAkgR6Bp8aCMwLmgRkbk/VF5k+rMv9MVeImB4g3TqQNq\\nB3QMObyGgI0wVKcBXn7bjkZTgEk6tB+lHukFa4JF74oCaPSCR16SicABAoGAJ9yX\\n8pmTW9lVBbTpmvFpJzfCPZmG1xr0MpoyDHvFfbdEw7F4jOsJ8Xj9mOW7pt30LWHn\\ntxTbgk8tIZI1ppjwXGcaaPjMRYhTG3czvTSoY1lSmsXY2kehzB92UnyDbFVpPDe9\\nqOz1sasTuAcVzOmF4Ht8u8W37G654uyURs97nCMCgYEA0lIoqnup9WeAqxt7kw3S\\n9V0eGt/xbb5GpWRUdTWU2iJ3BgqZyUfYj4l/Gx4/sUFLX6jYKij0Io/jztFiN6Or\\nuWvaR688Y/t6tftzJIqt0Pr9hRrvsye1rzfFx6Ih9S5tx85yRjjw20BtT+KnwvEn\\n5HQc0Dtn9UJTFGYrDgSpJeE=\\n-----END PRIVATE KEY-----",
					"server_cert":                   "-----BEGIN CERTIFICATE-----\\nMIIESTCCArGgAwIBAgIRAMYShh0RZPfqjGoYAcYkcGQwDQYJKoZIhvcNAQELBQAw\\nczEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSQwIgYDVQQLDBttdXhp\\nbkBCLVFUUFhIMzBOLTIzMTIubG9jYWwxKzApBgNVBAMMIm1rY2VydCBtdXhpbkBC\\nLVFUUFhIMzBOLTIzMTIubG9jYWwwHhcNMjUwOTA4MTQ0ODA2WhcNMjcxMjA4MTQ0\\nODA2WjBPMScwJQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1lbnQgY2VydGlmaWNhdGUx\\nJDAiBgNVBAsMG211eGluQEItUVRQWEgzME4tMjMxMi5sb2NhbDCCASIwDQYJKoZI\\nhvcNAQEBBQADggEPADCCAQoCggEBAMpHcZKOUW8/KRq0T3pmPXluxTWf5O0imAbB\\nVC907ecrxLReleE/70lCVVxxOf2gOkm0qhYxGLgO4c4V2rNI2u8Jepw9Rof7uojV\\n1WhglBHIjcl2h+MB2cqkDieBlM2WQDHfIjZCNPi9rLpu464H4og688cmIH+z8lmS\\nEucSbuadStmCu0Kornqo/zzxRUAWlUQ+xqdlTYgQdFEbsRNHBty8OIiaMiR3XmDd\\nvrDXEceJKu7drQ8gv4G2ECEYqQPfzjNlaaVkIhBv00cxV5k8gbVIgg9JB8NqQuxj\\nwqd5exCcJi/Cq87pt1Oovk1QeU/I1WT6njy2W0V55x7toqKg6PMCAwEAAaN8MHow\\nDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8GA1UdIwQYMBaA\\nFFNGK+jyz7hgjwBTROub+zLMn1CgMDIGA1UdEQQrMCmCCWxvY2FsaG9zdIcECJOK\\nKYcEfwAAAYcQAAAAAAAAAAAAAAAAAAAAATANBgkqhkiG9w0BAQsFAAOCAYEAdsG6\\n7NnFu+fYI9RQmfZG9imiBkkDWW1CY1ZYxBstYd974Lqig0IKU+waVQBAlPOut5Wy\\n6J2SLDvl165mlzvoRiz743Djq03CWaGk2mvrB5vJBTJN+Pe5bV4FgKC3xynuPR8T\\n1NQ8WmBxGhwBpNnpU5q8hxc6hh7oHjJH60g886P+79aZnACSec8hYrg6IxbdtT1c\\nge+Qs7kMBCBFVHuosDQgg2XSgUFf0V73G5IJXcHVa8yxV3/M69tA/ECSfq4qxhgz\\nbakrSdmGNzkqYYEX6GSXxGBDjnq79E6VQDB9oCFD1i9Db2yTImr0wLlvHDQWx21D\\nOEADBGPc8nYHkTSvLERYoz6PHd/OTrC6k6hFtSWINtRWlHanUZfmPun+5UKEibMI\\nhIX0Pjz/WRBfYJBqdxVENa5Tu0goVLxFD4/MtvQPqj9Y863l1vTqEVp5ymMoxE76\\nfeZbson9UP+vzcnY1QAsdC77/T3eqXkM8xrf1b+y8PyONo7dtbbmT1e+Vjb0\\n-----END CERTIFICATE-----",
					"ssl_enabled":                   "1",
					"ca_type":                       "custom",
					"dashboard_password":            "YourPassword123",
					"database_password":             "YourPassword123",
					"initialize_with_existing_data": "true",
					"public_endpoint_enabled":       "true",
					"public_network_access_enabled": "true",
					"status":                        "Running",
					"auth_config_list": []map[string]interface{}{
						{
							"name":  "API_EXTERNAL_URL",
							"value": "http://localhost:3001",
						},
					},
					"storage_config_list": []map[string]interface{}{
						{
							"name":  "FILE_SIZE_LIMIT",
							"value": "52428801",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":              name,
						"app_type":              "supabase",
						"db_instance_name":      CHECKSET,
						"server_key":            "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKR3GSjlFvPyka\ntE96Zj15bsU1n+TtIpgGwVQvdO3nK8S0XpXhP+9JQlVccTn9oDpJtKoWMRi4DuHO\nFdqzSNrvCXqcPUaH+7qI1dVoYJQRyI3JdofjAdnKpA4ngZTNlkAx3yI2QjT4vay6\nbuOuB+KIOvPHJiB/s/JZkhLnEm7mnUrZgrtCqK56qP888UVAFpVEPsanZU2IEHRR\nG7ETRwbcvDiImjIkd15g3b6w1xHHiSru3a0PIL+BthAhGKkD384zZWmlZCIQb9NH\nMVeZPIG1SIIPSQfDakLsY8KneXsQnCYvwqvO6bdTqL5NUHlPyNVk+p48tltFeece\n7aKioOjzAgMBAAECggEAJtGZIRgA1smXONG7ovC2AXTZkdXyl/OYm0tEvarB8Sg8\nIqU4PDcJ09RQD2KHT34NUZHDRmj7pm7stKELDHcB1PfLuOole/k6LgJjZxmJsPP9\nCdmecFktk67yLHC4vs+D2E5LAYCpK8cyu8CGHyLSPXSazfAqMne1Ha1jxUaLU+Qq\nPsmGX2O808R9DTrdcd7fzT9MwBiAbO3zX0fDQeAXF/NFeuty2YU3hMPxrgm5xirL\nZnbShZJ7uZPXNslOSK8qJTcGjYbBAPydAKWTS8b1fww78dqvUIEUPiqi9CfbH8T5\nvwq+pROIbBORhMBxgFfeNZg0PW3r0kRNYi7SGNNzcQKBgQDwBnXQUt9AG6rCaF4c\njiwNr3l4r0wA3ltrNSsxvgw8JJ4+mn0J6NfYpVFfXxteXInGAW1kWOT2U2rUIRLm\nk4nt48OVBiYMKxk1aqqWRjP+E51Kv3ZJjUGiBuzFLaAlELgr7TeH9M1s9EiujSgx\nmO44vNtZxnvhx2799wdM5XkT2QKBgQDXvd2xFzDuEM37Al/aAmmCJNvl8vGojiOt\nWIhKilke8h7NVEum3A6OV80htMp2Nrxa3fbuGv0g2NA3H2Ddh/ALqJ56DZlYRwxc\n/GkJ4aHr3gvMgmRBNcOSrw8zEutntzj3y/KRdxofbqloTpmHLy7T/dO25ycNGaax\nFuEOsCh/qwKBgH+C7xO88t0b1Ztx1o1U+hJLJjz242mStv49QLUsQVOyIF8hs0uQ\nZxqwuInx/JgkkQyftX2ZvAkgR6Bp8aCMwLmgRkbk/VF5k+rMv9MVeImB4g3TqQNq\nB3QMObyGgI0wVKcBXn7bjkZTgEk6tB+lHukFa4JF74oCaPSCR16SicABAoGAJ9yX\n8pmTW9lVBbTpmvFpJzfCPZmG1xr0MpoyDHvFfbdEw7F4jOsJ8Xj9mOW7pt30LWHn\ntxTbgk8tIZI1ppjwXGcaaPjMRYhTG3czvTSoY1lSmsXY2kehzB92UnyDbFVpPDe9\nqOz1sasTuAcVzOmF4Ht8u8W37G654uyURs97nCMCgYEA0lIoqnup9WeAqxt7kw3S\n9V0eGt/xbb5GpWRUdTWU2iJ3BgqZyUfYj4l/Gx4/sUFLX6jYKij0Io/jztFiN6Or\nuWvaR688Y/t6tftzJIqt0Pr9hRrvsye1rzfFx6Ih9S5tx85yRjjw20BtT+KnwvEn\n5HQc0Dtn9UJTFGYrDgSpJeE=\n-----END PRIVATE KEY-----",
						"server_cert":           "-----BEGIN CERTIFICATE-----\nMIIESTCCArGgAwIBAgIRAMYShh0RZPfqjGoYAcYkcGQwDQYJKoZIhvcNAQELBQAw\nczEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSQwIgYDVQQLDBttdXhp\nbkBCLVFUUFhIMzBOLTIzMTIubG9jYWwxKzApBgNVBAMMIm1rY2VydCBtdXhpbkBC\nLVFUUFhIMzBOLTIzMTIubG9jYWwwHhcNMjUwOTA4MTQ0ODA2WhcNMjcxMjA4MTQ0\nODA2WjBPMScwJQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1lbnQgY2VydGlmaWNhdGUx\nJDAiBgNVBAsMG211eGluQEItUVRQWEgzME4tMjMxMi5sb2NhbDCCASIwDQYJKoZI\nhvcNAQEBBQADggEPADCCAQoCggEBAMpHcZKOUW8/KRq0T3pmPXluxTWf5O0imAbB\nVC907ecrxLReleE/70lCVVxxOf2gOkm0qhYxGLgO4c4V2rNI2u8Jepw9Rof7uojV\n1WhglBHIjcl2h+MB2cqkDieBlM2WQDHfIjZCNPi9rLpu464H4og688cmIH+z8lmS\nEucSbuadStmCu0Kornqo/zzxRUAWlUQ+xqdlTYgQdFEbsRNHBty8OIiaMiR3XmDd\nvrDXEceJKu7drQ8gv4G2ECEYqQPfzjNlaaVkIhBv00cxV5k8gbVIgg9JB8NqQuxj\nwqd5exCcJi/Cq87pt1Oovk1QeU/I1WT6njy2W0V55x7toqKg6PMCAwEAAaN8MHow\nDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8GA1UdIwQYMBaA\nFFNGK+jyz7hgjwBTROub+zLMn1CgMDIGA1UdEQQrMCmCCWxvY2FsaG9zdIcECJOK\nKYcEfwAAAYcQAAAAAAAAAAAAAAAAAAAAATANBgkqhkiG9w0BAQsFAAOCAYEAdsG6\n7NnFu+fYI9RQmfZG9imiBkkDWW1CY1ZYxBstYd974Lqig0IKU+waVQBAlPOut5Wy\n6J2SLDvl165mlzvoRiz743Djq03CWaGk2mvrB5vJBTJN+Pe5bV4FgKC3xynuPR8T\n1NQ8WmBxGhwBpNnpU5q8hxc6hh7oHjJH60g886P+79aZnACSec8hYrg6IxbdtT1c\nge+Qs7kMBCBFVHuosDQgg2XSgUFf0V73G5IJXcHVa8yxV3/M69tA/ECSfq4qxhgz\nbakrSdmGNzkqYYEX6GSXxGBDjnq79E6VQDB9oCFD1i9Db2yTImr0wLlvHDQWx21D\nOEADBGPc8nYHkTSvLERYoz6PHd/OTrC6k6hFtSWINtRWlHanUZfmPun+5UKEibMI\nhIX0Pjz/WRBfYJBqdxVENa5Tu0goVLxFD4/MtvQPqj9Y863l1vTqEVp5ymMoxE76\nfeZbson9UP+vzcnY1QAsdC77/T3eqXkM8xrf1b+y8PyONo7dtbbmT1e+Vjb0\n-----END CERTIFICATE-----",
						"ssl_enabled":           "1",
						"ca_type":               "custom",
						"status":                "Running",
						"auth_config_list.#":    "1",
						"storage_config_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dashboard_password", "database_password", "initialize_with_existing_data", "public_endpoint_enabled", "public_network_access_enabled", "auth_config_list", "storage_config_list"},
			},
		},
	})
}

var AliCloudRdsAiInstanceMap12309 = map[string]string{
	"status": CHECKSET,
}

func AliCloudRdsAiInstanceBasicDependence12309(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vswitches" "default" {
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "17.0"
  db_instance_storage_type = "general_essd"
  instance_type            = "pg.n2.1c.1m"
  instance_storage         = 100
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
}
`, name)
}

func AliCloudRdsAiInstanceBasicDependence12309Twin(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_db_instances" "default" {
  name_regex = "^default-NODELETING$"
  status     = "Running"
  engine     = "PostgreSQL"
}
`, name)
}

// Test RdsAi Instance. <<< Resource test cases, automatically generated.
