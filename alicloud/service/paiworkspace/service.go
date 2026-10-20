package paiworkspace

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// PAI Workspace (POP product AIWorkSpace, version 2021-02-04) wire constants.
// CreateProductOrders is a ROA POST that places a paid order; ListProducts is a
// ROA GET that reads the product catalog without purchasing.
const (
	apiProduct = "AIWorkSpace"
	apiVersion = "2021-02-04"
	pathCreate = "/api/v1/productorders"
	pathList   = "/api/v1/products"
)

// instancePropertyInput is the plain-Go form of one InstanceProperties entry,
// used to build the CreateProductOrders body independent of framework types.
type instancePropertyInput struct {
	Code  string
	Name  string
	Value string
}

// productCreateInput is the plain-Go form of the create-only order parameters.
// It is populated from the framework plan in Create and fed to buildCreateBody,
// which keeps the request-body construction unit-testable without the cloud.
type productCreateInput struct {
	ProductCode        string
	ChargeType         string
	PricingCycle       string
	Duration           int64
	AutoRenew          bool
	OrderType          string
	AutoPay            bool
	InstanceProperties []instancePropertyInput
}

// buildCreateBody constructs the CreateProductOrders request body from the
// create-only parameters. ProductCode is always sent (required, constant
// PAI_share). AutoRenew/AutoPay are booleans and always sent; the remaining
// order fields are sent only when the caller set a non-zero value, so an empty
// plan does not push invalid empty strings to the order API. AutoPay lives at
// the body root (body.AutoPay), not inside Products, matching the CloudSpec
// attribute mapping $.AutoPay -> body.AutoPay.
func buildCreateBody(in productCreateInput) map[string]interface{} {
	product := map[string]interface{}{
		"ProductCode": in.ProductCode,
		"AutoRenew":   in.AutoRenew,
	}
	if in.ChargeType != "" {
		product["ChargeType"] = in.ChargeType
	}
	if in.PricingCycle != "" {
		product["PricingCycle"] = in.PricingCycle
	}
	if in.OrderType != "" {
		product["OrderType"] = in.OrderType
	}
	if in.Duration != 0 {
		product["Duration"] = in.Duration
	}
	if len(in.InstanceProperties) > 0 {
		items := make([]interface{}, 0, len(in.InstanceProperties))
		for _, ip := range in.InstanceProperties {
			// Drop entries that carry no property at all: an all-empty
			// InstanceProperties item is meaningless and would push an empty
			// object to the order API.
			if ip.Code == "" && ip.Name == "" && ip.Value == "" {
				continue
			}
			item := map[string]interface{}{}
			if ip.Code != "" {
				item["Code"] = ip.Code
			}
			if ip.Name != "" {
				item["Name"] = ip.Name
			}
			if ip.Value != "" {
				item["Value"] = ip.Value
			}
			items = append(items, item)
		}
		product["InstanceProperties"] = items
	}
	return map[string]interface{}{
		"Products": []interface{}{product},
		"AutoPay":  in.AutoPay,
	}
}

// createProduct places a product order via CreateProductOrders and returns the
// ProductId from the response (ProductIds[0]). The order is a real, chargeable
// purchase; tests must not drive it against the real cloud.
func createProduct(client *connectivity.AliyunClient, body map[string]interface{}) (string, error) {
	resp, err := client.RoaPost(apiProduct, apiVersion, pathCreate, nil, nil, body, true)
	if err != nil {
		return "", err
	}
	ids, _ := resp["ProductIds"].([]interface{})
	if len(ids) == 0 {
		return "", fmt.Errorf("CreateProductOrders returned no ProductIds in response: %v", resp)
	}
	id, _ := ids[0].(string)
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("CreateProductOrders returned an empty ProductId in response: %v", resp)
	}
	return id, nil
}

// readProduct returns the catalog product whose ProductId matches, or nil when
// no such product exists. GetProduct is not published as a distinct action in
// the public OpenAPI catalog for this version, so Read is implemented through
// ListProducts (GET /api/v1/products) with a ProductCodes filter, bounded by the
// fixed catalog value of ProductCode, followed by a client-side match on
// ProductId. When productCode is empty (e.g. on import with no prior state),
// the whole catalog is listed and filtered; PAI Workspace products are few, so
// the scale risk is contained.
func readProduct(client *connectivity.AliyunClient, productId, productCode string) (map[string]interface{}, error) {
	query := map[string]*string{}
	if productCode != "" {
		v := productCode
		query["ProductCodes"] = &v
	}
	verbose := "true"
	query["Verbose"] = &verbose

	resp, err := client.RoaGet(apiProduct, apiVersion, pathList, query, nil, nil)
	if err != nil {
		return nil, err
	}
	products, _ := resp["Products"].([]interface{})
	return findProduct(products, productId), nil
}

// findProduct returns the first list entry whose ProductId equals productId, or
// nil. Extracted so the filter is unit-testable without a client.
func findProduct(products []interface{}, productId string) map[string]interface{} {
	for _, p := range products {
		obj, _ := p.(map[string]interface{})
		if obj == nil {
			continue
		}
		if pid, _ := obj["ProductId"].(string); pid == productId {
			return obj
		}
	}
	return nil
}
