package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/UltimateSoftware/ultipkg/config"

	"github.com/Sirupsen/logrus"
)

var (
	// Version will be defined by our build and release Rake tasks.
	// We use ldflags to define this as the application compiles.
	Version = ""

	// Logger is our configured instance of Logrus.
	Logger = logrus.New()

	// VCSHost is the user, domain, and port that clients will need
	// to checkout the code from the server.
	VCSHost = config.Get("vcs.host", "git@git.example.com")

	// Domain is the actual domain this server is running.
	Domain = config.Get("domain", "pkg.example.com")
)

func init() {
	logLevel, err := logrus.ParseLevel(config.Get("log.level", "info"))
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

// ultipkg is the meat & potatoes of the entire application. This is where
// we figure out the requested repo and generate the correct response.
func ultipkg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Server", "Ultipkg")
	w.Header().Set("X-Server-Version", Version)

	// Using a logger with fields makes tracing things easier.
	L := Logger.WithFields(logrus.Fields{
		"secure":  r.TLS != nil,
		"agent":   r.UserAgent(),
		"remote":  r.RemoteAddr,
		"version": Version,
	})

	// For now, the index just returns a 404. This could just as easily return
	// an actual page, or a redirect to someplace else.
	if r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	// Split the URL path into its components so we can start extracting data
	// from it.
	fragments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	L.WithFields(logrus.Fields{
		"frags": fragments,
		"len":   len(fragments),
	}).Debug("Split path")

	repo := &Repo{
		Domain:  Domain,
		VCSHost: VCSHost,
	}

	if len(fragments) >= 1 {
		// example.com/hello
		repo.Organization = fragments[0]
	}
	if len(fragments) >= 2 {
		// example.com/hello/world
		repo.Project = fragments[1]
	}
	if len(fragments) > 2 {
		// example.com/hello/world/goodbye/moon
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
