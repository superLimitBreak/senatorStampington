package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	domain   = "superlimitbreak.uk."
	dnsPort  = ":53"
	catchAll = "192.168.0.1"
)

var catchAllIP net.IP

func main() {
	log.Info("Dns Server starts")

	dns.HandleFunc(".", handleDNS)

	catchAllIP = net.ParseIP(catchAll).To4()
	if catchAllIP == nil {
		log.Fatal("failed to parse the catch all ip address")
	}

	go serve("tcp")
	go serve("udp")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func serve(netType string) {
	server := &dns.Server{Addr: dnsPort, Net: netType, TsigSecret: nil}

	log.Infof("Serve %s", netType)

	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"NetType": netType,
		}).Fatal("Failed to setup server")
	}
}

func handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	var (
		rr dns.RR
	)

	msgResp := new(dns.Msg)
	msgResp.SetReply(r)
	msgResp.Compress = false

	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
	rr.(*dns.A).A = catchAllIP

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msgResp.Answer = append(msgResp.Answer, rr)
	default:
		//log
		return
	}

	w.WriteMsg(msgResp)
}
