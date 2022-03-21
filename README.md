Solution to a test task implementing a web service that combines two existing web services.

### Building
`$ go build .`

### Running
`$ ./test1 [[addr]:port]`
If ran with no arguments, service web server starts at :8080

### Further improvements
- Random name API has a request rate limit, which acts as a bottleneck. Consider caching previous results and picking from cached values in case of API failure.
- Add configuration options (address, timeouts etc.) via environment variables.
- Improve error handling according to business constraints.
- Add tests.
