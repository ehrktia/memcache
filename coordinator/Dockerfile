FROM debian:stable
COPY . .
RUN cp src/config.zgy ./zig-out/bin/
RUN chmod a+x zig-out/bin/coordinator && chmod a+x zig-out/bin/config.zgy && \
ls -lart zig-out/bin
CMD ["zig-out/bin/coordinator"]


