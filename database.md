Bolt (https://github.com/boltdb/bolt) has been selected as the inital design database due 
to a complete Go implementation and general "fit" for the project.

Archit has (2) main database:  Fileinfo.bolt and Peerinfo.bolt

FileInfo.bolt is only used R/W client side, and  contains storage information for a 
renters files.  Farming nodes MAY open this database R/O in order to spot check farmers
as part of reputation maintenance

PeerInfo.bolt contains the connectivity and reputation information for everyone known.
It is opened R/W by the farming instance.  Clients RPC request to the farming instance
in order to get peering information and to adjust a peers reputation.
