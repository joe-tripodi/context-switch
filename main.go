package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func activeWindowMacCmd() *exec.Cmd {
	return exec.Command(
		"osascript",
		"-e",
		"tell application \"System Events\" to get name of application processes whose frontmost is true and visible is true",
	)
}

func activeChromeTabTitleMacCmd() *exec.Cmd {
	return exec.Command(
		"osascript",
		"-e",
		"tell application \"Google Chrome\" to return title of active tab of front window",
	)
}

var cswitch = 0
var prev []byte
var prevTab []byte

func checkSwitchMac(db *sql.DB) {
	inc := false
	window, err := activeWindowMacCmd().Output()
	if err != nil {
		log.Fatal("There was an error matey: ", err)
		return
	}

	if bytes.Compare(prev, window) != 0 {
		cswitch++
		prev = window
		fmt.Print("Switched to: ", string(window))
		inc = true
    if string(window) != "Google Chrome\n" {
      record := Record{
        App: strings.TrimSpace(string(window)),
        Tab: "",
        When: time.Now(),
      }
      SaveRecord(db, &record)
    }
	}

	if string(window) == "Google Chrome\n" {
		activeTab, err := activeChromeTabTitleMacCmd().Output()
		if err != nil {
			log.Fatal("Error getting the tab: ", err)
			return
		}
		if bytes.Compare(prevTab, activeTab) != 0 {
			if !inc {
				cswitch++
			}
			prevTab = activeTab
			fmt.Print(string(activeTab))

      record := Record{
        App: strings.TrimSpace(string(window)),
        Tab: strings.TrimSpace(string(activeTab)),
        When: time.Now(),
      }

      SaveRecord(db, &record)
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Println("No. of context switches: ", cswitch)
}

func macWindows(db *sql.DB) {
	for {
		checkSwitchMac(db)
	}

}

func main() {
  cfg := mysql.Config{
    User:   os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "ctxswitches",
  }

  db, err := sql.Open("mysql", cfg.FormatDSN())
  if err != nil {
    log.Fatal(err)
  }

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

  log.Print("Connected to database")

	switch runtime.GOOS {
	case "darwin":
		macWindows(db)
	case "windows":
		fmt.Println("Not implemented for windows")
	case "freebsd":
		fmt.Println("Who knows?")
	case "linux":
		fmt.Println("Not implemented for linux yet")
	default:
		fmt.Println("Not implemented for ", runtime.GOOS)
	}
}
