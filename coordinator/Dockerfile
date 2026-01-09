FROM alpine:3.23.2
RUN mkdir -p /app
WORKDIR /app
COPY ./config.zgy .
COPY ./zig-out/bin/coordinator .
RUN ls -lart .
EXPOSE 9999/tcp
CMD ["/app/coordinator"]


