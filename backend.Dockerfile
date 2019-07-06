FROM golang:1.12

WORKDIR /downloadkubernetes
ADD go.mod .
ADD main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o downloadkubernetes .

WORKDIR /script
ADD ./scripts/go.mod .
ADD ./scripts/go.sum .
RUN go mod download
ADD ./scripts/generate-index.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o updateindex .

FROM node:12
WORKDIR /dlk8s
COPY . .
RUN npm install
COPY --from=0 /downloadkubernetes/downloadkubernetes .
COPY --from=0 /script/updateindex ./scripts
ENTRYPOINT [ "./downloadkubernetes" ]
