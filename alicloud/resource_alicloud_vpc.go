package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpcCreate,
		Read:   resourceAliyunVpcRead,
		Update: resourceAliyunVpcUpdate,
		Delete: resourceAliyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"tags": tagsSchema(),
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	var response map[string]interface{}
	action := "CreateVpc"
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"CidrBlock":   d.Get("cidr_block").(string),
		"ClientToken": buildClientToken("CreateVpc"),
	}

	if v := d.Get("name").(string); v != "" {
		request["VpcName"] = v
	}

	if v := d.Get("description").(string); v != "" {
		request["Description"] = v
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request["ResourceGroupId"] = v
	}

	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	// If the API supports
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["VpcId"]))

	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliyunVpcUpdate(d, meta)
}

func resourceAliyunVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpcWithTeadsl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", object["CidrBlock"])
	d.Set("name", object["VpcName"])
	d.Set("description", object["Description"])
	d.Set("router_id", object["VRouterId"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	tags, err := vpcService.ListTagResources(d.Id(), "VPC")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}

	// Retrieve all route tables and filter to get system
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	action := "DescribeRouteTables"
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"VpcId":           d.Id(),
		"VRouterId":       object["VRouterId"],
		"ResourceGroupId": object["ResourceGroupId"],
		"PageNumber":      1,
		"PageSize":        PageSizeLarge,
	}
	var routeTabls []interface{}
	for {
		total := 0
		err = resource.Retry(6*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil && IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			addDebug(action, response, request)

			v, err := jsonpath.Get("$.RouteTables.RouteTable", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.RouteTables.RouteTable", response))
			}

			routeTabls = append(routeTabls, v.([]interface{})...)
			total = len(v.([]interface{}))
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if total < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	// Generally, the system route table is the last one
	for i := len(routeTabls) - 1; i >= 0; i-- {
		object := routeTabls[i].(map[string]interface{})
		if object["RouteTableType"] == "System" {
			d.Set("router_table_id", object["RouteTableId"])
			d.Set("route_table_id", object["RouteTableId"])
		}
	}
	return nil
}

func resourceAliyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()

	if err != nil {
		return WrapError(err)
	}
	if err := vpcService.setInstanceTags(d, TagResourceVpc); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunVpcRead(d, meta)
	}
	attributeUpdate := false
	action := "ModifyVpcAttribute"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    d.Id(),
	}

	if d.HasChange("name") {
		request["VpcName"] = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		request["Description"] = d.Get("description").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
	}

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	action := "DeleteVpc"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    d.Id(),
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Pending"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
