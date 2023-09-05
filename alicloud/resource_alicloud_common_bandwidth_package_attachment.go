package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAliCloudCommonBandwidthPackageAttachmentRead,
		Update: resourceAliCloudCommonBandwidthPackageAttachmentUpdate,
		Delete: resourceAliCloudCommonBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_package_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("cancel_common_bandwidth_package_ip_bandwidth").(bool)
				},
			},
			"cancel_common_bandwidth_package_ip_bandwidth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "AddCommonBandwidthPackageIp"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["BandwidthPackageId"] = d.Get("bandwidth_package_id")
	request["IpInstanceId"] = d.Get("instance_id")
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Eip", "EipOperation.TooFrequently"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	//check the common bandwidth package attachment
	d.SetId(fmt.Sprint(request["BandwidthPackageId"], COLON_SEPARATED, request["IpInstanceId"]))
	cbwpServiceV2 := CbwpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 0*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudCommonBandwidthPackageAttachmentUpdate(d, meta)
}

func resourceAliCloudCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbwpServiceV2 := CbwpServiceV2{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]
	_, err = cbwpServiceV2.DescribeCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)

	vpcService := VpcService{client}
	object, err := vpcService.DescribeEipAddress(ipInstanceId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_package_bandwidth", object["Bandwidth"])

	return nil
}

func resourceAliCloudCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbwpServiceV2 := CbwpServiceV2{client}
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"BandwidthPackageId": parts[0],
		"EipId":              parts[1],
	}

	if d.HasChange("bandwidth_package_bandwidth") {
		update = true
		if v, ok := d.GetOk("bandwidth_package_bandwidth"); ok {
			request["Bandwidth"] = v
		}
	}

	if update {
		action := "ModifyCommonBandwidthPackageIpBandwidth"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false

	if d.HasChange("cancel_common_bandwidth_package_ip_bandwidth") {
		if v, ok := d.GetOkExists("cancel_common_bandwidth_package_ip_bandwidth"); ok && v.(bool) == true {
			update = true
		}
	}

	if update {
		action := "CancelCommonBandwidthPackageIpBandwidth"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAliCloudCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	action := "RemoveCommonBandwidthPackageIp"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCbwpClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["BandwidthPackageId"] = bandwidthPackageId
	request["IpInstanceId"] = ipInstanceId
	request["ClientToken"] = buildClientToken(action)

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)
		//Waiting for unassociate the common bandwidth package
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Eip"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 0*time.Second, cbwpServiceV2.CbwpCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
