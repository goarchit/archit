Q1)  What is this project about?
A1)  Read the readme.MD file!

Q2)  What coin is going to be used for payment?
A2)  IMACredit  (https://github.com/imacredit/imacredit - https://www.imacredit.org)

Q3)  Do I need an IMACredit wallet?
A3)  As of 1/2/2017, Yes, since no exchange list IMACredit.  Once one does, you could use an
     exchange address if they permit it, but we still recommend a local wallet.

Q4)  Why IMACredit?
A4)  #1 Developed by the same team 2+ years ago.  #2 Stable.  #3 ASIC resistant so little fear
     of a hold & releae (where a scoundral leverages a massive ASIC pool to overwhelm the coin,
     blocking normal users from finding a block for a month, then releasing all the locally
     found blocks while maintaining over half the network hash - just a variation on the 51%
     problem).  #4 A blockchain several years old that would be hard to hack from the 
     beginning. #5 Not currently on an exchange, which provides some flexibility come test time

Q5)  Why GoLang? (eg. Why write in Go)
A5)  Easy to read compiled language (more efficient than say something based on node.js) with
     built-in networking, multitasking, and inter-process communication capabilities.

Q6)  Why not just use Storj.IO?
A6)  Read the readme.MD file!

Q7)  Why did this start out in a private repository on Github?
A7)  Authors wanted to have something working before releasing to the world

Q8)  Does my wallet need to be on the same machine?  
A8)  No.  By default the code will look on the local machine, but you can override with a
     dns name or IP address in the configuration file.  Likewise for IMACredits default RPC
     port (64096)

Q9)  Do I need to use port 1958 for archit?
A9)  No, it is overrideable in the configuration file.

Q10) Is there a chat board for this project?
A10) Yes, currently hosted at http://farmbase.carpenter-farms.us:3000, plans are to host it
     at http://chat.archit.it eventually (domain pending on project progress)

Q11) What is this "Raptor" setting?
A11) Rator (short for Rapid Tornado) is a class of fountain codes that is used by Archit to provide
     storage redundancy.  Basically a file is broken into ~1GB slices, which are then further broken
     into (32) ~32GB shards.   Encodeing is done such that there is a >99% chance the file could be 
     decoded using just those (32) shards.  Extra shards are generated to improve that possibility of
     a successful decode and reconstruction.  With just (2) extra shards, the chance of failure jumps 
     from <1% to less than one-in-a-million - presuming those 32+2 shards are fetchable.  Archit defaults
     to (8) additional shards in order to provide protection from a farmer being offline at the time of
     file recreation and/or farmer file corruption.  Note that renters pay primarily for uploaded shards,
     so increasing the Raptor value does add some expense.

Q12) How am I protected from scamming farmers?  
A12) Inherrent in every interaction between a renter and a farmer is a reputation system.  This is 
     primarily designed to protect renters, but there is a design goal to make this a two way system that
     would protect farmers from unreasonable renters (such as excessive recalls of shards).
     Please see the file "dataflow.md" for gritty details.

Q13) How do I limit the number of cores Archit will use?
A13) Simply set the GOMAXPROCS environment variable.  There was little point in providing this as an 
     option flag since GoLang has it built in.
