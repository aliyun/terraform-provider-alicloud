package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostAccountShareKeyAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_account_share_key_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostAccountShareKeyAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAccountShareKeyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostaccountsharekeyattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostAccountShareKeyAttachmentBasicDependence0)
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
					"instance_id":       "${data.alicloud_bastionhost_instances.default.instances.0.id}",
					"host_share_key_id": "${alicloud_bastionhost_host_share_key.default.host_share_key_id}",
					"host_account_id":   "${alicloud_bastionhost_host_account.default.host_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_id":   CHECKSET,
						"host_share_key_id": CHECKSET,
						"instance_id":       CHECKSET,
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

var AlicloudBastionhostHostAccountShareKeyAttachmentMap0 = map[string]string{}

func AlicloudBastionhostHostAccountShareKeyAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_host_share_key" "default" {
  host_share_key_name = var.name
  instance_id         = data.alicloud_bastionhost_instances.default.instances.0.id
  pass_phrase         = "NTIxeGlubXU="
  private_key         = "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJEZGNBVEl1cQpla0lyNXIzMUY4Z0NEc0FBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FDNzNHK3JTTWRsCjhmbTkycTlzTUxRWFZ6c2loRDlIblV1SnFib0JMbGhjY0F4Nk5SbnVoSTFGaHF4WXFnZ0RvVUFXdWltd2pQOHdvczZWNk4Kc3krc0pTbExFckZUdkZsWTBiL1JlN2tmOWxVcFkrL0o4Y2xNYUJ0MjRpNWdxelVuS1dFQ0diS2s4QlhuVTJUOTRNR291ZQphYmNsNWwwTGVvMzk1R1ZNa2VvZ1M2YzZkYmppS2hyYUhqcUxUcDg4TG5RNWZ3QW9jbnViZTIxN3pxRk1YaFZoaFl2QmlNCnRqS0dJYm5ZNzhoR0daSjd0VjdUUFlvL2EzQ3h5dXJ1ODhHeHRJd1h0cXJXWFFpZnp2Tzd6SUhHSDd1bmo3K212aURqRWwKa0VhbHpIa2VvY1AvMnFMQmxIcTlKaWRIa1FPVG5KWVE4M3dCRUpTM3R0OFBjd1VKc3E0aHhVbks5eHBqb2ZLcEtBaWRyNgpjbHZHLzltVTFSZmpYQlFFanBya1N3UW43d1ZCalRpZWZlM3dqMUdwREJtcTNXYTRZRjRXT29JSUNabmtIa0c0R1FhOG9HCjFES1Z3N29HUFRZenZ0ajJKYmY1dFNjSkZGWTVndTJvUjZKZk9HbzdXYlFUazZVYW9OaFdSNkZ3TUs3RlR1cjJDdFBhY1gKMWhQc2xjeXJ0aThyVUFBQVdRTXBaTW1XVVBrQUdmRGdxUGw0anNQdWV0cWFkdFpPMzN3Vm4zUXJsdFhjUUpEeUtBdVY3UQo5K0ZhekU2TjZ1bXBiT3AvUGM5bGFJN2paSW11V3R3RTBJMUFxZG9wT1dQWmVFQjFGN2dDbEw5UHJHZUtQN01PdFd1U2s4CkpaNFBNRXIwVG1IT3dzN3NvZDZpYjMxWWVlaXBsdUFOcFM5UFJPV1BiM1VwQ09TazdDalByQ3dwdTdHOGtIUUVCOGt4ZFIKQml5VDByclplV0J3TExHYkgxTjd0aW5rYVJoUkNtRDhadTBHV0k3SXFINHMvWGlBTmszR0dTTXdJVVBkMnJwdGVqUFdOdwpnR2pYRnUwbHB2TktRYU5GZmJQUFJEK1BRcUs5S3kraFdkN2o2MWVyR1ZUbmJZZ1lDaloxazRSUW16SDJDeGQwQkM1ajNRCktZa2hYWGg0TnQyRTVBa2RzbEh4TGNVdkx6N0tyT3V5RnZtSm83RHhUeUJyT0RZQlV6K3JHM0xpQXg0MmJydVFBMVFGUXMKc3EzeFNUYk4vYzVqajRvUXp3WHJvQ0p5bEJUY3Q5RVhzckdtcmNlQmI4OXNQWElrZjhCaS9kQlRGVDl1VUZqc1hvVW50SwpiSU9nWGtUUXJmZEViNWFmUitnY21XWlkxSmFjZEF0V2FTUGFCb2JGQ3FMamUwc2tqNHRnTkxnczVCQjZwYW1ic1FVRmpGClIvVzk0K2R6bitqM0VpWUIxMDFOYnJ4T1dCMmRDUE45clZxN0pCR2NGK3FyUjYxQzJvWFR0TWxNK1NFOGpCc1QwZEVYeUEKbjB6SVd2dEhidTdZekFidzRHSUszRGVmOU9vU3pJWHZFZmZrcE1FQjBlMHBnaEw0ZmlJYTBzZkdEbTA1Wk05MlNTNGpCcgp4KytBVHRqcXg1VnNJRkg5VzFoa252Szg0REZlSjZsdUxtVnAremdzdE04R0s5RGJqYmxtRGdQNDYwSHRoOU5kVFMrOTE3CkY4ajRYQ3FJTVJsbG1TVEd6VUFURDh3OEZkNExUdytSYlRBMjZJZnJHNFlZWFVNakJlNWxQdHFybGFLL0JEc2NJdWxRUWYKcmVtdmk3dXFrWUI4UHpsNFAreHRxRE9KUzljY2lxL0R3V3ppTGU3ekpnQ1d4ZFNvTFVLaVFkSjZDSUliMDVicWhXR1FDSApDZ2JKRW5hM0luY3p4cHl1ZmkrZUVDbklMUFk3c01oT1RNTmdNRUdkbGI1dnZyUEEvY2RmTDNNRDBSUERCQnh0allxRUlSCnJQZ1BiMGRjeUNmNmJZVTJhMUM2ZnNGeE0rbEdiY3doMWVHbWVCbFRsUy95YmFuWk92NXdtckQyODVlbXpWeGJ2RU1PTlYKb0tsYnYzK0JlNHRJV0w3VHRsR0VLNGs2L0hIdUI0SkdMV0dxNmZ6TUp4V0lreENLUTFXOXNhcXVQckFVZzVVdGFQcG9KSApqb2t4ellRWVZpK1Bzd3psOExhMTlHNVRialJ5dmd4aUhUYlB1K2dFaGdDT3pmZnA4TS85UDhYanNGUi9sU1RsbnQzUGlwCjloNVdDUllEeStTUWxiZ0E2OXRMQ1pBdUl0aHJ5bk5QbmZRZHpMdytHVWhhTDQ2ai9yeGdJV2gwWmc3ZWJMZHJvVExad2YKVjdiWVVFeXdZbkMwQjlpODNKNGZxSkwxOHdxeWkyL2xQOXIxR0ZRKytvSTJaMjk1THN3N0NaWTZKaDhJLytQSVVoeUNlZgpoNHhNd2FEUGFKeUx4MEdVOW9XdE5oSXRuZXUyNUtOd1Q4ODVkL01xTmRKSklxdE0vN04vT1pKWDZnWGlXTE1TS2tEUE1lCkNyZmtSZ1daOUtwZnc0WUxVdWkzd2FTZWEvbGpuVjBsTW9Eb2VsWUh6Vk1sRW5FTWZsZzRsYi9Nb0hXSEV6ZFlJeElWMXoKTTZFWjZQMzNacTJmY1ZHNUNBZ2dYTE0yQkRlYXdvODRrU1lmckJDeHBISFJMTmFybzJyUmNQUk8yWWhPaHp1WDY5WnJqbQprUFNtNXZzMVB6bjBMblpRUXhJNzFNWkhOSzlobVh1UUxaOHZaalBwOHBSR0U5Zi85b0NhYVBTdzZNdGp6VXZlN0RpNFVrCkR1bXhiT1BaMkF1aUtnaEE4TnVENTF5SzI5NzJTN0ttQ2hDbEY1UGdkakVVSk1SY3huNkp2ZS83WFZleVhCWnFGQmFkNFkKcXhHbS9mV2c3MHovRnExTlZsaUpnTThGQ1FvdGY3M0E4RDZneUYzK2xPTStvZFpQYTFhVFhwRWlEQnYwNEpXQld3ajhMYwpwUndXMzBwOVJ3ODBsQW1yMlZlMkFUaVk3VGlOenNQTEFlK2UxcFFFWGwxTlpRTnRQbFlnRXpyQTBmWTNCUG5NUHF3dUpVClhiS25WK3pUQTdHS1VUS3FtOExyYkJxcy9PVmhDQlF3VHhRM1lVcktTN1psaHBjcjh0STJjY2UzVEFoTFNJZUxYaVF2S3oKMUhPbm5XbDg4NUtpVXkyZmUvQkQ0LzMybEpBPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0="
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = data.alicloud_bastionhost_instances.default.ids.0
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}

resource "alicloud_bastionhost_host_account" "default" {
  instance_id       =  data.alicloud_bastionhost_instances.default.ids.0
  host_account_name = var.name
  host_id           = alicloud_bastionhost_host.default.host_id
  protocol_name     = "SSH"
  password          = "YourPassword12345"
}
`, name)
}
