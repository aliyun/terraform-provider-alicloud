package alicloud

import (
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksRemind() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksRemindRead,
		Schema: map[string]*schema.Schema{
			"remind_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remind_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remind_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remind_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnd_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnd_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert_methods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"baselines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"baseline_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"alert_targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"useflag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"biz_processes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biz_process_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"max_alert_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alert_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"robots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"webhooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"founder": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudDataWorksRemindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataworksPublicService := DataworksPublicService{client}
	remindId := d.Get("remind_id").(string)
	object, err := dataworksPublicService.DescribeDataWorksRemind(remindId)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] dataSource alicloud_data_works_remind DescribeDataWorksRemind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(fmt.Sprint(object["remindId"]))
	d.Set("remind_id", fmt.Sprint(object["remindId"]))
	d.Set("remind_name", object["remindName"])
	d.Set("remind_unit", object["remindUnit"])
	d.Set("remind_type", object["remindType"])
	d.Set("alert_unit", object["alertUnit"])
	d.Set("dnd_end", object["dndEnd"])
	d.Set("dnd_start", object["dndStart"])
	d.Set("founder", object["founder"])
	d.Set("detail", object["detail"])
	d.Set("max_alert_times", object["maxAlertTimes"])
	d.Set("alert_interval", object["alertInterval"])
	d.Set("useflag", object["useflag"])
	d.Set("region_id", object["regionId"])

	if v, err := jsonpath.Get("$.alertMethods", object); err == nil {
		d.Set("alert_methods", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.alertTargets", object); err == nil {
		d.Set("alert_targets", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.webhooks", object); err == nil {
		d.Set("webhooks", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.nodes", object); err == nil {
		d.Set("nodes", flattenRemindNodes(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.baselines", object); err == nil {
		d.Set("baselines", flattenRemindBaselines(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.bizProcesses", object); err == nil {
		d.Set("biz_processes", flattenRemindBizProcesses(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.robots", object); err == nil {
		d.Set("robots", flattenRemindRobots(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.projects", object); err == nil {
		d.Set("projects", flattenRemindProjects(v.([]interface{})))
	}

	return nil
}
