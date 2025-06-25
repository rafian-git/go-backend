package utility

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.techetronventures.com/core/backend/pkg/apierror"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type RestApiPayload struct {
	Endpoint    string
	RequestBody interface{}
	QueryParams url.Values
	Dest        interface{}
	DestErr     interface{}
	Method      string
	Logger      *log.Logger
	Headers     map[string]string
}

func GetResponseFromRestApi(ctx context.Context, payload *RestApiPayload) error {
	var r *http.Request
	var err error
	if payload.RequestBody != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(payload.RequestBody)
		if err != nil {
			msg := "parsing error"
			payload.Logger.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
			return apierror.New(apierror.Internal, "Something went wrong")
		}
		r, err = http.NewRequest(payload.Method, payload.Endpoint, &buf)
	} else {
		r, err = http.NewRequest(payload.Method, payload.Endpoint, nil)
	}

	if err != nil {
		msg := "parsing error"
		payload.Logger.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return apierror.New(apierror.Internal, "Something went wrong")
	}

	if payload.QueryParams != nil {
		r.URL.RawQuery = payload.QueryParams.Encode()
	}

	r.Header.Add("Content-Type", "application/json")
	for key, val := range payload.Headers {
		r.Header.Add(key, val)
	}
	//r.Header.Add("Authorization", "Bearer "+utils.EncodeSHA256(p.cnf.AppToken))
	payload.Logger.Info(ctx, fmt.Sprintf("Request URL: %v\n", r.URL))
	payload.Logger.Info(ctx, fmt.Sprintf("Request Body: %v\n", r.Body))

	cli := http.Client{}
	resp, err := cli.Do(r)
	if err != nil {
		msg := "parsing error"
		payload.Logger.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return apierror.New(apierror.Internal, "Something went wrong")
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		msg := "Got status code.."
		payload.Logger.Error(ctx, fmt.Sprintf("%s : %d", msg, resp.StatusCode))
		err = json.NewDecoder(resp.Body).Decode(payload.DestErr)

		payload.Logger.Error(ctx, fmt.Sprintf("API Error: %v", payload.DestErr))
		if err != nil {
			msg := "parsing error"
			payload.Logger.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
			return apierror.New(apierror.Internal, "Something went wrong")
		}

		return apierror.New(apierror.Internal, "Something went wrong")
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.NewDecoder(resp.Body).Decode(payload.Dest)
	if err != nil {
		var builder strings.Builder
		builder.Write(bodyBytes)
		builder.WriteString("\n header \n")
		for key, value := range resp.Header {
			builder.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}

		payload.Logger.Error(ctx, "unavailable server "+payload.Endpoint, zap.String("content", builder.String()))
		return apierror.New(apierror.Internal, "Something went wrong")
	}

	return nil
}

func AddUrlValues(v url.Values, data interface{}) error {
	values := reflect.ValueOf(data)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct or a pointer to struct")
	}

	typ := values.Type()
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("url")

		if tag != "" {
			// Check if the field is zero for its type
			if isZero(field) {
				continue // Skip fields with zero values
			}
			v.Add(tag, fmt.Sprintf("%v", field.Interface()))
		}
	}

	return nil
}

// isZero checks if a value is zero for its type
func isZero(value reflect.Value) bool {
	zeroValue := reflect.Zero(value.Type())
	return reflect.DeepEqual(value.Interface(), zeroValue.Interface())
}
