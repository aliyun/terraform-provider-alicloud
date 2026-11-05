package polardb

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceApplicationEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationEndpointCreate,
		Read:   resourceApplicationEndpointRead,
		Delete: resourceApplicationEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Public"}, false),
			},
		},
	}
}

func resourceApplicationEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}

	action := "CreateApplicationEndpointAddress"
	request := map[string]interface{}{
		"ApplicationId": d.Get("application_id").(string),
		"EndpointId":    d.Get("endpoint_id").(string),
		"NetType":       d.Get("net_type").(string),
	}

	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		logs.AddDebug(action, response, request)
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_application_endpoint", action, errs.AlibabaCloudSdkGoERROR)
	}
	applicationId, _ := response["ApplicationId"].(string)
	if applicationId == "" {
		applicationId = d.Get("application_id").(string)
	}
	stateConf := helper.BuildStateConf([]string{"NetCreating"}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(applicationId, []string{""}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, applicationId)
	}
	d.SetId(applicationId)

	return resourceApplicationEndpointRead(d, meta)
}

func resourceApplicationEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	applicationAttribute, err := polarDBService.DescribePolarDBApplicationAttribute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := applicationAttribute["ApplicationId"].(string); ok {
		d.Set("application_id", v)
	}
	endpoints, _ := applicationAttribute["Endpoints"].([]interface{})
	if len(endpoints) == 0 {
		d.SetId("")
		return nil
	}
	if ep, ok := endpoints[0].(map[string]interface{}); ok {
		if epId, ok := ep["EndpointId"].(string); ok {
			d.Set("endpoint_id", epId)
		}
	}

	return nil
}

func resourceApplicationEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	action := "DeleteApplicationEndpointAddress"
	request := map[string]interface{}{
		"ApplicationId": d.Get("application_id").(string),
		"EndpointId":    d.Get("endpoint_id").(string),
		"NetType":       d.Get("net_type").(string),
	}
	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		logs.AddDebug(action, response, request)
		if errs.IsExpectedErrors(err, []string{"Endpoint.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_application_endpoint", action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"NetDeleting"}, []string{"Activated"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{""}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
