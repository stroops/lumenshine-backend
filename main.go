package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"

	"github.com/spf13/cobra"
)

var (
	messagingPath, notificationPath, mailPath, dbPath, tfaPath, jwtPath, userAPIPath, payAPIPath, payPath, ssePath, sseAPIPath string
)

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(ex)

	dbPath = path.Join(dir, "services", "db")
	tfaPath = path.Join(dir, "services", "2fa")
	jwtPath = path.Join(dir, "services", "jwt")
	userAPIPath = path.Join(dir, "api", "userapi")
	payAPIPath = path.Join(dir, "api", "payapi")
	mailPath = path.Join(dir, "services", "mail")
	notificationPath = path.Join(dir, "services", "notification")
	messagingPath = path.Join(dir, "services", "messaging")
	payPath = path.Join(dir, "services", "pay")
	ssePath = path.Join(dir, "services", "sse")
	sseAPIPath = path.Join(dir, "api", "sseapi")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(migrateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func buildApp(dir string, withRice bool) error {
	var err error

	log.Println("Building: ", dir)
	if withRice {
		//build rice-file first
		cmd := exec.Command("rice", "embed-go")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Start()
		if err != nil {
			log.Fatalf("error calling rice-embed for: %s, %v", dir, err)
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("go", "build")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Start()
	if err != nil {
		log.Fatalf("error calling build for: %s, %v", dir, err)
	}
	return cmd.Wait()
}

func runExe(dir string, exe string, args ...string) {
	cmdStr := path.Join(dir, exe)

	cmd := exec.Command(cmdStr, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		log.Fatalf("error calling start for: %s, %v", exe, err)
	}
	err = cmd.Wait()
}

func init() {
	migrateCmd.PersistentFlags().BoolVar(&migrateDown, "down", false, "will migrate down-to 0. If not specified, will migrate up")
}

var rootCmd = &cobra.Command{
	Use:   "lumenshine-backend",
	Short: "lumenshine-backend tasks",
	Long:  `All possible functionality for handling the microservices`,
	Run: func(cmd *cobra.Command, args []string) {
		// Run the services
		log.Println(`
You need to specify a command:
icop run [2fa, jwt, db, userapi, mail, payapi, pay, notification, messaging, sse, sseapi]
e.g. ./lumenshine-backend run 2fa jwt db userapi mail payapi pay notification messaging sse sseapi
lumenshine-backend build
lumenshine-backend migrate [--down]`)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "running the microservices",
	Long:  `This command runs all the defined microservices.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("*************** Press <ctrl> + <c> to stop ************************")
		// Run the services
		if len(args) == 0 {
			//start all
			go runExe(dbPath, "db")
			log.Println("Waiting 5 Sekonds for the DB service to start up, befor starting other services!")
			time.Sleep(5 * time.Second)

			go runExe(tfaPath, "2fa")
			go runExe(jwtPath, "jwt")
			go runExe(userAPIPath, "userapi")
			go runExe(mailPath, "mail")
			go runExe(payAPIPath, "payapi")
			go runExe(payPath, "pay")
			go runExe(notificationPath, "notification")
			go runExe(messagingPath, "messaging")
			go runExe(sseAPIPath, "sseapi")
			go runExe(ssePath, "sse")
		} else {
			if helpers.StringInSliceI("db", args) {
				go runExe(dbPath, "db")
				log.Println("Waiting 5 Sekonds for the DB service to start up, befor starting other services!")
				time.Sleep(5 * time.Second)
			}

			if helpers.StringInSliceI("2fa", args) {
				go runExe(tfaPath, "2fa")
			}

			if helpers.StringInSliceI("jwt", args) {
				go runExe(jwtPath, "jwt")
			}

			if helpers.StringInSliceI("userapi", args) {
				go runExe(userAPIPath, "userapi")
			}

			if helpers.StringInSliceI("payapi", args) {
				go runExe(payAPIPath, "payapi")
			}

			if helpers.StringInSliceI("mail", args) {
				go runExe(mailPath, "mail")
			}

			if helpers.StringInSliceI("pay", args) {
				go runExe(payPath, "pay")
			}

			if helpers.StringInSliceI("notification", args) {
				go runExe(notificationPath, "notification")
			}

			if helpers.StringInSliceI("messaging", args) {
				go runExe(messagingPath, "messaging")
			}

			if helpers.StringInSliceI("sseapi", args) {
				go runExe(sseAPIPath, "sseapi")
			}

			if helpers.StringInSliceI("sse", args) {
				go runExe(ssePath, "sse")
			}
		}

		done := make(chan bool)
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c          // got signal
			done <- true //stop main rutine
		}()

		<-done //wait for sigterm
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build all microservices",
	Long:  `This command rebuilds all microservices.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build the services
		buildApp(dbPath, true)
		buildApp(tfaPath, false)
		buildApp(jwtPath, false)
		buildApp(userAPIPath, true)
		buildApp(mailPath, false)
		buildApp(payAPIPath, true)
		buildApp(payPath, false)
		buildApp(notificationPath, false)
		buildApp(messagingPath, false)
		buildApp(sseAPIPath, true)
		buildApp(ssePath, true)
	},
}

var (
	migrateDown bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate all databases",
	Long:  `This command migrates all databases`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build the services
		if migrateDown {
			log.Println("migrating down-to 0")
			runExe(dbPath, "db", "migrate", "down-to", "0")
		} else {
			log.Println("migrating up")
			runExe(dbPath, "db", "migrate", "up")
		}
	},
}
