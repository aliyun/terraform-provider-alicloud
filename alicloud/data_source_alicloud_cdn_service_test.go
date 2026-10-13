package alicloud

import (
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudCDNServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cdn_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_cdn_service.current", "status", "Opened"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "internet_charge_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "opening_time"),
					//resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_charge_type"),
					//resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_affect_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCdnServiceDataSource = `
data "alicloud_cdn_service" "current" {
	enable = "On"
	internet_charge_type = "PayByTraffic"
}
`

func testUnitCdnServiceMockMeta(t *testing.T) *connectivity.AliyunClient {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: new(sync.Map), SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	meta, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return meta
}

func TestUnitAliCloudCdnServiceRead(t *testing.T) {
	openedResponse := func() map[string]interface{} {
		return map[string]interface{}{
			"OpeningTime":        "2024-01-01T16:00:00+08:00",
			"InternetChargeType": "PayByBandwidth",
			"ChangingChargeType": "",
			"ChangingAffectTime": "",
		}
	}
	notOpenedResponse := func() map[string]interface{} {
		return map[string]interface{}{
			"OpeningTime":        "",
			"InternetChargeType": "",
			"ChangingChargeType": "",
			"ChangingAffectTime": "",
		}
	}

	t.Run("OpenedChargeTypeDiffersDoesNotModify", func(t *testing.T) {
		var actions []string
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			actions = append(actions, *action)
			return openedResponse(), nil
		})
		defer patches.Reset()

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudCdnService().Schema, map[string]interface{}{
			"enable":               "On",
			"internet_charge_type": "PayByTraffic",
		})
		if err := dataSourceAliCloudCdnServiceRead(d, testUnitCdnServiceMockMeta(t)); err != nil {
			t.Fatal(err)
		}
		assert.NotContains(t, actions, "ModifyCdnService")
		assert.Equal(t, []string{"DescribeCdnService", "DescribeCdnService"}, actions)
		assert.Equal(t, "CdnServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
		assert.Equal(t, "PayByBandwidth", d.Get("internet_charge_type"))
	})

	t.Run("NotOpenedEnableOnOpensService", func(t *testing.T) {
		var actions []string
		var openBody map[string]interface{}
		describeCount := 0
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			actions = append(actions, *action)
			switch *action {
			case "DescribeCdnService":
				describeCount++
				if describeCount == 1 {
					return notOpenedResponse(), nil
				}
				return map[string]interface{}{
					"OpeningTime":        "2024-01-01T16:00:00+08:00",
					"InternetChargeType": "PayByTraffic",
					"ChangingChargeType": "",
					"ChangingAffectTime": "",
				}, nil
			case "OpenCdnService":
				openBody = body
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{}, nil
		})
		defer patches.Reset()

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudCdnService().Schema, map[string]interface{}{
			"enable":               "On",
			"internet_charge_type": "PayByTraffic",
		})
		if err := dataSourceAliCloudCdnServiceRead(d, testUnitCdnServiceMockMeta(t)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"DescribeCdnService", "OpenCdnService", "DescribeCdnService"}, actions)
		assert.Equal(t, "PayByTraffic", openBody["InternetChargeType"])
		assert.Equal(t, "CdnServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
	})

	t.Run("NotOpenedEnableOnHasOpenedTolerated", func(t *testing.T) {
		var actions []string
		describeCount := 0
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			actions = append(actions, *action)
			switch *action {
			case "DescribeCdnService":
				describeCount++
				if describeCount == 1 {
					return notOpenedResponse(), nil
				}
				return openedResponse(), nil
			case "OpenCdnService":
				return nil, &tea.SDKError{
					Code:       String("CdnService.HasOpened"),
					Data:       String("CdnService.HasOpened"),
					Message:    String("The CDN service has been opened."),
					StatusCode: tea.Int(400),
				}
			}
			return map[string]interface{}{}, nil
		})
		defer patches.Reset()

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudCdnService().Schema, map[string]interface{}{
			"enable":               "On",
			"internet_charge_type": "PayByTraffic",
		})
		if err := dataSourceAliCloudCdnServiceRead(d, testUnitCdnServiceMockMeta(t)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"DescribeCdnService", "OpenCdnService", "DescribeCdnService"}, actions)
		assert.Equal(t, "CdnServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
	})

	t.Run("EnableOffReadOnly", func(t *testing.T) {
		var actions []string
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			actions = append(actions, *action)
			return openedResponse(), nil
		})
		defer patches.Reset()

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudCdnService().Schema, map[string]interface{}{
			"enable": "Off",
		})
		if err := dataSourceAliCloudCdnServiceRead(d, testUnitCdnServiceMockMeta(t)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"DescribeCdnService"}, actions)
		assert.Equal(t, "CdnServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
	})
}
