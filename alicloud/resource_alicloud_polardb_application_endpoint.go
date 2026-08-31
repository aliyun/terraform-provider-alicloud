package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBApplicationEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBApplicationEndpointCreate,
		Read:   resourceAlicloudPolarDBApplicationEndpointRead,
		Delete: resourceAlicloudPolarDBApplicationEndpointDelete,
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
				ValidateFunc: StringInSlice([]string{"Public"}, false),
			},
		},
	}
}

func resourceAlicloudPolarDBApplicationEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	action := "CreateApplicationEndpointAddress"
	request := map[string]interface{}{
		"ApplicationId": d.Get("application_id").(string),
		"EndpointId":    d.Get("endpoint_id").(string),
		"NetType":       d.Get("net_type").(string),
	}

	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		addDebug(action, response, request)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_application_endpoint", action, AlibabaCloudSdkGoERROR)
	}
	applicationId, _ := response["ApplicationId"].(string)
	if applicationId == "" {
		applicationId = d.Get("application_id").(string)
	}
	stateConf := BuildStateConf([]string{"NetCreating"}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(applicationId, []string{""}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, applicationId)
	}
	d.SetId(applicationId)

	return resourceAlicloudPolarDBApplicationEndpointRead(d, meta)
}

func resourceAlicloudPolarDBApplicationEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	applicationAttribute, err := polarDBService.DescribePolarDBApplicationAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAlicloudPolarDBApplicationEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	action := "DeleteApplicationEndpointAddress"
	request := map[string]interface{}{
		"ApplicationId": d.Get("application_id").(string),
		"EndpointId":    d.Get("endpoint_id").(string),
		"NetType":       d.Get("net_type").(string),
	}
	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		addDebug(action, response, request)
		if IsExpectedErrors(err, []string{"Endpoint.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_application_endpoint", action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"NetDeleting"}, []string{"Activated"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, polarDBService.PolarDBApplicationStateRefreshFunc(d.Id(), []string{""}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
