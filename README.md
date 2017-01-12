# archit
Distributed archival environment

Initial concept creation as of 12/24/2016

Archit will be a PC based distributed environment whereby users can encrypt and archive objects to a shared infrastructure.  This differs
from current commercial implementations in that "farmers" will rent space on their PCs and "renters" will store their files distributed 
across a number of networked computers.

Initial coin usage will be via IMACREDIT (https://github.com/imacredit/imacredit).  Coin usage is required as a method of resolving 
deliberate sabotage of the network by request flooding, since each request will have a small transaction charge.

This project was inspired by the kind folks over at Storj (https://github.com/storj).  It differs in a few key areas:

1)  The project will be written primarily, if not exclusively, in the Go programming language in order to provide speed, portability, and a
relatively bug free environment (less chances for memory leaks and assocaited environmental issues).

2)  Go's inherrent ability to both multi-thread and communicate over the Internet will be heavily utilized, with an emphasis on IP 
connectivity between distinct processes.  For instance, the environment should be able to split off network management from other core 
features - to the point that such a service could run on a seperate machines if a farmer desired to do so.  Likewise, there should be 
little benefit to running more than one instance of archit - that instance should take full advantage of the server its running on behave
if restrictions are specified via confirguration (for instance, to only use 4 cores of an 8 core system).

3)  Self-tuning will be inherent in the code design, within the configuration limits.  e.g. issues like concurrency will be automatically
adjusted over time to attempt to best utilize offered resources.

4)  File stored in the archit ecosystem will have metadata search capabilities.  Files will be encrypted and broken into shards for storage 
across multiple systems.  Files themselves will either be public or private, based on the renters desire.  This is turn would allow projects
like the Climate Data preservation initiative to pay for storage that the world could access.

5)  Initial effort will be to create a linux based CLI.  However, since the project will be written in Go, a Windows port should be 
relatively easy to provide.  Support for a GUI interface is not intended until after Alpha testing, if then.

6)  This project will be invitation only/private until Beta release, when limited public support will be offered.

7)  This code is dependant upon one modification to the standard github.com/jessevdk/go-flags/option.go file.  Issue number #213 addresses this, but until that is mainstreamed, add:

// IsSetDefault returns true if option has been set via the default tag
func (option *Option) IsSetDefault() bool {
return option.isSetDefault
}

to option.go.  Recommend placing it immediately under IsSet().
