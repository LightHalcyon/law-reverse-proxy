# law-reverse-proxy
Golang NGINX reverse proxy implementation for LAW class assignment

## Usage
- Copy nginx.conf to your `nginx/conf` folder
- Execute `go run .` on `dl-server` and `ul-server`
- Run/restart NGINX & access `http://localhost:22006` from your browser

Alternatively, run `docker-compose build` and `docker-compose run`