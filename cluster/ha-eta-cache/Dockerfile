FROM debian:10

RUN apt-get update
RUN apt-get install -y \
  nginx

RUN mkdir -p /var/lib/nginx/

EXPOSE 80

ENTRYPOINT ["nginx"]
CMD ["-g", "daemon off;"]
