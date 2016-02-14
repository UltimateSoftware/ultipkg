package main

import (
	"net/http"
	"strings"
	"sync"

	config "pkg.ulti.io/ultimate/envconfig"

	"github.com/Sirupsen/logrus"
)

// Version will be defined by our build and release Rake tasks.
// We use ldflags to define this as the application compiles.
var Version = ""

// Logger is our configured instance of Logrus.
var Logger = logrus.New()

func init() {
	logLevel, err := logrus.ParseLevel(config.Get("logg.level", "info"))
	if err != nil {
		Logger.Panicf("Fatal error: %s \n", err)
	}

	Logger.Level = logLevel
	config.Logger = Logger
}

func main() {
	switch config.Get("log.format") {
	case "text":
		Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	default:
		Logger.Formatter = &logrus.JSONFormatter{}
	}

	http.HandleFunc("/", ultipkg)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		addr := config.Get("addr", "localhost:8080")
		Logger.WithFields(logrus.Fields{
			"addr":    addr,
			"version": Version,
		}).Info("starting http server")
		Logger.Fatal(http.ListenAndServe(addr, nil))

		wg.Done()
	}()

	go func() {
		certFile := config.Get("ssl.certificate")
		keyFile := config.Get("ssl.privatekey")

		L := Logger.WithFields(logrus.Fields{
			"certificate": certFile,
			"privatekey":  keyFile,
		})

		if certFile == "" || keyFile == "" {
			L.Warn("Skipping SSL")
		} else {
			addr := config.Get("addr_tls", "localhost:8081")

			Logger.WithFields(logrus.Fields{
				"addr":    addr,
				"version": Version,
			}).Info("starting https server")
			L.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, nil))
		}

		wg.Done()
	}()

	wg.Wait()

	Logger.Info("shutting down")
}

func ultipkg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Server", "Ultipkg")
	w.Header().Set("X-Server-Version", Version)

	L := Logger.WithFields(logrus.Fields{
		"secure":  r.TLS != nil,
		"agent":   r.UserAgent(),
		"remote":  r.RemoteAddr,
		"version": Version,
	})

	if r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	fragments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	L.WithFields(logrus.Fields{
		"frags": fragments,
		"len":   len(fragments),
	}).Debug("Split path")

	repo := &Repo{
		Domain: config.Get("domain", "pkg.ulti.io"),
	}

	if len(fragments) >= 1 {
		repo.Organization = fragments[0]
	}
	if len(fragments) >= 2 {
		repo.Project = fragments[1]
	}
	if len(fragments) > 2 {
		repo.SubPath = strings.Join(fragments[2:], "/")
	}

	L.WithField("repo", repo).Debug("Collected Repo information")

	w.Header().Set("Content-Type", "text/html")

	err := packageTemplate.Execute(w, repo)
	if err != nil {
		L.WithField("err", err).Error("Could not execute packageTemplate")
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}

	L.Infof("%s %s", r.Method, r.URL.Path)
}
