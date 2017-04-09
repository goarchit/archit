# archit
Distributed archival environment

Initial concept creation as of 12/24/2016

Archit will be a PC based distributed environment whereby users can encrypt and archive objects to 
a shared infrastructure.  This differs from current commercial implementations in that "farmers" 
will rent space on their PCs and "renters" will store their files distributed across a number of 
networked computers.

Initial coin usage will be via IMACREDIT (https://github.com/imacredit/imacredit).  Coin usage is 
required as a method of resolving deliberate sabotage of the network by request flooding, since 
each request will have a small transaction charge.

This project was inspired by the kind folks over at Storj (https://github.com/storj).  It differs 
in a few key areas:

1)  The project will be written primarily, if not exclusively, in the Go programming language in 
order to provide speed, portability, and a relatively bug free environment (less chances for memory
 leaks and assocaited environmental issues).

2)  Go's inherrent ability to both multi-thread and communicate over the Internet is heavily 
utilized, with an emphasis on IP connectivity between distinct processes.  For instance, the 
environment splits off network management from other core features - to the point that such a 
service could run on a seperate machines if that feature was enabled and a farmer desired to do so.
Likewise, there is little benefit to running more than one instance of archit - that instance 
should take full advantage of the server its running unless restrictions are specified via 
confirguration (for instance, to only use 4 cores of an 8 core system via the GOMAXPROCS variable).

3)  Self-tuning is an inherent code design goal, within the configuration limits.  e.g. issues like
 concurrency will be automatically adjusted over time to attempt to best utilize offered resources.

4)  Files stored in the archit ecosystem will be sliced, encrypted, and broken into shards for 
storage across multiple systems.  Shards will always be private and fetchable only from the renters
 machines (where the tracking database resides) when combined with a PIN (which is normally 
prompted for, but can be specified in a configuration file).  Archit itself will have renter side 
checks that prevent more than one shard per slide from existing on any one farmer, and will strive 
to distribute each shard to a seperate farmer.

5)  Raptor like coding techniques will be utilized to provide file recovery redundancy across 
farmers.

6)  Initial effort will be to create a linux based CLI.  However, since the project will be written
 in Go, a Windows port should be relatively easy to provide.  Support for a GUI interface is not 
intended until after Alpha testing, if then.

7)  This project will be invitation only/private until Beta release, when limited public support 
will be offered.

8)  This code is dependant upon one modification to the standard 
github.com/btcsuite/btcutil/amount.go file:  
    Lines 36->44 should be changed from "BTC" to "IMAC"
