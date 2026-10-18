package dataworks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// DataWorks RPC product/version constants. The connectivity client lower-cases
// and kebab→snake converts the product code, so "dataworks-public" is the key
// under which the endpoint is stored.
const (
	dataWorksProduct = "dataworks-public"
	dataWorksVersion = "2024-05-18"

	// errorCodeProjectNotUpgraded is returned when the workspace has not been
	// upgraded to the new DataWorks data-development model. The CreateResource/
	// ListResources calls fail with it; it is treated as a not-found-class
	// terminal state for Read/Delete so the resource drops out of state.
	errorCodeProjectNotUpgraded = "5801480744983330847"
	// errorCodeResourceGone is returned by GetResource/DeleteResource when the
	// resource no longer exists; treated as not-found for idempotent delete.
	errorCodeResourceGone = "5801488441564725318"

	notFoundResource = "Resource"

	// pageSize is the ListResources page size. Package-level so the unit tests
	// can drive pagination without duplicating the magic number.
	pageSize = 100
)

// getResource calls GetResource and returns the Resource object, or nil when the
// resource no longer exists (not-found or the project-not-upgraded terminal
// state). The composite id "<project_id>:<resource_id>" is split here so callers
// never need to reparse it.
func getResource(client *connectivity.AliyunClient, projectID, resourceID string) (map[string]interface{}, error) {
	response, err := client.RpcGet(dataWorksProduct, dataWorksVersion, "GetResource", map[string]interface{}{
		"Id":        resourceID,
		"ProjectId": projectID,
		"RegionId":  client.RegionId,
	}, nil)
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{errorCodeResourceGone, errorCodeProjectNotUpgraded}) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "Reading DataWorks Resource %s:%s: calling GetResource", projectID, resourceID)
	}
	object, _ := response["Resource"].(map[string]interface{})
	if object == nil {
		return nil, fmt.Errorf("the GetResource response carried no Resource object: %v", response)
	}
	return object, nil
}

// listResources pages ListResources and returns the matching Resource items. The
// owner/type/ids filters narrow the result set; ids is a set for O(1) membership.
func listResources(client *connectivity.AliyunClient, projectID, owner, typ string, idSet map[string]bool) ([]map[string]interface{}, error) {
	pageNumber := 1
	var objects []map[string]interface{}
	for {
		query := map[string]interface{}{
			"ProjectId":  projectID,
			"RegionId":   client.RegionId,
			"PageSize":   pageSize,
			"PageNumber": pageNumber,
		}
		if owner != "" {
			query["Owner"] = owner
		}
		if typ != "" {
			query["Type"] = typ
		}
		response, err := client.RpcGet(dataWorksProduct, dataWorksVersion, "ListResources", query, nil)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errorCodeProjectNotUpgraded}) {
				return nil, nil
			}
			return nil, errs.WrapErrorf(err, "Listing DataWorks Resources in project %s: calling ListResources", projectID)
		}
		items := resourcesFromResponse(response)
		for _, item := range items {
			if len(idSet) > 0 {
				if id, ok := item["Id"]; !ok || !idSet[fmt.Sprint(id)] {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(items) < pageSize {
			break
		}
		pageNumber++
	}
	return objects, nil
}

// resourcesFromResponse pulls the PagingInfo.Resources array out of a
// ListResources response, tolerating both json.Number and plain number shapes.
func resourcesFromResponse(response map[string]interface{}) []map[string]interface{} {
	paging, _ := response["PagingInfo"].(map[string]interface{})
	if paging == nil {
		return nil
	}
	raw, _ := paging["Resources"].([]interface{})
	out := make([]map[string]interface{}, 0, len(raw))
	for _, v := range raw {
		if m, ok := v.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// asString coerces an API value to a string, handling json.Number.
func asString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return val.String()
	case fmt.Stringer:
		return val.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(val)
	}
}

// asInt64 coerces an API numeric value to int64, handling json.Number and float64.
func asInt64(v interface{}) int64 {
	switch val := v.(type) {
	case nil:
		return 0
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			f, ferr := val.Float64()
			if ferr != nil {
				return 0
			}
			return int64(f)
		}
		return n
	case string:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

// splitResourceID splits the composite id "<project_id>:<resource_id>".
func splitResourceID(id string) (projectID, resourceID string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid DataWorks Resource id %q: expected <project_id>:<resource_id>", id)
	}
	return parts[0], parts[1], nil
}

// createDataWorksResource calls CreateResource and returns the new resource id.
// Spec and ResourceFile are optional; an empty spec is omitted so the server
// applies its default.
func createDataWorksResource(client *connectivity.AliyunClient, projectID, spec, resourceFile string) (string, error) {
	body := map[string]interface{}{
		"RegionId":  client.RegionId,
		"ProjectId": projectID,
	}
	if spec != "" {
		body["Spec"] = spec
	}
	if resourceFile != "" {
		body["ResourceFile"] = resourceFile
	}
	response, err := client.RpcPost(dataWorksProduct, dataWorksVersion, "CreateResource", nil, body, true)
	if err != nil {
		return "", errs.WrapErrorf(err, "Creating DataWorks Resource in project %s: calling CreateResource", projectID)
	}
	id := asString(response["Id"])
	if id == "" {
		return "", fmt.Errorf("the CreateResource response carried no Id: %v", response)
	}
	return id, nil
}

// updateDataWorksResource calls UpdateResource for the spec and/or resource_file
// fields whose send flag is true. Callers decide the send flags per the F4 rule.
func updateDataWorksResource(client *connectivity.AliyunClient, projectID, resourceID, spec string, sendSpec bool, resourceFile string, sendFile bool) error {
	body := map[string]interface{}{
		"Id":        resourceID,
		"ProjectId": projectID,
		"RegionId":  client.RegionId,
	}
	if sendSpec {
		body["Spec"] = spec
	}
	if sendFile {
		body["ResourceFile"] = resourceFile
	}
	if _, err := client.RpcPost(dataWorksProduct, dataWorksVersion, "UpdateResource", nil, body, true); err != nil {
		return errs.WrapErrorf(err, "Updating DataWorks Resource %s:%s: calling UpdateResource", projectID, resourceID)
	}
	return nil
}

// renameDataWorksResource calls RenameResource to set the resource name.
func renameDataWorksResource(client *connectivity.AliyunClient, projectID, resourceID, name string) error {
	body := map[string]interface{}{
		"Id":        resourceID,
		"Name":      name,
		"ProjectId": projectID,
		"RegionId":  client.RegionId,
	}
	if _, err := client.RpcPost(dataWorksProduct, dataWorksVersion, "RenameResource", nil, body, true); err != nil {
		return errs.WrapErrorf(err, "Renaming DataWorks Resource %s:%s: calling RenameResource", projectID, resourceID)
	}
	return nil
}

// moveDataWorksResource calls MoveResource to set the destination path.
func moveDataWorksResource(client *connectivity.AliyunClient, projectID, resourceID, p string) error {
	body := map[string]interface{}{
		"Id":        resourceID,
		"ProjectId": projectID,
		"RegionId":  client.RegionId,
		"Path":      p,
	}
	if _, err := client.RpcPost(dataWorksProduct, dataWorksVersion, "MoveResource", nil, body, true); err != nil {
		return errs.WrapErrorf(err, "Moving DataWorks Resource %s:%s: calling MoveResource", projectID, resourceID)
	}
	return nil
}

// deleteDataWorksResource calls DeleteResource. A not-found response (code
// 5801488441564725318 or a 404) is success so destroy is idempotent.
func deleteDataWorksResource(client *connectivity.AliyunClient, projectID, resourceID string) error {
	body := map[string]interface{}{
		"Id":        resourceID,
		"RegionId":  client.RegionId,
		"ProjectId": projectID,
	}
	_, err := client.RpcPost(dataWorksProduct, dataWorksVersion, "DeleteResource", nil, body, true)
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{errorCodeResourceGone}) {
			return nil
		}
		return errs.WrapErrorf(err, "Deleting DataWorks Resource %s:%s: calling DeleteResource", projectID, resourceID)
	}
	return nil
}

// notFoundHTTPStatus reports whether an HTTP status looks like not-found, used by
// the unit-test mock handlers to emulate the server's 4xx behavior.
func notFoundHTTPStatus(code int) bool {
	return code == http.StatusNotFound
}

var _ = notFoundHTTPStatus // retained for mock-handler parity
