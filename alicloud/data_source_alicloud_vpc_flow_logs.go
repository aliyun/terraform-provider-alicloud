package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVpcFlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcFlowLogsRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "VSwitch", "NetworkInterface"}, false),
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"All", "Allow", "Drop"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
				Default:      "Active",
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"flow_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcFlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeFlowLogsRequest()
	request.RegionId = client.RegionId

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	if v := d.Get("log_store_name").(string); v != "" {
		request.LogStoreName = v
	}
	if v := d.Get("project_name").(string); v != "" {
		request.ProjectName = v
	}
	if v := d.Get("resource_id").(string); v != "" {
		request.ResourceId = v
	}
	if v := d.Get("resource_type").(string); v != "" {
		request.ResourceType = v
	}
	if v := d.Get("traffic_type").(string); v != "" {
		request.TrafficType = v
	}
	if v := d.Get("status").(string); v != "" {
		request.Status = v
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []vpc.FlowLog
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeFlowLogs(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_flowlogs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.DescribeFlowLogsResponse)

		for _, item := range response.FlowLogs.FlowLog {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.FlowLogName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.FlowLogId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.FlowLogs.FlowLog) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, flowLog := range objects {
		mapping := map[string]interface{}{
			"id":             flowLog.FlowLogId,
			"flow_log_name":  flowLog.FlowLogName,
			"description":    flowLog.Description,
			"creation_time":  flowLog.CreationTime,
			"project_name":   flowLog.ProjectName,
			"log_store_name": flowLog.LogStoreName,
			"region_id":      flowLog.RegionId,
		}
		ids = append(ids, flowLog.FlowLogId)
		names = append(names, flowLog.FlowLogName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("flow_logs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
