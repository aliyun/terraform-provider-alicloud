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
  domain = "songwenpeng.alivetest.asia"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "songwenpeng.alivetest.asia"
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
							"cert_id": "22495571-cn-hangzhou",
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
  domain = "songwenpeng.alivetest.asia"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth"
  type        = "TXT"
  domain_name = "songwenpeng.alivetest.asia"
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
					"domain": "songwenpeng.${alicloud_alidns_record.defaultnHqm5p.domain_name}",
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
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIGLDCCBRSgAwIBAgIQC1haGkCG29WLl7YjhfFqfDANBgkqhkiG9w0BAQsFADBu\\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\\nd3cuZGlnaWNlcnQuY29tMS0wKwYDVQQDEyRFbmNyeXB0aW9uIEV2ZXJ5d2hlcmUg\\nRFYgVExTIENBIC0gRzIwHhcNMjUxMjI5MDAwMDAwWhcNMjYwMzI4MjM1OTU5WjAl\\nMSMwIQYDVQQDExpzb25nd2VucGVuZy5hbGl2ZXRlc3QuYXNpYTCCASIwDQYJKoZI\\nhvcNAQEBBQADggEPADCCAQoCggEBAM34HhxQW2URornO1y5EAG7wBIW/WxFzP81g\\nASejB5qHzQxeDSwu4QlSW7qFU5njf3toQ4HG5f9XxfNSpmirlpiPSlWLOaws7kIa\\nlmc57pLVsLCgA4OIp2SjgUXNlfPjdQ/uq/Z07BGOZ7EUaZo7qOd2MgGkwNlb4hZW\\ntK2LSTJp9lH2hlgs50LU6UOC8qTya9xpyhrzorZKpGX9Oe5F1OjByDsj21cIeeu3\\nzMXzB8XHX9vLvPf1tAUorlezRX4T2BlslNbHYpdbUpeIyRD5yoOb2I4zPw62nXV9\\nirLt33XhKpM+eMGuldG8wI9m5+WVr2vycswQ2zUjDE2uRCXDYPkCAwEAAaOCAw0w\\nggMJMB8GA1UdIwQYMBaAFHjfkZBf7t6s9sV169VMVVPvJEq2MB0GA1UdDgQWBBQt\\nONdtzvXRX4IKV94S99BzLBk4VzBFBgNVHREEPjA8ghpzb25nd2VucGVuZy5hbGl2\\nZXRlc3QuYXNpYYIed3d3LnNvbmd3ZW5wZW5nLmFsaXZldGVzdC5hc2lhMD4GA1Ud\\nIAQ3MDUwMwYGZ4EMAQIBMCkwJwYIKwYBBQUHAgEWG2h0dHA6Ly93d3cuZGlnaWNl\\ncnQuY29tL0NQUzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEG\\nCCsGAQUFBwMCMIGABggrBgEFBQcBAQR0MHIwJAYIKwYBBQUHMAGGGGh0dHA6Ly9v\\nY3NwLmRpZ2ljZXJ0LmNvbTBKBggrBgEFBQcwAoY+aHR0cDovL2NhY2VydHMuZGln\\naWNlcnQuY29tL0VuY3J5cHRpb25FdmVyeXdoZXJlRFZUTFNDQS1HMi5jcnQwDAYD\\nVR0TAQH/BAIwADCCAX4GCisGAQQB1nkCBAIEggFuBIIBagFoAHYAlpdkv1VYl633\\nQ4doNwhCd+nwOtX2pPM2bkakPw/KqcYAAAGbae/bnwAABAMARzBFAiBuc7rbM4gK\\ny87P7A5I0B5WQIvXOfgWhwG7u9ygCuRw2AIhAMGOgxpAYUF4rHc4HvZTDpg06iGq\\nqYLs9fqS9qNizksbAHYAFoMtq/CpJQ8P8DqlRf/Iv8gj0IdL9gQpJ/jnHzMT9foA\\nAAGbae/blwAABAMARzBFAiBzu8OJcdPBhEHdUGWPDllNX6AqPBqj1FQkKHohA2mp\\nBwIhAL1t0T+dSN4ZPBi+CmoWL7Uskcds1wBB8O1IrA+lGjT1AHYAZBHEbKQS7KeJ\\nHKICLgC8q08oB9QeNSer6v7VA8l9zfAAAAGbae/bmwAABAMARzBFAiAW/MJQij/I\\npiwZ3SDBr2/TNtHnNu26iCHLMvGRYmaWCAIhAONH4obDKa+Z7wd1eBZGK0r/Mm7i\\neUHKkAb4s3ETy5uZMA0GCSqGSIb3DQEBCwUAA4IBAQCGpNhlNt4NmMRXikaTpD32\\nlGpUV3EZ+XXhmSIt7p362UfB/T/GDozfX2aH009PJ4IIYHibNvFnZFXDn+pVPn9a\\nYm941Vu2Khzt4GGXx76oA5AML8ZOl7GBMBukPZMg53fCr0BLDqvH0BOyenfBPCYt\\nzpk5kdzWT/YufHpoDBRkaz4qE6mkcEt+wggzJWJLyhukFAVVLQPhj54OWX3dRe5W\\nQIe3TSzZWSuIto9+PGd+s93oN0OXi7PUkOQSoOrYMunwgopxcICi8mewAPZEpLxC\\nhiKCoQ7vBVFTQ6t0J+KGFOd9XpwuQg2BR4LGtnmOwnBSCnoOdZiG+dqomku3ztPQ\\n-----END CERTIFICATE-----",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEogIBAAKCAQEAzfgeHFBbZRGiuc7XLkQAbvAEhb9bEXM/zWABJ6MHmofNDF4N\\nLC7hCVJbuoVTmeN/e2hDgcbl/1fF81KmaKuWmI9KVYs5rCzuQhqWZznuktWwsKAD\\ng4inZKOBRc2V8+N1D+6r9nTsEY5nsRRpmjuo53YyAaTA2VviFla0rYtJMmn2UfaG\\nWCznQtTpQ4LypPJr3GnKGvOitkqkZf057kXU6MHIOyPbVwh567fMxfMHxcdf28u8\\n9/W0BSiuV7NFfhPYGWyU1sdil1tSl4jJEPnKg5vYjjM/DraddX2Ksu3fdeEqkz54\\nwa6V0bzAj2bn5ZWva/JyzBDbNSMMTa5EJcNg+QIDAQABAoIBABLFNFVHO3UD/Ozq\\n/TAxsUpq5DaeIDoAY0WfpKtMj7JVAupIHfIzWX3EfCiM4vgIxALmlxRaIHa7NIZ+\\ntzyduo0vrcoK9JgMxi/PBXrlzCikgcQu6PMRPpQM1IicejhuN6paiWBd+m+FJ0z7\\ne747BqMFYfxFW+/TEFER2MhiA6ss0/cvr5k5lK7ejPojIjyiObGaTx61bdIXhKMH\\nTQCf9BvqDm4bJnoorsjJpqGnZcpBtqUkl3VBXuCbuxvJCWjrWnCgRFhWm2ToHM9I\\nS9HiS185UPd4l18/5KBc02mf+hxCgf33vW0oKzjkI6GEIyaq3PAmWxMi1VtuLRER\\n76wfYwMCgYEA9xSujmMNX2x+4bs1qM4giQF2RgpAHlDzWBVin/SQQYv6H+JOaHBR\\nmNJ69omcTbAE5w65Q8tHoBNnT/UyF4M5btJ6mbYYYTtGjv4Gerc1UnsU21+PuH9R\\n8dOx/eDjfuOrz0yOv8+5uw1Qc8dkOIDdag7XYBSuxu5aapRVViFyy9sCgYEA1WeE\\nKFqSjsiAr0/jqFy/T/iBVyTYRWLcoq6a/GRwujk3xdRylNtkQ1eVt43XfXhNRFJC\\nWNy/Y5BCLzwSU8P/n1mPdhjy6XaBhRwWQ+JjheCXnJloiAC7wI9HwW/jgwBiDU8k\\nVgr8BJPi8bDEZdlDl89ztFUEJRxDF4BW3Kee6LsCgYAItITV5W2CMCtkPplMYj7J\\nNPD61L+fkdCRCOfZpN80P/9HAk0q5tIpJTlJ2F1Wa14w2dbzKYVTgXuBWK00IN50\\nJhxFsCG5w0HgJdkKl8vcJRP+CqbgpDO55nB99l9tiA30lsjsvx/XFEgCXEMOrpOe\\neflinDfwMFOlL6a2CyWlQwKBgGSvr4n+mdFmRljwv3/rKpSHsja0epnaODFFYnic\\nxxcF8guT3e/fx2GCjHALK1XWkdYfXZBhrqdCJAf3Nspw2kWL0wUsZkfCkv+Drfmf\\nccdznPTU6J3qgqsqrvdUXCqt3pVa9tDl49whDl1sQm2vYZXZ2kSGLCt6Nyl6cwEu\\n0OoLAoGAPJ4cDDs31nSgh4dTDh09QN51Y1ZxS+erY5jj+CzLhXteO5zLBkQVbOKs\\nBfD8wjvhus828NHNxQ9zsLO8aOZkd0yelOIV12a9kWkKduZpYua6Qjmgjri758GI\\nd4nIOAGBXPL/lvFCvgyGzf8KtRdBeau904i6hK1Pt6D87jFGjt0=\\n-----END RSA PRIVATE KEY-----",
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
							"cert_id": "22495571-cn-hangzhou",
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
							"cert_id": "22495586-cn-hangzhou",
						},
					},
					"previous_cert_id": "22495571-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force":            "false",
						"previous_cert_id": "22495571-cn-hangzhou",
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
  domain = "songwenpeng.alivetest.asia"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth.songwenpeng"
  type        = "TXT"
  domain_name = "alivetest.asia"
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
