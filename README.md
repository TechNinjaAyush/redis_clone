# redis_clone
Redis clone using golang
## Some  important things about redis 

1.Redis is an in-memory database, meaning it stores all data in RAM (Random Access Memory). This allows for extremely fast data access compared to disk-based databases.

2.Clients connect to the Redis server via TCP, which is a reliable protocol that uses a three-way handshake mechanism to establish connections.

3.Redis is commonly used for:

1.Caching mechanisms

2.Session storage (e.g., cookies)

3.Real-time data streams



##  Why redis is so Fast?

1.All Redis operations are performed in-memory, avoiding the latency of disk I/O.

2.Redis is single-threaded for command execution, which removes the complexity and overhead of multithreaded locking and synchronization.

3.Redis uses I/O multiplexing, meaning:

1.A single thread watches all client sockets.

2.It uses system calls like epoll (Linux) or kqueue (macOS/BSD) to detect which sockets are ready (i.e., have incoming data).

3.If a socket has data, Redis reads, processes the command, and sends the response immediately — all in one thread.

4. pls contact







