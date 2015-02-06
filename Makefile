PKG=github.com/bearprada/who_most_like_you

dev:
	rm -rf ${GOPATH}/pkg/${PKG}
	go get ${PKG}/bin/server
	${GOPATH}/bin/server
