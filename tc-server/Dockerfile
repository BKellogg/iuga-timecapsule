FROM alpine
RUN apk add --no-cache ca-certificates
COPY tc-server tc-server
EXPOSE 443
ENTRYPOINT ["/tc-server"]