package tracex

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/ooni/probe-cli/v3/internal/model"
	"github.com/ooni/probe-cli/v3/internal/netxlite"
)

func TestSaverTLSHandshakerSuccessWithReadWrite(t *testing.T) {
	// This is the most common use case for collecting reads, writes
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	nextprotos := []string{"h2"}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Config: &tls.Config{NextProtos: nextprotos},
		Dialer: netxlite.NewDialerWithResolver(
			model.DiscardLogger,
			netxlite.NewResolverStdlib(model.DiscardLogger),
			saver.NewReadWriteObserver(),
		),
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	// Implementation note: we don't close the connection here because it is
	// very handy to have the last event being the end of the handshake
	_, err := tlsdlr.DialTLSContext(context.Background(), "tcp", "www.google.com:443")
	if err != nil {
		t.Fatal(err)
	}
	ev := saver.Read()
	if len(ev) < 4 {
		// it's a bit tricky to be sure about the right number of
		// events because network conditions may influence that
		t.Fatal("unexpected number of events")
	}
	if ev[0].Name() != "tls_handshake_start" {
		t.Fatal("unexpected Name")
	}
	if ev[0].Value().TLSServerName != "www.google.com" {
		t.Fatal("unexpected TLSServerName")
	}
	if !reflect.DeepEqual(ev[0].Value().TLSNextProtos, nextprotos) {
		t.Fatal("unexpected TLSNextProtos")
	}
	if ev[0].Value().Time.After(time.Now()) {
		t.Fatal("unexpected Time")
	}
	last := len(ev) - 1
	for idx := 1; idx < last; idx++ {
		if ev[idx].Value().Data == nil {
			t.Fatal("unexpected Data")
		}
		if ev[idx].Value().Duration <= 0 {
			t.Fatal("unexpected Duration")
		}
		if ev[idx].Value().Err != nil {
			t.Fatal("unexpected Err")
		}
		if ev[idx].Value().NumBytes <= 0 {
			t.Fatal("unexpected NumBytes")
		}
		switch ev[idx].Name() {
		case netxlite.ReadOperation, netxlite.WriteOperation:
		default:
			t.Fatal("unexpected Name")
		}
		if ev[idx].Value().Time.Before(ev[idx-1].Value().Time) {
			t.Fatal("unexpected Time")
		}
	}
	if ev[last].Value().Duration <= 0 {
		t.Fatal("unexpected Duration")
	}
	if ev[last].Value().Err != nil {
		t.Fatal("unexpected Err")
	}
	if ev[last].Name() != "tls_handshake_done" {
		t.Fatal("unexpected Name")
	}
	if ev[last].Value().TLSCipherSuite == "" {
		t.Fatal("unexpected TLSCipherSuite")
	}
	if ev[last].Value().TLSNegotiatedProto != "h2" {
		t.Fatal("unexpected TLSNegotiatedProto")
	}
	if !reflect.DeepEqual(ev[last].Value().TLSNextProtos, nextprotos) {
		t.Fatal("unexpected TLSNextProtos")
	}
	if ev[last].Value().TLSPeerCerts == nil {
		t.Fatal("unexpected TLSPeerCerts")
	}
	if ev[last].Value().TLSServerName != "www.google.com" {
		t.Fatal("unexpected TLSServerName")
	}
	if ev[last].Value().TLSVersion == "" {
		t.Fatal("unexpected TLSVersion")
	}
	if ev[last].Value().Time.Before(ev[last-1].Value().Time) {
		t.Fatal("unexpected Time")
	}
}

func TestSaverTLSHandshakerSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	nextprotos := []string{"h2"}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Config:        &tls.Config{NextProtos: nextprotos},
		Dialer:        &netxlite.DialerSystem{},
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	conn, err := tlsdlr.DialTLSContext(context.Background(), "tcp", "www.google.com:443")
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	ev := saver.Read()
	if len(ev) != 2 {
		t.Fatal("unexpected number of events")
	}
	if ev[0].Name() != "tls_handshake_start" {
		t.Fatal("unexpected Name")
	}
	if ev[0].Value().TLSServerName != "www.google.com" {
		t.Fatal("unexpected TLSServerName")
	}
	if !reflect.DeepEqual(ev[0].Value().TLSNextProtos, nextprotos) {
		t.Fatal("unexpected TLSNextProtos")
	}
	if ev[0].Value().Time.After(time.Now()) {
		t.Fatal("unexpected Time")
	}
	if ev[1].Value().Duration <= 0 {
		t.Fatal("unexpected Duration")
	}
	if ev[1].Value().Err != nil {
		t.Fatal("unexpected Err")
	}
	if ev[1].Name() != "tls_handshake_done" {
		t.Fatal("unexpected Name")
	}
	if ev[1].Value().TLSCipherSuite == "" {
		t.Fatal("unexpected TLSCipherSuite")
	}
	if ev[1].Value().TLSNegotiatedProto != "h2" {
		t.Fatal("unexpected TLSNegotiatedProto")
	}
	if !reflect.DeepEqual(ev[1].Value().TLSNextProtos, nextprotos) {
		t.Fatal("unexpected TLSNextProtos")
	}
	if ev[1].Value().TLSPeerCerts == nil {
		t.Fatal("unexpected TLSPeerCerts")
	}
	if ev[1].Value().TLSServerName != "www.google.com" {
		t.Fatal("unexpected TLSServerName")
	}
	if ev[1].Value().TLSVersion == "" {
		t.Fatal("unexpected TLSVersion")
	}
	if ev[1].Value().Time.Before(ev[0].Value().Time) {
		t.Fatal("unexpected Time")
	}
}

func TestSaverTLSHandshakerHostnameError(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Dialer:        &netxlite.DialerSystem{},
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	conn, err := tlsdlr.DialTLSContext(
		context.Background(), "tcp", "wrong.host.badssl.com:443")
	if err == nil {
		t.Fatal("expected an error here")
	}
	if conn != nil {
		t.Fatal("expected nil conn here")
	}
	for _, ev := range saver.Read() {
		if ev.Name() != "tls_handshake_done" {
			continue
		}
		if ev.Value().NoTLSVerify == true {
			t.Fatal("expected NoTLSVerify to be false")
		}
		if len(ev.Value().TLSPeerCerts) < 1 {
			t.Fatal("expected at least a certificate here")
		}
	}
}

func TestSaverTLSHandshakerInvalidCertError(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Dialer:        &netxlite.DialerSystem{},
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	conn, err := tlsdlr.DialTLSContext(
		context.Background(), "tcp", "expired.badssl.com:443")
	if err == nil {
		t.Fatal("expected an error here")
	}
	if conn != nil {
		t.Fatal("expected nil conn here")
	}
	for _, ev := range saver.Read() {
		if ev.Name() != "tls_handshake_done" {
			continue
		}
		if ev.Value().NoTLSVerify == true {
			t.Fatal("expected NoTLSVerify to be false")
		}
		if len(ev.Value().TLSPeerCerts) < 1 {
			t.Fatal("expected at least a certificate here")
		}
	}
}

func TestSaverTLSHandshakerAuthorityError(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Dialer:        &netxlite.DialerSystem{},
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	conn, err := tlsdlr.DialTLSContext(
		context.Background(), "tcp", "self-signed.badssl.com:443")
	if err == nil {
		t.Fatal("expected an error here")
	}
	if conn != nil {
		t.Fatal("expected nil conn here")
	}
	for _, ev := range saver.Read() {
		if ev.Name() != "tls_handshake_done" {
			continue
		}
		if ev.Value().NoTLSVerify == true {
			t.Fatal("expected NoTLSVerify to be false")
		}
		if len(ev.Value().TLSPeerCerts) < 1 {
			t.Fatal("expected at least a certificate here")
		}
	}
}

func TestSaverTLSHandshakerNoTLSVerify(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	saver := &Saver{}
	tlsdlr := &netxlite.TLSDialerLegacy{
		Config:        &tls.Config{InsecureSkipVerify: true},
		Dialer:        &netxlite.DialerSystem{},
		TLSHandshaker: saver.WrapTLSHandshaker(&netxlite.TLSHandshakerConfigurable{}),
	}
	conn, err := tlsdlr.DialTLSContext(
		context.Background(), "tcp", "self-signed.badssl.com:443")
	if err != nil {
		t.Fatal(err)
	}
	if conn == nil {
		t.Fatal("expected non-nil conn here")
	}
	conn.Close()
	for _, ev := range saver.Read() {
		if ev.Name() != "tls_handshake_done" {
			continue
		}
		if ev.Value().NoTLSVerify != true {
			t.Fatal("expected NoTLSVerify to be true")
		}
		if len(ev.Value().TLSPeerCerts) < 1 {
			t.Fatal("expected at least a certificate here")
		}
	}
}