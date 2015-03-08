FROM dockerfile/go

MAINTAINER Rui Lopes <rgl@ruilopes.com>

COPY Makefile *.go ./

RUN make
