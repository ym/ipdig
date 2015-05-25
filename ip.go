package main

import (
	dnslib "github.com/miekg/dns"
	"github.com/wangtuanjie/ip17mon"
	"reflect"
	"strconv"
)

func getIpInfo(ip string) string {
	loc, err := ip17mon.Find(ip)
	if err != nil {
		return "ERROR"
	}
	return loc.Country + " " + loc.Region + " " + loc.City + " " + loc.Isp
}

func showIpInfo(rr dnslib.RR) string {
	if reflect.TypeOf(rr).Elem().String() != "dns.A" {
		return ""
	}

	ip := rr.(*dnslib.A).A.String()

	return " (" + getIpInfo(ip) + ")"
}

func ipMsg(dns *dnslib.Msg) string {
	if dns == nil {
		return "<nil> MsgHdr"
	}
	s := dns.MsgHdr.String() + " "
	s += "QUERY: " + strconv.Itoa(len(dns.Question)) + ", "
	s += "ANSWER: " + strconv.Itoa(len(dns.Answer)) + ", "
	s += "AUTHORITY: " + strconv.Itoa(len(dns.Ns)) + ", "
	s += "ADDITIONAL: " + strconv.Itoa(len(dns.Extra)) + "\n"
	if len(dns.Question) > 0 {
		s += "\n;; QUESTION SECTION:\n"
		for i := 0; i < len(dns.Question); i++ {
			s += dns.Question[i].String() + "\n"
		}
	}
	if len(dns.Answer) > 0 {
		s += "\n;; ANSWER SECTION:\n"
		for i := 0; i < len(dns.Answer); i++ {
			if dns.Answer[i] != nil {
				s += dns.Answer[i].String() + showIpInfo(dns.Answer[i]) + "\n"
			}
		}
	}
	if len(dns.Ns) > 0 {
		s += "\n;; AUTHORITY SECTION:\n"
		for i := 0; i < len(dns.Ns); i++ {
			if dns.Ns[i] != nil {
				s += dns.Ns[i].String() + showIpInfo(dns.Ns[i]) + "\n"
			}
		}
	}
	if len(dns.Extra) > 0 {
		s += "\n;; ADDITIONAL SECTION:\n"
		for i := 0; i < len(dns.Extra); i++ {
			if dns.Extra[i] != nil {
				s += dns.Extra[i].String() + showIpInfo(dns.Extra[i]) + "\n"
			}
		}
	}
	return s
}

func initIpLibrary(path string) {
	if err := ip17mon.Init(path); err != nil {
		panic(err)
	}
}
