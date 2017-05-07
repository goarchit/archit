// Status command:  Display Network and Wallet status
// Originally work created on 1/8/2017
//

package cmd

import (
	"fmt"
	"github.com/goarchit/archit/config"
	"github.com/goarchit/archit/log"
	"github.com/goarchit/archit/util"
	"github.com/valyala/gorpc"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type BackupCommand struct {
}

var backupCmd BackupCommand

func init() {
	_, err := config.Parser.AddCommand("backup", "Flushes all active databases and creates a backup of each", "", &backupCmd)
	if err != nil {
		fmt.Println("Internal error parsing Backup command:", err)
		os.Exit(1)
	}
}

func (ec *BackupCommand) Execute(args []string) error {
	config.Conf(false)

	// Copy file information database
	src := filepath.Join(util.DBDir, util.FileDBName)
	dst := filepath.Join(util.DBDir, util.FileDBName+".bkup")
	err := util.CopyFile(src, dst)
	if err != nil {
		log.Critical("Backup failed:", err)
	}

	// Insert RPC code to query the farmer

	port := util.PortBase + 1
	serverIP := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	c := gorpc.NewTCPClient(serverIP)
	c.Start()
	defer c.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("PeerDBSync", func() {})
	dc := d.NewFuncClient(c)
	errmsg, err := dc.Call("PeerDBSync", nil)
	if err != nil {
		log.Error("Backup failed:", err)
	}
	if errmsg != nil {
		log.Error("Backup may have failed:", err)
	}
	// Copy peers information database
	src = filepath.Join(util.DBDir, util.PeerDBName)
	dst = filepath.Join(util.DBDir, util.PeerDBName+".bkup")
	err = util.CopyFile(src, dst)
	if err != nil {
		log.Critical("Backup failed:", err)
	}
	log.Console("Backup complete: you can now do whatever to FileInfo.bolt.bkup and PeerInfo.bolt.bkup in", util.DBDir)
	return nil
}
