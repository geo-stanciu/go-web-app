package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"encoding/gob"

	"github.com/geo-stanciu/go-utils/utils"
	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
)

var (
	log                 = logrus.New()
	audit               = utils.AuditLog{}
	templateDelims      = []string{"{{%", "%}}"}
	templates           *template.Template
	addr                *string
	db                  *sql.DB
	dbutl               *utils.DbUtils
	config              = Configuration{}
	timezone            *time.Location
	hs                  *http.Server
	stop                chan os.Signal
	appName             = "GoWebsiteExample"
	appVersion          = "0.0.0.5"
	authCookieStoreName = strings.Replace(appName, " ", "", -1)
	errCookieStoreName  = strings.Replace(appName, " ", "", -1) + "Err"
	cookieStore         *sessions.CookieStore
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.DebugLevel

	// init databaseutils
	dbutl = new(utils.DbUtils)

	// init exit
	stop = make(chan os.Signal, 1)

	// register SessionData for cookie use
	gob.Register(&SessionData{})

	// initialize the templates,
	// since we have custom delimiters.
	basePath := "templates/"

	err := filepath.Walk(basePath, parseTemplate(basePath))
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	var err error
	var wg sync.WaitGroup

	cfgFile := "./app.config"
	err = config.ReadFromFile(cfgFile)
	if err != nil {
		log.Println(err)
		return
	}

	time2wait := 5
	startTime := time.Now()
	endWait := startTime.Add(time.Duration(time2wait) * time.Minute)

	for {
		err = dbutl.Connect2Database(&db, config.Database.DbType, config.Database.DbURL)
		if err != nil {
			log.Println(err)

			now := time.Now()

			if now.After(endWait) {
				log.Printf("Failed to connect to the database after trying for %d minute(s).\n", time2wait)
				return
			}

			time.Sleep(30 * time.Second)
		} else {
			break
		}
	}
	defer db.Close()

	timezone, err = time.LoadLocation(config.General.Timezone)
	if err != nil {
		log.Println(err)
		return
	}

	audit.SetLogger(appName, appVersion, log, dbutl)
	audit.SetWaitGroup(&wg)
	defer audit.Close()

	mw := io.MultiWriter(os.Stdout, audit)
	log.Out = mw

	cookieStore, err = getNewCookieStore()
	if err != nil {
		audit.Log(err, "get cookie store", "error while initializing the cookie store")
		return
	}

	err = parseArguments()
	if err != nil {
		audit.Log(err, "parse arguments", "error while parsing the command line arguments")
		return
	}

	err = initializeDatabase()
	if err != nil {
		audit.Log(err, "initialize database", "error while initializing the database")
		return
	}

	// server flags
	addr = flag.String("addr", ":"+config.General.Port, "http service address")

	flag.Parse()

	log.WithField("port", *addr).Info("Starting listening...")

	router := pat.New()

	// Normal resources
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("public/static"))))
	router.PathPrefix("/images/").Handler(
		http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))
	router.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))

	router.PathPrefix("/favicon.ico").Handler(http.NotFoundHandler())

	router.Get("/", handler)
	router.Get("/{url}", handler)
	router.Post("/", handler)
	router.Post("/{url}", handler)

	encodeKeys, _ := getCookiesEncodeKeys()
	if encodeKeys == nil || len(encodeKeys) == 0 {
		log.Error("No encode Keys")
		wg.Wait()
		return
	}

	var key []byte
	l1 := len(encodeKeys)
	if l1 >= 4 {
		key = encodeKeys[3]
	} else if l1 >= 2 {
		key = encodeKeys[1]
	} else {
		key = encodeKeys[0]
	}

	hs = &http.Server{
		Addr: ":" + config.General.Port,
		Handler: csrf.Protect(
			key,
			csrf.Secure(config.General.IsHTTPS),
			csrf.Path("/"),
			csrf.FieldName("csrfToken"),
			csrf.CookieName("csrfCookie"),
			csrf.HttpOnly(true),
			csrf.MaxAge(24*3600),
			csrf.RequestHeader("X-CSRF-Token"))(router),
	}

	go func() {
		if err = hs.ListenAndServe(); err != http.ErrServerClosed {
			audit.Log(err, "Server Start", "Error listening on "+config.General.Port)
			log.Fatal(err)
			return
		}
	}()

	serverStop()

	// wait for all logs to be written
	wg.Wait()
}

func serverStop() {
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hs.Shutdown(ctx); err != nil {
		audit.Log(nil, "Server Stop", "Error closing application")
	} else {
		audit.Log(nil, "Server Stop", "Closing application...")
	}

	time.Sleep(5 * time.Second)
}

func parseArguments() error {
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		switch arg {
		case "--stop":
			audit.Log(nil, "process-stop", "Stop process requested...")

			err := stopProcess(utils.String2int(config.General.Port))
			if err != nil {
				audit.Log(err, "process-stop", "Error encountered trying to stop the process...")
				return err
			}

			os.Exit(0)
		default:
			err := fmt.Errorf("unknown argument \"%s\"", arg)
			audit.Log(err, "parse-arguments", "Unknown argument encountered...")
			return err
		}
	}

	return nil
}

func stopProcess(port int) error {
	schema := "http"
	if config.General.IsHTTPS {
		schema = "https"
	}

	url := fmt.Sprintf("%s://localhost:%d/stop-process", schema, port)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	audit.Log(nil, "process-stop", "Stop process requested", "server response", string(buf))

	return nil
}
