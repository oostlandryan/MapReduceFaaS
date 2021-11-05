# Google Cloud Functions
## RyoostMapHttp
Performs the map operation on the given files  
Input:  
```json
{
"Files" : ["filename1", "filename2", ...]
}
```
Output: `Inverted Index Built` or `Inverted Index Not Stored` depending on if the entire inverted index was uploaded to Firestore or not.
```json
{
    "mapresult" : [
        {
            "wordfile" : "word1_filename1"
            "count" : 1
        },
        {
            "wordfile" : "word1_filename1"
            "count" : 1
        },
        {
            "wordfile" : "word2_filename1"
            "count" : 1
        }, ...
    ]
}
```
## RyoostReduceHttp
Performs the reduce operation on the given key, value tuples  
Input:  
```json
{
    "mapresult" : [
        {
            "wordfile" : "word1_filename1"
            "count" : 1
        },
        {
            "wordfile" : "word1_filename1"
            "count" : 1
        },
        {
            "wordfile" : "word2_filename1"
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
            "wordfile" : "word1_filename1"
            "count" : 2
        },
        {
            "wordfile" : "word2_filename1"
            "count" : 1
        }, ...
    ]
}
```
## RyoostCreateIndexHttp
Creates the inverted index, storing it in Firestore.  
Input:  
```json
{
    "Files": ["filename1", "filename2", ...]
    "Mappers": <number of desired map processes>
    "Reducers": <number of desired reduce processes>
}
```
Output:  

## RyoostSearchHttp
Searches the inverted index for the given search term  
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
## Map
## Reduce
## Barrier
## Storage