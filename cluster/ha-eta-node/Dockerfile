FROM ha-eta-build as build

FROM debian:10

RUN mkdir -p /var/log/ha-eta-node

COPY --from=build /usr/local/src/ha-eta/ha-eta /usr/local/bin/ha-eta
COPY --from=build /usr/local/src/ha-eta/ha-eta-node /usr/local/bin/ha-eta-node

EXPOSE 80

RUN ha-eta --version
RUN ha-eta-node --version

ENTRYPOINT ["ha-eta"]
CMD ["node"]
