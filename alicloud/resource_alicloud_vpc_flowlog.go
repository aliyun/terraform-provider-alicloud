package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcFlowLogCreate,
		Read:   resourceAlicloudVpcFlowLogRead,
		Update: resourceAlicloudVpcFlowLogUpdate,
		Delete: resourceAlicloudVpcFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "VSwitch", "NetworkInterface"}, false),
			},
			"traffic_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"All", "Allow", "Drop"}, false),
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
				Default:      "Active",
			},
		},
	}
}

func resourceAlicloudVpcFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateCreateFlowLogRequest()

	request.ResourceId = d.Get("resource_id").(string)
	request.ResourceType = d.Get("resource_type").(string)
	request.TrafficType = d.Get("traffic_type").(string)
	request.RegionId = client.RegionId
	request.LogStoreName = d.Get("log_store_name").(string)
	request.ProjectName = d.Get("project_name").(string)
	if v, ok := d.GetOk("flow_log_name"); ok && v.(string) != "" {
		request.FlowLogName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	var response *vpc.CreateFlowLogResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateFlowLog(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", "TOKEN_PROCESSING", "OperationConflict", Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateFlowLogResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_flow_log", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.FlowLogId)
	if err := vpcService.WaitForVpcFlowLog(d.Id(), Active, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudVpcFlowLogUpdate(d, meta)
}

func resourceAlicloudVpcFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*connectivity.AliyunClient)}
	object, err := vpcService.DescribeVpcFlowlog(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_id", object.ResourceId)
	d.Set("resource_type", object.ResourceType)
	d.Set("traffic_type", object.TrafficType)
	d.Set("flow_log_name", object.FlowLogName)
	d.Set("description", object.Description)
	d.Set("log_store_name", object.LogStoreName)
	d.Set("project_name", object.ProjectName)
	d.Set("status", object.Status)

	return nil
}

func resourceAlicloudVpcFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	d.Partial(true)

	update := false
	request := vpc.CreateModifyFlowLogAttributeRequest()

	request.FlowLogId = d.Id()
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if !d.IsNewResource() && d.HasChange("flow_log_name") {
		request.FlowLogName = d.Get("flow_log_name").(string)
		update = true
	}

	if update {
		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ModifyFlowLogAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("flow_log_name")
	}

	if d.HasChange("status") {
		object, err := vpcService.DescribeVpcFlowlog(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object.Status != target {
			if target == "Active" {
				request := vpc.CreateActiveFlowLogRequest()
				request.FlowLogId = d.Id()
				request.RegionId = client.RegionId
				err := resource.Retry(30*time.Second, func() *resource.RetryError {
					raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.ActiveFlowLog(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				if err := vpcService.WaitForVpcFlowLog(d.Id(), Active, DefaultTimeoutMedium); err != nil {
					return WrapError(err)
				}
			}
			if target == "Inactive" {
				request := vpc.CreateDeactiveFlowLogRequest()
				request.FlowLogId = d.Id()
				request.RegionId = client.RegionId
				err := resource.Retry(30*time.Second, func() *resource.RetryError {
					raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.DeactiveFlowLog(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				if err := vpcService.WaitForVpcFlowLog(d.Id(), Inactive, DefaultTimeoutMedium); err != nil {
					return WrapError(err)
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)

	return resourceAlicloudVpcFlowLogRead(d, meta)
}

func resourceAlicloudVpcFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteFlowLogRequest()

	request.FlowLogId = d.Id()
	request.RegionId = client.RegionId

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteFlowLog(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(vpcService.WaitForVpcFlowLog(d.Id(), Deleted, DefaultTimeoutMedium))
}
