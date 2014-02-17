go-translate-cli
================

a cli tool for translation between Chinese and English


##install:
1. you'd be sure to have go 1.2(golang) installed, and put $GOPATH/bin in
your $PATH 

2. go get github.com/rnoldo/go-translate-cli  
   or  
	* git clone github.com/rnoldo/go-translate-cli.git  
	* cd go-translate-cli  
	* go install

##usage:

translate hello --> 你好  

translate 你好 --> hello

translate "你好 golang" --> hello golang  

translate "hello golang" --> 你好 golang

##and you can do more:
if you are familar with go, you can fork it and hack on it to translate between any language pair supported by translate.google.com.
