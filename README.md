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

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack](http://slack.k8s.io/)
- [Mailing List](https://groups.google.com/forum/#!forum/kubernetes-dev)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
