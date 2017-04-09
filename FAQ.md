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
     into (32) ~32MB shards.   Encodeing is done such that there is a >99% chance the file could be 
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

Q14) Can I run Archit on a Raspberry Pi
A14) Yes, but not really.  Raspberry Pis and similar computers simply do not have the memory available
     to run the code efficiently.  We recommend at LEAST 2GB of memory.  When we say "yes" we are being
     technically correct, presuming you have plenty of swap space, but the program is going to run 
     hundreds of times slower than on a machine with that much memory.

Q15) Why is UPnP not supported?
A15) This was a tough one, but the bottom line is your dealing with cyber-currency, and you should 
     always have control over any environment with a wallet in it.  If Archit was to support UPnP,
     folks would demand IMAC support UPnP, which we won't do for the users benefit.  Way to likely
     someone would publish, and others would use, a common configuration file - which could then be
     easily hacked and give IMAC a bad reputation.  Of course, even without UPnP, folks can act foolishly,
     but we believe those that can manage their modems are likely to manage their security.

     Yes, we realize this greatly reduces the number of people who might use the software, but 
     security comes first.

Q16) Can a farmer scam me of coins?
A16) Yes, but...
     A fundemental goal of the project is to make scamming as hard and unprofitable as possible.  For
     instance, farmers are only provided a filename and a data shard to store.  All details around
     which block that shard represents, keys used, who storing the shard, reputation of the farmer, 
     etc. are stored exclusively by the renter.

     Tools are planned that will allow a "trust relationship" to be built, with some responsibility
     on the renter to behave responsibly.  Some tools will be automatic, such as random checks of
     known farmers to assure they are pingable.  Other tools will be created to allow a farmer to
     do some "data mining" of the reputation scores - such as the number of postively rated farmers
     vs. negatively rated, and allow them set minimum and maximum scores to work with.  

Q17) How can a farmer scam me?
A17) Since Archit is open source, a farmer could modify the source with the implicit goal of scamming
     renters.  For instance: they could store all data in /dev/null, destorying it the moment it
     is written, in hopes of collecting the fees and never actually having to store anything.  
     Fighting this type of behavior is the balance of renter responsbility and code checks. Good
     renters will not upload massive amounts of critical data to totally unknown farmers.  Instead
     they should join the network, and perhaps upload some short term data in order to start
     building a farmer reputation matrix - and gradually favor those farmers that have proven 
     themselves.  Likewise, a good renter will spend some trivial amount of coins to audit farmers,
     randomly pulling shards and checking their authentication codes.  Code checks will be
     implemented to make it troublesome to run multiple farming nodes on the same server.  Alas,
     depending on the verosity of the scammer farmer, any farm side checks could be disabled,
     which leads us back to renter caution.  Code checks will be added when doing so will likely
     prevent classes of non-coding scammers from being successful, which should help.

Q18) Why would I not ALWAYS want to send my shards to the highest rated farmers?
A18) It is in a renters best interest to spread their data across as many reasonable farmers as 
     possible.  Doing so reduces the risk of having to recover multiple files due to a single (the
     farmer with the best reputation?) going offline.
Q19) Why does Archit generate certification files?
A19) This is a one time thing and required if your public IP ever changes.
