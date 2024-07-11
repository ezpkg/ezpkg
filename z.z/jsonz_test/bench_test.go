package test

import (
	stdjson "encoding/json"
	"io"
	"sync/atomic"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/bytedance/sonic"
	"github.com/iOliverNguyen/ujson"
	jsoniter "github.com/json-iterator/go"
	pkgjson "github.com/pkg/json"
	"github.com/valyala/fastjson"
)

var fastjsonPool fastjson.ParserPool

func BenchmarkParse(b *testing.B) {
	for _, parallel := range []bool{false, true} {
		for _, tc := range testcases {
			tc = tc.WithParallel(parallel)
			{
				var nTokens atomic.Int64
				err := bench(b, "ujson", tc, func(tc *Testcase) error {
					n := 0
					err := ujson.Walk(tc.data, func(level int, key, value []byte) ujson.WalkFuncRtnType {
						if len(key) != 0 {
							n++
						}
						if len(value) != 0 {
							n++
						}
						return ujson.WalkRtnValDefault
					})
					nTokens.Store(int64(n))
					return err
				})
				must(0, err)
				if n := nTokens.Load(); n > 0 {
					debugf(b, "%10v: ntokens=%v", "ujson", n)
				}
			}
			{
				err := bench(b, "stdjson", tc, func(tc *Testcase) error {
					var value any
					dec := stdjson.NewDecoder(tc.Reader())
					return dec.Decode(&value)
				})
				must(0, err)
			}
			{
				type State struct {
					buf []byte
				}
				var nTokens atomic.Int64
				tc = tc.WithInitFunc(func(tc *Testcase) any {
					return State{buf: make([]byte, 8<<10)}
				})
				err := bench(b, "pkgjson", tc, func(tc *Testcase) error {
					st := tc.State().(State)
					dec := pkgjson.NewDecoderBuffer(tc.Reader(), st.buf[:])
					n := 0
					for {
						_, err := dec.NextToken()
						if err == io.EOF {
							break
						}
						if err != nil {
							return err
						}
						n++
					}
					nTokens.Store(int64(n))
					return nil
				})
				must(0, err)
				if n := nTokens.Load(); n > 0 {
					debugf(b, "%10v: tokens=%v\n", "pkgjson", n)
				}
			}
			{
				var result atomic.Pointer[fastjson.Value]
				err := bench(b, "fastjson", tc, func(tc *Testcase) error {
					parser := fastjsonPool.Get()
					value, err := parser.ParseBytes(tc.data)
					fastjsonPool.Put(parser)
					result.Store(value)
					return err
				})
				if err != nil {
					b.Logf("%10v: error %v", "fastjson", shortenError(err))
				} else if res := result.Load(); res != nil {
					debugf(b, "%10v: keys=%v", "fastjson", res.GetObject().Len())
				}
			}
			{
				var nKeys atomic.Int64
				path := []string{"foo"}
				err := bench(b, "jsonparser", tc, func(tc *Testcase) (outErr error) {
					n := 0
					jsonparser.EachKey(tc.data, func(i int, bytes []byte, valueType jsonparser.ValueType, err error) {
						if err != nil && outErr == nil {
							outErr = err
						}
						n++
					}, path)
					nKeys.Store(int64(n))
					return outErr
				})
				must(0, err)
				if n := nKeys.Load(); n > 0 {
					debugf(b, "%10v: keys=%v", "jsonparser", n)
				}
			}
			{
				err := bench(b, "jsoniter", tc, func(tc *Testcase) error {
					var value any
					dec := jsoniter.NewDecoder(tc.Reader())
					return dec.Decode(&value)
				})
				must(0, err)
			}
			{
				err := bench(b, "sonic", tc, func(tc *Testcase) error {
					var value any
					dec := sonic.ConfigDefault.NewDecoder(tc.Reader())
					return dec.Decode(&value)
				})
				must(0, err)
			}
		}
	}
}
