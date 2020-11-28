// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package m3tsz

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"testing"

	"github.com/m3db/m3/src/dbnode/encoding"

	"github.com/stretchr/testify/require"
)

// BenchmarkM3TSZDecode-12    	   10000	    108797 ns/op
func BenchmarkM3TSZDecode(b *testing.B) {
	var (
		encodingOpts = encoding.NewOptions()
		reader       = bytes.NewReader(nil)
		seriesRun    = prepareSampleSeriesRun(b)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Reset(seriesRun[i])
		iter := NewReaderIterator(reader, DefaultIntOptimizationEnabled, encodingOpts)
		for iter.Next() {
			_, _, _ = iter.Current()
		}
		require.NoError(b, iter.Err())
	}
}

// BenchmarkM3TSZDecode64-12    	   17222	     69151 ns/op
func BenchmarkM3TSZDecode64(b *testing.B) {
	var (
		encodingOpts = encoding.NewOptions()
		seriesRun    = prepareSampleSeriesRun(b)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := NewReaderIterator64(seriesRun[i], DefaultIntOptimizationEnabled, encodingOpts)
		for iter.Next() {
			_, _, _ = iter.Current()
		}
		require.NoError(b, iter.Err())
	}
}

func prepareSampleSeriesRun(b *testing.B) [][]byte {
	var (
		rnd          = rand.New(rand.NewSource(42))
		sampleSeries = make([][]byte, 0, len(sampleSeriesBase64))
		seriesRun    = make([][]byte, 0, b.N)
	)

	for _, b64 := range sampleSeriesBase64 {
		data, err := base64.StdEncoding.DecodeString(b64)
		require.NoError(b, err)
		sampleSeries = append(sampleSeries, data)
	}

	for i := 0; i < b.N; i++ {
		seriesRun = append(seriesRun, sampleSeries[rnd.Intn(len(sampleSeries))])
	}

	return seriesRun
}
