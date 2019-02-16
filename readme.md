# local run harness for Firestore triggered golang functions


https://cloud.google.com/functions/docs/calling/cloud-firestore#functions_eventdata-go

Develop locally by modifying and then running harness/main.go

It will create a listen stream of snapshots, and invoke your function code in a way that seems pretty much identical to the way the runtime works deployed.

When you are happy with it, deploy:

```
gcloud functions deploy DocChange --runtime go111 \
  --trigger-event providers/cloud.firestore/eventTypes/document.write \
  --trigger-resource projects/ptone-serverless/databases/(default)/documents/cities/{any}
```