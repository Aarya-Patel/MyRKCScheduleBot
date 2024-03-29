FROM golang:1.21

WORKDIR /bot

ENV TOKEN="<TELEGRAM_API_TOKEN>"
ENV WEBHOOK_DOMAIN="<WEBHOOK_DOMAIN>"
ENV WEBHOOK_SECRET="<WEBHOOK_SECRET>"
ENV PROJECT_ID="<GCP_PROJECT_ID>"

EXPOSE 8080/tcp

COPY . .
RUN go mod download && go mod verify
RUN go build ./cmd/main.go

CMD [ "./main" ]