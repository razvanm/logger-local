export GOPATH=$(CURDIR)
export VDLPATH=$(CURDIR)/src
PREFIX := PATH=$(JIRI_ROOT)/release/go/bin:$(PATH)

V23_BINARIES := v.io/x/ref/cmd/principal
V23_BINARIES += v.io/x/ref/services/mounttable/mounttabled
V23_BINARIES += v.io/x/ref/services/syncbase/syncbased
V23_BINARIES += v.io/x/ref/services/agent/v23agentd

.PHONY: all
all: cred binaries

binaries:
	jiri go install logger/server logger/client logger/setup $(V23_BINARIES)

cred:
	$(PREFIX) principal create -with-passphrase=false $(PWD)/cred/me me
	$(PREFIX) principal -v23.credentials=$(PWD)/cred/me fork -require-caveats=false -with-passphrase=false $(PWD)/cred/ext ext

.PHONY: cred-dump
cred-dump:
	$(PREFIX) principal -v23.credentials=$(PWD)/cred/ext dump

start-servers: mounttabled.pid syncbase-server.pid syncbase-client.pid

mounttabled.pid:
	$(PREFIX) mounttabled \
	    -v23.credentials=$(PWD)/cred/ext \
	    -v23.tcp.address=:8101 \
	    >/dev/null & echo $$! > mounttabled.pid

syncbase-server.pid:
	$(PREFIX) syncbased \
	    -v23.credentials=$(PWD)/cred/ext \
		-v23.namespace.root=/localhost:8101 \
		-name=syncbase/server \
		-root-dir=./syncbase-server >/dev/null & echo $$! > syncbase-server.pid

syncbase-client.pid:
	$(PREFIX) syncbased \
	    -v23.credentials=$(PWD)/cred/ext \
		-v23.namespace.root=/localhost:8101 \
		-name=syncbase/client \
		-root-dir=./syncbase-client >/dev/null & echo $$! > syncbase-client.pid

.PHONY: setup
setup: binaries
	$(PREFIX) ./bin/setup \
	    -v23.credentials=$(PWD)/cred/ext \
	    -mountpoint=/localhost:8101/syncgroups \
		/localhost:8101/syncbase/server \
		/localhost:8101/syncbase/client

.PHONY: server
server: binaries
	$(PREFIX) ./bin/server \
	    -v23.credentials=$(PWD)/cred/ext \
		-syncbase=/localhost:8101/syncbase/server

.PHONY: client
client: binaries
	$(PREFIX) ./bin/client \
	    -v23.credentials=$(PWD)/cred/ext \
		-syncbase=/localhost:8101/syncbase/client

.PHONY: kill-servers
kill-servers:
	kill $(shell cat *.pid)
	rm *.pid

.PHONY: clean
clean:
	rm -rf syncbase-server/ syncbase-client/ bin/ cred/ *.pid
