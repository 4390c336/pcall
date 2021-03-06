package pcall

import (
	"context"
	"net"
	"os/exec"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Pcall struct {
	Next        plugin.Handler
	CommandPath string
}

// ServeDNS implements the plugin.Handler interface.
func (p Pcall) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	if len(r.Question) > 1 {
		return 0, nil
	}

	m := new(dns.Msg)
	m.SetReply(r)

	qname := r.Question[0].Name
	qtype := r.Question[0].Qtype
	class := r.Question[0].Qclass

	if qtype != dns.TypeA && qtype != dns.TypeAAAA {
		//response with nxdomain
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return dns.RcodeNameError, nil
	}

	cmd := exec.Command(p.CommandPath, dns.TypeToString[qtype], qname)
	stdout, err := cmd.Output()

	if err != nil {
		//response with nxdomain
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return dns.RcodeNameError, nil
	}

	//@TODO: support multiple answers

	ip := net.ParseIP(strings.Trim(string(stdout), "\n\t "))

	if ip == nil {
		//response with nxdomain
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return dns.RcodeNameError, nil
	}

	var rr dns.RR

	if qtype == dns.TypeA {
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: class}
		rr.(*dns.A).A = ip.To4()
	}

	if qtype == dns.TypeAAAA {
		rr = new(dns.AAAA)
		rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: class}
		rr.(*dns.AAAA).AAAA = ip.To16()
	}

	//add the answer
	m.Answer = append(m.Answer, rr)

	w.WriteMsg(m)

	return 0, nil
}

// Name implements the Handler interface.
func (p Pcall) Name() string { return "pcall" }
