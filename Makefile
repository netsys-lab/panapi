CC       = go build
BUILDDIR = ./bin
PRGS = server client

all: $(PRGS)

.PHONY: server
server:
	$(CC) -o $(BUILDDIR)/$@ cmd/server/*.go

.PHONY: client
client:
	$(CC) -o $(BUILDDIR)/$@ cmd/client/*.go

clean:
	rm -f $(BUILDDIR)/*
