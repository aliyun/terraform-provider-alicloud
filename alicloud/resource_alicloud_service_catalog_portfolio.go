package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudServiceCatalogPortfolio() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudServiceCatalogPortfolioCreate,
		Read:   resourceAlicloudServiceCatalogPortfolioRead,
		Update: resourceAlicloudServiceCatalogPortfolioUpdate,
		Delete: resourceAlicloudServiceCatalogPortfolioDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"portfolio_arn": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"portfolio_name": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"provider_name": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
		},
	}
}

func resourceAlicloudServiceCatalogPortfolioCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewSrvcatalogClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("portfolio_name"); ok {
		request["PortfolioName"] = v
	}
	if v, ok := d.GetOk("provider_name"); ok {
		request["ProviderName"] = v
	}

	var response map[string]interface{}
	action := "CreatePortfolio"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_catalog_portfolio", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PortfolioId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_service_catalog_portfolio")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudServiceCatalogPortfolioRead(d, meta)
}

func resourceAlicloudServiceCatalogPortfolioRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceCatalogService := ServicecatalogService{client}

	object, err := serviceCatalogService.DescribeServiceCatalogPortfolio(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_catalog_portfolio servicecatalogService.DescribeServiceCatalogPortfolio Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("description", object["Description"])
	d.Set("portfolio_arn", object["PortfolioArn"])
	d.Set("portfolio_name", object["PortfolioName"])
	d.Set("provider_name", object["ProviderName"])

	return nil
}

func resourceAlicloudServiceCatalogPortfolioUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewSrvcatalogClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"PortfolioId": d.Id(),
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("portfolio_name") {
		update = true
	}
	request["PortfolioName"] = d.Get("portfolio_name")
	if d.HasChange("provider_name") {
		update = true
	}
	request["ProviderName"] = d.Get("provider_name")

	if update {
		action := "UpdatePortfolio"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudServiceCatalogPortfolioRead(d, meta)
}

func resourceAlicloudServiceCatalogPortfolioDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewSrvcatalogClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"PortfolioId": d.Id(),
	}

	action := "DeletePortfolio"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidPortfolio.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
