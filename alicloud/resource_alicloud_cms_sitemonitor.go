package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCmsSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsSiteMonitorCreate,
		Read:   resourceAlicloudCmsSiteMonitorRead,
		Update: resourceAlicloudCmsSiteMonitorUpdate,
		Delete: resourceAlicloudCmsSiteMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{SiteMonitorHTTP, SiteMonitorDNS, SiteMonitorFTP, SiteMonitorPOP3,
					SiteMonitorPing, SiteMonitorSMTP, SiteMonitorTCP, SiteMonitorUDP}, false),
			},
			"alert_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 15}),
			},
			"isp_cities": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsSiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreateCreateSiteMonitorRequest()
	request.Address = d.Get("address").(string)
	request.AlertIds = d.Get("alert_ids").(string)
	request.TaskName = d.Get("task_name").(string)
	request.TaskType = d.Get("task_type").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.IspCities = d.Get("isp_cities").(string)
	request.OptionsJson = d.Get("options_json").(string)

	_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.CreateSiteMonitor(request)
	})
	if err != nil {
		return fmt.Errorf("Creating site monitor got an error: %#v", err)
	}

	return resourceAlicloudCmsSiteMonitorRead(d, meta)
}

func resourceAlicloudCmsSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreateDescribeSiteMonitorListRequest()
	request.Keyword = d.Get("task_name").(string)
	request.Scheme = "https"

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorList(request)
	})
	if err != nil {
		return fmt.Errorf("Error calling DescribeSiteMonitorList: %#v", err)
	}
	response, _ := raw.(*cms.DescribeSiteMonitorListResponse)
	if len(response.SiteMonitors.SiteMonitor) == 0 {
		return fmt.Errorf("Task '%s' does not exist: %#v", request.Keyword, err)
	}
	siteMonitor := response.SiteMonitors.SiteMonitor[0]

	d.SetId(siteMonitor.TaskId)
	d.Set("address", siteMonitor.Address)
	d.Set("task_name", siteMonitor.TaskName)
	d.Set("task_type", siteMonitor.TaskType)
	d.Set("task_state", siteMonitor.TaskState)
	d.Set("interval", siteMonitor.Interval)
	d.Set("options_json", siteMonitor.OptionsJson)
	d.Set("create_time", siteMonitor.CreateTime)
	d.Set("update_time", siteMonitor.UpdateTime)

	return nil
}

func resourceAlicloudCmsSiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	request := cms.CreateModifySiteMonitorRequest()
	request.TaskId = d.Id()
	request.Address = d.Get("address").(string)
	request.AlertIds = d.Get("alert_ids").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.IspCities = d.Get("isp_cities").(string)
	request.OptionsJson = d.Get("options_json").(string)
	request.TaskName = d.Get("task_name").(string)

	_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.ModifySiteMonitor(request)
	})
	if err != nil {
		return fmt.Errorf("ModifySiteMonitor got an error: %#v", err)
	}

	d.SetPartial("address")
	d.SetPartial("alert_ids")
	d.SetPartial("interval")
	d.SetPartial("isp_cities")
	d.SetPartial("options_json")

	d.Partial(false)

	return resourceAlicloudCmsSiteMonitorRead(d, meta)
}

func resourceAlicloudCmsSiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cms.CreateDeleteSiteMonitorsRequest()

	request.TaskIds = d.Id()
	request.IsDeleteAlarms = "false"
	taskName := d.Get("task_name").(string)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteSiteMonitors(request)
		})

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting Site Monitor got an error: %#v", err))
		}

		listRequest := cms.CreateDescribeSiteMonitorListRequest()
		listRequest.Keyword = taskName
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorList(listRequest)
		})
		list := raw.(*cms.DescribeSiteMonitorListResponse)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Describe site monitor got an error: %#v", err))
		}
		if list.TotalCount == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Deleting site monitor got an error: %#v", err))
	})
}
