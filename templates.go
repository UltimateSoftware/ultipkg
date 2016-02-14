package main

import "html/template"

// These are templates used when serving UltiPkg pages.
const (
	goPackage = `
<!doctype>
<html>
<head>
  <meta name="go-import" content="{{.PkgPath}} git ssh://{{.VCSHost}}/{{.PkgRoot}}.git" />

  <title>UltiPkg - ultimate/dynamite</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.4.0/css/font-awesome.min.css">
  <style>
  body {
    margin: 0; padding: 0;
    font: 62.5%/1.25 Helvetica,Arial,sans-serif;
    -webkit-font-smoothing: antialiased;
  }
  .container {
    width: 600px;
    border: 1px solid #ccc;
    margin: 0px auto;
    margin-top: 4%;
    border-radius: 2px;
  }
  header {
    position: relative;

    background: #8dc641;

    padding: 10px;
    padding-left: 40px;

    color: #fff;
    font: 1.6em/1.25 Helvetica,Arial,sans-serif;

    -webkit-box-shadow: inset 0 0 0 0 #00757d;
    -moz-box-shadow: inset 0 0 0 0 #00757d;
    box-shadow: inset 0 0 0 0 #00757d;
  }

  h1 {
    margin: 0;
    margin-top: 5px;
  }

  header span {
    position: absolute;
    right: 10px;
    top: 50%;
    margin-top:-10px
  }

  a {
    text-decoration: none;
    color: #8dc641;
  }

  a:hover {
    text-decoration: underline;
  }

  a.btn {
    border: 1px solid #8dc641;
    background: #b0d87c;
    border-radius: 4px;
    text-decoration: none;
    padding: 4px 9px;
    color: #fff;
  }

  section {
    font-size: 1.75em;
    color: #8c847c;
    margin-top: 30px;
    margin-bottom: 40px;
    padding: 0 40px;
  }

  pre {
    display: block;
    padding: 9.5px;
    margin: 0 0 10px;
    font-size: 13px;
    line-height: 1.42857143;
    word-break: break-all;
    word-wrap: break-word;
    color: #666;
    background-color: #f5f5f5;
    border: 1px solid #ccc;
    border-radius: 2px;
  }
  footer {
    color: #bbb;
    text-align: center;
    font-size: 1.5em;
  }
  </style>
</head>
<body>
  <div class="container">
    <section>
      <p>To get the package, execute:</p>
      <pre>go get {{.ImportPath}}</pre>
      <p>To import this package, add the following line to your code:</p>
      <pre>import "{{.ImportPath}}"</pre>
    </section>
  </div>
  <footer>
    <p>Brought to you with <i class="fa fa-heart-o"></i> by <a href="http://www.ultimatesoftware.com">Ultimate Software</a></p>
  </footer>
</body>
</html>

`
)

// Parse templates
var (
	packageTemplate = template.Must(template.New("").Parse(goPackage))
)
