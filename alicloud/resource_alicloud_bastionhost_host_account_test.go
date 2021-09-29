package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_account.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostAccountBasicDependence0)
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
					"protocol_name":     "SSH",
					"host_id":           "${alicloud_bastionhost_host.default.host_id}",
					"instance_id":       "${alicloud_bastionhost_host.default.instance_id}",
					"host_account_name": "tf-testAcc-sYQ45HFBO7j9ACfiBBxxOj5M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_name":     "SSH",
						"host_id":           CHECKSET,
						"instance_id":       CHECKSET,
						"host_account_name": "tf-testAcc-sYQ45HFBO7j9ACfiBBxxOj5M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_name": "tf-testAcc-UyFi4eO0cKhGXvIXPLJmXiQj",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_name": "tf-testAcc-UyFi4eO0cKhGXvIXPLJmXiQj",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_key": "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUEzeTlYYkh2VHVsaCt2OUU3bk5BSXI3TEdTRmdNblU5VUxiYWFwYXFZVE0vMmpnbmdlNWhzClF1S3FWQWFVNXFaeFRORmIrek5BNzM4a3F4cmlrZGI0S1JxUVlxWTFudnBKNkpiK1RxQzVGMEQ1SXprd0pVa0N6WVdjdWcKU3NxTFdUYTQwSVNQdlQwcllwbEVraTl4RnRmSGpVd1hvVFEwb0pacW91VnZQRXUySUJjZ0FVSm1CUEI4NzVmVzZsWEdCSApONU43V1NEQmwyN3pjUlNjSXAvT3htemVWK052NjNtR0ZTdmxOTjM1aUZ0OTdCbXhIWHdQR3ZsU0laVmFNZFNTVS9pb3BWClgvUEJrZFRYdGJ3b3FEcUNULzNSbk1icDFHMFFXOEtqTUZ0K0NpWkhTV0RRdEhyeWNkVC82RHRYK1JneWM2b1NtN0d3Q3kKdlAzcmVrUHNqU01DclBmZkc4bWlTakt5bktsRDhnSTNHS2VhdUora2pyVm1hMTdLczdCckQwWGNabktidDFDY1VXR2VKbwpRb0QrdjlxZjVmdnMyNE0yVlFkQlBqVHNnWm90QnNrWmp5bFprSnl1M3lVRDRoWmg3UWJEQ294YVR2SDE1V0UwalU3dW5GClZWdURIdmw1VVcyYUdHSWltNnY5RzE4MWZDTldiaWh3a2RJQm01TlpBQUFGb0txUmU4NnFrWHZPQUFBQUIzTnphQzF5YzIKRUFBQUdCQU44dlYyeDcwN3BZZnIvUk81elFDSyt5eGtoWURKMVBWQzIybXFXcW1FelA5bzRKNEh1WWJFTGlxbFFHbE9hbQpjVXpSVy9zelFPOS9KS3NhNHBIVytDa2FrR0ttTlo3NlNlaVcvazZndVJkQStTTTVNQ1ZKQXMyRm5Mb0VyS2kxazJ1TkNFCmo3MDlLMktaUkpJdmNSYlh4NDFNRjZFME5LQ1dhcUxsYnp4THRpQVhJQUZDWmdUd2ZPK1gxdXBWeGdSemVUZTFrZ3daZHUKODNFVW5DS2Z6c1pzM2xmamIrdDVoaFVyNVRUZCtZaGJmZXdac1IxOER4cjVVaUdWV2pIVWtsUDRxS1ZWL3p3WkhVMTdXOApLS2c2Z2svOTBaekc2ZFJ0RUZ2Q296QmJmZ29tUjBsZzBMUjY4bkhVLytnN1Yva1lNbk9xRXB1eHNBc3J6OTYzcEQ3STBqCkFxejMzeHZKb2tveXNweXBRL0lDTnhpbm1yaWZwSTYxWm10ZXlyT3dhdzlGM0daeW03ZFFuRkZobmlhRUtBL3IvYW4rWDcKN051RE5sVUhRVDQwN0lHYUxRYkpHWThwV1pDY3J0OGxBK0lXWWUwR3d3cU1Xazd4OWVWaE5JMU83cHhWVmJneDc1ZVZGdAptaGhpSXB1ci9SdGZOWHdqVm00b2NKSFNBWnVUV1FBQUFBTUJBQUVBQUFHQVpQQThVY3dQRGhCSUF1aldWUzJoUUJWU3FCClZxWHhzcHJ5TU8vaTRSZzJ2cXpvS1pERXo3YWFTcDlDYWw0VXNWb3ZCczhVZFU3dnhKMFRqdmo1WHgxbVUxTitpRUI4cWEKOHA2WGxXZ0xUZ0VNckdtSTVOUUllSHNkVHVRZVVvOE1oVy9iZDJhdGZuYjBoeVFzdENFbHEwM2FxMFpTdi9RVUhHS0xZcgpnTkdkSlJaUVcwRjBjbmR4aWNyYVlGRTZwWGkwTWdYa2I2UjByZXZ4M2JINjIzRHZiZCtGSGNwRnJwMFZsdzZHQWNJeGQ5CjgvSUtzSk1USWRmdzFVOSs0MXFIbXhKb1NPdmhjRFkzcUZqM1N3clZ2MUJwVG9QazRENnNGMksrTDhscDNyMjdaU1hnQlQKSjRMekY0dVVDU1M3TVU5SGlLRk9zZXZLUXlKVjBhVzhZUmNJd3Zpc1BoK1BGeCs2a3Z1anVxb2xrOHk2M09YbmRzSUVTUQpvRDlPM2JWU0JCRC9IZ0E2NGNNQ3lacmY1dHg5ZGpmbE44OXl6bytFZlRUQWYrMDRIWXpWY0RodEhLd0RTcW16U2NGc3U1CjVGeTEwb1VWalFsaWZ2L1AzYUc0b0pwanRQalo1WldSdmpvbTJTQzVwaHJzRWNpaFByVGVVMzFFTUpneTFQczBTQkFBQUEKd0hpMENsQm41a1J1RkgvbkYvQTI0QnRXL1pwbXllZGljYzhpMFIyOU16V1ZoU2hvZ1RaeHM5MURVYkdweGhoL1dvZUJ1OAo3U3pyWkR5QkVscGRkN1JZMGdtanJzazUweW16NUxpTTVwa2VKZEdZWmZsUXZyQ0g4NUoyVDJaSStXVTlZRDdCYXZmVG9sCm41Tjh3eDlUa3M4Q0xBUjhGckNLeVYvb2M2WkZ0eFZxZ043QWo4S0FmcnE2T1lwVjZSUExoTjN4VFN0YXlCTG4rUEs1ejEKVy9tb0g5QTNxRWgrN1dUQ25qQUFFM0krTXBRekFOSmwrMitEdWhUSEgya3k4UDF3QUFBTUVBK2pkK0w3VFhKK3RjV1dEQgpoRjFKRlMxNlFCM1VqbUdxWW92R1FoU01RSEpuVVlYZGxYTTlFYncxU1FvSXFQSnBDSDF3SGZmczVRQWYrRnI5ODUzcDZICjhYaEhJVU54MTE5Qld2M0FHakdYOUloeDB1R1lBRnBGVERrZSt0WVp0QWUySVdMRDlPcTB6QUM5WkMvbmozWFVOZ0dvLzUKTHAyc0VncnF6QnFxbTNEZjlJTFY2MnA4THRZdk5NNkNkYUxmdmYzL0w4dkRqazQvS0ZxYkFqemFxdVYyWmJGTmJQNU9raQpYMTJVK3RvY3Y5NUppaGV2ZlNKSUExWUpuZnBEZEZBQUFBd1FEa1YrZG9aVy80NDR2REhBNS9HNzJoV3M2b0lnTEIzb0JsCk9rWnpXYjZQbmhpQVROcUJQdTR6d1h2d0VyYkhyMzQ3SGVKeEMralFuUVJIY1BsVmN0V2F4N1dlTzNrLzJYbUttVWpzbzcKZUd2TmR5czVaMDdOQWRTNDNmSzJpRnpKNDdhOHZtNEdEY09BNjlwekJOVVVqL0NFaVRNUE9za09TMGNvZ0doQ3c5U0VaWgowVDB4ODk2Yi9hd09HZjNOWWhvZXE4NXFlS1BjUkF3cVRXVUh3SXlidWNYcFFkQ2MzSlh0Nkx3QlZNLzNKL2RnNjVFbjdhCkJsbmI5MUF5NHc4d1VBQUFBa2FHVm5kV2x0YVc1QWFHVm5kV2x0YVc1a1pVMWhZMEp2YjJzdFVISnZMbXh2WTJGc0FRSUQKQkFVR0J3PT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_key": "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUEzeTlYYkh2VHVsaCt2OUU3bk5BSXI3TEdTRmdNblU5VUxiYWFwYXFZVE0vMmpnbmdlNWhzClF1S3FWQWFVNXFaeFRORmIrek5BNzM4a3F4cmlrZGI0S1JxUVlxWTFudnBKNkpiK1RxQzVGMEQ1SXprd0pVa0N6WVdjdWcKU3NxTFdUYTQwSVNQdlQwcllwbEVraTl4RnRmSGpVd1hvVFEwb0pacW91VnZQRXUySUJjZ0FVSm1CUEI4NzVmVzZsWEdCSApONU43V1NEQmwyN3pjUlNjSXAvT3htemVWK052NjNtR0ZTdmxOTjM1aUZ0OTdCbXhIWHdQR3ZsU0laVmFNZFNTVS9pb3BWClgvUEJrZFRYdGJ3b3FEcUNULzNSbk1icDFHMFFXOEtqTUZ0K0NpWkhTV0RRdEhyeWNkVC82RHRYK1JneWM2b1NtN0d3Q3kKdlAzcmVrUHNqU01DclBmZkc4bWlTakt5bktsRDhnSTNHS2VhdUora2pyVm1hMTdLczdCckQwWGNabktidDFDY1VXR2VKbwpRb0QrdjlxZjVmdnMyNE0yVlFkQlBqVHNnWm90QnNrWmp5bFprSnl1M3lVRDRoWmg3UWJEQ294YVR2SDE1V0UwalU3dW5GClZWdURIdmw1VVcyYUdHSWltNnY5RzE4MWZDTldiaWh3a2RJQm01TlpBQUFGb0txUmU4NnFrWHZPQUFBQUIzTnphQzF5YzIKRUFBQUdCQU44dlYyeDcwN3BZZnIvUk81elFDSyt5eGtoWURKMVBWQzIybXFXcW1FelA5bzRKNEh1WWJFTGlxbFFHbE9hbQpjVXpSVy9zelFPOS9KS3NhNHBIVytDa2FrR0ttTlo3NlNlaVcvazZndVJkQStTTTVNQ1ZKQXMyRm5Mb0VyS2kxazJ1TkNFCmo3MDlLMktaUkpJdmNSYlh4NDFNRjZFME5LQ1dhcUxsYnp4THRpQVhJQUZDWmdUd2ZPK1gxdXBWeGdSemVUZTFrZ3daZHUKODNFVW5DS2Z6c1pzM2xmamIrdDVoaFVyNVRUZCtZaGJmZXdac1IxOER4cjVVaUdWV2pIVWtsUDRxS1ZWL3p3WkhVMTdXOApLS2c2Z2svOTBaekc2ZFJ0RUZ2Q296QmJmZ29tUjBsZzBMUjY4bkhVLytnN1Yva1lNbk9xRXB1eHNBc3J6OTYzcEQ3STBqCkFxejMzeHZKb2tveXNweXBRL0lDTnhpbm1yaWZwSTYxWm10ZXlyT3dhdzlGM0daeW03ZFFuRkZobmlhRUtBL3IvYW4rWDcKN051RE5sVUhRVDQwN0lHYUxRYkpHWThwV1pDY3J0OGxBK0lXWWUwR3d3cU1Xazd4OWVWaE5JMU83cHhWVmJneDc1ZVZGdAptaGhpSXB1ci9SdGZOWHdqVm00b2NKSFNBWnVUV1FBQUFBTUJBQUVBQUFHQVpQQThVY3dQRGhCSUF1aldWUzJoUUJWU3FCClZxWHhzcHJ5TU8vaTRSZzJ2cXpvS1pERXo3YWFTcDlDYWw0VXNWb3ZCczhVZFU3dnhKMFRqdmo1WHgxbVUxTitpRUI4cWEKOHA2WGxXZ0xUZ0VNckdtSTVOUUllSHNkVHVRZVVvOE1oVy9iZDJhdGZuYjBoeVFzdENFbHEwM2FxMFpTdi9RVUhHS0xZcgpnTkdkSlJaUVcwRjBjbmR4aWNyYVlGRTZwWGkwTWdYa2I2UjByZXZ4M2JINjIzRHZiZCtGSGNwRnJwMFZsdzZHQWNJeGQ5CjgvSUtzSk1USWRmdzFVOSs0MXFIbXhKb1NPdmhjRFkzcUZqM1N3clZ2MUJwVG9QazRENnNGMksrTDhscDNyMjdaU1hnQlQKSjRMekY0dVVDU1M3TVU5SGlLRk9zZXZLUXlKVjBhVzhZUmNJd3Zpc1BoK1BGeCs2a3Z1anVxb2xrOHk2M09YbmRzSUVTUQpvRDlPM2JWU0JCRC9IZ0E2NGNNQ3lacmY1dHg5ZGpmbE44OXl6bytFZlRUQWYrMDRIWXpWY0RodEhLd0RTcW16U2NGc3U1CjVGeTEwb1VWalFsaWZ2L1AzYUc0b0pwanRQalo1WldSdmpvbTJTQzVwaHJzRWNpaFByVGVVMzFFTUpneTFQczBTQkFBQUEKd0hpMENsQm41a1J1RkgvbkYvQTI0QnRXL1pwbXllZGljYzhpMFIyOU16V1ZoU2hvZ1RaeHM5MURVYkdweGhoL1dvZUJ1OAo3U3pyWkR5QkVscGRkN1JZMGdtanJzazUweW16NUxpTTVwa2VKZEdZWmZsUXZyQ0g4NUoyVDJaSStXVTlZRDdCYXZmVG9sCm41Tjh3eDlUa3M4Q0xBUjhGckNLeVYvb2M2WkZ0eFZxZ043QWo4S0FmcnE2T1lwVjZSUExoTjN4VFN0YXlCTG4rUEs1ejEKVy9tb0g5QTNxRWgrN1dUQ25qQUFFM0krTXBRekFOSmwrMitEdWhUSEgya3k4UDF3QUFBTUVBK2pkK0w3VFhKK3RjV1dEQgpoRjFKRlMxNlFCM1VqbUdxWW92R1FoU01RSEpuVVlYZGxYTTlFYncxU1FvSXFQSnBDSDF3SGZmczVRQWYrRnI5ODUzcDZICjhYaEhJVU54MTE5Qld2M0FHakdYOUloeDB1R1lBRnBGVERrZSt0WVp0QWUySVdMRDlPcTB6QUM5WkMvbmozWFVOZ0dvLzUKTHAyc0VncnF6QnFxbTNEZjlJTFY2MnA4THRZdk5NNkNkYUxmdmYzL0w4dkRqazQvS0ZxYkFqemFxdVYyWmJGTmJQNU9raQpYMTJVK3RvY3Y5NUppaGV2ZlNKSUExWUpuZnBEZEZBQUFBd1FEa1YrZG9aVy80NDR2REhBNS9HNzJoV3M2b0lnTEIzb0JsCk9rWnpXYjZQbmhpQVROcUJQdTR6d1h2d0VyYkhyMzQ3SGVKeEMralFuUVJIY1BsVmN0V2F4N1dlTzNrLzJYbUttVWpzbzcKZUd2TmR5czVaMDdOQWRTNDNmSzJpRnpKNDdhOHZtNEdEY09BNjlwekJOVVVqL0NFaVRNUE9za09TMGNvZ0doQ3c5U0VaWgowVDB4ODk2Yi9hd09HZjNOWWhvZXE4NXFlS1BjUkF3cVRXVUh3SXlidWNYcFFkQ2MzSlh0Nkx3QlZNLzNKL2RnNjVFbjdhCkJsbmI5MUF5NHc4d1VBQUFBa2FHVm5kV2x0YVc1QWFHVm5kV2x0YVc1a1pVMWhZMEp2YjJzdFVISnZMbXh2WTJGc0FRSUQKQkFVR0J3PT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-sHK3VVqXUCIcusBVF3LTEdBU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-sHK3VVqXUCIcusBVF3LTEdBU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_name": "tf-testAcc-sHK3VVqXUCIcusBVF3LTEdBU",
					"private_key":       "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUFvMS83ZXdFQzU5QTJMZXhsTWwvQWRlcFA0Z1ZYWGVhN041azJUTFBQTmtXRkVqNTBWUXNjCkt3ay8wbjg1WVdKdnpZUml5NVh2R2tmVUZjd2dUNElyNDc1dlBkelF5U2EraEZDUUNIYVhRamlINzAwWXJpMEkyZU5va3cKT0djem1JTGpWdWM0eGxXRXBHQlB2V0V4eVIzM2dvOGZOQlJoWU16a0dCVmEwLzFWM2NxSFh5T0dzWUFJR1V2YkUwNGpVNwpVZklvTForYVAzOUgzRC9UMEpJVUF3QXpwTjRmT095MFdkZmJwWUV2SFdsdUhTVWE5eEZRM21uMk5xSC9rR1A2VFlNWS9WCkNITjJqZVl3Tk1PT3luNExPQWRTVFUrM1dhWjZaY0J6MG5KL2pQaUNKYjFueEErWnFDSTJjSVdKV21XaWVoaERkaXdZRW8KVW9CSGxnVFJLTnhORE02UXVDSjBCOXY2bi82VHpyU1VGR0dkcjFVSEdUNHAycWgva0ZvazVtbTBXYUJ5WHEybSswOVc4SQpDS25mYzVqcTNyMlhlc3VMd2FCVmo1eGNXZ2ZkbXBjZFZVK3BaVHdBdXN4NjhkeWhqVmJvSjNSbU5qZ05LcGpzb1ZuWGNRCk9KUnZHUnluL21MZCtmNllPOXBwcmg4UnRNdkRHeWdUTW9zSnZWNVZBQUFGb0tsbHhnT3BaY1lEQUFBQUIzTnphQzF5YzIKRUFBQUdCQUtOZiszc0JBdWZRTmkzc1pUSmZ3SFhxVCtJRlYxM211emVaTmt5enp6WkZoUkkrZEZVTEhDc0pQOUovT1dGaQpiODJFWXN1Vjd4cEgxQlhNSUUrQ0srTytiejNjME1rbXZvUlFrQWgybDBJNGgrOU5HSzR0Q05uamFKTURobk01aUM0MWJuCk9NWlZoS1JnVDcxaE1ja2Q5NEtQSHpRVVlXRE01QmdWV3RQOVZkM0toMThqaHJHQUNCbEwyeE5PSTFPMUh5S0MyZm1qOS8KUjl3LzA5Q1NGQU1BTTZUZUh6anN0Rm5YMjZXQkx4MXBiaDBsR3ZjUlVONXA5amFoLzVCaitrMkRHUDFRaHpkbzNtTURURApqc3ArQ3pnSFVrMVB0MW1tZW1YQWM5SnlmNHo0Z2lXOVo4UVBtYWdpTm5DRmlWcGxvbm9ZUTNZc0dCS0ZLQVI1WUUwU2pjClRRek9rTGdpZEFmYitwLytrODYwbEJSaG5hOVZCeGsrS2Rxb2Y1QmFKT1pwdEZtZ2NsNnRwdnRQVnZDQWlwMzNPWTZ0NjkKbDNyTGk4R2dWWStjWEZvSDNacVhIVlZQcVdVOEFMck1ldkhjb1kxVzZDZDBaalk0RFNxWTdLRloxM0VEaVVieGtjcC81aQozZm4rbUR2YWFhNGZFYlRMd3hzb0V6S0xDYjFlVlFBQUFBTUJBQUVBQUFHQWFlQ0NNYXp1SFIwcWY0aDc3TEZ4SVBuQTIxCkZxMVVmNmZJV21VdjhVZ3E5N0ZkK3p0SW1HcjcxR3h6djhDOGluZkNFWGhhaWRWQUxJeDNlS1dQeWJSUFRkVXRJUDNNeG4KRzRpNlQwSEx0UGE5NGErdEZ2UElrS3gzMFE2dnkyeTFmSHpVSDc0VXo2c1N4WmdQbkVNZnBodFJMYnZmeVhQd3lKcGJIeApNd0V5N0pHY09XUGtucFBDcStJbEQ5WEx2eXZhQ1p1VGQ5MXppOWNWZE1CaGNsTU95b2kyZ0lBL1FpelRhKytmdEdkV0VtCkkyMHdsSDE4VXllSWtNOWJ1RStSbGZ3TG1YZVRHK0d3aGhlU1RkSFVUd3pHWklQVXNPSUJONnRjMmQxd0pVR1haemc0V1MKd2RCZ0RaOHRXMWJZTUdRaFkvdjJ4M1FTci8wYlBBVjBvSWhPM09PZGVpTnZZcWYxeWExZEZSaGlYbm15eTJnZklyUGhvTwpqMFhXVGI3UFhCSm85WC8xYU1WL0gzdnlGYUFLcVR0eWJxWWp4YndabkpmY1dmVHlmL2JwV0hnalBoalR1T0wxWnk0ZnFSCmpkSEs0SE1PZ1lHaUcwZXhRTytKY25Cbk5Jcjk4UlNCWG9nRmFGUHpaVVloL3d4RFcrTTlsVnVEandNNzVJdFl6aEFBQUEKd0JOeEJucnNET0hHa0g5aDVsQ0wyTnFIaWUwVm9tL3ZLMjBkWmRFSVRmK3VWM3hwQTFVVEIwVkJUSjJmdXVQM0p4NnpuMQpJMWhLVitCZmh4N3NaWkVtbzk2WloxOHRkbDBXSkhjcVgxaU9nMHh4ZXAyd2p4bUlkOGhLc2FKWE1TdHJLeDVUUG5md296CklrT2g5Vm02VVFUQ1VxakdFRG9IVEtsc0ZwTjI3VVZ1WnJNV3pOTjZ4aHFockdQWmIwS2RkbEgvTEhMVTJmeFBUdldQRmMKSkVkUkRrMldBQmRGOUxEWmIxT2RMSXdHdm9paVhlOVk4U3NudnBZZXF6QnNXcmtnQUFBTUVBMC9Vby96dXNLSUtTazdXTgpiQXlVc0E3TWJnSUpYdDM2WFk5YkVYUE83TGc1ZnZlS2ZxdFo3TlU2SzYxYVNTcG9lb1JnWWhPd1R6MWM4Q3REMHB6SHY2CmpvQ3dZTTF4RXVha2tYeXV0T0ZpTmNzZHo1TnRYN3dYc0xwbHQwcy9HYlVITWpkTnlDdUpEWDhiWlg5aFlwdUhNZS9hY0gKWmdRam1oTXdNSER1ck1rZ1p6UnY4WWgya2RPemJxRTFPUXR3V0FGa21YRW5QZTBNZ1JvUjJRQmVUWExFZGVhaGdaUGtTbgptUzRNc296TUkvdzBMaG9Gc1FCQXRBand6SFFJRHpBQUFBd1FERlVvVllTQzZ5OXAyOHBGaVc4TDhLY2xRb29xcmozb3RMCmN0WHFsZTRWUHBTNm5vQUpKRzR5b3RveVhPY3VlRWxlQ3h2NFFsTWhNWU5kQWFGSG5JQlQ1Y29nWGg0TUIzbDNpSm9hR2UKNUF5czM0S3NBUGhHbm1qdEhFZ0I4OGh6a0tnQThpVmpmVWJ6ZWxGckl6MVM4Z05nL05BbGhaT0d2OW5zZ1JFOENGRmtwSApWck00RnlRSFVNb0Q5QXV4bUlCVG1Qd014QnMzL3cxS043eUZyUHF3dTNEU2w4N05jdmNxQWd0Mk1oeVcyWTZIRGI1UnV3CnVjNGdwQm5aL0ZOWmNBQUFBa2FHVm5kV2x0YVc1QWFHVm5kV2x0YVc1a1pVMWhZMEp2YjJzdFVISnZMbXh2WTJGc0FRSUQKQkFVR0J3PT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
					"password":          "tf-testAcc-pozAPSdUuQEYszUUdz0bffpO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_name": "tf-testAcc-sHK3VVqXUCIcusBVF3LTEdBU",
						"private_key":       "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUFvMS83ZXdFQzU5QTJMZXhsTWwvQWRlcFA0Z1ZYWGVhN041azJUTFBQTmtXRkVqNTBWUXNjCkt3ay8wbjg1WVdKdnpZUml5NVh2R2tmVUZjd2dUNElyNDc1dlBkelF5U2EraEZDUUNIYVhRamlINzAwWXJpMEkyZU5va3cKT0djem1JTGpWdWM0eGxXRXBHQlB2V0V4eVIzM2dvOGZOQlJoWU16a0dCVmEwLzFWM2NxSFh5T0dzWUFJR1V2YkUwNGpVNwpVZklvTForYVAzOUgzRC9UMEpJVUF3QXpwTjRmT095MFdkZmJwWUV2SFdsdUhTVWE5eEZRM21uMk5xSC9rR1A2VFlNWS9WCkNITjJqZVl3Tk1PT3luNExPQWRTVFUrM1dhWjZaY0J6MG5KL2pQaUNKYjFueEErWnFDSTJjSVdKV21XaWVoaERkaXdZRW8KVW9CSGxnVFJLTnhORE02UXVDSjBCOXY2bi82VHpyU1VGR0dkcjFVSEdUNHAycWgva0ZvazVtbTBXYUJ5WHEybSswOVc4SQpDS25mYzVqcTNyMlhlc3VMd2FCVmo1eGNXZ2ZkbXBjZFZVK3BaVHdBdXN4NjhkeWhqVmJvSjNSbU5qZ05LcGpzb1ZuWGNRCk9KUnZHUnluL21MZCtmNllPOXBwcmg4UnRNdkRHeWdUTW9zSnZWNVZBQUFGb0tsbHhnT3BaY1lEQUFBQUIzTnphQzF5YzIKRUFBQUdCQUtOZiszc0JBdWZRTmkzc1pUSmZ3SFhxVCtJRlYxM211emVaTmt5enp6WkZoUkkrZEZVTEhDc0pQOUovT1dGaQpiODJFWXN1Vjd4cEgxQlhNSUUrQ0srTytiejNjME1rbXZvUlFrQWgybDBJNGgrOU5HSzR0Q05uamFKTURobk01aUM0MWJuCk9NWlZoS1JnVDcxaE1ja2Q5NEtQSHpRVVlXRE01QmdWV3RQOVZkM0toMThqaHJHQUNCbEwyeE5PSTFPMUh5S0MyZm1qOS8KUjl3LzA5Q1NGQU1BTTZUZUh6anN0Rm5YMjZXQkx4MXBiaDBsR3ZjUlVONXA5amFoLzVCaitrMkRHUDFRaHpkbzNtTURURApqc3ArQ3pnSFVrMVB0MW1tZW1YQWM5SnlmNHo0Z2lXOVo4UVBtYWdpTm5DRmlWcGxvbm9ZUTNZc0dCS0ZLQVI1WUUwU2pjClRRek9rTGdpZEFmYitwLytrODYwbEJSaG5hOVZCeGsrS2Rxb2Y1QmFKT1pwdEZtZ2NsNnRwdnRQVnZDQWlwMzNPWTZ0NjkKbDNyTGk4R2dWWStjWEZvSDNacVhIVlZQcVdVOEFMck1ldkhjb1kxVzZDZDBaalk0RFNxWTdLRloxM0VEaVVieGtjcC81aQozZm4rbUR2YWFhNGZFYlRMd3hzb0V6S0xDYjFlVlFBQUFBTUJBQUVBQUFHQWFlQ0NNYXp1SFIwcWY0aDc3TEZ4SVBuQTIxCkZxMVVmNmZJV21VdjhVZ3E5N0ZkK3p0SW1HcjcxR3h6djhDOGluZkNFWGhhaWRWQUxJeDNlS1dQeWJSUFRkVXRJUDNNeG4KRzRpNlQwSEx0UGE5NGErdEZ2UElrS3gzMFE2dnkyeTFmSHpVSDc0VXo2c1N4WmdQbkVNZnBodFJMYnZmeVhQd3lKcGJIeApNd0V5N0pHY09XUGtucFBDcStJbEQ5WEx2eXZhQ1p1VGQ5MXppOWNWZE1CaGNsTU95b2kyZ0lBL1FpelRhKytmdEdkV0VtCkkyMHdsSDE4VXllSWtNOWJ1RStSbGZ3TG1YZVRHK0d3aGhlU1RkSFVUd3pHWklQVXNPSUJONnRjMmQxd0pVR1haemc0V1MKd2RCZ0RaOHRXMWJZTUdRaFkvdjJ4M1FTci8wYlBBVjBvSWhPM09PZGVpTnZZcWYxeWExZEZSaGlYbm15eTJnZklyUGhvTwpqMFhXVGI3UFhCSm85WC8xYU1WL0gzdnlGYUFLcVR0eWJxWWp4YndabkpmY1dmVHlmL2JwV0hnalBoalR1T0wxWnk0ZnFSCmpkSEs0SE1PZ1lHaUcwZXhRTytKY25Cbk5Jcjk4UlNCWG9nRmFGUHpaVVloL3d4RFcrTTlsVnVEandNNzVJdFl6aEFBQUEKd0JOeEJucnNET0hHa0g5aDVsQ0wyTnFIaWUwVm9tL3ZLMjBkWmRFSVRmK3VWM3hwQTFVVEIwVkJUSjJmdXVQM0p4NnpuMQpJMWhLVitCZmh4N3NaWkVtbzk2WloxOHRkbDBXSkhjcVgxaU9nMHh4ZXAyd2p4bUlkOGhLc2FKWE1TdHJLeDVUUG5md296CklrT2g5Vm02VVFUQ1VxakdFRG9IVEtsc0ZwTjI3VVZ1WnJNV3pOTjZ4aHFockdQWmIwS2RkbEgvTEhMVTJmeFBUdldQRmMKSkVkUkRrMldBQmRGOUxEWmIxT2RMSXdHdm9paVhlOVk4U3NudnBZZXF6QnNXcmtnQUFBTUVBMC9Vby96dXNLSUtTazdXTgpiQXlVc0E3TWJnSUpYdDM2WFk5YkVYUE83TGc1ZnZlS2ZxdFo3TlU2SzYxYVNTcG9lb1JnWWhPd1R6MWM4Q3REMHB6SHY2CmpvQ3dZTTF4RXVha2tYeXV0T0ZpTmNzZHo1TnRYN3dYc0xwbHQwcy9HYlVITWpkTnlDdUpEWDhiWlg5aFlwdUhNZS9hY0gKWmdRam1oTXdNSER1ck1rZ1p6UnY4WWgya2RPemJxRTFPUXR3V0FGa21YRW5QZTBNZ1JvUjJRQmVUWExFZGVhaGdaUGtTbgptUzRNc296TUkvdzBMaG9Gc1FCQXRBand6SFFJRHpBQUFBd1FERlVvVllTQzZ5OXAyOHBGaVc4TDhLY2xRb29xcmozb3RMCmN0WHFsZTRWUHBTNm5vQUpKRzR5b3RveVhPY3VlRWxlQ3h2NFFsTWhNWU5kQWFGSG5JQlQ1Y29nWGg0TUIzbDNpSm9hR2UKNUF5czM0S3NBUGhHbm1qdEhFZ0I4OGh6a0tnQThpVmpmVWJ6ZWxGckl6MVM4Z05nL05BbGhaT0d2OW5zZ1JFOENGRmtwSApWck00RnlRSFVNb0Q5QXV4bUlCVG1Qd014QnMzL3cxS043eUZyUHF3dTNEU2w4N05jdmNxQWd0Mk1oeVcyWTZIRGI1UnV3CnVjNGdwQm5aL0ZOWmNBQUFBa2FHVm5kV2x0YVc1QWFHVm5kV2x0YVc1a1pVMWhZMEp2YjJzdFVISnZMbXh2WTJGc0FRSUQKQkFVR0J3PT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
						"password":          "tf-testAcc-pozAPSdUuQEYszUUdz0bffpO",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"pass_phrase", "private_key", "password"},
			},
		},
	})
}

var AlicloudBastionhostHostAccountMap0 = map[string]string{
	"host_account_id": CHECKSET,
	"instance_id":     CHECKSET,
}

func AlicloudBastionhostHostAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
zone_id = local.zone_id
vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
vswitch_name = var.name
vpc_id       = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.ids.0
cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
name   = var.name
}
locals {
vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
description        = var.name
license_code       = "bhah_ent_50_asset"
period             = "1"
vswitch_id         = local.vswitch_id
security_group_ids = [alicloud_security_group.default.id]
}
resource "alicloud_bastionhost_host" "default" {
 instance_id          = alicloud_bastionhost_instance.default.id
 host_name            = var.name
 active_address_type  = "Private"
 host_private_address = "172.16.0.10"
 os_type              = "Linux"
 source               = "Local"
}
`, name)
}
func TestAccAlicloudBastionhostHostAccount_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_account.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostAccountMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostAccountBasicDependence1)
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
					"protocol_name":     "RDP",
					"host_id":           "${alicloud_bastionhost_host.default.host_id}",
					"instance_id":       "${alicloud_bastionhost_host.default.instance_id}",
					"host_account_name": "tf-testAcc-sYQ45HFBO7j9ACfiBBxxOj5M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_name":     "RDP",
						"host_id":           CHECKSET,
						"instance_id":       CHECKSET,
						"host_account_name": "tf-testAcc-sYQ45HFBO7j9ACfiBBxxOj5M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_name": "tf-testAcc-wXhJa78yERkMiiAGRz3qVNVL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_name": "tf-testAcc-wXhJa78yERkMiiAGRz3qVNVL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-sOYnEN3xJpcLfWhcSq5j7LaS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-sOYnEN3xJpcLfWhcSq5j7LaS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_name": "tf-testAcc-sOYnEN3xJpcLfWhcSq5j7LaS",
					"password":          "tf-testAcc-yICeFGQtCmVUS07CgTLcNMuf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_name": "tf-testAcc-sOYnEN3xJpcLfWhcSq5j7LaS",
						"password":          "tf-testAcc-yICeFGQtCmVUS07CgTLcNMuf",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password", "private_key", "pass_phrase"},
			},
		},
	})
}

var AlicloudBastionhostHostAccountMap1 = map[string]string{
	"host_account_id": CHECKSET,
	"instance_id":     CHECKSET,
}

func AlicloudBastionhostHostAccountBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
zone_id = local.zone_id
vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
vswitch_name = var.name
vpc_id       = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.ids.0
cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
name   = var.name
}
locals {
vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
description        = var.name
license_code       = "bhah_ent_50_asset"
period             = "1"
vswitch_id         = local.vswitch_id
security_group_ids = [alicloud_security_group.default.id]
}
resource "alicloud_bastionhost_host" "default" {
 instance_id          = alicloud_bastionhost_instance.default.id
 host_name            = var.name
 active_address_type  = "Private"
 host_private_address = "172.16.0.10"
 os_type              = "Linux"
 source               = "Local"
}
`, name)
}
