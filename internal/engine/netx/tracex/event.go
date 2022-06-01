package tracex

import (
	"crypto/x509"
	"net/http"
	"time"

	"github.com/ooni/probe-cli/v3/internal/netxlite"
)

// Event is one of the events within a trace.
type Event interface {
	// Value returns the event value
	Value() *EventValue

	// Name returns the event name
	Name() string
}

// EventTLSHandshakeStart is the beginning of the TLS handshake.
type EventTLSHandshakeStart struct {
	V *EventValue
}

func (ev *EventTLSHandshakeStart) Value() *EventValue {
	return ev.V
}

func (ev *EventTLSHandshakeStart) Name() string {
	return "tls_handshake_start"
}

// EventTLSHandshakeDone is the end of the TLS handshake.
type EventTLSHandshakeDone struct {
	V *EventValue
}

func (ev *EventTLSHandshakeDone) Value() *EventValue {
	return ev.V
}

func (ev *EventTLSHandshakeDone) Name() string {
	return "tls_handshake_done"
}

// EventResolveStart is the beginning of a DNS lookup operation.
type EventResolveStart struct {
	V *EventValue
}

func (ev *EventResolveStart) Value() *EventValue {
	return ev.V
}

func (ev *EventResolveStart) Name() string {
	return "resolve_start"
}

// EventResolveDone is the end of a DNS lookup operation.
type EventResolveDone struct {
	V *EventValue
}

func (ev *EventResolveDone) Value() *EventValue {
	return ev.V
}

func (ev *EventResolveDone) Name() string {
	return "resolve_done"
}

// EventDNSRoundTripStart is the start of a DNS round trip.
type EventDNSRoundTripStart struct {
	V *EventValue
}

func (ev *EventDNSRoundTripStart) Value() *EventValue {
	return ev.V
}

func (ev *EventDNSRoundTripStart) Name() string {
	return "dns_round_trip_start"
}

// EventDNSRoundTripDone is the end of a DNS round trip.
type EventDNSRoundTripDone struct {
	V *EventValue
}

func (ev *EventDNSRoundTripDone) Value() *EventValue {
	return ev.V
}

func (ev *EventDNSRoundTripDone) Name() string {
	return "dns_round_trip_done"
}

// EventQUICHandshakeStart is the start of a QUIC handshake.
type EventQUICHandshakeStart struct {
	V *EventValue
}

func (ev *EventQUICHandshakeStart) Value() *EventValue {
	return ev.V
}

func (ev *EventQUICHandshakeStart) Name() string {
	return "quic_handshake_start"
}

// EventQUICHandshakeDone is the end of a QUIC handshake.
type EventQUICHandshakeDone struct {
	V *EventValue
}

func (ev *EventQUICHandshakeDone) Value() *EventValue {
	return ev.V
}

func (ev *EventQUICHandshakeDone) Name() string {
	return "quic_handshake_done"
}

// EventWriteToOperation summarizes the WriteTo operation.
type EventWriteToOperation struct {
	V *EventValue
}

func (ev *EventWriteToOperation) Value() *EventValue {
	return ev.V
}

func (ev *EventWriteToOperation) Name() string {
	return netxlite.WriteToOperation
}

// EventReadFromOperation summarizes the ReadFrom operation.
type EventReadFromOperation struct {
	V *EventValue
}

func (ev *EventReadFromOperation) Value() *EventValue {
	return ev.V
}

func (ev *EventReadFromOperation) Name() string {
	return netxlite.ReadFromOperation
}

// EventHTTPRequestMetadata contains HTTP request metadata.
type EventHTTPRequestMetadata struct {
	V *EventValue
}

func (ev *EventHTTPRequestMetadata) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPRequestMetadata) Name() string {
	return "http_request_metadata"
}

// EventHTTPResponseMetadata contains HTTP response metadata.
type EventHTTPResponseMetadata struct {
	V *EventValue
}

func (ev *EventHTTPResponseMetadata) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPResponseMetadata) Name() string {
	return "http_response_metadata"
}

// EventHTTPTransactionStart is the beginning of an HTTP transaction.
type EventHTTPTransactionStart struct {
	V *EventValue
}

func (ev *EventHTTPTransactionStart) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPTransactionStart) Name() string {
	return "http_transaction_start"
}

// EventHTTPTransactionDone is the end of an HTTP transaction.
type EventHTTPTransactionDone struct {
	V *EventValue
}

func (ev *EventHTTPTransactionDone) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPTransactionDone) Name() string {
	return "http_transaction_done"
}

// EventHTTPRequestBodySnapshot contains a snapshot of the request body.
type EventHTTPRequestBodySnapshot struct {
	V *EventValue
}

func (ev *EventHTTPRequestBodySnapshot) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPRequestBodySnapshot) Name() string {
	return "http_request_body_snapshot"
}

// EventHTTPResponseBodySnapshot contains a snapshot of the response body.
type EventHTTPResponseBodySnapshot struct {
	V *EventValue
}

func (ev *EventHTTPResponseBodySnapshot) Value() *EventValue {
	return ev.V
}

func (ev *EventHTTPResponseBodySnapshot) Name() string {
	return "http_response_body_snapshot"
}

// EventConnectOperation contains information about the connect operation.
type EventConnectOperation struct {
	V *EventValue
}

func (ev *EventConnectOperation) Value() *EventValue {
	return ev.V
}

func (ev *EventConnectOperation) Name() string {
	return netxlite.ConnectOperation
}

// EventReadOperation contains information about a read operation.
type EventReadOperation struct {
	V *EventValue
}

func (ev *EventReadOperation) Value() *EventValue {
	return ev.V
}

func (ev *EventReadOperation) Name() string {
	return netxlite.ReadOperation
}

// EventWriteOperation contains information about a write operation.
type EventWriteOperation struct {
	V *EventValue
}

func (ev *EventWriteOperation) Value() *EventValue {
	return ev.V
}

func (ev *EventWriteOperation) Name() string {
	return netxlite.WriteOperation
}

// Event is one of the events within a trace
type EventValue struct {
	Addresses          []string            `json:",omitempty"`
	Address            string              `json:",omitempty"`
	DNSQuery           []byte              `json:",omitempty"`
	DNSReply           []byte              `json:",omitempty"`
	DataIsTruncated    bool                `json:",omitempty"`
	Data               []byte              `json:",omitempty"`
	Duration           time.Duration       `json:",omitempty"`
	Err                error               `json:",omitempty"`
	HTTPHeaders        http.Header         `json:",omitempty"`
	HTTPMethod         string              `json:",omitempty"`
	HTTPStatusCode     int                 `json:",omitempty"`
	HTTPURL            string              `json:",omitempty"`
	Hostname           string              `json:",omitempty"`
	NoTLSVerify        bool                `json:",omitempty"`
	NumBytes           int                 `json:",omitempty"`
	Proto              string              `json:",omitempty"`
	TLSServerName      string              `json:",omitempty"`
	TLSCipherSuite     string              `json:",omitempty"`
	TLSNegotiatedProto string              `json:",omitempty"`
	TLSNextProtos      []string            `json:",omitempty"`
	TLSPeerCerts       []*x509.Certificate `json:",omitempty"`
	TLSVersion         string              `json:",omitempty"`
	Time               time.Time           `json:",omitempty"`
	Transport          string              `json:",omitempty"`
}