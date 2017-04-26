// Copyright 2017 the gousb Authors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package usb

func (e *endpoint) newStream(size, count uint, submit bool) (*stream, error) {
	var ts []transferIntf
	for i := uint(0); i < count; i++ {
		t, err := newUSBTransfer(e.h, &e.Info, make([]byte, size), e.timeout)
		if err != nil {
			for _, t := range ts {
				t.free()
			}
			return nil, err
		}
		ts = append(ts, t)
	}
	return newStream(ts, submit), nil
}

// NewStream prepares a new read stream that will keep reading data from the
// endpoint until closed.
// Size defines a buffer size for a single read transaction and count
// defines how many transactions should be active at any time.
// By keeping multiple transfers active at the same time, a Stream reduces
// the latency between subsequent transfers and increases reading throughput.
func (e *InEndpoint) NewStream(size, count uint) (ReadStream, error) {
	s, err := e.newStream(size, count, true)
	if err != nil {
		return ReadStream{}, err
	}
	return ReadStream{s}, nil
}