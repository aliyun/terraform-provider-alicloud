package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

const (
	cmsProduct   = "Cms"
	cmsVersion   = "2024-03-30"
	contactPath  = "/contact"
	contactsPath = "/contacts"
	sourceOBS    = "OBS"
)

// CmsService holds the Alibaba Cloud client and exposes the CMS Contact API
// operations shared by the managed resource and the data source. Methods return
// raw API objects so the resource/data source layers own Terraform state
// mapping; this keeps a single API code path that the no-cloud unit tests
// exercise directly against a mocked ROA transport.
type CmsService struct {
	client *connectivity.AliyunClient
}

// NewService returns a CmsService bound to client.
func NewService(client *connectivity.AliyunClient) *CmsService {
	return &CmsService{client: client}
}

// Client exposes the bound client for the resource/data source Configure path.
func (s *CmsService) Client() *connectivity.AliyunClient { return s.client }

// retry runs fn until it succeeds, returns a non-retryable error, or the
// bounded attempt budget is exhausted. Retryable errors (throttling, eventual
// consistency) back off with the provider's incremental waiter; the context
// is honored between attempts. The ROA runtime already auto-retries 5xx when
// autoRetry=true; this loop covers business-level retryable codes.
func retry(ctx context.Context, fn func() (map[string]interface{}, error)) (map[string]interface{}, error) {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	maxAttempts := 5
	var response map[string]interface{}
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		response, err = fn()
		if err == nil {
			return response, nil
		}
		if !helper.NeedRetry(err) && !isExpectedErrors(err, []string{"ResourceExist", "InvalidParameter", "FrequencyLimit", "SmsQuotaExceeded"}) {
			return response, err
		}
		logs.AddDebug("cms contact retry", err)
		wait()
	}
	return response, errs.WrapErrorf(err, "retrying CMS Contact operation after %d attempts", maxAttempts)
}

// isExpectedErrors reports whether err carries one of codes, matching the
// provider's expected-error convention across wrapped SDK errors.
func isExpectedErrors(err error, codes []string) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, c := range codes {
		if msg == c || contains(msg, c) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// contactIdsQuery builds the ROA query map carrying the contactIds JSON array
// and the fixed OBS source used by the CMS ListContacts / DeleteContacts API.
func contactIdsQuery(ids []string) map[string]*string {
	raw, _ := json.Marshal(ids)
	idsStr := string(raw)
	source := sourceOBS
	return map[string]*string{
		"contactIds": &idsStr,
		"source":     &source,
	}
}

// CreateContact creates a CMS alarm contact and returns the raw response.
func (s *CmsService) CreateContact(ctx context.Context, body map[string]interface{}) (map[string]interface{}, error) {
	response, err := retry(ctx, func() (map[string]interface{}, error) {
		return s.client.RoaPost(cmsProduct, cmsVersion, contactPath, nil, nil, body, true)
	})
	logs.AddDebug("CreateContact", response, body)
	if err != nil {
		return nil, errs.WrapErrorf(err, "creating CMS Contact: POST %s", contactPath)
	}
	contactId := fmt.Sprint(response["contactId"])
	if contactId == "" || contactId == "<nil>" {
		return nil, errs.Error("CreateContact response missing contactId: %v", response)
	}
	return response, nil
}

// GetContact queries a single CMS alarm contact by id via ListContacts. It
// returns (nil, nil) when the contact no longer exists so the caller can drop
// state. The first element of the returned data array is type-guarded before
// the map conversion: a bare .(map[string]interface{}) assertion would panic
// if the API ever returned a non-map element, so the comma-ok form returns a
// descriptive not-found error instead.
func (s *CmsService) GetContact(ctx context.Context, id string) (map[string]interface{}, error) {
	response, err := retry(ctx, func() (map[string]interface{}, error) {
		return s.client.RoaGet(cmsProduct, cmsVersion, contactPath, contactIdsQuery([]string{id}), nil, nil)
	})
	logs.AddDebug("GetContact", response, contactIdsQuery([]string{id}))
	if err != nil {
		if isExpectedErrors(err, []string{"ResourceNotFound", "404"}) || errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "reading CMS Contact %s: GET %s", id, contactPath)
	}
	arr := extractContactsArray(response)
	if len(arr) == 0 {
		return nil, nil
	}
	contact, ok := arr[0].(map[string]interface{})
	if !ok {
		// AC2: a malformed (non-map) element must not panic; surface a clear
		// error so the resource Read reports failure rather than crashing.
		return nil, errs.WrapErrorf(errs.NotFoundErr("Contact", id), "CMS Contact %s returned a non-map element: %T", id, arr[0])
	}
	return contact, nil
}

// ListContacts paginates the CMS ListContacts API. It accepts the data-source
// filter values (empty strings mean "unset") and returns the raw contact
// objects plus the reported total.
func (s *CmsService) ListContacts(ctx context.Context, name, phone, email string, queryUngrouped bool) ([]map[string]interface{}, error) {
	pageNumber := 1
	pageSize := 100
	var objects []map[string]interface{}
	for {
		pageNumberStr := fmt.Sprintf("%d", pageNumber)
		pageSizeStr := fmt.Sprintf("%d", pageSize)
		source := sourceOBS
		query := map[string]*string{
			"pageNumber": &pageNumberStr,
			"pageSize":   &pageSizeStr,
			"source":     &source,
		}
		if name != "" {
			n := name
			query["name"] = &n
		}
		if phone != "" {
			p := phone
			query["phone"] = &p
		}
		if email != "" {
			e := email
			query["email"] = &e
		}
		if queryUngrouped {
			quc := "true"
			query["queryUngroupedContacts"] = &quc
		}
		response, err := retry(ctx, func() (map[string]interface{}, error) {
			return s.client.RoaGet(cmsProduct, cmsVersion, contactPath, query, nil, nil)
		})
		logs.AddDebug("ListContacts", response, query)
		if err != nil {
			return nil, errs.WrapErrorf(err, "listing CMS Contacts: GET %s", contactPath)
		}
		arr := extractContactsArray(response)
		for _, raw := range arr {
			if contact, ok := raw.(map[string]interface{}); ok {
				objects = append(objects, contact)
			}
		}
		total := contactsTotalToInt(response["total"], response)
		if len(objects) >= total || len(arr) == 0 {
			return objects, nil
		}
		pageNumber++
	}
}

// UpdateContact patches a CMS alarm contact by id.
func (s *CmsService) UpdateContact(ctx context.Context, id string, body map[string]interface{}) (map[string]interface{}, error) {
	response, err := retry(ctx, func() (map[string]interface{}, error) {
		return s.client.RoaPatch(cmsProduct, cmsVersion, fmt.Sprintf("/contact/%s", id), nil, nil, body, true)
	})
	logs.AddDebug("UpdateContact", response, body)
	if err != nil {
		return nil, errs.WrapErrorf(err, "updating CMS Contact %s: PATCH /contact/%s", id, id)
	}
	return response, nil
}

// DeleteContact deletes a CMS alarm contact by id. A not-found error is
// treated as success so Terraform destroy is idempotent.
func (s *CmsService) DeleteContact(ctx context.Context, id string) error {
	response, err := retry(ctx, func() (map[string]interface{}, error) {
		return s.client.RoaDelete(cmsProduct, cmsVersion, contactsPath, contactIdsQuery([]string{id}), nil, nil, true)
	})
	logs.AddDebug("DeleteContact", response, contactIdsQuery([]string{id}))
	if err != nil {
		if isExpectedErrors(err, []string{"ResourceNotFound", "404"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, "deleting CMS Contact %s: DELETE %s", id, contactsPath)
	}
	return nil
}

// extractContactsArray finds the contact list inside a ListContacts response.
// The CMS RESTful API wraps the array under data.data or data.contacts, or at
// the top-level contacts, depending on the endpoint variant; each is probed
// with comma-ok so an unexpected shape yields an empty slice, not a panic.
func extractContactsArray(response map[string]interface{}) []interface{} {
	if dataRaw, ok := response["data"].(map[string]interface{}); ok {
		if arr, ok := dataRaw["data"].([]interface{}); ok {
			return arr
		}
		if arr, ok := dataRaw["contacts"].([]interface{}); ok {
			return arr
		}
	}
	if arr, ok := response["contacts"].([]interface{}); ok {
		return arr
	}
	return nil
}

// contactsTotalToInt converts the ListContacts total field to int safely. The
// provider JSON decoder enables UseNumber, so total arrives as json.Number
// (not float64); a bare t.(float64) assertion would panic. Both forms are
// handled.
func contactsTotalToInt(total interface{}, response map[string]interface{}) int {
	if total == nil {
		if dataRaw, ok := response["data"].(map[string]interface{}); ok {
			total = dataRaw["total"]
		}
	}
	if n, ok := total.(json.Number); ok {
		if v, err := n.Int64(); err == nil {
			return int(v)
		}
	}
	if f, ok := total.(float64); ok {
		return int(f)
	}
	return 0
}
