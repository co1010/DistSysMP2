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
To exit the program simply type EXIT and hit enter.

# System Overview
A central server handles all connections using goroutines. When a client connects they register their username with their connection. When the server receives a message from a client it looks up the username to see if it has a connection. When the server gets the exit command it messages all clients to close their connections, then closes the server. When clients get the exit command it messages the server which closes that single connection and removes the user from the register. The client then exits. The functionality of registration and exiting is built into the Message struct.

# The Message Struct
In addition to the obvious To, From, and Content fields in the Message struct there are two bool fields: Register and Exit. These fields act as flags so the system knows what to do. The Register field is true when a client wants to "register" its username with the server. The Exit field is true when either the client or the server wants to let the other know to close the connection. These fields have to be in the Message struct because when a gob encoded message is decoded it needs to be put into the correct type. For this reason I needed a struct that allows for normal message format but also holds additional information for actions. I found the best way to do this was to put boolean fields in the struct. One potential problem this implementation has is security. If a user was somehow able to set the Exit or Register fields on a normal message they could potentially cause issues with the server.

# Registering
Clients "register" their usernames with the server after establishing a connection. On the backend the server has a global map variable called connections that takes string keys and net.Conn values. This is necessary due to the implementation of the central server and multiple clients. The server does not know which connection is which username. Therefore, the map connects usernames to connections so the server knows where to send messages when given a username. One potential problem is that maps are not safe for concurrent use accoring to the golang blog. I could have used some sort of synchronization mechanism like RWMutex which would prevent any simultaeneous reads/writes. While this may seem bad, in testing I did not run into any of these errors, which I think can only really happen in the system when multiple users are using it at the same time.

# Exiting
There are two types of exiting in this system, one from the client and one from the server. When a client exits, it first sends a message with the Exit field true to the server. When the server receives a message with the Exit field set to true it deletes the user that sent it from the connections map, closes the connection, then breaks the infinite for loop. Client side, after sending the message it closes the connection and exits. When the server exits it first closes all active connections by sending a message to each connection with the exit field set to true. When a client receieves a message with the Exit flag it knows that the server is closing and therefore closes it's connection, prints that the server is closing, and then exits the program. On the server side, after closing all client connections the server is closed and the program exits.

# Packages and Code Compartmentalization
After writing most of the code using main.go, server.go, and client.go in the main package I attempted to refactor my code to reflect the application and network layers. I hoped this would make it more readable, professional, and understandable. I made a client and server file in an application and network package, and moved functions around depending on whether they were part of the "application" layer or "network" layer. However, after attempting to run my code after the refactor I got the error "Import cycle not allowed". Go does not allow for cyclic imports, and in my refactor the application package imported the network package, and vice versa. This created an import cycle that caused the error. After some attempted debugging and thinking I reverted to just using the main package. Since information travels both into and out of the network layer I couldn't find a way to organize the packages such that imports were "one way". This project was simple enough to where the code is easily understandable all in one package, but for a larger project it might be helpful to compartmentalize the code in separate packages while avoiding creating a cycle.

# Reading Input
To avoid bottlenecks with command line inputs I made a utility function meant to be called as a goroutine that puts input into a string channel. This allows the utility function to keep reading input while another function processes it, which helps reduce bottlenecking.

# Future Considerations and Concerns
Some concerns I have with this code are security, safety, and privacy. The server decodes all messages and reads the To field of the message. If a malicious person got into the server they would be able to see and read all messages coming through it. Further, the Register and Exit flags are attached to the Message struct. If a client was somehow able to edit these fields in a message they could theoretically kick anyone they wanted offline by sending a message from them with the Exit flag as true. A future consideration would be the compartmentalization of code that I talked about two paragraphs above. One benefit to this code design is that it can very easily be made into a chat room rather than a provate messaging system. The "To" field could be removed and the server could just multicast each message to all connected clients. Multicast is already more or less implemented for the server exit functionality.







