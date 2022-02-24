FROM alpine:latest

COPY ./clingo /usr/local/bin/clingo
RUN chmod a+x /usr/local/bin/clingo
