FROM alpine:3.10

RUN apk add --no-cache curl tzdata

COPY ./firebase-adminsdk.key.json /.gcp/firebase-adminsdk.key.json

ENV GOOGLE_APPLICATION_CREDENTIALS=/.gcp/firebase-adminsdk.key.json

WORKDIR /app/svc
COPY ./dist/playground-linux-x64 .

CMD ./playground-linux-x64
