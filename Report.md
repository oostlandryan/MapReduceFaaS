# Running Instructions
## Webpage on my project
I created a basic front-end running on Google's Cloud Run, which is accessible [here](https://wordsearch-ojbmp4ml5a-uc.a.run.app). The webpage queries the RyoostSearchHttp function and RyoostCreateIndexHttp function. It required a bit of extra work to set up and it's extra credit for the assignment, so I did not include it in the deployment script. However, the files themselves are still included with the assignment. If there is an inverted index in my Firestore, the page allows you to search it or recreate it. If there is no inverted index, it allows you to create one and then search it.
## Deploying to a new project
Running [deploy.sh](deploy.sh) will deploy all of the needed cloud functions and upload the corpus from Project Gutenberg to Firestore. The Firestore CLI might be necessary and you'll have to sign into it separately from gcloud.  
Prior to running the deploy and cleanup scripts, you will have to modify the ProjectID in the following files:  
* [fileupload.go](fileupload.go#L41) 
* [deleteFirestoreCollection.go](deleteFirestoreCollection.go#L13)
* [search_index.go](./searchindex/search_index.go#L12)
* [map_wrapper.go](./mapper/map_wrapper.go#L13)
* [inverse_index_wrapper.go](./inverseindex/inverse_index_wrapper.go#L15)

The links should take you directly to the line where the ProjectID is defined, but if not they are all defined before the first function in each file. You'll also have to modify `mapFuncUrl` and `reduceFuncUrl` in [inverse_index.go](./inverseindex/inverse_index.go#L15). If we were running these functions locally, I could have made these changes part of the script, but I could not figure out how to pass the command-line arguments to Cloud Functions.  

The functions can be tested using Postman or your favorite method of making post requests. Alternatively, [demo.sh](demo.sh) creates an inverted index, assuming the corpus is in your project's Firestore, and then searches it. To do so, `searchFuncUrl` and `createFuncUrl` will need to be changed in [mapreduce_demo.go](mapreduce_demo.go#L13). As with the deployment script, you'll need to have the Firestore CLI and be signed in for it to work. If you've already built the inverted index and want to search for another word, [maprduce_demo.go](mapreduce_demo.go) can be ran with 
```console
go run mapreduce_demo.go -createindex=false -word=<search term>
``` 

# Google Cloud Functions
## RyoostMapHttp
RyoostMapHttp performs the map function on the files that are given to it. It uses the default 256 MB of memory.
### Function Cost
99% of calls take at most 6,551ms. With the execution price of $0.000000648/100ms, this function costs $0.000031 per execution. This function also reads from Firestore. After the first 20,000 document reads in a day, it costs $0.036 per 100,000 document reads. This function typically only makes one or two document reads, so this price is negligible.
### Function Signatures
RyoostMapHttp is triggered by post requests. The body of the request must contain the following input and the response body is the following output.  
Input:  
```json
{
"Files" : ["filename1", "filename2", ...]
}
```
Output:
```json
{
    "mapresult" : [
        {
            "wordfile" : "word1_filename1",
            "count" : 1
        },
        {
            "wordfile" : "word1_filename1",
            "count" : 1
        },
        {
            "wordfile" : "word2_filename1",
            "count" : 1
        }, ...
    ]
}
```
## RyoostReduceHttp
RyoostReduceHttp performs the reduce operation on the resulting tuples from RyoostMapHttp. It uses the default 256 MB of memory.
### Function Cost
99% of calls take at most 970ms. With the execution price of $0.000000463/100ms, this function costs $0.00000449 per execution.
### Function Signature
RyoostReduceHttp is triggered by post requests. The body of the request must contain the following input and the response body is the following output.  
Input:  
```json
{
    "mapresult" : [
        {
            "wordfile" : "word1_filename1",
            "count" : 1
        },
        {
            "wordfile" : "word1_filename1",
            "count" : 1
        },
        {
            "wordfile" : "word2_filename1",
            "count" : 1
        }, ...
    ]
}
```
Output:  
```json
{
    "reduceresult" : [
        {
            "wordfile" : "word1_filename1",
            "count" : 2
        },
        {
            "wordfile" : "word2_filename1",
            "count" : 1
        }, ...
    ]
}
```
## RyoostCreateIndexHttp
RyoostCreateIndexHttp creates an inverted index from the given files using the given number of parallel map and reduce functions. Because it assembles the inverted index from the outputs of the reduce of functions, it must store the entire inverted index in memory. With my corpus of 25 books, it requires 1 GB of memory.
### Function Cost
99% of calls take at most 1.32 minutes. With the execution price of $0.000001650/100ms, this function costs $0.00152 per execution for computation time. Additionally, this function writes to Firestore. After 20,000 document writes a day, it costs $0.108 per 100,000 documents. Creating an inverted index of my corpus of 25 books took 37,718 writes. Assuming the free quota has already been met, writing to Firestore costs $0.040735, giving a total per execution cost of $0.042253. This function also calls the RyoostMapHttp and RyoostReduceHttp functions and therefore would have additional costs for each of those depending on the number of mappers and reducers used.
### Function Signature
RyoostCreateIndexHttp is triggered by post requests. The body of the request must contain the following input and the response body is the following output.  
Input:  
```json
{
    "Files": ["filename1", "filename2", ...],
    "Mappers": <number of desired map processes>,
    "Reducers": <number of desired reduce processes>
}
```
Output:  
```json
{
    "status": "true" or "false"
}
```

## RyoostSearchHttp
RyoostSearchHttp searches the inverted index on Firestore for the given search term. The function uses the default 256 MB of memory.   
### Function Cost
99% of calls take at most 449.32ms. With the execution price of $0.000000463/100ms, this function costs $0.00000208 per execution for computation time. It also reads one document from Firestore, which has a negligible cost.
### Function Signature
RyoostSearchIndexHttp is triggered by post requests. The body of the request must contain the following input and the response body is the following output.
Input:  
```json
{
    "Term" : "<search term>"
}
```
Output:  
```json
{
    "results" : [
        {
            "title": "filename1",
            "count": 4345
        },
        {
            "title": "filename2",
            "count": 2347
        }, ...
    ]
}
```
# Implementation Details
Note: halfway through I realized I was using "inverse index" instead of "inverted index" and I didn't want to spend time refactoring since the assignment was already taking much longer than expected.
## Corpus
The corpus is made up of 25 of the top downloaded books from Project Gutenberg. Books that are approximately 1MB or larger couldn't fit into a single Firestore document and so they were not included. The corpus is composed of the following:  
"Frankenstein","Pride and Prejudice","The Legend of Sleepy Hollow","Alice's Adventures in Wonderland","Dracula","The Scarlet Letter","A Christmas Carol","The Adventures of Sherlock Holmes","The Yellow Wallpaper","The Picture of Dorian Gray","A Tale of Two Cities","The Strange Case of Dr. Jekyll And Mr. Hyde","The Great Gatsby","A Doll's House","A Modest Proposal","Metamorphosis","The Prince","Heart of Darkness","The Odyssey","Grimms' Fairy Tales","Beowulf","The Adventures of Tom Sawyer","Emma","The Communist Manifesto", and "Anthem".  
The corpus is stored in the ryoost-mapreduce collection on Firestore. Each book is a document whose id is the title of the book and has a field called "text" that stores the text of the book.
## Barrier and Intermediate Communication
[inverse_index.go](./inverseindex/inverse_index.go) partitions the corpus into *m* files, where *m* is the number of specified map instances, and then creates a go routine to call RyoostMapHttp for each of *m* groups of files. The barrier is implemented using channels, which is Go's method of inter go routine communication and essentially act as a queued future. *m* channels are created and RyoostCreateIndexHttp waits until it gets the results from all *m* channels. Once the results of all *m* map functions are receieved by RyoostCreateIndexHttp, it partitions them into *r* groups of word-file pairs where *r* is the number of specified reduce instances. Once the result of all reduce instances are recieved, the (word-file, count) pairs are grouped based on their words to build the inverted index and uploaded to Firestore 500 words at a time. While a normal MapReduce would have the intermediate results be sent directly from the map processes to the reduce processes, Google Cloud Functions lacks the required inter-process communication to coordinate this, so I store the intermediate results in memory of the function orchestrating the entire process.
## Storage
All cloud storage is done through Firestore in a collection called ryoost-mapreduce. The corpus storage is described above. The inverted index is stored as a subcollection "Words" in a document "inverseindex" of the main ryoost-mapreduce collection. Each word in the inverted index is it's own document whose fields are the names of the books and associated values the count of that word. While this means more writes when uploading the inverted index, it improves search performance and Firestore is optimized for many small documents instead of few large documents (also it wouldn't fit in the 1MB document size limit).