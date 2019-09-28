# cloudflare-go
Dynamic DNS record updater written in Golang

## Pre-requisite

To actually interact with Cloudflare, this needs a separate file provided by you, with 3 lines. The name of the file should be `credfile` and have the following in order:

```
email address
global API key
zone identifier
```

### How to find the attributes for credfile

  1. Account Email -- This is just the email address your Cloudflare account is under.

  2. Global API Key -- This is found in the Cloudflare website under "My Profile > API Keys > Global API Key"

  3. Zone ID -- This ID is a random string of letters and numbers specific to your domain name. Found in the Cloudflare website on your domain's "Overview" page written as "Zone ID" 


## Use with Docker

Start your container: `docker container run -d -p 8080:8080 -v "$PWD/credfile:/config/credfile" starttoaster/cloudflare-go`

View the webUI in your web browser to see it run: `localhost:8080`

![Image of webUI](https://i.imgur.com/qziKSxC.png)
