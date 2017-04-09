In order to avoid import cycles, anything in util MUST be very basic and not call things in 
higher level packages

If the reference to a higher level package is just to grab a variable, the simple fix is to
add a copy in ../util/include.go and have the higher level package copy to that, or simply
have the higher level function reference a util.<variable>.

Bottom line:  About the only safe local code in include is github.com/goarchit/archit/log
