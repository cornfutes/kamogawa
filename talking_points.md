Talking points with YC


===============
Date: Sept 4 

- inquired initial impressions of Go
  - more sane library than C++ 
  - forces us to code one way 
  - single library ecosystem means we don't have debates about which library to use 
  - structs and slices were kind of unintuitive, but we're past that 
  - JSON in go is painful, however there is a dynamic library 
- david: full-text search might not make sense
  - trigram 
  - louis: product debate about suffix searching is all we need, for now
  - product idea to have different searching algorithm
  - we have answer to our FOSS story 
- caching
  - debated caching vs scheduling, how to keep fresh
- HTTP
  - when content-length not specified, Go defaults to chunked 
    Transfer-Encoding: chunked. This prevents using gzip middleware from working.
