package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostHostShareKeysDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bastionhost_host_share_key.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_bastionhost_host_share_key.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bastionhost_host_share_key.default.host_share_key_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_bastionhost_host_share_key.default.host_share_key_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_bastionhost_host_share_key.default.id}"]`,
			"name_regex": `"${alicloud_bastionhost_host_share_key.default.host_share_key_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_bastionhost_host_share_key.default.id}_fake"]`,
			"name_regex": `"${alicloud_bastionhost_host_share_key.default.host_share_key_name}_fake"`,
		}),
	}
	var existAlicloudBastionhostHostShareKeysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"keys.#":                          "1",
			"keys.0.host_share_key_name":      fmt.Sprintf("tf-testAccHostShareKey-%d", rand),
			"keys.0.id":                       CHECKSET,
			"keys.0.host_share_key_id":        CHECKSET,
			"keys.0.instance_id":              CHECKSET,
			"keys.0.private_key_finger_print": CHECKSET,
		}
	}
	var fakeAlicloudBastionhostHostShareKeysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudBastionhostHostShareKeysCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_bastionhost_host_share_keys.default",
		existMapFunc: existAlicloudBastionhostHostShareKeysDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudBastionhostHostShareKeysDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudBastionhostHostShareKeysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudBastionhostHostShareKeysDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccHostShareKey-%d"
}

data "alicloud_bastionhost_instances" "default" {
}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = var.name
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "NTIxeGlubXU="
  private_key         = "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJEZGNBVEl1cQpla0lyNXIzMUY4Z0NEc0FBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FDNzNHK3JTTWRsCjhmbTkycTlzTUxRWFZ6c2loRDlIblV1SnFib0JMbGhjY0F4Nk5SbnVoSTFGaHF4WXFnZ0RvVUFXdWltd2pQOHdvczZWNk4Kc3krc0pTbExFckZUdkZsWTBiL1JlN2tmOWxVcFkrL0o4Y2xNYUJ0MjRpNWdxelVuS1dFQ0diS2s4QlhuVTJUOTRNR291ZQphYmNsNWwwTGVvMzk1R1ZNa2VvZ1M2YzZkYmppS2hyYUhqcUxUcDg4TG5RNWZ3QW9jbnViZTIxN3pxRk1YaFZoaFl2QmlNCnRqS0dJYm5ZNzhoR0daSjd0VjdUUFlvL2EzQ3h5dXJ1ODhHeHRJd1h0cXJXWFFpZnp2Tzd6SUhHSDd1bmo3K212aURqRWwKa0VhbHpIa2VvY1AvMnFMQmxIcTlKaWRIa1FPVG5KWVE4M3dCRUpTM3R0OFBjd1VKc3E0aHhVbks5eHBqb2ZLcEtBaWRyNgpjbHZHLzltVTFSZmpYQlFFanBya1N3UW43d1ZCalRpZWZlM3dqMUdwREJtcTNXYTRZRjRXT29JSUNabmtIa0c0R1FhOG9HCjFES1Z3N29HUFRZenZ0ajJKYmY1dFNjSkZGWTVndTJvUjZKZk9HbzdXYlFUazZVYW9OaFdSNkZ3TUs3RlR1cjJDdFBhY1gKMWhQc2xjeXJ0aThyVUFBQVdRTXBaTW1XVVBrQUdmRGdxUGw0anNQdWV0cWFkdFpPMzN3Vm4zUXJsdFhjUUpEeUtBdVY3UQo5K0ZhekU2TjZ1bXBiT3AvUGM5bGFJN2paSW11V3R3RTBJMUFxZG9wT1dQWmVFQjFGN2dDbEw5UHJHZUtQN01PdFd1U2s4CkpaNFBNRXIwVG1IT3dzN3NvZDZpYjMxWWVlaXBsdUFOcFM5UFJPV1BiM1VwQ09TazdDalByQ3dwdTdHOGtIUUVCOGt4ZFIKQml5VDByclplV0J3TExHYkgxTjd0aW5rYVJoUkNtRDhadTBHV0k3SXFINHMvWGlBTmszR0dTTXdJVVBkMnJwdGVqUFdOdwpnR2pYRnUwbHB2TktRYU5GZmJQUFJEK1BRcUs5S3kraFdkN2o2MWVyR1ZUbmJZZ1lDaloxazRSUW16SDJDeGQwQkM1ajNRCktZa2hYWGg0TnQyRTVBa2RzbEh4TGNVdkx6N0tyT3V5RnZtSm83RHhUeUJyT0RZQlV6K3JHM0xpQXg0MmJydVFBMVFGUXMKc3EzeFNUYk4vYzVqajRvUXp3WHJvQ0p5bEJUY3Q5RVhzckdtcmNlQmI4OXNQWElrZjhCaS9kQlRGVDl1VUZqc1hvVW50SwpiSU9nWGtUUXJmZEViNWFmUitnY21XWlkxSmFjZEF0V2FTUGFCb2JGQ3FMamUwc2tqNHRnTkxnczVCQjZwYW1ic1FVRmpGClIvVzk0K2R6bitqM0VpWUIxMDFOYnJ4T1dCMmRDUE45clZxN0pCR2NGK3FyUjYxQzJvWFR0TWxNK1NFOGpCc1QwZEVYeUEKbjB6SVd2dEhidTdZekFidzRHSUszRGVmOU9vU3pJWHZFZmZrcE1FQjBlMHBnaEw0ZmlJYTBzZkdEbTA1Wk05MlNTNGpCcgp4KytBVHRqcXg1VnNJRkg5VzFoa252Szg0REZlSjZsdUxtVnAremdzdE04R0s5RGJqYmxtRGdQNDYwSHRoOU5kVFMrOTE3CkY4ajRYQ3FJTVJsbG1TVEd6VUFURDh3OEZkNExUdytSYlRBMjZJZnJHNFlZWFVNakJlNWxQdHFybGFLL0JEc2NJdWxRUWYKcmVtdmk3dXFrWUI4UHpsNFAreHRxRE9KUzljY2lxL0R3V3ppTGU3ekpnQ1d4ZFNvTFVLaVFkSjZDSUliMDVicWhXR1FDSApDZ2JKRW5hM0luY3p4cHl1ZmkrZUVDbklMUFk3c01oT1RNTmdNRUdkbGI1dnZyUEEvY2RmTDNNRDBSUERCQnh0allxRUlSCnJQZ1BiMGRjeUNmNmJZVTJhMUM2ZnNGeE0rbEdiY3doMWVHbWVCbFRsUy95YmFuWk92NXdtckQyODVlbXpWeGJ2RU1PTlYKb0tsYnYzK0JlNHRJV0w3VHRsR0VLNGs2L0hIdUI0SkdMV0dxNmZ6TUp4V0lreENLUTFXOXNhcXVQckFVZzVVdGFQcG9KSApqb2t4ellRWVZpK1Bzd3psOExhMTlHNVRialJ5dmd4aUhUYlB1K2dFaGdDT3pmZnA4TS85UDhYanNGUi9sU1RsbnQzUGlwCjloNVdDUllEeStTUWxiZ0E2OXRMQ1pBdUl0aHJ5bk5QbmZRZHpMdytHVWhhTDQ2ai9yeGdJV2gwWmc3ZWJMZHJvVExad2YKVjdiWVVFeXdZbkMwQjlpODNKNGZxSkwxOHdxeWkyL2xQOXIxR0ZRKytvSTJaMjk1THN3N0NaWTZKaDhJLytQSVVoeUNlZgpoNHhNd2FEUGFKeUx4MEdVOW9XdE5oSXRuZXUyNUtOd1Q4ODVkL01xTmRKSklxdE0vN04vT1pKWDZnWGlXTE1TS2tEUE1lCkNyZmtSZ1daOUtwZnc0WUxVdWkzd2FTZWEvbGpuVjBsTW9Eb2VsWUh6Vk1sRW5FTWZsZzRsYi9Nb0hXSEV6ZFlJeElWMXoKTTZFWjZQMzNacTJmY1ZHNUNBZ2dYTE0yQkRlYXdvODRrU1lmckJDeHBISFJMTmFybzJyUmNQUk8yWWhPaHp1WDY5WnJqbQprUFNtNXZzMVB6bjBMblpRUXhJNzFNWkhOSzlobVh1UUxaOHZaalBwOHBSR0U5Zi85b0NhYVBTdzZNdGp6VXZlN0RpNFVrCkR1bXhiT1BaMkF1aUtnaEE4TnVENTF5SzI5NzJTN0ttQ2hDbEY1UGdkakVVSk1SY3huNkp2ZS83WFZleVhCWnFGQmFkNFkKcXhHbS9mV2c3MHovRnExTlZsaUpnTThGQ1FvdGY3M0E4RDZneUYzK2xPTStvZFpQYTFhVFhwRWlEQnYwNEpXQld3ajhMYwpwUndXMzBwOVJ3ODBsQW1yMlZlMkFUaVk3VGlOenNQTEFlK2UxcFFFWGwxTlpRTnRQbFlnRXpyQTBmWTNCUG5NUHF3dUpVClhiS25WK3pUQTdHS1VUS3FtOExyYkJxcy9PVmhDQlF3VHhRM1lVcktTN1psaHBjcjh0STJjY2UzVEFoTFNJZUxYaVF2S3oKMUhPbm5XbDg4NUtpVXkyZmUvQkQ0LzMybEpBPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0="
}

data "alicloud_bastionhost_host_share_keys" "default" {	
   enable_details = true
   instance_id = data.alicloud_bastionhost_instances.default.instances.0.id
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
