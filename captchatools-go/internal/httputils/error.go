package httputils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"syscall"
)

/*
Possible errors returned from `.Do()`:

	*url.Error:
		*net.OpError:
			*net.AddrError:
			*net.DNSError:
				// do something more here
			*net.InvalidAddrError:
			*net.ParseError:
			*net.UnknownNetworkError:
			*os.SyscallError:
				// syscall errors

source: https://haisum.github.io/2021/08/14/2021-golang-http-errors/golang
*/

// Will be used to represent an HTTP client related error
type ClientErr struct {
	TimedOut bool // If the connection timedout
}

func (c ClientErr) Error() string {
	var message string
	if c.TimedOut {
		message = "timed out"
	} else {
		message = "unknown http client error"
	}
	return message
}

// Will be used to represent an proxy related error
type ProxyErr struct {
	InvalidAuth bool // Invalid username or password
	InvalidHost bool // Invalid proxy host name
	InvalidPort bool // Invalid proxy port
}

func (p ProxyErr) Error() string {
	var message string
	if p.InvalidAuth {
		message = "Proxy auth is incorrect"
	} else if p.InvalidHost {
		message = "Proxy host is incorrect"
	} else if p.InvalidPort {
		message = "Proxy port is incorrect"
	} else {
		message = "Proxy error"
	}
	return message
}

type RequestErr struct {
	Message string
}

func (r RequestErr) Error() string { return r.Message }

// Parses an error that was returned from the `.Do()`
// method and returns a ProxyErr or ClientErr
func ParseError(err error) error {

	var addrError *net.AddrError
	if errors.As(err, &addrError) && addrError.Err == "invalid port" {
		return &ProxyErr{InvalidPort: true}
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsNotFound {
			return &ProxyErr{InvalidHost: true} // host could not be found (usually proxy error)
		}
	}

	// Check if it's connection refused
	//
	// It is usually a client-side issue with a number of possible causes,
	// including an unreliable internet connection, Chrome extension issues,
	// antivirus and firewall interference, and incorrect internet settings
	if errors.Is(err, syscall.ECONNREFUSED) {
		return &ProxyErr{}
	}

	// According to Go's docs, any error returned by the
	// .Do method will always be of type url.Error
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		// Check for Client related errors
		if urlErr.Timeout() { // This also triggers when context is done so we have to becareful
			return &ClientErr{TimedOut: true}
		}
		if urlErr.Err.Error() == "EOF" {
			return &RequestErr{Message: "EOF"}
		}

		// Check for Proxy related errors
		errStr := urlErr.Unwrap().Error()
		switch errStr {
		case "Proxy responded with non 200 code: 407 Proxy Authentication Required":
			return &ProxyErr{InvalidAuth: true}
		case "Proxy responded with non 200 code: 502 Proxy Error (The selected relay is offline or busy processing other threads)":
			return &ProxyErr{}
		case "invalid URL scheme: []":
			return &ClientErr{}
		}
	}

	return err // if everything fails, return a client error
}

// Returns if the error was a proxy related error.
func IsProxyError(err error) bool {
	err = ParseError(err)
	_, ok := err.(*ProxyErr)
	_, ok2 := errors.Unwrap(err).(*ProxyErr)
	// fmt.Printf("ok = %t || ok2 = %t\n", ok, ok2)
	return ok || ok2
}

// Returns if the error was a http client related error
func IsClientError(err error) bool {
	err = ParseError(err)
	_, ok := err.(*ClientErr)
	_, ok2 := errors.Unwrap(err).(*ClientErr)
	return ok || ok2
}

// NetworkError is a custom type to categorize network-related errors.
type NetworkError string

const (
	// DNSNameNotFound indicates the requested name does not contain any
	// records of the requested type (data not found), or the name
	// itself was not found (NXDOMAIN)
	DNSNameNotFound NetworkError = "DNS name not found"

	// ConnectionRefused indicates that the connection was refused by the target machine
	ConnectionRefused NetworkError = "connection refused"

	// ConnectionTimedOut indicates that the connection has timedout
	ConnectionTimedOut NetworkError = "connection timed out"

	// HostUnreachable indicates that the target host was unreachable
	HostUnreachable NetworkError = "host unreachable"

	// NetworkUnreachable indicates that the network was unreachable
	NetworkUnreachable NetworkError = "network unreachable"

	// EOF indicates that connection was most likely closed before or while the headers are read.
	EOF NetworkError = "EOF error"

	// ProxyAuthRequired indicates that the request did not succeed because it
	// lacks valid authentication credentials for the proxy server that sits
	// between the client and the server with access to the requested resource.
	ProxyAuthRequired NetworkError = "proxy error - proxy authentication required"

	// ProxyRelayOffline indicates that the proxy responded with non 200 code:
	// 502 Proxy Error (The selected relay is offline or busy processing other threads)
	ProxyRelayOffline NetworkError = "proxy error - relay is offline or busy processing other threads"

	// BadURL indicates a URL error. This can be for various reasons such as a
	// missing scheme, blank URI, etc.
	BadURL NetworkError = "bad URL"

	// GenericError indicates a generic network error. This can be for various
	// reasons, most commonly that the error has not been added for validation.
	GenericError NetworkError = "network error"
)

// ClassifyNetworkError attempts to classify the provided error into one of the predefined NetworkError types.
func ClassifyNetworkError(err error) NetworkError {
	// Unwrap the error to find the root cause
	cause := errors.Unwrap(err)

	// Check for DNS errors, particularly "no such host"
	var dnsError *net.DNSError
	if errors.As(cause, &dnsError) {
		if dnsError.IsNotFound {
			return DNSNameNotFound
		}
		return NetworkError(fmt.Sprintf("unknown DNS error: %s", dnsError.Err))
	}

	// Handle specific syscall errors, such as connection refused
	var sysCallErr *syscall.Errno
	if errors.As(cause, &sysCallErr) {
		switch *sysCallErr {
		case syscall.ECONNREFUSED, 10061: // 10061 is specific to Windows
			return ConnectionRefused
		case syscall.ETIMEDOUT:
			return ConnectionTimedOut
		case syscall.EHOSTUNREACH:
			return HostUnreachable
		case syscall.ENETUNREACH:
			return NetworkUnreachable
		}
	}

	// According to Go's docs, any error returned by the
	// .Do method will always be of type url.Error
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		// Check for Client related errors
		if urlErr.Timeout() { // This also triggers when context is done so we have to becareful
			return ConnectionTimedOut
		}
		if urlErr.Err.Error() == "EOF" {
			return EOF
		}

		// Check for Proxy related errors
		errStr := urlErr.Unwrap().Error()
		switch errStr {
		case "Proxy responded with non 200 code: 407 Proxy Authentication Required":
			return ProxyAuthRequired
		case "Proxy responded with non 200 code: 502 Proxy Error (The selected relay is offline or busy processing other threads)":
			return ProxyRelayOffline
		case "invalid URL scheme: []":
			return BadURL
		}
	}

	// Handle generic network errors with a timeout condition
	var genericErr net.Error
	if errors.As(cause, &genericErr) {
		if genericErr.Timeout() {
			return ConnectionTimedOut
		}
		return GenericError
	}

	// Handle other cases as unknown errors
	return NetworkError(fmt.Sprintf("unknown network error: %s", err))
}
