FROM envoyproxy/envoy-alpine

# add protobuf descriptor
COPY descriptor.pb /
COPY envoy.yaml /

ADD server /bin
ADD entrypoint.sh /bin

RUN chmod +x /bin/entrypoint.sh && mkdir -p /var/log/envoy

ENTRYPOINT ["entrypoint.sh"]

EXPOSE 8080 8090 9901