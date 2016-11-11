#!/usr/bin/vim -S

args *.go teas/*.go Dockerfile README.md

for s:p in [ 'teas/main', 'db', 'types', 'tsv' ]
    execute "tabnew " . s:p . "_test.go"
    topleft vsplit
    execute "edit " . s:p . ".go"
endfor

tabfirst
tabclose
