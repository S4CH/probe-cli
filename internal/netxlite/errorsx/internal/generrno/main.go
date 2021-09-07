package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iancoleman/strcase"
	"golang.org/x/sys/execabs"
)

// ErrorSpec specifies the error we care about.
type ErrorSpec struct {
	// errno is the error name as an errno value (e.g., ECONNREFUSED).
	errno string

	// failure is the error name according to OONI (e.g., FailureConnectionRefused).
	failure string
}

// AsErrnoName returns the name of the corresponding errno, if this
// is a system error, or panics otherwise.
func (es *ErrorSpec) AsErrnoName() string {
	if !es.IsSystemError() {
		panic("not a system error")
	}
	return es.errno
}

// AsFailureVar returns the name of the failure var.
func (es *ErrorSpec) AsFailureVar() string {
	return "Failure" + strcase.ToCamel(es.failure)
}

// AsFailureString returns the OONI failure string.
func (es *ErrorSpec) AsFailureString() string {
	return strcase.ToSnake(es.failure)
}

// NewSystemError constructs a new ErrorSpec representing a system
// error, i.e., an error returned by a system call.
func NewSystemError(errno, failure string) *ErrorSpec {
	return &ErrorSpec{errno: errno, failure: failure}
}

// NewLibraryError constructs a new ErrorSpec representing a library
// error, i.e., an error returned by the Go standard library or by other
// dependecies written typicall in Go (e.g., quic-go).
func NewLibraryError(failure string) *ErrorSpec {
	return &ErrorSpec{failure: failure}
}

// IsSystemError returns whether this ErrorSpec describes a system
// error, i.e., an error returned by a syscall.
func (es *ErrorSpec) IsSystemError() bool {
	return es.errno != ""
}

// Specs contains all the error specs.
var Specs = []*ErrorSpec{
	NewSystemError("ECANCELED", "operation_canceled"),
	NewSystemError("ECONNREFUSED", "connection_refused"),
	NewSystemError("ECONNRESET", "connection_reset"),
	NewSystemError("EHOSTUNREACH", "host_unreachable"),
	NewSystemError("ETIMEDOUT", "timed_out"),
	NewSystemError("EAFNOSUPPORT", "address_family_not_supported"),
	NewSystemError("EADDRINUSE", "address_in_use"),
	NewSystemError("EADDRNOTAVAIL", "address_not_available"),
	NewSystemError("EISCONN", "already_connected"),
	NewSystemError("EFAULT", "bad_address"),
	NewSystemError("EBADF", "bad_file_descriptor"),
	NewSystemError("ECONNABORTED", "connection_aborted"),
	NewSystemError("EALREADY", "connection_already_in_progress"),
	NewSystemError("EDESTADDRREQ", "destination_address_required"),
	NewSystemError("EINTR", "interrupted"),
	NewSystemError("EINVAL", "invalid_argument"),
	NewSystemError("EMSGSIZE", "message_size"),
	NewSystemError("ENETDOWN", "network_down"),
	NewSystemError("ENETRESET", "network_reset"),
	NewSystemError("ENETUNREACH", "network_unreachable"),
	NewSystemError("ENOBUFS", "no_buffer_space"),
	NewSystemError("ENOPROTOOPT", "no_protocol_option"),
	NewSystemError("ENOTSOCK", "not_a_socket"),
	NewSystemError("ENOTCONN", "not_connected"),
	NewSystemError("EWOULDBLOCK", "operation_would_block"),
	NewSystemError("EACCES", "permission_denied"),
	NewSystemError("EPROTONOSUPPORT", "protocol_not_supported"),
	NewSystemError("EPROTOTYPE", "wrong_protocol_type"),

	// Implementation note: we need to specify acronyms we
	// want to be upper case in uppercase here. For example,
	// we must write "DNS" rather than writing "dns".
	NewLibraryError("DNS_bogon_error"),
	NewLibraryError("DNS_NXDOMAIN_error"),
	NewLibraryError("EOF_error"),
	NewLibraryError("generic_timeout_error"),
	NewLibraryError("QUIC_incompatible_version"),
	NewLibraryError("SSL_failed_handshake"),
	NewLibraryError("SSL_invalid_hostname"),
	NewLibraryError("SSL_unknown_authority"),
	NewLibraryError("SSL_invalid_certificate"),
	NewLibraryError("JSON_parse_error"),
}

func fileCreate(filename string) *os.File {
	filep, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	return filep
}

func fileWrite(filep *os.File, content string) {
	if _, err := filep.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

func fileClose(filep *os.File) {
	if err := filep.Close(); err != nil {
		log.Fatal(err)
	}
}

func filePrintf(filep *os.File, format string, v ...interface{}) {
	fileWrite(filep, fmt.Sprintf(format, v...))
}

func gofmt(filename string) {
	cmd := execabs.Command("go", "fmt", filename)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func writeSystemSpecificFile(kind string) {
	filename := "errno_" + kind + ".go"
	filep := fileCreate(filename)
	fileWrite(filep, "// Code generated by go generate; DO NOT EDIT.\n")
	filePrintf(filep, "// Generated: %+v\n\n", time.Now())
	fileWrite(filep, "package errorsx\n\n")
	filePrintf(filep, "import \"golang.org/x/sys/%s\"\n\n", kind)
	fileWrite(filep, "const (\n")
	for _, spec := range Specs {
		if !spec.IsSystemError() {
			continue
		}
		filePrintf(filep, "\t%s = %s.%s\n",
			spec.AsErrnoName(), kind, spec.AsErrnoName())
	}
	fileWrite(filep, ")\n\n")
	fileClose(filep)
	gofmt(filename)
}

func writeGenericFile() {
	filename := "errno.go"
	filep := fileCreate(filename)
	fileWrite(filep, "// Code generated by go generate; DO NOT EDIT.\n")
	filePrintf(filep, "// Generated: %+v\n\n", time.Now())
	fileWrite(filep, "package errorsx\n\n")
	fileWrite(filep, "//go:generate go run ./internal/generrno/\n\n")
	fileWrite(filep, "import (\n")
	fileWrite(filep, "\t\"errors\"\n")
	fileWrite(filep, "\t\"syscall\"\n")
	fileWrite(filep, ")\n\n")

	fileWrite(filep, "// This enumeration lists the failures defined at\n")
	fileWrite(filep, "// https://github.com/ooni/spec/blob/master/data-formats/df-007-errors.md\n")
	fileWrite(filep, "const (\n")
	fileWrite(filep, "//\n")
	fileWrite(filep, "// System errors\n")
	fileWrite(filep, "//\n")
	for _, spec := range Specs {
		if !spec.IsSystemError() {
			continue
		}
		filePrintf(filep, "\t%s = \"%s\"\n",
			spec.AsFailureVar(),
			spec.AsFailureString())
	}
	fileWrite(filep, "\n")
	fileWrite(filep, "//\n")
	fileWrite(filep, "// Library errors\n")
	fileWrite(filep, "//\n")
	for _, spec := range Specs {
		if spec.IsSystemError() {
			continue
		}
		filePrintf(filep, "\t%s = \"%s\"\n",
			spec.AsFailureVar(),
			spec.AsFailureString())
	}
	fileWrite(filep, ")\n\n")

	fileWrite(filep, "// classifySyscallError converts a syscall error to the\n")
	fileWrite(filep, "// proper OONI error. Returns the OONI error string\n")
	fileWrite(filep, "// on success, an empty string otherwise.\n")
	fileWrite(filep, "func classifySyscallError(err error) string {\n")
	fileWrite(filep, "\t// filter out system errors: necessary to detect all windows errors\n")
	fileWrite(filep, "\t// https://github.com/ooni/probe/issues/1526 describes the problem\n")
	fileWrite(filep, "\t// of mapping localized windows errors.\n")
	fileWrite(filep, "\tvar errno syscall.Errno\n")
	fileWrite(filep, "\tif !errors.As(err, &errno) {\n")
	fileWrite(filep, "\t\treturn \"\"\n")
	fileWrite(filep, "\t}\n")
	fileWrite(filep, "\tswitch errno {\n")
	for _, spec := range Specs {
		if !spec.IsSystemError() {
			continue
		}
		filePrintf(filep, "\tcase %s:\n", spec.AsErrnoName())
		filePrintf(filep, "\t\treturn %s\n", spec.AsFailureVar())
	}
	fileWrite(filep, "\t}\n")
	fileWrite(filep, "\treturn \"\"\n")
	fileWrite(filep, "}\n\n")

	fileClose(filep)
	gofmt(filename)
}

func writeGenericTestFile() {
	filename := "errno_test.go"
	filep := fileCreate(filename)

	fileWrite(filep, "// Code generated by go generate; DO NOT EDIT.\n")
	filePrintf(filep, "// Generated: %+v\n\n", time.Now())
	fileWrite(filep, "package errorsx\n\n")
	fileWrite(filep, "import (\n")
	fileWrite(filep, "\t\"io\"\n")
	fileWrite(filep, "\t\"syscall\"\n")
	fileWrite(filep, "\t\"testing\"\n")
	fileWrite(filep, ")\n\n")

	fileWrite(filep, "func TestToSyscallErr(t *testing.T) {\n")
	fileWrite(filep, "\tif v := classifySyscallError(io.EOF); v != \"\" {\n")
	fileWrite(filep, "\t\tt.Fatalf(\"expected empty string, got '%s'\", v)\n")
	fileWrite(filep, "\t}\n")

	for _, spec := range Specs {
		if !spec.IsSystemError() {
			continue
		}
		filePrintf(filep, "\tif v := classifySyscallError(%s); v != %s {\n",
			spec.AsErrnoName(), spec.AsFailureVar())
		filePrintf(filep, "\t\tt.Fatalf(\"expected '%%s', got '%%s'\", %s, v)\n",
			spec.AsFailureVar())
		fileWrite(filep, "\t}\n")
	}

	fileWrite(filep, "\tif v := classifySyscallError(syscall.Errno(0)); v != \"\" {\n")
	fileWrite(filep, "\t\tt.Fatalf(\"expected empty string, got '%s'\", v)\n")
	fileWrite(filep, "\t}\n")
	fileWrite(filep, "}\n")

	fileClose(filep)
	gofmt(filename)
}

func main() {
	writeSystemSpecificFile("unix")
	writeSystemSpecificFile("windows")
	writeGenericFile()
	writeGenericTestFile()
}