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
	resource.AddTestSweepers("alicloud_ssl_certificates_service_certificate", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_certificate",
		F:    testSweepSslCertificatesServiceCertificate,
	})
}

func testSweepSslCertificatesServiceCertificate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"cert-tf-testacc",
		"cert-tf_testacc",
	}

	action := "DescribeUserCertificateList"
	request := make(map[string]interface{})
	request["ShowSize"] = PageSizeXLarge
	request["CurrentPage"] = 1
	ids := make([]string, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_certificates_service_certificates", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.CertificateList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CertificateList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["name"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping ssl certificate: %s ", item["name"])
					continue
				}
			}
			ids = append(ids, fmt.Sprint(item["id"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, sslId := range ids {
		log.Printf("[INFO] Deleting ssl centrficate: %s ", sslId)
		action = "DeleteUserCertificate"
		request = map[string]interface{}{
			"CertId": sslId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ssl centrficate (%s): %s", sslId, err)
		}
	}
	return nil
}

func TestAccAliCloudSslCertificatesServiceCertificate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence0)
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
					"cert":             "${var.cert}",
					"key":              "${var.key}",
					"certificate_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert":             CHECKSET,
						"key":              CHECKSET,
						"certificate_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServiceCertificate_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence0)
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
					"cert":              "${var.cert}",
					"key":               "${var.key}",
					"certificate_name":  name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert":              CHECKSET,
						"key":               CHECKSET,
						"certificate_name":  name,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServiceCertificate_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence0)
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
					"cert": "${var.cert}",
					"key":  "${var.key}",
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert": CHECKSET,
						"key":  CHECKSET,
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServiceCertificate_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence0)
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
					"cert":              "${var.cert}",
					"key":               "${var.key}",
					"name":              name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cert":              CHECKSET,
						"key":               CHECKSET,
						"name":              name,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudSslCertificatesServiceCertificateMap0 = map[string]string{
	"certificate_name":  CHECKSET,
	"resource_group_id": CHECKSET,
	"name":              CHECKSET,
}

func AliCloudSslCertificatesServiceCertificateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	variable "cert" {
  		default = <<EOF
-----BEGIN CERTIFICATE-----
MIID1jCCAr6gAwIBAgIQQ7/8/QOOTbywxdgSX9aMqDANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNTA5MjIwNTU3NDVaFw0zMDA5MjEwNTU3NDVaMCAxCzAJBgNVBAYTAkNOMREw
DwYDVQQDEwgxNjg4LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AMEl04gKBqJxV+8KideZb7S4mPysehPzr/cXu4i1RXT7UFtNVZuqc4IdIzOja2SU
6uNn8mY6Pfc5FNybg98bYx0ADbub55TUaw2Pz1CFEbiMvLpzMkp4EZadvmJWZk8t
dNb+ClKqdXUWhxApS3Lz+wjCNYQnlODk4KmxmM8/U/CyQS7lgWS/1G72UFB09Skg
sfvWdoHLrFfIlbVkp9XVELCtOkjj8Nn/rPOhc31NbstrwV4Whl6jngGAkaEtImJ7
//sL+sPPsutefCgfZPrC+Zwru2En1BuIo5KW02NYLdjXbABH8xjkUobqRoro7eY3
VySBr7adD6QmNv5hWohOuykCAwEAAaOBzTCByjAOBgNVHQ8BAf8EBAMCBaAwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiBJgXRNBo/
wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYVaHR0cDov
L29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlzc2wuY29t
L215c3NsdGVzdHJzYS5jcnQwEwYDVR0RBAwwCoIIMTY4OC5jb20wDQYJKoZIhvcN
AQELBQADggEBAHa0ATVeHtPPw1+a6kajlW6OQUjhiJg+Sk9fVA1eJ2Hzl1yDDw3K
yAyl1gkxGI6BwWdX/C8IE6PuPYcG2CmJGoFoEAAIbAE76AKABvHoA8I6wyDruxFz
06bNM8104TxAHTxe2zaHgBQnYIRk07uA8gxjZKFp1//eYbxj8HiP0Q9zXqYjF79G
Le4PDw7Q6U22CP+cT9Sz5ZEoJCzmUtx3uQWhLzNxvyISrXeSqAFJzjtL0KKSR1cr
8he6FoeU37oKdmrnweLeBe+no3OMChETa2JN4VAzXj/nPpQcyB7nXDfLUHe01+BB
ZBXKFLD2H38e97mFl/7mgNP5Nc1sycI5Sp4=
-----END CERTIFICATE-----
EOF
	}

	variable "key" {
  		default = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEowIBAAKCAQEAwSXTiAoGonFX7wqJ15lvtLiY/Kx6E/Ov9xe7iLVFdPtQW01V
m6pzgh0jM6NrZJTq42fyZjo99zkU3JuD3xtjHQANu5vnlNRrDY/PUIURuIy8unMy
SngRlp2+YlZmTy101v4KUqp1dRaHEClLcvP7CMI1hCeU4OTgqbGYzz9T8LJBLuWB
ZL/UbvZQUHT1KSCx+9Z2gcusV8iVtWSn1dUQsK06SOPw2f+s86FzfU1uy2vBXhaG
XqOeAYCRoS0iYnv/+wv6w8+y6158KB9k+sL5nCu7YSfUG4ijkpbTY1gt2NdsAEfz
GORShupGiujt5jdXJIGvtp0PpCY2/mFaiE67KQIDAQABAoIBAAKF9CZTUd8zvDKE
azo/Ur0Zf5omxgOBC/vzj0DLyXKr89KgMdhHmPG1YBKFIIU0XYCHXkclR05LAcbu
BdeCJpXS5zBbwDdAB9P/XHXQqeNvfJRc++ZgJ4QAXzkuqBssXK87ALcwFeUShxot
cphiWpW0inlwVkVn3WLUzfUV0+ARljn8VOf+aAmfCiQMl4gsBpvD3dxF84aihS+1
blqar5dE1GCJWHW67R1uSaAqHf7nwbBkZY8nTWF8n4+ELAAtlOgQKZlrQ+JxB3Ar
rWzgMj4M6F1/man1y/XPR56px9Xv3DwBZHuLufsqPr10q/nI9VIIQHe49sFgnN4+
48Q7wIECgYEAwxlrgBJI8gua4mJZxJRT8gBv2Mb1Kk1k7HVX11I+yF4eXr+cm+24
Cq7MjqmBXSnqvdQkwGFZ+C3cTKXJBPONWGF8NgiXaHSKjPEoFuHLdKBpgZMAax/L
aZBQRw6g12nz3XUCK0DE0wGgPkoDxc65s4NEWS+ua43LZ4TUOzWwwWECgYEA/XB1
ARNHyARy+P3iTeebh3t7qJoNoptLWHMlKjSjIZ1VZ4+9ilKsi5ZKVkPaLIjo8MGv
Ank3vzSrFSYhId0XfmSqoWySWc0eBkc6NERvopxuIV1WwRKf/18lLhxiEjHIcgds
G2KmfeiXdCKSgGlWvJmLITY4gJpOYMjpEDxipskCgYAdxnljmGbNmfvPZRcyKzkM
jAiF2wd7p0gp1lbLo9+1ELgt2ax7F7Ko3riVZUU7BLSwt/nL6o+iks02XW7qdIkz
3dzpGjKRXIfwrrVhmKBGclzny5mav8V5nO7DiXX+qkrvl3X3R/FCCtN77ivZOo2Y
2gXKXr6N55wNdnY1eyI4wQKBgQDXjZo2O+vFVuNimqyrjd1eMcxO7hfCwUooBGcL
qpFEucg1uK+Awig24LCBBly9nARjIJh1Bhw/58/KwQ9U+fJNcdkeSnV/I1HyDQqY
AczhBSM2BWkP9YNXc9jvivxudSECuwVblV/9nqGSCQWJag53gjAvIyqTVqpq7vYq
9PEC4QKBgGY2pj0ZNqGkq16jD3iS+DDBpX+TPnoHzu5GZCM/1GLZ6xXbpNWtZQt4
/m+6koRWeGvNAULnp8RSnhBzm+ZglpbwYcvsqRNDqIPGhJ2JruVA/bY3S0ebkRlD
xDn0dJVMvNyRR83ZpjTQhxoq5l56TN5xk1vdJ9nZdwJMmXiz2TrA
-----END PRIVATE KEY-----
EOF
	}
`, name)
}

// Test SslCertificatesService Certificate. >>> Resource test cases, automatically generated.
// Case Certificate资源用例 10998
func TestAccAliCloudSslCertificatesServiceCertificate_basic10998(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap10998)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence10998)
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
					"encrypt_cert":        "${var.encrypt_cert}",
					"encrypt_private_key": "${var.encrypt_private_key}",
					"sign_cert":           "${var.sign_cert}",
					"sign_private_key":    "${var.sign_private_key}",
					"certificate_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_cert":        CHECKSET,
						"encrypt_private_key": CHECKSET,
						"sign_cert":           CHECKSET,
						"sign_private_key":    CHECKSET,
						"certificate_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServiceCertificate_basic10998_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServiceCertificateMap10998)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServiceCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServiceCertificateBasicDependence10998)
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
					"encrypt_cert":        "${var.encrypt_cert}",
					"encrypt_private_key": "${var.encrypt_private_key}",
					"sign_cert":           "${var.sign_cert}",
					"sign_private_key":    "${var.sign_private_key}",
					"certificate_name":    name,
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_cert":        CHECKSET,
						"encrypt_private_key": CHECKSET,
						"sign_cert":           CHECKSET,
						"sign_private_key":    CHECKSET,
						"certificate_name":    name,
						"resource_group_id":   CHECKSET,
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudSslCertificatesServiceCertificateMap10998 = map[string]string{
	"certificate_name":  CHECKSET,
	"resource_group_id": CHECKSET,
	"name":              CHECKSET,
}

func AliCloudSslCertificatesServiceCertificateBasicDependence10998(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	variable "encrypt_cert" {
  		default = <<EOF
-----BEGIN CERTIFICATE-----
MIICOTCCAd+gAwIBAgIRAMyLoXyyYUjinX19tBRq2NswCgYIKoEcz1UBg3UwgYgx
CzAJBgNVBAYTAkNOMREwDwYDVQQIDAhTaGFuZ2hhaTERMA8GA1UEBwwIU2hhbmdo
YWkxDjAMBgNVBAoMBU15U1NMMS8wLQYDVQQLDCZNeVNTTCBTTTIgVGVzdCBNaWQg
LSBGb3IgdGVzdCB1c2Ugb25seTESMBAGA1UEAwwJTXlTU0wuY29tMB4XDTI1MDky
MjA5MzAwOFoXDTMwMDkyMTA5MzAwOFowIDELMAkGA1UEBhMCQ04xETAPBgNVBAMT
CDE2ODguY29tMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEEUV6G/mqnVAX8w5Y
pOS+CT78Qxv+p5XaOf/j/dE+nl50wvjNG90hNbc/xXXfncyaxGFia8IEv02OfyoI
7wH/nKOBkDCBjTAOBgNVHQ8BAf8EBAMCAzgwZgYIKwYBBQUHAQEEWjBYMCQGCCsG
AQUFBzABhhhodHRwOi8vb2NzcC5nbS5teXNzbC5jb20wMAYIKwYBBQUHMAKGJGh0
dHA6Ly9jYS5teXNzbC5jb20vbXlzc2x0ZXN0c20yLmNydDATBgNVHREEDDAKgggx
Njg4LmNvbTAKBggqgRzPVQGDdQNIADBFAiAmIZDik4VOlSdGNI7JnZb5qNQRxZ7I
3M7HVHCPLG/AzQIhAPw/gCRrotbanc4BXLNDqjmASDF6Rr3yMN85zRgkmyOM
-----END CERTIFICATE-----
EOF
	}

	variable "encrypt_private_key" {
  		default = <<EOF
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgeZV0qGwTsp+HRWoz
p5RjfAnjf23C0BHwH30i5EdxeHygCgYIKoEcz1UBgi2hRANCAAQRRXob+aqdUBfz
Dlik5L4JPvxDG/6nldo5/+P90T6eXnTC+M0b3SE1tz/Fdd+dzJrEYWJrwgS/TY5/
KgjvAf+c
-----END PRIVATE KEY-----
EOF
	}

	variable "sign_cert" {
  		default = <<EOF
-----BEGIN CERTIFICATE-----
MIICNzCCAd6gAwIBAgIQFMHRdrBlR1uG8XvOAKCqRzAKBggqgRzPVQGDdTCBiDEL
MAkGA1UEBhMCQ04xETAPBgNVBAgMCFNoYW5naGFpMREwDwYDVQQHDAhTaGFuZ2hh
aTEOMAwGA1UECgwFTXlTU0wxLzAtBgNVBAsMJk15U1NMIFNNMiBUZXN0IE1pZCAt
IEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDDAlNeVNTTC5jb20wHhcNMjUwOTIy
MDkyMjQwWhcNMzAwOTIxMDkyMjQwWjAgMQswCQYDVQQGEwJDTjERMA8GA1UEAxMI
MTY4OC5jb20wWTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAARd6pakVALwxocIsNN8
o8ZzKOLUmvMJ/lrvc6Jymkj4kON/o2deuZqdLJTuoImIbAElReEpSfLmb4kfhZND
hzy/o4GQMIGNMA4GA1UdDwEB/wQEAwIHgDBmBggrBgEFBQcBAQRaMFgwJAYIKwYB
BQUHMAGGGGh0dHA6Ly9vY3NwLmdtLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0
cDovL2NhLm15c3NsLmNvbS9teXNzbHRlc3RzbTIuY3J0MBMGA1UdEQQMMAqCCDE2
ODguY29tMAoGCCqBHM9VAYN1A0cAMEQCIGecCrkkVF6MhM+ZY157+3QiWjzCOicx
qJmwZd05e3++AiBeKMZ1UUuKVDJ/oxSy2i6U1oxOlzAxEN9GVdAZ1Q8lYA==
-----END CERTIFICATE-----
EOF
	}

	variable "sign_private_key" {
  		default = <<EOF
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgTCW/gskueHncKfwU
T+d32+jtTxNny9i6ZBGeea5+6XqgCgYIKoEcz1UBgi2hRANCAARd6pakVALwxocI
sNN8o8ZzKOLUmvMJ/lrvc6Jymkj4kON/o2deuZqdLJTuoImIbAElReEpSfLmb4kf
hZNDhzy/
-----END PRIVATE KEY-----
EOF
	}
`, name)
}

// Test SslCertificatesService Certificate. <<< Resource test cases, automatically generated.
