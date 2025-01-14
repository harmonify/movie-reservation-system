# Idempotency

## Client generated idempotency key

Client generates UUID and set it to HTTP header `Idempotency-Key` for some POST request, i.e. reservation process.

The server then will process the HTTP request and stores the HTTP response alongside idempotency key in a persistent storage.

If the server receives the same idempotency key, it returns the stored response.

Since POST requests is typically not read-heavy, there is no need for caching.

For some cases, i.e. reservation process, we could periodically clean up the records that are older than its movie showtime.
