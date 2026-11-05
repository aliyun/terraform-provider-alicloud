package market

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrderCreate,
		Read:   resourceOrderRead,
		Delete: resourceOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(8, 12),
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PostPaid",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 256),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Day", "Month", "Year"}, false),
			},
			"package_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"coupon_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// lintignore: S006
			"components": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
				ForceNew: true,
			},
		},
	}
}

func resourceOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildOrderArgs(d, meta)
	if err != nil {
		return errs.WrapError(err)
	}
	request.RegionId = client.RegionId
	raw, err := client.WithMarketClient(func(marketClient *market.Client) (interface{}, error) {
		return marketClient.CreateOrder(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_market_order", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*market.CreateOrderResponse)

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(response.OrderId)
	return resourceOrderRead(d, meta)
}

func resourceOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	marketService := NewMarketService(client)

	object, err := marketService.DescribeMarketOrder(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	err = d.Set("product_code", object.ProductCode)
	if err != nil {
		return errs.WrapError(err)
	}
	orderPayType := "PostPaid"
	if strings.Split(object.ProductSkuCode, "-")[1] == "prepay" {
		orderPayType = "PrePaid"
	}
	err = d.Set("pay_type", orderPayType)
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("quantity", object.Quantity)
	if err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func resourceOrderDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildOrderArgs(d *schema.ResourceData, meta interface{}) (*market.CreateOrderRequest, error) {
	packageVersion := d.Get("package_version").(string)
	couponId := ""
	if v, ok := d.GetOk("coupon_id"); ok {
		couponId = v.(string)
	}
	duration := d.Get("duration").(int)
	pricingCycle := d.Get("pricing_cycle").(string)
	productCode := d.Get("product_code").(string)
	quantity := d.Get("quantity").(int)
	orderPayType := "postpay"
	if d.Get("pay_type").(string) == "PrePaid" {
		orderPayType = "prepay"
	}
	skuCode := fmt.Sprintf("%s-%s", productCode, orderPayType)
	components := d.Get("components").(map[string]interface{})
	components["package_version"] = packageVersion
	componentsJson, _ := json.Marshal(components)
	commodity := fmt.Sprintf(`{
    "components": %s,
    "couponId": "%s", 
    "duration": %d, 
    "pricingCycle": "%s",  
    "productCode": "%s",
    "quantity": %d,	
    "skuCode": "%s"  
	}`, componentsJson, couponId, duration, pricingCycle, productCode, quantity, skuCode)

	request := market.CreateCreateOrderRequest()
	request.ClientToken = helper.BuildClientToken(request.GetActionName())
	request.OrderType = "INSTANCE_BUY"
	request.PaymentType = "AUTO"
	request.Commodity = commodity
	return request, nil
}
