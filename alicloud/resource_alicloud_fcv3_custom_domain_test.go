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

// Test Fcv3 CustomDomain. >>> Resource test cases, automatically generated.
// Case TestCustomDomain_CORS 12462
func TestAccAliCloudFcv3CustomDomain_basic12462(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3CustomDomainMap12462)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, defaultRegionToTest)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3CustomDomainBasicDependence12462)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{\\n			\\\"jwks\\\":{\\n			  \\\"keys\\\":[\\n				{\\n					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n					\\\"kty\\\": \\\"RSA\\\",\\n					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n					\\\"e\\\": \\\"AQAB\\\",\\n					\\\"use\\\": \\\"sig\\\",\\n					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n					\\\"alg\\\": \\\"RS256\\\",\\n					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n				}\\n			  ]\\n			},\\n			\\\"tokenLookup\\\":\\\"header:auth\\\",\\n			\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICKzCCAZQCCQDsPBz0/KsfiDANBgkqhkiG9w0BAQsFADBaMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMRMwEQYDVQQDDApzZXJ2ZXIuY29tMB4XDTE5MDQwNDAzMzM1MloX\\nDTIwMDQwMzAzMzM1MlowWjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3Rh\\ndGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDETMBEGA1UEAwwK\\nc2VydmVyLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAs6Kp0R11i+cM\\nDyf/Wl0hH+43oooU8GhIRTTjimSypEdXTbAnPA6TZ2QVDTeT43cF8a8hkhRYepo6\\n+mvgYz+lqK7O4xX2zQ/lkG5+pE2iVC+TWySC0UcGXZrArKzA0u0wrBjJZpE2jfhd\\nhDWrlgHfcHTuegec7juSS505JQDCfMsCAwEAATANBgkqhkiG9w0BAQsFAAOBgQAq\\nB/Ia1nZi4N+H9WH9vbYqt1fOLBoKvmNblpd0W2bYTYEbPLCya3Kf4UQQJjWzfJHS\\nEzS7JiQRBAdmAckr5FaKRLwLeV0tJb4pui+VnwQUZgMJBSvXsWq8yPBKA9ZBEvdh\\n+ZU2kJvt+1KLefUycMtdcP6ZU9cBf5TZ9rmwceUNAg==\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXAIBAAKBgQCzoqnRHXWL5wwPJ/9aXSEf7jeiihTwaEhFNOOKZLKkR1dNsCc8\\nDpNnZBUNN5PjdwXxryGSFFh6mjr6a+BjP6Wors7jFfbND+WQbn6kTaJUL5NbJILR\\nRwZdmsCsrMDS7TCsGMlmkTaN+F2ENauWAd9wdO56B5zuO5JLnTklAMJ8ywIDAQAB\\nAoGAPMBCVip0WoAlH+sS/OiKD1ZtEldIhZV++4jLez5a/Bv0dp2gZzs2tryuMe4d\\n4cubAwWLgO/IjI4kbBSXqnkX+MbxqqGCeeLslY2ugHZ5jhI+rEqq0eyspbiTsbwj\\njb7QxC9W52VbokOQchyinmpupl/EniMGFP1FXPxz+pwyokECQQDkLzgrZuJoeYaR\\nhj8BI7aUzCSVgT7EmoDwu5z1YOACZifXhf0IBpxXa3RjU4jcehTq/4LfRgZb6lHG\\nWR+Fdj4TAkEAyYhpCwdaccEK/fULYVqqOETZxXzIneGZzA24ccf4zbOcyVX/6tJA\\npEl0b190WNoBguBYC5mS3wdDO1npHLW9aQJAajxRumM8JcfujvIhgzZNWxlwLurt\\nfjswrOOsP9HKeVN2WTFYjNQHFexBU70giwWLl50+IRVJAKInUGFN+6UBYQJAcfKl\\n6e1zbwQGMgceMyJvQjdzphzi1ZncOqq7UeIOREg86v2sIFpW4E0D/4DKKP7CgfxU\\n6+IeT+osUl+I1YnQmQJBALPkYOx1/Bvs8BnsazQBEnwNI679VW5r6lfiX9GgHBxt\\n6tViPMTigDRaFMK1P4oMzRsX0h4Enk4v9awP0swjdJQ=\\n-----END RSA PRIVATE KEY-----",
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
					"cors_config": []map[string]interface{}{
						{
							"max_age": "3600",
							"allow_origins": []string{
								"*"},
							"allow_methods": []string{
								"*"},
							"allow_headers": []string{
								"Content-Type", "X-Request-ID"},
							"expose_headers": []string{
								"Content-Type", "x-fc-request-id"},
							"allow_credentials": "false",
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
							"auth_info": "{\\n	\\\"jwks\\\":{\\n	  \\\"keys\\\":[\\n		{\\n			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n			\\\"kty\\\": \\\"RSA\\\",\\n			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n			\\\"e\\\": \\\"AQAB\\\",\\n			\\\"use\\\": \\\"sig\\\",\\n			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n			\\\"alg\\\": \\\"RS256\\\",\\n			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n		}\\n	  ]\\n	},\\n	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\",\\n	\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n}",
							"auth_type": "jwt",
						},
					},
					"protocol": "HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name2",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICJzCCAZACCQD6tPfk33WlKzANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMREwDwYDVQQDDAh0ZXN0LmNvbTAeFw0xOTA0MTIwMzM2MzdaFw0y\\nMDA0MTEwMzM2MzdaMFgxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRl\\nMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxETAPBgNVBAMMCHRl\\nc3QuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCp+aW2nAitEyxgCr2\\n0OBq8eOByyrT3sR7uP8yqTxYZcXuATuBZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j\\n6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2SbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ\\n4F7AAcXtvnfiAdXIJxxfUbildwIDAQABMA0GCSqGSIb3DQEBCwUAA4GBAGnAwXKh\\nmadocDIDvz20YDxLi3wPs40I+8QgJwxdpnVVMoxBVsRLGPR82wxz7QrsUPYb5Z0T\\nb2kvU05U2Hxcjzc+beYrQ6eRd8fjXJ9hck7K8brZV5+nAbkpUZd34wXa/nFVoNPX\\neD7NXLSla3i5zucLnPqsnoliroo4AbyWxDrV\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXQIBAAKBgQDCp+aW2nAitEyxgCr20OBq8eOByyrT3sR7uP8yqTxYZcXuATuB\\nZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2\\nSbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ4F7AAcXtvnfiAdXIJxxfUbildwIDAQAB\\nAoGAFT1a1NUK7U59G9UfWwUZp7GzIGN5zdp91/4somuC8SZRvZGW25zYL+o4wvGS\\ndWldbEd3PmDhtvCLT1oVtZeWaDc7f4msoPOO2KaEW6gEd4bOx5nUk4iwiHOV1q5D\\nsV4+JuZH8EBq/ANCPdceIpoiIjKwPD1QGotCXLtJo9GoR1ECQQDte7vCDy3mE5h6\\ni46QISkHK+Qt25LGAGMx+8kkUeui1xTJJ6ie/tx/nsZM4uQqAlRQUDswi8dhVfI0\\nR2P/oRupAkEA0dVQJk4yeHgqvuRz/iNGZBtZK73Zv+nbUFEUPp6yW6p8TU+UNEjo\\nQV1NBxnNOe64WHu23dBUNwCkPGv6mBNsHwJAURrG3tmsRT0//+oVgCezCV32CatJ\\njxGmzvU8lojbvrtRv/kpX1OPHo6tDqkWXzp4bQ1ZiZTTPOzLUQtonW76MQJBAKbi\\ni9NbYAK2N/EIy0P1lDdsFNiYLwXWnancQkineN0005W9U/bdgXLzHJ8oIzQPK6ic\\nBE2YMlJofTbc/jpTQCsCQQCIVIZR9HV3EUKv7yEyDy37UF2BtoAd+6tv5yla6bIw\\nCb+ZX4mwxt1/YWF8Zf8OI6ZW8KG0wMHqq1FM++EOCw3W\\n-----END RSA PRIVATE KEY-----",
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
					"cors_config": []map[string]interface{}{
						{
							"allow_credentials": "true",
							"max_age":           "360",
							"allow_origins": []string{
								"https://api.example.com", "https://api.test.com"},
							"allow_methods": []string{
								"GET", "POST", "PATCH", "DELETE"},
							"allow_headers": []string{
								"X-Custom-Header", "X-Request-ID"},
							"expose_headers": []string{
								"Date", "X-Custom-Header"},
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
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

var AlicloudFcv3CustomDomainMap12462 = map[string]string{
	"api_version":        CHECKSET,
	"account_id":         CHECKSET,
	"create_time":        CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFcv3CustomDomainBasicDependence12462(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name1" {
  default = "terraform-custom-domain-cors"
}

variable "function_name2" {
  default = "terraform-custom-domain-cors"
}


`, name)
}

// Case TestCustomDomain_Waf_Online 7338
func TestAccAliCloudFcv3CustomDomain_basic7338(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3CustomDomainMap7338)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, "cn-shanghai")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3CustomDomainBasicDependence7338)
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{\\n			\\\"jwks\\\":{\\n			  \\\"keys\\\":[\\n				{\\n					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n					\\\"kty\\\": \\\"RSA\\\",\\n					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n					\\\"e\\\": \\\"AQAB\\\",\\n					\\\"use\\\": \\\"sig\\\",\\n					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n					\\\"alg\\\": \\\"RS256\\\",\\n					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n				}\\n			  ]\\n			},\\n			\\\"tokenLookup\\\":\\\"header:auth\\\",\\n			\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"waf_config": []map[string]interface{}{
						{
							"enable_waf": "true",
						},
					},
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICKzCCAZQCCQDsPBz0/KsfiDANBgkqhkiG9w0BAQsFADBaMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMRMwEQYDVQQDDApzZXJ2ZXIuY29tMB4XDTE5MDQwNDAzMzM1MloX\\nDTIwMDQwMzAzMzM1MlowWjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3Rh\\ndGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDETMBEGA1UEAwwK\\nc2VydmVyLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAs6Kp0R11i+cM\\nDyf/Wl0hH+43oooU8GhIRTTjimSypEdXTbAnPA6TZ2QVDTeT43cF8a8hkhRYepo6\\n+mvgYz+lqK7O4xX2zQ/lkG5+pE2iVC+TWySC0UcGXZrArKzA0u0wrBjJZpE2jfhd\\nhDWrlgHfcHTuegec7juSS505JQDCfMsCAwEAATANBgkqhkiG9w0BAQsFAAOBgQAq\\nB/Ia1nZi4N+H9WH9vbYqt1fOLBoKvmNblpd0W2bYTYEbPLCya3Kf4UQQJjWzfJHS\\nEzS7JiQRBAdmAckr5FaKRLwLeV0tJb4pui+VnwQUZgMJBSvXsWq8yPBKA9ZBEvdh\\n+ZU2kJvt+1KLefUycMtdcP6ZU9cBf5TZ9rmwceUNAg==\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXAIBAAKBgQCzoqnRHXWL5wwPJ/9aXSEf7jeiihTwaEhFNOOKZLKkR1dNsCc8\\nDpNnZBUNN5PjdwXxryGSFFh6mjr6a+BjP6Wors7jFfbND+WQbn6kTaJUL5NbJILR\\nRwZdmsCsrMDS7TCsGMlmkTaN+F2ENauWAd9wdO56B5zuO5JLnTklAMJ8ywIDAQAB\\nAoGAPMBCVip0WoAlH+sS/OiKD1ZtEldIhZV++4jLez5a/Bv0dp2gZzs2tryuMe4d\\n4cubAwWLgO/IjI4kbBSXqnkX+MbxqqGCeeLslY2ugHZ5jhI+rEqq0eyspbiTsbwj\\njb7QxC9W52VbokOQchyinmpupl/EniMGFP1FXPxz+pwyokECQQDkLzgrZuJoeYaR\\nhj8BI7aUzCSVgT7EmoDwu5z1YOACZifXhf0IBpxXa3RjU4jcehTq/4LfRgZb6lHG\\nWR+Fdj4TAkEAyYhpCwdaccEK/fULYVqqOETZxXzIneGZzA24ccf4zbOcyVX/6tJA\\npEl0b190WNoBguBYC5mS3wdDO1npHLW9aQJAajxRumM8JcfujvIhgzZNWxlwLurt\\nfjswrOOsP9HKeVN2WTFYjNQHFexBU70giwWLl50+IRVJAKInUGFN+6UBYQJAcfKl\\n6e1zbwQGMgceMyJvQjdzphzi1ZncOqq7UeIOREg86v2sIFpW4E0D/4DKKP7CgfxU\\n6+IeT+osUl+I1YnQmQJBALPkYOx1/Bvs8BnsazQBEnwNI679VW5r6lfiX9GgHBxt\\n6tViPMTigDRaFMK1P4oMzRsX0h4Enk4v9awP0swjdJQ=\\n-----END RSA PRIVATE KEY-----",
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
							"auth_info": "{\\n	\\\"jwks\\\":{\\n	  \\\"keys\\\":[\\n		{\\n			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n			\\\"kty\\\": \\\"RSA\\\",\\n			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n			\\\"e\\\": \\\"AQAB\\\",\\n			\\\"use\\\": \\\"sig\\\",\\n			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n			\\\"alg\\\": \\\"RS256\\\",\\n			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n		}\\n	  ]\\n	},\\n	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\",\\n	\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n}",
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
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICJzCCAZACCQD6tPfk33WlKzANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMREwDwYDVQQDDAh0ZXN0LmNvbTAeFw0xOTA0MTIwMzM2MzdaFw0y\\nMDA0MTEwMzM2MzdaMFgxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRl\\nMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxETAPBgNVBAMMCHRl\\nc3QuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCp+aW2nAitEyxgCr2\\n0OBq8eOByyrT3sR7uP8yqTxYZcXuATuBZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j\\n6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2SbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ\\n4F7AAcXtvnfiAdXIJxxfUbildwIDAQABMA0GCSqGSIb3DQEBCwUAA4GBAGnAwXKh\\nmadocDIDvz20YDxLi3wPs40I+8QgJwxdpnVVMoxBVsRLGPR82wxz7QrsUPYb5Z0T\\nb2kvU05U2Hxcjzc+beYrQ6eRd8fjXJ9hck7K8brZV5+nAbkpUZd34wXa/nFVoNPX\\neD7NXLSla3i5zucLnPqsnoliroo4AbyWxDrV\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXQIBAAKBgQDCp+aW2nAitEyxgCr20OBq8eOByyrT3sR7uP8yqTxYZcXuATuB\\nZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2\\nSbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ4F7AAcXtvnfiAdXIJxxfUbildwIDAQAB\\nAoGAFT1a1NUK7U59G9UfWwUZp7GzIGN5zdp91/4somuC8SZRvZGW25zYL+o4wvGS\\ndWldbEd3PmDhtvCLT1oVtZeWaDc7f4msoPOO2KaEW6gEd4bOx5nUk4iwiHOV1q5D\\nsV4+JuZH8EBq/ANCPdceIpoiIjKwPD1QGotCXLtJo9GoR1ECQQDte7vCDy3mE5h6\\ni46QISkHK+Qt25LGAGMx+8kkUeui1xTJJ6ie/tx/nsZM4uQqAlRQUDswi8dhVfI0\\nR2P/oRupAkEA0dVQJk4yeHgqvuRz/iNGZBtZK73Zv+nbUFEUPp6yW6p8TU+UNEjo\\nQV1NBxnNOe64WHu23dBUNwCkPGv6mBNsHwJAURrG3tmsRT0//+oVgCezCV32CatJ\\njxGmzvU8lojbvrtRv/kpX1OPHo6tDqkWXzp4bQ1ZiZTTPOzLUQtonW76MQJBAKbi\\ni9NbYAK2N/EIy0P1lDdsFNiYLwXWnancQkineN0005W9U/bdgXLzHJ8oIzQPK6ic\\nBE2YMlJofTbc/jpTQCsCQQCIVIZR9HV3EUKv7yEyDy37UF2BtoAd+6tv5yla6bIw\\nCb+ZX4mwxt1/YWF8Zf8OI6ZW8KG0wMHqq1FM++EOCw3W\\n-----END RSA PRIVATE KEY-----",
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
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

var AlicloudFcv3CustomDomainMap7338 = map[string]string{
	"api_version":        CHECKSET,
	"account_id":         CHECKSET,
	"create_time":        CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFcv3CustomDomainBasicDependence7338(name string) string {
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

// Case TestCustomDomain_Waf 6974
func TestAccAliCloudFcv3CustomDomain_basic6974(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3CustomDomainMap6974)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, "cn-shanghai")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3CustomDomainBasicDependence6974)
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{\\n			\\\"jwks\\\":{\\n			  \\\"keys\\\":[\\n				{\\n					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n					\\\"kty\\\": \\\"RSA\\\",\\n					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n					\\\"e\\\": \\\"AQAB\\\",\\n					\\\"use\\\": \\\"sig\\\",\\n					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n					\\\"alg\\\": \\\"RS256\\\",\\n					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n				}\\n			  ]\\n			},\\n			\\\"tokenLookup\\\":\\\"header:auth\\\",\\n			\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"waf_config": []map[string]interface{}{
						{
							"enable_waf": "true",
						},
					},
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICKzCCAZQCCQDsPBz0/KsfiDANBgkqhkiG9w0BAQsFADBaMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMRMwEQYDVQQDDApzZXJ2ZXIuY29tMB4XDTE5MDQwNDAzMzM1MloX\\nDTIwMDQwMzAzMzM1MlowWjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3Rh\\ndGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDETMBEGA1UEAwwK\\nc2VydmVyLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAs6Kp0R11i+cM\\nDyf/Wl0hH+43oooU8GhIRTTjimSypEdXTbAnPA6TZ2QVDTeT43cF8a8hkhRYepo6\\n+mvgYz+lqK7O4xX2zQ/lkG5+pE2iVC+TWySC0UcGXZrArKzA0u0wrBjJZpE2jfhd\\nhDWrlgHfcHTuegec7juSS505JQDCfMsCAwEAATANBgkqhkiG9w0BAQsFAAOBgQAq\\nB/Ia1nZi4N+H9WH9vbYqt1fOLBoKvmNblpd0W2bYTYEbPLCya3Kf4UQQJjWzfJHS\\nEzS7JiQRBAdmAckr5FaKRLwLeV0tJb4pui+VnwQUZgMJBSvXsWq8yPBKA9ZBEvdh\\n+ZU2kJvt+1KLefUycMtdcP6ZU9cBf5TZ9rmwceUNAg==\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXAIBAAKBgQCzoqnRHXWL5wwPJ/9aXSEf7jeiihTwaEhFNOOKZLKkR1dNsCc8\\nDpNnZBUNN5PjdwXxryGSFFh6mjr6a+BjP6Wors7jFfbND+WQbn6kTaJUL5NbJILR\\nRwZdmsCsrMDS7TCsGMlmkTaN+F2ENauWAd9wdO56B5zuO5JLnTklAMJ8ywIDAQAB\\nAoGAPMBCVip0WoAlH+sS/OiKD1ZtEldIhZV++4jLez5a/Bv0dp2gZzs2tryuMe4d\\n4cubAwWLgO/IjI4kbBSXqnkX+MbxqqGCeeLslY2ugHZ5jhI+rEqq0eyspbiTsbwj\\njb7QxC9W52VbokOQchyinmpupl/EniMGFP1FXPxz+pwyokECQQDkLzgrZuJoeYaR\\nhj8BI7aUzCSVgT7EmoDwu5z1YOACZifXhf0IBpxXa3RjU4jcehTq/4LfRgZb6lHG\\nWR+Fdj4TAkEAyYhpCwdaccEK/fULYVqqOETZxXzIneGZzA24ccf4zbOcyVX/6tJA\\npEl0b190WNoBguBYC5mS3wdDO1npHLW9aQJAajxRumM8JcfujvIhgzZNWxlwLurt\\nfjswrOOsP9HKeVN2WTFYjNQHFexBU70giwWLl50+IRVJAKInUGFN+6UBYQJAcfKl\\n6e1zbwQGMgceMyJvQjdzphzi1ZncOqq7UeIOREg86v2sIFpW4E0D/4DKKP7CgfxU\\n6+IeT+osUl+I1YnQmQJBALPkYOx1/Bvs8BnsazQBEnwNI679VW5r6lfiX9GgHBxt\\n6tViPMTigDRaFMK1P4oMzRsX0h4Enk4v9awP0swjdJQ=\\n-----END RSA PRIVATE KEY-----",
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
							"auth_info": "{\\n	\\\"jwks\\\":{\\n	  \\\"keys\\\":[\\n		{\\n			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n			\\\"kty\\\": \\\"RSA\\\",\\n			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n			\\\"e\\\": \\\"AQAB\\\",\\n			\\\"use\\\": \\\"sig\\\",\\n			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n			\\\"alg\\\": \\\"RS256\\\",\\n			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n		}\\n	  ]\\n	},\\n	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\",\\n	\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n}",
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
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICJzCCAZACCQD6tPfk33WlKzANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMREwDwYDVQQDDAh0ZXN0LmNvbTAeFw0xOTA0MTIwMzM2MzdaFw0y\\nMDA0MTEwMzM2MzdaMFgxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRl\\nMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxETAPBgNVBAMMCHRl\\nc3QuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCp+aW2nAitEyxgCr2\\n0OBq8eOByyrT3sR7uP8yqTxYZcXuATuBZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j\\n6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2SbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ\\n4F7AAcXtvnfiAdXIJxxfUbildwIDAQABMA0GCSqGSIb3DQEBCwUAA4GBAGnAwXKh\\nmadocDIDvz20YDxLi3wPs40I+8QgJwxdpnVVMoxBVsRLGPR82wxz7QrsUPYb5Z0T\\nb2kvU05U2Hxcjzc+beYrQ6eRd8fjXJ9hck7K8brZV5+nAbkpUZd34wXa/nFVoNPX\\neD7NXLSla3i5zucLnPqsnoliroo4AbyWxDrV\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXQIBAAKBgQDCp+aW2nAitEyxgCr20OBq8eOByyrT3sR7uP8yqTxYZcXuATuB\\nZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2\\nSbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ4F7AAcXtvnfiAdXIJxxfUbildwIDAQAB\\nAoGAFT1a1NUK7U59G9UfWwUZp7GzIGN5zdp91/4somuC8SZRvZGW25zYL+o4wvGS\\ndWldbEd3PmDhtvCLT1oVtZeWaDc7f4msoPOO2KaEW6gEd4bOx5nUk4iwiHOV1q5D\\nsV4+JuZH8EBq/ANCPdceIpoiIjKwPD1QGotCXLtJo9GoR1ECQQDte7vCDy3mE5h6\\ni46QISkHK+Qt25LGAGMx+8kkUeui1xTJJ6ie/tx/nsZM4uQqAlRQUDswi8dhVfI0\\nR2P/oRupAkEA0dVQJk4yeHgqvuRz/iNGZBtZK73Zv+nbUFEUPp6yW6p8TU+UNEjo\\nQV1NBxnNOe64WHu23dBUNwCkPGv6mBNsHwJAURrG3tmsRT0//+oVgCezCV32CatJ\\njxGmzvU8lojbvrtRv/kpX1OPHo6tDqkWXzp4bQ1ZiZTTPOzLUQtonW76MQJBAKbi\\ni9NbYAK2N/EIy0P1lDdsFNiYLwXWnancQkineN0005W9U/bdgXLzHJ8oIzQPK6ic\\nBE2YMlJofTbc/jpTQCsCQQCIVIZR9HV3EUKv7yEyDy37UF2BtoAd+6tv5yla6bIw\\nCb+ZX4mwxt1/YWF8Zf8OI6ZW8KG0wMHqq1FM++EOCw3W\\n-----END RSA PRIVATE KEY-----",
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
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

var AlicloudFcv3CustomDomainMap6974 = map[string]string{
	"api_version":        CHECKSET,
	"account_id":         CHECKSET,
	"create_time":        CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFcv3CustomDomainBasicDependence6974(name string) string {
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

// Case TestCustomDomain_Base 7241
func TestAccAliCloudFcv3CustomDomain_basic7241(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_custom_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3CustomDomainMap7241)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3CustomDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("flask-07ap.fcv3.%d.%s.fc.devsapp.net", 1511928242963727, "cn-shanghai")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3CustomDomainBasicDependence7241)
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
								{
									"function_name": "${var.function_name1}",
									"methods": []string{
										"POST"},
									"path":      "/c",
									"qualifier": "1",
									"rewrite_config": []map[string]interface{}{
										{},
									},
								},
							},
						},
					},
					"auth_config": []map[string]interface{}{
						{
							"auth_type": "jwt",
							"auth_info": "{\\n			\\\"jwks\\\":{\\n			  \\\"keys\\\":[\\n				{\\n					\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n					\\\"kty\\\": \\\"RSA\\\",\\n					\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n					\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n					\\\"e\\\": \\\"AQAB\\\",\\n					\\\"use\\\": \\\"sig\\\",\\n					\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n					\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n					\\\"alg\\\": \\\"RS256\\\",\\n					\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n					\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n				}\\n			  ]\\n			},\\n			\\\"tokenLookup\\\":\\\"header:auth\\\",\\n			\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n		}",
						},
					},
					"protocol": "HTTP,HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICKzCCAZQCCQDsPBz0/KsfiDANBgkqhkiG9w0BAQsFADBaMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMRMwEQYDVQQDDApzZXJ2ZXIuY29tMB4XDTE5MDQwNDAzMzM1MloX\\nDTIwMDQwMzAzMzM1MlowWjELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3Rh\\ndGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDETMBEGA1UEAwwK\\nc2VydmVyLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAs6Kp0R11i+cM\\nDyf/Wl0hH+43oooU8GhIRTTjimSypEdXTbAnPA6TZ2QVDTeT43cF8a8hkhRYepo6\\n+mvgYz+lqK7O4xX2zQ/lkG5+pE2iVC+TWySC0UcGXZrArKzA0u0wrBjJZpE2jfhd\\nhDWrlgHfcHTuegec7juSS505JQDCfMsCAwEAATANBgkqhkiG9w0BAQsFAAOBgQAq\\nB/Ia1nZi4N+H9WH9vbYqt1fOLBoKvmNblpd0W2bYTYEbPLCya3Kf4UQQJjWzfJHS\\nEzS7JiQRBAdmAckr5FaKRLwLeV0tJb4pui+VnwQUZgMJBSvXsWq8yPBKA9ZBEvdh\\n+ZU2kJvt+1KLefUycMtdcP6ZU9cBf5TZ9rmwceUNAg==\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXAIBAAKBgQCzoqnRHXWL5wwPJ/9aXSEf7jeiihTwaEhFNOOKZLKkR1dNsCc8\\nDpNnZBUNN5PjdwXxryGSFFh6mjr6a+BjP6Wors7jFfbND+WQbn6kTaJUL5NbJILR\\nRwZdmsCsrMDS7TCsGMlmkTaN+F2ENauWAd9wdO56B5zuO5JLnTklAMJ8ywIDAQAB\\nAoGAPMBCVip0WoAlH+sS/OiKD1ZtEldIhZV++4jLez5a/Bv0dp2gZzs2tryuMe4d\\n4cubAwWLgO/IjI4kbBSXqnkX+MbxqqGCeeLslY2ugHZ5jhI+rEqq0eyspbiTsbwj\\njb7QxC9W52VbokOQchyinmpupl/EniMGFP1FXPxz+pwyokECQQDkLzgrZuJoeYaR\\nhj8BI7aUzCSVgT7EmoDwu5z1YOACZifXhf0IBpxXa3RjU4jcehTq/4LfRgZb6lHG\\nWR+Fdj4TAkEAyYhpCwdaccEK/fULYVqqOETZxXzIneGZzA24ccf4zbOcyVX/6tJA\\npEl0b190WNoBguBYC5mS3wdDO1npHLW9aQJAajxRumM8JcfujvIhgzZNWxlwLurt\\nfjswrOOsP9HKeVN2WTFYjNQHFexBU70giwWLl50+IRVJAKInUGFN+6UBYQJAcfKl\\n6e1zbwQGMgceMyJvQjdzphzi1ZncOqq7UeIOREg86v2sIFpW4E0D/4DKKP7CgfxU\\n6+IeT+osUl+I1YnQmQJBALPkYOx1/Bvs8BnsazQBEnwNI679VW5r6lfiX9GgHBxt\\n6tViPMTigDRaFMK1P4oMzRsX0h4Enk4v9awP0swjdJQ=\\n-----END RSA PRIVATE KEY-----",
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
							"auth_info": "{\\n	\\\"jwks\\\":{\\n	  \\\"keys\\\":[\\n		{\\n			\\\"p\\\": \\\"8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM\\\",\\n			\\\"kty\\\": \\\"RSA\\\",\\n			\\\"q\\\": \\\"x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us\\\",\\n			\\\"d\\\": \\\"bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ\\\",\\n			\\\"e\\\": \\\"AQAB\\\",\\n			\\\"use\\\": \\\"sig\\\",\\n			\\\"qi\\\": \\\"MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM\\\",\\n			\\\"dp\\\": \\\"J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU\\\",\\n			\\\"alg\\\": \\\"RS256\\\",\\n			\\\"dq\\\": \\\"H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas\\\",\\n			\\\"n\\\": \\\"u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ\\\"\\n		}\\n	  ]\\n	},\\n	\\\"tokenLookup\\\":\\\"header:auth:noneSensePrefix\\\",\\n	\\\"claimPassBy\\\":\\\"header:name:name\\\"\\n}",
							"auth_type": "jwt",
						},
					},
					"protocol": "HTTPS",
					"cert_config": []map[string]interface{}{
						{
							"cert_name":   "cert-name2",
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIICJzCCAZACCQD6tPfk33WlKzANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJB\\nVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0\\ncyBQdHkgTHRkMREwDwYDVQQDDAh0ZXN0LmNvbTAeFw0xOTA0MTIwMzM2MzdaFw0y\\nMDA0MTEwMzM2MzdaMFgxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRl\\nMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxETAPBgNVBAMMCHRl\\nc3QuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCp+aW2nAitEyxgCr2\\n0OBq8eOByyrT3sR7uP8yqTxYZcXuATuBZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j\\n6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2SbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ\\n4F7AAcXtvnfiAdXIJxxfUbildwIDAQABMA0GCSqGSIb3DQEBCwUAA4GBAGnAwXKh\\nmadocDIDvz20YDxLi3wPs40I+8QgJwxdpnVVMoxBVsRLGPR82wxz7QrsUPYb5Z0T\\nb2kvU05U2Hxcjzc+beYrQ6eRd8fjXJ9hck7K8brZV5+nAbkpUZd34wXa/nFVoNPX\\neD7NXLSla3i5zucLnPqsnoliroo4AbyWxDrV\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIICXQIBAAKBgQDCp+aW2nAitEyxgCr20OBq8eOByyrT3sR7uP8yqTxYZcXuATuB\\nZ0YDvN1shO/RdGm88bL0XqgI9YUP1l4j6s6PpHT7DWxBaOJ5xwYM83Cz11TJT5T2\\nSbQDcc6leMpxK9i1AoUTPUpKONBL8IAQ4F7AAcXtvnfiAdXIJxxfUbildwIDAQAB\\nAoGAFT1a1NUK7U59G9UfWwUZp7GzIGN5zdp91/4somuC8SZRvZGW25zYL+o4wvGS\\ndWldbEd3PmDhtvCLT1oVtZeWaDc7f4msoPOO2KaEW6gEd4bOx5nUk4iwiHOV1q5D\\nsV4+JuZH8EBq/ANCPdceIpoiIjKwPD1QGotCXLtJo9GoR1ECQQDte7vCDy3mE5h6\\ni46QISkHK+Qt25LGAGMx+8kkUeui1xTJJ6ie/tx/nsZM4uQqAlRQUDswi8dhVfI0\\nR2P/oRupAkEA0dVQJk4yeHgqvuRz/iNGZBtZK73Zv+nbUFEUPp6yW6p8TU+UNEjo\\nQV1NBxnNOe64WHu23dBUNwCkPGv6mBNsHwJAURrG3tmsRT0//+oVgCezCV32CatJ\\njxGmzvU8lojbvrtRv/kpX1OPHo6tDqkWXzp4bQ1ZiZTTPOzLUQtonW76MQJBAKbi\\ni9NbYAK2N/EIy0P1lDdsFNiYLwXWnancQkineN0005W9U/bdgXLzHJ8oIzQPK6ic\\nBE2YMlJofTbc/jpTQCsCQQCIVIZR9HV3EUKv7yEyDy37UF2BtoAd+6tv5yla6bIw\\nCb+ZX4mwxt1/YWF8Zf8OI6ZW8KG0wMHqq1FM++EOCw3W\\n-----END RSA PRIVATE KEY-----",
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
									"rewrite_config": []map[string]interface{}{
										{},
									},
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

var AlicloudFcv3CustomDomainMap7241 = map[string]string{
	"api_version":        CHECKSET,
	"account_id":         CHECKSET,
	"create_time":        CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFcv3CustomDomainBasicDependence7241(name string) string {
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

// Test Fcv3 CustomDomain. <<< Resource test cases, automatically generated.
