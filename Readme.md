Run go run Server/server.go in one terminal

In the other terminal, 

## Command to make Get Request

# For getting status for all the urls in memory

    curl -X GET "http://localhost:3421/getStatus"

# For getting status for requested urls

    curl -X GET -H "Content-Type: application/json" -d '{"website": <List of URLS Requested>}' "http://localhost:3421/getStatus"


## Command to make Post Request

# Adding URLS to the Memory

    curl -X POST -H "Content-Type: application/json" -d '{"website":<List of URLS added>}' "http://localhost:3421/addToUrlList"