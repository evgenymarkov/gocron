package gocron

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

func callJobFuncWithParams(jobFunc interface{}, params ...interface{}) error {
	if jobFunc == nil {
		return nil
	}
	f := reflect.ValueOf(jobFunc)
	if f.IsZero() {
		return nil
	}
	if len(params) != f.Type().NumIn() {
		return nil
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	vals := f.Call(in)
	for _, val := range vals {
		i := val.Interface()
		if err, ok := i.(error); ok {
			return err
		}
	}
	return nil
}

func requestJob(id uuid.UUID, ch chan jobOutRequest, ctx context.Context) *internalJob {
	resp := make(chan internalJob, 1)
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	ch <- jobOutRequest{
		id:      id,
		outChan: resp,
	}
	var j internalJob
	select {
	case <-ctx.Done():
		return nil
	case jobReceived := <-resp:
		j = jobReceived
	}
	return &j
}

func removeSliceDuplicatesInt(in []int) []int {
	m := make(map[int]struct{})

	for _, i := range in {
		m[i] = struct{}{}
	}
	return maps.Keys(m)
}