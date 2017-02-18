package challenges

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lost-mountain/isard/domain"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme"
	"gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// DNSResolver configures the DNS Record
// entry to validate the DNS challenge
type DNSResolver struct {
	acmeClient *acme.Client
	ns1Client  *rest.Client
}

// Cleanup removes the TXT record from the domain zone.
func (r *DNSResolver) Cleanup(d *domain.Domain) error {
	zone, err := r.getHostedZone(d.Name)
	if err != nil {
		return err
	}

	name := recordName(d.Name)
	_, err = r.ns1Client.Records.Delete(zone.Zone, name, "TXT")
	return errors.Wrapf(err, "error removing DNS record for domain challenge: %s", d.Name)
}

// Resolve uses the NS1 API to setup a TXT record
// for the ACME challenge.
func (r *DNSResolver) Resolve(d *domain.Domain, challenge *acme.Challenge) error {
	value, err := r.acmeClient.DNS01ChallengeRecord(challenge.Token)
	if err != nil {
		return errors.Wrapf(err, "error getting the DNS challenge record for %s", d.Name)
	}

	zone, err := r.getHostedZone(d.Name)
	if err != nil {
		return errors.Wrapf(err, "error getting the hosted zone for domain: %s", d.Name)
	}

	record := r.newTxtRecord(zone, d.Name, value)
	_, err = r.ns1Client.Records.Create(record)
	if err != nil && err != rest.ErrRecordExists {
		return errors.Wrapf(err, "error creating DNS record for domain challenge: %s", d.Name)
	}

	return nil
}

// NewDNSResolver uses NS1's api to resolve DNS challenges.
func NewDNSResolver(key string, ac *acme.Client) *DNSResolver {
	httpClient := &http.Client{Timeout: time.Second * 10}
	ns1Client := rest.NewClient(httpClient, rest.SetAPIKey(key))

	return &DNSResolver{
		acmeClient: ac,
		ns1Client:  ns1Client,
	}
}

func (r *DNSResolver) getHostedZone(domain string) (*dns.Zone, error) {
	zone, _, err := r.ns1Client.Zones.Get(domain)
	if err != nil {
		if err != rest.ErrZoneMissing {
			return nil, err
		}

		_, err = r.ns1Client.Zones.Create(&dns.Zone{Zone: domain})
		if err != nil {
			return nil, err
		}

		zone, _, err = r.ns1Client.Zones.Get(domain)
		fmt.Println(zone)
		if err != nil {
			return nil, err
		}
	}

	return zone, nil
}

func (r *DNSResolver) newTxtRecord(zone *dns.Zone, domain, value string) *dns.Record {
	return &dns.Record{
		Type:   "TXT",
		Zone:   zone.Zone,
		Domain: recordName(domain),
		Answers: []*dns.Answer{
			{Rdata: []string{value}},
		},
	}
}

func recordName(domain string) string {
	return fmt.Sprintf("_acme-challenge.%s", domain)
}
