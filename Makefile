all: linux windows macos freebsd openbsd netbsd solaris

clean:
	echo "cleaning"
	rm ./bin/* 
	rmdir ./bin
linux:	
	./build.sh "linux" "amd64"
	./build.sh "linux" "386"
windows:
	./build.sh "windows" "386"
	./build.sh "windows" "amd64"
macos:
	./build.sh "darwin" "amd64"
	./build.sh "darwin" "386"

freebsd:
	./build.sh "freebsd" "386"
	./build.sh "freebsd" "amd64"

openbsd:
	./build.sh "openbsd" "386"
	./build.sh "openbsd" "amd64"

netbsd:
	./build.sh "netbsd" "386"
	./build.sh "netbsd" "amd64"

solaris:
	./build.sh "solaris" "amd64"
