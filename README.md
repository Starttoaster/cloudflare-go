# cloudflare-go
Dynamic DNS record updater written in Golang

[![pipeline status](https://gitlab.com/brandonbutler/cloudflare-go/badges/master/pipeline.svg)](https://gitlab.com/brandonbutler/cloudflare-go/commits/master)
[![coverage report](https://gitlab.com/brandonbutler/cloudflare-go/badges/master/coverage.svg)](https://gitlab.com/brandonbutler/cloudflare-go/commits/master)

# Use

You can either start this container via `docker run` or `docker-compose`.

**Quick Start:** 
```
docker container run \
  -p 8080:8080 \
  -e CF_EMAIL="email@email.com" \
  -e CF_KEY="my_global_api_key" \
  -e CF_ZONE="my_zone_id" \
  starttoaster/cloudflare-go:latest
```

View the webUI in your web browser to see it run: `localhost:8080`

![Image of webUI](https://i.imgur.com/qziKSxC.png)


Example docker-compose.yml:

```
version: '2'
services:

 cloudflare:
   container_name: cloudflare-dns
   image: starttoaster/cloudflare-go:latest
   environment:
     - CF_EMAIL="email@email.com"
     - CF_KEY="my_global_api_key"
     - CF_ZONE="my_zone_id"
```

**NOTE:** The bare minimum required details to interact with Cloudflare's DNS API is the account email, Global API Key, and Zone ID attributes. This container finds other necessary attributes 
without requiring user input by itself, but if you need help finding the 3 required attributes I listed, view the section below regarding "API Attributes".

# Environment Variables

| Variable | Function |
| ---- | ---- | 
| -e CF_EMAIL | The email address associated with your Cloudflare account | 
| -e CF_KEY | A unique 'Global API Key' is assigned to every Cloudflare DNS user | 
| -e CF_ZONE | A unique 'Zone ID' is assigned to each domain registered in a Cloudflare account | 

### How to find the attributes for credfile

  1. Account Email -- This is just the email address your Cloudflare account is under.

  2. Global API Key -- This is found in the Cloudflare website under "My Profile > API Keys > Global API Key"

  3. Zone ID -- This ID is a random string of letters and numbers specific to your domain name. Found in the Cloudflare website on your domain's "Overview" page written as "Zone ID" 