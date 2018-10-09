package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunApigatewayApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayApiCreate,
		Read:   resourceAliyunApigatewayApiRead,
		Update: resourceAliyunApigatewayApiUpdate,
		Delete: resourceAliyunApigatewayApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"request_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"request_parameters": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_parameters": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_parameters_map": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewayApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := buildAliyunApiArgs(d, meta)
	fmt.Println("_______" + d.Get("group_id").(string))
	fmt.Println(request.GroupId)
	resp, err := client.cloudapiconn.CreateApi(request)

	if err != nil {
		return fmt.Errorf("Creating apigatway api error: %#v", err)
	}

	d.SetId(resp.ApiId)

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cloudapi.CreateDescribeApiRequest()
	request.ApiId = d.Id()
	request.GroupId = d.Get("group_id").(string)
	resp, err := client.cloudapiconn.DescribeApi(request)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("api_id", resp.ApiId)
	d.Set("group_id", resp.GroupId)
	d.Set("name", resp.ApiName)
	d.Set("description", resp.Description)
	d.Set("auth_type", resp.AuthType)
	d.Set("request_config", resp.RequestConfig)
	d.Set("request_parameters", resp.RequestParameters)
	d.Set("service_parameters", resp.ServiceParameters)
	d.Set("service_parameters_map", resp.ServiceParametersMap)

	return nil
}

func resourceAliyunApigatewayApiUpdate(d *schema.ResourceData, meta interface{}) error {

	cloudapiconn := meta.(*AliyunClient).cloudapiconn

	if d.HasChange("groupName") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") || d.HasChange("description") {
		req := cloudapi.CreateModifyApiRequest()
		req.ApiId = d.Id()
		req.GroupId = d.Get("group_id").(string)
		req.ApiName = d.Get("name").(string)
		req.Description = d.Get("description").(string)
		req.AuthType = d.Get("auth_type").(string)
		req.RequestConfig = d.Get("request_config").(string)
		req.RequestParameters = d.Get("request_parameters").(string)
		req.ServiceParameters = d.Get("service_parameters").(string)
		req.ServiceParametersMap = d.Get("service_parameters_map").(string)

		if _, err := cloudapiconn.ModifyApi(req); err != nil {
			return fmt.Errorf("Modify Api got an error: %#v", err)
		}
	}
	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayApiDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	req := cloudapi.CreateDeleteApiRequest()
	req.GroupId = d.Get("group_id").(string)
	req.ApiId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.cloudapiconn.DeleteApi(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting Api failed: %#v", err))
		}

		if _, err := client.DescribeApi(d.Get("group_id").(string), d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing api failed when deleting api: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete api %s timeout.", d.Id()))
	})
}

func buildAliyunApiArgs(d *schema.ResourceData, meta interface{}) *cloudapi.CreateApiRequest {
	request := cloudapi.CreateCreateApiRequest()

	request.GroupId = d.Get("group_id").(string)
	request.Description = d.Get("description").(string)
	request.ApiName = d.Get("name").(string)
	request.AuthType = d.Get("auth_type").(string)

	request.RequestConfig = d.Get("request_config").(string)
	request.ServiceConfig = d.Get("service_config").(string)
	request.RequestParameters = d.Get("request_parameters").(string)
	request.ServiceParameters = d.Get("service_parameters").(string)
	request.ServiceParametersMap = d.Get("service_parameters_map").(string)

	request.ResultType = "JSON"
	request.ResultSample = "Result Sample"
	request.Visibility = "PRIVATE"
	request.AllowSignatureMethod = "HmacSHA256"
	request.WebSocketApiType = "COMMON"

	return request
}
