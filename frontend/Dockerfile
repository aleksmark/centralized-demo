FROM partlab/ubuntu-golang

WORKDIR /app

ADD src/bin/frontend /app/
ADD src/index.html /app/

ENTRYPOINT ["./frontend"]
