module github.com/sjsafranek/gosimpleserver

go 1.19

replace github.com/sjsafranek/httpfileserver v0.0.3 => github.com/sjsafranek/httpfileserver v0.0.4
replace github.com/sjsafranek/logger v1.0.3 => ../logger

require github.com/fsnotify/fsnotify v1.6.0

require (
	github.com/sjsafranek/logger v1.0.3 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)
