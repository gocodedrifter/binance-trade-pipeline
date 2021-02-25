FROM golang:1.12.14-alpine as build-env
# All these steps will be cached
RUN apk add git
RUN mkdir /coinpanel
WORKDIR /coinpanel
COPY go.mod . 

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
WORKDIR /coinpanel/cmd/myapp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/coinpanel
FROM golang:1.12.14-alpine
COPY --from=build-env /go/bin/coinpanel /go/bin/coinpanel
ENTRYPOINT ["/go/bin/coinpanel"]