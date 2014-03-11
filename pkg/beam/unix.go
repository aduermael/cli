package beam

import (
	"syscall"
	"fmt"
	"net"
	"os"
)

// Send sends a new message on conn with data and f as payload and
// attachment, respectively.
func Send(conn *net.UnixConn, data []byte, f *os.File) error {
	return sendUnix(conn, data, int(f.Fd()))
}

// Receive waits for a new message on conn, and receives its payload
// and attachment, or an error if any.
//
// If more than 1 file descriptor is sent in the message, they are all
// closed except for the first, which is the attachment.
// It is legal for a message to have no attachment or an empty payload.
func Receive(conn *net.UnixConn) ([]byte, *os.File, error) {
	for {
		data, fds, err := receiveUnix(conn)
		if err != nil {
			return nil, nil, fmt.Errorf("receive: %v", err)
		}
		if len(fds) == 0 {
			continue
		}
		if len(fds) > 1 {
			for _, fd := range fds {
				syscall.Close(fd)
			}
		}
		return data, os.NewFile(uintptr(fds[0]), ""), nil
	}
	panic("impossibru")
	return nil, nil, nil
}

// SendPipe creates a new unix socket pair, sends one end as the attachment
// to a beam message with the payload `data`, and returns the other end.
//
// This is a common pattern to open a new service endpoint.
// For example, a service wishing to advertise its presence to clients might
// open an endpoint with:
//
//  endpoint, _ := SendPipe(conn, []byte("sql"))
//  defer endpoint.Close()
//  for {
//  	conn, _ := endpoint.Receive()
//	go func() {
//		Handle(conn)
//		conn.Close()
//	}()
//  }
//
// Note that beam does not distinguish between clients and servers in the logical
// sense: any program wishing to establishing a communication with another program
// may use SendPipe() to create an endpoint.
// For example, here is how an application might use it to connect to a database client.
//
//  endpoint, _ := SendPipe(conn, []byte("userdb"))
//  defer endpoint.Close()
//  conn, _ := endpoint.Receive()
//  defer conn.Close()
//  db := NewDBClient(conn)
//
// In this example note that we only need the first connection out of the endpoint,
// but we could open new ones to retry after a broken connection.
// Note that, because the underlying service transport is abstracted away, this
// allows for arbitrarily complex service discovery and retry logic to take place,
// without complicating application code.
//
func SendPipe(conn *net.UnixConn, data []byte) (*os.File, error) {
	local, remote, err := SocketPair()
	if err != nil {
		return nil, err
	}
	if err := Send(conn, data, remote); err != nil {
		remote.Close()
		return nil, err
	}
	return local, nil
}

func receiveUnix(conn *net.UnixConn) ([]byte, []int, error) {
	buf := make([]byte, 4096)
	oob := make([]byte, 4096)
	bufn, oobn, _, _, err := conn.ReadMsgUnix(buf, oob)
	if err != nil {
		return nil, nil, err
	}
	return buf[:bufn], extractFds(oob[:oobn]), nil
}

func sendUnix(conn *net.UnixConn, data []byte, fds ...int) error {
	_, _, err := conn.WriteMsgUnix(data, syscall.UnixRights(fds...), nil)
	return err
}


func extractFds(oob []byte) (fds []int) {
	scms, err := syscall.ParseSocketControlMessage(oob)
	if err != nil {
		return
	}
	for _, scm := range scms {
		gotFds, err := syscall.ParseUnixRights(&scm)
		if err != nil {
			continue
		}
		fds = append(fds, gotFds...)
	}
	return
}

func socketpair() ([2]int, error) {
	return syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM|syscall.FD_CLOEXEC, 0)
}

// SocketPair is a convenience wrapper around the socketpair(2) syscall.
// It returns a unix socket of type SOCK_STREAM in the form of 2 file descriptors
// not bound to the underlying filesystem.
// Messages sent on one end are received on the other, and vice-versa.
// It is the caller's responsibility to close both ends.
func SocketPair() (*os.File, *os.File, error) {
	pair, err := socketpair()
	if err != nil {
		return nil, nil, err
	}
	return os.NewFile(uintptr(pair[0]), ""), os.NewFile(uintptr(pair[1]), ""), nil
}

// FdConn wraps a file descriptor in a standard *net.UnixConn object, or
// returns an error if the file descriptor does not point to a unix socket.
func FdConn(fd int) (*net.UnixConn, error) {
	f := os.NewFile(uintptr(fd), fmt.Sprintf("%d", fd))
	conn, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}
	uconn, ok := conn.(*net.UnixConn)
	if !ok {
		return nil, fmt.Errorf("%d: not a unix connection", fd)
	}
	return uconn, nil
}
