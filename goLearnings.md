Things I learned about Go developing this project:

1)  Strings are read-only, []byte arrays are not.
2)  Slices are passed by pointer
3)  Maps need to be made before used
4)  2D arrays arn't.  [x][y]byte != [x,y]byte (format used by most other languages), rather its
    an array of arrays, and since arrays are simply slices of fixed length, its really a
    fixed length slice or fixed length slices.  Therefore "var z [x][y]byte" will return
    differnet values for &z and &z[0]
6)  By all that is Holy, casing matters bigtime (variables that start in lower case are local
    variables that start in upper case are global).
7)  What you think is global (var db string) becomes very local if you do a db := something.
8)  A map cannot be referenced concurrently - all map activity needs to be wrapped in mutex
    protection if referenced in a goroutine.
9)  Don't be stupid and define your mutex variable locally in the goroutine.
10) GOBing a structure takes a lot less buffer space than JSON Marshalling it - best if your
    only going to be using the data interally. 
11) Gzip and the other Compress utilities all use Flate, so unless your going to be processing
    data externally, might as well use Flate directly.
