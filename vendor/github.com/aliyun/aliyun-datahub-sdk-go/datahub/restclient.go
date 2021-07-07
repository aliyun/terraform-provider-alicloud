package datahub

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
    "errors"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "os"
    "sort"
    "strconv"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"
)

const (
    httpHeaderAcceptEncoding     = "Accept-Encoding"
    httpHeaderAuthorization      = "Authorization"
    httpHeadercacheControl       = "Cache-Control"
    httpHeaderChunked            = "chunked"
    httpHeaderClientVersion      = "x-datahub-client-version"
    httpHeaderContentDisposition = "Content-Disposition"
    httpHeaderContentEncoding    = "Content-Encoding"
    httpHeaderContentLength      = "Content-Length"
    httpHeaderContentMD5         = "Content-MD5"
    httpHeaderContentType        = "Content-Type"
    httpHeaderDate               = "Date"
    httpHeaderETAG               = "ETag"
    httpHeaderEXPIRES            = "Expires"
    httpHeaderHost               = "Host"
    httpHeaderlastModified       = "Last-Modified"
    httpHeaderLocation           = "Location"
    httpHeaderRange              = "Range"
    httpHeaderRawSize            = "x-datahub-content-raw-size"
    httpHeaderRequestAction      = "x-datahub-request-action"
    httpHeaderRequestId          = "x-datahub-request-id"
    httpHeaderSecurityToken      = "x-datahub-security-token"
    httpHeaderTransferEncoding   = "Transfer-Encoding"
    httpHeaderUserAgent          = "User-Agent"
)

const (
    datahubHeadersPrefix = "x-datahub-"
)

func init() {
    // Log as JSON instead of the default ASCII formatter.
    log.SetFormatter(&log.TextFormatter{})

    // Output to stdout instead of the default stderr
    // Can be any io.Writer, see below for File examples
    log.SetOutput(os.Stdout)

    // Only log the level severity or above.
    dev := strings.ToLower(os.Getenv("GODATAHUB_DEV"))
    switch dev {
    case "true":
        log.SetLevel(log.DebugLevel)
    default:
        log.SetLevel(log.WarnLevel)
    }
}

// DialContextFn was defined to make code more readable.
type DialContextFn func(ctx context.Context, network, address string) (net.Conn, error)

// TraceDialContext implements our own dialer in order to trace conn info.
func TraceDialContext(ctimeout time.Duration) DialContextFn {
    dialer := &net.Dialer{
        Timeout:   ctimeout,
        KeepAlive: ctimeout,
    }
    return func(ctx context.Context, network, addr string) (net.Conn, error) {
        conn, err := dialer.DialContext(ctx, network, addr)
        if err != nil {
            return nil, err
        }

        log.Debug("connect done, use", conn.LocalAddr().String())
        return conn, nil
    }
}

// RestClient rest客户端
type RestClient struct {
    // Endpoint datahub服务的endpint
    Endpoint string
    // Useragent user agent
    Useragent string
    // HttpClient http client
    HttpClient *http.Client
    // Account
    Account        Account
    CompressorType CompressorType
}

// NewRestClient create a new rest client
func NewRestClient(endpoint string, useragent string, httpclient *http.Client, account Account, ctype CompressorType) *RestClient {
    if strings.HasSuffix(endpoint, "/") {
        endpoint = endpoint[0 : len(endpoint)-1]
    }
    return &RestClient{
        Endpoint:       endpoint,
        Useragent:      useragent,
        HttpClient:     httpclient,
        Account:        account,
        CompressorType: ctype,
    }
}

// Get send HTTP Get method request
func (client *RestClient) Get(resource string) ([]byte, error) {
    return client.request(http.MethodGet, resource, &EmptyRequest{})
}

// Post send HTTP Post method request
func (client *RestClient) Post(resource string, model RequestModel) ([]byte, error) {
    return client.request(http.MethodPost, resource, model)
}

// Put send HTTP Put method request
func (client *RestClient) Put(resource string, model RequestModel) (interface{}, error) {
    return client.request(http.MethodPut, resource, model)
}

// Delete send HTTP Delete method request
func (client *RestClient) Delete(resource string) (interface{}, error) {
    return client.request(http.MethodDelete, resource, &EmptyRequest{})
}

func (client *RestClient) request(method, resource string, requestModel RequestModel) ([]byte, error) {
    url := fmt.Sprintf("%s%s", client.Endpoint, resource)

    header := map[string]string{
        httpHeaderClientVersion: DATAHUB_CLIENT_VERSION,
        httpHeaderDate:          time.Now().UTC().Format(http.TimeFormat),
        httpHeaderUserAgent:     client.Useragent,

        //TODO ContentType 和 RequestAction 应该以参数传进来
        //httpHeaderContentType:   "application/x-protobuf",
        //httpHeaderRequestAction: "pub",
    }

    //serialization
    reqBody, err := requestModel.requestBodyEncode(header)
    if err != nil {
        return nil, err
    }

    //compress
    client.compress(header, &reqBody)

    if client.Account.GetSecurityToken() != "" {
        header[httpHeaderSecurityToken] = client.Account.GetSecurityToken()
    }
    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }

    for k, v := range header {
        req.Header.Add(k, v)
    }
    client.buildSignature(&req.Header, method, resource)

    resp, err := client.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    //decompress
    if err := client.decompress(&respBody, &resp.Header); err != nil {
        return nil, err
    }

    //detect error
    respResult, err := newCommonResponseResult(resp.StatusCode, &resp.Header, respBody)
    log.Debug(fmt.Sprintf("request id: %s\nrequest url: %s\nrequest headers: %v\nrequest body: %s\nresponse headers: %v\nresponse body: %s",
        respResult.RequestId, url, req.Header, string(reqBody), resp.Header, string(respBody)))
    if err != nil {
        return nil, err
    }

    return respBody, nil
}

func (client *RestClient) buildSignature(header *http.Header, method, resource string) {
    builder := make([]string, 0, 5)
    builder = append(builder, method)
    builder = append(builder, header.Get(httpHeaderContentType))
    builder = append(builder, header.Get(httpHeaderDate))

    headersToSign := make(map[string][]string)
    for k, v := range *header {
        lower := strings.ToLower(k)
        if strings.HasPrefix(lower, datahubHeadersPrefix) {
            headersToSign[lower] = v
        }
    }

    keys := make([]string, len(headersToSign))
    for k := range headersToSign {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        for _, v := range headersToSign[k] {
            builder = append(builder, fmt.Sprintf("%s:%s", k, v))
        }
    }

    builder = append(builder, resource)

    canonString := strings.Join(builder, "\n")

    log.Debug(fmt.Sprintf("canonString: %s, accesskey: %s", canonString, client.Account.GetAccountKey()))

    hash := hmac.New(sha1.New, []byte(client.Account.GetAccountKey()))
    hash.Write([]byte(canonString))
    crypto := hash.Sum(nil)
    signature := base64.StdEncoding.EncodeToString(crypto)
    authorization := fmt.Sprintf("DATAHUB %s:%s", client.Account.GetAccountId(), signature)

    header.Add(httpHeaderAuthorization, authorization)
}

func (client *RestClient) compress(header map[string]string, reqBody *[]byte) {
    if client.CompressorType == NOCOMPRESS {
        return
    }
    compressor := getCompressor(client.CompressorType)
    if compressor != nil {
        compressedReqBody, err := compressor.Compress(*reqBody)
        header[httpHeaderAcceptEncoding] = client.CompressorType.String()
        //compress is valid
        if err == nil && len(compressedReqBody) < len(*reqBody) {
            header[httpHeaderContentEncoding] = client.CompressorType.String()
            //header[httpHeaderAcceptEncoding] = client.CompressorType.String()
            header[httpHeaderRawSize] = strconv.Itoa(len(*reqBody))
            *reqBody = compressedReqBody
        } else {
            //print warning and give up compress when compress failed
            log.Warning("compress failed or compress invalid, give up compression")
        }
    }
    header[httpHeaderContentLength] = strconv.Itoa(len(*reqBody))
    return
}

func (client *RestClient) decompress(respBody *[]byte, header *http.Header) error {
    encoding := header.Get(httpHeaderContentEncoding)
    if encoding == "" {
        return nil
    }
    compressor := getCompressor(CompressorType(encoding))
    if compressor == nil {
        return errors.New(fmt.Sprintf("not support the compress mode %s ", encoding))
    }
    rawSize := header.Get(httpHeaderRawSize)
    //str convert to int64
    size, err := strconv.ParseInt(rawSize, 10, 64)
    if err != nil {
        return err
    }

    buf, err := compressor.DeCompress(*respBody, size)
    if err != nil {
        return err
    }
    *respBody = buf
    return nil
}
