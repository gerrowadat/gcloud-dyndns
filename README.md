# gcloud-dyndns
dyndns-style Google Cloud DNS Updater

### Getting Started

This wee utility will update an A record in a DNS zone managed by cloud DNS. Let's kind of assume you have that setup already.

Get yourself a credentials JSON file (or don't if you're just going to run this interactively and you've got the ```gcloud``` command all set up). More complete information on how to do this setup is in the [README.md for clouddns-sync](https://github.com/gerrowadat/clouddns-sync#readme).

As part of that setup, you'll have a 'project' and 'zone' for the GCP project and DNS zone identifier (this is different than the actual domain name, because of course). 

Now, from the network you want to 'do dyndns' from, run something like:

```gcloud-dyndns --cloud-project=mycloudproject --cloud-dns-zone=mydnszone --cloud-dns-record-name=hostname.domain.tld.```

You can specify ```--json-keyfile=mykeyfle.json``` to give a credentials JSON file, and ```--dry-run``` if you're a bit scared and want to see what we'd do to your precious dns records.

That's it.