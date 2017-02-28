Basic data flow documentation (presumes a Raptor value of 8):

Sending
1) User specifies a file of size N to be sent
2) File is read in and broken into (1GB-sum(headers)) byte SLICES at a time
3) File is Blake2 hashed to generate a unique filename for the SLICE
4) Read buffer is split into (32) 32mb-aesHeaders-fountainCodeHeaders PLAINTEXT BLOCKS
5) Each PLAINTEXT BLOCK is aes-crt encrypted, in parallel ("go encrypt()")
6) Systems waits for all (32) ENCRYPTED BLOCKS to be ready
7) Fountain code is applied to ENCRYPTED BLOCKS to generate 40 LT BLOCKs
8) In parallel for each block ("go sendblock()"):
	8A)  HMAC generated for each encryption block
	8B)  SLICE# * BlockNumber (0-39) used to index highest reputation known farmer
	8C)  ENCRYPTED BLOCK attempted to be sent to identified farmer as Blake2 filename
	8Ca)  If successful, payment initiated, farmer reputation boosted by 1
	8Cb)  If failed, farmer reputation decreased by 1, 40 added to index, new farmer selected.
	8Cc)  If farmer refused download due to duplicate filename, restore reputation (should never
		occur but deals with mathematical possibility of a collision)
	8D)  Original file id record (filename + timestamp) updated to reflect slice name, block number,
             and block HMAC

Notes:  Since the Blake2 hash will generate a unique name, so should not be a conflict storing on each	
        farmer.
	Utilizing farmer from 2 block not an issue since a difference SLICE name will be used on that
	path.
	Farmer/network table is large (several thousand entries)

Receiving
1) User specifies file id to be recovered
2) By slice:
   2A) Parallel request sent to each of 40 farmers ("go getblock()")
   2B) Each received block is HMAC hashed and compared to recorded HMAC
   2Ba)  If match, farmer gains 5 reputation points, sent payment
   2Bb)  If failed, farmer loses 5 reputation points, payment withheld (famer notified?)
   2C) If valid, Block entered into LT decode process
   2Ca)  If LT decode process indicates SLICE is complete, any remaining farmers who have not returned
         a block lose 1 reputation point for being slow and notified of such
   2Ca1)  If transmitted block is less than 10% downloaded, no payment is sent
   2Ca2)  If transmitted block is 10%+ downloaded, proportional payment is queued pending SLICE 
	  confirmation
3) Reconstructed SLICE is Blake2 hashed and results compared to filename for final confirmation
3a) Should final confirmation fail, recovery is attempted by finishing retreval of remaining blocks 
    (subject to reasonable timeouts), and alternative block reconstruction attempted.  If successful,
    and results are different, SLICE is again Blake2 hashed.  This continues (and is logged) until
    either success or alternatives are exhausted
4) File slice is written to disk and next slice started.

Notes:  Failure of more than 8 farmers will result in failure notification to user - a very bad thing
