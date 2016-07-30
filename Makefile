all:
		go build -o arkor .
		make -C dataserver

install:
		cp arkor /usr/local/bin/arkor
clean:
		go clean
		@rm -rf dataserver/*.o
	    @rm -rf dataserver/spy_server