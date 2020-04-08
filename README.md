HTTP API server with a single endpoint: `/fizzbuzz`

It only responds to POST requests with the correct content-type and the response is a JSON object with the following format: 

```json
{ 
  "response": "<response string>"
}
```

The FizzBuzz program supports 2 commands:
* `client` - Runs a few concurrent clients that continuously poll the fizzbuzz endpoint and prints the output
* `server` - Runs a FizzBuzzServer

It can be configured through the following environment variables:
* `FIZZBUZZ_PORT` - Port that the server listens on (default: `4343`)
* `FIZZBUZZ_REMOTE_ADDR` - Remote address that the FizzBuzz client connects to (default: `http://localhost:4343`)
