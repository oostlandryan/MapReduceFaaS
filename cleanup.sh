#!/bin/bash

# Delete Cloud Functions
gcloud functions delete RyoostCreateIndexHttp --quiet
gcloud functions delete RyoostSearchHttp --quiet
gcloud functions delete RyoostMapHttp --quiet
gcloud functions delete RyoostReduceHttp --quiet

# Delete entire ryoost-mapreduce Firestore Collection
go mod download
go run deleteFirestoreCollection.go