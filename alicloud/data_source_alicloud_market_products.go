package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudProductsRead,

		Schema: map[string]*schema.Schema{
			"sort": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"user_count-desc", "created_on-desc", "price-desc", "score-desc"}, false),
			},
			"category_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"APP", "SERVICE", "MIRROR", "DOWNLOAD", "API_SERVICE"}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"supplier_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"supplier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"suggested_price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"score": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warranty_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_way": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := market.CreateDescribeProductsRequest()
	request.RegionId = client.RegionId
	var productsFilter []market.DescribeProductsFilter
	var product market.DescribeProductsFilter
	if v, ok := d.GetOk("sort"); ok && v.(string) != "" {
		product.Key = "sort"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	if v, ok := d.GetOk("category_id"); ok && v.(string) != "" {
		product.Key = "categoryId"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	if v, ok := d.GetOk("product_type"); ok && v.(string) != "" {
		product.Key = "productType"
		product.Value = v.(string)
		productsFilter = append(productsFilter, product)
	}
	request.Filter = &productsFilter
	var allProduct []market.ProductItem
	for {
		raw, err := client.WithMarketClient(func(marketClient *market.Client) (interface{}, error) {
			return marketClient.DescribeProducts(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_market_products", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*market.DescribeProductsResponse)

		allProduct = append(allProduct, response.ProductItems.ProductItem...)

		if len(allProduct) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	return productsDescriptionAttributes(d, allProduct)
}

func productsDescriptionAttributes(d *schema.ResourceData, allProduct []market.ProductItem) error {
	var ids []string
	var s []map[string]interface{}
	for _, p := range allProduct {
		mapping := map[string]interface{}{
			"code":              p.Code,
			"name":              p.Name,
			"category_id":       p.CategoryId,
			"supplier_id":       p.SupplierId,
			"supplier_name":     p.SupplierName,
			"short_description": p.ShortDescription,
			"tags":              p.Tags,
			"suggested_price":   p.SuggestedPrice,
			"target_url":        p.TargetUrl,
			"image_url":         p.ImageUrl,
			"score":             p.Score,
			"operation_system":  p.OperationSystem,
			"warranty_date":     p.WarrantyDate,
			"delivery_date":     p.DeliveryDate,
			"delivery_way":      p.DeliveryWay,
		}

		ids = append(ids, p.Code)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("products", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
