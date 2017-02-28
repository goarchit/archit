This document describes, at a high level, the theory behind Archit's replication approach.

Back to Highschool!

You may recall from Highschool math class that in order to solve a system of N variables you
need N independant equations.  So:

A=B+C
B=C+D
D=5
B=D-A

Should be solvable.  In fact, the process was something like:

A=B+C
B=C+5
B=5-A
therefore:
C+5=5-A
C=-A
therefore:
2A=B
2A=5-A
3A=5
A=5/3, B=10/3, C=-5/3, D=5
Cross checking
5/3 = 10/3 + -5/3 (A=B+C)
10/3 = -5/3 + 15/3 (B=C+D)
10/3 = 15/3 - 5/3 (B=D-A)

In exactly the same way, we can solve the values for any 6 variables given 6 independant equations.

Onward to Archit!

Archit slices a file into ~1 GB chucks, and subdivides that ~1GB into (32) ~32MB shards.  Why
approximately 1GB and approximately 32MB?  There is a few bytes of header information in each shard,
but more on that in a bit.

The GB chuck is hashed to generate a 512bit hash (64 byte) which is used as a filename on farmer 
machines.  Since no farmer machine is allowed to own more than one shard, there should never be a 
filenanme conflict.  If there is, that farmer indicates the problem and a different farmer is 
selected (should be a very rare event)

Each of the ~32MB shards are encrypted in memory and then XORed with R-1 adjoining shards.  The 
result is then hashed for authentication purposes. R, which is called the "Raptor factor", is the 
number of additional blocks that will be generated.  By default R is 8, which will be used in this 
example.

In fact, Archit does not use a Raptor code, but was inspired by fountain codes, one of which is 
called Raptor.  The code Archit does use is more deterministic and predictable, allowing for 
selective recovery based on known missing blocks (e.g. a desired farmer is not responding (bad
for their reputation) or worse, responds with bad data.  Although computationally a bit more
intensive, by allowing selecting shard recovery, data bandwidth is potentially preserved, and
data bandwidth tends to be the limiting factor in network speed.

Instead of A,B,C&D for symbols, we can now think of each block as a symbol, and call it by its block
number.  So we now have B0, B1, B2... B31, plus 8 more blocks: R0-R7

We then bitwise compute transmission blocks T0-T(31+N) as follows:

T0 = B0 + B1 + B2 + B3 + B4 + B5 + B6
T1 = B1 + B2 + B3 + B4 + B5 + B6 + B7
.
.
.
T31 = B31 + R0 + R1 + R2 + R3 + R4 + R5
T32 = B0
T33 = B1
.
.
.
T39 = B7

Note that R0-R7 are simply copies of B0-B7.

What does this leave us with?  Well, T0-T31 are (32) independant equations with 32 variables, which
is mathematically solvable if all 32 were available.  T32-T39 are very simple redundant equations.
If fact, should B0 -> B7 go missing, a direct read of T32->T39 would suffice for a replacement 
value.  More importantly, we can lose up to R farmers and still be able to theoritically recover
the data.  Alas, recovering from a R farmer loss is troublesome for two reasons:  1)  All remaining
blocks will need to be read, and 2) its mathematically tough (think 32,768 XOR operations against
32MB symbol blocks).  Therefore, Archit strives to be able to survive the lost of any R-1 blocks.

Instead, the renter can use the knowledge of HOW the blocks were put together, along with the 
extra redundancy, to shortcut the process.

Note that with a R of 8, any 7 blocks blocks can be lost, guarenteeing (presuming ONLY 7 blocks
are lost) that at least one of T32-T39 survived.  This gives us a critical launching hint, much
like the D=5 in the introductory sample.

Also note that with an R of 8, only an additional 25% of data need to transmitted and stored on
the network, functionally providing the equivalent of 7 redundant copies.  Should a file exceed
~1GB in size, best effort is made to distribute the resulting shards across different nodes per
slice, to minimize the chances of a loss of any one farmer from owning multiple shards of a
renters file.

It should also be noted that a farmer is blind of most of this.  All they get is the filename to
store a ~32mb file under.  The renter maintains which farmer contains which blocks of their
file.  A farmer simply responds to two request:  1) An unpaid "are you there" Ping and 2) a
please send filename X.  Consideration was given to requesting a farmer to do a hash of shard 
as a relatively low-cost sanity check data was actually being stored - but that would require
the renter to pass the authentication seed stored.  That is problematical, since anyone who
wanted to scam the system could wait until such a request came in, save the result, and purge
the file (downside of being open sourced).  Renters, on the other hand, will eventually be 
able to randomly sample blocks and severely downgrade any farmer caught cheating, to the point
of being able to recall/recover any block stored by such a farmer.  e.g. farmer cheating 
would quickly destory the farmers reputation, at least with that renter.  Please see the FAQ
for a more detailed discussion around scammers.
