---
title: "My adaptation of the \"A Better Web Server\" project to Go from \"Beej's Guide to Network Concepts\""
date: 2025-12-22
tags: ["go", "network"]
draft: true
---

<!--more-->

{{< toc >}}

{{< note title="Author note" >}}

Most of the text is copied from the [original article](https://beej.us/guide/bgnet0/html/split/project-a-better-web-server.html).
Only the examples have been changed to be Go-specific.

{{< /note >}}

## 9 Project: A Better Web Server

Time to improve the web server so that it serves actual files!

We’re going to make it so that when a web client (in this case we’ll use a browser) requests a specific file, the webserver will return that file.

There are some interesting details to be found along the way.

### 9.1 Restrictions

In order to better understand the sockets API at a lower level, this project may not use any of the following:

- Anything from the [`net/http`](https://pkg.go.dev/net/http) package.
- Any third-party HTTP client/server libraries.

You should use only:

- [`net`](https://pkg.go.dev/net) (for `net.Dial`, `net.Listen`, `net.Conn`, etc.)
- `bufio`, `io`, `os`, `fmt`, `bytes`, etc., as needed.

The idea is: after you finish this, Go’s `net/http` will look like a thin, friendly layer over what you just did yourself.

### 9.2 Running the Server

Just like in the [previous project](../2025-12-18-beej-project-http-client-server/#55-the-server), the server should start listening on port 28333 unless the user specifies a port on the command line. E.g.:

```sh
$ go run webserver.go       # Listens on port 28333
$ go run webserver.go 3490  # Listens on port 3490
```

### 9.3 The Process

If you go to your browser and enter a URL like this (substituting the port number of your running server):

```
http://localhost:33490/file2.html
```

The client will send a request to your server that looks like this:

```
GET /file2.html HTTP/1.1
Host: localhost
Connection: close
```

Notice the file name is right there in the GET request on the first line!

Your server will:

1. Parse that request header to get the file name.
2. Strip the path off for security reasons.
3. Read the data from the named file.
4. Determine the type of data in the file, HTML or text.
5. Build an HTTP response packet with the file data in the payload.
6. Send that HTTP response back to the client.

The response will look like this example file:

```
HTTP/1.1 200 OK
Content-Type: text/html
Content-Length: 373
Connection: close

<!DOCtype html>

<html>
<head>
...
```

[The rest of the HTML file has been truncated in this example.]

At this point, the browser should display the file.

Notice a couple things in the header that need to be computed: the Content-Type will be set according to the type of data in the file being served, and the Content-Length will be set to the length in bytes of that data.

We’re going to want to be able to display at least two different types of files: HTML and text files.

### 9.4 Parsing the Request Header
