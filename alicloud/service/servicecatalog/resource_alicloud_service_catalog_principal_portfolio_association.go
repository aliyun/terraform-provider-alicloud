// Package servicecatalog. This file is generated automatically. Please do not modify it manually, thank you!
package servicecatalog

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudServiceCatalogPrincipalPortfolioAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudServiceCatalogPrincipalPortfolioAssociationCreate,
		Read:   resourceAliCloudServiceCatalogPrincipalPortfolioAssociationRead,
		Delete: resourceAliCloudServiceCatalogPrincipalPortfolioAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"portfolio_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudServiceCatalogPrincipalPortfolioAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociatePrincipalWithPortfolio"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PortfolioId"] = d.Get("portfolio_id")
	request["PrincipalId"] = d.Get("principal_id")
	request["PrincipalType"] = d.Get("principal_type")

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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_service_catalog_principal_portfolio_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["PrincipalId"], request["PrincipalType"], request["PortfolioId"]))

	return resourceAliCloudServiceCatalogPrincipalPortfolioAssociationRead(d, meta)
}

func resourceAliCloudServiceCatalogPrincipalPortfolioAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceCatalogServiceV2 := NewService(client)

	objectRaw, err := serviceCatalogServiceV2.DescribeServiceCatalogPrincipalPortfolioAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_catalog_principal_portfolio_association DescribeServiceCatalogPrincipalPortfolioAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["PrincipalId"] != nil {
		d.Set("principal_id", objectRaw["PrincipalId"])
	}
	if objectRaw["PrincipalType"] != nil {
		d.Set("principal_type", objectRaw["PrincipalType"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("principal_id", parts[0])
	d.Set("principal_type", parts[1])
	d.Set("portfolio_id", parts[2])

	return nil
}

func resourceAliCloudServiceCatalogPrincipalPortfolioAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DisassociatePrincipalFromPortfolio"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PortfolioId"] = parts[2]
	request["PrincipalId"] = parts[0]
	request["PrincipalType"] = parts[1]

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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
