FROM golang:1.19.4
LABEL version="1.0" maintainer="Andrei <koodJõhvi@koodJõhvi.ee>"

# RUN useradd -m koodJõhvi 
# RUN  adduser koodJõhvi sudo

# USER koodJõhvi

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . . 

RUN go build -o /docker-server

EXPOSE 8080

CMD [ "/docker-server" ]
