// Package v4 implements signing for AWS V4 signer
//
// Provides request signing for request that need to be signed with
// AWS V4 Signatures.
//
// Standalone Signer
//
// Generally using the signer outside of the SDK should not require any additional
//  The signer does this by taking advantage of the URL.EscapedPath method. If your request URI requires
// additional escaping you many need to use the URL.Opaque to define what the raw URI should be sent
// to the service as.
//
// The signer will first check the URL.Opaque field, and use its value if set.
// The signer does require the URL.Opaque field to be set in the form of:
//
//     "//<hostname>/<path>"
//
//     // e.g.
//     "//example.com/some/path"
//
// The leading "//" and hostname are required or the URL.Opaque escaping will
// not work correctly.
//
// If URL.Opaque is not set the signer will fallback to the URL.EscapedPath()
// method and using the returned value.
//
// AWS v4 signature validation requires that the canonical string's URI path
// element must be the URI escaped form of the HTTP request's path.
// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
//
// The Go HTTP client will perform escaping automatically on the request. Some
// of these escaping may cause signature validation errors because the HTTP
// request differs from the URI path or query that the signature was generated.
// https://golang.org/pkg/net/url/#URL.EscapedPath
//
// Because of this, it is recommended that when using the signer outside of the
// SDK that explicitly escaping the request prior to being signed is preferable,
// and will help prevent signature validation errors. This can be done by setting
// the URL.Opaque or URL.RawPath. The SDK will use URL.Opaque first and then
// call URL.EscapedPath() if Opaque is not set.
//
// Test `TestStandaloneSign` provides a complete example of using the signer
// outside of the SDK and pre-escaping the URI path.
package v4

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4Internal "github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4"
	"github.com/awslabs/smithy-go/httpbinding"
)

const (
	signingAlgorithm    = "AWS4-HMAC-SHA256"
	authorizationHeader = "Authorization"
)

// HTTPSigner is an interface to a SigV4 signer that can sign HTTP requests
type HTTPSigner interface {
	SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error
}

type keyDerivator interface {
	DeriveKey(credential aws.Credentials, service, region string, signingTime v4Internal.SigningTime) []byte
}

// Signer applies AWS v4 signing to given request. Use this to sign requests
// that need to be signed with AWS V4 Signatures.
type Signer struct {
	// Disables the Signer's moving HTTP header key/value pairs from the HTTP
	// request header to the request's query string. This is most commonly used
	// with pre-signed requests preventing headers from being added to the
	// request's query string.
	DisableHeaderHoisting bool

	// Disables the automatic escaping of the URI path of the request for the
	// siganture's canonical string's path. For services that do not need additional
	// escaping then use this to disable the signer escaping the path.
	//
	// S3 is an example of a service that does not need additional escaping.
	//
	// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
	DisableURIPathEscaping bool

	keyDerivator keyDerivator
}

// NewSigner returns a new SigV4 Signer
func NewSigner(optFns ...func(signer *Signer)) *Signer {
	s := &Signer{keyDerivator: v4Internal.NewSigningKeyDeriver()}

	for _, fn := range optFns {
		fn(s)
	}

	return s
}

type httpSigner struct {
	Request      *http.Request
	ServiceName  string
	Region       string
	Time         v4Internal.SigningTime
	ExpireTime   time.Duration
	Credentials  aws.Credentials
	KeyDerivator keyDerivator
	IsPreSign    bool

	// PayloadHash is the hex encoded SHA-256 hash of the request payload
	// If len(PayloadHash) == 0 the signer will attempt to send the request
	// as an unsigned payload. Note: Unsigned payloads only work for a subset of services.
	PayloadHash string

	DisableHeaderHoisting  bool
	DisableURIPathEscaping bool
}

func (s *httpSigner) Build() (signedRequest, error) {
	req := s.Request

	query := req.URL.Query()
	headers := req.Header

	s.setRequiredSigningFields(headers, query)

	// Sort Each Query Key's Values
	for key := range query {
		sort.Strings(query[key])
	}

	v4Internal.SanitizeHostForHeader(req)

	credentialScope := s.buildCredentialScope()
	credentialStr := s.Credentials.AccessKeyID + "/" + credentialScope
	if s.IsPreSign {
		query.Set(v4Internal.AmzCredentialKey, credentialStr)
	}

	unsignedHeaders := headers
	if s.IsPreSign && !s.DisableHeaderHoisting {
		var urlValues url.Values
		urlValues, unsignedHeaders = buildQuery(v4Internal.AllowedQueryHoisting, headers)
		for k := range urlValues {
			query[k] = urlValues[k]
		}
	}

	host := req.URL.Host
	if len(req.Host) > 0 {
		host = req.Host
	}

	signedHeaders, signedHeadersStr, canonicalHeaderStr := s.buildCanonicalHeaders(host, v4Internal.IgnoredHeaders, unsignedHeaders, s.Request.ContentLength)

	if s.IsPreSign {
		query.Set(v4Internal.AmzSignedHeadersKey, signedHeadersStr)
	}

	var rawQuery strings.Builder
	rawQuery.WriteString(strings.Replace(query.Encode(), "+", "%20", -1))

	canonicalURI := v4Internal.GetURIPath(req.URL)
	if !s.DisableURIPathEscaping {
		canonicalURI = httpbinding.EscapePath(canonicalURI, false)
	}

	canonicalString := s.buildCanonicalString(
		req.Method,
		canonicalURI,
		rawQuery.String(),
		signedHeadersStr,
		canonicalHeaderStr,
	)

	strToSign := s.buildStringToSign(credentialScope, canonicalString)
	signingSignature, err := s.buildSignature(strToSign)
	if err != nil {
		return signedRequest{}, err
	}

	if s.IsPreSign {
		rawQuery.WriteString("&X-Amz-Signature=")
		rawQuery.WriteString(signingSignature)
	} else {
		headers[authorizationHeader] = append(headers[authorizationHeader][:0], buildAuthorizationHeader(credentialStr, signedHeadersStr, signingSignature))
	}

	req.URL.RawQuery = rawQuery.String()

	return signedRequest{
		Request:         req,
		SignedHeaders:   signedHeaders,
		CanonicalString: canonicalString,
		StringToSign:    strToSign,
		PreSigned:       s.IsPreSign,
	}, nil
}

func buildAuthorizationHeader(credentialStr, signedHeadersStr, signingSignature string) string {
	const credential = "Credential="
	const signedHeaders = "SignedHeaders="
	const signature = "Signature="
	const commaSpace = ", "

	var parts strings.Builder
	parts.Grow(len(signingAlgorithm) + 1 +
		len(credential) + len(credentialStr) + 2 +
		len(signedHeaders) + len(signedHeadersStr) + 2 +
		len(signature) + len(signingSignature),
	)
	parts.WriteString(signingAlgorithm)
	parts.WriteRune(' ')
	parts.WriteString(credential)
	parts.WriteString(credentialStr)
	parts.WriteString(commaSpace)
	parts.WriteString(signedHeaders)
	parts.WriteString(signedHeadersStr)
	parts.WriteString(commaSpace)
	parts.WriteString(signature)
	parts.WriteString(signingSignature)
	return parts.String()
}

// SignHTTP signs AWS v4 requests with the provided payload hash, service name, region the
// request is made to, and time the request is signed at. The signTime allows
// you to specify that a request is signed for the future, and cannot be
// used until then.
//
// Sign differs from Presign in that it will sign the request using HTTP
// header values. This type of signing is intended for http.Request values that
// will not be shared, or are shared in a way the header values on the request
// will not be lost.
//
// The passed in request will be modified in place.
func (v4 Signer) SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error {
	signer := &httpSigner{
		Request:                r,
		PayloadHash:            payloadHash,
		ServiceName:            service,
		Region:                 region,
		Credentials:            credentials,
		Time:                   v4Internal.NewSigningTime(signingTime.UTC()),
		DisableHeaderHoisting:  v4.DisableHeaderHoisting,
		DisableURIPathEscaping: v4.DisableURIPathEscaping,
		KeyDerivator:           v4.keyDerivator,
	}

	_, err := signer.Build()
	if err != nil {
		return err
	}

	return nil
}

// PresignHTTP signs AWS v4 requests with the payload hash, service name, region
// the request is made to, and time the request is signed at. The signTime
// allows you to specify that a request is signed for the future, and cannot
// be used until then.
//
// Returns the signed URL and the map of HTTP headers that were included in the signature or an
// error if signing the request failed. For presigned requests these headers
// and their values must be included on the HTTP request when it is made. This
// is helpful to know what header values need to be shared with the party the
// presigned request will be distributed to.
//
// PresignHTTP differs from SignHTTP in that it will sign the request using query string
// instead of header values. This allows you to share the Presigned Request's
// URL with third parties, or distribute it throughout your system with minimal
// dependencies.
//
// PresignHTTP also takes an exp value which is the duration the
// signed request will be valid after the signing time. This is allows you to
// set when the request will expire.
//
// This method does not modify the provided request.
func (v4 *Signer) PresignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, expireTime time.Duration, signingTime time.Time) (signedURI string, signedHeaders http.Header, err error) {
	signer := &httpSigner{
		Request:                r.Clone(r.Context()),
		PayloadHash:            payloadHash,
		ServiceName:            service,
		Region:                 region,
		Credentials:            credentials,
		Time:                   v4Internal.NewSigningTime(signingTime.UTC()),
		IsPreSign:              true,
		ExpireTime:             expireTime,
		DisableHeaderHoisting:  v4.DisableHeaderHoisting,
		DisableURIPathEscaping: v4.DisableURIPathEscaping,
		KeyDerivator:           v4.keyDerivator,
	}

	signedRequest, err := signer.Build()
	if err != nil {
		return "", nil, err
	}

	return signedRequest.Request.URL.String(), signedRequest.SignedHeaders, nil
}

func (s *httpSigner) buildCredentialScope() string {
	return strings.Join([]string{
		s.Time.ShortTimeFormat(),
		s.Region,
		s.ServiceName,
		"aws4_request",
	}, "/")
}

func buildQuery(r v4Internal.Rule, header http.Header) (url.Values, http.Header) {
	query := url.Values{}
	unsignedHeaders := http.Header{}
	for k, h := range header {
		if r.IsValid(k) {
			query[k] = h
		} else {
			unsignedHeaders[k] = h
		}
	}

	return query, unsignedHeaders
}

func (s *httpSigner) buildCanonicalHeaders(host string, rule v4Internal.Rule, header http.Header, length int64) (signed http.Header, signedHeaders, canonicalHeadersStr string) {
	signed = make(http.Header)

	var headers []string
	const hostHeader = "host"
	headers = append(headers, hostHeader)
	signed[hostHeader] = append(signed[hostHeader], host)

	if length > 0 {
		const contentLengthHeader = "content-length"
		headers = append(headers, contentLengthHeader)
		signed[contentLengthHeader] = append(signed[contentLengthHeader], strconv.FormatInt(length, 10))
	}

	for k, v := range header {
		if !rule.IsValid(k) {
			continue // ignored header
		}

		lowerCaseKey := strings.ToLower(k)
		if _, ok := signed[lowerCaseKey]; ok {
			// include additional values
			signed[lowerCaseKey] = append(signed[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		signed[lowerCaseKey] = v
	}
	sort.Strings(headers)

	signedHeaders = strings.Join(headers, ";")

	var canonicalHeaders strings.Builder
	n := len(headers)
	const colon = ':'
	for i := 0; i < n; i++ {
		if headers[i] == hostHeader {
			canonicalHeaders.WriteString(hostHeader)
			canonicalHeaders.WriteRune(colon)
			canonicalHeaders.WriteString(v4Internal.StripExcessSpaces(host))
		} else {
			canonicalHeaders.WriteString(headers[i])
			canonicalHeaders.WriteRune(colon)
			canonicalHeaders.WriteString(strings.Join(signed[headers[i]], ","))
		}
		canonicalHeaders.WriteRune('\n')
	}
	canonicalHeadersStr = canonicalHeaders.String()

	return signed, signedHeaders, canonicalHeadersStr
}

func (s *httpSigner) buildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders,
		signedHeaders,
		s.PayloadHash,
	}, "\n")
}

func (s *httpSigner) buildStringToSign(credentialScope, canonicalRequestString string) string {
	return strings.Join([]string{
		signingAlgorithm,
		s.Time.TimeFormat(),
		credentialScope,
		hex.EncodeToString(makeHash(sha256.New(), []byte(canonicalRequestString))),
	}, "\n")
}

func makeHash(hash hash.Hash, b []byte) []byte {
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func (s *httpSigner) buildSignature(strToSign string) (string, error) {
	key := s.KeyDerivator.DeriveKey(s.Credentials, s.ServiceName, s.Region, s.Time)
	return hex.EncodeToString(v4Internal.HMACSHA256(key, []byte(strToSign))), nil
}

func (s *httpSigner) setRequiredSigningFields(headers http.Header, query url.Values) {
	amzDate := s.Time.TimeFormat()

	if s.IsPreSign {
		query.Set(v4Internal.AmzAlgorithmKey, signingAlgorithm)
		if sessionToken := s.Credentials.SessionToken; len(sessionToken) > 0 {
			query.Set("X-Amz-Security-Token", sessionToken)
		}

		duration := int64(s.ExpireTime / time.Second)
		query.Set(v4Internal.AmzDateKey, amzDate)
		query.Set(v4Internal.AmzExpiresKey, strconv.FormatInt(duration, 10))
		return
	}

	headers[v4Internal.AmzDateKey] = append(headers[v4Internal.AmzDateKey][:0], amzDate)

	if len(s.Credentials.SessionToken) > 0 {
		headers[v4Internal.AmzSecurityTokenKey] = append(headers[v4Internal.AmzSecurityTokenKey][:0], s.Credentials.SessionToken)
	}
}

type signedRequest struct {
	Request         *http.Request
	SignedHeaders   http.Header
	CanonicalString string
	StringToSign    string
	PreSigned       bool
}
