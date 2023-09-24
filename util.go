package gocron

import (
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

func strPtr(in string) *string {
	return &in
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return
	case <-time.After(timeout):
		return
	}
}

//func callJobFunc(jobFunc interface{}) {
//	if jobFunc == nil {
//		return
//	}
//	f := reflect.ValueOf(jobFunc)
//	if !f.IsZero() {
//		f.Call([]reflect.Value{})
//	}
//}

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

func requestJob(id uuid.UUID, ch chan jobOutRequest) job {
	resp := make(chan job)
	ch <- jobOutRequest{
		id:      id,
		outChan: resp,
	}
	var j job
	for jobReceived := range resp {
		j = jobReceived
	}
	return j
}
