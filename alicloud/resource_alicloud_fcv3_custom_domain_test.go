package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testFc3ServerCertificate = `-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----`
const testFc3PrivateKey = `-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----`

// Test Fc3 CustomDomain. >>> Resource test cases, automatically generated.

var AlicloudFc3CustomDomainMap6974 = map[string]string{
	"custom_domain_name": CHECKSET,
	"create_time":        CHECKSET,
}

func AlicloudFc3CustomDomainBasicDependence6974(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name1" {
  default = "terraform-custom-domain-t1"
}

variable "function_name2" {
  default = "terraform-custom-domain-t1"
}


`, name)
}

var AlicloudFc3CustomDomainMap7241 = map[string]string{
	"custom_domain_name": CHECKSET,
	"create_time":        CHECKSET,
}

func AlicloudFc3CustomDomainBasicDependence7241(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name1" {
  default = "terraform-custom-domain-t1"
}

variable "function_name2" {
  default = "terraform-custom-domain-t1"
}


`, name)
}

// Case TestCustomDomain_Waf 6974  raw
func TestAccAliCloudFcv3CustomDomain_basic6974_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3CustomDomainMap6974)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, defaultRegionToTest)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3CustomDomainBasicDependence6974)
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
					"custom_domain_name": name,
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									"function_name": "${var.function_name1}",
									"rewrite_config": []map[string]interface{}{
										{
											"equal_rules": []map[string]interface{}{
												{
													"match":       "/old",
													"replacement": "/new",
												},
												{
													"match":       "/old1",
													"replacement": "/new1",
												},
												{
													"match":       "/old2",
													"replacement": "/new2",
												},
											},
											"regex_rules": []map[string]interface{}{
												{
													"match":       "/api/*",
													"replacement": "$1",
												},
												{
													"match":       "/api1/*",
													"replacement": "$1",
												},
												{
													"match":       "/api2/*",
													"replacement": "$1",
												},
											},
											"wildcard_rules": []map[string]interface{}{
												{
													"match":       "^/api1/.+?/(.*)",
													"replacement": "/api/v1/$1",
												},
												{
													"match":       "^/api2/.+?/(.*)",
													"replacement": "/api/v2/$1",
												},
												{
													"match":       "^/api2/.+?/(.*)",
													"replacement": "/api/v3/$1",
												},
											},
										},
									},
									"methods": []string{
										"GET", "POST", "DELETE", "HEAD"},
									"path":      "/a",
									"qualifier": "LATEST",
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"GET"},
									"path":      "/b",
									"qualifier": "LATEST",
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{ 			\\\"jwks\\\":{ 			  \\\"keys\\\":[ 				{ 					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\", 					\\\"kty\\\": \\\"RSA\\\", 					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\", 					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\", 					\\\"e\\\": \\\"AQAB\\\", 					\\\"use\\\": \\\"sig\\\", 					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\", 					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\", 					\\\"alg\\\": \\\"RS256\\\", 					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\", 					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\" 				} 			  ] 			}, 			\\\"tokenLookup\\\":\\\"header:auth\\\", 			\\\"claimPassBy\\\":\\\"header:name:name\\\" 		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"waf_config": []map[string]interface{}{
						{
							"enable_waf": "false",
						},
					},
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": testFcCertificate,
							"private_key": testFcPrivateKey,
						},
					},
					"tls_config": []map[string]interface{}{
						{
							"cipher_suites": []string{
								"TLS_RSA_WITH_AES_128_CBC_SHA", "TLS_RSA_WITH_AES_256_CBC_SHA", "TLS_RSA_WITH_AES_128_GCM_SHA256", "TLS_RSA_WITH_AES_256_GCM_SHA384"},
							"max_version": "TLSv1.3",
							"min_version": "TLSv1.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_domain_name": name,
						"protocol":           "HTTP,HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									"function_name": "${var.function_name2}",
									"rewrite_config": []map[string]interface{}{
										{
											"equal_rules": []map[string]interface{}{
												{
													"match":       "/old3",
													"replacement": "/new3",
												},
											},
											"regex_rules": []map[string]interface{}{
												{
													"match":       "/api3/*",
													"replacement": "/api3/$1",
												},
											},
											"wildcard_rules": []map[string]interface{}{
												{
													"match":       "^/api4/.+?/(.*)",
													"replacement": "/api/v4/$1",
												},
											},
										},
									},
									"methods": []string{
										"GET", "POST"},
									"path":      "/b",
									"qualifier": "1",
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_info": "{ 	\\\"jwks\\\":{ 	  \\\"keys\\\":[ 		{ 			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\", 			\\\"kty\\\": \\\"RSA\\\", 			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\", 			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\", 			\\\"e\\\": \\\"AQAB\\\", 			\\\"use\\\": \\\"sig\\\", 			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\", 			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\", 			\\\"alg\\\": \\\"RS256\\\", 			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\", 			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\" 		} 	  ] 	}, 	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\", 	\\\"claimPassBy\\\":\\\"header:name:name\\\" }",
							"auth_type": "jwt",
						},
					},
					"protocol": "HTTPS",
					"waf_config": []map[string]interface{}{
						{
							"enable_waf": "false",
						},
					},
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name2",
							"certificate": testFcCertificate,
							"private_key": testFcPrivateKey,
						},
					},
					"tls_config": []map[string]interface{}{
						{
							"cipher_suites": []string{
								"TLS_RSA_WITH_RC4_128_SHA", "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"},
							"max_version": "TLSv1.2",
							"min_version": "TLSv1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
									"function_name": "${var.function_name2}",
									"methods": []string{
										"GET"},
									"path":      "/*",
									"qualifier": "LATEST",
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "function",
						},
					},
					"protocol":    "HTTP",
					"cert_config": REMOVEKEY,
					"tls_config":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_config"},
			},
		},
	})
}

// Case TestCustomDomain_Base 7241  raw
func TestAccAliCloudFcv3CustomDomain_basic7241_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFc3CustomDomainMap7241)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, defaultRegionToTest)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFc3CustomDomainBasicDependence7241)
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
					"custom_domain_name": name,
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									"function_name": "${var.function_name1}",
									"rewrite_config": []map[string]interface{}{
										{
											"equal_rules": []map[string]interface{}{
												{
													"match":       "/old",
													"replacement": "/new",
												},
												{
													"match":       "/old1",
													"replacement": "/new1",
												},
												{
													"match":       "/old2",
													"replacement": "/new2",
												},
											},
											"regex_rules": []map[string]interface{}{
												{
													"match":       "/api/*",
													"replacement": "$1",
												},
												{
													"match":       "/api1/*",
													"replacement": "$1",
												},
												{
													"match":       "/api2/*",
													"replacement": "$1",
												},
											},
											"wildcard_rules": []map[string]interface{}{
												{
													"match":       "^/api1/.+?/(.*)",
													"replacement": "/api/v1/$1",
												},
												{
													"match":       "^/api2/.+?/(.*)",
													"replacement": "/api/v2/$1",
												},
												{
													"match":       "^/api2/.+?/(.*)",
													"replacement": "/api/v3/$1",
												},
											},
										},
									},
									"methods": []string{
										"GET", "POST", "DELETE", "HEAD"},
									"path":      "/a",
									"qualifier": "LATEST",
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"GET"},
									"path":      "/b",
									"qualifier": "LATEST",
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{ 			\\\"jwks\\\":{ 			  \\\"keys\\\":[ 				{ 					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\", 					\\\"kty\\\": \\\"RSA\\\", 					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\", 					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\", 					\\\"e\\\": \\\"AQAB\\\", 					\\\"use\\\": \\\"sig\\\", 					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\", 					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\", 					\\\"alg\\\": \\\"RS256\\\", 					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\", 					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\" 				} 			  ] 			}, 			\\\"tokenLookup\\\":\\\"header:auth\\\", 			\\\"claimPassBy\\\":\\\"header:name:name\\\" 		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": testFcCertificate,
							"private_key": testFcPrivateKey,
						},
					},
					"tls_config": []map[string]interface{}{
						{
							"cipher_suites": []string{
								"TLS_RSA_WITH_AES_128_CBC_SHA", "TLS_RSA_WITH_AES_256_CBC_SHA", "TLS_RSA_WITH_AES_128_GCM_SHA256", "TLS_RSA_WITH_AES_256_GCM_SHA384"},
							"max_version": "TLSv1.3",
							"min_version": "TLSv1.0",
						},
					},
					"waf_config": []map[string]interface{}{
						{
							"enable_waf": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_domain_name": name,
						"protocol":           "HTTP,HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									"function_name": "${var.function_name2}",
									"rewrite_config": []map[string]interface{}{
										{
											"equal_rules": []map[string]interface{}{
												{
													"match":       "/old3",
													"replacement": "/new3",
												},
											},
											"regex_rules": []map[string]interface{}{
												{
													"match":       "/api3/*",
													"replacement": "/api3/$1",
												},
											},
											"wildcard_rules": []map[string]interface{}{
												{
													"match":       "^/api4/.+?/(.*)",
													"replacement": "/api/v4/$1",
												},
											},
										},
									},
									"methods": []string{
										"GET", "POST"},
									"path":      "/b",
									"qualifier": "1",
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_info": "{ 	\\\"jwks\\\":{ 	  \\\"keys\\\":[ 		{ 			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\", 			\\\"kty\\\": \\\"RSA\\\", 			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\", 			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\", 			\\\"e\\\": \\\"AQAB\\\", 			\\\"use\\\": \\\"sig\\\", 			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\", 			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\", 			\\\"alg\\\": \\\"RS256\\\", 			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\", 			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\" 		} 	  ] 	}, 	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\", 	\\\"claimPassBy\\\":\\\"header:name:name\\\" }",
							"auth_type": "jwt",
						},
					},
					"protocol": "HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name2",
							"certificate": testFc3ServerCertificate,
							"private_key": testFc3PrivateKey,
						},
					},
					"tls_config": []map[string]interface{}{
						{
							"cipher_suites": []string{
								"TLS_RSA_WITH_RC4_128_SHA", "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"},
							"max_version": "TLSv1.2",
							"min_version": "TLSv1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_config": []map[string]interface{}{
						{
							"routes": []map[string]interface{}{
								{
									//"rewrite_config": []map[string]interface{}{
									//	{
									//		"equal_rules":    []map[string]interface{}{},
									//		"regex_rules":    []map[string]interface{}{},
									//		"wildcard_rules": []map[string]interface{}{},
									//	},
									//},
									"function_name": "${var.function_name2}",
									"methods": []string{
										"GET"},
									"path":      "/*",
									"qualifier": "LATEST",
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "function",
						},
					},
					"cert_config": REMOVEKEY,
					"tls_config":  REMOVEKEY,
					"protocol":    "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_config"},
			},
		},
	})
}

// Test Fc3 CustomDomain. <<< Resource test cases, automatically generated.
