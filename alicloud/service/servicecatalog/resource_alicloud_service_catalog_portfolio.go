// Package servicecatalog. This file is generated automatically. Please do not modify it manually, thank you!
package servicecatalog

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudServiceCatalogPortfolio() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudServiceCatalogPortfolioCreate,
		Read:   resourceAliCloudServiceCatalogPortfolioRead,
		Update: resourceAliCloudServiceCatalogPortfolioUpdate,
		Delete: resourceAliCloudServiceCatalogPortfolioDelete,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"portfolio_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"portfolio_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudServiceCatalogPortfolioCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePortfolio"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["PortfolioName"] = d.Get("portfolio_name")
	request["ProviderName"] = d.Get("provider_name")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_service_catalog_portfolio", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PortfolioId"]))

	return resourceAliCloudServiceCatalogPortfolioRead(d, meta)
}

func resourceAliCloudServiceCatalogPortfolioRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceCatalogServiceV2 := NewService(client)

	objectRaw, err := serviceCatalogServiceV2.DescribeServiceCatalogPortfolio(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_catalog_portfolio DescribeServiceCatalogPortfolio Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["PortfolioArn"] != nil {
		d.Set("portfolio_arn", objectRaw["PortfolioArn"])
	}
	if objectRaw["PortfolioName"] != nil {
		d.Set("portfolio_name", objectRaw["PortfolioName"])
	}
	if objectRaw["ProviderName"] != nil {
		d.Set("provider_name", objectRaw["ProviderName"])
	}

	return nil
}

func resourceAliCloudServiceCatalogPortfolioUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdatePortfolio"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PortfolioId"] = d.Id()
	query["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudServiceCatalogPortfolioRead(d, meta)
}

func resourceAliCloudServiceCatalogPortfolioDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePortfolio"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PortfolioId"] = d.Id()
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidPortfolio.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
