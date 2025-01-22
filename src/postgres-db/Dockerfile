FROM postgres:latest

ENV POSTGRES_USER=dbuser
ENV POSTGRES_PASSWORD=dbpass
ENV POSTGRES_DB=portfolio_simulator

EXPOSE 5432

COPY init.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]
