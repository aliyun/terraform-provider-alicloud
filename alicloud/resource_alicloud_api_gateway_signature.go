package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunApigatewaySignature() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewaySignatureCreate,
		Read:   resourceAliyunApigatewaySignatureRead,
		Update: resourceAliyunApigatewaySignatureUpdate,
		Delete: resourceAliyunApigatewaySignatureDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"signature_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"signature_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"signature_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewaySignatureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateCreateSignatureRequest()
	request.RegionId = client.RegionId
	request.SignatureName = d.Get("signature_name").(string)
	request.SignatureKey = d.Get("signature_key").(string)
	request.SignatureSecret = d.Get("signature_secret").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateSignature(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateSignatureResponse)
		d.SetId(response.SignatureId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_signature", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resourceAliyunApigatewaySignatureRead(d, meta)
}

func resourceAliyunApigatewaySignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := cloudApiService.DescribeApiGatewaySignature(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("signature_name", object.SignatureName)
		d.Set("signature_key", object.SignatureKey)
		d.Set("signature_secret", object.SignatureSecret)
		d.Set("created_time", object.CreatedTime)
		d.Set("modified_time", object.ModifiedTime)
		return nil
	}); err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	return nil
}

func resourceAliyunApigatewaySignatureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("signature_name") || d.HasChange("signature_key") || d.HasChange("signature_secret") {
		request := cloudapi.CreateModifySignatureRequest()
		request.RegionId = client.RegionId
		request.SignatureId = d.Id()
		request.SignatureName = d.Get("signature_name").(string)
		request.SignatureKey = d.Get("signature_key").(string)
		request.SignatureSecret = d.Get("signature_secret").(string)

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifySignature(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	time.Sleep(3 * time.Second)
	return resourceAliyunApigatewaySignatureRead(d, meta)
}

func resourceAliyunApigatewaySignatureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDeleteSignatureRequest()
	request.RegionId = client.RegionId
	request.SignatureId = d.Id()

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteSignature(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFoundSignature"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(cloudApiService.WaitForApiGatewaySignature(d.Id(), Deleted, DefaultTimeout))
}
