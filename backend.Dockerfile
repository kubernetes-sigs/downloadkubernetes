FROM golang:1.12

WORKDIR /downloadkubernetes
ADD go.mod .
ADD main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o downloadkubernetes .

WORKDIR /script
ADD ./scripts/go.mod .
ADD ./scripts/go.sum .
RUN go mod download
ADD ./scripts .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o updateindex ./scripts

FROM node:12
COPY . /dlk8s
COPY --from=0 /downloadkubernetes/downloadkubernetes /dlk8s/
COPY --from=0 /script/updateindex /dlk8s/scripts
ENV PATH=/dlk8s/scripts:$PATH
ENTRYPOINT [ "/dlk8s/downloadkubernetes" ]
