FROM alpine
RUN apk add --no-cache --update ca-certificates
ADD app /app
CMD ["./app"]