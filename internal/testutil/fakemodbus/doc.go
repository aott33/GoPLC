// Package fakemodbus is a scriptable TCP server for testing Modbus clients
// against well-behaved and misbehaving devices.
//
// The server works only in raw bytes. It does not import the Modbus codec,
// so tests that use it are an independent check on the codec: the bytes a
// test expects and replies with are written out by hand.
//
// # Script format
//
// A script is a []Step run in order on one TCP connection. The server is
// started with one script per connection: the first accepted connection
// runs scripts[0], the second runs scripts[1], and so on. A connection
// accepted after the scripts run out is closed at once, which looks like a
// device that accepts and then drops the connection.
//
// Steps:
//
//   - Expect(b): read exactly len(b) bytes and compare them with b.
//   - ExpectAnyTID(b): like Expect, but bytes 0 and 1 (the MBAP transaction
//     ID) are not compared. The received transaction ID is remembered.
//   - Reply(b): write b as-is.
//   - ReplyEchoTID(b): write b with bytes 0 and 1 replaced by the transaction
//     ID remembered by the last ExpectAnyTID on this connection.
//   - Garbage(b): write b as-is. Same as Reply; use it to make intent clear.
//   - Delay(d): wait for d, or until the server is closed.
//   - Close(): close the connection and end the script.
//
// Repeat(n, steps...) returns steps repeated n times, for building scripts
// such as "answer three polls, then drop the connection":
//
//	script := append(
//		fakemodbus.Repeat(3,
//			fakemodbus.ExpectAnyTID(req),
//			fakemodbus.ReplyEchoTID(resp),
//		),
//		fakemodbus.Close(),
//	)
//
// A reply can be split across TCP writes with Reply, Delay, Reply.
//
// When a script finishes without Close, the connection stays open and
// anything the client sends is read and discarded without a reply, so the
// client sees a timeout. The connection ends when the client closes it or
// the server is closed.
//
// # Errors
//
// If an Expect step reads different bytes, or the connection ends before
// all expected bytes arrive, the server records an error, closes that
// connection and stops its script. Err returns the recorded errors. Write
// errors are not recorded: a client closing the connection early is normal
// in failure tests.
//
// Use Start in tests; it registers Close with t.Cleanup and fails the test
// if any expectation failed.
package fakemodbus
