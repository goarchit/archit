1) First install Go per:  http://ask.xmodulo.com/install-go-language-linux.html
2) To define your workspace:  https://golang.org/doc/code.html
   For this project, you should end up with a ~/src/github.com/goarchit directory
3) Clone the source for github:  git clone https://github.com/goarchit/archit into the
   goarchit directory you created in step 2
4) Make sure you installed Go in the step 2 (for gentoo:  "emerge go")
5) Do "go install archit" in ./archit.  You will likely get a lot of unresolved 
   references like github.com/userx/proga.  For each do a "go get github.com/userx/proga" to 
   resolve.  This is a one time thing to bring in all the dependencies.
   As of this writing:
      go get github.com/Unknwon/goconfig
      go get github.com/astaxie/beego/logs
      go get github.com/boltdb/bolt
      go get github.com/briandowns/spinner
      go get github.com/btcsuite/btcd/chaincfg
      go get github.com/btcsuite/btcrpcclient
      go get github.com/btcsuite/btcutil
      go get github.com/jessevdk/go-flags
      go get github.com/minio/blake2b-simd
      go get github.com/valyala/gorpc
      go get golang.org/x/crypto/scrypt

      Each of the above will only take a few seconds to run.  Hint:  cut & paste them as a batch.

      A "go install archit" should now compile clean.

      Note: for a clean "archit status" response:
            edit github.com/btcsuite/btcutil/amount.go 
            and change lines 36->44 from "BTC" to "IMAC"
 
6) You will need an IMACredit wallet running:
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
7) Alternatively download the Windows wallet from http://www.imacredit.org - there is no need
   for the wallet to be running on the same machine as GoArchit
8) Although you can run archit without a configuration file, it gets tedious to do so...
   Suggest creating a ~/.archit/archit.conf file with the following:
        KeyPass = "Some very private key string... make it long and unique."
        WalletAddr = <a valid IMAC Deposit Address from the wallet you instead in step 6>
              This should look something like 9KsqKMgLjzBWKidhes356Kjhdwbd9BT4Te
        WalletUser = <the RPC user defined in imacredit.conf>
        WalletPassword = <the RPC password defined in imacredit.conf>
        RPCuser = architrpc  <or anything else you like>
        RPCPassword = SomeRandomPasswordWithoutSpaces902q34890u328490
        
9)  Do an "archit --help" to see options
10) Do an "archit farm" to start farming
11) Have Fun!

Come /join us on freenode channels #imacredit and/or #goarchit

To upgrade to the latest push:

cd ~/src/github.com/goarchit/archit
git pull https://github.com/goarchit/archit
cd archit
go install

(and restart your farmer)

Special note for GRSEC users (typically hardened kernels):
   Go programs will generate "denied RWX mprotect of" type messages
   Issue a "paxctl-ng -m /home/goarchit/bin/archit" to fix it (using your binaries location)
   You may need to do this each time you regenerate archit
