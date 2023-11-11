package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	google_oauth "golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type CloudDNSSpec struct {
	svc         *dns.Service
	project     *string
	zone        *string
	default_ttl *int
	dry_run     *bool
}

func getMyIP() (string, error) {
	res, err := http.Get("http://whatismyip.akamai.com")
	if err != nil {
		log.Fatal("HTTP Error getting our IP: ", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	ip := net.ParseIP(string(resBody))
	if ip == nil {
		log.Fatalf("Non-IP returned from request: %s", ip)
	}

	return string(resBody), nil
}

func updateRecord(dns_spec *CloudDNSSpec, record_name string, old_ip, new_ip string) error {

	log.Printf("Updating Cloud DNS: %s : %s -> %s", record_name, old_ip, new_ip)

	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			{
				Name:    record_name,
				Type:    "A",
				Rrdatas: []string{new_ip},
			},
		},
	}

	// Gcloud DNS shits the bed if you try to delete a record that's not there.
	if old_ip != "" {
		new_rr := dns.ResourceRecordSet{
			Name:    record_name,
			Type:    "A",
			Rrdatas: []string{old_ip},
		}
		change.Deletions = append(change.Deletions, &new_rr)
	}

	call := dns_spec.svc.Changes.Create(*dns_spec.project, *dns_spec.zone, change)

	if !*dns_spec.dry_run {
		out, err := call.Do()
		if err != nil {
			log.Fatal("Error updating Cloud DNS: ", err)
		}
		log.Printf("Added [%d] and deleted [%d] records.", len(out.Additions), len(out.Deletions))
	}

	return nil

}

func main() {
	var jsonKeyfile = flag.String("json-keyfile", "", "json credentials file for Cloud DNS")
	var cloudProject = flag.String("cloud-project", "", "Google Cloud Project")
	var cloudZone = flag.String("cloud-dns-zone", "", "Cloud DNS zone to operate on")
	var cloudDnsRecordName = flag.String("cloud-dns-record-name", "", "Cloud DNS zone to operate on")
	var defaultCloudTtl = flag.Int("cloud-dns-default-ttl", 300, "Default TTL for Cloud DNS records")

	var dryRun = flag.Bool("dry-run", false, "Do not update Cloud DNS, print what would be done")

	flag.Parse()

	// These are required in all cases
	if *cloudZone == "" {
		log.Fatal("--cloud-dns-zone is required")
	}
	if *cloudDnsRecordName == "" {
		log.Fatal("--cloud-dns-record-name is required")
	}

	ctx := context.Background()

	creds := &google_oauth.Credentials{}

	if *jsonKeyfile != "" {
		jsonData, ioerror := os.ReadFile(*jsonKeyfile)
		if ioerror != nil {
			log.Fatal(*jsonKeyfile, ioerror)
		}
		creds, _ = google_oauth.CredentialsFromJSON(ctx, jsonData, "https://www.googleapis.com/auth/cloud-platform")
	} else {
		creds, _ = google_oauth.FindDefaultCredentials(ctx)
	}

	// Get project from json keyfile if present.
	if creds.ProjectID != "" {
		*cloudProject = creds.ProjectID
	}

	if *cloudProject == "" {
		log.Fatal("--cloud-project is required")
	}

	dnsservice, err := dns.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatal("Cloud DNS Error: ", err)
	}

	dns_spec := &CloudDNSSpec{
		svc:         dnsservice,
		project:     cloudProject,
		zone:        cloudZone,
		default_ttl: defaultCloudTtl,
		dry_run:     dryRun,
	}

	my_ip, err := getMyIP()
	if err != nil {
		log.Fatalf("Error getting our IP: %s", err)
	}

	log.Printf("Detected IP: %s", my_ip)

	current_dns, err := net.LookupIP(*cloudDnsRecordName)

	current_ip := ""

	if err != nil {
		log.Print("Error in DNS resolution: ", err)
		log.Print("Continuing...")
	} else {
		if len(current_dns) > 1 {
			log.Fatalf("%s resolves to multiple IPs. Weird.", *cloudDnsRecordName)
		}

		current_ip = current_dns[0].To4().String()

		if my_ip == current_ip {
			log.Print("My IP matches DNS. Nothing to do.")
			return
		}
	}

	err = updateRecord(dns_spec, *cloudDnsRecordName, current_ip, my_ip)

	if err != nil {
		log.Fatalf("Error Updating GCloud: %s", err)
	}

}
