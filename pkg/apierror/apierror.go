package apierror

import (
	"fmt"
	"regexp"

	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	// "github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"

	"github.com/golang/protobuf/proto" // nolint
)

// compiler time checks
var (
	_ error                                        = (*Status)(nil)
	_ interface{ GRPCStatus() *grpcstatus.Status } = (*Status)(nil)
)

const (
	UpdateDataNotfound  = "No data found to be updated"
	DataAlreadyExists   = "Already exits"
	PermissionDeniedMsg = "Permission denied"
)

// New API error.
func New(code Code, msg string) (st *Status) {
	st = new(Status)
	st.Code = code
	st.Message = msg
	return
}

// Newf uses formated message.
func Newf(code Code, format string, args ...interface{}) (st *Status) {
	return New(code, fmt.Sprintf(format, args...))
}

// apierror object from the rpc error message
func NewFromErrorMsg(message string) (st *Status) {
	pattern := `code = (\w+) desc = (.+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		msg := fmt.Sprintf("Invalid rpc response: %s", message)
		return New(Code(Code_value["InvalidArgument"]), msg)
	}

	errObj := Newf(Code(Code_value[matches[1]]), "%s", matches[2])

	return errObj
}

// Error implements standard error interface.
func (s *Status) Error() string {
	return fmt.Sprintf("%s: %s", grpccodes.Code(s.Code), s.Message)
}

// GRPCStatus used for gRPC error response.
func (s *Status) GRPCStatus() (st *grpcstatus.Status) {
	st = grpcstatus.New(grpccodes.Code(s.Code), s.Message)
	if len(s.Details) == 0 {
		return st // no details, done
	}
	var details = make([]proto.Message, 0, len(s.Details))
	for _, any := range s.Details {
		var to, anyerr = types.EmptyAny(any)
		if anyerr != nil {
			continue // silently drop
		}
		if anyerr = types.UnmarshalAny(any, to); anyerr != nil {
			continue // silently drop
		}
		details = append(details, to)
	}
	var ds, detailserr = st.WithDetails(details...)
	if detailserr != nil {
		return // return it without details (silently drop them)
	}
	return ds // with details
}

// WithDetails returns the Status with given details.
func (s *Status) WithDetails(details ...proto.Message) (ds *Status, err error) {
	ds = New(s.Code, s.Message)
	if len(details) > 0 && s.Code == OK {
		return nil, fmt.Errorf("can't add details to successful Status")
	}
	ds.Details = make([]*types.Any, 0, len(details))
	for _, det := range details {
		var any *types.Any
		if any, err = types.MarshalAny(det); err != nil {
			return
		}
		ds.Details = append(ds.Details, any)
	}
	return
}

// AddDetails works exactly like the WithDetails. But it returns
// the Status itself with given details. And if a detail can't be
// added the error just ignored.
func (s *Status) AddDetails(details ...proto.Message) *Status {
	if len(details) > 0 && s.Code == OK {
		return s // don't modify
	}
	var err error
	for _, det := range details {
		var any *types.Any
		if any, err = types.MarshalAny(det); err != nil {
			continue // just ignore error, try next
		}
		s.Details = append(s.Details, any)
	}
	return s // may be with details
}

// Useless, errors.As resolves it via reflection.
//
// // As used for errors.As(err, &status) purpose.
// func (s *Status) As(err interface{}) (ok bool) {
// 	var x *Status
// 	if x, ok = err.(*Status); ok {
// 		(*s) = (*x)
// 	}
// 	return
// }
