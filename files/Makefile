tDir = test
dir = mkdir test;
rDir = rm -rf $(tDir)

tDir2 = test2
dir2 = mkdir $(tDir2);
test2 = $(dir2) mkdir $(tDir2)/$(tDir2); touch $(tDir2)/t1 $(tDir2)/$(tDir2)/t2 $(tDir2)/$(tDir2)/t3;
clearTest2 = rm -rf $(tDir2);
append = head -c 64 </dev/urandom > logFile.txt;
make_reconstruct = mkdir $(reconstruct_dir)
mockFiles = touch $(tDir)/f1 $(tDir)/f2 $(tDir)/f3;
gTest = go test -v -tags=vartime;
rmFiles = rm -rf $(tDir);
fillFile = head -c 50000 </dev/random > appendFull.txt
fillFile2 = head -c 50000 </dev/random > segmentBuffer.txt
writeTest = touch writeFile.txt; touch appendFile.txt; touch logFile.txt; touch appendOnly.txt; touch appendFull.txt; touch segmentBuffer.txt; $(fillFile);$(fillFile2);
endWriteTest = rm writeFile.txt appendFile.txt logFile.txt appendOnly.txt appendFull.txt segmentBuffer.txt;
read = cat appendFile.txt;
removeReconstruction = rm -rf reconstructed
test: 
	$(dir) $(mockFiles) $(test2) $(writeTest) $(append) $(gTest) $(rmFiles) $(clearTest2) $(endWriteTest) $(removeReconstruction)