package awsauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Credentials stores the information necessary to authorize
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SecurityToken   string `json:"Token"`
	Expiration      time.Time
}

type metadata struct {
	algorithm       string
	credentialScope string
	signedHeaders   string
	date            string
	region          string
	service         string
}

// Sign signs a request with Signed Signature Version 4.
func Sign(request *http.Request, credentials ...Credentials) *http.Request {
	keys := credentials[0]

	// Add the X-Amz-Security-Token header when using STS
	if keys.SecurityToken != "" {
		request.Header.Set("X-Amz-Security-Token", keys.SecurityToken)
	}

	prepareRequestV4(request)
	meta := new(metadata)

	// Task 1
	hashedCanonReq := hashedCanonicalRequestV4(request, meta)

	// Task 2
	stringToSign := stringToSignV4(request, hashedCanonReq, meta)

	// Task 3
	signingKey := signingKeyV4(keys.SecretAccessKey, meta.date, meta.region, meta.service)
	signature := signatureV4(signingKey, stringToSign)

	request.Header.Set("Authorization", buildAuthHeaderV4(signature, meta, keys))

	return request
}

func serviceAndRegion(host string) (service string, region string) {
	// These are the defaults if the hostname doesn't suggest something else
	region = "us-east-1"
	service = "s3"

	parts := strings.Split(host, ".")
	if len(parts) == 4 {
		// Either service.region.amazonaws.com or virtual-host.region.amazonaws.com
		if parts[1] == "s3" {
			service = "s3"
		} else if strings.HasPrefix(parts[1], "s3-") {
			region = parts[1][3:]
			service = "s3"
		} else {
			service = parts[0]
			region = parts[1]
		}
	} else if len(parts) == 5 {
		service = parts[2]
		region = parts[1]
	} else {
		// Either service.amazonaws.com or s3-region.amazonaws.com
		if strings.HasPrefix(parts[0], "s3-") {
			region = parts[0][3:]
		} else {
			service = parts[0]
		}
	}

	if region == "external-1" {
		region = "us-east-1"
	}

	return
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func hashSHA256(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func readAndReplaceBody(request *http.Request) []byte {
	if request.Body == nil {
		return []byte{}
	}
	payload, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewReader(payload))
	return payload
}

func concat(delim string, str ...string) string {
	return strings.Join(str, delim)
}

var now = func() time.Time {
	return time.Now().UTC()
}

func normuri(uri string) string {
	parts := strings.Split(uri, "/")
	for i := range parts {
		parts[i] = encodePathFrag(parts[i])
	}
	return strings.Join(parts, "/")
}

func encodePathFrag(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}
	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		} else {
			t[j] = c
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		return false
	}
	if '0' <= c && c <= '9' {
		return false
	}
	if c == '-' || c == '_' || c == '.' || c == '~' {
		return false
	}
	return true
}

func normquery(v url.Values) string {
	queryString := v.Encode()
	return strings.Replace(queryString, "+", "%20", -1)
}

func hashedCanonicalRequestV4(request *http.Request, meta *metadata) string {
	// TASK 1. http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html

	payload := readAndReplaceBody(request)
	payloadHash := hashSHA256(payload)
	request.Header.Set("X-Amz-Content-Sha256", payloadHash)

	// Set this in header values to make it appear in the range of headers to sign
	request.Header.Set("Host", request.Host)

	var sortedHeaderKeys []string
	for key, _ := range request.Header {
		switch key {
		case "Content-Type", "Content-Md5", "Host":
		default:
			if !strings.HasPrefix(key, "X-Amz-") {
				continue
			}
		}
		sortedHeaderKeys = append(sortedHeaderKeys, strings.ToLower(key))
	}
	sort.Strings(sortedHeaderKeys)

	var headersToSign string
	for _, key := range sortedHeaderKeys {
		value := strings.TrimSpace(request.Header.Get(key))
		if key == "host" {
			//AWS does not include port in signing request.
			if strings.Contains(value, ":") {
				split := strings.Split(value, ":")
				port := split[1]
				if port == "80" || port == "443" {
					value = split[0]
				}
			}
		}
		headersToSign += key + ":" + value + "\n"
	}
	meta.signedHeaders = concat(";", sortedHeaderKeys...)
	canonicalRequest := concat("\n", request.Method, normuri(request.URL.Path), normquery(request.URL.Query()), headersToSign, meta.signedHeaders, payloadHash)

	return hashSHA256([]byte(canonicalRequest))
}

func stringToSignV4(request *http.Request, hashedCanonReq string, meta *metadata) string {
	// TASK 2. http://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html

	requestTs := request.Header.Get("X-Amz-Date")

	meta.algorithm = "AWS4-HMAC-SHA256"
	meta.service, meta.region = serviceAndRegion(request.Host)
	meta.date = tsDateV4(requestTs)
	meta.credentialScope = concat("/", meta.date, meta.region, meta.service, "aws4_request")

	return concat("\n", meta.algorithm, requestTs, meta.credentialScope, hashedCanonReq)
}

func signatureV4(signingKey []byte, stringToSign string) string {
	// TASK 3. http://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html

	return hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
}

func prepareRequestV4(request *http.Request) *http.Request {
	necessaryDefaults := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
		"X-Amz-Date":   timestampV4(),
	}

	for header, value := range necessaryDefaults {
		if request.Header.Get(header) == "" {
			request.Header.Set(header, value)
		}
	}

	if request.URL.Path == "" {
		request.URL.Path += "/"
	}

	return request
}

func signingKeyV4(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

func buildAuthHeaderV4(signature string, meta *metadata, keys Credentials) string {
	credential := keys.AccessKeyID + "/" + meta.credentialScope

	return meta.algorithm +
		" Credential=" + credential +
		", SignedHeaders=" + meta.signedHeaders +
		", Signature=" + signature
}

func timestampV4() string {
	return now().Format(timeFormatV4)
}

func tsDateV4(timestamp string) string {
	return timestamp[:8]
}

const timeFormatV4 = "20060102T150405Z"
