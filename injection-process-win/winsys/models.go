package winsys

import "syscall"

type Inject struct {
	Pid              uint32
	DllPath          string
	DLLSize          uint32
	Privilege        string
	RemoteProcHandle uintptr
	Lpaddr           uintptr
	LoadLibAddr      uintptr
	RThread          uintptr
	Token            TOKEN
}
type TOKEN struct {
	tokenHandle syscall.Token
}
