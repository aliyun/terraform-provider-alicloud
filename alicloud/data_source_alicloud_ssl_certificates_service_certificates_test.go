package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudSslCertificatesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "data.alicloud_ssl_certificates_service_certificates.default"
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAliCloudSslCertificatesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
		}),
	}

	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword": "${alicloud_ssl_certificates_service_certificate.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword": "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"keyword":    "${alicloud_ssl_certificates_service_certificate.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"name_regex": "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"keyword":    "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
		}),
	}

	var existAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.#":                  "1",
			"certificates.0.id":               CHECKSET,
			"certificates.0.cert_id":          CHECKSET,
			"certificates.0.certificate_name": CHECKSET,
			"certificates.0.fingerprint":      CHECKSET,
			"certificates.0.common":           CHECKSET,
			"certificates.0.sans":             CHECKSET,
			"certificates.0.org_name":         CHECKSET,
			"certificates.0.issuer":           CHECKSET,
			"certificates.0.country":          CHECKSET,
			"certificates.0.province":         CHECKSET,
			"certificates.0.city":             CHECKSET,
			"certificates.0.expired":          CHECKSET,
			"certificates.0.start_date":       CHECKSET,
			"certificates.0.end_date":         CHECKSET,
			"certificates.0.cert":             "",
			"certificates.0.key":              "",
			"certificates.0.name":             CHECKSET,
		}
	}

	var fakeAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"certificates.#": "0",
		}
	}

	var aliCloudSslCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ssl_certificates_service_certificates.default",
		existMapFunc: existAliCloudSslCertificatesDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudSslCertificatesDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudSslCertificatesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, allConf)
}

func TestAccAliCloudSslCertificatesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ssl_certificates_service_certificates.default"
	name := fmt.Sprintf("tf-testacc%ssslcertificatesservicecertificate%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAliCloudSslCertificatesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"enable_details": "true",
		}),
	}

	keyWordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}"},
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}",
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ssl_certificates_service_certificate.default.id}_fake"},
			"name_regex":     "${alicloud_ssl_certificates_service_certificate.default.certificate_name}_fake",
			"keyword":        "${alicloud_ssl_certificates_service_certificate.default.id}0_fake",
			"enable_details": "true",
		}),
	}

	var existAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"certificates.#":                  "1",
			"certificates.0.id":               CHECKSET,
			"certificates.0.cert_id":          CHECKSET,
			"certificates.0.certificate_name": CHECKSET,
			"certificates.0.fingerprint":      CHECKSET,
			"certificates.0.common":           CHECKSET,
			"certificates.0.sans":             CHECKSET,
			"certificates.0.org_name":         CHECKSET,
			"certificates.0.issuer":           CHECKSET,
			"certificates.0.country":          CHECKSET,
			"certificates.0.province":         CHECKSET,
			"certificates.0.city":             CHECKSET,
			"certificates.0.expired":          CHECKSET,
			"certificates.0.start_date":       CHECKSET,
			"certificates.0.end_date":         CHECKSET,
			"certificates.0.cert":             CHECKSET,
			"certificates.0.key":              CHECKSET,
			"certificates.0.buy_in_aliyun":    CHECKSET,
			"certificates.0.name":             CHECKSET,
		}
	}

	var fakeAliCloudSslCertificatesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"certificates.#": "0",
		}
	}

	var aliCloudSslCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ssl_certificates_service_certificates.default",
		existMapFunc: existAliCloudSslCertificatesDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudSslCertificatesDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudSslCertificatesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, keyWordConf, allConf)
}

func dataSourceAliCloudSslCertificatesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIIDeDCCAmCgAwIBAgIEN3ZT6zANBgkqhkiG9w0BAQsFADBVMQswCQYDVQQGEwJD
TjEVMBMGA1UEAwwMKi50ZnRlc3QudG9wMRAwDgYDVQQIDAdCZWlKaW5nMRAwDgYD
VQQHDAdCZWlKaW5nMQswCQYDVQQKDAJBQTAeFw0yMzA4MjgwNjQ5NDNaFw0yNTA4
MjcwNjQ5NDNaMFUxCzAJBgNVBAYTAkNOMRUwEwYDVQQDDAwqLnRmdGVzdC50b3Ax
EDAOBgNVBAgMB0JlaUppbmcxEDAOBgNVBAcMB0JlaUppbmcxCzAJBgNVBAoMAkFB
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzkk9NJUH7PLSQK4RRrGQ
Y5dVsftkhnKh9HhI6yrnlowWIDPS1PZHOU/5gQ7xPUPGdKQV5S7x8wROnAaXEimx
N4GdQw25pGhRJvlwme9fzJJiSe6lG49NCxmuBiEdJAzPKaTPpK1cG1f1TqdgCfHR
HAL6Jxb3ylHG2LlTNFLXikubUi5RT6/9C8psr713Zm4HveCI/cx0WdgZ+fmsc9ft
rkIB1DdyV1kQ51m8r2rLi3J7aC5ggGOiex/VlGSd4e6SOQLpdQEdDbodtOJ4LgVk
+arFNCMinUWIOPGFzXhdm6lssPbh4MOBrz8c/M9TcF4hoMn5/o/9johZIZ/DOvXt
ZQIDAQABo1AwTjAdBgNVHQ4EFgQUOnWiddgeZj17IeysatqhE361o5YwHwYDVR0j
BBgwFoAUOnWiddgeZj17IeysatqhE361o5YwDAYDVR0TBAUwAwEB/zANBgkqhkiG
9w0BAQsFAAOCAQEAfh3cnOszHM/5wXjY7BIkmgDOReksS+87ibhBz7T2ddZj+yCF
9GdIBzXCiHpQFXpW8a3kc3I7l3nGfMTkmF6ld3ot/6SXP17QKJwxtvUA4ib8QkWD
S7FT+UcHCUHv42Sh1e5uAlQ5pMSul7iKcR7jwlwZGZ0905HOqrmdyUGJ+Ud2uZWD
AC0dJF6Bv9VhNtci8Imp05PaPH09deXLZu8LRrKRZFy9qLW5R6Swv7nzxckOAqDk
TTc40xwvQROekWUyxeJL7xaHuylUHE0bxsiIfx5bZsBizRjprIwGlj85CSPuTZyP
DPfaiZAN/61h5HNAnxLltOZfqabKYYw7l9LBDg==
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDOST00lQfs8tJA
rhFGsZBjl1Wx+2SGcqH0eEjrKueWjBYgM9LU9kc5T/mBDvE9Q8Z0pBXlLvHzBE6c
BpcSKbE3gZ1DDbmkaFEm+XCZ71/MkmJJ7qUbj00LGa4GIR0kDM8ppM+krVwbV/VO
p2AJ8dEcAvonFvfKUcbYuVM0UteKS5tSLlFPr/0LymyvvXdmbge94Ij9zHRZ2Bn5
+axz1+2uQgHUN3JXWRDnWbyvasuLcntoLmCAY6J7H9WUZJ3h7pI5Aul1AR0Nuh20
4nguBWT5qsU0IyKdRYg48YXNeF2bqWyw9uHgw4GvPxz8z1NwXiGgyfn+j/2OiFkh
n8M69e1lAgMBAAECggEAevPgTTT+0lYwx2h416ACJboP09O5KQGuUl5XaAPcoTjB
/1OkOFbKQPjQCAJ1+0QoR2F9w2plv6kziX/MD4FWJXVV3J+TpNCgfhBy8u1gNjiR
6Osa8gBJtXIK7ZBTJCeWWoXnVYoWuh2FEupkLck6D+4eV6oy6x4u3QIo+6jc24n9
dIXQG6/v/Iao34kB0LUdp/4WNaUDvfI6NDhEwchpKE95dtWIDlIN/YhfiYAdjrnl
YmH2VDbAGgsdEiHP4wLZfjgsGPPDGS0+qBHoSiJGH0E6wWEZdAE4TsYGRFsO86n3
LfjEPFGfPlcnZe2cTTe3kmyKq/DTjxtu2rh3I8o2CQKBgQD/5Xe7cenaOBefzPlx
GOEsB+qv49UmzEPOXDNZe9hmAawuuuxPUM+xlE++P+mEgQm1LPT4WWgtFLPVuwmx
ncxt4CJNZh+ZGFyAZ4dm4M4ZhIBXNonyIP+yGyAJUUVF9Iy3TYcJNiGzv2Rx9JRQ
XWJMQnTDILmZbmU+ltTea7/zqwKBgQDOXqCqb17MuLt7OcKWSgthm79OlaOdzZvl
i9qU6VzZKG7Axc5gA9yq6tHp3vWPI4bNdvwqIIa/nzVILjGA5fcYFbRN+7gHwo8s
rNAgi5PAoKWqQRovyJRY9Eq/sn6l1jbJZAOUAMZMWDm8z89OqK7PNQSIAtfFSneo
2QxJkGeTLwKBgGJkafBB8af9b1/7YWISLepPNPbihH/BhMThAMGEdAVs2TaymtA4
g1OFck/1pSVUtFXcbmjbf8ntruQcYbLQuNz6lFXsUXP9QPwCUrbE85ouL2bZSps2
AvsJoPzUKe2nBUAp6CUrkjPaAJYsc6ae8X/fAaRRfeu33ef9+OV4yrq3AoGAYFZo
ZmfrN2Kdkt7Z6dLTEVPlsMfGQ6pyNmxdM9rkzzNC0JcGymfDIb7RE35T3+hTy6La
AMiCXv3xn6qAzY2NFh87tpPlyymWMOLTnf3Kkcfszlfp45idOBGCu46V9NDVbppT
2UmrSIR/H5dbTXsNcAlt/hhlpeInjhkU1VqmH10CgYEA7Kk+QhWq705SczpWjm5J
9kHqfFzJLwAWNBduiia0WypgPhLe/4wT1rYQkBtKMVKrgFo7Cvi4YKlrtlDnXyeU
CIFqfEL5NriQelqrFsvgHsmD+MpvDoSWm5C8IrTubtlNyWUzXSVT4OIwzPobzPqG
LILJ+e7bLw8RrM0HfgFnl8c=
-----END PRIVATE KEY-----
EOF
	}
`, name)
}
