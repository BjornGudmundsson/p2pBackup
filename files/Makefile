tDir = test
dir = mkdir test;
rDir = rm -rf $(tDir)

tDir2 = test2
dir2 = mkdir $(tDir2);
test2 = $(dir2) mkdir $(tDir2)/$(tDir2); touch $(tDir2)/t1 $(tDir2)/$(tDir2)/t2 $(tDir2)/$(tDir2)/t3;
clearTest2 = rm -rf $(tDir2);


mockFiles = touch $(tDir)/f1 $(tDir)/f2 $(tDir)/f3;
gTest = go test -v -tags=vartime;
rmFiles = rm -rf $(tDir);
writeTest = touch writeFile.txt; touch appendFile.txt; touch logFile.txt;
endWriteTest = rm writeFile.txt appendFile.txt logFile.txt;
read = cat appendFile.txt;
test: 
	$(dir) $(mockFiles) $(test2) $(writeTest) $(gTest) $(rmFiles) $(clearTest2) $(endWriteTest)