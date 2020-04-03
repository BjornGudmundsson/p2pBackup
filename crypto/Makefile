setFile = set.txt
secretKey = "30944f2b07b37b792b4484a3b191447936bfa598bb84770f0da1a5019df5eee4"
key1 = "5c18b6874b243aeffbc6f6bead03f34e37439328c65586a440995030b112909c"
key2 = "cf70bbc0d670225d0ec79d58ebb1dc0c86438798a699062a49b34851775b7741"
mkSet = touch $(setFile);
addToSet = echo "$(key1)" >> "$(setFile)"; echo "$(key2)" >> "$(setFile)";
gTest = go test -v -tags=vartime;
clear = rm -rf $(setFile);


test:
	$(mkSet) $(addToSet) $(gTest) $(clear)