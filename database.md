Bolt (https://github.com/boltdb/bolt) has been selected as the inital design database due to a complete Go implementation and
general "fit" for the project.

Archit has (2) main buckets, both stored in the same physical file (archit.db).

Bucket("Reputation") contains the connectivity and reputation information for everyone known.

Reputation struct {
	address string // KEY:  IMAC address of partner (renter or farmer),
	ip_Addr string // TCP/IP address to communicate over
	value int      // Actual reputation value
}

Bucket("File") contains storage information for a renters files

File struct {
	key struct {
		filename string
		timestamp string
	}
	map slice[string] struct {	// 1 - N (1)GB slices of file, where string is the hash of the 1GB slice
		segmentCount int		
		segment [40]struct {  
			address string  // IMAC address of farmer
			hmac string	// segment Hash value
		}
	}
}
