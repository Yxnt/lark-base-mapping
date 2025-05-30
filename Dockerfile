# Start by building the application.
FROM docker.ispider.io/golang:1.23 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.ispider.io/distroless/static-debian12
COPY --from=build /go/bin/app /
CMD ["/app"]