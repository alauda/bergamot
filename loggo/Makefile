default: check

check:
	go test

docs:
	godoc2md github.com/alauda/bergamot/loggo > README.md
	sed -i 's|\[godoc-link-here\]|[![GoDoc](https://godoc.org/github.com/alauda/bergamot/loggo?status.svg)](https://godoc.org/github.com/alauda/bergamot/loggo)|' README.md 


.PHONY: default check docs
