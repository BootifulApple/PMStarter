package main

import (
	"os"
	"os/exec"
	"strings"
	"flag"
	"time"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("PMStarter")

func main() {
	logging.SetFormatter(logging.MustStringFormatter(
		`%{color}%{level:.4s}> %{color:reset} %{message}`,
	))
	
	loop := flag.Bool("loop", false, "Loop startup")
	delay := flag.Int("delay", 5, "Delay next loop by n seconds")
	php_bin := flag.String("php", findPHP(), "Use custom PHP Binary")
	pm_file := flag.String("pm", findPM(), "Use custom PM phar / php script")
	flag.Parse()
	
	if *php_bin == "" {
		log.Error("Cannot find PHP Binary")
		os.Exit(1)
	}
	
	if *pm_file == "" {
		log.Error("Cannot find PocketMine-MP File")
		os.Exit(1)
	}
	
	n := 0
	for run := true; run; run = *loop {
		log.Info("Starting Server...")
		cmd := exec.Command(*php_bin, *pm_file)

		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		if err := cmd.Run() ; err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				log.Errorf("Server exited with code: %d", exitError.ExitCode())
				os.Exit(exitError.ExitCode())
			}
		}
		if *loop && *delay > 0 {
			if n > 0 {
				log.Infof("Restarted %d times", n)
			}
			log.Infof("To escape the loop, press CTRL+C now. Otherwise, wait %d seconds for the server to restart.", *delay)
			time.Sleep(time.Duration(*delay) * time.Second)
			n++
		}
	}
}

func findPHP() string {
	if fileExists("./bin/php7/bin/php") {
		err := os.Setenv("PHPRC", "")
		if err != nil {
			log.Error(err)
		}
		return "./bin/php7/bin/php"
	}
	out, err := exec.Command("bash", "-c", "type -p php").CombinedOutput()
	if err != nil {
		log.Error(err)
	}
	out_str := strings.TrimSpace(string(out))
	if out_str != "" {
		return out_str
	}
	return ""
}

func findPM() string {
	if fileExists("./PocketMine-MP.phar") {
		return "./PocketMine-MP.phar"
	}
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
