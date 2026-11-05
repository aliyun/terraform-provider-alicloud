package ecd

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNetworkPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkPackageCreate,
		Read:   resourceNetworkPackageRead,
		Update: resourceNetworkPackageUpdate,
		Delete: resourceNetworkPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNetworkPackage"
	request := make(map[string]interface{})
	ecdService := NewService(client)
	var err error

	request["Bandwidth"] = d.Get("bandwidth")

	request["OfficeSiteId"] = d.Get("office_site_id")

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_network_package", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["NetworkPackageId"]))

	stateConf := helper.BuildStateConf([]string{"Creating"}, []string{"InUse"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceNetworkPackageRead(d, meta)
}
func resourceNetworkPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdNetworkPackage(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_network_package ecdService.DescribeEcdNetworkPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := object["Bandwidth"]; ok && fmt.Sprint(v) != "0" {
		bandwidthVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("bandwidth", bandwidthVal)
	}
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("status", object["NetworkPackageStatus"])
	return nil
}
func resourceNetworkPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"NetworkPackageId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		request["Bandwidth"] = d.Get("bandwidth")
	}
	action := "ModifyNetworkPackage"
	ecdService := NewService(client)
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectNetworkPackageStatus.ModificationNotSupport"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"Creating", "Updating"}, []string{"InUse"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceNetworkPackageRead(d, meta)
}
func resourceNetworkPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNetworkPackages"
	var response map[string]interface{}
	ecdService := NewService(client)
	var err error
	request := map[string]interface{}{
		"NetworkPackageId": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectNetworkPackageStatus.DeletionNotSupport"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"Releasing"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
