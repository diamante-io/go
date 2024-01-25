FROM golang:1.17 as build

ADD . /src/recoverysigner
WORKDIR /src/recoverysigner
RUN go build -o /bin/recoverysigner ./exp/services/recoverysigner


FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/recoverysigner /app/
EXPOSE 8000
ENTRYPOINT ["/app/recoverysigner"]
CMD ["serve"]
