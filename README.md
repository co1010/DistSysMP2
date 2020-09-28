# DistSysMP2

# To Run
First create a server. To do so, navigate to /DistSysMP2 and type
> go run . -s server -s [PORT]

Where [PORT] is the port you want to run the server on.
To shut down the server type EXIT and hit enter.

To create clients navigate to /DistSysMP2 and type
> go run . -s [IP] -s [PORT] -s [USERNAME]

Where [IP] and [PORT] are the ip and port of the server you want to connect to, and [USERNAME] is the username you wish to use.
To send a message to another user type
> send [USERNAME] [MESSAGE]

Where [USERNAME] is the username of the person you want to send the message to and [MESSAGE] is the message you want to send.

# System Overview
A central server handles all connections using goroutines. When a client connects they register their username with their connection. When the server receives a message from a client it looks up the username to see if it has a connection. When the server gets the exit command it messages all clients to close their connections, then closes the server. When clients get the exit command it messages the server which closes that single connection and removes the user from the register. The client then exits. The functionality of registration and exiting is built into the Message struct.

# The Message Struct
In addition to the obvious To, From, and Content fields in the Message struct there are two bool fields: Register and Exit. These fields act as flags so the system knows what to do. The Register field is true when a client wants to "register" its username with the server. The Exit field is true when either the client or the server wants to let the other know to close the connection. These fields have to be in the Message struct because when a gob encoded message is decoded it needs to be put into the correct type. For this reason I needed a struct that allows for normal message format but also holds additional information for actions. I found the best way to do this was to put boolean fields in the struct. One potential problem this implementation has is security. If a user was somehow able to set the Exit or Register fields on a normal message they could potentially cause issues with the server.

# Registering
Clients must "register" their usernames with the server. On the backend the server has a global map variable called connections that takes string keys and net.Conn values. This is necessary due to the implementation of the central server and multiple clients. The server does not know which connection is which username. Therefore, the map connects usernames to connections so the server knows where to send messages when given a username.

# Exiting
When a client receieves a message with the Exit flag it knows that the server is closing and therefore closes it's
