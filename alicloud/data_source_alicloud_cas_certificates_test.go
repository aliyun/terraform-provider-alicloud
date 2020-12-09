package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCasCertificatesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cas_certificates.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccCasDataSource_%d", rand),
		dataSourceCasCertificatesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cas_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cas_certificate.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cas_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cas_certificate.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_cas_certificate.default.id}"},
			"name_regex": "${alicloud_cas_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_cas_certificate.default.id}"},
			"name_regex": "${alicloud_cas_certificate.default.name}_fake",
		}),
	}

	var existCasCertificatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"names.#":                      "1",
			"names.0":                      fmt.Sprintf("tf_testAccCasDataSource_%d", rand),
			"certificates.#":               "1",
			"certificates.0.name":          fmt.Sprintf("tf_testAccCasDataSource_%d", rand),
			"certificates.0.org_name":      "Internet Widgits Pty Ltd",
			"certificates.0.province":      "Some-State",
			"certificates.0.city":          "",
			"certificates.0.country":       "AU",
			"certificates.0.common":        "",
			"certificates.0.finger_print":  CHECKSET,
			"certificates.0.issuer":        CHECKSET,
			"certificates.0.start_date":    CHECKSET,
			"certificates.0.end_date":      CHECKSET,
			"certificates.0.sans":          "",
			"certificates.0.expired":       CHECKSET,
			"certificates.0.buy_in_aliyun": CHECKSET,
		}
	}

	var fakeCasCertificatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"certificates.#": "0",
		}
	}

	var casCertificatesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCasCertificatesMapFunc,
		fakeMapFunc:  fakeCasCertificatesMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CasClassicSupportedRegions)
	}
	casCertificatesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCasCertificatesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cas_certificate" "default" {
  name = "%s"
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIICWDCCAcGgAwIBAgIJAJYuytFmHxV0MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjAxMDIwMDY1MTMzWhcNMjAxMTE5MDY1MTMzWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDHXfC9Yk9WqrPkkXp3nF/n2mz88dGre4o6aRsfbXANSqnmHgPeoL6G3o3CslfS
TLV/UNRZypkB9W2Ie5SVoRt5n2Ow5keB/kEWutOGvAOfuyyBuXgGFPQZYjuBVZ1u
tpedtgWfJpyTxwrpV6vHKvIiISW0fjrO+UWu4DpLjQ8SSwIDAQABo1AwTjAdBgNV
HQ4EFgQUfKP4kMb1+5nyVzXqwyWvr9Qxh34wHwYDVR0jBBgwFoAUfKP4kMb1+5ny
VzXqwyWvr9Qxh34wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQAG87o5
ME68YN4qfTCvzY0TmWYKzBXzNEIDvgmtDPZZZt0BuQ9HgKGtuEnTEGnltsUBga2f
VQNKO3DPTyoFe7gXpDoR3Z5L+GdhwkEMkcLDLnn6E4hVJq6mtwbI6fY8Wc+kdh3Z
ivxUXMEAwR2Z64LvZu4bw+U/yofaBpozQwedOQ==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDHXfC9Yk9WqrPkkXp3nF/n2mz88dGre4o6aRsfbXANSqnmHgPe
oL6G3o3CslfSTLV/UNRZypkB9W2Ie5SVoRt5n2Ow5keB/kEWutOGvAOfuyyBuXgG
FPQZYjuBVZ1utpedtgWfJpyTxwrpV6vHKvIiISW0fjrO+UWu4DpLjQ8SSwIDAQAB
AoGAf+QSP9rl3EnK9rAgKMSmfTwQOD8D6oZoiMnN/V4dyFkCHj1Y7CKftjLkK2Zu
kdhlgZOfdS5S8v+20Ru9mDLuRPb7AHRLBaS25+eSXsltsMGRwiCvyphptlVll0O+
vSkU4eAWvAsg6W2KlT/0/WGjD+UXt6v96iBsXU6oLHZu67kCQQD2Y+zPrY6GGNBw
PktuyzU+mtEhkqWleju87QGuDzU7ZE8iz1QIlcn97mDdagNd5JSEuod4mRC4XKFq
QSzqP2alAkEAzySDIDtS5XNh2stEDyx0mW3haiKtCAwpoiQCdnIE/h0ngyzdzApg
pSRVs6WjrqcJYxvMGbkQf+UgaqsglLMyLwJBAJWrCXT2Fnd6p1MnZCb/JW7MGHFu
ZVTptVQEHFshPdLAEhoxGvjEFJk9rnWRKk5kxZsCu7wULsXu7tZelOwOa0kCQERv
d7LqZpTP7gBvL2kj8tHN7681DZ2fBxI+e2HOgb/Cug8of46tzwhAXOAhVVaacQuB
X4kQD1dxx6f2Kal3GpkCQQDGl0jzWlzTgOj5Ac4rSVOb/E1EABZA+298elqrEgQQ
ukV0sc58CvXgvu3dP9XzCQnBFeNYZxO71dqOF8e8fp9r
-----END RSA PRIVATE KEY-----
EOF
}
`, name)
}
