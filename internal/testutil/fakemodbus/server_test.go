package fakemodbus

import (
	"bytes"
	"errors"
	"go/parser"
	"go/token"
	"io"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// FC03 request: TID 1, unit 1, start 0, qty 2, and a response holding 0x1234, 0x5678.
var (
	req  = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02}
	resp = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x07, 0x01, 0x03, 0x04, 0x12, 0x34, 0x56, 0x78}
)

func dial(t *testing.T, s *Server) net.Conn {
	t.Helper()
	c, err := net.DialTimeout("tcp", s.Addr(), time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	return c
}

func write(t *testing.T, c net.Conn, b []byte) {
	t.Helper()
	if _, err := c.Write(b); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readN(t *testing.T, c net.Conn, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := io.ReadFull(c, b); err != nil {
		t.Fatalf("read %d bytes: %v", n, err)
	}
	return b
}

// expectEOF waits for the server to close the connection.
func expectEOF(t *testing.T, c net.Conn) {
	t.Helper()
	var b [1]byte
	n, err := c.Read(b[:])
	if !errors.Is(err, io.EOF) {
		t.Fatalf("read: got n=%d err=%v, want EOF", n, err)
	}
}

func TestExpectReply(t *testing.T) {
	s := Start(t, []Step{Expect(req), Reply(resp)})
	c := dial(t, s)

	write(t, c, req)
	if got := readN(t, c, len(resp)); !bytes.Equal(got, resp) {
		t.Errorf("got % X, want % X", got, resp)
	}
}

func TestExpectMismatchRecorded(t *testing.T) {
	s, err := New([]Step{Expect(req), Reply(resp)})
	if err != nil {
		t.Fatal(err)
	}
	c := dial(t, s)

	bad := bytes.Clone(req)
	bad[11] = 0x04
	write(t, c, bad)
	expectEOF(t, c) // the server drops the connection on a mismatch

	s.Close()
	err = s.Err()
	if err == nil || !strings.Contains(err.Error(), "step 0") {
		t.Fatalf("Err() = %v, want a step 0 mismatch", err)
	}
}

func TestShortRequestRecorded(t *testing.T) {
	s, err := New([]Step{Expect(req)})
	if err != nil {
		t.Fatal(err)
	}
	c := dial(t, s)
	write(t, c, req[:5])
	_ = c.Close()

	// Wait for the server to notice the closed connection.
	deadline := time.Now().Add(2 * time.Second)
	for s.Err() == nil && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	s.Close()
	if err := s.Err(); err == nil || !strings.Contains(err.Error(), "got 5 bytes") {
		t.Fatalf("Err() = %v, want a short read", err)
	}
}

func TestEchoTID(t *testing.T) {
	s := Start(t, Repeat(2, ExpectAnyTID(req), ReplyEchoTID(resp)))
	c := dial(t, s)

	for _, tid := range []uint16{0x1234, 0xBEEF} {
		r := bytes.Clone(req)
		r[0], r[1] = byte(tid>>8), byte(tid)
		write(t, c, r)

		got := readN(t, c, len(resp))
		if got[0] != byte(tid>>8) || got[1] != byte(tid) {
			t.Errorf("TID = %02X%02X, want %04X", got[0], got[1], tid)
		}
		if !bytes.Equal(got[2:], resp[2:]) {
			t.Errorf("body = % X, want % X", got[2:], resp[2:])
		}
	}
}

func TestReplyEchoTIDWithoutCapture(t *testing.T) {
	s := Start(t, []Step{ReplyEchoTID(resp)})
	c := dial(t, s)
	if got := readN(t, c, len(resp)); !bytes.Equal(got, resp) {
		t.Errorf("got % X, want % X", got, resp)
	}
}

func TestDelayAndSplitReply(t *testing.T) {
	const d = 100 * time.Millisecond
	s := Start(t, []Step{Expect(req), Reply(resp[:4]), Delay(d), Reply(resp[4:])})
	c := dial(t, s)

	write(t, c, req)
	start := time.Now()
	got := readN(t, c, len(resp))
	if elapsed := time.Since(start); elapsed < d {
		t.Errorf("reply took %v, want at least %v", elapsed, d)
	}
	if !bytes.Equal(got, resp) {
		t.Errorf("got % X, want % X", got, resp)
	}
}

func TestGarbage(t *testing.T) {
	junk := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	s := Start(t, []Step{Expect(req), Garbage(junk), Close()})
	c := dial(t, s)

	write(t, c, req)
	if got := readN(t, c, len(junk)); !bytes.Equal(got, junk) {
		t.Errorf("got % X, want % X", got, junk)
	}
	expectEOF(t, c)
}

func TestScriptPerConnection(t *testing.T) {
	s := Start(t,
		[]Step{Close()},                  // first connection: dropped at once
		[]Step{Expect(req), Reply(resp)}, // second connection: answers
	)

	c1 := dial(t, s)
	expectEOF(t, c1)

	c2 := dial(t, s)
	write(t, c2, req)
	if got := readN(t, c2, len(resp)); !bytes.Equal(got, resp) {
		t.Errorf("got % X, want % X", got, resp)
	}

	// Scripts have run out: the third connection is closed at once.
	c3 := dial(t, s)
	expectEOF(t, c3)

	if got := s.Accepted(); got != 3 {
		t.Errorf("Accepted() = %d, want 3", got)
	}
}

func TestHoldAfterScript(t *testing.T) {
	s := Start(t, []Step{Expect(req), Reply(resp)})
	c := dial(t, s)

	write(t, c, req)
	readN(t, c, len(resp))

	// A second request gets no reply.
	write(t, c, req)
	_ = c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	var b [1]byte
	_, err := c.Read(b[:])
	var ne net.Error
	if !errors.As(err, &ne) || !ne.Timeout() {
		t.Fatalf("read after script: err = %v, want timeout", err)
	}
}

func TestCloseInterruptsDelay(t *testing.T) {
	s, err := New([]Step{Delay(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	dial(t, s)

	// Wait until the connection is being served.
	for s.Accepted() == 0 {
		time.Sleep(time.Millisecond)
	}

	done := make(chan struct{})
	go func() {
		s.Close()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return while a Delay step was running")
	}
	if err := s.Err(); err != nil {
		t.Errorf("Err() = %v, want nil", err)
	}
}

func TestCloseInterruptsExpect(t *testing.T) {
	s, err := New([]Step{Expect(req)})
	if err != nil {
		t.Fatal(err)
	}
	dial(t, s)
	for s.Accepted() == 0 {
		time.Sleep(time.Millisecond)
	}
	s.Close()
	s.Close() // second call is a no-op
	if err := s.Err(); err != nil {
		t.Errorf("Err() = %v, want nil when the server is closed during Expect", err)
	}
}

func TestRepeat(t *testing.T) {
	steps := Repeat(3, Expect(req), Reply(resp))
	if len(steps) != 6 {
		t.Fatalf("len = %d, want 6", len(steps))
	}
	if got := Repeat(0, Close()); len(got) != 0 {
		t.Errorf("Repeat(0) len = %d, want 0", len(got))
	}
}

// TestNoCodecImport keeps this package independent of the Modbus codec.
func TestNoCodecImport(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, name := range files {
		f, err := parser.ParseFile(fset, name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imp := range f.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			if strings.HasSuffix(path, "/internal/modbus") || strings.Contains(path, "/internal/modbus/") {
				t.Errorf("%s imports %s; fakemodbus must not depend on the codec", name, path)
			}
		}
	}
}
