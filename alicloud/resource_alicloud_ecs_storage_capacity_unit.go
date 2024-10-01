package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsStorageCapacityUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsStorageCapacityUnitCreate,
		Read:   resourceAliCloudEcsStorageCapacityUnitRead,
		Update: resourceAliCloudEcsStorageCapacityUnitUpdate,
		Delete: resourceAliCloudEcsStorageCapacityUnitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntInSlice([]int{20, 40, 100, 200, 500, 1024, 2048, 5120, 10240, 20480, 51200}),
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 5, 6}),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"storage_capacity_unit_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with 'http://', 'https://'."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9:_-]{1,127}$`), `The name must start with a letter. It must be 2 to 128 characters in length. It can contain digits, colons (:), underscores (_), and hyphens (-).`)),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEcsStorageCapacityUnitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "PurchaseStorageCapacityUnit"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["Amount"] = 1
	request["ClientToken"] = buildClientToken("PurchaseStorageCapacityUnit")
	request["Capacity"] = d.Get("capacity")

	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOk("storage_capacity_unit_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_storage_capacity_unit", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.StorageCapacityUnitIds", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_storage_capacity_unit")
	} else {
		storageCapacityUnitId := resp.(map[string]interface{})["StorageCapacityUnitId"].([]interface{})[0]
		d.SetId(fmt.Sprint(storageCapacityUnitId))
	}

	time.Sleep(5 * time.Minute)
	stateConf := BuildStateConf([]string{}, []string{"Pending", "Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsStorageCapacityUnitStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsStorageCapacityUnitRead(d, meta)
}

func resourceAliCloudEcsStorageCapacityUnitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsStorageCapacityUnit(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_storage_capacity_unit ecsService.DescribeEcsStorageCapacityUnit Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if capacity, ok := object["Capacity"]; ok && fmt.Sprint(capacity) != "0" {
		d.Set("capacity", formatInt(capacity))
	}

	d.Set("start_time", object["StartTime"])
	d.Set("storage_capacity_unit_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudEcsStorageCapacityUnitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"StorageCapacityUnitId": d.Id(),
	}

	if d.HasChange("storage_capacity_unit_name") {
		update = true
	}
	if v, ok := d.GetOk("storage_capacity_unit_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "ModifyStorageCapacityUnitAttribute"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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

	return resourceAliCloudEcsStorageCapacityUnitRead(d, meta)
}

func resourceAliCloudEcsStorageCapacityUnitDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAliCloudEcsStorageCapacityUnit. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
