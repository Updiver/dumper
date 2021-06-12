module github.com/updiver/dumper

go 1.15

require (
	github.com/go-git/go-git/v5 v5.4.2
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
)

// this is only while we are developing users package
replace github.com/updiver/dumper/users => /home/klaus/dumper/users
