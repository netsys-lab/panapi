CC       = go build
BUILDDIR = ./bin
PRGS = server client #test

all: $(PRGS)

.PHONY: server
server:
	$(CC) -o $(BUILDDIR)/$@ cmd/server/*.go
	cp -f $(BUILDDIR)/$@ /home/ubuntu/Documents/API

.PHONY: client
client:
	$(CC) -o $(BUILDDIR)/$@ cmd/client/*.go
	cp -f $(BUILDDIR)/$@ /home/ubuntu/Documents/ATN/L3
	
.PHONY: test
test:
	$(CC) -o $(BUILDDIR)/$@ cmd/test/*.go

clean:
	rm -f $(BUILDDIR)/*
