FROM golang:1.19 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/api
CMD ["air"]

FROM base as built
WORKDIR /go/app/api
COPY . .
ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go build -o /tmp/kamogawa-server ./*.go

FROM postgres:13.8-alpine
COPY --from=built /tmp/kamogawa-server /usr/bin/kamogawa-server
WORKDIR /usr/bin

USER postgres
ENV SHIMOGAWA_URL postgres://postgres@0.0.0.0:5432/shimogawa_db
RUN chmod 0700 /var/lib/postgresql/data &&\
    initdb /var/lib/postgresql/data &&\
    echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf
CMD pg_ctl start; psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'shimogawa_db'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE shimogawa_db"; psql -U postgres "CREATE EXTENSION pg_trgm;"; psql -U postgres "CREATE INDEX gce_instance_dbs_idx ON gce_instance_dbs USING GIN ((name || ' ' || id || ' ' || project_id || ' ' || zone) gin_trgm_ops);"; CMD psql -U postgres "CREATE INDEX project_dbs_idx ON project_dbs USING GIN ((name || ' ' || project_number || ' ' || project_id) gin_trgm_ops);"; kamogawa-server; 
