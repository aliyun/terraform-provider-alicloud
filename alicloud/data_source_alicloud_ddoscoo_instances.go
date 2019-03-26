package alicloud

import (
	"encoding/json"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDdoscoo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdoscooInstancesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_band_width": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band_width": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_band_width": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDdoscooInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = string(PageSizeLarge)
	request.PageNo = "1"
	request.InstanceIds = "[\"" + d.Get("id").(string) + "\"]"

	var instances []ddoscoo.Instance

	for {
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "ddoscoo_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*ddoscoo.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances) < 1 {
			break
		}

		for _, item := range resp.Instances {
			instances = append(instances, item)
		}

		if len(resp.Instances) < PageSizeLarge {
			break
		}

		currentPageNo, err := strconv.Atoi(request.PageNo)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "ddoscoo_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if page, err := getNextpageNumber(requests.NewInteger(currentPageNo)); err != nil {
			return WrapError(err)
		} else {
			request.PageNo = string(page)
		}
	}

	if len(instances) < 1 {
		var defaultSpec = ddoscoo.InstanceSpec{
			InstanceId:       "defaultInstanceSpec",
			BaseBandwidth:    30,
			ElasticBandwidth: 30,
			PortLimit:        50,
			DomainLimit:      50,
			BandwidthMbps:    100,
		}

		var defaultSpecs = []ddoscoo.InstanceSpec{defaultSpec}

		return WrapError(extractDdoscooInstance(d, defaultSpecs))
	}

	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}

	return WrapError(extractDdoscooInstance(d, describeInstanceInfo(instanceIds, meta)))
}

func describeInstanceInfo(instanceIds []string, meta interface{}) []ddoscoo.InstanceSpec {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateDescribeInstanceSpecsRequest()
	request.RegionId = client.RegionId
	instanceIdsStr, _ := json.Marshal(instanceIds)
	request.InstanceIds = string(instanceIdsStr)

	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeInstanceSpecs(request)
	})

	if err != nil {
		return nil
	}

	resp, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)

	return resp.InstanceSpecs
}

func extractDdoscooInstance(d *schema.ResourceData, instanceSpecs []ddoscoo.InstanceSpec) error {
	var instanceId string
	var s []map[string]interface{}

	for _, item := range instanceSpecs {
		mapping := map[string]interface{}{
			"id":                 item.InstanceId,
			"band_width":         strconv.Itoa(item.ElasticBandwidth),
			"base_band_width":    strconv.Itoa(item.BaseBandwidth),
			"service_band_width": strconv.Itoa(item.BandwidthMbps),
			"port_count":         strconv.Itoa(item.PortLimit),
			"domain_count":       strconv.Itoa(item.DomainLimit),
		}
		instanceId = item.InstanceId
		s = append(s, mapping)
	}

	d.SetId(instanceId)
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("id", instanceId); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
