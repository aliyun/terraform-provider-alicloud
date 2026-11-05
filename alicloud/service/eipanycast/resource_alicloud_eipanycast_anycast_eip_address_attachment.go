// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package eipanycast

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAnycastEipAddressAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAnycastEipAddressAttachmentCreate,
		Read:   resourceAnycastEipAddressAttachmentRead,
		Update: resourceAnycastEipAddressAttachmentUpdate,
		Delete: resourceAnycastEipAddressAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"anycast_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"association_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Default", "Normal"}, false),
			},
			"bind_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"SlbInstance", "NetworkInterface"}, false),
			},
			"bind_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pop_locations": {
				Type:     schema.TypeSet,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("association_mode"); ok && v.(string) == "Default" {
						return true
					}
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pop_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAnycastEipAddressAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "AssociateAnycastEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["BindInstanceType"] = d.Get("bind_instance_type")
	request["AnycastId"] = d.Get("anycast_id")
	request["BindInstanceId"] = d.Get("bind_instance_id")
	request["BindInstanceRegionId"] = d.Get("bind_instance_region_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	if v, ok := d.GetOk("association_mode"); ok {
		request["AssociationMode"] = v
	}
	if v, ok := d.GetOk("pop_locations"); ok {
		popLocationsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["PopLocation"] = dataLoopTmp["pop_location"]
			popLocationsMaps = append(popLocationsMaps, dataLoopMap)
		}
		request["PopLocations"] = popLocationsMaps
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Eipanycast", "2020-03-09", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Anycast"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_eipanycast_anycast_eip_address_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["AnycastId"], request["BindInstanceId"], request["BindInstanceRegionId"], request["BindInstanceType"]))

	eipanycastService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"BINDED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAnycastEipAddressAttachmentUpdate(d, meta)
}

func resourceAnycastEipAddressAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := NewService(client)

	objectRaw, err := eipanycastService.DescribeEipanycastAnycastEipAddressAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eipanycast_anycast_eip_address_attachment DescribeEipanycastAnycastEipAddressAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("association_mode", objectRaw["AssociationMode"])
	d.Set("bind_time", objectRaw["BindTime"])
	d.Set("private_ip_address", objectRaw["PrivateIpAddress"])
	d.Set("status", objectRaw["Status"])
	d.Set("bind_instance_id", objectRaw["BindInstanceId"])
	d.Set("bind_instance_region_id", objectRaw["BindInstanceRegionId"])
	d.Set("bind_instance_type", objectRaw["BindInstanceType"])
	popLocations1Raw := objectRaw["PopLocations"]
	popLocationsMaps := make([]map[string]interface{}, 0)
	if popLocations1Raw != nil {
		for _, popLocationsChild1Raw := range popLocations1Raw.([]interface{}) {
			popLocationsMap := make(map[string]interface{})
			popLocationsChild1Raw := popLocationsChild1Raw.(map[string]interface{})
			popLocationsMap["pop_location"] = popLocationsChild1Raw["PopLocation"]
			popLocationsMaps = append(popLocationsMaps, popLocationsMap)
		}
	}
	d.Set("pop_locations", popLocationsMaps)
	parts, _ := helper.ParseResourceId(d.Id(), 4)
	d.Set("anycast_id", parts[0])

	return nil
}

func resourceAnycastEipAddressAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateAnycastEipAddressAssociations"
	var err error
	request = make(map[string]interface{})
	request["BindInstanceId"] = parts[1]
	request["AnycastId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("association_mode") {
		update = true
		request["AssociationMode"] = d.Get("association_mode")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Eipanycast", "2020-03-09", action, nil, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Anycast"}) || helper.NeedRetry(err) {
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
		eipanycastService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"BINDED"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("pop_locations") {
		oldEntry, newEntry := d.GetChange("pop_locations")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UpdateAnycastEipAddressAssociations"
			request = make(map[string]interface{})
			request["BindInstanceId"] = parts[1]
			request["AnycastId"] = parts[0]
			request["RegionId"] = client.RegionId
			request["ClientToken"] = helper.BuildClientToken(action)
			localData := removed.List()
			popLocationDeleteListMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["PopLocation"] = dataLoopTmp["pop_location"]
				popLocationDeleteListMaps = append(popLocationDeleteListMaps, dataLoopMap)
			}
			request["PopLocationDeleteList"] = popLocationDeleteListMaps

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Eipanycast", "2020-03-09", action, nil, request, true)
				request["ClientToken"] = helper.BuildClientToken(action)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Anycast"}) || helper.NeedRetry(err) {
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
			eipanycastService := NewService(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"BINDED"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "UpdateAnycastEipAddressAssociations"
			request = make(map[string]interface{})
			request["BindInstanceId"] = parts[1]
			request["AnycastId"] = parts[0]
			request["RegionId"] = client.RegionId
			request["ClientToken"] = helper.BuildClientToken(action)
			localData := added.List()
			popLocationAddListMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["PopLocation"] = dataLoopTmp["pop_location"]
				popLocationAddListMaps = append(popLocationAddListMaps, dataLoopMap)
			}
			request["PopLocationAddList"] = popLocationAddListMaps

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Eipanycast", "2020-03-09", action, nil, request, true)
				request["ClientToken"] = helper.BuildClientToken(action)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Anycast"}) || helper.NeedRetry(err) {
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
			eipanycastService := NewService(client)
			stateConf := helper.BuildStateConf([]string{}, []string{"BINDED"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}

		}

	}
	return resourceAnycastEipAddressAttachmentRead(d, meta)
}

func resourceAnycastEipAddressAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnassociateAnycastEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["BindInstanceType"] = parts[3]
	request["AnycastId"] = parts[0]
	request["BindInstanceId"] = parts[1]
	request["BindInstanceRegionId"] = parts[2]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Eipanycast", "2020-03-09", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Anycast"}) || helper.NeedRetry(err) {
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

	eipanycastService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
