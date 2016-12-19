# go.counter
This code counts nubmer of "Go" on provided urls.

## Building
go.counteruses standard go building routine, checkout [golang page](https://golang.org/doc/code.html).  
Code should be located at `$GOPATH/src/github.com/wizpnz/go.counter`.
go get github.com/wizpnz/go.counter

### Step by step
0. Install 'go', see [golang installation guide](https://golang.org/doc/install).
1. Clone repository to `$GOPATH/src/github.com/wizpnz/go.counter` 
2. Build with `go install`
3. Executable will be in `bin` folder

```
cd %GOPATH/src
git clone https://github.com/wizpnz/go.counter.git
cd github.com/wizpnz/go.counter
go get
go install
```

##Run
Run from command line
```
Use -log to enable logging to stdout (stderr is default)
```
put urls to stdin

