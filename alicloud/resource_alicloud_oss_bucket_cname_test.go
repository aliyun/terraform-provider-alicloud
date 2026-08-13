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
