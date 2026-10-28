package scdn_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func SkipTestAccAlicloudScdnDomainDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	testAccConfig := acctest.DataSourceTestAccConfigFunc("data.alicloud_scdn_domains.default", strconv.FormatInt(int64(rand), 10), dataSourceScdnDomainsConfigDependence)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_scdn_domain.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_scdn_domain.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_scdn_domain.default.id}"},
			"enable_details": "true",
			"status":         "online",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_scdn_domain.default.id}"},
			"enable_details": "true",
			"status":         "offline",
		}),
	}
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_scdn_domain.default.id}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_scdn_domain.default.id}-fake",
			"enable_details": "true",
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_scdn_domain.default.id}"},
			"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
			"enable_details":    "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_scdn_domain.default.id}"},
			"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID") + "fake",
			"enable_details":    "true",
		}),
	}

	var existScdnDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "1",
			"ids.0":     acctest.CHECKSET,
			"names.#":   "1",
			"names.0":   fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", scdnTestRegion(), rand),
			"domains.#": "1",
			//"domains.0.cert_infos.#":      "1",
			"domains.0.cname":             acctest.CHECKSET,
			"domains.0.create_time":       acctest.CHECKSET,
			"domains.0.description":       "",
			"domains.0.id":                fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", scdnTestRegion(), rand),
			"domains.0.domain_name":       fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", scdnTestRegion(), rand),
			"domains.0.gmt_modified":      acctest.CHECKSET,
			"domains.0.resource_group_id": acctest.CHECKSET,
			"domains.0.sources.#":         "1",
			"domains.0.status":            "online",
		}
	}

	var fakeScdnDomainsMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"domains.#": "0",
		}
	}

	var ScdnDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_scdn_domains.default",
		ExistMapFunc: existScdnDomainsMapCheck,
		FakeMapFunc:  fakeScdnDomainsMapCheck,
	}

	ScdnDomainsCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf, nameConf, resourceGroupIdConf)
}

func dataSourceScdnDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc%s%s.alicloud-provider.cn"
}
resource "alicloud_scdn_domain" "default" {
  domain_name = var.name
  resource_group_id = "%s"
  sources {
    content  = "xxx.aliyuncs.com"
    enabled  = "online"
    port     = 80
    priority = "20"
    type     = "oss"
  }
  cert_infos {
	cert_name = var.name
    cert_type ="upload"
    ssl_protocol = "on"
    ssl_pri = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n"
    ssl_pub = "-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n"
  }
}`, scdnTestRegion(), name, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}
