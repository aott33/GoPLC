package fakemodbus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

type stepKind int

const (
	stepExpect stepKind = iota
	stepExpectAnyTID
	stepWrite
	stepWriteEchoTID
	stepDelay
	stepClose
)

// Step is one action in a connection script. Build steps with Expect,
// ExpectAnyTID, Reply, ReplyEchoTID, Garbage, Delay and Close.
type Step struct {
	kind  stepKind
	data  []byte
	delay time.Duration
}

// String describes the step for error messages.
func (s Step) String() string {
	switch s.kind {
	case stepExpect:
		return fmt.Sprintf("Expect(% X)", s.data)
	case stepExpectAnyTID:
		return fmt.Sprintf("ExpectAnyTID(% X)", s.data)
	case stepWrite:
		return fmt.Sprintf("Reply(% X)", s.data)
	case stepWriteEchoTID:
		return fmt.Sprintf("ReplyEchoTID(% X)", s.data)
	case stepDelay:
		return fmt.Sprintf("Delay(%v)", s.delay)
	case stepClose:
		return "Close()"
	default:
		return "unknown step"
	}
}

// Expect reads exactly len(b) bytes and records an error if they differ from b.
func Expect(b []byte) Step {
	return Step{kind: stepExpect, data: bytes.Clone(b)}
}

// ExpectAnyTID is like Expect but does not compare bytes 0 and 1, and
// remembers them for ReplyEchoTID. b must be at least 2 bytes long.
func ExpectAnyTID(b []byte) Step {
	if len(b) < 2 {
		panic("fakemodbus: ExpectAnyTID needs at least 2 bytes")
	}
	return Step{kind: stepExpectAnyTID, data: bytes.Clone(b)}
}

// Reply writes b to the client.
func Reply(b []byte) Step {
	return Step{kind: stepWrite, data: bytes.Clone(b)}
}

// ReplyEchoTID writes b with bytes 0 and 1 replaced by the transaction ID
// captured by the last ExpectAnyTID on this connection. If none was
// captured, b is written unchanged. b must be at least 2 bytes long.
func ReplyEchoTID(b []byte) Step {
	if len(b) < 2 {
		panic("fakemodbus: ReplyEchoTID needs at least 2 bytes")
	}
	return Step{kind: stepWriteEchoTID, data: bytes.Clone(b)}
}

// Garbage writes b to the client. It behaves exactly like Reply.
func Garbage(b []byte) Step {
	return Step{kind: stepWrite, data: bytes.Clone(b)}
}

// Delay waits for d, or until the server is closed.
func Delay(d time.Duration) Step {
	return Step{kind: stepDelay, delay: d}
}

// Close closes the connection and ends the script.
func Close() Step {
	return Step{kind: stepClose}
}

// Repeat returns steps repeated n times.
func Repeat(n int, steps ...Step) []Step {
	out := make([]Step, 0, n*len(steps))
	for range n {
		out = append(out, steps...)
	}
	return out
}

// Server is a scripted TCP server listening on 127.0.0.1.
type Server struct {
	ln      net.Listener
	scripts [][]Step
	done    chan struct{}
	wg      sync.WaitGroup

	closeOnce sync.Once

	mu       sync.Mutex
	conns    map[net.Conn]struct{}
	accepted int
	errs     []error
}

// New starts a server on 127.0.0.1 with a random port. Connection i runs
// scripts[i]. Call Close when done.
func New(scripts ...[]Step) (*Server, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("fakemodbus: listen: %w", err)
	}

	s := &Server{
		ln:      ln,
		scripts: scripts,
		done:    make(chan struct{}),
		conns:   make(map[net.Conn]struct{}),
	}

	s.wg.Add(1)
	go s.acceptLoop()

	return s, nil
}

// Start is New for tests. It fails the test if the server cannot start,
// closes the server in t.Cleanup and then reports any recorded errors.
func Start(t testing.TB, scripts ...[]Step) *Server {
	t.Helper()

	s, err := New(scripts...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		s.Close()
		if err := s.Err(); err != nil {
			t.Errorf("fakemodbus: %v", err)
		}
	})

	return s
}

// Addr returns the server address as host:port.
func (s *Server) Addr() string {
	return s.ln.Addr().String()
}

// Accepted returns the number of connections accepted so far.
func (s *Server) Accepted() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.accepted
}

// Err returns the expectation failures recorded so far, joined, or nil.
func (s *Server) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return errors.Join(s.errs...)
}

// Close stops the listener, closes open connections and waits for all
// connection goroutines to finish. It is safe to call more than once.
func (s *Server) Close() {
	s.closeOnce.Do(func() {
		close(s.done)
		_ = s.ln.Close()

		s.mu.Lock()
		for c := range s.conns {
			_ = c.Close()
		}
		s.mu.Unlock()
	})
	s.wg.Wait()
}

func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		c, err := s.ln.Accept()
		if err != nil {
			return // listener closed
		}

		s.mu.Lock()
		select {
		case <-s.done:
			s.mu.Unlock()
			_ = c.Close()
			return
		default:
		}
		idx := s.accepted
		s.accepted++
		s.conns[c] = struct{}{}
		s.wg.Add(1)
		s.mu.Unlock()

		go s.serve(idx, c)
	}
}

func (s *Server) recordf(format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errs = append(s.errs, fmt.Errorf(format, args...))
}

func (s *Server) serve(idx int, c net.Conn) {
	defer s.wg.Done()
	defer func() {
		s.mu.Lock()
		delete(s.conns, c)
		s.mu.Unlock()
		_ = c.Close()
	}()

	if idx >= len(s.scripts) {
		return
	}

	var tid [2]byte
	haveTID := false

	for i, step := range s.scripts[idx] {
		switch step.kind {
		case stepExpect, stepExpectAnyTID:
			got := make([]byte, len(step.data))
			n, err := io.ReadFull(c, got)
			if err != nil {
				if s.closing() {
					return
				}
				s.recordf("conn %d step %d %v: got %d bytes (% X) before %v", idx, i, step, n, got[:n], err)
				return
			}
			from := 0
			if step.kind == stepExpectAnyTID {
				from = 2
				copy(tid[:], got[:2])
				haveTID = true
			}
			if !bytes.Equal(got[from:], step.data[from:]) {
				s.recordf("conn %d step %d %v: got % X", idx, i, step, got)
				return
			}

		case stepWrite, stepWriteEchoTID:
			out := step.data
			if step.kind == stepWriteEchoTID && haveTID {
				out = bytes.Clone(step.data)
				copy(out[:2], tid[:])
			}
			if _, err := c.Write(out); err != nil {
				return // client gone; not an expectation failure
			}

		case stepDelay:
			t := time.NewTimer(step.delay)
			select {
			case <-t.C:
			case <-s.done:
				t.Stop()
				return
			}

		case stepClose:
			return
		}
	}

	// Script finished: hold the connection and ignore anything sent.
	_, _ = io.Copy(io.Discard, c)
}

func (s *Server) closing() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}
