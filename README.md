# pipeline
A Go implementation of a Unix-like pipeline (analog of cat emails.txt | SelectUsers | SelectMessages | CheckSpam | CombineResults). Uses goroutines and channels to pass data between pipeline stages. Implements RunPipeline, SelectUsers, SelectMessages, CheckSpam and CombineResults functions with proper parallelism and batching.  
# Notes  
Only spammer.go was implemented by me. All other files (common.go, main_test.go) were provided as part of the original task.

Task
The original task required implementing the following functions in spammer.go: https://gitlab.com/vk-golang/lectures/-/blob/master/02_async/99_hw/spammer/README.md?ref_type=heads
