package main

import "path"

// Repo represents the target repository that "go get" is attempting to get.
type Repo struct {
	Domain       string
	Organization string
	Project      string
	SubPath      string
	VCSHost      string
}

// PkgRoot generates the root import path, sans sub-package.
func (r *Repo) PkgRoot() string {
	return path.Join(r.Organization, r.Project)
}

// ImportPath is the full import path including sub-package.
func (r *Repo) ImportPath() string {
	return path.Join(r.PkgPath(), r.SubPath)
}

// PkgPath is the import path to the root package in the form: domain/organization/project
func (r *Repo) PkgPath() string {
	return path.Join(r.Domain, r.PkgRoot())
}
