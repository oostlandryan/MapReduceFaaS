# Running Instructions
## Webpage on my project
I created a basic front-end running on Google's Cloud Run, which is accessible [here](https://wordsearch-ojbmp4ml5a-uc.a.run.app). The webpage queries the RyoostSearchHttp function and RyoostCreateIndexHttp function. It required a bit of extra work to set up and it's extra credit for the assignment, so I did not include it in the deployment script. However, the files themselves are still included with the assignment.
## Deploying to a new project
Running [deploy.sh](deploy.sh) will deploy all of the needed cloud functions and upload the corpus from Project Gutenberg to Firestore. The Firestore CLI might be necessary and you'll have to sign into it separately from gcloud.  
Prior to running the deploy script, you will have to modify the ProjectID in the following files:  
* [fileupload.go](fileupload.go#L41) 
* [search_index.go](./searchindex/search_index.go#L12)
* [map_wrapper.go](./mapper/map_wrapper.go#L13)
* [inverse_index_wrapper.go](./inverseindex/inverse_index_wrapper.go#L15)

The links should take you directly to the line where the ProjectID is defined, but if not they are all defined before the first function in each file. You'll also have to modify `mapFuncUrl` and `reduceFuncUrl` in [inverse_index.go](./inverseindex/inverse_index.go#L15). If we were running these functions locally, I could have made these changes part of the script, but I could not figure out how to pass the command-line arguments to Cloud Functions.  

The functions can be tested using Postman or your favorite method of making post requests. Alternatively, [demo.sh](demo.sh) creates an inverted index, assuming the corpus is in your project's Firestore, and then searches it. To do so, `searchFuncUrl` and `createFuncUrl` will need to be changed in [mapreduce_demo.go](mapreduce_demo.go#L13). As with the deployment script, you'll need to have the Firestore CLI and be signed in for it to work.

# Google Cloud Functions
## RyoostMapHttp
Performs the map operation on the given files  
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
Performs the reduce operation on the given key, value tuples  
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
Creates the inverted index, storing it in Firestore.  
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
Searches the inverted index for the given search term   
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
## Corpus
The corpus is made up of 25 of the top downloaded books from Project Gutenberg. Books that are approximately 1MB or larger couldn't fit into a single Firestore document and so they were not included. The corpus is composed of the following:  
"Frankenstein","Pride and Prejudice","The Legend of Sleepy Hollow","Alice's Adventures in Wonderland","Dracula","The Scarlet Letter","A Christmas Carol","The Adventures of Sherlock Holmes","The Yellow Wallpaper","The Picture of Dorian Gray","A Tale of Two Cities","The Strange Case of Dr. Jekyll And Mr. Hyde","The Great Gatsby","A Doll's House","A Modest Proposal","Metamorphosis","The Prince","Heart of Darkness","The Odyssey","Grimms' Fairy Tales","Beowulf","The Adventures of Tom Sawyer","Emma","The Communist Manifesto", and "Anthem".  
The corpus is stored in the ryoost-mapreduce collection on Firestore. Each book is a document whose id is the title of the book and has a field called "text" that stores the text of the book.
## Map
## Reduce
## Barrier and Intermediate Communication
## Storage