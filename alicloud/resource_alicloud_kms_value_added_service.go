package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// kmsValueAddedServiceAlreadyHeldMsg explains an order that is accepted but never becomes an
// effective instance because the account already holds the service.
const kmsValueAddedServiceAlreadyHeldMsg = "This account may already hold this value added service - " +
	"it can be purchased only once while it is in effect, and a duplicate order is refunded immediately. " +
	"Check the KMS console and, if the service is already in effect, import it with " +
	"'terraform import alicloud_kms_value_added_service.<name> <instance_id>' instead of creating it."

// Default key rotation is the one value added service this resource supports, and it mints instance
// ids under a prefix of its own, so an existing instance of it can be recognised before ordering a
// duplicate that would only be refunded. The prefix is checked only for this service: the expert
// service shares the 'kst-' prefix with alicloud_kms_instance, so a prefix match there would refuse a
// create over an unrelated KMS instance, and a wrong refusal leaves the user no way forward.
const (
	// The status an instance reaches once the service is actually in effect. Shared by the create wait and
	// the already-held pre-flight so the two cannot disagree about what "held" means.
	kmsValueAddedServiceEffectiveStatus = "Normal"
)

func resourceAliCloudKmsValueAddedService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsValueAddedServiceCreate,
		Read:   resourceAliCloudKmsValueAddedServiceRead,
		Update: resourceAliCloudKmsValueAddedServiceUpdate,
		Delete: resourceAliCloudKmsValueAddedServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3),
			},
			"renew_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_added_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudKmsValueAddedServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v.(string)
	}
	if v, ok := d.GetOk("renew_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok && v.(int) > 0 {
		request["RenewPeriod"] = convertKmsValueAddedServiceRenewPeriodRequest(v.(int))
	}
	if v, ok := d.GetOkExists("period"); ok && v.(int) > 0 {
		request["Period"] = convertKmsValueAddedServicePeriodRequest(v.(int))
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("value_added_service"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ValueAddedService",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "ProductVersion",
		"Value": "4",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList

	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	kmsServiceV2 := KmsServiceV2{client}

	// Ordering a value added service the account already holds is accepted and then refunded rather than
	// rejected, so nothing about the order itself reveals the duplicate. Recognising the existing instance
	// by its id prefix turns that into an immediate, actionable failure instead of a refunded order and a
	// wait that can only end in the same message.
	//
	// Only an instance that is actually in effect counts. The prefix matches every instance the service
	// ever minted, including ones that hold nothing: a refunded order from an earlier apply lingers in
	// Creating for minutes, and an expired one stays listed. Refusing on those would block a create that
	// the account is entitled to make, and a wrong refusal leaves the user no way forward - which is worse
	// than the refunded order this check exists to avoid, because the wait below still catches that.
	if d.Get("value_added_service").(string) == "2" {
		instances, err := kmsServiceV2.ListKmsValueAddedServiceInstancesByPrefix("dkr-")
		if err != nil {
			return WrapError(err)
		}
		now := time.Now()
		heldIds := make([]string, 0)
		for _, instance := range instances {
			if kmsValueAddedServiceInEffect(instance, now) {
				heldIds = append(heldIds, fmt.Sprint(instance["InstanceID"]))
			}
		}
		if len(heldIds) > 0 {
			return WrapError(Error("the account already holds this value added service in %s as instance %s. "+
				kmsValueAddedServiceAlreadyHeldMsg, client.RegionId, strings.Join(heldIds, ", ")))
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_value_added_service", action, AlibabaCloudSdkGoERROR)
	}

	// BssOpenApi reports business failures in the response body with a 200 status code, so err alone
	// does not tell whether the order succeeded.
	if fmt.Sprint(response["Success"]) != "true" {
		return WrapError(Error("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	if id == nil || fmt.Sprint(id) == "" {
		return WrapError(Error("%s did not return an instance id, response: %v", action, response))
	}

	// A duplicate order for a service the account already holds is accepted and then refunded, and the
	// refunded instance stays in Creating and later drops out of the available instance list, which would
	// make Read clear the id and the next plan order yet another refunded instance. Waiting for the
	// instance to actually become effective is what keeps that loop from starting, and it stays as the
	// backstop for whatever the prefix pre-flight above cannot see: a service ordered under a different
	// code, or a concurrent create in the same region.
	//
	// The wait runs on the ordered id rather than d.Id(), which is only set once the instance is known to
	// be effective, so a refunded order never reaches the state and leaves nothing to taint.
	orderedId := fmt.Sprint(id)
	refresh := kmsValueAddedServiceRefundedRefreshFunc(kmsServiceV2.KmsValueAddedServiceStateRefreshFunc(orderedId, "$.Status", []string{}))
	stateConf := BuildStateConf([]string{}, []string{kmsValueAddedServiceEffectiveStatus}, d.Timeout(schema.TimeoutCreate), 5*time.Second, refresh)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(Error("The ordered value added service %s did not become effective: %s. "+kmsValueAddedServiceAlreadyHeldMsg, orderedId, err))
	}

	d.SetId(orderedId)
	return resourceAliCloudKmsValueAddedServiceRead(d, meta)
}

func resourceAliCloudKmsValueAddedServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsValueAddedService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_value_added_service DescribeKmsValueAddedService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("region_id", objectRaw["Region"])
	// RenewalDuration is omitted unless the instance is automatically renewed, and formatInt panics on
	// the missing value, so keep whatever the state already holds.
	if v, ok := objectRaw["RenewalDuration"]; ok && v != nil {
		d.Set("renew_period", formatInt(convertKmsValueAddedServiceDataInstanceListRenewalDurationResponse(v)))
	}
	d.Set("renew_status", objectRaw["RenewStatus"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudKmsValueAddedServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if d.HasChange("renew_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renew_status")
	if d.HasChange("renew_period") {
		update = true
		request["RenewalPeriod"] = d.Get("renew_period")
	}

	request["RenewalPeriodUnit"] = "Y"
	request["SubscriptionType"] = "Subscription"
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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

	return resourceAliCloudKmsValueAddedServiceRead(d, meta)
}

func resourceAliCloudKmsValueAddedServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 60*time.Second, kmsServiceV2.KmsValueAddedServiceStateRefreshFunc(d.Id(), "$.InstanceID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// kmsValueAddedServiceRefundedRefreshFunc wraps a state refresh so that an instance which disappears
// after having been listed is reported as an error instead of being waited out. An automatically
// refunded order leaves exactly that trace: the instance is listed as Creating for a few minutes, then
// drops out once the refund settles, and it never reaches Normal. Failing on the disappearance keeps a
// create from sitting on the timeout for an outcome that is already decided, which matters because
// BssOpenApi gives no bound on how long a genuine order takes to become effective, so that timeout has
// to stay generous.
//
// The instance has to have been seen first. An id that is not listed yet right after the order is
// indistinguishable from one that never will be, and treating that as failure would break a create
// whose instance simply takes a moment to appear.
// kmsValueAddedServiceInEffect reports whether a listed instance holds the service right now. Status is
// not enough on its own: BssOpenApi keeps listing an instance after its subscription lapses, so an
// expired one can still read as Normal, and treating that as held would refuse a create the account is
// entitled to make. The subscription end has to be a real timestamp in the future as well.
//
// An EndTime that is absent or unparseable counts as not in effect. That direction is deliberate: it
// can cost at most a duplicate order, which the create wait catches and which BssOpenApi refunds by
// itself, whereas the opposite would block the create with nothing the user can do about it.
func kmsValueAddedServiceInEffect(instance map[string]interface{}, now time.Time) bool {
	if fmt.Sprint(instance["Status"]) != kmsValueAddedServiceEffectiveStatus {
		return false
	}
	endTime, ok := instance["EndTime"]
	if !ok || endTime == nil || fmt.Sprint(endTime) == "" {
		log.Printf("[WARN] alicloud_kms_value_added_service instance %v reports no EndTime, treating it as not in effect", instance["InstanceID"])
		return false
	}
	// EndTime arrives as RFC3339, for example 2027-08-04T16:00:00Z.
	end, err := time.Parse(time.RFC3339, fmt.Sprint(endTime))
	if err != nil {
		log.Printf("[WARN] alicloud_kms_value_added_service could not parse the EndTime %v of instance %v, treating it as not in effect: %s", endTime, instance["InstanceID"], err)
		return false
	}
	return end.After(now)
}

func kmsValueAddedServiceRefundedRefreshFunc(refresh resource.StateRefreshFunc) resource.StateRefreshFunc {
	listed := false
	return func() (interface{}, string, error) {
		object, state, err := refresh()
		if err != nil || state != "" {
			listed = listed || state != ""
			return object, state, err
		}
		if listed {
			return object, state, WrapError(Error("the instance is no longer available, which means its order was refunded"))
		}
		return object, state, nil
	}
}

func convertKmsValueAddedServiceDataInstanceListRenewalDurationResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "12":
		return "1"
	case "24":
		return "2"
	case "36":
		return "3"
	}
	return source
}
func convertKmsValueAddedServiceDataInstanceListRenewalDurationUnitResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "M":
		return "Y"
	}
	return source
}
func convertKmsValueAddedServiceRenewPeriodRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "12"
	case "2":
		return "24"
	case "3":
		return "36"
	}
	return source
}
func convertKmsValueAddedServicePeriodRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "12"
	case "2":
		return "24"
	case "3":
		return "36"
	}
	return source
}
