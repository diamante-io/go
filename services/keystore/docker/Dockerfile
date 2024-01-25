FROM golang:1.17 as build

ADD . /src/keystore
WORKDIR /src/keystore
RUN go build -o /bin/keystored ./services/keystore/cmd/keystored


FROM ubuntu:18.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
ADD ./services/keystore/migrations/ /app/migrations/
COPY --from=build /bin/keystored /app/
EXPOSE 8000
ENTRYPOINT ["/app/keystored"]
CMD ["serve"]
