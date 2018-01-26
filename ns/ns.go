package ns

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Log() {
	name := time.Now().Format("tcp20060102.log")
	logfile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	fmt.Fprintln(logfile, "DATE,TIME,ESTABLISHED,CLOSE_WAIT,SYN_SENT")

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for range t.C {
		x, y, z := numEstablished()
		fmt.Fprintf(
			logfile,
			"%s,%d,%d,%d\n",
			time.Now().Format("2006/01/02,15:04:05"),
			x,
			y,
			z,
		)
	}
}

func numEstablished() (x, y, z int) {
	out, err := exec.Command("netstat", "-na").Output()
	if err != nil {
		log.Fatal(err)
	}

	r := strings.NewReader(string(out))

	s := bufio.NewScanner(r)
	for s.Scan() {
		if strings.Index(s.Text(), "ESTABLISHED") != -1 {
			x++
		}
		if strings.Index(s.Text(), "CLOSE_WAIT") != -1 {
			y++
		}
		if strings.Index(s.Text(), "SYN_SENT") != -1 {
			z++
		}
	}
	return x, y, z
}
