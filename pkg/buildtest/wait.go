/*
Copyright 2018 Google, Inc. All rights reserved.

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

package buildtest

import (
	"sync"
	"time"
)

// This extends WaitGroup with a few useful methods to avoid tests having a failure mode of hanging.
type Wait struct {
	sync.WaitGroup
}

// Override Add so folks don't call it.
func (w *Wait) Add(delta int) {
	panic("Don't call Add(int) on wait!")
}

// This calls Done() after the specified duration.
func (w *Wait) In(d time.Duration) {
	go func() {
		time.Sleep(d)
		w.Done()
	}()
}

// A convenience for passing into WaitUntil
func WaitNop() {}

// This waits until Done() has been called or the specified duration has elapsed.
// If Done() is called, then onSuccess is called.
// If the duration elapses without Done() being called, then onTimeout is called.
func (w *Wait) WaitUntil(d time.Duration, onSuccess func(), onTimeout func()) {
	ch := make(chan struct{})
	go func() {
		w.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		onSuccess()
	case <-time.After(d):
		onTimeout()
	}
}

// A convenience for creating a "wait" object.
func NewWait() *Wait {
	w := Wait{}
	w.WaitGroup.Add(1)
	return &w
}
