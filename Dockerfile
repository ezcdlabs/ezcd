FROM scratch
COPY ezcd-server /ezcd-server
ENTRYPOINT ["/ezcd-server"]
EXPOSE 3923