package domain

import (
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func listARecords(domain string) ([]dns.RR, error) {
	return listAnswers(domain, dns.TypeA)
}

func listAnswers(domain string, questionType uint16) ([]dns.RR, error) {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, errors.Wrapf(err, "error querying DNS servers for domain: %s", domain)
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), questionType)

	r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return nil, errors.Wrapf(err, "error querying DNS servers for domain: %s", domain)
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, errors.Wrapf(err, "error querying DNS servers for domain: %s", domain)
	}

	return r.Answer, nil
}
