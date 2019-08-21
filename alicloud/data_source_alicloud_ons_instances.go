package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudOnsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
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
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"release_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOnsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	request := ons.CreateOnsInstanceInServiceListRequest()
	request.RegionId = client.RegionId
	request.PreventCache = onsService.GetPreventCache()

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceInServiceList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.OnsInstanceInServiceListResponse)

	var filteredInstances []ons.InstanceVO
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, instance := range response.Data.InstanceVO {
			if r != nil && !r.MatchString(instance.InstanceName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[instance.InstanceId]; !ok {
					continue
				}
			}

			filteredInstances = append(filteredInstances, instance)
		}
	} else {
		filteredInstances = response.Data.InstanceVO
	}

	return onsInstancesDecriptionAttributes(d, filteredInstances, meta)
}

func onsInstancesDecriptionAttributes(d *schema.ResourceData, instancesInfo []ons.InstanceVO, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range instancesInfo {
		mapping := map[string]interface{}{
			"id":              item.InstanceId,
			"instance_id":     item.InstanceId,
			"instance_status": item.InstanceStatus,
			"release_time":    item.ReleaseTime,
			"instance_type":   item.InstanceType,
			"instance_name":   item.InstanceName,
		}

		names = append(names, item.InstanceName)
		ids = append(ids, item.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
