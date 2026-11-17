package slb_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	slbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSLBDomainExtensionBasic(t *testing.T) {
	var v *slbsdk.DescribeDomainExtensionAttributeResponse
	resourceId := "alicloud_slb_domain_extension.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return slb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDomainExtensionAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccDomainExtensionConfigSpot%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSlbDomainExtensionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		// module name
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":      "${alicloud_slb_load_balancer.instance.id}",
					"frontend_port":         "${alicloud_slb_listener.https.frontend_port}",
					"domain":                "www.test.com",
					"server_certificate_id": "${alicloud_slb_server_certificate.foo.id}",
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":      CHECKSET,
						"listener_port":         "443",
						"frontend_port":         "443",
						"domain":                "www.test.com",
						"server_certificate_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":      "${alicloud_slb_load_balancer.instance.id}",
					"frontend_port":         "${alicloud_slb_listener.https.frontend_port}",
					"domain":                "www.test.com",
					"server_certificate_id": "${alicloud_slb_server_certificate.new.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":      CHECKSET,
						"listener_port":         "443",
						"frontend_port":         "443",
						"domain":                "www.test.com",
						"server_certificate_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAliCloudSLBDomainExtensionListenerPort(t *testing.T) {
	var v *slbsdk.DescribeDomainExtensionAttributeResponse
	resourceId := "alicloud_slb_domain_extension.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return slb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDomainExtensionAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccDomainExtensionListenerPort%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSlbDomainExtensionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		// module name
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":             "${alicloud_slb_load_balancer.instance.id}",
					"listener_port":                "${alicloud_slb_listener.https.frontend_port}",
					"domain":                       "www.test.com",
					"server_certificate_id":        "${alicloud_slb_server_certificate.foo.id}",
					"delete_protection_validation": false,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":      CHECKSET,
						"listener_port":         "443",
						"frontend_port":         "443",
						"domain":                "www.test.com",
						"server_certificate_id": CHECKSET,
					}),
					func(state *terraform.State) error {
						resourceState, ok := state.RootModule().Resources[resourceId]
						if !ok || resourceState.Primary == nil || resourceState.Primary.ID == "" {
							return fmt.Errorf("G1 evidence: %s has no primary ID", resourceId)
						}
						client := acctest.Provider.Meta().(*connectivity.AliyunClient)
						if client.RegionId != "cn-beijing" {
							return fmt.Errorf("G1 evidence: expected cn-beijing, got %q", client.RegionId)
						}
						accountId, err := client.AccountId()
						if err != nil {
							return fmt.Errorf("G1 evidence: resolve account ID: %w", err)
						}
						request := slbsdk.CreateDescribeDomainExtensionAttributeRequest()
						request.DomainExtensionId = resourceState.Primary.ID
						requestStartedAt := time.Now().UTC()
						raw, err := client.WithSlbClient(func(slbClient *slbsdk.Client) (interface{}, error) {
							return slbClient.DescribeDomainExtensionAttribute(request)
						})
						responseReceivedAt := time.Now().UTC()
						if err != nil {
							return fmt.Errorf("G1 evidence: DescribeDomainExtensionAttribute: %w", err)
						}
						response, ok := raw.(*slbsdk.DescribeDomainExtensionAttributeResponse)
						if !ok || response == nil || response.BaseResponse == nil {
							return fmt.Errorf("G1 evidence: unexpected SDK response %T", raw)
						}
						if response.GetHttpStatus() < 200 || response.GetHttpStatus() >= 300 {
							return fmt.Errorf("G1 evidence: unexpected HTTP status %d", response.GetHttpStatus())
						}
						rawBody := response.GetHttpContentBytes()
						var wireResponse struct {
							DomainExtensionId string `json:"DomainExtensionId"`
							ListenerPort      int    `json:"ListenerPort"`
							RequestId         string `json:"RequestId"`
						}
						if err := json.Unmarshal(rawBody, &wireResponse); err != nil {
							return fmt.Errorf("G1 evidence: decode raw HTTP response: %w", err)
						}
						if wireResponse.DomainExtensionId != resourceState.Primary.ID || wireResponse.ListenerPort != 443 || wireResponse.RequestId == "" {
							return fmt.Errorf("G1 evidence: raw HTTP response does not contain the target identity, ListenerPort=443 and RequestId")
						}
						if response.DomainExtensionId != resourceState.Primary.ID || response.ListenerPort != 443 || response.RequestId == "" {
							return fmt.Errorf("G1 evidence: SDK response does not contain the target identity, ListenerPort=443 and RequestId")
						}
						if response.RequestId != wireResponse.RequestId {
							return fmt.Errorf("G1 evidence: raw HTTP and SDK RequestId values differ")
						}
						contextJson, err := json.Marshal(struct {
							AccountId          string `json:"accountId"`
							Action             string `json:"action"`
							DomainExtensionId  string `json:"domainExtensionId"`
							HttpStatus         int    `json:"httpStatus"`
							Region             string `json:"region"`
							RequestId          string `json:"requestId"`
							RequestStartedAt   string `json:"requestStartedAt"`
							ResponseReceivedAt string `json:"responseReceivedAt"`
						}{
							AccountId:          accountId,
							Action:             request.GetActionName(),
							DomainExtensionId:  resourceState.Primary.ID,
							HttpStatus:         response.GetHttpStatus(),
							Region:             client.RegionId,
							RequestId:          response.RequestId,
							RequestStartedAt:   requestStartedAt.Format(time.RFC3339Nano),
							ResponseReceivedAt: responseReceivedAt.Format(time.RFC3339Nano),
						})
						if err != nil {
							return fmt.Errorf("G1 evidence: encode request context: %w", err)
						}
						decodedJson, err := json.Marshal(response)
						if err != nil {
							return fmt.Errorf("G1 evidence: encode SDK response: %w", err)
						}
						fmt.Printf("G1_DESCRIBE_REQUEST_CONTEXT=%s\n", contextJson)
						fmt.Printf("G1_DESCRIBE_RAW_HTTP_RESPONSE=%s\n", rawBody)
						fmt.Printf("G1_DESCRIBE_SDK_DECODED_RESPONSE=%s\n", decodedJson)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

// listener_port and frontend_port are mutually exclusive and at least one is
// required; both rules are enforced by schema ExactlyOneOf at plan time, so a
// config-level unit test locks them without creating cloud resources.
func TestUnitAlicloudSLBDomainExtensionPortExactlyOneOf(t *testing.T) {
	baseConfig := map[string]interface{}{
		"load_balancer_id":      "lb-bp1unitcheck",
		"domain":                "www.test.com",
		"server_certificate_id": "cert-bp1unitcheck",
	}

	configWith := func(ports map[string]interface{}) map[string]interface{} {
		config := make(map[string]interface{}, len(baseConfig)+len(ports))
		for key, value := range baseConfig {
			config[key] = value
		}
		for key, value := range ports {
			config[key] = value
		}
		return config
	}

	errorDetails := func(config map[string]interface{}) []string {
		var details []string
		for _, diagnostic := range slb.ResourceDomainExtension().Validate(terraform.NewResourceConfigRaw(config)) {
			if diagnostic.Severity == diag.Error {
				details = append(details, diagnostic.Detail)
			}
		}
		return details
	}

	t.Run("both ports configured", func(t *testing.T) {
		details := errorDetails(configWith(map[string]interface{}{
			"listener_port": 443,
			"frontend_port": 443,
		}))
		want := "only one of `frontend_port,listener_port` can be specified, but `frontend_port,listener_port` were specified."
		if !strings.Contains(strings.Join(details, "\n"), want) {
			t.Fatalf("want rejection reason %q, got %v", want, details)
		}
	})

	t.Run("neither port configured", func(t *testing.T) {
		details := errorDetails(configWith(nil))
		want := "one of `frontend_port,listener_port` must be specified"
		if !strings.Contains(strings.Join(details, "\n"), want) {
			t.Fatalf("want rejection reason %q, got %v", want, details)
		}
	})

	t.Run("single port configured", func(t *testing.T) {
		if details := errorDetails(configWith(map[string]interface{}{"listener_port": 443})); len(details) != 0 {
			t.Fatalf("want no errors, got %v", details)
		}
	})
}

func resourceSlbDomainExtensionConfigDependence(name string) string {
	return fmt.Sprintf(`
   variable "name" {
		default = "%s"
	}
   resource "alicloud_slb_load_balancer" "instance" {
        load_balancer_name                 = "${var.name}"
        internet_charge_type = "PayByTraffic"
        address_type             = "internet"
		load_balancer_spec        = "slb.s2.small"
	}
	resource "alicloud_slb_server_certificate" "foo" {
  		name               = "${var.name}"
  		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm+erkcKN7DANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQGEwJj\nbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmFsaXl1\nbjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0LmNvbTEaMBgGCSqGSIb3\nDQEJARYLMTIzQDEyMy5jb20wHhcNMTkwNDI2MDM0ODAxWhcNMjQwNDI1MDM0ODAx\nWjB9MQswCQYDVQQGEwJjbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcx\nDzANBgNVBAoMBmFsaXl1bjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0\nLmNvbTEaMBgGCSqGSIb3DQEJARYLMTIzQDEyMy5jb20wggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDKMKF5qmN/uoMjdH3D8aPRcUOA0s8rZpYhG8zbkF1j\n8gHYoB/FDvM7G7dfVsyjbMwLOxKvAhWvHHSpEz/t7gB+QdwrAMiMJwGmtCnXrh2E\nWiXgalMe1y4a/T5R7q+m4T1zFATf+kbnHWfkSGF4W7b6UBoaH+9StQ95CnqzNf/2\np/Of7+S0XzCxFXw8GIVzZk0xFe6lHJzaq06f3mvzrD+4rpO56tTUvrgTY/n61gsF\nZP7f0CJ2JQh6eNRFOEUSfxKu/Dy/+IsQxorCJY2Q59ZAf3rXrqDN104jw9PlwnLl\nqfZz3RMODN6BWjxE8rvRtT0qMfuAfv1gjBdWZN0hUYBRAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAABzo82TxGp5poVkd5pLWj5ACgcBv8Cs6oH9D+4Jz9BmyuBUsQXh\n2aG0hQAe1mU61C9konsl/GTW8umJQ4M4lYEztXXwMf5PlBMGwebM0ZbSGg6jKtZg\nWCgJ3eP/FMmyXGL5Jji5+e09eObhUDVle4tdi0On97zBoz85W02rgWFAqZJwiEAP\nt+c7jX7uOSBq2/38iGStlrX5yB1at/gJXXiA5CL5OtlR3Okvb0/QH37efO1Nu39m\nlFi0ODPAVyXjVypAiLguDxPn6AtDTdk9Iw9B19OD4NrzNRWgSSX5vuxo/VcRcgWk\n3gEe9Ca0ZKN20q9XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
  		private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyjCheapjf7qDI3R9w/Gj0XFDgNLPK2aWIRvM25BdY/IB2KAf\nxQ7zOxu3X1bMo2zMCzsSrwIVrxx0qRM/7e4AfkHcKwDIjCcBprQp164dhFol4GpT\nHtcuGv0+Ue6vpuE9cxQE3/pG5x1n5EhheFu2+lAaGh/vUrUPeQp6szX/9qfzn+/k\ntF8wsRV8PBiFc2ZNMRXupRyc2qtOn95r86w/uK6TuerU1L64E2P5+tYLBWT+39Ai\ndiUIenjURThFEn8Srvw8v/iLEMaKwiWNkOfWQH96166gzddOI8PT5cJy5an2c90T\nDgzegVo8RPK70bU9KjH7gH79YIwXVmTdIVGAUQIDAQABAoIBAE1J4a/8biR5S3/W\nG+03BYQeY8tuyjqw8FqfoeOcf9agwAvqybouSNQjeCk9qOQfxq/UWQQFK/zQR9gJ\nv7pX7GBXFK5rkj3g+0SaQhRsPmRFgY0Tl8qGPt2aSKRRNVv5ZeADmwlzRn86QmiF\nMp0rkfqFfDTYWEepZszCML0ouzuxsW/9tq7rvtSjsgATNt31B3vFa3D3JBi31jUh\n5nfR9A3bATze7mQw3byEDiVl5ASRDgYyur403P1fDnMy9DBHZ8NaPOsFF6OBpJal\nBJsG5z00hll5PFN2jfmBQKlvAeU7wfwqdaSnGHOfqf2DeTTaFjIQ4gUhRn/m6pLo\n6kXttLECgYEA9sng0Qz/TcPFfM4tQ1gyvB1cKnnGIwg1FP8sfUjbbEgjaHhA224S\nk3BxtX2Kq6fhTXuwusAFc6OVMAZ76FgrQ5K4Ci7+DTsrF28z4b8td+p+lO/DxgP9\nlTgN+ddsiTOV4fUef9Z3yY0Zr0CnBUMbQYRaV2UIbCdiB0G4V/bt9TsCgYEA0bya\nOo9wGI0RJV0bYP7qwO74Ra1/i1viWbRlS7jU37Q+AZstrlKcQ5CTPzOjKFKMiUzl\n4miWacZ0/q2n+Mvd7NbXGXTLijahnyOYKaHJYyh4oBymfkgAifRstE0Ki9gdvArb\n/I+emC0GvLSyfGN8UUeDJs4NmqdEXGqjo2JOV+MCgYALFv1MR5o9Y1u/hQBRs2fs\nPiGDIx+9OUQxYloccyaxEfjNXAIGGkcpavchIbgWiJ++PJ2vdquIC8TLeK8evL+M\n9M3iX0Q5UfxYvD2HmnCvn9D6Xl/cyRcfGnq+TGjrLW9BzSMGuZt+aiHKV0xqFx7l\nbc4leTvMqGRmURS4lzcQOwKBgQCDzA/i4sYfN25h21tcHXSpnsG3D2rJyQi5NCo/\nZjunA92/JqOTGuiFcLGHEszhhtY3ZXJET1LNz18vtzKJnpqrvOnYXlOVW/U+SqDQ\n8JDb1c/PVZGuY1KrXkR9HLiW3kz5IJ3S3PFdUVYdeTN8BQxXCyg4V12nJJtJs912\ny0zN3wKBgGDS6YttCN6aI4EOABYE8fI1EYQ7vhfiYsaWGWSR1l6bQey7KR6M1ACz\nZzMASNyytVt12yXE4/Emv6/pYqigbDLfL1zQJSLJ3EHJYTh2RxjR+AaGDudYFG/T\nliQ9YXhV5Iu2x1pNwrtFnssDdaaGpfA7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
	}
	resource "alicloud_slb_server_certificate" "new" {
  		name               = "${var.name}"
  		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm+erkcKN7DANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQGEwJj\nbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmFsaXl1\nbjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0LmNvbTEaMBgGCSqGSIb3\nDQEJARYLMTIzQDEyMy5jb20wHhcNMTkwNDI2MDM0ODAxWhcNMjQwNDI1MDM0ODAx\nWjB9MQswCQYDVQQGEwJjbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcx\nDzANBgNVBAoMBmFsaXl1bjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0\nLmNvbTEaMBgGCSqGSIb3DQEJARYLMTIzQDEyMy5jb20wggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDKMKF5qmN/uoMjdH3D8aPRcUOA0s8rZpYhG8zbkF1j\n8gHYoB/FDvM7G7dfVsyjbMwLOxKvAhWvHHSpEz/t7gB+QdwrAMiMJwGmtCnXrh2E\nWiXgalMe1y4a/T5R7q+m4T1zFATf+kbnHWfkSGF4W7b6UBoaH+9StQ95CnqzNf/2\np/Of7+S0XzCxFXw8GIVzZk0xFe6lHJzaq06f3mvzrD+4rpO56tTUvrgTY/n61gsF\nZP7f0CJ2JQh6eNRFOEUSfxKu/Dy/+IsQxorCJY2Q59ZAf3rXrqDN104jw9PlwnLl\nqfZz3RMODN6BWjxE8rvRtT0qMfuAfv1gjBdWZN0hUYBRAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAABzo82TxGp5poVkd5pLWj5ACgcBv8Cs6oH9D+4Jz9BmyuBUsQXh\n2aG0hQAe1mU61C9konsl/GTW8umJQ4M4lYEztXXwMf5PlBMGwebM0ZbSGg6jKtZg\nWCgJ3eP/FMmyXGL5Jji5+e09eObhUDVle4tdi0On97zBoz85W02rgWFAqZJwiEAP\nt+c7jX7uOSBq2/38iGStlrX5yB1at/gJXXiA5CL5OtlR3Okvb0/QH37efO1Nu39m\nlFi0ODPAVyXjVypAiLguDxPn6AtDTdk9Iw9B19OD4NrzNRWgSSX5vuxo/VcRcgWk\n3gEe9Ca0ZKN20q9XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
  		private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyjCheapjf7qDI3R9w/Gj0XFDgNLPK2aWIRvM25BdY/IB2KAf\nxQ7zOxu3X1bMo2zMCzsSrwIVrxx0qRM/7e4AfkHcKwDIjCcBprQp164dhFol4GpT\nHtcuGv0+Ue6vpuE9cxQE3/pG5x1n5EhheFu2+lAaGh/vUrUPeQp6szX/9qfzn+/k\ntF8wsRV8PBiFc2ZNMRXupRyc2qtOn95r86w/uK6TuerU1L64E2P5+tYLBWT+39Ai\ndiUIenjURThFEn8Srvw8v/iLEMaKwiWNkOfWQH96166gzddOI8PT5cJy5an2c90T\nDgzegVo8RPK70bU9KjH7gH79YIwXVmTdIVGAUQIDAQABAoIBAE1J4a/8biR5S3/W\nG+03BYQeY8tuyjqw8FqfoeOcf9agwAvqybouSNQjeCk9qOQfxq/UWQQFK/zQR9gJ\nv7pX7GBXFK5rkj3g+0SaQhRsPmRFgY0Tl8qGPt2aSKRRNVv5ZeADmwlzRn86QmiF\nMp0rkfqFfDTYWEepZszCML0ouzuxsW/9tq7rvtSjsgATNt31B3vFa3D3JBi31jUh\n5nfR9A3bATze7mQw3byEDiVl5ASRDgYyur403P1fDnMy9DBHZ8NaPOsFF6OBpJal\nBJsG5z00hll5PFN2jfmBQKlvAeU7wfwqdaSnGHOfqf2DeTTaFjIQ4gUhRn/m6pLo\n6kXttLECgYEA9sng0Qz/TcPFfM4tQ1gyvB1cKnnGIwg1FP8sfUjbbEgjaHhA224S\nk3BxtX2Kq6fhTXuwusAFc6OVMAZ76FgrQ5K4Ci7+DTsrF28z4b8td+p+lO/DxgP9\nlTgN+ddsiTOV4fUef9Z3yY0Zr0CnBUMbQYRaV2UIbCdiB0G4V/bt9TsCgYEA0bya\nOo9wGI0RJV0bYP7qwO74Ra1/i1viWbRlS7jU37Q+AZstrlKcQ5CTPzOjKFKMiUzl\n4miWacZ0/q2n+Mvd7NbXGXTLijahnyOYKaHJYyh4oBymfkgAifRstE0Ki9gdvArb\n/I+emC0GvLSyfGN8UUeDJs4NmqdEXGqjo2JOV+MCgYALFv1MR5o9Y1u/hQBRs2fs\nPiGDIx+9OUQxYloccyaxEfjNXAIGGkcpavchIbgWiJ++PJ2vdquIC8TLeK8evL+M\n9M3iX0Q5UfxYvD2HmnCvn9D6Xl/cyRcfGnq+TGjrLW9BzSMGuZt+aiHKV0xqFx7l\nbc4leTvMqGRmURS4lzcQOwKBgQCDzA/i4sYfN25h21tcHXSpnsG3D2rJyQi5NCo/\nZjunA92/JqOTGuiFcLGHEszhhtY3ZXJET1LNz18vtzKJnpqrvOnYXlOVW/U+SqDQ\n8JDb1c/PVZGuY1KrXkR9HLiW3kz5IJ3S3PFdUVYdeTN8BQxXCyg4V12nJJtJs912\ny0zN3wKBgGDS6YttCN6aI4EOABYE8fI1EYQ7vhfiYsaWGWSR1l6bQey7KR6M1ACz\nZzMASNyytVt12yXE4/Emv6/pYqigbDLfL1zQJSLJ3EHJYTh2RxjR+AaGDudYFG/T\nliQ9YXhV5Iu2x1pNwrtFnssDdaaGpfA7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
	}
	resource "alicloud_slb_listener" "https" {
  		load_balancer_id          = "${alicloud_slb_load_balancer.instance.id}"
  		backend_port              = 80
  		frontend_port             = 443
  		protocol                  = "https"
  		sticky_session            = "on"
  		sticky_session_type       = "insert"
  		cookie                    = "testslblistenercookie"
  		cookie_timeout            = 86400
  		health_check              = "on"
  		health_check_uri          = "/cons"
  		health_check_connect_port = 20
  		healthy_threshold         = 8
  		unhealthy_threshold       = 8
  		health_check_timeout      = 8
  		health_check_interval     = 5
  		health_check_http_code    = "http_2xx,http_3xx"
  		bandwidth                 = 10
  		ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
	}
	`, name)
}
