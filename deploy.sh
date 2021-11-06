#!/bin/bash

# NOTE change the following before running:
#   projectID in fileupload.go, ./searchindex/search_index.go, ./mapper/map_wrapper.go, ./inverseindex/inverse_index_wrapper.go 
#   mapFuncUrl in ./inverseindex/inverse_index.go
#   reduceFuncUrl in ./inverseindex/inverse_index.go
#
# All changes should be at the top of the files, just under the import statements.
# I could automate this if running locally, but couldn't figure it out with running them on Cloud Functions

# Deploy Mapper
cd mapper
gcloud functions deploy RyoostMapHttp --runtime go113 --trigger-http --allow-unauthenticated

# Deploy Reducer
cd ../reducer
gcloud functions deploy RyoostReduceHttp --runtime go113 --trigger-http --allow-unauthenticated

# Deploy Create Index Function
cd ../inverseindex
gcloud functions deploy RyoostCreateIndexHttp --runtime go113 --trigger-http --allow-unauthenticated --timeout 300 --memory 1024

# Deploy Search Index Function
cd ../searchindex
gcloud functions deploy RyoostSearchHttp --runtime go113 --trigger-http --allow-unauthenticated

# Upload Files to Firestore
cd ../
go mod download
go run fileupload.go