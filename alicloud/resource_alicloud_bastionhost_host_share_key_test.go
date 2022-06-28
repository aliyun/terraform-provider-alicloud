package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostShareKey_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_share_key.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostShareKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostShareKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostsharekey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostShareKeyBasicDependence0)
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
					"instance_id":         "${data.alicloud_bastionhost_instances.default.instances.0.id}",
					"host_share_key_name": name,
					"private_key":         "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBc25oc29SSVVwVXltSG1FVHJXUGxDbkhMa3c3N0JYTm44ZHcvbDg3eG10SUhjd2syCkRybjFDZk5jZkpJV0tSdkFaYkdKMlZTS1RiRDhPTmcyT3JvUHFGUHBLOHJ5QjJRb1NYQVRsaUVHWFhNeW1tRm8KeDBmem12THFscUxpNGZnOExhcTc5UC85aGxLU1djTWhJU0pYVTNHMS9KdEFBUmEyQXc4cXEzSVQvMkZ5NktrdwowMU9MdDdLN2pGUFRPaHhtdmNoSkZ1SVo1YXI0cW5HUlFHQnpCL2hoRHVIWEMwRlhJZ2ozd0NXMDZ4R2V2WjJyCmNCWWwwN1luL2lvZk95MU0wRjZZN0JrMU95N21vYndzM1JsalUyL2FpZlhLMmNOUlk2Qjl5WXNvd1RBZmQ5OTQKQ2YxSlF3TlhsaUZCeTZueEJLQk1YbDhIU1grS1o3L29PUlIwVXdJREFRQUJBb0lCQVFDbU5JSXR5ckhSY3oxdApJMGo0L0FQc295ZE1EL0owRkJMa2FoSUxKWjFaYW1tbmx4ZHh4WHBQUndXRnVXTEw2OTFVbDI5aUoxb1ptazU1Ci9ka2EvZlhnOUN3OUxXWVN2aExLdVlaMEZOTmhxZ3VoUEVBZy9uLytlR1ZCM2ZYZkxaZVZpK0E0L1VHMG15ZFMKVXVlQ2ZRSElZeWh4VkgvWnc3WER5WmNhVFVZVVdMUWlYcVN0Y2JRbnZFOXpwOGc5TWh5UkhBcWYwbEt2UTRqdwphUnNKTnlob3lhZWcvUXlFeHVYNGdBR1lIc1lTSDRFVmtXOHl5WE1aOHpRdk1OSUNiYXhmUkRuSngybUh6a09rCnFHczVXbFp5L3VrQk5jWTQwd2Y0eTY2bEVJaVpKbiswaFhtSTF4Tk5SdHRwMjZnY3ROOXZWbmVicTdLTnpjTDgKeFQrMXZJaEpBb0dCQU9iMVM1YlE4NVRFWDBoZTRmdXc2R3ExbnhRLzJUSU03emZhK2VhZThPQlh2eVNFV3JpdwpPZzM3RFhVUDFNVU1iTEJRenE0STl1dE5YSVZadEFLR0h6ZDR6WmtQeGxORjZPN0FyWnJEWElDNEdKZHdmSEhxCjJZcDkxUDlWekJlOVhkTVdZVGFCNkMzWVdtYzQwM08vYWdyRCtNb2JnL0hqMSt0d2xZR2hjdlV2QW9HQkFNWFMKT2VnWHc5VUF3VEZabFBtZzZKeDI3TzNXUFBHd1E3QzRnYitFZzZkR1pLRnJVR1ZId2VUUG1HaGtwN1BmYU5ESwplaFVoUWFnNm9XOTF4dkE2YldZZ29SQmczUWkxc01MblRWeTExeVN1UEVFSCtMT2s1N3d2akNLSk5XZnM0SjVyCmg1NGw0QXZ6UVhyWWN0UlZkSmYrNjFacGFnTkdZMVBvWVJMTHJMSWRBb0dBTndydzErRzJtNWJ0YW04S2hwU1QKMzVLbmRnajlkM3N6cStrcE03aGZpZWYvcXZGTU9jWHVJQlRjRVRFVHNWNlRyTFdsZkQ2d3NrVitybDFCbEhSbwpqaXpoT3dCU2NOZ3hlbTA3TXE0cXBwYTViYVltVW5QNUlwTjRwdDNJeFVPaFQ4UitxS0h2TnJYZ1hjZGlSYXl4CjFoejhkeFoxckxselpTNHd3M001MVlzQ2dZRUFpUDEwTEUySXg5Q2wrTTdZWTZZU2I0Zkx1MGhKRy9XOGFuemIKSFExZlBrOTVFRytJVlJyRUl2ZS95MHNvOTE4VzdyL0lteWxVbG5ORHFEUWZkK3grSmVNaXBuenRsRUorRGZxdgprQ3c4dUtJUUI5akZXV0l4T0JpVktyVnB6bll6ZG9Gd2dRd3BneDBKazFDZzlIblpMQWpVWUJyUDEwUy9ORFFRClJUdldjK0VDZ1lBeGRIZWxQNG1RdkJaS1oxMlNKbHlLbFVLeW43aHhzSHVkMkphMVNtS3FWeHBERDNlR0w0Y3QKZXA1QTZ5NkF4eGViZkI0aDdYNEZ0QTBBRURPdkZDR0J1QlRvZ3ZBdUNDVUtqK2JIUG1SNG53UVYzcWZ2M3loRAp0TGkwU2FHVElta2wvbUNCUDhZaW9JUys2N0xjby9kbHphUTNGVDlxTnJieFdFWjJlaS9LVlE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_share_key_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pass_phrase": "eGlubXU1MjE=",
					"private_key": "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJCTE9KUW1UTQoyZk5scGRmSWRMRHhyeEFBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FERjlMRVU2ZzdxCkdnUmttWFlScDUyU0RxMGIxTmh4Q1NYZzVnWlIvU3pFRXUyR2c5a3RJb2huOTlmTytYK2J2N0JnUm5ib0svYzlFT1dyckQKYm1CeHcvOTh3NVlSQTZQS2JhUllhejdIMUJhRDlJaGFuc0JXaDNLWFducFcyV1RUeXpDaVZLOEJCSHhFeWVXdy9xK29iMwpWVXdIRER5QzdWU0lRTklrbS9WbUtlajdvSTk2eEhpSENtWHh4dkUzNVpVcHFvdkdqbXVETUhhYnB4cnpJZHlZejFkTlN0Cm43dlc1bDY4Skt3V2ZPQmkyWFVuYjE1VHJIWTJSTHdRN3QzUysyR1pWclRBZVNiTVJETXM4YXI5cHgwd0RzMFlZdkE4QjgKejBzQ25ERU5sZWd1cHBSMS9jdTFmRys3NXJXRTY1bmZwUnUyNktmcGhaZ2o1ckJDOVFMN0pHbWJ0WmpicTZHRHZ3bmJ4cwpBU2ZqNjFyN0JJejVLcFFCZE5aR3hEeDFmL1B3Nmg2R3c4dzBab25EajVKZkpSTkdkcEpCSnFpcUpaL3FCRUZmeVQ0WldZClVuNlc3UndTRnVrVVllMDQvQjZiMTVmNTBkWDJsRmNDbDMwZGUycW9HZmg3ZkVJa3FHTEZqSllNWUJnZ1ZXVDZ4SEdUZHUKZXVkM2NydVg3YTJsOEFBQVdRWXpuaEFhZDgvN2lCNHVFUlJUSHN1eW9tZ2JmVHFJZWhVN01ZVUVOSGJ1QUs2cDVhd3lxaQovTFluaFZkemtxS1RvaXYyS3BKYU5UNk1GTmJoUW9vNWFMdW0vRERVc0tnR2xmN0xiWEd2clNyYjZzd0w4Z2VGY3h4MHh5ClJ6YVN6dWNMbHlVZmZqaHpoaGFxU2pacEsrbHBaeUZ0L01mODJpZ3RCYUo3VlhBR3c2OERhNHFpRGVPUCtueUpWS3gvWmcKZERKV05CQU9NeDJyaitaZHJWb3R1aU5iTy81ejZzSm5BWnlHS0srdDFyd2hxb3djb09lWU1DR3VWMnFwY0tLOEhrNEtidwp1OVZlMDBNc0tzbnhFYS9ZYTQvZ1c2aFhUeGJ2N3NVR1pYaUNhcExtMS8zN1VlV3JncFZ4Q0ZwZXdZaWlGOHlBRUFCa1VqCk1PeXhiNDl1SUdFanZlSk8xMnFNWlBZU2lqS3JEV2RLMUxIQzl1TUNGbVo4MTRCN08rbGlSVnFKZVoxVHdEYmU1c3BoeGcKamV3UkpyZnA5SjZndXovdTJXVkhqeUVKUEJlSnFPSWZKbXFGU3BaakFlMStGanVLRks0Wi9pbTdnb3owbVRGK0t0Nmk4Nwo1YWVSa0ZCT1Rjbkd6V0tiVGtQL0VkRUtkK0dORkhtZjJ0dVhBSzRYaVRLUU90ek9kNnFoYjcycHRBTWlYTWxOVm9keHVJClZYUGZubzNUeG9YOEg4b2g1UlZKcEttY3BkRENuMm5kV21EL2VURHFxcEhrampibVlwcXM0SHJmWjVtTjNhNTJGRTZEN0wKd1JQcFZrREtnVHdJRnRXK1VxTnRpNDNKYTBrWWZIcHZJREVDRER3aVpmOGdUaHI4SVBXZ1dwOWxxM1BzaHZrUlNMYUQ0dApPY3REUDgrSWVMeTlhbERKRGRhYW00MG5lWnNaTDBsZm5FTWF3RHV6TE5QQ1pETWt0YVp1MDhlRXRyY3VsdE50TWp0N1F2CkltVDJvcVA5M0FWVTlnRm01NVo3d3V6c01Nb0xmOFlhYzVoUDVUVkVaOGVXMUNyY0J5YkRrdEM3SDFBN1ZFTFZrQS9vWkgKYlBjcmpyOXBtRit3bTE1SVVsbU1TWXVRYUpiNUYxTExxRzNBZTZHRjNMVHlUTktkMXg5Nk9IeUNXOWRML216dHJSM0xPbApidzVyK0FpOGpZcldKaWMySHFPVkNwRU03ZHBsTlFLRzVHcnNGZ0RnK0tDQ1dVTlpzM3RnOGJCbDBrMVJ5dzE5UjAramRrCmFpbWp4eHhmYkVvUDV2aVhpMnJ4ZE1iRTlLRnhXSkw5ZFJsUFYvc1BkN0VOQmNLSVdRNFVSL28rS2tBVW5HVjFjdWp4aHkKcEVRNXFDdVBOeTNZR3RacHZoa0ExZWF1NVpyenBFVlJLNlJBWldwc3c2QUJYd1duQnA2NnJEUWdBU0gvRFdmbXBiQmlaaAphaHJFcXlyY0RlS2xWZlhHRHBIMmExeFoyRXRoNHJjSThwazlacG9uVlJjTTlFUVAzaWxUUWR2Z3RzMUhzeDhIVkJCMlk0Cng1TzJsWjVGbEdBTFJacXpQSzk1TW1vazMxZWUyN09RMngwWC9HUjZTa2lKbnVNT3JTSlRScDlYa2ZQRmdYWU9YWExFancKa0N5dStUTkxuTzJlQjlQQUVYZjB1K2gyM0UyM2xTdDFobHhXVjhXNzRXemtxdWtJd1MrTk1zWWQxcjA5cjNLV2I5YjhhRQpsTkxZbkltOGhzZ1NUVDZ6aVJJT05mVmVyRU83VFlRZW9tdU41Q05OcGt6UmFrQzRIUEtVM0J6V1VYR2lGV01sSHJRdlZNCnNJZkpRaFpNUWdaNnNsdHZaUWNtVENDbmFsM0NXeWpPM2RFYUt6TTYxUE4rYmQyeGZoRDVKVnBvUVpQWkNaa0VlUENSQlEKdnFjQVorV3ZyN1dIcmhSbmI0bG84cmlQWmFLZDVmWk9xRTJVa1FYanVaN3JDMXFYQml5dDBsUVBwM0tkVW1yK3Q1d1RPRwpEYzhIOHQvOVhEcmJnclQvMFRZRFBUMlg2VURyS2dqMkNkamVZVlhLK2htUWFIb2F6NUhpeFBodytRb3o4cHhFd2VNSHpXCnRZZERvZDF6UzZSZ2YwUXhQQitOQWl4OW5IcDRuWm1UdmVheEo0SlBKK0NiVXI3ek9BM1JxMDRXd2FlcDFBa3p3SG1TQ0EKTWRIYkkxQ2FGUEpLdlEvVjRqMzR0dm41TDR3azJWU2p1NUNjdStyY2RMMG9aS1BidmtvK2hONEVOd2Vvd1FJUWwzL2Ewcwozb0NHZ2tIcS9mOHl0a1FMNzBpZUdzYVRTaUxpanYwc29BbGpXNkIvUjdwY3BDeEZCR0UvTWtSMzRXaHBuK052MFRqaUpPCnUxUWgySmIyOFFkNlA2RmZwWVVVeXhnTllPYk94WkJybUdmTU10anI4eGJ3M0dvZjZ0SkpDQzFQUnVxUndNRm9QZVBoZHcKUE9LUlBUSUxCYmk2YUV1RjZ2eGxUajVGMGdNPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_share_key_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_share_key_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pass_phrase":         "NTIxeGlubXU=",
					"private_key":         "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJEZGNBVEl1cQpla0lyNXIzMUY4Z0NEc0FBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FDNzNHK3JTTWRsCjhmbTkycTlzTUxRWFZ6c2loRDlIblV1SnFib0JMbGhjY0F4Nk5SbnVoSTFGaHF4WXFnZ0RvVUFXdWltd2pQOHdvczZWNk4Kc3krc0pTbExFckZUdkZsWTBiL1JlN2tmOWxVcFkrL0o4Y2xNYUJ0MjRpNWdxelVuS1dFQ0diS2s4QlhuVTJUOTRNR291ZQphYmNsNWwwTGVvMzk1R1ZNa2VvZ1M2YzZkYmppS2hyYUhqcUxUcDg4TG5RNWZ3QW9jbnViZTIxN3pxRk1YaFZoaFl2QmlNCnRqS0dJYm5ZNzhoR0daSjd0VjdUUFlvL2EzQ3h5dXJ1ODhHeHRJd1h0cXJXWFFpZnp2Tzd6SUhHSDd1bmo3K212aURqRWwKa0VhbHpIa2VvY1AvMnFMQmxIcTlKaWRIa1FPVG5KWVE4M3dCRUpTM3R0OFBjd1VKc3E0aHhVbks5eHBqb2ZLcEtBaWRyNgpjbHZHLzltVTFSZmpYQlFFanBya1N3UW43d1ZCalRpZWZlM3dqMUdwREJtcTNXYTRZRjRXT29JSUNabmtIa0c0R1FhOG9HCjFES1Z3N29HUFRZenZ0ajJKYmY1dFNjSkZGWTVndTJvUjZKZk9HbzdXYlFUazZVYW9OaFdSNkZ3TUs3RlR1cjJDdFBhY1gKMWhQc2xjeXJ0aThyVUFBQVdRTXBaTW1XVVBrQUdmRGdxUGw0anNQdWV0cWFkdFpPMzN3Vm4zUXJsdFhjUUpEeUtBdVY3UQo5K0ZhekU2TjZ1bXBiT3AvUGM5bGFJN2paSW11V3R3RTBJMUFxZG9wT1dQWmVFQjFGN2dDbEw5UHJHZUtQN01PdFd1U2s4CkpaNFBNRXIwVG1IT3dzN3NvZDZpYjMxWWVlaXBsdUFOcFM5UFJPV1BiM1VwQ09TazdDalByQ3dwdTdHOGtIUUVCOGt4ZFIKQml5VDByclplV0J3TExHYkgxTjd0aW5rYVJoUkNtRDhadTBHV0k3SXFINHMvWGlBTmszR0dTTXdJVVBkMnJwdGVqUFdOdwpnR2pYRnUwbHB2TktRYU5GZmJQUFJEK1BRcUs5S3kraFdkN2o2MWVyR1ZUbmJZZ1lDaloxazRSUW16SDJDeGQwQkM1ajNRCktZa2hYWGg0TnQyRTVBa2RzbEh4TGNVdkx6N0tyT3V5RnZtSm83RHhUeUJyT0RZQlV6K3JHM0xpQXg0MmJydVFBMVFGUXMKc3EzeFNUYk4vYzVqajRvUXp3WHJvQ0p5bEJUY3Q5RVhzckdtcmNlQmI4OXNQWElrZjhCaS9kQlRGVDl1VUZqc1hvVW50SwpiSU9nWGtUUXJmZEViNWFmUitnY21XWlkxSmFjZEF0V2FTUGFCb2JGQ3FMamUwc2tqNHRnTkxnczVCQjZwYW1ic1FVRmpGClIvVzk0K2R6bitqM0VpWUIxMDFOYnJ4T1dCMmRDUE45clZxN0pCR2NGK3FyUjYxQzJvWFR0TWxNK1NFOGpCc1QwZEVYeUEKbjB6SVd2dEhidTdZekFidzRHSUszRGVmOU9vU3pJWHZFZmZrcE1FQjBlMHBnaEw0ZmlJYTBzZkdEbTA1Wk05MlNTNGpCcgp4KytBVHRqcXg1VnNJRkg5VzFoa252Szg0REZlSjZsdUxtVnAremdzdE04R0s5RGJqYmxtRGdQNDYwSHRoOU5kVFMrOTE3CkY4ajRYQ3FJTVJsbG1TVEd6VUFURDh3OEZkNExUdytSYlRBMjZJZnJHNFlZWFVNakJlNWxQdHFybGFLL0JEc2NJdWxRUWYKcmVtdmk3dXFrWUI4UHpsNFAreHRxRE9KUzljY2lxL0R3V3ppTGU3ekpnQ1d4ZFNvTFVLaVFkSjZDSUliMDVicWhXR1FDSApDZ2JKRW5hM0luY3p4cHl1ZmkrZUVDbklMUFk3c01oT1RNTmdNRUdkbGI1dnZyUEEvY2RmTDNNRDBSUERCQnh0allxRUlSCnJQZ1BiMGRjeUNmNmJZVTJhMUM2ZnNGeE0rbEdiY3doMWVHbWVCbFRsUy95YmFuWk92NXdtckQyODVlbXpWeGJ2RU1PTlYKb0tsYnYzK0JlNHRJV0w3VHRsR0VLNGs2L0hIdUI0SkdMV0dxNmZ6TUp4V0lreENLUTFXOXNhcXVQckFVZzVVdGFQcG9KSApqb2t4ellRWVZpK1Bzd3psOExhMTlHNVRialJ5dmd4aUhUYlB1K2dFaGdDT3pmZnA4TS85UDhYanNGUi9sU1RsbnQzUGlwCjloNVdDUllEeStTUWxiZ0E2OXRMQ1pBdUl0aHJ5bk5QbmZRZHpMdytHVWhhTDQ2ai9yeGdJV2gwWmc3ZWJMZHJvVExad2YKVjdiWVVFeXdZbkMwQjlpODNKNGZxSkwxOHdxeWkyL2xQOXIxR0ZRKytvSTJaMjk1THN3N0NaWTZKaDhJLytQSVVoeUNlZgpoNHhNd2FEUGFKeUx4MEdVOW9XdE5oSXRuZXUyNUtOd1Q4ODVkL01xTmRKSklxdE0vN04vT1pKWDZnWGlXTE1TS2tEUE1lCkNyZmtSZ1daOUtwZnc0WUxVdWkzd2FTZWEvbGpuVjBsTW9Eb2VsWUh6Vk1sRW5FTWZsZzRsYi9Nb0hXSEV6ZFlJeElWMXoKTTZFWjZQMzNacTJmY1ZHNUNBZ2dYTE0yQkRlYXdvODRrU1lmckJDeHBISFJMTmFybzJyUmNQUk8yWWhPaHp1WDY5WnJqbQprUFNtNXZzMVB6bjBMblpRUXhJNzFNWkhOSzlobVh1UUxaOHZaalBwOHBSR0U5Zi85b0NhYVBTdzZNdGp6VXZlN0RpNFVrCkR1bXhiT1BaMkF1aUtnaEE4TnVENTF5SzI5NzJTN0ttQ2hDbEY1UGdkakVVSk1SY3huNkp2ZS83WFZleVhCWnFGQmFkNFkKcXhHbS9mV2c3MHovRnExTlZsaUpnTThGQ1FvdGY3M0E4RDZneUYzK2xPTStvZFpQYTFhVFhwRWlEQnYwNEpXQld3ajhMYwpwUndXMzBwOVJ3ODBsQW1yMlZlMkFUaVk3VGlOenNQTEFlK2UxcFFFWGwxTlpRTnRQbFlnRXpyQTBmWTNCUG5NUHF3dUpVClhiS25WK3pUQTdHS1VUS3FtOExyYkJxcy9PVmhDQlF3VHhRM1lVcktTN1psaHBjcjh0STJjY2UzVEFoTFNJZUxYaVF2S3oKMUhPbm5XbDg4NUtpVXkyZmUvQkQ0LzMybEpBPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=",
					"host_share_key_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_share_key_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pass_phrase", "private_key"},
			},
		},
	})
}

func TestAccAlicloudBastionhostHostShareKey_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_share_key.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostShareKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostShareKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostsharekey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostShareKeyBasicDependence0)
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
					"instance_id":         "${data.alicloud_bastionhost_instances.default.instances.0.id}",
					"pass_phrase":         "NTIxeGlubXU=",
					"private_key":         "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJEZGNBVEl1cQpla0lyNXIzMUY4Z0NEc0FBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FDNzNHK3JTTWRsCjhmbTkycTlzTUxRWFZ6c2loRDlIblV1SnFib0JMbGhjY0F4Nk5SbnVoSTFGaHF4WXFnZ0RvVUFXdWltd2pQOHdvczZWNk4Kc3krc0pTbExFckZUdkZsWTBiL1JlN2tmOWxVcFkrL0o4Y2xNYUJ0MjRpNWdxelVuS1dFQ0diS2s4QlhuVTJUOTRNR291ZQphYmNsNWwwTGVvMzk1R1ZNa2VvZ1M2YzZkYmppS2hyYUhqcUxUcDg4TG5RNWZ3QW9jbnViZTIxN3pxRk1YaFZoaFl2QmlNCnRqS0dJYm5ZNzhoR0daSjd0VjdUUFlvL2EzQ3h5dXJ1ODhHeHRJd1h0cXJXWFFpZnp2Tzd6SUhHSDd1bmo3K212aURqRWwKa0VhbHpIa2VvY1AvMnFMQmxIcTlKaWRIa1FPVG5KWVE4M3dCRUpTM3R0OFBjd1VKc3E0aHhVbks5eHBqb2ZLcEtBaWRyNgpjbHZHLzltVTFSZmpYQlFFanBya1N3UW43d1ZCalRpZWZlM3dqMUdwREJtcTNXYTRZRjRXT29JSUNabmtIa0c0R1FhOG9HCjFES1Z3N29HUFRZenZ0ajJKYmY1dFNjSkZGWTVndTJvUjZKZk9HbzdXYlFUazZVYW9OaFdSNkZ3TUs3RlR1cjJDdFBhY1gKMWhQc2xjeXJ0aThyVUFBQVdRTXBaTW1XVVBrQUdmRGdxUGw0anNQdWV0cWFkdFpPMzN3Vm4zUXJsdFhjUUpEeUtBdVY3UQo5K0ZhekU2TjZ1bXBiT3AvUGM5bGFJN2paSW11V3R3RTBJMUFxZG9wT1dQWmVFQjFGN2dDbEw5UHJHZUtQN01PdFd1U2s4CkpaNFBNRXIwVG1IT3dzN3NvZDZpYjMxWWVlaXBsdUFOcFM5UFJPV1BiM1VwQ09TazdDalByQ3dwdTdHOGtIUUVCOGt4ZFIKQml5VDByclplV0J3TExHYkgxTjd0aW5rYVJoUkNtRDhadTBHV0k3SXFINHMvWGlBTmszR0dTTXdJVVBkMnJwdGVqUFdOdwpnR2pYRnUwbHB2TktRYU5GZmJQUFJEK1BRcUs5S3kraFdkN2o2MWVyR1ZUbmJZZ1lDaloxazRSUW16SDJDeGQwQkM1ajNRCktZa2hYWGg0TnQyRTVBa2RzbEh4TGNVdkx6N0tyT3V5RnZtSm83RHhUeUJyT0RZQlV6K3JHM0xpQXg0MmJydVFBMVFGUXMKc3EzeFNUYk4vYzVqajRvUXp3WHJvQ0p5bEJUY3Q5RVhzckdtcmNlQmI4OXNQWElrZjhCaS9kQlRGVDl1VUZqc1hvVW50SwpiSU9nWGtUUXJmZEViNWFmUitnY21XWlkxSmFjZEF0V2FTUGFCb2JGQ3FMamUwc2tqNHRnTkxnczVCQjZwYW1ic1FVRmpGClIvVzk0K2R6bitqM0VpWUIxMDFOYnJ4T1dCMmRDUE45clZxN0pCR2NGK3FyUjYxQzJvWFR0TWxNK1NFOGpCc1QwZEVYeUEKbjB6SVd2dEhidTdZekFidzRHSUszRGVmOU9vU3pJWHZFZmZrcE1FQjBlMHBnaEw0ZmlJYTBzZkdEbTA1Wk05MlNTNGpCcgp4KytBVHRqcXg1VnNJRkg5VzFoa252Szg0REZlSjZsdUxtVnAremdzdE04R0s5RGJqYmxtRGdQNDYwSHRoOU5kVFMrOTE3CkY4ajRYQ3FJTVJsbG1TVEd6VUFURDh3OEZkNExUdytSYlRBMjZJZnJHNFlZWFVNakJlNWxQdHFybGFLL0JEc2NJdWxRUWYKcmVtdmk3dXFrWUI4UHpsNFAreHRxRE9KUzljY2lxL0R3V3ppTGU3ekpnQ1d4ZFNvTFVLaVFkSjZDSUliMDVicWhXR1FDSApDZ2JKRW5hM0luY3p4cHl1ZmkrZUVDbklMUFk3c01oT1RNTmdNRUdkbGI1dnZyUEEvY2RmTDNNRDBSUERCQnh0allxRUlSCnJQZ1BiMGRjeUNmNmJZVTJhMUM2ZnNGeE0rbEdiY3doMWVHbWVCbFRsUy95YmFuWk92NXdtckQyODVlbXpWeGJ2RU1PTlYKb0tsYnYzK0JlNHRJV0w3VHRsR0VLNGs2L0hIdUI0SkdMV0dxNmZ6TUp4V0lreENLUTFXOXNhcXVQckFVZzVVdGFQcG9KSApqb2t4ellRWVZpK1Bzd3psOExhMTlHNVRialJ5dmd4aUhUYlB1K2dFaGdDT3pmZnA4TS85UDhYanNGUi9sU1RsbnQzUGlwCjloNVdDUllEeStTUWxiZ0E2OXRMQ1pBdUl0aHJ5bk5QbmZRZHpMdytHVWhhTDQ2ai9yeGdJV2gwWmc3ZWJMZHJvVExad2YKVjdiWVVFeXdZbkMwQjlpODNKNGZxSkwxOHdxeWkyL2xQOXIxR0ZRKytvSTJaMjk1THN3N0NaWTZKaDhJLytQSVVoeUNlZgpoNHhNd2FEUGFKeUx4MEdVOW9XdE5oSXRuZXUyNUtOd1Q4ODVkL01xTmRKSklxdE0vN04vT1pKWDZnWGlXTE1TS2tEUE1lCkNyZmtSZ1daOUtwZnc0WUxVdWkzd2FTZWEvbGpuVjBsTW9Eb2VsWUh6Vk1sRW5FTWZsZzRsYi9Nb0hXSEV6ZFlJeElWMXoKTTZFWjZQMzNacTJmY1ZHNUNBZ2dYTE0yQkRlYXdvODRrU1lmckJDeHBISFJMTmFybzJyUmNQUk8yWWhPaHp1WDY5WnJqbQprUFNtNXZzMVB6bjBMblpRUXhJNzFNWkhOSzlobVh1UUxaOHZaalBwOHBSR0U5Zi85b0NhYVBTdzZNdGp6VXZlN0RpNFVrCkR1bXhiT1BaMkF1aUtnaEE4TnVENTF5SzI5NzJTN0ttQ2hDbEY1UGdkakVVSk1SY3huNkp2ZS83WFZleVhCWnFGQmFkNFkKcXhHbS9mV2c3MHovRnExTlZsaUpnTThGQ1FvdGY3M0E4RDZneUYzK2xPTStvZFpQYTFhVFhwRWlEQnYwNEpXQld3ajhMYwpwUndXMzBwOVJ3ODBsQW1yMlZlMkFUaVk3VGlOenNQTEFlK2UxcFFFWGwxTlpRTnRQbFlnRXpyQTBmWTNCUG5NUHF3dUpVClhiS25WK3pUQTdHS1VUS3FtOExyYkJxcy9PVmhDQlF3VHhRM1lVcktTN1psaHBjcjh0STJjY2UzVEFoTFNJZUxYaVF2S3oKMUhPbm5XbDg4NUtpVXkyZmUvQkQ0LzMybEpBPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0=",
					"host_share_key_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_share_key_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pass_phrase", "private_key"},
			},
		},
	})
}

var AlicloudBastionhostHostShareKeyMap0 = map[string]string{
	"instance_id":              CHECKSET,
	"host_share_key_id":        CHECKSET,
	"private_key_finger_print": CHECKSET,
}

func AlicloudBastionhostHostShareKeyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_bastionhost_instances" "default" {}

`, name)
}

func TestUnitAlicloudBastionhostHostShareKey(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id":         "CreateHostShareKeyValue",
		"pass_phrase":         "CreateHostShareKeyValue",
		"host_share_key_name": "CreateHostShareKeyValue",
		"private_key":         "CreateHostShareKeyValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetHostShareKey
		"HostShareKey": map[string]interface{}{
			"InstanceId":            "CreateHostShareKeyValue",
			"HostShareKeyId":        "CreateHostShareKeyValue",
			"HostShareKeyName":      "CreateHostShareKeyValue",
			"PrivateKeyFingerPrint": "CreateHostShareKeyValue",
		},
		"HostShareKeyId": "CreateHostShareKeyValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateHostShareKey
		"HostShareKeyId": "CreateHostShareKeyValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_bastionhost_host_share_key", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBastionhostHostShareKeyCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetHostShareKey Response
		"HostShareKeyId": "CreateHostShareKeyValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHostShareKey" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBastionhostHostShareKeyCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBastionhostHostShareKeyUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyHostShareKey
	attributesDiff := map[string]interface{}{
		"host_share_key_name": "ModifyHostShareKeyValue",
		"private_key":         "ModifyHostShareKeyValue",
		"pass_phrase":         "ModifyHostShareKeyValue",
	}
	diff, err := newInstanceDiff("alicloud_bastionhost_host_share_key", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetHostShareKey Response
		"HostShareKey": map[string]interface{}{
			"HostShareKeyName": "ModifyHostShareKeyValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyHostShareKey" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBastionhostHostShareKeyUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_bastionhost_host_share_key", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetHostShareKey" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBastionhostHostShareKeyRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBastionhostHostShareKeyDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_bastionhost_host_share_key", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_host_share_key"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteHostShareKey" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBastionhostHostShareKeyDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
