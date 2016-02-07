package dnsServer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
	"net"
)

// Catch all is a struct type that represents data to send to a client that
// requests an A record for any host.
// It also has methods for running a dns server with these credentials.
type CatchAll struct {
	Domain      string
	Port        string
	IP          net.IP
	SpoofDomain bool
}

// Serve starts a sever for a given net type (udp or tcp)
func (c *CatchAll) Serve(netType string) error {
	dns.HandleFunc(".", c.handleDNS)

	server := &dns.Server{Addr: c.Port, Net: netType, TsigSecret: nil}

	log.Infof("Serve %s", netType)

	return server.ListenAndServe()
}

// handleDNS is a handler function to actualy perform the dns querey response
func (c *CatchAll) handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	var rr dns.RR

	domainSpoof := r.Question[0].Name

	msgResp := new(dns.Msg)
	msgResp.SetReply(r)
	msgResp.Compress = false

	rr = new(dns.A)

	if c.SpoofDomain {
		rr.(*dns.A).Hdr = dns.RR_Header{Name: domainSpoof, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
	} else {
		rr.(*dns.A).Hdr = dns.RR_Header{Name: c.Domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
	}

	rr.(*dns.A).A = c.IP

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msgResp.Answer = append(msgResp.Answer, rr)
	default:
		log.Warnf("Unknown dns type %T", r.Question[0].Qtype)
		return
	}

	w.WriteMsg(msgResp)
}
