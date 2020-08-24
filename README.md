# Download Kubernetes

## Build locally

1. Run `go run ./cmd/update-index/main.go`
2. Build assets with `npm run build`
3. Open `dist/index.html` with a browser, from the command line: `open dist/index.html`

## Architecture

The published artifacts are static HTML/CSS/JavaScript files. They are updated offline by the `update-index` command and
then published to the web via docker.

The docker image has an nginx in it to serve static files.

Assumed HTTPS is terminated at the ingress level before reaching the docker image.
