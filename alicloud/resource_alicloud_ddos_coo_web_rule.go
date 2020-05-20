package alicloud

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDdosCooWebRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosCooWebRuleCreate,
		Read:   resourceAlicloudDdosCooWebRuleRead,
		Update: resourceAlicloudDdosCooWebRuleUpdate,
		Delete: resourceAlicloudDdosCooWebRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"proxy_types": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"real_servers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rs_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDdosCooWebRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateCreateWebRuleRequest()
	request.Domain = d.Get("domain").(string)
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	request.RsType = requests.NewInteger(d.Get("rs_type").(int))
	request.Rules = d.Get("rules").(string)
	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.CreateWebRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddos_coo_web_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", request.Domain))

	return resourceAlicloudDdosCooWebRuleUpdate(d, meta)
}
func resourceAlicloudDdosCooWebRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdosCooWebRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain", d.Id())
	d.Set("proxy_types", object.ProxyTypes)
	d.Set("real_servers", object.RealServers)
	return nil
}
func resourceAlicloudDdosCooWebRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := ddoscoo.CreateModifyWebRuleRequest()
	request.Domain = d.Id()
	if d.HasChange("proxy_types") {
		update = true
	}
	request.ProxyTypes = d.Get("proxy_types").(string)
	if d.HasChange("real_servers") {
		update = true
	}
	val := d.Get("real_servers").(string)
	realServer := make([]string, 0, 3)
	err := json.Unmarshal([]byte(val), &realServer)
	if err != nil {
		return WrapError(err)
	}
	request.RealServers = &realServer
	if !d.IsNewResource() && d.HasChange("rs_type") {
		update = true
	}
	request.RsType = requests.NewInteger(d.Get("rs_type").(int))
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request.ResourceGroupId = d.Get("resource_group_id").(string)
	}
	if update {
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ModifyWebRule(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudDdosCooWebRuleRead(d, meta)
}
func resourceAlicloudDdosCooWebRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ddoscoo.CreateDeleteWebRuleRequest()
	request.Domain = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DeleteWebRule(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
