package slb

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDomainExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainExtensionCreate,
		Read:   resourceDomainExtensionRead,
		Update: resourceDomainExtensionUpdate,
		Delete: resourceDomainExtensionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 65535),
				Required:     true,
				ForceNew:     true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_protection_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDomainExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateCreateDomainExtensionRequest()
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))
	request.Domain = d.Get("domain").(string)
	request.ServerCertificateId = d.Get("server_certificate_id").(string)

	var response *slb.CreateDomainExtensionResponse
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.CreateDomainExtension(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DomainExtensionProcessing"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*slb.CreateDomainExtensionResponse)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_domain_extension", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.DomainExtensionId)
	return resourceDomainExtensionRead(d, meta)
}

func resourceDomainExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	domainExtension, err := slbService.DescribeDomainExtensionAttribute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("load_balancer_id", domainExtension.LoadBalancerId)
	d.Set("domain", domainExtension.Domain)
	d.Set("server_certificate_id", domainExtension.ServerCertificateId)
	d.Set("frontend_port", domainExtension.ListenerPort)

	return nil
}

func resourceDomainExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("server_certificate_id") {
		request := slb.CreateSetDomainExtensionAttributeRequest()
		request.DomainExtensionId = d.Id()
		request.ServerCertificateId = d.Get("server_certificate_id").(string)
		client := meta.(*connectivity.AliyunClient)
		err := retry.Retry(3*time.Minute, func() *retry.RetryError {
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.SetDomainExtensionAttribute(request)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"BackendServer.configuring", "DomainExtensionProcessing"}) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_domain_extension", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceDomainExtensionRead(d, meta)
}

func resourceDomainExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	if d.Get("delete_protection_validation").(bool) {
		lbId := d.Get("load_balancer_id").(string)
		lbInstance, err := slbService.DescribeSlb(lbId)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapError(err)
		}
		if lbInstance.DeleteProtection == "on" {
			return errs.WrapError(fmt.Errorf("Current domain extension's SLB Instance %s has enabled DeleteProtection. Please set delete_protection_validation to false to delete the resource.", lbId))
		}
	}

	request := slb.CreateDeleteDomainExtensionRequest()
	request.DomainExtensionId = d.Id()

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteDomainExtension(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DomainExtensionProcessing", "InternalError"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)

	}
	return errs.WrapError(slbService.WaitForSlbDomainExtension(d.Id(), Deleted, DefaultTimeoutMedium))
}
