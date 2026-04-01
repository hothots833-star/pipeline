# pipeline
A Go implementation of a Unix-like pipeline (analog of cat emails.txt | SelectUsers | SelectMessages | CheckSpam | CombineResults). Uses goroutines and channels to pass data between pipeline stages. Implements RunPipeline, SelectUsers, SelectMessages, CheckSpam and CombineResults functions with proper parallelism and batching.
