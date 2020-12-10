CC       = go build
BUILDDIR = ./bin
PRGS = server client #test

all: $(PRGS)

.PHONY: server
server:
	$(CC) -o $(BUILDDIR)/$@ cmd/server/*.go

.PHONY: client
client:
	$(CC) -o $(BUILDDIR)/$@ cmd/client/*.go
	
.PHONY: test
test:
	$(CC) -o $(BUILDDIR)/$@ cmd/test/*.go

clean:
	rm -f $(BUILDDIR)/*
