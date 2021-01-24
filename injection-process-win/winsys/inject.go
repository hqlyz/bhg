package winsys

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var nullRef int

func OpenProcessHandle(i *Inject) error {
	var rights uint32 = PROCESS_CREATE_THREAD |
		PROCESS_QUERY_INFORMATION |
		PROCESS_VM_OPERATION |
		PROCESS_VM_WRITE |
		PROCESS_VM_READ
	var inheritHandle uint32 = 0
	var processID uint32 = i.Pid
	remoteProcHandle, _, lastErr := ProcOpenProcess.Call(
		uintptr(rights),        // DWORD dwDesiredAccess
		uintptr(inheritHandle), // BOOL bInheritHandle
		uintptr(processID))     // DWORD dwProcessId
	if remoteProcHandle == 0 {
		return errors.Wrap(lastErr, `[!] ERROR : Can't Open Remote Process. Maybe running elevated integrity?`)
	}
	i.RemoteProcHandle = remoteProcHandle
	fmt.Printf("[-] Input PID: %v\n", i.Pid)
	fmt.Printf("[-] Input DLL: %v\n", i.DllPath)
	fmt.Printf("[+] Process handle: %v\n", unsafe.Pointer(i.RemoteProcHandle))
	return nil
}

func VirtualAllocEx(i *Inject) error {
	var flAllocationType uint32 = MEM_COMMIT | MEM_RESERVE
	var flProtect uint32 = PAGE_EXECUTE_READWRITE
	lpBaseAddress, _, lastErr := ProcVirtualAllocEx.Call(
		i.RemoteProcHandle,        // HANDLE hProcess
		uintptr(nullRef),          // LPVOID lpAddress
		uintptr(i.DLLSize),        // SIZE_T dwSize
		uintptr(flAllocationType), // DWORD flAllocationType
		// https://docs.microsoft.com/en-us/windows/desktop/Memory/memory-protection-constants
		uintptr(flProtect)) // DWORD flProtect
	if lpBaseAddress == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Can't Allocate Memory On Remote Process.")
	}
	i.Lpaddr = lpBaseAddress
	fmt.Printf("[+] Base memory address: %v\n", unsafe.Pointer(i.Lpaddr))
	return nil
}

func WriteProcessMemory(i *Inject) error {
	var nBytesWritten *byte
	dllPathBytes, err := syscall.BytePtrFromString(i.DllPath)
	if err != nil {
		return err
	}
	writeMem, _, lastErr := ProcWriteProcessMemory.Call(
		i.RemoteProcHandle,                     // HANDLE hProcess
		i.Lpaddr,                               // LPVOID lpBaseAddress
		uintptr(unsafe.Pointer(dllPathBytes)),  // LPCVOID lpBuffer
		uintptr(i.DLLSize),                     // SIZE_T nSize
		uintptr(unsafe.Pointer(nBytesWritten))) // SIZE_T *lpNumberOfBytesWritten
	if writeMem == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Can't write to process memory.")
	}
	return nil
}

func GetLoadLibAddress(i *Inject) error {
	var llibBytePtr *byte
	llibBytePtr, err := syscall.BytePtrFromString("LoadLibraryA")
	if err != nil {
		return err
	}
	lladdr, _, lastErr := ProcGetProcAddress.Call(
		ModKernel32.Handle(),                 // HMODULE hModule
		uintptr(unsafe.Pointer(llibBytePtr))) // LPCSTR lpProcName
	if &lladdr == nil {
		return errors.Wrap(lastErr, "[!] ERROR : Can't get process address.")
	}
	i.LoadLibAddr = lladdr
	fmt.Printf("[+] Kernel32.Dll memory address: %v\n", unsafe.Pointer(ModKernel32.Handle()))
	fmt.Printf("[+] Loader memory address: %v\n", unsafe.Pointer(i.LoadLibAddr))
	return nil
}

func CreateRemoteThread(i *Inject) error {
	var threadId uint32 = 0
	var dwCreationFlags uint32 = 0
	remoteThread, _, lastErr := ProcCreateRemoteThread.Call(
		i.RemoteProcHandle,                 // HANDLE hProcess
		uintptr(nullRef),                   // LPSECURITY_ATTRIBUTES lpThreadAttributes
		uintptr(nullRef),                   // SIZE_T dwStackSize
		i.LoadLibAddr,                      // LPTHREAD_START_ROUTINE lpStartAddress
		i.Lpaddr,                           // LPVOID lpParameter
		uintptr(dwCreationFlags),           // DWORD dwCreationFlags
		uintptr(unsafe.Pointer(&threadId)), // LPDWORD lpThreadId
	)
	if remoteThread == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Can't Create Remote Thread.")
	}
	i.RThread = remoteThread
	fmt.Printf("[+] Thread identifier created: %v\n", unsafe.Pointer(&threadId))
	fmt.Printf("[+] Thread handle created: %v\n", unsafe.Pointer(i.RThread))
	return nil
}

func WaitForSingleObject(i *Inject) error {
	var dwMilliseconds uint32 = INFINITE
	var dwExitCode uint32
	rWaitValue, _, lastErr := ProcWaitForSingleObject.Call(
		i.RThread,               // HANDLE hHandle
		uintptr(dwMilliseconds)) // DWORD dwMilliseconds
	if rWaitValue != 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Error returning thread wait state.")
	}
	success, _, lastErr := ProcGetExitCodeThread.Call(
		i.RThread,                            // HANDLE hThread
		uintptr(unsafe.Pointer(&dwExitCode))) // LPDWORD lpExitCode
	if success == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Error returning thread exit code.")
	}
	closed, _, lastErr := ProcCloseHandle.Call(i.RThread) // HANDLE hObject
	if closed == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Error closing thread handle.")
	}
	return nil
}

func VirtualFreeEx(i *Inject) error {
	var dwFreeType uint32 = MEM_RELEASE
	var size uint32 = 0 //Size must be 0 to MEM_RELEASE all of the region
	rFreeValue, _, lastErr := ProcVirtualFreeEx.Call(
		i.RemoteProcHandle,  // HANDLE hProcess
		i.Lpaddr,            // LPVOID lpAddress
		uintptr(size),       // SIZE_T dwSize x
		uintptr(dwFreeType)) // DWORD dwFreeType
	if rFreeValue == 0 {
		return errors.Wrap(lastErr, "[!] ERROR : Error freeing process memory.")
	}
	fmt.Println("[+] Success: Freed memory region")
	return nil
}
