FROM gcr.io/google.com/cloudsdktool/cloud-sdk:alpine

RUN gcloud components install kubectl

RUN apk add --no-cache jq

ADD entrypoint.sh /var/entrypoint.sh

ENTRYPOINT [ "/var/entrypoint.sh" ]

CMD ["gcloud"]