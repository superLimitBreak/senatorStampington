package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/SuperLimitBreak/senatorStampington/dnsServer"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create catch all dns server
	server := dnsServer.CatchAll{
		Domain: "superlimitbreak.uk.",
		Port:   ":53",
		IP:     net.ParseIP("192.168.0.1").To4(),
	}

	// Wrap the serve method with exit logging
	serveWrap := func(netType string) {
		if err := server.Serve(netType); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"NetType": netType,
			}).Fatal("Failed to setup server")
		}
	}

	// Start the two servers async
	go serveWrap("tcp")
	go serveWrap("udp")

	// create chanels to listen for os signals
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Seperate channel as we are expected to provide a core dump for a SIGQUIT
	dump := make(chan os.Signal)
	signal.Notify(dump, syscall.SIGQUIT)

	// Wait for a signal to exit
	select {
	case <-dump:
		panic(errors.New("SIGQUIT received: panicing!"))
	case <-exit:
		return
	}
}
