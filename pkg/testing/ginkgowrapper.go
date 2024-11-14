/*
Copyright 2024 The KServe Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var errInterface = reflect.TypeOf((*error)(nil)).Elem()

// IgnoreNotFound can be used to wrap an arbitrary function in a call to
// [ginkgo.DeferCleanup]. When the wrapped function returns an error that
// `apierrors.IsNotFound` considers as "not found", the error is ignored
// instead of failing the test during cleanup. This is useful for cleanup code
// that just needs to ensure that some object does not exist anymore.
func IgnoreNotFound(in any) any {
	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)
	return reflect.MakeFunc(inType, func(args []reflect.Value) []reflect.Value {
		var out []reflect.Value
		if inType.IsVariadic() {
			out = inValue.CallSlice(args)
		} else {
			out = inValue.Call(args)
		}
		if len(out) > 0 {
			lastValue := out[len(out)-1]
			last := lastValue.Interface()
			if last != nil && lastValue.Type().Implements(errInterface) && apierrors.IsNotFound(last.(error)) {
				out[len(out)-1] = reflect.Zero(errInterface)
			}
		}
		return out
	}).Interface()
}