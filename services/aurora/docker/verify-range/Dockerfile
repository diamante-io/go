FROM ubuntu:20.04

MAINTAINER Bartek Nowotarski <bartek@diamcircle.org>

ARG diamcircle_CORE_VERSION 
ENV diamcircle_CORE_VERSION=${diamcircle_CORE_VERSION:-*}
# to remove tzdata interactive flow
ENV DEBIAN_FRONTEND=noninteractive

ADD dependencies /
RUN ["chmod", "+x", "dependencies"]
RUN /dependencies

ADD diamcircle-core.cfg /

ADD start /
RUN ["chmod", "+x", "start"]

ENTRYPOINT ["/start"]
