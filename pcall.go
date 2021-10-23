package pcall

import (
	"context"
	"log"
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
		log.Fatalln("Multi-Question not supported")
		return 0, nil
	}

	qname := r.Question[0].Name
	qtype := r.Question[0].Qtype
	class := r.Question[0].Qclass

	//int is just a tmp solution for now

	cmd := exec.Command(p.CommandPath, dns.TypeToString[qtype], qname)

	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalln("command err", p.CommandPath, err)
		return 0, nil
	}

	ip := net.ParseIP(strings.Trim(string(stdout), "\n\t"))

	log.Println("Qname:", qname, "QType", qtype, ip)

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

	ans := new(dns.Msg)

	//add the answer
	ans.Answer = append(ans.Answer, rr)

	w.WriteMsg(ans)

	return 0, nil
}

// Name implements the Handler interface.
func (a Pcall) Name() string { return "pcall" }
