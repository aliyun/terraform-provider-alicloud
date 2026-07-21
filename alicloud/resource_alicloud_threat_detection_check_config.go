package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionCheckConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionCheckConfigCreate,
		Read:   resourceAliCloudThreatDetectionCheckConfigRead,
		Update: resourceAliCloudThreatDetectionCheckConfigUpdate,
		Delete: resourceAliCloudThreatDetectionCheckConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"configure": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cycle_days": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"enable_add_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"selected_checks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"section_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"system_config": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vendors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudThreatDetectionCheckConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ChangeCheckConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("cycle_days"); ok {
		cycleDaysMapsArray := convertToInterfaceArray(v)

		request["CycleDays"] = cycleDaysMapsArray
	}

	if v, ok := d.GetOkExists("enable_auto_check"); ok {
		request["EnableAutoCheck"] = v
	}

	if v, ok := d.GetOk("system_config"); ok {
		request["SystemConfig"] = v
	}
	if v, ok := d.GetOk("configure"); ok {
		request["Configure"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("vendors"); ok {
		vendorsMapsArray := convertToInterfaceArray(v)

		request["Vendors"] = vendorsMapsArray
	}

	if v, ok := d.GetOkExists("enable_add_check"); ok {
		request["EnableAddCheck"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}

	if _, ok := d.GetOk("selected_checks"); ok {
		threatDetectionServiceV2 := ThreatDetectionServiceV2{client}
		objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionCheckConfig(d.Id())
		if err != nil {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_check_config DescribeThreatDetectionCheckConfig Failed!!! %s", err)
		}
		selectedChecksListRaw := objectRaw["SelectedChecks"]
		if selectedChecksListRaw != nil {
			removedCheckList := make([]interface{}, 0)
			for _, selectedChecksChildRaw := range selectedChecksListRaw.([]interface{}) {
				attackPathAssetListChildRaw := selectedChecksChildRaw.(map[string]interface{})
				removedCheckList = append(removedCheckList, map[string]interface{}{
					"CheckId":   attackPathAssetListChildRaw["CheckId"],
					"SectionId": attackPathAssetListChildRaw["SectionId"],
				})
			}
			request["RemovedCheck"] = removedCheckList
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if _, ok := d.GetOk("selected_checks"); ok {
		_, newEntry := d.GetChange("selected_checks")
		added := newEntry

		if len(added.([]interface{})) > 0 {
			addedCheckList := make([]interface{}, 0)
			for _, item := range added.([]interface{}) {
				itemMap := item.(map[string]interface{})
				addedCheckList = append(addedCheckList, map[string]interface{}{
					"CheckId":   itemMap["check_id"],
					"SectionId": itemMap["section_id"],
				})
			}
			request["AddedCheck"] = addedCheckList
		}
		request["RemovedCheck"] = nil
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
	}

	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_check_config", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudThreatDetectionCheckConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionCheckConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionCheckConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_check_config DescribeThreatDetectionCheckConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable_add_check", objectRaw["EnableAddCheck"])
	d.Set("enable_auto_check", objectRaw["EnableAutoCheck"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("start_time", objectRaw["StartTime"])

	selectedChecksListRaw := objectRaw["SelectedChecks"]
	selectedChecksListMaps := make([]map[string]interface{}, 0)
	if selectedChecksListRaw != nil {
		for _, selectedChecksChildRaw := range selectedChecksListRaw.([]interface{}) {
			selectedChecksListMap := make(map[string]interface{})
			attackPathAssetListChildRaw := selectedChecksChildRaw.(map[string]interface{})
			selectedChecksListMap["check_id"] = attackPathAssetListChildRaw["CheckId"]
			selectedChecksListMap["section_id"] = attackPathAssetListChildRaw["SectionId"]

			selectedChecksListMaps = append(selectedChecksListMaps, selectedChecksListMap)
		}
	}
	d.Set("selected_checks", selectedChecksListMaps)

	cycleDaysRaw := make([]interface{}, 0)
	if objectRaw["CycleDays"] != nil {
		cycleDaysRaw = convertToInterfaceArray(objectRaw["CycleDays"])
	}

	d.Set("cycle_days", cycleDaysRaw)

	return nil
}

func resourceAliCloudThreatDetectionCheckConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAliCloudThreatDetectionCheckConfigCreate(d, meta)
}

func resourceAliCloudThreatDetectionCheckConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Check Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
