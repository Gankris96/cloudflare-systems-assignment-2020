# Cloudflare Systems Assignment Submission

## Setup
1. Make sure you have GO installed. https://golang.org/
2. Run the command `go run profiler.go --url <url> --profile <int>`
3. Here is the usage information
  ```
    go run profiler.go <options> 
    Here are the available options
     -help 
        	URL to make the request to
      -profile int
        	URL to make the request to
      -url string
        	URL to make the request to (default "https://cf-assgn.gramadur.workers.dev/links")

  ```
## Screenshots
1. Execution screenshots can be found in screenshots folder

## Metrics 

 ```
| Url                   | Request Count | Fastest Time | Slowest Time | Mean      | Median    | Success Percentage | Smallest Response | Largest Response |
|-----------------------|---------------|--------------|--------------|-----------|-----------|--------------------|-------------------|------------------|
| My cloudflare website | 30            | 27.019 ms    | 103.51 ms    | 49.404 ms | 48.380 ms | 100%               | 2610  bytes       | 2758  bytes      |
| facebook.com          | 30            | 59.34 ms     | 350.78 ms    | 81.26 ms  | 67.73 ms  | 100%               | 369 bytes         | 369 bytes        |
| nytimes.com           | 30            | 19.037 ms    | 58.29 ms     | 26.313 ms | 24.50 ms  | 100%               | 462  bytes        | 462  bytes       |
| reddit.com            | 30            | 27.386 ms    | 467.98 ms    | 82.34 ms  | 38.38 ms  | 100%               | 256 bytes         | 256 bytes        |
| github.com            | 30            | 19.71 ms     | 712.70 ms    | 72.38 ms  | 26.43 ms  | 100%               | 107 bytes         | 107 bytes        |
```
1. My cloudflare website performs on par with other websites despite the fact that it has a significantly larger response size in bytes.
2. My General Programming assignment can be found here - https://cf-assgn.gramadur.workers.dev