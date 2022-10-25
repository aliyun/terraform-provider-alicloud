package fc

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SignURLInput ...
type SignURLInput struct {
	Method         string
	Expires        time.Time
	ServiceName    string
	Qualifier      string
	FunctionName   string
	CustomEndpoint string
	EscapedPath    string
	Header         http.Header
	Queries        url.Values
}

// NewSignURLInput : create sign url request
func NewSignURLInput(method, serviceName, functionName string, expires time.Time) *SignURLInput {
	return &SignURLInput{
		ServiceName:  serviceName,
		FunctionName: functionName,
		Method:       method,
		Expires:      expires,
	}
}

func (s *SignURLInput) signURL(apiVersion, endpoint, akID, akSec, stsToken string) (string, error) {
	orinalEscapedPath := s.EscapedPath
	if orinalEscapedPath == "" {
		orinalEscapedPath = "/"
	}
	unescapedPath, err := url.PathUnescape(orinalEscapedPath)
	if err != nil {
		return "", err
	}

	var queries url.Values = map[string][]string{}
	for k, values := range s.Queries {
		for _, v := range values {
			queries.Add(k, v)
		}
	}
	queries.Set(AuthQueryKeyExpires, fmt.Sprintf("%d", s.Expires.Unix()))
	if stsToken != "" {
		queries.Set(AuthQueryKeySecurityToken, stsToken)
	}
	queries.Set(AuthQueryKeyAccessKeyID, akID)

	singleHeaders := map[string]string{}
	for k, v := range s.Header {
		if len(v) > 0 {
			singleHeaders[http.CanonicalHeaderKey(k)] = v[0]
		}
	}

	serviceWithQualifier := s.ServiceName
	if s.Qualifier != "" {
		serviceWithQualifier = s.ServiceName + "." + s.Qualifier
	}

	unescapedPathForSign := fmt.Sprintf("/%s/proxy/%s/%s%s", apiVersion,
		serviceWithQualifier, s.FunctionName, unescapedPath)

	escapedPath := fmt.Sprintf("/%s/proxy/%s/%s%s", apiVersion,
		serviceWithQualifier, s.FunctionName, orinalEscapedPath)

	if s.CustomEndpoint != "" {
		_, err := url.Parse(s.CustomEndpoint)
		if s.CustomEndpoint == endpoint || err != nil {
			return "", fmt.Errorf("invalid custom endpoint: %s", s.CustomEndpoint)
		}
		endpoint = s.CustomEndpoint
		escapedPath = orinalEscapedPath
		unescapedPathForSign = unescapedPath
	}

	fcResource := GetSignResourceWithQueries(unescapedPathForSign, queries)
	signature := GetSignature(akSec, s.Method, singleHeaders, fcResource)

	queries.Set(AuthQueryKeySignature, signature)

	return fmt.Sprintf("%s%s?%s", endpoint, escapedPath, queries.Encode()), nil
}

// WithQualifier ...
func (s *SignURLInput) WithQualifier(qualifier string) *SignURLInput {
	s.Qualifier = qualifier
	return s
}

// WithCustomEndpoint ...
func (s *SignURLInput) WithCustomEndpoint(customEndpoint string) *SignURLInput {
	s.CustomEndpoint = customEndpoint
	return s
}

// WithEscapedPath ...
func (s *SignURLInput) WithEscapedPath(esapedPath string) *SignURLInput {
	s.EscapedPath = esapedPath
	return s
}

// WithHeader ...
func (s *SignURLInput) WithHeader(header http.Header) *SignURLInput {
	s.Header = header
	return s
}

// WithQueries ...
func (s *SignURLInput) WithQueries(queries url.Values) *SignURLInput {
	s.Queries = queries
	return s
}
