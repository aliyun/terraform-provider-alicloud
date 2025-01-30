// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BandwidthPackageId"] = d.Get("bandwidth_package_id")
	query["IpInstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ip_type"); ok {
		request["IpType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "TaskConflict", "EipOperation.TooFrequently", "IncorrectStatus.Eip"}) || NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v:%v", query["BandwidthPackageId"], query["IpInstanceId"]))

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

	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["BandwidthPackageId"] != nil {
		d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	}

	publicIpAddresse1RawObj, _ := jsonpath.Get("$.PublicIpAddresses.PublicIpAddresse[*]", objectRaw)
	publicIpAddresse1Raw := make([]interface{}, 0)
	if publicIpAddresse1RawObj != nil {
		publicIpAddresse1Raw = publicIpAddresse1RawObj.([]interface{})
	}

	publicIpAddresseChild1Raw := publicIpAddresse1Raw[0].(map[string]interface{})
	if publicIpAddresseChild1Raw["AllocationId"] != nil {
		d.Set("instance_id", publicIpAddresseChild1Raw["AllocationId"])
	}

	objectRaw, err = cbwpServiceV2.DescribeDescribeEipAddresses(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["Bandwidth"] != nil {
		d.Set("bandwidth_package_bandwidth", objectRaw["Bandwidth"])
	}
	if objectRaw["BandwidthPackageId"] != nil {
		d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("bandwidth_package_id", parts[0])
	d.Set("instance_id", parts[1])

	return nil
}

func resourceAliCloudCbwpCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyCommonBandwidthPackageIpBandwidth"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BandwidthPackageId"] = parts[0]
	query["EipId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("bandwidth_package_bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth_package_bandwidth")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
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
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["BandwidthPackageId"] = parts[0]
			query["EipId"] = parts[1]
			request["RegionId"] = client.RegionId
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
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
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BandwidthPackageId"] = parts[0]
	query["IpInstanceId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
