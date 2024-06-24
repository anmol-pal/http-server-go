# http-server-go

USAGE
```
./your_server.sh
```

### HTTP response
An HTTP response is made up of three parts, each separated by a CRLF (\r\n):
* Status line.
* Zero or more headers, each ending with a CRLF.
* Optional response body.
```
HTTP/1.1 200 OK\r\n\r\n
```

```
// Status line
HTTP/1.1  // HTTP version
200       // Status code
OK        // Optional reason phrase
\r\n      // CRLF that marks the end of the status line

// Headers (empty)
\r\n      // CRLF that marks the end of the headers

// Response body (empty)
```

> Tests
```
curl -v http://localhost:4221
```


### Extract URL path
HTTP Request : ````GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n```

Here's a breakdown of the request:

```
// Request line
GET                          // HTTP method
/index.html                  // Request target
HTTP/1.1                     // HTTP version
\r\n                         // CRLF that marks the end of the request line

// Headers
Host: localhost:4221\r\n     // Header that specifies the server's host and port
User-Agent: curl/7.64.1\r\n  // Header that describes the client's user agent
Accept: */*\r\n              // Header that specifies which media types the client can accept
\r\n                         // CRLF that marks the end of the headers

// Request body (empty)
```
> Tests

HTTP/1.1 404 Not Found\r\n\r\n
```
curl -v http://localhost:4221/abcdefg

```

HTTP/1.1 200 OK\r\n\r\n
```
curl -v http://localhost:4221
```


### Respond with body

Your /echo/{str} endpoint must return a 200 response, with the response body set to given string, and with a Content-Type and Content-Length header.

Here's an example of an /echo/{str} request:

```GET /echo/abc HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n```

And here's the expected response:

```HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 3\r\n\r\nabc```

Here's a breakdown of the response:
```
// Status line
HTTP/1.1 200 OK
\r\n                          // CRLF that marks the end of the status line

// Headers
Content-Type: text/plain\r\n  // Header that specifies the format of the response body
Content-Length: 3\r\n         // Header that specifies the size of the response body, in bytes
\r\n                          // CRLF that marks the end of the headers

// Response body
abc      
```

> Tests
```
curl -v http://localhost:4221/echo/abc
Response: HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 3\r\n\r\nabc
```

### Read header
Your /user-agent endpoint must read the User-Agent header, and return it in your response body. Here's an example of a /user-agent request:

```
// Request line
GET
/user-agent
HTTP/1.1
\r\n

// Headers
Host: localhost:4221\r\n
User-Agent: foobar/1.2.3\r\n  // Read this value
Accept: */*\r\n
\r\n

// Request body (empty)
```

Response
```
// Status line
HTTP/1.1 200 OK               // Status code must be 200
\r\n

// Headers
Content-Type: text/plain\r\n
Content-Length: 12\r\n
\r\n

// Response body
foobar/1.2.3                  // The value of `User-Agent`
```

> Tests

```curl -v --header "User-Agent: foobar/1.2.3" http://localhost:4221/user-agent```
Your server must respond with a 200 response that contains the following parts:

Content-Type header set to text/plain.
Content-Length header set to the length of the User-Agent value.
Message body set to the User-Agent value.

```HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nfoobar/1.2.3```

### Return a file
/files/{filename} endpoint, which returns a requested file to the client.
usage : ```./your_server.sh --directory /tmp/```

The tester will then send two GET requests to the /files/{filename} endpoint on your server.

* First request
The first request will ask for a file that exists in the files directory:

```
$ echo -n 'Hello, World!' > /tmp/foo
$ curl -i http://localhost:4221/files/foo
```
Your server must respond with a 200 response that contains the following parts:

Content-Type header set to application/octet-stream.
Content-Length header set to the size of the file, in bytes.
Response body set to the file contents.
```HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: 14\r\n\r\nHello, World!```

* Second request
The second request will ask for a file that doesn't exist in the files directory:

```$ curl -i http://localhost:4221/files/non_existant_file```
Your server must respond with a 404 response:

```HTTP/1.1 404 Not Found\r\n\r\n```


### Read request body
Request body
A request body is used to send data from the client to the server.

Here's an example of a POST /files/{filename} request:
```
// Request line
POST /files/number HTTP/1.1
\r\n

// Headers
Host: localhost:4221\r\n
User-Agent: curl/7.64.1\r\n
Accept: */*\r\n
Content-Type: application/octet-stream  // Header that specifies the format of the request body
Content-Length: 5\r\n                   // Header that specifies the size of the request body, in bytes
\r\n

// Request Body
12345
```

> Tests
The tester will execute your program with a ```--directory``` flag. The --directory flag specifies the directory to create the file in, as an absolute path.

```$ ./your_server.sh --directory /tmp/```

The tester will then send a POST request to the /files/{filename} endpoint on your server, with the following parts:

Content-Type header set to application/octet-stream.
Content-Length header set to the size of the request body, in bytes.
Request body set to some random text.

```$ curl -v --data "12345" -H "Content-Type: application/octet-stream" http://localhost:4221/files/file_123```
Your server must return a 201 response:

```HTTP/1.1 201 Created\r\n\r\n```
Your server must also create a new file in the files directory, with the following requirements:

The filename must equal the filename parameter in the endpoint.
The file must contain the contents of the request body.

### Compression headers

support for the Accept-Encoding and Content-Encoding headers.
```Accept-Encoding and Content-Encoding```

An HTTP client uses the Accept-Encoding header to specify the compression schemes it supports. In the following example, the client specifies that it supports the gzip compression scheme:
```
> GET /echo/foo HTTP/1.1
> Host: localhost:4221
> User-Agent: curl/7.81.0
> Accept: */*
> Accept-Encoding: gzip  // Client specifies it supports the gzip compression scheme.
```
The server then chooses one of the compression schemes listed in Accept-Encoding and compresses the response body with it.

Then, the server sends a response with the compressed body and a Content-Encoding header. Content-Encoding specifies the compression scheme that was used.

In the following example, the response body is compressed with gzip:

```
< HTTP/1.1 200 OK
< Content-Encoding: gzip    // Server specifies that the response body is compressed with gzip.
< Content-Type: text/plain  // Original media type of the body.
< Content-Length: 23        // Size of the compressed body.
< ...                       // Compressed body.
```
If the server doesn't support any of the compression schemes specified by the client, then it will not compress the response body. Instead, it will send a standard response and omit the Content-Encoding header.

For this extension, assume that your server only supports the gzip compression scheme.

For this stage, you don't need to compress the body. You'll implement compression in a later stage.

> Tests

First request
First, the tester will send a request with this header: Accept-Encoding: gzip.

```$ curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/abc```
Your server's response must contain this header: Content-Encoding: gzip.

```
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Encoding: gzip

...  // Body omitted.
```
Second request
Next, the tester will send a request with this header: Accept-Encoding: invalid-encoding.

```$ curl -v -H "Accept-Encoding: invalid-encoding" http://localhost:4221/echo/abc```
Your server's response must not contain a Content-Encoding header:

```
HTTP/1.1 200 OK
Content-Type: text/plain

...  // Body omitted.
```
Notes
You'll add support for Accept-Encoding headers with multiple compression schemes in a later stage.
There's another method for HTTP compression that uses the TE and Transfer-Encoding headers. We won't cover that method in this extension.

### Gzip compression

The request will contain an Accept-Encoding header that includes gzip.

```$ curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/abc | hexdump -C```
Your server's response must contain the following:

```
200 response code.
Content-Type header set to text/plain.
Content-Encoding header set to gzip.
Content-Length header set to the size of the compressed body.
Response body set to the gzip-compressed str parameter.
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Type: text/plain
Content-Length: 23

1F 8B 08 00 00 00 00 00  // Hexadecimal representation of the response body
00 03 4B 4C 4A 06 00 C2
41 24 35 03 00 00 00
```
