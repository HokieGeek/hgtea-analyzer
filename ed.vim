#!/usr/bin/vim -S

args *.go Dockerfile ed.vim

edit main.go

tabnew db_test.go
topleft vsplit
edit db.go

tabnew types_test.go
topleft vsplit
edit types.go

tabfirst
