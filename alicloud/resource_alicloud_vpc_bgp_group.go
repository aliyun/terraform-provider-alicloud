package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcBgpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcBgpGroupCreate,
		Read:   resourceAliCloudVpcBgpGroupRead,
		Update: resourceAliCloudVpcBgpGroupUpdate,
		Delete: resourceAliCloudVpcBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_asn": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"local_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_fake_asn": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bgp_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{1,127}$`), "The name must be `2` to `128` characters in length and can contain digits, periods (.), underscores (_), and hyphens (-)."), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcBgpGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateBgpGroup"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBgpGroup")
	request["RouterId"] = d.Get("router_id")
	request["PeerAsn"] = d.Get("peer_asn")

	if v, ok := d.GetOk("local_asn"); ok {
		request["LocalAsn"] = v
	}

	if v, ok := d.GetOkExists("is_fake_asn"); ok {
		request["IsFakeAsn"] = v
	}

	if v, ok := d.GetOk("auth_key"); ok {
		request["AuthKey"] = v
	}

	if v, ok := d.GetOk("bgp_group_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_bgp_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BgpGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcBgpGroupRead(d, meta)
}

func resourceAliCloudVpcBgpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpcBgpGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_bgp_group vpcService.DescribeVpcBgpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("router_id", object["RouterId"])
	d.Set("peer_asn", formatInt(object["PeerAsn"]))
	d.Set("local_asn", formatInt(object["LocalAsn"]))
	d.Set("is_fake_asn", object["IsFake"])
	d.Set("auth_key", object["AuthKey"])
	d.Set("bgp_group_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudVpcBgpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("ModifyBgpGroupAttribute"),
		"BgpGroupId":  d.Id(),
	}

	if d.HasChange("peer_asn") {
		update = true

		request["PeerAsn"] = d.Get("peer_asn")
	}

	if d.HasChange("local_asn") {
		update = true

		if v, ok := d.GetOkExists("local_asn"); ok {
			request["LocalAsn"] = v
		}
	}

	if d.HasChange("is_fake_asn") {
		update = true

		if v, ok := d.GetOkExists("is_fake_asn"); ok {
			request["IsFakeAsn"] = v
		}
	}

	if d.HasChange("auth_key") {
		update = true
	}
	if v, ok := d.GetOk("auth_key"); ok {
		request["AuthKey"] = v
	}

	if d.HasChange("bgp_group_name") {
		update = true
	}
	if v, ok := d.GetOk("bgp_group_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyBgpGroupAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudVpcBgpGroupRead(d, meta)
}

func resourceAliCloudVpcBgpGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteBgpGroup"
	var response map[string]interface{}

	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("DeleteBgpGroup"),
		"BgpGroupId":  d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
