package alicloud

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const defaultClusterEventPageSize int64 = 100

func dataSourceAlicloudCSClusterLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSClusterLogsRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "log",
				ValidateFunc: validation.StringInSlice([]string{"log", "event", "task"}, false),
			},
			"entries": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// more filter for logs
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudCSClusterLogsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapError(err)
	}
	csClient := CsClient{client}

	clusterId := d.Get("cluster_id").(string)
	logType := d.Get("type").(string)

	entries := 100
	if v, ok := d.GetOk("entries"); ok {
		entries = v.(int)
	}

	switch logType {
	case "event":
		return csClusterEventsDescriptionAttributes(d, clusterId, entries, csClient)
	case "task":
		return csClusterTaskDescriptionAttributes(d, clusterId, entries, csClient)
	default:
		return csClusterLogsDescriptionAttributes(d, clusterId, entries, csClient)
	}
}

func csClusterEventsDescriptionAttributes(d *schema.ResourceData, clusterId string, entries int, csClient CsClient) error {
	attrLogs := make([]map[string]interface{}, 0, entries)

	totalPage := int64(entries) / defaultClusterEventPageSize
	totalPage += 1
	pageNumber := int64(1)
	cnt := 0
	for pageNumber <= totalPage {
		if entries <= cnt {
			break
		}

		events, err := csClient.DescribeClusterEvents(clusterId, defaultClusterEventPageSize, pageNumber)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_cluster_logs", "DescribeClusterEvents", AlibabaCloudSdkGoERROR)
		}
		if len(events.Events) <= 0 {
			break
		}

		for _, event := range events.Events {
			data := structs.Map(event)
			attrLogs = append(attrLogs, data)
			cnt += 1
		}
		pageNumber += 1
	}

	d.Set("entries", cnt)
	d.SetId(dataResourceIdHash([]string{clusterId}))

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), attrLogs)
	}

	return nil
}

func csClusterTaskDescriptionAttributes(d *schema.ResourceData, clusterId string, entries int, csClient CsClient) error {
	attrLogs := make([]map[string]interface{}, 0, entries)
	cnt := 0

	tasks, err := csClient.DescribeClusterTasks(clusterId)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_cluster_logs", "DescribeClusterTasks", AlibabaCloudSdkGoERROR)
	}

	for _, task := range tasks.Tasks {
		if entries <= cnt {
			break
		}
		data := structs.Map(task)
		attrLogs = append(attrLogs, data)
		cnt += 1
	}

	d.Set("entries", cnt)
	d.SetId(dataResourceIdHash([]string{clusterId}))

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), attrLogs)
	}

	return nil
}

func csClusterLogsDescriptionAttributes(d *schema.ResourceData, clusterId string, entries int, csClient CsClient) error {
	attrLogs := make([]string, 0, entries)
	cnt := 0

	logs, err := csClient.DescribeClusterLogs(clusterId)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_cluster_logs", "DescribeClusterLogs", AlibabaCloudSdkGoERROR)
	}

	for _, log := range logs {
		if entries <= cnt {
			break
		}

		attrLogs = append(attrLogs, tea.StringValue(log.ClusterLog))
		cnt += 1
	}

	d.Set("entries", cnt)
	d.SetId(dataResourceIdHash([]string{clusterId}))

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), attrLogs)
	}

	return nil
}
