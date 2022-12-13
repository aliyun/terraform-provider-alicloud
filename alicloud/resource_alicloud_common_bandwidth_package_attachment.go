package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAliyunCommonBandwidthPackageAttachmentRead,
		Update: resourceAliyunCommonBandwidthPackageAttachmentUpdate,
		Delete: resourceAliyunCommonBandwidthPackageAttachmentDelete,
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
			},
		},
	}
}

func resourceAliyunCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateAddCommonBandwidthPackageIpRequest()
	request.RegionId = client.RegionId
	request.BandwidthPackageId = Trim(d.Get("bandwidth_package_id").(string))
	request.IpInstanceId = Trim(d.Get("instance_id").(string))
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AddCommonBandwidthPackageIp(request)
	})
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AddCommonBandwidthPackageIp(request)
		})
		//Waiting for unassociate the common bandwidth package
		if err != nil && !IsExpectedErrors(err, []string{"IpInstanceId.AlreadyInBandwidthPackage"}) {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Eip", "EipOperation.TooFrequently"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_common_bandwidth_package_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//check the common bandwidth package attachment
	d.SetId(request.BandwidthPackageId + COLON_SEPARATED + request.IpInstanceId)
	if err := vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Available, 5*DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunCommonBandwidthPackageAttachmentUpdate(d, meta)
}

func resourceAliyunCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]
	_, err = vpcService.DescribeCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)

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

func resourceAliyunCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
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
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

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

		if err = vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Available, 5*DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliyunCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.RemoveCommonBandwidthPackageIp(request)
		})
		//Waiting for unassociate the common bandwidth package
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Eip"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"OperationUnsupported.IpNotInCbwp"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
