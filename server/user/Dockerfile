FROM scratch

WORKDIR /
COPY . .

EXPOSE 8222 8224
ENTRYPOINT ["./user", "server", "-c", "./config/config-docker.ini"]