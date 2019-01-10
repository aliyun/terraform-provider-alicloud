package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudLogtailToMachineGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogtailToMachineGrouRead,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"offset": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  10,
			},
			// Computed value
			"machine_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"logtail_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAlicloudLogtailToMachineGrouRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	project := d.Get("project").(string)
	_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		_, err := slsClient.GetProject(project)
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("get logs got an error: %#v", err)
	}
	machine_group, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		machine_group, _, err := slsClient.ListMachineGroup(project, d.Get("offset").(int), d.Get("size").(int))
		if err != nil {
			return nil, err
		}
		return machine_group, err
	})
	logtail_config, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		logtail_config, _, err := slsClient.ListConfig(project, d.Get("offset").(int), d.Get("size").(int))
		if err != nil {
			return nil, err
		}
		return logtail_config, err
	})
	ids := []string{fmt.Sprintln(project)}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("machine_group", machine_group); err != nil {
		return nil
	}
	if err := d.Set("logtail_config", logtail_config); err != nil {
		return nil
	}
	var project_detail = map[string]interface{}{"project_machine": machine_group}
	project_detail["project_logtail"] = logtail_config
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), project_detail)
	}
	return nil
}
