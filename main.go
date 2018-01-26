package main

import (
	"database/sql"
	"flag"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
	dbUtils             *utils.DbUtils
	config              = Configuration{}
	timezone            *time.Location
	appName             = "GoWebsiteExample"
	appVersion          = "0.0.0.2"
	authCookieStoreName = strings.Replace(appName, " ", "", -1)
	errCookieStoreName  = strings.Replace(appName, " ", "", -1) + "Err"
	cookieStore         *sessions.CookieStore
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.DebugLevel

	// init databaseutils
	dbUtils = new(utils.DbUtils)

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

	err = dbUtils.Connect2Database(&db, config.Database.DbType, config.Database.DbURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	timezone, err = time.LoadLocation(config.General.Timezone)
	if err != nil {
		log.Println(err)
		return
	}

	audit.SetLogger(appName+"/"+appVersion, log, dbUtils)
	audit.SetWaitGroup(&wg)

	mw := io.MultiWriter(os.Stdout, audit)
	log.Out = mw

	cookieStore, err = getNewCookieStore()
	if err != nil {
		log.Println(err)
		return
	}

	err = initializeDatabase()
	if err != nil {
		log.Println(err)
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

	err = http.ListenAndServe(*addr,
		csrf.Protect(
			key,
			csrf.Secure(config.General.IsHTTPS),
			csrf.Path("/"),
			csrf.FieldName("csrfToken"),
			csrf.CookieName("csrfCookie"),
			csrf.HttpOnly(true),
			csrf.MaxAge(24*3600),
			csrf.RequestHeader("X-CSRF-Token"))(router))

	if err != nil {
		log.Println(err)
		return
	}

	log.Info("Closing application...")
	wg.Wait()
}
