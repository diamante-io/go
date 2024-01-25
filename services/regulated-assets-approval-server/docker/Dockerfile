FROM golang:1.17 as build

ADD . /src/regulated-assets-approval-server
WORKDIR /src/regulated-assets-approval-server
RUN go build -o /bin/regulated-assets-approval-server ./services/regulated-assets-approval-server


FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/regulated-assets-approval-server /app/
EXPOSE 8000
ENTRYPOINT ["/app/regulated-assets-approval-server"]
CMD ["serve"]
