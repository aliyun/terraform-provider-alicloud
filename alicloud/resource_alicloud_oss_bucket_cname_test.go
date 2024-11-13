package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudOssBucketCname_basic8544(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cname.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCnameMap8544)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCname")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcname%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCnameBasicDependence8544)
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
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${alicloud_alidns_record.defaultnHqm5p.domain_name}",
					"certificate": []map[string]interface{}{
						{
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIGBTCCBO2gAwIBAgIQBZ8CvYikSGggdmowncLwVTANBgkqhkiG9w0BAQsFADBu\\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\\nd3cuZGlnaWNlcnQuY29tMS0wKwYDVQQDEyRFbmNyeXB0aW9uIEV2ZXJ5d2hlcmUg\\nRFYgVExTIENBIC0gRzIwHhcNMjQxMTA1MDAwMDAwWhcNMjUwMjAzMjM1OTU5WjAY\\nMRYwFAYDVQQDEw10ZnRlc3RhY2MuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A\\nMIIBCgKCAQEArYwC1TpEMONgxUd6ZdRBmI2G1RzcgUb88bTn//PUWIU7V3kTvzcJ\\ntozkCGhZ3Bl1Kh2srnSqvOTbU8wii1RCRBELCVRAovVBEGa544gQ+UFH92kkRVLs\\nF4lRq+cPjm1fQp3zYzONeLEp8obgoiMYNNgWvB/2Q3/VmSwz0JK1lIaqxDvooFih\\nzQgxgKiYcEYEGDoNW9VcqEIqF3zdWzvBC3eaSH837MSOWiPDvkw3AJBDSsUGBKuc\\n1JQtw+HvXUafBhd+vetbp+5CHM0K8iXm2zJspoBufaDz3CjNVEbi9HKBnyf2b4zb\\nMnlaVX16CPwUIeiQoEKhxdEfNePWlkR7jwIDAQABo4IC8zCCAu8wHwYDVR0jBBgw\\nFoAUeN+RkF/u3qz2xXXr1UxVU+8kSrYwHQYDVR0OBBYEFAU1XHnIPizNK5TKSUQU\\nXJx37hcpMCsGA1UdEQQkMCKCDXRmdGVzdGFjYy5jb22CEXd3dy50ZnRlc3RhY2Mu\\nY29tMD4GA1UdIAQ3MDUwMwYGZ4EMAQIBMCkwJwYIKwYBBQUHAgEWG2h0dHA6Ly93\\nd3cuZGlnaWNlcnQuY29tL0NQUzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\\nKwYBBQUHAwEGCCsGAQUFBwMCMIGABggrBgEFBQcBAQR0MHIwJAYIKwYBBQUHMAGG\\nGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBKBggrBgEFBQcwAoY+aHR0cDovL2Nh\\nY2VydHMuZGlnaWNlcnQuY29tL0VuY3J5cHRpb25FdmVyeXdoZXJlRFZUTFNDQS1H\\nMi5jcnQwDAYDVR0TAQH/BAIwADCCAX4GCisGAQQB1nkCBAIEggFuBIIBagFoAHYA\\nTnWjJ1yaEMM4W2zU3z9S6x3w4I4bjWnAsfpksWKaOd8AAAGS+gk+7wAABAMARzBF\\nAiAviRYG+a8hzAC0fGZGq0cAP+1Tv5Y4XbwZKTKJi+2opAIhAIVtcUcdyLgaJMeL\\n8Bqf9SgVNqCtgU4QNys9dOj+rVpxAHYA5tIxY0B3jMEQQQbXcbnOwdJA9paEhvu6\\nhzId/R43jlAAAAGS+gk/JgAABAMARzBFAiEAwRlbBYiQC2WuKIqwIZQ+nqI81Z97\\nNpkcuXLqhRCTFisCIE7gzDiq17Mnp1H/CQyhNpNB/26E0xt/Bg4Ti1X1oBpPAHYA\\nzxFW7tUufK/zh1vZaS6b6RpxZ0qwF+ysAdJbd87MOwgAAAGS+gk/JwAABAMARzBF\\nAiA6wPovzGoaWdyg5Fh/S1aBDwAJDWoqHG3F4t1hFPYEcwIhALKv5QACKohe8tDr\\nm2Z9GeSoQ/jiqf8jhxXQVz5GUxNkMA0GCSqGSIb3DQEBCwUAA4IBAQDYIG5tyi0s\\njmcGd8dEtViPzAt4gGvZ2RnRffla1r8u2HqFWQb9C1xXdKkjEPfD0M0amvc5FO+n\\nxS4WJco4A6WpB26FgoSVybdrNh8GnZkfcvLdXKoOvNYFkofPYd+tZH1DZfCfipBp\\n2FnV/RVndI2LH16YG4VhoLWwK3NRh6wdwj+qCqWJ2BhRaHdFOjpcYZb44cvh4huW\\nkr56ZwMlFOINfyWfUEtBKpRWduJH40vNwEq3fdKrB4/jC2YkpmnNpkwl3gHigLZ6\\nJHBTDy4JP+k6nnVX90nPYV2grvpQutX2leOy6K6ebGmpp8i3Bevw66PnfJDGABWe\\nmqhAsFrMDJrv\\n-----END CERTIFICATE-----\\n",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEogIBAAKCAQEArYwC1TpEMONgxUd6ZdRBmI2G1RzcgUb88bTn//PUWIU7V3kT\\nvzcJtozkCGhZ3Bl1Kh2srnSqvOTbU8wii1RCRBELCVRAovVBEGa544gQ+UFH92kk\\nRVLsF4lRq+cPjm1fQp3zYzONeLEp8obgoiMYNNgWvB/2Q3/VmSwz0JK1lIaqxDvo\\noFihzQgxgKiYcEYEGDoNW9VcqEIqF3zdWzvBC3eaSH837MSOWiPDvkw3AJBDSsUG\\nBKuc1JQtw+HvXUafBhd+vetbp+5CHM0K8iXm2zJspoBufaDz3CjNVEbi9HKBnyf2\\nb4zbMnlaVX16CPwUIeiQoEKhxdEfNePWlkR7jwIDAQABAoIBAADGzF4eEK5FvsIW\\n/mJErNWMcQslnFH4Z5zRepV7PnWYqP/vPk8HsMmeSUTV8EznPtLbZz3ZDEFV7S1X\\nO2lZzfCfPlAcXiwxNG768tXuipZ0iehsfPoUDvy0/GyNdtsCPqTQIcgLMOMkzxhO\\nqNg0Fiyr9Z141qWbvtiTokaRZVllWnmEVs14xJhPtMowzipYb7+9aW9NTYjT65N9\\niK/aruE2HVQ9yjD6gY09N8Bz6AQSS3jSqtzOIqfoXjKdAFuBq8C00AsdGnBR0qwx\\n6sTseW/rMaIklxDF7KYmOi+JgIaMJTBMFPW8fmYMlpz9LFudVhMMHtP2gqJGJFZL\\nWnsqPFkCgYEA7KjXdISqFvUFUQoK/JzIZ9gbNIv6Kp64S+BxgYo8A3HHKi2Mhnwa\\ngbVL8BrdmXPPOI2k7yfYKZi2RbG/CQzc9WJcsV96c1WBg/D3oHACwlsgYbPR+cFK\\nLzxNNlMHRESMCNeM+uuviksReIOdF3C2KO3T9Z1HS8x9KdYsgBrrnpcCgYEAu7rK\\nChw7NNFz45aV/EkiUX0XxwhN5vn3pyhUIz47tZ68omfz7tZXaifZwqQqxF65Vluu\\n7Kpr9JSGJJwGVT7vDTFr+zZsbh4TOW3AafbgX80gVaBZutNvZF1CxlXYM2W73Xz0\\n37VExz2gC/KMKNNGgOVFg1+5Y3jw2gn1Qfi4ockCgYAcxtcUCwGnsvmHhiIZ33Ka\\n9fMw64hq4Evmpg8HQmjTvmUKYumAfNy4QvRN6OZjP2rGJKsWjZDCVhhr1xY0ooTH\\nrcM5qjN4jMAn7AggUR50xaHlX3k71l4P6lQ1M9lhWrhwZs10wW3h8gjYz6Atdn4f\\n8fNhHVPLCr15ddqJZTybVwKBgAsLoneV0aX57OenJIwDgZFp2sxLIMpGStv683hf\\nYQP+ovqrQx76XYpRbe6V2i5TpHQAUPp7zH5Hft0IkRbS7R3JmqDdQuP3wQnP+1JA\\nxFLertha5uynJBazpgolYuMjSTpu77l54OIYLiKF0tlUFQHge4aPS0kfBIzPqS6I\\ng9SBAoGAOhJIB60Kqu5xIIgAoNi6hCyAB0Ma0snXKUqjNKboyRoeF0KYC8VwHz/h\\nFA2VEYgV+XsLZEM6NdZ4pF0EMLKFt9ZUyD/i1AbZG0twtwNrhSc7R41BRsvOayJE\\nLNAglAQEqlqerZB6N/DzRIf7rpmlt971Ir1O6v3X4atfTj0W09A=\\n-----END RSA PRIVATE KEY-----\\n",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate": []map[string]interface{}{
						{
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIF8TCCBNmgAwIBAgIQDEMtbL1pb2R4Y+6G6rzbuDANBgkqhkiG9w0BAQsFADBu\\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\\nd3cuZGlnaWNlcnQuY29tMS0wKwYDVQQDEyRFbmNyeXB0aW9uIEV2ZXJ5d2hlcmUg\\nRFYgVExTIENBIC0gRzIwHhcNMjQxMDE1MDAwMDAwWhcNMjUwMTEyMjM1OTU5WjAY\\nMRYwFAYDVQQDEw1iai5kaW5hcnkudG9wMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A\\nMIIBCgKCAQEAsDUC/ob21cD3xDekIf8ioL9H4S5X1x+NOQ/+/6YAbo0l0KDbXJJg\\n/+PgTV/ikwJsAqVHzBzd2uIxyYNHxD2XpxRvlOnoXS2gFSBxCI8BPdc1nlepyGB6\\nPFvpuQWfnNrrcuOSkQy7UctVd5ARtDc+OHOj+aADuGHg3ssqbPIvLQtF7shzRwN3\\nbRqZ4nXEqs12tQW3pfi8uMhZdITwPJZ5lHCQB/j+wvBOvJkL0Lpdh2qjdM+iI/0z\\nKPvinMy2rDwifJfTGP93cyl+9iUESfLVDfVPvkClhmeVFbo6mnmfR4DbNUBfIyh3\\ndrX/dyzE0v7KvG6FTPBnDr73wH7GZlSzuwIDAQABo4IC3zCCAtswHwYDVR0jBBgw\\nFoAUeN+RkF/u3qz2xXXr1UxVU+8kSrYwHQYDVR0OBBYEFLimtiPuZwC1es0dRTTg\\n1ihtUmx3MBgGA1UdEQQRMA+CDWJqLmRpbmFyeS50b3AwPgYDVR0gBDcwNTAzBgZn\\ngQwBAgEwKTAnBggrBgEFBQcCARYbaHR0cDovL3d3dy5kaWdpY2VydC5jb20vQ1BT\\nMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw\\ngYAGCCsGAQUFBwEBBHQwcjAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNl\\ncnQuY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20v\\nRW5jcnlwdGlvbkV2ZXJ5d2hlcmVEVlRMU0NBLUcyLmNydDAMBgNVHRMBAf8EAjAA\\nMIIBfQYKKwYBBAHWeQIEAgSCAW0EggFpAWcAdgBOdaMnXJoQwzhbbNTfP1LrHfDg\\njhuNacCx+mSxYpo53wAAAZKOzEEbAAAEAwBHMEUCIHLM0Fq8RCrKyRnSNhC0YZBS\\n4UxjCap1mGydKFj0a40cAiEAg/R3ydG6gIqBTpzv7quAqSyTbiUnqyj/oGW2kvsq\\ndS4AdQDm0jFjQHeMwRBBBtdxuc7B0kD2loSG+7qHMh39HjeOUAAAAZKOzED+AAAE\\nAwBGMEQCIAERiAIyJsQLCL66rIZ9ThhOFYo07f8v5g2EYUE8T0tZAiBLC8FViH2r\\niqrb9Nut9khpywtJ7604XnmCe7QMlT71DQB2AMz7D2qFcQll/pWbU87psnwi6YVc\\nDZeNtql+VMD+TA2wAAABko7MQMUAAAQDAEcwRQIhALnLscXz37q/qBOar6Ws2WWY\\n96uQtd9fpvCITRMs4/U7AiAg6TdZWOSiDNnbRNojFqdsAwKu/ucheTaXO0wbOQfg\\nDzANBgkqhkiG9w0BAQsFAAOCAQEAuLe2m6HknIlhc4hemlVRGmLtcITimFdn89MP\\nZj+t2v15nAmztEIeMsWTp5wFITV1yVMCivk6mr7W5N0fBLkQfxDGwwNjJk93mlQo\\n5uTyvuJQ4OtWbxVYgz9hyJVVFhHpV3+GLse2n4bkxZcMA7MjdrtXEqZDFtJbXFmr\\nEkPNuRBfo/qbKazNZg2EWvjpkWY7OY8QbjhsrijHDpUTc0Ma/+xk1OnUXFjzs04q\\nVxQHPg7xnPaiMKAwxp/B1blLmeMqsTq0cupkg7bjAOEGE9CN7kGRRDt5Z9SO5j7a\\nSrtl9fMxrgGxQ5tubEl+MVcr1IRJM2nIcSSXvPFPoTvlVgqAbg==\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEowIBAAKCAQEAsDUC/ob21cD3xDekIf8ioL9H4S5X1x+NOQ/+/6YAbo0l0KDb\\nXJJg/+PgTV/ikwJsAqVHzBzd2uIxyYNHxD2XpxRvlOnoXS2gFSBxCI8BPdc1nlep\\nyGB6PFvpuQWfnNrrcuOSkQy7UctVd5ARtDc+OHOj+aADuGHg3ssqbPIvLQtF7shz\\nRwN3bRqZ4nXEqs12tQW3pfi8uMhZdITwPJZ5lHCQB/j+wvBOvJkL0Lpdh2qjdM+i\\nI/0zKPvinMy2rDwifJfTGP93cyl+9iUESfLVDfVPvkClhmeVFbo6mnmfR4DbNUBf\\nIyh3drX/dyzE0v7KvG6FTPBnDr73wH7GZlSzuwIDAQABAoIBADwZ4MsTGsce2gO1\\n3MhxvwxoIerHBVQNYXx4ncfyBYyvnRnTe+7PyMEPJzcNAPmWpmOin2IZ6HwbkdLD\\nceuX/I2TFVoMDGMXyFXcamF6cXh32sSG7xS2/4pt6ULgDaiRLSTTRW8vEgdcnOq6\\nm6dF/nV/0Aq5TvuJewtS7cYaNwgcER0LhChtP9Uj1Ui/84TW27ZhwuTcWl09vFJd\\nkyrULzru1E6Wi0VruzcNBL41mLQrxL+FF/zwRpTBAVVsz4ynUMOmTbfl/n2rqIs4\\nZTXuORhuSA7DHtQYpxzyfBNfkWXEd5NQSn7fHtEYdKZtbH2nF3ZF7jB//lmY9riB\\no+NR7V0CgYEA+AR+PkDBzAZhKm+EkeEDo4SEpvCPP1vVReAkwUGV87DqrgHbfnxs\\npHjjmSnVcFdLghq0U6lKksKTrNLdd8PbWusAlZBEyhsZpAHXBr3BW2XqPnnHbuSd\\nBlajzynF3OSkVJ04sUu7gMszOtqsvjOQGJWXkC6MVFqFGZ6h9xsgKrcCgYEAteDY\\njrVaPNGVlDRP9S7IRFAf19iqUCjbfVowAkQR2BBP+s1OvfbCmZM/lSqsKIIQPj+o\\nxvx5A+uqYvDDF20Zz3FOS9GO1rySXR+A63dkDeBrNztDhhX8qMJhBpY36pGQzf2E\\nOEOezsdqxE2vKK/ZMKCS2FhyTHChdkP8hEr2Cx0CgYEAhE1hSrQgrUV577ktbuQp\\nnMDEQolw4MuMKYo4ER97blOh3NEA1ahqDBKw1rOKODNZBD5ak4ZrUX6aaEbT/V9t\\nVEKoPSCIkYeDVgnlOqNe0fK70jgEOxOY8BinqYsPEZamUrzL0Ugk7b93xJ2CKLQ4\\n2eRyxWcPVLA08EW/AKJntmECgYACa4N2IqOYu5Ep76hAsuanQgmqbY+WkXSaLmEF\\nJrK2FUF7LNAnZukf8f2elnrD7zcYHPC59RIHI1OZDWsLHMCDKhbIm3kzEj9ATfMB\\nLw19wcarbXZwikpaVHvGAqmrzVQH6Z+gwAWU6sJY6k+yUuSo6PoLNuIOclEzqaPq\\nfrTXYQKBgG0BlTPZkIFQfPl0VkuY7TEbtt+Yr7M2I9OjkPZjxHCVU7FWIIF8v/4D\\nmieN/8xaKL23hqgAugm8STzAVEJqppfwhET9pt4n8vMlW2be6GqWLMqyFUnPuLSz\\na51yor6QxtSjulPAZ15z15lTskxK0u8XNnVV1b57LSra48MuY9XS\\n-----END RSA PRIVATE KEY-----",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_certificate": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_certificate": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "delete_certificate", "force", "previous_cert_id"},
			},
		},
	})
}

var AlicloudOssBucketCnameMap8544 = map[string]string{}

func AlicloudOssBucketCnameBasicDependence8544(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "tftestacc.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "tftestacc.com"
  priority    = "1"
  value       = alicloud_oss_bucket_cname_token.defaultZaWJfG.token
  ttl         = "600"
  lifecycle {
    ignore_changes = [
      value,
    ]
  }
}


`, name)
}

// Case 自定义域名测试（初始化带CertId） 8542
func TestAccAliCloudOssBucketCname_basic8542(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cname.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCnameMap8542)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCname")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcname%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCnameBasicDependence8542)
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
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${alicloud_alidns_record.defaultnHqm5p.domain_name}",
					"certificate": []map[string]interface{}{
						{
							"cert_id": "15687876-cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_certificate": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_certificate": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "delete_certificate", "force", "previous_cert_id"},
			},
		},
	})
}

var AlicloudOssBucketCnameMap8542 = map[string]string{
	"status": CHECKSET,
}

func AlicloudOssBucketCnameBasicDependence8542(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "tftestacc.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "tftestacc.com"
  priority    = "1"
  value       = alicloud_oss_bucket_cname_token.defaultZaWJfG.token
  ttl         = "600"
  lifecycle {
    ignore_changes = [
      value,
    ]
  }
}


`, name)
}

// Case 自定义域名测试 8386
func TestAccAliCloudOssBucketCname_basic8386(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cname.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCnameMap8386)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCname")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcname%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCnameBasicDependence8386)
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
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${alicloud_alidns_record.defaultnHqm5p.domain_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
					"certificate": []map[string]interface{}{
						{
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIFaTCCA1ECFGaJhRg/DJ8cZEsVlMhSd/fiqAqbMA0GCSqGSIb3DQEBCwUAMHEx\\nCzAJBgNVBAYTAkNOMREwDwYDVQQIDAhaaGVqaWFuZzERMA8GA1UEBwwISGFuZ3po\\nb3UxETAPBgNVBAoMCEFsaUNMb3VkMREwDwYDVQQLDAhzaGFucXVhbjEWMBQGA1UE\\nAwwNYmouZGluYXJ5LnRvcDAeFw0yNDEwMTUwNjAyMTdaFw0yNTEwMTUwNjAyMTda\\nMHExCzAJBgNVBAYTAkNOMREwDwYDVQQIDAhaaGVqaWFuZzERMA8GA1UEBwwISGFu\\nZ3pob3UxETAPBgNVBAoMCEFsaUNMb3VkMREwDwYDVQQLDAhzaGFucXVhbjEWMBQG\\nA1UEAwwNYmouZGluYXJ5LnRvcDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC\\nggIBAKRUVwPmfHhbXbV4VE07M6+2ANSgjYwZ6z2T3X5APsMR0o5daquh6rSsxI7P\\nSWv4a4lYHjM+4uiyQTYXDjnkDgNXMD+oIA7Dgmamw014pmRzECcERkw70LmK6vBS\\n5exyfjfF1Uo6dqVcoaxUWzBhufImPvPuYVGKF0pQfr8QKGmjMi6w5USXSSo1DRV3\\nqmvLQjlm/sVjkeTcvc9FryxQm00gzPfTRbrETOd1Urk76ZGw41CfldlRXiDGt6Xo\\nwWPUyJZpQu2afzuoXj3PW3DDWMNmJ3QyIXq91l7z/HiLshxwAxfajjn1jqrqFgG/\\nE2/61/WOnb1Fg9unheeZcEH4sBstHx5TbIX//r5RFNZtyRpkLyuQZhaESb6sB9A7\\nVAoUCFoYJZQF0tnKj5MjwYhkcwnjB1H6o+QeOb4EWEW3ZYiD2ijS1LW/fdiDH4N4\\n8H5wQoYLt7MtBoNw6G5qR9pk2X0xt3SaPjSRzxgLEMHfAu8UEbspf3l/OVeM7GGn\\n9km3s3k/dKWBOHOcczlhRPFBqffkeDxnW4f7q9TXwUR3gf2e8Cp4IvLgGPJGbrgY\\nYApM6w2VOR6nVThIq7fVnjSbAUZNTjMnfrZe6KZ59s194UreFoPKJfgr1RXvoTyc\\n8LGYpAG1qSG/Zo8IrAquyWL68D/9T0c/HtX8qK4YM8aHOs35AgMBAAEwDQYJKoZI\\nhvcNAQELBQADggIBAImsy4pnNMc9D3jDd1GlzRSwvHX4NRgLxAT1a+yNQDZgylFO\\n35En9HzZLnqJPzl8uG+we4UoxQ3oBA4hwfUZRHrrORTApW7D+YP4h52q5PJALOpI\\njaGnY+zJcXUA081b2eF7MZCMqUVeum9S/UudHXFXn20WgMJ69rcutvYTGmEKAZ4a\\nwyMiAvzWsGl5nkSfxDjOwt01KUQbjt7e/Cnm/snHCg4zokLuJm/G4AEw159jyQt0\\nqJqZGGENfFPFRJvkxMJF7qtIIPYR1D1R2jaJ7d88r7N2MzruZSLpukbcrFXkU4t/\\n3FaZQdAnn/B3JjkljR6UdCn/gd9UrTqGPsaIMMJTZxbtgTtG3kQWJd3KAsQhVkeo\\nyJnnhSk664JuzwdisXweyl7EvBKheACflSHDFX+77su4/+f5Kt+xnNqwVHQTGR1E\\n/inIX2evRUM5XqlOE9oBiz3+KpiMf7RHVl7HCdI5cAX7oOuMcoEyleR86h9BL1kj\\nZ7f6LuaJwSuh+15rnR1ULSRaUUfIq6Csvc6ck4Rme2Mj38FO3njF3M8yJwPQOZS8\\nZKfuCKKdX2isgCLr7NXmpUbctjYvbvVfM+UgDHNW/XB61iRwhi0qYd7leO6bMwnK\\nCjrC5/2JaR/9AUmVWg7RXy+5i4zM9RDqfSvhLAn23qR/lUFeMFQyU7MjKaYT\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN PRIVATE KEY-----\\n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\\nrh0lVGuFTYg5jMnM3UyUCGuDmFw3Zw==\\n-----END PRIVATE KEY-----",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate": []map[string]interface{}{
						{
							"cert_id": "15687876-cn-hangzhou",
						},
					},
					"delete_certificate": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_certificate": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "false",
					"certificate": []map[string]interface{}{
						{
							"cert_id": "15687876-cn-hangzhou",
						},
					},
					"previous_cert_id": "15582268-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force":            "false",
						"previous_cert_id": "15582268-cn-hangzhou",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "delete_certificate", "force", "previous_cert_id"},
			},
		},
	})
}

var AlicloudOssBucketCnameMap8386 = map[string]string{
	"status": CHECKSET,
}

func AlicloudOssBucketCnameBasicDependence8386(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "tftestacc.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "tftestacc.com"
  priority    = "1"
  value       = alicloud_oss_bucket_cname_token.defaultZaWJfG.token
  ttl         = "600"
  lifecycle {
    ignore_changes = [
      value,
    ]
  }
}

`, name)
}
