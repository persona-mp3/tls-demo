### Project
This is a simple tcp server that can be configured with or without using TLS. 
There's both a `client.go` and `server.go` file.
```bash
go run server.go

# In another terminal
go run client.go
```

Now, any message you type on the client side will respond back to you with 
a simple echo message. This was done to demonstrate how unencrypted messages 
over a network are a vulnerabilty and how any one with a computer can see those 
messages. 


## To run a TLS Server
You need to first generate tls certificates that the TCP Server 
would use, or you could use the defaults provided
```bash
chmod u+x gen-certs.sh # changes user permissions to executable

./gen-certs.sh
```

Now you can run the TCP Server over TLS
```bash
go run server.go -secure 
```

And the client over TLS
```bash
go run client.go --secure 
```

Now every data sent over this connection is encrypted using 
TLS. Upon first connection, the server and client negotiatie 
over using TLS. The client is configured to use v1.3, and then
they exchange keys, and complete the handshake.

During the DEMO, I was able to show the state of this data 
using `Wireshark`, to analyse tcp traffic on the computer


### Installation
To be able to use this, you need to have [Go](https://go.dev) compiler installed 
on your local machine.


