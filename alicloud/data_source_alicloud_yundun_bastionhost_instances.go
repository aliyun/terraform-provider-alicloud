package alicloud

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudBastionhostInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBastionhostInstancesRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_domain": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"license_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_network_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudBastionhostInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := yundun_bastionhost.CreateDescribeInstanceBastionhostRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_bastionhost.Instance

	// get name Regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	if v, ok := d.GetOk("ids"); ok {
		ids, _ := v.([]interface{})
		var ids_str []string
		for _, v_instance_id := range ids {
			ids_str = append(ids_str, v_instance_id.(string))
		}
		request.InstanceId = &ids_str
	}
	for {
		raw, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.DescribeInstanceBastionhost(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_bastionhost_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*yundun_bastionhost.DescribeInstanceBastionhostResponse)
		if len(response.Instances) < 1 {
			break
		}

		for _, e := range response.Instances {
			if nameRegex != nil && !nameRegex.MatchString(e.Description) {
				continue
			}
			instances = append(instances, e)
		}

		if len(response.Instances) < PageSizeSmall {
			break
		}

		currentPageNo := request.CurrentPage
		if page, err := getNextpageNumber(currentPageNo); err != nil {
			return WrapError(err)
		} else {
			request.CurrentPage = page
		}
	}

	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}
	if len(instanceIds) < 1 {
		return WrapError(extractBastionhostInstance(d, nil))
	}
	var specs []yundun_bastionhost.InstanceAttribute
	for _, instanceId := range instanceIds {
		request := yundun_bastionhost.CreateDescribeInstanceAttributeRequest()
		request.InstanceId = instanceId
		raw, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.DescribeInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_bastionhost_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		res, _ := raw.(*yundun_bastionhost.DescribeInstanceAttributeResponse)
		specs = append(specs, res.InstanceAttribute)
	}
	return WrapError(extractBastionhostInstance(d, specs))
}

func extractBastionhostInstance(d *schema.ResourceData, specs []yundun_bastionhost.InstanceAttribute) error {
	var instanceIds []string
	var descriptions []string
	var instances []map[string]interface{}

	for _, item := range specs {
		instanceMap := map[string]interface{}{
			"id":                    item.InstanceId,
			"description":           item.Description,
			"user_vswitch_id":       item.VswitchId,
			"private_domain":        item.IntranetEndpoint,
			"public_domain":         item.InternetEndpoint,
			"instance_status":       item.InstanceStatus,
			"license_code":          item.LicenseCode,
			"public_network_access": item.PublicNetworkAccess,
			"security_group_ids":    item.ReferredSecurityGroups,
		}
		instanceIds = append(instanceIds, item.InstanceId)
		descriptions = append(descriptions, item.Description)
		instances = append(instances, instanceMap)
	}

	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", instances); err != nil {
		return WrapError(err)
	}
	// storage locally
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), instances)
	}
	return nil
}
