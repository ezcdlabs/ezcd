FROM scratch
COPY ezcd-server /ezcd-server
ENTRYPOINT ["/ezcd-server"]