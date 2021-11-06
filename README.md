# MapReduceFaaS
Implements MapReduce using FaaS [assignment](https://cgi.luddy.indiana.edu/~prateeks/cloud/faas-mr.html) for Enginering Cloud Computing ENGR-E 516 at Indiana University.

## Assignment Summary
The core of the assignment is to create an inverse index of word counts for a corpus of books from Project Gutenberg using map and reduce functions deployed on Google Cloud Functions. This allows for a parallel and distributed version of mapreduce using functions as a service. Extra credit was also given for creating a simple frontend for search the inverted index, which can be found in the [webpage](./webpage) folder. Additional details required for the assignment can be found in [Report.md](Report.md).

## Scripts
[deploy.sh](deploy.sh) deploys the required Google Cloud Functions and uploads my selected books to Firestore. [demo.sh](demo.sh) creates the inverted index and searches two words as demonstration. [cleanup.sh](cleanup.sh) deletes all of the created documents in Firestore and also deletes the functions that were deployed in [deploy.sh](deploy.sh).
