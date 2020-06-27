GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

PROGRAMS = pihole-stats-exporter

depend:
	env GOPATH=$(GOPATH) go get -u github.com/sirupsen/logrus
	env GOPATH=$(GOPATH) go get -u gopkg.in/ini.v1
	env GOPATH=$(GOPATH) go get -u github.com/gorilla/mux

build:
	env GOPATH=$(GOPATH) go install $(PROGRAMS)

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/sbin

strip: build
	strip --strip-all $(BINDIR)/pihole-stats-exporter

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/pihole-stats-exporter $(DESTDIR)/usr/sbin

clean:
	/bin/rm -f bin/pihole-stats-exporter

distclean: clean
	/bin/rm -rf src/gopkg.in/
	/bin/rm -rf src/github.com/

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: build strip install

