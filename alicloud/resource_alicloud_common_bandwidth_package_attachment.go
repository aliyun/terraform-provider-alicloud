// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCbwpCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCbwpCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAliCloudCbwpCommonBandwidthPackageAttachmentRead,
		Update: resourceAliCloudCbwpCommonBandwidthPackageAttachmentUpdate,
		Delete: resourceAliCloudCbwpCommonBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("cancel_common_bandwidth_package_ip_bandwidth").(bool)
				},
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cancel_common_bandwidth_package_ip_bandwidth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"EIP"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCbwpCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddCommonBandwidthPackageIp"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["BandwidthPackageId"] = d.Get("bandwidth_package_id")
	request["IpInstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ip_type"); ok {
		request["IpType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "TaskConflict", "EipOperation.TooFrequently"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_common_bandwidth_package_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["BandwidthPackageId"], request["IpInstanceId"]))

	cbwpServiceV2 := CbwpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCbwpCommonBandwidthPackageAttachmentUpdate(d, meta)
}

func resourceAliCloudCbwpCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbwpServiceV2 := CbwpServiceV2{client}

	objectRaw, err := cbwpServiceV2.DescribeCbwpCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_common_bandwidth_package_attachment DescribeCbwpCommonBandwidthPackageAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]
	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)

	objectRaw, err = cbwpServiceV2.DescribeDescribeEipAddresses(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("bandwidth_package_bandwidth", objectRaw["Bandwidth"])

	return nil
}

func resourceAliCloudCbwpCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyCommonBandwidthPackageIpBandwidth"
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["BandwidthPackageId"] = parts[0]
	request["EipId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("bandwidth_package_bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth_package_bandwidth")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cbwpServiceV2 := CbwpServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("cancel_common_bandwidth_package_ip_bandwidth") {

		target := d.Get("cancel_common_bandwidth_package_ip_bandwidth").(bool)
		if target == true {
			parts = strings.Split(d.Id(), ":")
			action = "CancelCommonBandwidthPackageIpBandwidth"
			conn, err = client.NewCbwpClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["BandwidthPackageId"] = parts[0]
			request["EipId"] = parts[1]
			request["RegionId"] = client.RegionId
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

				if err != nil {
					if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			cbwpServiceV2 := CbwpServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}

	return resourceAliCloudCbwpCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAliCloudCbwpCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RemoveCommonBandwidthPackageIp"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["BandwidthPackageId"] = parts[0]
	request["IpInstanceId"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "TaskConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationUnsupported.IpNotInCbwp"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cbwpServiceV2 := CbwpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
