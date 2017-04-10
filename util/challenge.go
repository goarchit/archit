// Challenge - a simple routine to interactively challenge the user before they do something
// potentially stupid
// Originally work created on 4/9/2017
//

package util

import (
	"github.com/goarchit/archit/log"
	"math/rand"
	"fmt"
	"time"
)

func Challenge() {
	var enum int

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(899999)+100000

	fmt.Println("*** WARNING!  You are about to display very sensitive data! ***")
	fmt.Print("If you wish to proceed enter the number ",num,": ")
	_, err := fmt.Scanf("%d", &enum)
	if err != nil {
		log.Critical("Error getting response number:", err)
	}
	if num != enum {
		log.Critical("Response does not match")
	}
	fmt.Print("\n")
}
