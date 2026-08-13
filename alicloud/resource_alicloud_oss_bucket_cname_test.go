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
			// Binding a custom domain to a bucket located in a mainland China region
			// requires the domain to have completed ICP filing, otherwise PutCname
			// fails with NoSuchCnameInRecord. The test domain is not filed, so run
			// this case in a region where the ICP requirement does not apply.
			checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${var.name}.${alicloud_alidns_record.defaultnHqm5p.domain_name}",
					"certificate": []map[string]interface{}{
						{
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIDQTCCAimgAwIBAgIJAIvaxVoOGHWMMA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDIx\\nMjU5WhcNMzYwODEwMDIxMjU5WjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\\nxMpKrFIX/x1iF/i51l2t5KKVWU0YRXX6DM5gSSeT6wVhsV2Y0LiS7+TE+WS/8tPf\\noDuNcPxj/6OxvPovlBADRnpV2LNxj2KkdAA1IBH+tqfdpE6pAtsYxrpxy5DNhriv\\nnSsX3ubnOOJeJN94yRmeZPBudq/EeMdNrBjTqPZgzTKEBnT19IsVyklsmgTVEaJH\\nNf21pAccvWwVxVQ+W9h2+8S1BM9HXZzTlu/toe7ETAB5Rzb++dApA+ioB9wZpbxG\\nlbymu0pYS+VLqa8irvuD2h8sC6WkBa1FUzlgP6AJyHT0SoNuMvUfzZhqiU+GwLBd\\nB7tKd6tdb9fRxULnbqmY5QIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\\nA4IBAQBLSn4SUeP2WfcE/k93+SB5hEZ+m7hxaA/sH8GYuJ1Uy/9X0nlmVgX040/I\\nfr7Zc9QG8c66TvvzscaWCAqwyeLAhHH0QzoIyX3S0k603G+XsTUDxfAuFaGN36aq\\n4x4SBhzEzQZI8rJj44mDOhkJRkg2QvemQl2aQAiAMgY6Ex5mayQVfwCPY5/9Aca7\\nXtyFNyT+DMVRS9eOIfkMKmp0kBIUDa3aDaEcxdj5L4nGcd8pyysOTyCBPDwSH4jy\\nmXV6aTGjwDlGxD0K56NXvk+FDL2eebaFbDUWbKF8qlovME3Td9u6ipdiRis+huVs\\n+LX3Ao6dfokhStA4mB8ULe7DP5Gh\\n-----END CERTIFICATE-----\\n",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEAxMpKrFIX/x1iF/i51l2t5KKVWU0YRXX6DM5gSSeT6wVhsV2Y\\n0LiS7+TE+WS/8tPfoDuNcPxj/6OxvPovlBADRnpV2LNxj2KkdAA1IBH+tqfdpE6p\\nAtsYxrpxy5DNhrivnSsX3ubnOOJeJN94yRmeZPBudq/EeMdNrBjTqPZgzTKEBnT1\\n9IsVyklsmgTVEaJHNf21pAccvWwVxVQ+W9h2+8S1BM9HXZzTlu/toe7ETAB5Rzb+\\n+dApA+ioB9wZpbxGlbymu0pYS+VLqa8irvuD2h8sC6WkBa1FUzlgP6AJyHT0SoNu\\nMvUfzZhqiU+GwLBdB7tKd6tdb9fRxULnbqmY5QIDAQABAoIBAESbQe1RsYv/cnNp\\nA2D4x+ctx2OavRt6RfKxAGCAq9EDz0tGlkAuGQwJdaJ8vW6q7wutt2Hsm/BD4XNA\\nxdWYv4uSmtsxtCWI/kxyxhKoM2T6oQrnYYTdXYXq+kE9+mk9efwRSgEr/vCV+rxg\\nLHvvsoj+SYSXQqfY0/trrF77hkQC+64/ECOW2dfs7GTMQ3ntypoA0fmH7m2vhPth\\nzR2MdeGdO96zaLUz8GW9wyQvlFx7NR8/33EHEtTttoerHa0sf7NzPhncTSOhLGQw\\nqN9xXoU+cNEWAIgbJigLv3I819O2/FySDU8Y98U5cwlVhRE0EQ9rjNvXcTizzSPT\\np/88OAECgYEA75oG1MvteNoSm+slydfWAbfH9DkSjYfm+fnCSNrAATxnaApU5XMT\\nmyGd3t7bOiQZsvm6oOLM3P2MWw+byHibUs357VpbPKCDR80w+K0OyZYqBPMzXqWK\\n1BTAPzx3iVsg9XZZDo/i2qq8RoWaWABL6hWEd5CIoheXsZtksdOIXGkCgYEA0kIu\\nU679uz+8JRnMqlmDy4A7t3Eea3iXL2SXmbl1dqpa0a9fDNhJshNHj3UQCamFDbcF\\nWIW60/Uk9OMUKAe55/LuoNnTa/vlD/3e0w9BXd3Oipi6f9bcXdjOBRGKoTo/vKM9\\n+5crSojaE32EQmVmN2ONRkViZ/tU2w735F0N+R0CgYEAiBB/Kp74J3YntTWPSxVv\\n6Z/VREKY35i6uWB1TWw0Nz93NaUQWxDDpIgtn+AMvPK9SV759d12G1U9PIUboXek\\nNRzVfk2enEpG4yKKWd3lFONaz17Q4EHAGfoCxqZu96ixidOAdX2OhUEKFD5QzQK2\\nGaPIiyGgBfTB51FomHeY62kCgYEAr1Y5Q9feB9SylU3aewSC/6VEJ0nR0FWT3hXI\\nxoz+A6M0cUAJx7BmZHXnax537VbMeg9yCcwbbL41v3HOUUOAKIlRrhl4UciR0LAo\\nrWE/ZYOexb1vaURIKIqv41IphSIKHMkU20XI+DL/iNlW/feJMg92tG8QDR5uOO9W\\nkb139ZECgYA82itIfQWhNQ19tMkg0wBvqu/oRpjXrfkdVrSYS3zxmRuXrSXP2Nl8\\nTl3aGCOuo9GSMP0DLQ5tKdmTc6Phw5uEyUAZqy2OJeroioz+lIRAEks48uFWoaMe\\nAaCmsgSJ6VnVtJxGXsDNxT3wA4vuv6MBHXnU+ME0KuBXKwHfoMkXSQ==\\n-----END RSA PRIVATE KEY-----\\n",
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
					"force": "true",
					"certificate": []map[string]interface{}{
						{
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIDQTCCAimgAwIBAgIJAOC5ZWJJiL1xMA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDIx\\nMzAwWhcNMzYwODEwMDIxMzAwWjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\\n3X9ntOcOxpCbvyC8zoU2oFqMLBjvc4tjlosgku4FJJ7CG4sq1fkuPwK7Vz+6KO6y\\ntuBiUsS/unqf8VGZgOKAhBQalOxxSbVAfdm9+Ngko7romJHnxGckfQE3rfOcYjXa\\nkV7gzydI2JaIi5E522RjqOuX8w9ytaBxVxoy0klYkv7IaPNOnX4wEF8j0d2ho1CO\\nHvh1cIGv+cATKBp+RkV2s3RkaqbtFLH/qPtZcaOoNdhc1u6LGlmlPN+0cMHN7qzV\\nfBgjjKLRxT7QqYFdS6U56JzIEL1dB64UH3c4Jdzn9QaU/lzvJkWkWUegaV3Xx1rl\\n3nThNtKTnqtLBbZfDeVNRQIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\\nA4IBAQBElss9qWStN5sfU/93EgQmBuuqv8rVxr8td99YN8J6Ss6/wgfMx+v88xZh\\n4E2SNIg1upD12MLHbAq3DLjMZ521AUm+r3KYEd+03kv6K9/hw+nA00PP5cvkq3Gj\\n15whGer9RVbIdYi/TZyHrJMmvNhEfWcwGSYTEXdhQAD7CmBt1i+om0ElDQ7TvhPy\\n0ptDhC94faYkMeEgUDJVoEBFcTy7NRla0Tw/yOeDB58T4QyaKwm8K30zalTAvLB2\\nFCrR8235FBvvF0pvU+Cs73H3NdJs8bdfWHULEU1/+L6Z3tBo34nXJf/26V2tDvqd\\nZe74hDbF7Vl//uWV1amR+g1ClQKG\\n-----END CERTIFICATE-----\\n",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA3X9ntOcOxpCbvyC8zoU2oFqMLBjvc4tjlosgku4FJJ7CG4sq\\n1fkuPwK7Vz+6KO6ytuBiUsS/unqf8VGZgOKAhBQalOxxSbVAfdm9+Ngko7romJHn\\nxGckfQE3rfOcYjXakV7gzydI2JaIi5E522RjqOuX8w9ytaBxVxoy0klYkv7IaPNO\\nnX4wEF8j0d2ho1COHvh1cIGv+cATKBp+RkV2s3RkaqbtFLH/qPtZcaOoNdhc1u6L\\nGlmlPN+0cMHN7qzVfBgjjKLRxT7QqYFdS6U56JzIEL1dB64UH3c4Jdzn9QaU/lzv\\nJkWkWUegaV3Xx1rl3nThNtKTnqtLBbZfDeVNRQIDAQABAoIBAQCuEeI+mRdTlXHQ\\n0rmO08IKYx6lyTLlazXoqY3/6m7ASMPjQYt4fUuK2WrBNqPmZzCr58tdoKHMu3HX\\nBHnOgDLfma0KPIcLlhYI0YYqejLROaJxxLiP8T8Lvlkzq6/Kvuf2NsoWApmNHUBR\\n7t+5OzvXFM9lhU5wzpZEDaLDAEFLwtgv2tOPH9UvmdxIyk8aQJkHj/znZRU1yMJR\\nZdZf//pHu8oIsfos+r7CZBv+6wFsfBB85XWwWCdPi+fujeE7VzE/czTAp698HlGK\\nH32r3CK0i5eZKwFsz0Cxz/cpZ9zSrx2wkkLW4ViUH3m4iTk+LEH1+bJZtBV+k1Lx\\nlOJ5qCRBAoGBAP8y7CX5e5hhIEYryS2dd05zMQ60SXQ3BGwVx2vjHEKG1j4VEiCi\\nawdty1M6S4dEu8pyZq5oWhCvSwPJ9hxSc5hpBWBIDRtLdYratgNcxD020e9Aqvnf\\nPXhFhhtoxfRn4vYLqoboXZdXHMvkdLJCzm8h8I0RgF0FrwfR8FAFx1FXAoGBAN4x\\nZn7g/bhLNE8fKv/4F4CvJ16B1hzt7TzZ+nBwkc+xPW9l535Gl41wiVeO3pvGauv5\\n2rebQ543mGQDp7+PsiNGRF4Cg9F/C8UCpC9UTr4aTOTgmcVBpZYzFtfmh0sPG5Yz\\nl4KFJH6gYzngl3ZwuAtc3dXrEybWfv+opAAq0WjDAoGBAMFj1ZDxfrf64npKtCnd\\nKoxIvuRlu0RWbQN7faREdyXzCGgDj7krW/BFQ8/OXW4kqCrChw2kBpyeOjqk0dyk\\nnvTgoTJVZ5lHlcuj8kqaAhxhbrXgS7EPe4WpKfebbmfIUjYioRea/1GwsiHQ/p4Y\\nAlg1YBWHLb9Qj1NdxL7foiwBAoGAesLOf1FtvRIH76Mnzc7TpWygks2nb8pg5dsF\\nTHRVi2vAprilwxXbi/DeYPr1sRlaX9Bm8ESfgl3zG2cNmoAZCvY6tbor/GZ2KT5B\\nWkj5TH0ZeOdC7kJL64WEnHqoy2aodj9A+YL4W+HfkM2uwWibtuNzSUqdBTtDZZtW\\nKSV/F6MCgYAhk2E7ubsSwlXTNxO7E60Dv83cYZaRhvdRTkllHHs2MiTyFMgyuU8u\\nX5FOHtoU4pGmiHK6yG60MPlN8FK147wS2+Ja+UjwH3EeEYDoN18ToYQsfI2oXcIW\\n5JbbQD9dT35VNue2agB8NsVdijcSFuczYtP/Jxdp8+31GfOwO8NIvQ==\\n-----END RSA PRIVATE KEY-----\\n",
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
  domain = "${var.name}.aliterraform.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth.${var.name}"
  type        = "TXT"
  domain_name = "aliterraform.com"
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
			// Binding a custom domain to a bucket located in a mainland China region
			// requires the domain to have completed ICP filing, otherwise PutCname
			// fails with NoSuchCnameInRecord. The test domain is not filed, so run
			// this case in a region where the ICP requirement does not apply.
			checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${var.name}.${alicloud_alidns_record.defaultnHqm5p.domain_name}",
					"certificate": []map[string]interface{}{
						{
							"cert_id": "${alicloud_ssl_certificates_service_certificate.default.id}-cn-hangzhou",
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

# The cname is bound by referencing a certificate already hosted in the
# certificate service, so provision one here rather than pinning an id that only
# ever exists in a single account. The id OSS expects is suffixed with the region
# of the certificate service itself, which is cn-hangzhou here and is unrelated to
# the region the bucket lives in.
resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = var.name
  cert             = "-----BEGIN CERTIFICATE-----\nMIIDQTCCAimgAwIBAgIJAKOyh837u808MA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDcy\nOTM5WhcNMzYwODEwMDcyOTM5WjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\nw2OvqiujcOzfz0vPPldCqDRJYIkc3jN8zhvvinK4/Am4wXlSvdbNTgoRlVfSnI3g\nW0zuxnmtRLP5hZ+0Pn1H4UYOGI6T6+wtnBtGPlSoO9y2siV5eBOGC2pDGT3yfxlo\ncwnQP7/Jre8QDsGyJyxJ+7+gpW56rnhzqP3cO2YAuiJU/t39bsmrIaTq+9sCiwPA\n0E+w8RasIKzXLgk8tQMtj4r5ffm2tQ9pPmE106v4MTc5gTFjgvgA+IIDjR2PN1RY\nFgEx5MLSDZWrE0OvzuV661XhBAq2Q7DrTdOOPSLUSC0YXb215iD0DVysjKJW8a9A\nYuoHX7QmWSOjrBxphmkgnQIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\nA4IBAQCZ3LOL3v1B3wKarbf+gEGXmUoa2BDmB+Wk6xEO4Z1BvkSPUqNPRMvuHxTr\nWnIg1962YObwtlWJdyi8UHJmXAbSV2D7Tn6ZZWFHPZqIOFXalZ0ZzsRKZsoda+Of\n2UmGzYjfNmBDk1eN7Qj7Lb7IWDVNIS2jFrCP8wEGwOmnUmr4c0ESH72Jdk3iojTR\nez6yZrZ2HndO8zdDoIYh2TgQPJG/tyHlgs3QP6RCy7jGeXZb0cYoecmD2RqT2NZe\ngR8kBTh7IyylTKFaVrAXiA5j46mYL0Sz7wnCTMYonxwl83G+FeJcmzXQgTYAPre5\np+HW1WeYLykCZvcVTv9CaLZHNoVE\n-----END CERTIFICATE-----\n"
  key              = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAw2OvqiujcOzfz0vPPldCqDRJYIkc3jN8zhvvinK4/Am4wXlS\nvdbNTgoRlVfSnI3gW0zuxnmtRLP5hZ+0Pn1H4UYOGI6T6+wtnBtGPlSoO9y2siV5\neBOGC2pDGT3yfxlocwnQP7/Jre8QDsGyJyxJ+7+gpW56rnhzqP3cO2YAuiJU/t39\nbsmrIaTq+9sCiwPA0E+w8RasIKzXLgk8tQMtj4r5ffm2tQ9pPmE106v4MTc5gTFj\ngvgA+IIDjR2PN1RYFgEx5MLSDZWrE0OvzuV661XhBAq2Q7DrTdOOPSLUSC0YXb21\n5iD0DVysjKJW8a9AYuoHX7QmWSOjrBxphmkgnQIDAQABAoIBAFrsK5WTowXspKWR\nUIphDtq5IiAzDeT0rrI227xgcGaQm5Ikw/UlXPpgwxfs+0vw1aOG5GIlwxSCb63X\nyId/wxA4ilyxFHKnv/2xz3k36eWZasbxm1neM/Vh6IF5izvL9gf6XBceR1qSMbW8\nOwvxlyf4X2g8RgikcqYEJBTb/aCfgWN4QxupH+mmdP5AjggvYsT/6+ABF3UbOWD9\nDj9Ynna42cazkiTwQFB9+zFpjsk24FJpprIBidtOF3A9Ol9PEN0MKJJj2+nR2cHQ\nHY9TNe/G/oBldaQJqxTeC7I3SeZ7mlUldgWK9BmLVBESeE7Plmh5pP6RG5wXgIxd\nhJxiTaECgYEA/s8/FVfltpK81kPAeyBBSBj2zK6o+LAa804p1TAQSrte4gyGPHBX\n5q5NZkPguqtBarRVA+eU+j19TD8EnzpwxlRi49RuccGvGv5ccmB8UeACOkg2kdTz\nGv8acWcOxNC0ysQg9wnLRUedoA1fH2jj/esYZwibc7mjMz/YWTQqCYkCgYEAxE1f\neYMqibN6L6inULSTSanASO686ezd0hOtruJPmjLfrKNk4Zsh7An307JfOse/uBMU\nQNHWpmymZEy9hvIj/nHBT+5c3l1DEBSNQjiBNSyPbK8everoKWqkMHwq4AQLZb57\nTSbiaD9+P5KoP0WBo9uZ4ETp48uFM6QhEXjcXXUCgYEAz6vPKTEDGmLLnwGHDZKD\nQiR+eOFc+5pjzKqGs6bBkHbXZPp6KSYSrgKfOFrX/Kt43GNu6ojCxZR52zt9I9z4\nbtv14OOQxAvsD98BL4Ltr7kXd7LFLuPU4srJHWW2BrhmsN9aUpzb23H7yKc9QJc3\nQgpqUAcW0yGYHjvJsyItpKkCgYAFfxEcSuLnBiJ2sSc2KEgzeNBMenrJpfs0BZ8I\nVYfbDm+a2txZQMm7XTAWOllWQP+KPOaFRhrXgBVMm6V24NLHLhI2lbr98uiMy7aE\n0yYzAfNmHKUkti4X8sd0IBXnPdW/3IyBRYRzXMvBJe8WDnEp0F1HnUZbPXiWUJMo\ndRTefQKBgGwgahqoWnHa3reMtxoGfcrJ+WxLSAr4Jt+15sS63+JT9r2grRMo2p8B\nWpsHZDqDIMEP3IqTiqLuaSUnN5Cx9YOXhIFFYgDZOunEeYmei5E5y0BCOn+zBHYS\npM2shXcGhQX8vp8ZIUlUsnPEXoL/gVi66Ao1/h41lGC7eFqEAAAB\n-----END RSA PRIVATE KEY-----\n"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "${var.name}.aliterraform.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth.${var.name}"
  type        = "TXT"
  domain_name = "aliterraform.com"
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
			// Binding a custom domain to a bucket located in a mainland China region
			// requires the domain to have completed ICP filing, otherwise PutCname
			// fails with NoSuchCnameInRecord. The test domain is not filed, so run
			// this case in a region where the ICP requirement does not apply.
			checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"domain": "${var.name}.${alicloud_alidns_record.defaultnHqm5p.domain_name}",
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
							"certificate": "-----BEGIN CERTIFICATE-----\\nMIIDQTCCAimgAwIBAgIJAO/c/EfBUd+MMA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDcy\\nOTM5WhcNMzYwODEwMDcyOTM5WjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\\n3JV1LxzcQ4wXMiGB+tk5UU9WtDdJQTx0I6QjicMMRr4gkuzmWUl+5riVXGlquoS4\\nbuBYbWA9L/6o9oDVXMcGoXqGTpqnBK1RbD3/dMkuCU32ovBWJt4ybfkhw6x3mJwy\\nHKalAcaFj3OCCzOvaIXGCj6WXxSUniGUxs3M3znVUzaWIUJEwp7q0K3l4QZmddAD\\noGi+12VZYOszQX7r8e/ETDEOXxpO6eEZP/+vE6uX938F3KaYqInyuoxHQnyin3I5\\nFuvUjaod/NDwF6kKyBG3GFWOYmT7dpjftxVrY+lmjRvlY/uyRDhAKZIr5dBFGkhw\\ngZFI1davKmnfkADO9ZIwewIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\\nA4IBAQAjl5SrxrXCl7gHM2gOHo219D1yXiUs3xjv0BmxFvNpnXjoUuDWUXvFRYu2\\nxPvRo/8QRrgUY4jWKqYBmmdPBUivmz5/wYBMTkezwMurqs8t1Bqd+aqtxAzUHwi9\\n0IzvEN1gk9uCHe82+a2bg6Nm3ljlhDefT4zuRwxtEdiRZOrPLyCxphZaxNOdD/i8\\n9WtZYe0nyh3jgMUlqesOyUzDRJzA/dRKDsraI5TSwb8rwBTJQCCai5/4UKKDNxhg\\nwL+Y0juEQKrlPhiL9pRqynHJv9iXOlf+brEIU+VE1zYooXNLB3Y/X37JXjlNLstg\\nYb78/x6S8d8XhBvwpn2eIr14bqsK\\n-----END CERTIFICATE-----\\n",
							"private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEogIBAAKCAQEA3JV1LxzcQ4wXMiGB+tk5UU9WtDdJQTx0I6QjicMMRr4gkuzm\\nWUl+5riVXGlquoS4buBYbWA9L/6o9oDVXMcGoXqGTpqnBK1RbD3/dMkuCU32ovBW\\nJt4ybfkhw6x3mJwyHKalAcaFj3OCCzOvaIXGCj6WXxSUniGUxs3M3znVUzaWIUJE\\nwp7q0K3l4QZmddADoGi+12VZYOszQX7r8e/ETDEOXxpO6eEZP/+vE6uX938F3KaY\\nqInyuoxHQnyin3I5FuvUjaod/NDwF6kKyBG3GFWOYmT7dpjftxVrY+lmjRvlY/uy\\nRDhAKZIr5dBFGkhwgZFI1davKmnfkADO9ZIwewIDAQABAoIBADsw19smkWyGwQqw\\ntyJK+/h3o7qEQ2IACOIvf2HONxMcnb0PWNiIwkbDLUE5AGzAhIUsKk5fTsv8N/a9\\np4NX3M2kBTo+gabdo0W6dTwvZ+0TQKWEfHm9kia0fXz2YLlQ4JmTlh+d1+Ugh7rd\\nyanwi63gEZW9/gtY04VtYBZefIHxWjr2NboEZDt8xvpWlcttufpCTBYggxLjoZMy\\ngLGJHk5R9gJZDKrdnkg2L99xNhOu/7N80k6+Hi8/jDNAOW9hoRr6fo1gK1XqXjFL\\naBsh190LgA/VYo5MR3r1SLmTVrnZya1j+xqGBY4ZjtJIzh4D+n3h7giTQ8/tr8dH\\nyzsJsUECgYEA+HNopUB37x8sIucarCV7V2Nw0oz07LAASuYt7AaHXNNuphwyP/no\\nJdaWP2cm+PDyNJWTFvXCXJyvwEjAXfM8+nrl2dPLUeU1uQ9BkQZ3iVsNeA+4+niK\\nJZ3nXZRDSph4TLXi5RSbrV+jksAfCzHP59iFP5jK2GPDiTovnrrcIJECgYEA40lI\\noAACu4caI0xoRp1B0MbKWr5N76AyJJxbKX5vxbjdXLLVrC8HhTA5mb/GmcWjn21h\\n+NbnbHaxPh7ZhJV0U+IMa6Heq6D4hiQCmqW2M3ZApaHP3eKtmwt0v12FUob1eNdb\\nQwDWidHkoLgcxImOvZAIKkaMMn89gHL+49ndRksCgYBBStsaapnaPp/zwDZTPTpv\\n2dNBkgef2BULmfhBiemy7GGsx8Yw5/UpVH6BxRMJ4xBT32cbZpSgkBDkAHqFdjH1\\nRaz4FN/e8tSugKLjgQaTE1mzzrX3JQxxHFE8V4VjqjQbPMWXHFZZNsQfAdxmrb2M\\nmWtTLk1IltdBTghLt6G38QKBgGkRf5k3aAv4sISQ1cOO/tXcj77TKoQTshpqjVnp\\nMRJeGza3FT+7neZcHMSOeuirDLCuiBPYhLMHS3hEGpnH3TbJ0KQQ+Dau+zRHgUys\\nPkYb7FalLsqL92UtLpMoUHGOIfvy0iVvRb4AYYhKlEHmtS28X4nrgvP1DiFLB7md\\nBUVxAoGAIF4XBdrk7QrV9MU7UR7SONwiCO7dWwh1bHreT7I5XSVNcWJb/zDmOH9w\\nPRhOxmn5ywYByj0W2haKJWaMw/MnBRrs3z307/6dBqx3rksYeEcTq2bh7Njw7GK7\\nCDaVYSPtm/QbUhpY3hk0UWp/omgFuc+KURoBPqMORPfC67davWk=\\n-----END RSA PRIVATE KEY-----\\n",
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
							"cert_id": "${alicloud_ssl_certificates_service_certificate.default0.id}-cn-hangzhou",
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
							"cert_id": "${alicloud_ssl_certificates_service_certificate.default1.id}-cn-hangzhou",
						},
					},
					"previous_cert_id": "${alicloud_ssl_certificates_service_certificate.default0.id}-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force":            "false",
						"previous_cert_id": CHECKSET,
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

# The cname is bound by referencing a certificate already hosted in the
# certificate service, so provision one here rather than pinning an id that only
# ever exists in a single account. The id OSS expects is suffixed with the region
# of the certificate service itself, which is cn-hangzhou here and is unrelated to
# the region the bucket lives in.
resource "alicloud_ssl_certificates_service_certificate" "default0" {
  certificate_name = "${var.name}-0"
  cert             = "-----BEGIN CERTIFICATE-----\nMIIDQTCCAimgAwIBAgIJALR2yS78+SqiMA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDcy\nOTM5WhcNMzYwODEwMDcyOTM5WjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\nn0La3EtHPbZVmz5gVY0tQKUR415YmKKfCblRzV538GPvdzekLmNIUdcLiTPw6OY5\nMAM9B4iCIBBGmOPXqK5NvoojL1nlWH7Rbu6DbUMTLeDkDNVT0pkz3m08VDgmjEpV\n/95qJNegsElHnz2eJz+kfhe6bV8lRkfINySDPPy0ovM3qSgLikNkR+U0si8eCqHt\ntY5iEsz0zFlsp5iFpvE6QVf7UdIk1FiF3JPt/UoKn1OY88TDvn5B4n+hezrYZzZs\nSCsDztPMXio2NIlGtv51FG6FYRQL31md+ZWWb6yGy4EPSmtWHz2urDxRvujjJnvQ\nHC/CJLX1Idhr1KbfbkP9CQIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\nA4IBAQAR0Ps522+c05JT/JzTA0UMkB9w9knu6ffQjI3FzzeUEKF+dGIeaN22Z/sw\n+jNAEkOzXJvIRbAMoULOWCFqh3XvVEyS+DhEKPRr8Lk+jBeto619DPPQIldFWE0C\nGp2vIYmkKFiBEOLBBOZkEQtOuVuuBGt75x7GUg4v4UYhV5LflylCnIn2unW+Xr32\n+3gWIB4NVI3b90FITMS/TvwXEW7JXLdDcROCVgLkXVIAG1vjwPudftVa7EIkYgNp\n3D/1bTX/b+EqlWLAeDwFTqL18wk3c7LE9pPzj8zYL/ERPV9fPf57O1D6/0reiwK0\nJH0ikZMJ1YrKc95e7b3NFnRMTZD1\n-----END CERTIFICATE-----\n"
  key              = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAn0La3EtHPbZVmz5gVY0tQKUR415YmKKfCblRzV538GPvdzek\nLmNIUdcLiTPw6OY5MAM9B4iCIBBGmOPXqK5NvoojL1nlWH7Rbu6DbUMTLeDkDNVT\n0pkz3m08VDgmjEpV/95qJNegsElHnz2eJz+kfhe6bV8lRkfINySDPPy0ovM3qSgL\nikNkR+U0si8eCqHttY5iEsz0zFlsp5iFpvE6QVf7UdIk1FiF3JPt/UoKn1OY88TD\nvn5B4n+hezrYZzZsSCsDztPMXio2NIlGtv51FG6FYRQL31md+ZWWb6yGy4EPSmtW\nHz2urDxRvujjJnvQHC/CJLX1Idhr1KbfbkP9CQIDAQABAoIBAQCfNIvo8G/VJzLI\nsEBJBYoZN2p8alISs25coB9AN5Gag6xc9whvPtyKw3hKvdu0VoEQmAwoPbQnLV4F\ndK6fdy9MrHaj3S/BmXTvegtz7Dt9/3S5x3+15WTOk1BduIwAbkcuMz7UeaGu2HJ6\no3Q4NAzR6BJ7R0PRz+w8A4oWK2DACuLp0adF+tGMls3OM1Unc9oVeZbDSdmWdjBk\nGP1zPC8E2c9UsLfJ93/XejlIV0CuFcRqYLvmobbx04h42KFprwZQqpq2Frq1UInb\n8Xm9L3jy5uVT1gsswBE7794Djsp1pBNYSUPd6LGEyZX/N2Bhx5hQsHLlmXXaAQRa\nrxuK+rOpAoGBAM6IV8f3zdMrwBjBEakj2eYuf5XRkxqSm1cNq0TSe5VfPo3+KoYF\ngPHlYEXrkXW9+M+MuyxrxXnXWBc6TZsQeuIrrK6X03psqE+GpefUOQB7Qp1okibz\ndKrkuFUQHpezTHjnnPrdderA9XIHISM3nakOnVGjKJHmpCBBb0XygT/jAoGBAMVo\nCdFtBEHX7OS3XVsJxgfRwYtgxfNjl+d3rO3BGFicgp4vWnOyc9mWNUFyVllq6CkN\nHLawDS+4hfYR0BaTFknpSQNBQ8LxLKAtWP1PuZtyIrTzZBVZQxswToMnyFYYoBaj\nd02fb0DiiQSeJixzCYivrZtBLX+KLDaFtPoH+IsjAoGBAIUGxKOENQpjD6PiF2H+\nOYdNQ9hX2IwxCeUUZNA7UmZvpncG0pToTpl/yHbAuDxCVFQ6rQR7lgJYdeDgKMRL\n5RpwTxVVrV0ZR3+RlqKvytdIjSueAyUbgnXgQ+pmK45Camslo7LhmeXOy0ja1rk8\nRUxyoVnH4YW4LNapzuYawK1JAoGAdhjKnt5wSI/L6fyEvhz3ut/SwPZRFk2Dp/ch\nnk8BqKlhPw8nNsYQpqBFFfU4EWByqXRttCFYki76/X4klgzCrc8BXhAiYLJ1txHK\nBik26fb7KnPdcSQokFBy9+XJ5S/wPfrnOanjHdcoj3mpbrgXgQ1Qd+wjMwTPdILD\nBT3VhC8CgYEAo1GBY/YbyhKtx7fsVmqIqH7Nz6IKU/50kSZDv0a4tnM27uYLstqu\nm3eNvh0Esj33yKUokUlADqkLgpZmWMEZEsLfp+l64sphKrvcH4bn4fvQrvsNeIsp\nzFj5DKCcbUZU90NFt21HLP8U4YTFPeajjmABdHHHUkcsWBHNjx0woOE=\n-----END RSA PRIVATE KEY-----\n"
}

resource "alicloud_ssl_certificates_service_certificate" "default1" {
  certificate_name = "${var.name}-1"
  cert             = "-----BEGIN CERTIFICATE-----\nMIIDQTCCAimgAwIBAgIJAM/fKNeEka6PMA0GCSqGSIb3DQEBCwUAMCoxCzAJBgNV\nBAYTAkNOMRswGQYDVQQDDBIqLmFsaXRlcnJhZm9ybS5jb20wHhcNMjYwODEzMDcy\nOTM5WhcNMzYwODEwMDcyOTM5WjAqMQswCQYDVQQGEwJDTjEbMBkGA1UEAwwSKi5h\nbGl0ZXJyYWZvcm0uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\n1XPjTUtP49nkD7MwYefLtioD71heZPVN0AibgefBHANatcQqpu1RAX7cfzugbicb\nw0zBnKHtKTqR9c41YH87q05mO8SBzGbvrXqUH3Yb98COle0rQQVuXz4WHz5fIMlV\noAqGXbJoBTq/1vGLk8DMzgIBVNUn/oyG8W8uzcj/AIy6gqFawIFDwbP0z25ayuk0\nh5dmTearBoXlfhyM966+PlPAJqByT0ysHxp9aAFAKdl2F+J9s93DOaMlppCd5Y5A\nTrq9iEqFRSb2Rxnw3vtt9ggcsHf36E9R7+6dntNNwIH/0Ff0i3JmlH0IzNtukxLN\n5JKZ+xwnXjqYqhAhtT4D0QIDAQABo2owaDAvBgNVHREEKDAmghIqLmFsaXRlcnJh\nZm9ybS5jb22CEGFsaXRlcnJhZm9ybS5jb20wCQYDVR0TBAIwADALBgNVHQ8EBAMC\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUA\nA4IBAQAwBGi/pGt9+IFCnBAecktii5gY9zsI4NH++DCDy+E49v/nQNjQb339iatX\nrqLhMsoF4x7NQeN7AxtpCKeykRR4A/TARDc8ZLzzTfDpUgnyFOxtXqzS+BmMyEru\nqG/75V9W3C1s5rAzKGMNN1ROZDpbnkCEKZlz18C7YHAMADurRVz9QdAKTtvQ1vTC\n99jeW5FquFpRg8Q7EHmuvxTntHXjhb+g5UdpnxiGaZk9wQo8Oq0qDkfv24cSmoRY\nPGGUi2jQFwZ1shBcLR9MBfOx7zqTTwRPA1u8FGn7pQXUkwzx+1gqGk87LRVL5yEY\nQEHfwPfhmnJjbC+ErkEvh4Y94OHW\n-----END CERTIFICATE-----\n"
  key              = "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEA1XPjTUtP49nkD7MwYefLtioD71heZPVN0AibgefBHANatcQq\npu1RAX7cfzugbicbw0zBnKHtKTqR9c41YH87q05mO8SBzGbvrXqUH3Yb98COle0r\nQQVuXz4WHz5fIMlVoAqGXbJoBTq/1vGLk8DMzgIBVNUn/oyG8W8uzcj/AIy6gqFa\nwIFDwbP0z25ayuk0h5dmTearBoXlfhyM966+PlPAJqByT0ysHxp9aAFAKdl2F+J9\ns93DOaMlppCd5Y5ATrq9iEqFRSb2Rxnw3vtt9ggcsHf36E9R7+6dntNNwIH/0Ff0\ni3JmlH0IzNtukxLN5JKZ+xwnXjqYqhAhtT4D0QIDAQABAoIBAHiBR3cQqJajIYz3\nhb4QRcKe77/FLO1kS7zBz0FEnJH7Fs/9YnMBEbV9cHBoMkddzt+wSrHp/OFEzrht\n5VaIHiC1TyQ46WqDRpay2EL2xA1X6WedEMlRjqE4hPa4mK4C3FNQ/dCR8wXYyAtK\nLJmKxFUdbrD88epUXa6aLVtCOSyOPecMomZehPDZopFr4dLgjImJbJlaavV13Qzt\nGHuor3Uh8lE3YKdFjr5Cf2A2amIF2ukv0Uc977MLJ5VFWcb6ZgqKEdiqv2SukTCv\njFdYtsbl0HOtMLDNOOqWdv9NB0H0bpNC/w+gT7YdejCd7x+Ga6FZ8e7rq25CwMGz\ndi3N4AECgYEA+7YnFejET41jQ0cbl99F/6yyz5k6nQ/rk/DtQmiIIzEzpVAIJvHw\nrTQA6/HhVmnQnJVBEDjf8RPCfzjB1cXXEnGt6HGOJP1KYvSo3UiLM324tWQZ8ZvS\nox2Gy5EDEy1tjJWLDL7E/kuUbxscqXpml0lPhFvZXqtqCVQiXA4y2dECgYEA2Rbe\nPW6t0BH6hudMeidlMAnsH2axkuaiflAKIgBG6sZDmMPOAiDdOUvOekjQGNo1OOrj\nuSsKLAGrvgbtkki3ie+7HzqWJ4/T2C3iCaOCGX8GR6yyktNsjnEZenWuzj/Gp6ov\n6KHmoEOrZhWEkXS2vHcQcQnMDn3zkgmzEjs5CgECgYABUEQH8z0DBUPdWAOm2T1u\nRiJwvuX1Z93c2ccDL7R2Ko2QcUh5m42b+cd/c7WvU8II7yZ1xTY19dpv+4XXbb7f\nk8RKkD0jqEa5GXnAHd7MF/3cxHb2Mc/5le/cJBeWBAisUSN2n5A7m31czxFpOQBM\nDc/iavBJdC+LeOrs/A374QKBgGDCs54YJfrW+J6Gm+zagFyQH6HDaSS8DfNVA58y\nFmnwoxKFO95w/YnbQxX4PGDHae+LqqLPD0KcIAucFOod5UjjBLmfqGvLzLXPha+c\nJJHur0LlM9cDy6AVwzB1IcwmWwpCbgY3m48VemEO+D7JEeYg/8ASiNRwyU7vadSX\ndw4BAoGAJHFvJ9g5JIvHquvXhd/w0clctn6DCM1TrrQFwDbQhl5TV/qr43q9qyAI\nuGe4Lhrlta/eJeUj0ciQf1up8B5dyUgxk2djRhhr87XtEPglWAPSmT/xQfjNs/Hb\nPhj7wwGEujShgkrVtQNCV5QzCZCHxLROVgoAucsJ9LBrj7k/QOo=\n-----END RSA PRIVATE KEY-----\n"
}

resource "alicloud_oss_bucket_cname_token" "defaultZaWJfG" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
  domain = "${var.name}.aliterraform.com"
}

resource "alicloud_alidns_record" "defaultnHqm5p" {
  status      = "ENABLE"
  line        = "default"
  rr          = "_dnsauth.${var.name}"
  type        = "TXT"
  domain_name = "aliterraform.com"
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
