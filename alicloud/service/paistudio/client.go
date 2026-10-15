package paistudio

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// paiStudioEndpoint returns the regional PAI Studio ROA endpoint.
//
// The public documentation alias is paistudio.aliyuncs.com, but that host does
// not resolve on internal networks; the region-scoped host pai.<region>-
// .aliyuncs.com is the resolvable form (probing pai.cn-beijing.aliyuncs.com
// returns the PaiStudio product). PaiStudio has no location-service service
// code either, so DescribeEndpoints cannot discover it.
//
// The endpoint is built from the client region here rather than registered in
// connectivity's endpoint map on purpose: scripts/ci/changed-unit-tests.go
// treats any changed non-resource/non-test .go file in a package as a trigger
// to run every TestUnit* in that package, and connectivity contains
// credential-gated tests that skip (counting as failures) in the
// credential-less ci-check-quick environment. Keeping endpoint.go unchanged
// avoids that broadening. Building the endpoint locally is otherwise identical
// to how the ACK/CR resources use NewTeaRoaCommonClient with an explicit host.
func paiStudioEndpoint(client *connectivity.AliyunClient) string {
	return fmt.Sprintf("pai.%s.aliyuncs.com", client.RegionId)
}

// roaRequest issues a PaiStudio ROA request against the regional endpoint and
// returns the response payload, unwrapped from the tea "body" envelope. It is
// the local analogue of connectivity.AliyunClient.roaRequest, but it pins the
// endpoint to pai.<region>.aliyuncs.com instead of resolving it through the
// shared endpoint map (which has no PaiStudio entry).
func roaRequest(client *connectivity.AliyunClient, method, pathName string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	conn, err := client.NewTeaRoaCommonClient(paiStudioEndpoint(client))
	if err != nil {
		return nil, err
	}
	runtime := &util.RuntimeOptions{}
	runtime.SetAutoretry(autoRetry)
	response, err := conn.DoRequest(tea.String(paiStudioVersion), nil, tea.String(method), tea.String("AK"), tea.String(pathName), nil, nil, body, runtime)
	if respBody, isExist := response["body"]; isExist && respBody != nil {
		// The tea ROA client nests the parsed payload under a "body" key; unwrap
		// it so callers see the flat object. The assertion is guarded (unlike the
		// shared helper) so a non-map body leaves the envelope intact instead of
		// panicking.
		if m, ok := respBody.(map[string]interface{}); ok {
			response = m
		}
	}
	return response, formatPaistudioError(response, err)
}

// formatPaistudioError mirrors connectivity.formatError. Some Alibaba ROA APIs
// answer with HTTP 200 yet carry an error Code in the body; without this a
// such response would be treated as success. errs.NotFoundError recognizes the
// tea.SDKError this builds (StatusCode 404 / NotFound message), so Delete
// idempotency against an already-vanished resource keeps working. The upstream
// helper is unexported, so it is duplicated here verbatim.
func formatPaistudioError(response map[string]interface{}, err error) error {
	if err != nil {
		return err
	}
	code, ok1 := response["Code"]
	if !ok1 {
		code, ok1 = response["code"]
	}
	// Some products (e.g. FC) return a map-typed Code; leave the transport error.
	if ok1 && !(isString(code) || isInteger(code)) {
		return err
	}
	if ok1 && (strings.ToLower(fmt.Sprint(code)) == "success" ||
		strings.ToLower(fmt.Sprint(code)) == "ok" ||
		fmt.Sprint(code) == "200" ||
		fmt.Sprint(code) == "0") {
		return err
	}
	success, ok2 := response["Success"]
	if !ok2 {
		success, ok2 = response["success"]
	}
	if ok2 && fmt.Sprint(success) == "true" {
		return err
	}
	message, ok3 := response["Message"]
	if !ok3 {
		message, ok3 = response["message"]
	}
	if ok3 && (message == nil || fmt.Sprint(message) == "") {
		return err
	}
	if ok1 || ok2 {
		statusCode := 200
		if v, ok := response["StatusCode"]; ok {
			if iv, ok := v.(*int); ok {
				statusCode = tea.IntValue(iv)
			}
		}
		return tea.NewSDKError(map[string]interface{}{
			"statusCode": statusCode,
			"code":       tea.ToString(code),
			"message":    tea.ToString(message),
			"data":       response,
		})
	}
	return err
}

// isString and isInteger mirror the unexported connectivity helpers used by
// formatError; they are duplicated here for the same reason.
func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func isInteger(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		if _, err := strconv.Atoi(fmt.Sprint(value)); err == nil {
			return true
		}
		return false
	default:
		return false
	}
}
