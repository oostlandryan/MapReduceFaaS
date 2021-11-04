#!/bin/bash
cd mapper
gcloud functions deploy RyoostMapHttp --runtime go113 --trigger-http --allow-unauthenticated
cd ../reducer
gcloud functions deploy RyoostReduceHttp --runtime go113 --trigger-http --allow-unauthenticated
cd ../inverseindex
gcloud functions deploy RyoostCreateIndexHttp --runtime go113 --trigger-http --allow-unauthenticated --timeout 300 --memory 1024
cd ../searchindex
gcloud functions deploy RyoostSearchHttp --runtime go113 --trigger-http --allow-unauthenticated
cd ../
go run fileupload.go