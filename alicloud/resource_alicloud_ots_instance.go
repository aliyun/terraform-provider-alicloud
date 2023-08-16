package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOtsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceCreate,
		Read:   resourceAliyunOtsInstanceRead,
		Update: resourceAliyunOtsInstanceUpdate,
		Delete: resourceAliyunOtsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},

			"accessed_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  AnyNetwork,
				ValidateFunc: validation.StringInSlice([]string{
					string(AnyNetwork), string(VpcOnly), string(VpcOrConsole),
				}, false),
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  OtsHighPerformance,
				ValidateFunc: validation.StringInSlice([]string{
					string(OtsCapacity), string(OtsHighPerformance),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	instanceTypeStr := d.Get("instance_type").(string)
	instanceType, err := parseAndCheckInstanceType(instanceTypeStr, otsService)
	if err != nil {
		return WrapError(err)
	}

	actionPath, instanceName, request := buildCreateInstanceRoaRequest(d, client.RegionId, instanceType)

	_, err = OtsRestApiPostWithRetry(d, client, "tablestore", "2020-12-09", actionPath, request)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(instanceName)
	if err := otsService.WaitForOtsInstance(instanceName, Running, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func parseAndCheckInstanceType(instanceTypeStr string, otsService OtsService) (string, error) {
	instanceType := convertInstanceType(OtsInstanceType(instanceTypeStr))
	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return "", WrapError(err)
	}
	valid := false
	for _, t := range types {
		if instanceType == t {
			valid = true
			break
		}
	}
	if valid {
		return instanceType, nil
	}
	return instanceType, WrapError(Error("The instance type %s is not available in the region %s.", instanceTypeStr, otsService.client.RegionId))

}

func buildCreateInstanceRoaRequest(d *schema.ResourceData, regionId string, instanceType string) (string, string, map[string]*string) {
	actionPath := "/v2/openapi/createinstance"
	request := make(map[string]*string)
	request["RegionId"] = StringPointer(regionId)
	request["ClusterType"] = StringPointer(instanceType)
	instanceName := d.Get("name").(string)
	request["InstanceName"] = StringPointer(instanceName)
	request["InstanceDescription"] = StringPointer(d.Get("description").(string))
	request["Network"] = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string))))
	return actionPath, instanceName, request
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.InstanceName)
	d.Set("accessed_by", convertInstanceAccessedByRevert(object.Network))
	d.Set("instance_type", convertInstanceTypeRevert(object.ClusterType))
	d.Set("description", object.Description)
	d.Set("tags", otsTagsToMap(object.TagInfos.TagInfo))
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("accessed_by") {
		actionPath := "/v2/openapi/updateinstance"
		request := make(map[string]*string)
		request["RegionId"] = StringPointer(client.RegionId)
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())
		request["Network"] = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string))))

		response, err := OtsRestApiPostWithRetry(d, client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("accessed_by")
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			request := ots.CreateDeleteTagsRequest()
			request.RegionId = client.RegionId
			request.InstanceName = d.Id()
			var tags []ots.DeleteTagsTagInfo
			for _, t := range remove {
				tags = append(tags, ots.DeleteTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.DeleteTags(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := ots.CreateInsertTagsRequest()
			request.RegionId = client.RegionId
			request.InstanceName = d.Id()
			var tags []ots.InsertTagsTagInfo
			for _, t := range create {
				tags = append(tags, ots.InsertTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.InsertTags(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
		d.SetPartial("tags")
	}
	if err := otsService.WaitForOtsInstance(d.Id(), Running, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actionPath := "/v2/openapi/deleteinstance"
	request := make(map[string]*string)
	request["RegionId"] = StringPointer(client.RegionId)
	// id is instanceName
	request["InstanceName"] = StringPointer(d.Id())

	_, err := OtsRestApiPostWithRetry(d, client, "tablestore", "2020-12-09", actionPath, request)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
	}

	otsService := OtsService{client}
	return WrapError(otsService.WaitForOtsInstance(d.Id(), Deleted, DefaultLongTimeout))
}
