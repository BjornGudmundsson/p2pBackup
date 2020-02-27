tDir = test
dir = mkdir test;
rDir = rm -rf $(tDir)

tDir2 = test2
dir2 = mkdir $(tDir2);
test2 = $(dir2) mkdir $(tDir2)/$(tDir2); touch $(tDir2)/t1 $(tDir2)/$(tDir2)/t2 $(tDir2)/$(tDir2)/t3;
clearTest2 = rm -rf $(tDir2);


mockFiles = touch $(tDir)/f1 $(tDir)/f2 $(tDir)/f3;
gTest = go test;
rmFiles = rm -rf $(tDir);

test: 
	$(dir) $(mockFiles) $(test2) $(gTest) $(rmFiles) $(clearTest2)