package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash"
	"io"
	"mime"
	"net/url"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
)

var characters = map[string]bool{"-": true, "_": true, ".": true, "~": true}

// signKeyList is a list about params without value
var signKeyList = []string{"location", "cors", "objectMeta",
	"uploadId", "partNumber", "security-token",
	"position", "img", "style", "styleName",
	"replication", "replicationProgress",
	"replicationLocation", "cname", "qos",
	"startTime", "endTime", "symlink",
	"x-oss-process", "response-content-type",
	"response-content-language", "response-expires",
	"response-cache-control", "response-content-disposition",
	"response-content-encoding", "udf", "udfName", "udfImage",
	"udfId", "udfImageDesc", "udfApplication",
	"udfApplicationLog", "restore", "callback", "callback-var",
	"policy", "encryption", "versions", "versioning", "versionId"}

var extToMimeType = map[string]string{
	".xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xltx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
	".potx":    "application/vnd.openxmlformats-officedocument.presentationml.template",
	".ppsx":    "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	".pptx":    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".sldx":    "application/vnd.openxmlformats-officedocument.presentationml.slide",
	".docx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".dotx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
	".xlam":    "application/vnd.ms-excel.addin.macroEnabled.12",
	".xlsb":    "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
	".apk":     "application/vnd.android.package-archive",
	".hqx":     "application/mac-binhex40",
	".cpt":     "application/mac-compactpro",
	".doc":     "application/msword",
	".ogg":     "application/ogg",
	".pdf":     "application/pdf",
	".rtf":     "text/rtf",
	".mif":     "application/vnd.mif",
	".xls":     "application/vnd.ms-excel",
	".ppt":     "application/vnd.ms-powerpoint",
	".odc":     "application/vnd.oasis.opendocument.chart",
	".odb":     "application/vnd.oasis.opendocument.database",
	".odf":     "application/vnd.oasis.opendocument.formula",
	".odg":     "application/vnd.oasis.opendocument.graphics",
	".otg":     "application/vnd.oasis.opendocument.graphics-template",
	".odi":     "application/vnd.oasis.opendocument.image",
	".odp":     "application/vnd.oasis.opendocument.presentation",
	".otp":     "application/vnd.oasis.opendocument.presentation-template",
	".ods":     "application/vnd.oasis.opendocument.spreadsheet",
	".ots":     "application/vnd.oasis.opendocument.spreadsheet-template",
	".odt":     "application/vnd.oasis.opendocument.text",
	".odm":     "application/vnd.oasis.opendocument.text-master",
	".ott":     "application/vnd.oasis.opendocument.text-template",
	".oth":     "application/vnd.oasis.opendocument.text-web",
	".sxw":     "application/vnd.sun.xml.writer",
	".stw":     "application/vnd.sun.xml.writer.template",
	".sxc":     "application/vnd.sun.xml.calc",
	".stc":     "application/vnd.sun.xml.calc.template",
	".sxd":     "application/vnd.sun.xml.draw",
	".std":     "application/vnd.sun.xml.draw.template",
	".sxi":     "application/vnd.sun.xml.impress",
	".sti":     "application/vnd.sun.xml.impress.template",
	".sxg":     "application/vnd.sun.xml.writer.global",
	".sxm":     "application/vnd.sun.xml.math",
	".sis":     "application/vnd.symbian.install",
	".wbxml":   "application/vnd.wap.wbxml",
	".wmlc":    "application/vnd.wap.wmlc",
	".wmlsc":   "application/vnd.wap.wmlscriptc",
	".bcpio":   "application/x-bcpio",
	".torrent": "application/x-bittorrent",
	".bz2":     "application/x-bzip2",
	".vcd":     "application/x-cdlink",
	".pgn":     "application/x-chess-pgn",
	".cpio":    "application/x-cpio",
	".csh":     "application/x-csh",
	".dvi":     "application/x-dvi",
	".spl":     "application/x-futuresplash",
	".gtar":    "application/x-gtar",
	".hdf":     "application/x-hdf",
	".jar":     "application/x-java-archive",
	".jnlp":    "application/x-java-jnlp-file",
	".js":      "application/x-javascript",
	".ksp":     "application/x-kspread",
	".chrt":    "application/x-kchart",
	".kil":     "application/x-killustrator",
	".latex":   "application/x-latex",
	".rpm":     "application/x-rpm",
	".sh":      "application/x-sh",
	".shar":    "application/x-shar",
	".swf":     "application/x-shockwave-flash",
	".sit":     "application/x-stuffit",
	".sv4cpio": "application/x-sv4cpio",
	".sv4crc":  "application/x-sv4crc",
	".tar":     "application/x-tar",
	".tcl":     "application/x-tcl",
	".tex":     "application/x-tex",
	".man":     "application/x-troff-man",
	".me":      "application/x-troff-me",
	".ms":      "application/x-troff-ms",
	".ustar":   "application/x-ustar",
	".src":     "application/x-wais-source",
	".zip":     "application/zip",
	".m3u":     "audio/x-mpegurl",
	".ra":      "audio/x-pn-realaudio",
	".wav":     "audio/x-wav",
	".wma":     "audio/x-ms-wma",
	".wax":     "audio/x-ms-wax",
	".pdb":     "chemical/x-pdb",
	".xyz":     "chemical/x-xyz",
	".bmp":     "image/bmp",
	".gif":     "image/gif",
	".ief":     "image/ief",
	".png":     "image/png",
	".wbmp":    "image/vnd.wap.wbmp",
	".ras":     "image/x-cmu-raster",
	".pnm":     "image/x-portable-anymap",
	".pbm":     "image/x-portable-bitmap",
	".pgm":     "image/x-portable-graymap",
	".ppm":     "image/x-portable-pixmap",
	".rgb":     "image/x-rgb",
	".xbm":     "image/x-xbitmap",
	".xpm":     "image/x-xpixmap",
	".xwd":     "image/x-xwindowdump",
	".css":     "text/css",
	".rtx":     "text/richtext",
	".tsv":     "text/tab-separated-values",
	".jad":     "text/vnd.sun.j2me.app-descriptor",
	".wml":     "text/vnd.wap.wml",
	".wmls":    "text/vnd.wap.wmlscript",
	".etx":     "text/x-setext",
	".mxu":     "video/vnd.mpegurl",
	".flv":     "video/x-flv",
	".wm":      "video/x-ms-wm",
	".wmv":     "video/x-ms-wmv",
	".wmx":     "video/x-ms-wmx",
	".wvx":     "video/x-ms-wvx",
	".avi":     "video/x-msvideo",
	".movie":   "video/x-sgi-movie",
	".ice":     "x-conference/x-cooltalk",
	".3gp":     "video/3gpp",
	".ai":      "application/postscript",
	".aif":     "audio/x-aiff",
	".aifc":    "audio/x-aiff",
	".aiff":    "audio/x-aiff",
	".asc":     "text/plain",
	".atom":    "application/atom+xml",
	".au":      "audio/basic",
	".bin":     "application/octet-stream",
	".cdf":     "application/x-netcdf",
	".cgm":     "image/cgm",
	".class":   "application/octet-stream",
	".dcr":     "application/x-director",
	".dif":     "video/x-dv",
	".dir":     "application/x-director",
	".djv":     "image/vnd.djvu",
	".djvu":    "image/vnd.djvu",
	".dll":     "application/octet-stream",
	".dmg":     "application/octet-stream",
	".dms":     "application/octet-stream",
	".dtd":     "application/xml-dtd",
	".dv":      "video/x-dv",
	".dxr":     "application/x-director",
	".eps":     "application/postscript",
	".exe":     "application/octet-stream",
	".ez":      "application/andrew-inset",
	".gram":    "application/srgs",
	".grxml":   "application/srgs+xml",
	".gz":      "application/x-gzip",
	".htm":     "text/html",
	".html":    "text/html",
	".ico":     "image/x-icon",
	".ics":     "text/calendar",
	".ifb":     "text/calendar",
	".iges":    "model/iges",
	".igs":     "model/iges",
	".jp2":     "image/jp2",
	".jpe":     "image/jpeg",
	".jpeg":    "image/jpeg",
	".jpg":     "image/jpeg",
	".kar":     "audio/midi",
	".lha":     "application/octet-stream",
	".lzh":     "application/octet-stream",
	".m4a":     "audio/mp4a-latm",
	".m4p":     "audio/mp4a-latm",
	".m4u":     "video/vnd.mpegurl",
	".m4v":     "video/x-m4v",
	".mac":     "image/x-macpaint",
	".mathml":  "application/mathml+xml",
	".mesh":    "model/mesh",
	".mid":     "audio/midi",
	".midi":    "audio/midi",
	".mov":     "video/quicktime",
	".mp2":     "audio/mpeg",
	".mp3":     "audio/mpeg",
	".mp4":     "video/mp4",
	".mpe":     "video/mpeg",
	".mpeg":    "video/mpeg",
	".mpg":     "video/mpeg",
	".mpga":    "audio/mpeg",
	".msh":     "model/mesh",
	".nc":      "application/x-netcdf",
	".oda":     "application/oda",
	".ogv":     "video/ogv",
	".pct":     "image/pict",
	".pic":     "image/pict",
	".pict":    "image/pict",
	".pnt":     "image/x-macpaint",
	".pntg":    "image/x-macpaint",
	".ps":      "application/postscript",
	".qt":      "video/quicktime",
	".qti":     "image/x-quicktime",
	".qtif":    "image/x-quicktime",
	".ram":     "audio/x-pn-realaudio",
	".rdf":     "application/rdf+xml",
	".rm":      "application/vnd.rn-realmedia",
	".roff":    "application/x-troff",
	".sgm":     "text/sgml",
	".sgml":    "text/sgml",
	".silo":    "model/mesh",
	".skd":     "application/x-koan",
	".skm":     "application/x-koan",
	".skp":     "application/x-koan",
	".skt":     "application/x-koan",
	".smi":     "application/smil",
	".smil":    "application/smil",
	".snd":     "audio/basic",
	".so":      "application/octet-stream",
	".svg":     "image/svg+xml",
	".t":       "application/x-troff",
	".texi":    "application/x-texinfo",
	".texinfo": "application/x-texinfo",
	".tif":     "image/tiff",
	".tiff":    "image/tiff",
	".tr":      "application/x-troff",
	".txt":     "text/plain",
	".vrml":    "model/vrml",
	".vxml":    "application/voicexml+xml",
	".webm":    "video/webm",
	".wrl":     "model/vrml",
	".xht":     "application/xhtml+xml",
	".xhtml":   "application/xhtml+xml",
	".xml":     "application/xml",
	".xsl":     "application/xml",
	".xslt":    "application/xslt+xml",
	".xul":     "application/vnd.mozilla.xul+xml",
	".webp":    "image/webp",
}

func typeByExtension(filePath string) string {
	typ := mime.TypeByExtension(path.Ext(filePath))
	if typ == "" {
		typ = extToMimeType[strings.ToLower(path.Ext(filePath))]
	}
	if typ == "" {
		typ = "text/plain"
	}
	return typ
}

// Sorter defines the key-value structure for storing the sorted data in signHeader.
type Sorter struct {
	Keys []string
	Vals []string
}

// newSorter is an additional function for function Sign.
func newSorter(m map[string]string) *Sorter {
	hs := &Sorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}

// Sort is an additional function for function SignHeader.
func (hs *Sorter) Sort() {
	sort.Sort(hs)
}

// Len is an additional function for function SignHeader.
func (hs *Sorter) Len() int {
	return len(hs.Vals)
}

// Less is an additional function for function SignHeader.
func (hs *Sorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Swap is an additional function for function SignHeader.
func (hs *Sorter) Swap(i, j int) {
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}

// Determine whether the parameters are in signKeyList
// signKeyList is a list about params without value
func isParamSign(paramKey string) bool {
	for _, k := range signKeyList {
		if paramKey == k {
			return true
		}
	}
	return false
}

// Fill in the values in dataValue for result
func flatRepeatedList(dataValue reflect.Value, result map[string]string, prefix string) {
	if !dataValue.IsValid() {
		return
	}

	dataType := dataValue.Type()
	if dataType.Kind().String() == "slice" {
		handleRepeatedParams(dataValue, result, prefix)
	} else if dataType.Kind().String() == "map" {
		handleMap(dataValue, result, prefix)
	} else {
		result[prefix] = fmt.Sprintf("%v", dataValue.Interface())
	}
}

func handleRepeatedParams(repeatedFieldValue reflect.Value, result map[string]string, prefix string) {
	if repeatedFieldValue.IsValid() && !repeatedFieldValue.IsNil() {
		for m := 0; m < repeatedFieldValue.Len(); m++ {
			elementValue := repeatedFieldValue.Index(m)
			key := prefix + "." + strconv.Itoa(m+1)
			fieldValue := reflect.ValueOf(elementValue.Interface())
			if fieldValue.Kind().String() == "map" {
				handleMap(fieldValue, result, key)
			} else {
				result[key] = fmt.Sprintf("%v", fieldValue.Interface())
			}
		}
	}
}

func handleMap(valueField reflect.Value, result map[string]string, prefix string) {
	if valueField.IsValid() && valueField.String() != "" {
		valueFieldType := valueField.Type()
		if valueFieldType.Kind().String() == "map" {
			var byt []byte
			byt, _ = json.Marshal(valueField.Interface())
			cache := make(map[string]interface{})
			_ = json.Unmarshal(byt, &cache)
			for key, value := range cache {
				pre := ""
				if prefix != "" {
					pre = prefix + "." + key
				} else {
					pre = key
				}
				fieldValue := reflect.ValueOf(value)
				flatRepeatedList(fieldValue, result, pre)
			}
		}
	}
}

// Get XMl StartElement
func getStartElement(body []byte) string {
	d := xml.NewDecoder(bytes.NewReader(body))
	for {
		tok, err := d.Token()
		if err != nil {
			return ""
		}
		if t, ok := tok.(xml.StartElement); ok {
			return t.Name.Local
		}
	}
}

func getSignedStrV1(req *tea.Request, canonicalizedResource, accessKeySecret string) string {
	// Find out the "x-oss-"'s address in header of the request
	temp := make(map[string]string)

	for k, v := range req.Headers {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			temp[strings.ToLower(k)] = tea.StringValue(v)
		}
	}
	hs := newSorter(temp)

	// Sort the temp by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	// when sign URL, date is expires
	date := tea.StringValue(req.Headers["date"])
	contentType := tea.StringValue(req.Headers["content-type"])
	contentMd5 := tea.StringValue(req.Headers["content-md5"])

	signStr := tea.StringValue(req.Method) + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}

func getSignedStrV2(req *tea.Request, canonicalizedResource, accessKeySecret string, additionalHeaders []string) string {
	// Find out the "x-oss-"'s address in header of the request
	temp := make(map[string]string)
	for _, value := range additionalHeaders {
		if req.Headers[value] != nil {
			temp[strings.ToLower(value)] = tea.StringValue(req.Headers[value])
		}
	}
	for k, v := range req.Headers {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			temp[strings.ToLower(k)] = tea.StringValue(v)
		}
	}
	hs := newSorter(temp)

	// Sort the temp by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}
	// Give other parameters values
	// when sign URL, date is expires
	date := tea.StringValue(req.Headers["date"])
	contentType := tea.StringValue(req.Headers["content-type"])
	contentMd5 := tea.StringValue(req.Headers["content-md5"])
	signStr := tea.StringValue(req.Method) + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + strings.Join(additionalHeaders, ";") + "\n" + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha256.New() }, []byte(accessKeySecret))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signedStr
}

func uriEncode(rawStr string, encodeSlash bool) string {
	res := ""
	for i := 0; i < len(rawStr); i++ {
		tmp := string(rawStr[i])
		if (tmp >= "a" && tmp <= "z") || (tmp >= "A" && tmp <= "Z") ||
			(tmp >= "0" && tmp <= "9") || characters[tmp] {
			res = res + tmp
		} else if tmp == "/" {
			if encodeSlash {
				res = res + "%2F"
			} else {
				res = res + tmp
			}
		} else {
			res = res + "%" + fmt.Sprintf("%02x", tmp)
		}
	}
	return res
}

func XmlUnmarshal(body []byte, result interface{}) (interface{}, error) {
	start := getStartElement(body)
	dataValue := reflect.ValueOf(result).Elem()
	dataType := dataValue.Type()
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		name, containsNameTag := field.Tag.Lookup("xml")
		if containsNameTag {
			if name == start {
				realType := dataValue.Field(i).Type()
				realValue := reflect.New(realType).Interface()
				err := xml.Unmarshal(body, realValue)
				if err != nil {
					return nil, err
				}
				return realValue, nil
			}
		}
	}
	return nil, nil
}

func getSignatureV1(request *tea.Request, bucketName, accessKeySecret string) string {
	resource := ""
	if bucketName != "" {
		resource = "/" + bucketName
	}
	resource = resource + tea.StringValue(request.Pathname)
	if !strings.Contains(resource, "?") && len(request.Query) > 0 {
		resource += "?"
	}

	for key, value := range request.Query {
		if isParamSign(key) {
			if value != nil {
				if strings.HasSuffix(resource, "?") {
					resource = resource + key + "=" + tea.StringValue(value)
				} else {
					resource = resource + "&" + key + "=" + tea.StringValue(value)
				}
			}
		}
	}
	return getSignedStrV1(request, resource, accessKeySecret)
}

func getSignatureV2(request *tea.Request, bucketName, accessKeySecret string, additionalHeaders []string) string {
	resource := ""
	pathName := tea.StringValue(request.Pathname)
	if bucketName != "" {
		pathName = "/" + bucketName + tea.StringValue(request.Pathname)
	}

	strs := strings.Split(pathName, "?")
	resource += uriEncode(strs[0], true)
	tmp := make(map[string]string)
	for k, v := range request.Query {
		tmp[k] = tea.StringValue(v)
	}
	hs := newSorter(tmp)
	if strings.Contains(pathName, "?") {
		hs.Keys = append(hs.Keys, strs[1])
		hs.Vals = append(hs.Vals, "")
	}

	// Sort the temp by the ascending order
	hs.Sort()
	if len(hs.Keys) > 0 {
		resource += "?"
	}
	for i := range hs.Keys {
		if !strings.HasSuffix(resource, "?") {
			resource += "&"
		}
		if hs.Vals[i] == "" {
			resource += uriEncode(hs.Keys[i], true)
		} else {
			resource += uriEncode(hs.Keys[i], true) + "=" + uriEncode(hs.Vals[i], true)
		}
	}
	return getSignedStrV2(request, resource, accessKeySecret, additionalHeaders)
}

// Decryption
func base64Decode(value string) string {
	strs := strings.Split(value, "/")
	result, err := base64.StdEncoding.DecodeString(strs[len(strs)-1])
	if err != nil {
		return ""
	}
	strs[len(strs)-1] = string(result)
	return strings.Join(strs, "/")
}

// Decryption
func urlDecode(value string) string {
	strs := strings.Split(value, "/")
	result, err := url.QueryUnescape(strs[len(strs)-1])
	if err != nil {
		return ""
	}
	strs[len(strs)-1] = result
	return strings.Join(strs, "/")
}

func listToString(a []string, sep string) string {
	return strings.Join(a, sep)
}
