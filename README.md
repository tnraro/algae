# Algae

Algae simplifies the Docker deployment process

Create alga:

`POST /algae { name: string, compose: string, env: string }`

Update alga:

`PATCH /algae/:name { compose?: string, env?: string }`

Delete alga:

`DELETE /algae/:name`

To start Algae:

```sh
docker run --rm -p 41943:41943 -v /var/run/docker.sock:/var/run/docker.sock -v ./data:/app/data tnraro/algae:latest
```

[compose.yml example](compose.yml)