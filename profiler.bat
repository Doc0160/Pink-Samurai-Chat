@echo off

go tool pprof -svg http://localhost:8080/debug/pprof/trace?seconds=5 > "profile/trace.svg"
go tool pprof -svg http://localhost:8080/debug/pprof/heap > "profile/heap.svg"
go tool pprof -svg http://localhost:8080/debug/pprof/block > "profile/block.svg"
go tool pprof -svg http://localhost:8080/debug/pprof/goroutine > "profile/goroutine.svg"
go tool pprof -svg http://localhost:8080/debug/pprof/threadcreate > "profile/threadcreate.svg"
go tool pprof -svg http://localhost:8080/debug/pprof/profile > "profile/profile.svg"

:exit
