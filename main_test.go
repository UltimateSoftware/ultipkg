package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/nbio/st"
)

func TestMain(m *testing.M) {
	flag.Parse()
	Domain = "pkg.example.com"
	VCSHost = "git@git.example.com:7999"
	os.Exit(m.Run())
}

type outcome struct {
	getPath          string
	importPrefixPath string
	repoRootPath     string
}

func (o outcome) importPrefix() string {
	return Domain + o.importPrefixPath
}

func (o outcome) repoRoot() string {
	return "ssh://" + VCSHost + o.repoRootPath + ".git"
}

var scenarios = map[string]outcome{
	"/hello/world":       outcome{"/hello/world", "/hello/world", "/hello/world"},
	"/hello/world/a":     outcome{"/hello/world/a", "/hello/world", "/hello/world"},
	"/hello/world/a/b/c": outcome{"/hello/world/a/b/c", "/hello/world", "/hello/world"},
	"/hello/world/":      outcome{"/hello/world", "/hello/world", "/hello/world"},
}

func get(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	ultipkg(w, req)

	return w
}

func TestUltipkg(t *testing.T) {
	scenarioNum := 0
	for path, scenario := range scenarios {
		scenarioNum++

		w := get(path)

		st.Expect(t, w.Code, http.StatusOK, scenarioNum)

		document, _ := goquery.NewDocumentFromReader(w.Body)
		metaContent, exists := document.Find("meta[name='go-import']").First().Attr("content")
		st.Expect(t, exists, true, scenarioNum)
		st.Expect(t, metaContent, scenario.importPrefix()+" git "+scenario.repoRoot(), scenarioNum)
	}
}

func BenchmarkRealPackage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		get("/uc/tenantmanagement")
	}
}

func BenchmarkRealSubPackage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		get("/uc/tenantmanagement/a/b/c")
	}
}

// func BenchmarkFauxPackage(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		b.ReportAllocs()
// 		subject.Get("/ultimate/ultipkg")
// 	}
// }
// func BenchmarkFauxSubPackage(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		b.ReportAllocs()
// 		subject.Get("/ultimate/ultipkg/a/b/c")
// 	}
// }
