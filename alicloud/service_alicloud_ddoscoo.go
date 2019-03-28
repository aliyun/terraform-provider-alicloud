package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DdoscooService struct {
	client *connectivity.AliyunClient
}

func (s *DdoscooService) DescribeDdoscooInstance(instanceId string) (v ddoscoo.DescribeInstancesResponse, err error) {
	request := ddoscoo.CreateDescribeInstancesRequest()
	request.InstanceIds = "[\"" + instanceId + "\"]"
	request.PageNo = "1"
	request.PageSize = "10"

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdoscooInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return err
		}

		resp, _ := raw.(*ddoscoo.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances) == 0 || resp.Instances[0].InstanceId != instanceId {
			return WrapErrorf(Error(GetNotFoundMessage("Ddoscoo Instance", instanceId)), NotFoundMsg, ProviderERROR)
		}

		v = *resp
		return nil
	})

	return
}

func (s *DdoscooService) DescribeDdoscooInstanceSpec(instanceId string) (v ddoscoo.DescribeInstanceSpecsResponse, err error) {
	request := ddoscoo.CreateDescribeInstanceSpecsRequest()
	request.InstanceIds = "[\"" + instanceId + "\"]"

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstanceSpecs(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdoscooInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return err
		}

		resp, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)
		if resp == nil || len(resp.InstanceSpecs) == 0 || resp.InstanceSpecs[0].InstanceId != instanceId {
			return WrapErrorf(Error(GetNotFoundMessage("Ddoscoo Instance", instanceId)), NotFoundMsg, ProviderERROR)
		}

		v = *resp
		return nil
	})

	return
}

func (s *DdoscooService) UpdateDdoscooInstanceName(instanceId string, name string) error {
	request := ddoscoo.CreateModifyInstanceRemarkRequest()
	request.InstanceId = instanceId
	request.Remark = name

	invoker := NewInvoker()
	err := invoker.Run(func() error {
		_, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ModifyInstanceRemark(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdoscooInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return err
		}

		return nil
	})

	return err
}

func UpdateBandwidth(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	domainCount, _ := d.GetChange("domain_count")
	odomainCount := domainCount.(string)

	serviceBandwidth, _ := d.GetChange("service_bandwidth")
	oserviceBandwidth := serviceBandwidth.(string)

	portCount, _ := d.GetChange("port_count")
	oportCount := portCount.(string)

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Upgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "DomainCount",
			Value: odomainCount,
		},
		{
			Code:  "PortCount",
			Value: oportCount,
		},
		{
			Code:  "ServiceBandwidth",
			Value: oserviceBandwidth,
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func DowngradeDomainCount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	serviceBandwidth, _ := d.GetChange("service_bandwidth")
	oserviceBandwidth := serviceBandwidth.(string)

	portCount, _ := d.GetChange("port_count")
	oportCount := portCount.(string)

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Downgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: oportCount,
		},
		{
			Code:  "ServiceBandwidth",
			Value: oserviceBandwidth,
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func UpgradeDomainCount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	serviceBandwidth, _ := d.GetChange("service_bandwidth")
	oserviceBandwidth := serviceBandwidth.(string)

	portCount, _ := d.GetChange("port_count")
	oportCount := portCount.(string)

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Upgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: oportCount,
		},
		{
			Code:  "ServiceBandwidth",
			Value: oserviceBandwidth,
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func DowngradePortCount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	serviceBandwidth, _ := d.GetChange("service_bandwidth")
	oserviceBandwidth := serviceBandwidth.(string)

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Downgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: oserviceBandwidth,
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func UpgradePortCount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	serviceBandwidth, _ := d.GetChange("service_bandwidth")
	oserviceBandwidth := serviceBandwidth.(string)

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Upgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: oserviceBandwidth,
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func DowngradeServiceBandwidth(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Downgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: d.Get("service_bandwidth").(string),
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func UpgradeServiceBandwidth(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"
	request.ModifyType = "Upgrade"
	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  "BaseBandwidth",
			Value: d.Get("base_bandwidth").(string),
		},
		{
			Code:  "Bandwidth",
			Value: d.Get("bandwidth").(string),
		},
		{
			Code:  "DomainCount",
			Value: d.Get("domain_count").(string),
		},
		{
			Code:  "PortCount",
			Value: d.Get("port_count").(string),
		},
		{
			Code:  "ServiceBandwidth",
			Value: d.Get("service_bandwidth").(string),
		},
		{
			Code:  "NormalQps",
			Value: "3000",
		},
	}

	if _, err := client.WithBssopenapiClient(d.Get("business_endpoint").(string), func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
