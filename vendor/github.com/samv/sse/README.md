# README for sse

This repository contains some code for Server-Sent Events (SSE).

The focus at first will be on providing an SSE client which is completely
compliant to the [spec](https://www.w3.org/TR/2009/WD-eventsource-20091029/),
and more-or-less API equivalent to the corresponding EventSource API available
in JavaScript.

The idea is that using this library as an SSE client in your code should result
in server event implementations which are also useful from JavaScript clients,
without excessive hacks.

Later, it may include functions for writing SSE server handlers which support
things like retry negotation, keep-alives and resuming events by ID.

## Current Status

Past proof-of-concept stage.

Client code: compliant protocol support is complete, minimal client in
progress.  From there, protocol features and specified API functionality like
connection closing, event listeners and callback APIs will come in, in order of
usefulness.

Server code: I have basic proof of concept code which is not in this repository
yet, and will probably feature in the client test mocks first before the server
API is designed and implemented.

See the TODO.org file for working notes.
