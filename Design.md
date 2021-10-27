# Google Cloud Functions
The map and reduces functions are implemented as Google Cloud Functions and are therefore deployed through gcloud.
## Map Function
Input: List of files in the firestore db to be processed: \[file1, file2, ...\]  
Output: List of tuples: \[(word1_filename1, 1), (word2_filename1, 1), ...\] for each time \<word\> appears in \<filename\>.  
## Reduce Function
Input: List of tuples: \[(word1_filename1, 1), (word2_filename1, 1), ...\]  
Output: List of tuples: \[(word1_filename1, # of occurences), (word2_filename1, # of occurences), ...\]  
# Orchestrator
The orchestrator acts as the main hub of this assignment. It determines what documents are added to the inverted index, how many mapper functions are deployed, and how many reducer functions are deployed. The orchestrator also acts as the barrier between the mappers and reducers. Finally, the orchestrator builds the generated inverted index from the reducer functions' output and acts as the interface for searching it.
# Database
The Project Gutenberg files are stored in a single collection on Google's Firestore. Each file is stored as a single document where the document name is the file name and the one key-value in the document is "text":\<contents of text file\>. While not strictly a key-value database, this assignment can easily be done within the free quota of Firestore while using Redis on Google Cloud does not have a free tier.
# Project Gutenberg Upload Script
