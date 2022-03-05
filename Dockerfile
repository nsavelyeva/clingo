FROM alpine:latest

RUN mkdir -p /opt/clingo /opt/clingo/weather /opt/clingo/currency

COPY ./weather/conditions.csv /opt/clingo/weather/conditions.csv
COPY ./currency/details.json /opt/clingo/currency/details.json
COPY ./clingo /opt/clingo/clingo

RUN chmod a+x /opt/clingo/clingo \
    && ln -s /opt/clingo/clingo /usr/local/bin/clingo

WORKDIR /opt/clingo
