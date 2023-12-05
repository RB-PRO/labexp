all: run

run:
	go run cmd/main/main.go

push:
	git push git@github.com:RB-PRO/labexp.git

pull:
	git pull git@github.com:RB-PRO/labexp.git

pushW:
	git push https://github.com/RB-PRO/labexp.git

pullW:
	git pull https://github.com/RB-PRO/labexp.git

pushCar:
	scp main root@194.87.107.129:go/labexp/

build-config:
	go env GOOS GOARCH

build-windows-to-linux:
	set GOARCH=amd64 set GOOS=linux go build cmd/main/main.go  

build-linux-to-windows:
	export GOARCH=amd64 export GOOS=windows go build cmd/main/main.go 

build-car:
	set GOARCH=amd64
	set GOOS=linux
	set CGO_ENABLED=0
	go env GOOS GOARCH
	go build cmd/main/main.go
	scp main keys.json tg.json root@194.87.107.129:go/labexp/

build-carW:
	set GOARCH=amd64
	set GOOS=windows 
	go env GOOS GOARCH
	go build -o main.exe cmd/main/main.go
	scp main keys.json tg.json root@194.87.107.129:go/labexp/