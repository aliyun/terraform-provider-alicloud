package alicloud

import (
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudApiGatewayService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudApigatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("ApiGatewayServicHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)

	response, err := client.NewTeaCommonClientWithEndpoint(connectivity.OpenApiGatewayService, func(teaClient *rpc.Client) (map[string]interface{}, error) {
		return teaClient.DoRequest(StringPointer("OpenApiGatewayService"), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	})
	addDebug("OpenApiGatewayService", response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("ApiGatewayServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_service", "OpenApiGatewayService", AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}
