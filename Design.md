# Google Cloud Functions
## RyoostMapHttp
Performs the map operation on the given files  
Input:  
```
{
"Files" : ["filename1", "filename2", ...]
}
```
Output:  
```
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
```
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
```
{
    "reduceresult" : [
        {
            "wordfile" : "word1_filename1"
            "count" : 2
        },
        {
            "wordfile" : "word2_filename1"
            "count" : 1
        },
        {
            "wordfile" : "word3_filename2"
            "count" : 3
        }, ...
    ]
}
```
## RyoostCreateIndexHttp
Creates the inverted index, storing it in Firestore.
## RyoostSearchHttp
Searches the inverted index for the given search term  
Input:  
```
{
    "Term" : "<search term>"
}
```
Output:  
```
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
