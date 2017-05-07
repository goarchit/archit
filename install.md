1) First install Go per:  http://ask.xmodulo.com/install-go-language-linux.html
2) To define your workspace:  https://golang.org/doc/code.html
3) Make sure you installed Go in the step 2 (for gentoo:  "emerge go")
4) "go get github.com/goarchit/archit" should create all directorys and pull ind dependancies
   If you get a "git" error, manually clone it:
      cd src/github.com/goarchit
      git clone https://github.com/goarchit/archit
      then repeat the "go get github.com/goarchit/archit" to pull in the dependancies
    A "go install archit" should now compile clean.

    Note: for a clean "archit status" response:
       edit github.com/btcsuite/btcutil/amount.go 
       and change lines 36->44 from "BTC" to "IMAC"
 
5) You will need an IMACredit wallet running:
   In the directory of your choice:  git clone https://github.com/imacredit/imacredit
   cd imacredit/src
   make -f makefile.unix
   Note IMACredit has a few dependencies common to most coins (openssl, berkley db, boost,...)
   Start via ./imacreditd
   After a few moment you will get the standard new coin message about needing a config file
      with some security settings.  Create per the messages and restart.
   ./imacreditd&
   and just let it run - it should sync with the IMAC network within a few hours, or perhaps
   overnight, depending on network speeds
   After a few minutes, you should be able to issue a command like:
   ./imacreditd getaccountaddress goarchit - which will create a valid IMAC wallet address and 
   lable it "goarchit".  Of course, you can pick your own label.
6) Alternatively download the Windows wallet from http://www.imacredit.org - there is no need
   for the wallet to be running on the same machine as GoArchit
7) Although you can run archit without a configuration file, it gets tedious to do so...
   Suggest creating a ~/.archit/archit.conf file with the following:
        KeyPass = Some very private key string... make it long and unique, quotes not needed
        WalletAddr = <a valid IMAC Deposit Address from the wallet you instead in step 6>
              This should look something like 9KsqKMgLjzBWKidhes356Kjhdwbd9BT4Te
        WalletUser = <the RPC user defined in imacredit.conf>
        WalletPassword = <the RPC password defined in imacredit.conf>
        
    And... since we like to make life easy... you can use the "archit create" command to 
    help you do so.  It will even query your wallet and create an address for you.
    
    Do an "archit create --help" for detail, but here is an example:

    archit create --PortBase 3000 --WalletUser myrpcid --WalletPassword WayTooShort

    Make sure you have your port number DNATed through your firewall to the machine
    your running this on.

8)  Do an "archit --help" to see options
9)  Do an "archit farm" to start farming
10) Have Fun!

Come /join us on freenode channels #imacredit and/or #goarchit

To upgrade to the latest push:

cd src/github.com/goarchit/archit
git pull
go install  (and resolve any missing packages with go get <package>)

Then restart your farmer

Special note for GRSEC users (typically hardened kernels):
   Go programs will generate "denied RWX mprotect of" type messages
   Issue a "paxctl-ng -m /home/goarchit/bin/archit" to fix it (using your binaries location)
   You will need to do this each time you regenerate archit
   Alternatively, you can disable "Restrict mprotect()" in your kernel (Grsecurity ->
      Customize Configuration -> PaX -> Non-executable pages -> (uncheck) Restrict mprotect()
