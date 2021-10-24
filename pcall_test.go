package pcall

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

// testHandler
type testHandler struct {
	Response *test.Case
	Next     plugin.Handler
}

type testcase struct {
	Expected int
	test     test.Case
	config   string
}

func (t *testHandler) Name() string { return "test-handler" }

func (t *testHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	d := new(dns.Msg)
	d.SetReply(r)
	if t.Response != nil {
		d.Answer = t.Response.Answer
		d.Rcode = t.Response.Rcode
	}
	w.WriteMsg(d)
	return 0, nil
}

func TestPcall(t *testing.T) {
	var tests = []testcase{
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.A("linux.example.org. 0 IN A 10.10.1.1")},
				Qname:  "linux.example.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("linux.example.org. 0 IN AAAA 2a00:1450:4009:823::200e")},
				Qname:  "linux.example.org.",
				Qtype:  dns.TypeAAAA,
			},
		},
		{
			Expected: dns.RcodeNameError,
			test: test.Case{
				Answer: []dns.RR{test.MX("example.org. 585 IN MX 50 mx01.example.org.")},
				Qname:  "linux.example.org.",
				Qtype:  dns.TypeMX,
			},
		},
	}

	for _, tc := range tests {

		m := new(dns.Msg)
		m.SetQuestion(tc.test.Qname, tc.test.Qtype)

		tHandler := &testHandler{
			Response: &tc.test,
			Next:     nil,
		}
		o := &Pcall{Next: tHandler}
		w := dnstest.NewRecorder(&test.ResponseWriter{})

		//default value for the binary
		if tc.config == "" {
			o.CommandPath = "./test/resolver"
		} else {
			o.CommandPath = tc.config
		}

		_, err := o.ServeDNS(context.TODO(), w, m)

		if err != nil {
			t.Errorf("Error %q", err)
		}

		if w.Rcode != tc.Expected {
			t.Fatal("Expected:", tc.Expected, "Got:", w.Rcode)
		}

		//nothing to check if nxdomain
		if tc.Expected == dns.RcodeNameError {
			continue
		}

		if tc.test.Answer[0].String() != w.Msg.Answer[0].String() {
			t.Error("Expected:", tc.test.Answer[0], "Got:", w.Msg.Answer[0], "Rcode:", w.Rcode)
		}

	}
}
