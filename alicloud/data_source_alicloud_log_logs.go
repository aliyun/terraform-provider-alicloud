package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudLogLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogLogsRead,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"from": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"to": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"offset": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"lines": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  100,
			},
			"reverse": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"output_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceAlicloudLogLogsRead performs the Alicloud logs lookup.
func dataSourceAlicloudLogLogsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var from int64
	var to int64
	if f, ok := d.GetOk("from"); ok && f.(int) > 0 {
		from = int64(f.(int))
	} else {
		from = time.Now().Unix() - 15*60
	}

	if t, ok := d.GetOk("to"); ok && t.(int) > 0 {
		to = int64(t.(int))
	} else {
		to = time.Now().Unix()
	}

	result, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		result, err := slsClient.GetLogs(
			d.Get("project").(string),
			d.Get("logstore").(string),
			d.Get("topic").(string),
			from,
			to,
			d.Get("query").(string),
			(int64)(d.Get("lines").(int)),
			(int64)(d.Get("offset").(int)),
			d.Get("reverse").(bool),
		)
		return result, err
	})
	if err != nil {
		return fmt.Errorf("get logs got an error: %#v", err)
	}

	return logLogsDescriptionAttributes(d, result.(*sls.GetLogsResponse), meta, from, to)
}

// populate the numerous fields that the image description returns.
func logLogsDescriptionAttributes(d *schema.ResourceData, result *sls.GetLogsResponse, meta interface{}, from int64, to int64) error {

	ids := []string{
		fmt.Sprintln(
			d.Get("project").(string),
			d.Get("logstore").(string),
			d.Get("topic").(string),
			from,
			to,
			d.Get("query").(string),
			d.Get("lines").(int),
			d.Get("offset").(int),
			d.Get("reverse").(bool),
			time.Now(),
		)}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("progress", result.Progress); err != nil {
		return err
	}
	if err := d.Set("logs", result.Logs); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), result.Logs)
	}
	return nil
}
