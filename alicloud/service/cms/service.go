// Package cms provides the Framework-native resources and data sources for the
// CMS (CloudMonitor) product, API version 2024-03-30.
package cms

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

const (
	cmsProduct  = "Cms"
	cmsVersion  = "2024-03-30"
	cmsTypeName = "alicloud_cms_transformer"
	cmsDSName   = "alicloud_cms_transformers"

	// pageSizeLarge is the default MaxResults used by the transformers list
	// (data source) paging loop. It mirrors the root provider PageSizeLarge
	// constant; the cms package cannot import the root package without a
	// cycle, so the value is declared locally.
	pageSizeLarge = 50
)

// stringPointer returns a pointer to s. It is the local equivalent of the
// root provider StringPointer helper; kept in-package so the cms package
// builds *string request/query fields without importing the root package.
func stringPointer(s string) *string {
	return &s
}

// describeCmsTransformer fetches a single Transformer by its composite id
// "<transformerId>:<workspace>" and returns the API "transformer" object, or
// (nil, nil) when the transformer no longer exists.
func describeCmsTransformer(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	transformerId, workspace, ok := splitTransformerId(id)
	if !ok {
		return nil, fmt.Errorf("invalid resource id %q: expected <transformerId>:<workspace>", id)
	}

	action := fmt.Sprintf("/transformers/%s", transformerId)
	query := map[string]*string{
		"regionId":  stringPointer(client.RegionId),
		"workspace": stringPointer(workspace),
	}

	response, err := client.RoaGet(cmsProduct, cmsVersion, action, query, nil, nil)
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
			return nil, nil
		}
		return nil, err
	}

	object, ok := response["transformer"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("the GetTransformer response carried no transformer object: %v", response)
	}
	return object, nil
}

// toggleCmsTransformer enables or disables a transformer via the
// /transformers/{id}/{enable|disable} sub-operation.
func toggleCmsTransformer(client *connectivity.AliyunClient, transformerId, workspace, op string) error {
	action := fmt.Sprintf("/transformers/%s/%s", transformerId, op)
	query := map[string]*string{
		"regionId":  stringPointer(client.RegionId),
		"workspace": stringPointer(workspace),
	}
	_, err := client.RoaPut(cmsProduct, cmsVersion, action, query, nil, nil, true)
	if err != nil {
		return errs.WrapErrorf(err, "calling %s %s", cmsProduct, action)
	}
	return nil
}

// splitTransformerId splits the composite id "<transformerId>:<workspace>".
func splitTransformerId(id string) (transformerId, workspace string, ok bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}
