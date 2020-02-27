build = go build
BINARY = a
port = 8081
buildFlags = -o $(BINARY)
binaryFlags = -p $(port)



MAKE = make



build:
	$(build) $(buildFlags)

run: 
	$(MAKE); ./$(BINARY) $(binaryFlags)

clean:
	rm $(BINARY)