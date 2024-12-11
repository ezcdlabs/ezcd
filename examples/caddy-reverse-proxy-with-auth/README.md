# Caddy Reverse Proxy With Basic Auth

This example includes a very simple Caddyfile that demonstrates how you can add authentication to EZCD.

Caddyfile

```caddy
# you can replace this with your hostname or port that you want to listen on
ezcd.localhost:443

basic_auth /* {
	# user:password created using caddy hash-password
	user "$2a$14$kbCcVZablzYeFEz2Q3pq5uIyHJTK.Sogfih1YydwDTz5ZmvmvF8Xa"
}
reverse_proxy localhost:3923
```

Once you have set up postgres

Run the ezcd-server locally, e.g. for version 0.1.15:

Make sure your database is set up and accessible using your database url and then you can run the dashboard like this:

```sh
docker run -p 3923:3923 -e EZCD_DATABASE_URL=$EZCD_DATABASE_URL ghcr.io/ezcdlabs/ezcd-server:0.1.15
```
