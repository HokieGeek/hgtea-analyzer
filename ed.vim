#!/usr/bin/vim -S

args *.go Dockerfile README.md ed.vim

edit teas/main.go

tabnew db_test.go
topleft vsplit
edit db.go

tabnew types_test.go
topleft vsplit
edit types.go

tabfirst
