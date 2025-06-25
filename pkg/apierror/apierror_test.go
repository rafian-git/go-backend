package apierror

import (
	"errors"
	"testing"

	// "github.com/gogo/protobuf/types"
	// "github.com/golang/protobuf/proto" //nolint

	// grpcstatus "google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var apierr = New(OK, "pew")
	require.NotNil(t, apierr)
	assert.Equal(t, OK, apierr.Code)
	assert.Equal(t, "pew", apierr.Message)
}

func TestNewf(t *testing.T) {
	var apierr = Newf(OK, "pew-%d", 10)
	require.NotNil(t, apierr)
	assert.Equal(t, OK, apierr.Code)
	assert.Equal(t, "pew-10", apierr.Message)
}

func TestStatus_Error(t *testing.T) {
	var apierr = New(OK, "pew")
	assert.Equal(t, OK.String()+": pew", apierr.Error())
}

func TestStatus_WithDetails(t *testing.T) {
	var apierr = New(Unauthenticated, "куда прёшь!")
	var axx, err = apierr.WithDetails(&DebugInfo{
		Detail: "nothing",
	})
	assert.NoError(t, err)
	assert.NotNil(t, axx)

	apierr = New(OK, "welcome")
	_, err = apierr.WithDetails(&DebugInfo{Detail: "nothing"})
	assert.Error(t, err)
}

func TestStatus_AddDetails(t *testing.T) {
	//
}

func TestStatus_As(t *testing.T) {
	var (
		apierr = New(PermissionDenied, "denied")
		stat   *Status
	)
	require.True(t, errors.As(apierr, &stat))
}

func TestStatus_GRPCStatus(t *testing.T) {
	var (
		apierr = New(Unauthenticated, "pew").
			AddDetails(&DebugInfo{Detail: "debug-pew"})
		status = apierr.GRPCStatus()
	)
	assert.EqualValues(t, apierr.Code, status.Code())
	assert.EqualValues(t, apierr.Message, status.Message())
	assert.Len(t, status.Details(), 1)
}
