FROM golang:1.16.0-stretch
 
WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV CGO_ENABLED=0
 
# Copie o arquivo Go principal para dentro do contêiner
COPY main.go .
 
# Execute a aplicação Go em segundo plano e mantenha o contêiner em execução
CMD go run main.go & tail -f /dev/null
