package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEaisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEaisInstanceCreate,
		Read:   resourceAliCloudEaisInstanceRead,
		Update: resourceAliCloudEaisInstanceUpdate,
		Delete: resourceAliCloudEaisInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_var": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"image": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"eais.ei-a6.2xlarge", "eais.ei-a6.medium"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'force' is deprecated and will be removed in a future release.",
			},
		},
	}
}

func resourceAliCloudEaisInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if InArray(fmt.Sprint(d.Get("category")), []string{"eais", ""}) {
		action := "CreateEai"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
		if v, ok := d.GetOk("tags"); ok {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			request = expandTagsToMap(request, tagsMap)
		}

		request["InstanceType"] = d.Get("instance_type")
		if v, ok := d.GetOk("instance_name"); ok {
			request["InstanceName"] = v
		}
		request["SecurityGroupId"] = d.Get("security_group_id")
		request["VSwitchId"] = d.Get("vswitch_id")
		if v, ok := d.GetOk("image"); ok {
			request["Image"] = v
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eais_instance", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprint(response["ElasticAcceleratedInstanceId"]))

		eaisServiceV2 := EaisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, eaisServiceV2.EaisInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	if v, ok := d.GetOk("category"); ok && InArray(fmt.Sprint(v), []string{"jupyter"}) {
		action := "CreateEaiJupyter"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
		if v, ok := d.GetOk("tags"); ok {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			request = expandTagsToMap(request, tagsMap)
		}

		request["SecurityGroupId"] = d.Get("security_group_id")
		request["VSwitchId"] = d.Get("vswitch_id")
		if v, ok := d.GetOk("instance_name"); ok {
			request["EaisName"] = v
		}
		if v, ok := d.GetOk("environment_var"); ok {
			environmentVarMapsArray := make([]interface{}, 0)
			for _, dataLoop1 := range v.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Value"] = dataLoop1Tmp["value"]
				dataLoop1Map["Key"] = dataLoop1Tmp["key"]
				environmentVarMapsArray = append(environmentVarMapsArray, dataLoop1Map)
			}
			environmentVarMapsJson, err := json.Marshal(environmentVarMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["EnvironmentVar"] = string(environmentVarMapsJson)
		}

		request["EaisType"] = d.Get("instance_type")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eais_instance", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprint(response["ElasticAcceleratedInstanceId"]))

		eaisServiceV2 := EaisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, eaisServiceV2.EaisInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	if v, ok := d.GetOk("category"); ok && InArray(fmt.Sprint(v), []string{"ei"}) {
		action := "CreateEaisEi"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
		if v, ok := d.GetOk("tags"); ok {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			request = expandTagsToMap(request, tagsMap)
		}

		if v, ok := d.GetOk("instance_name"); ok {
			request["InstanceName"] = v
		}
		request["InstanceType"] = d.Get("instance_type")
		request["VSwitchId"] = d.Get("vswitch_id")
		request["SecurityGroupId"] = d.Get("security_group_id")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eais_instance", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprint(response["EiInstanceId"]))

		eaisServiceV2 := EaisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Bindable"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, eaisServiceV2.EaisInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	return resourceAliCloudEaisInstanceUpdate(d, meta)
}

func resourceAliCloudEaisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eaisServiceV2 := EaisServiceV2{client}

	objectRaw, err := eaisServiceV2.DescribeEaisInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eais_instance DescribeEaisInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("security_group_id", objectRaw["SecurityGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEaisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if d.HasChange("status") {
		eaisServiceV2 := EaisServiceV2{client}
		object, err := eaisServiceV2.DescribeEaisInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			enableAction := false
			if v, ok := d.GetOk("category"); ok && InArray(fmt.Sprint(v), []string{"jupyter"}) {
				enableAction = true
			}
			if target == "Stopped" && enableAction {
				action := "StopEaiJupyter"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
				eaisServiceV2 := EaisServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, eaisServiceV2.EaisInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "InUse" && enableAction {
				action := "StartEaiJupyter"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
				eaisServiceV2 := EaisServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, eaisServiceV2.EaisInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
	}

	if d.HasChange("tags") {
		eaisServiceV2 := EaisServiceV2{client}
		if err := eaisServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudEaisInstanceRead(d, meta)
}

func resourceAliCloudEaisInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("category"); ok {
		if InArray(fmt.Sprint(v), []string{"ei"}) {
			enableDelete = false
			log.Printf("[WARN] Cannot destroy resource alicloud_eais_instance which category valued ei. Terraform will remove this resource from the state file, however resources may remain.")
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := "DeleteEai"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["ElasticAcceleratedInstanceId"] = d.Id()
		request["RegionId"] = client.RegionId

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}

	enableDelete = false
	if v, ok := d.GetOk("category"); ok {
		if InArray(fmt.Sprint(v), []string{"ei"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := "DeleteEaisEi"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["EiInstanceId"] = d.Id()
		request["RegionId"] = client.RegionId

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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
			if IsExpectedErrors(err, []string{"InvalidParameter.InstanceId.NotFound"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}
