# Shared CGO Library

The shared functionality for the export tool is written in go to take advantage of existing libraries. 

To make sure it can be easily integrated with a myriad of other toolkits, a C interface is exported.

## Notes

### Go Pointers
Due to the nature of CGo, Go pointers can't be shared across the C boundary. We instead return a handle to a global map
which contains the active instance pointers. See [internal/handle.go](internal/handle.go) for more details.

Do note that the current implementation is not designed to be usable in all possible scenarios. It's tailored
to the needs for the export tool. More specifically, freed handles are not reused. If one were to allocate
and deallocate many small instances, the handle array would grow in size for each new instance.

If such a scenario is desirable, then handles should be replaced with a something similar to a generational handle map.

Finally, the handles are encoded into the C pointers returned by the API, so we must make sure never to return a 0 
handle as the null pointer checks would fail.

### Memory Allocation

It's recommended to free all allocated memory by this shared library(strings, arrays, etc..) using the `etFree` function 
or other equivalents, to prevent issues on windows, where the CGO shared library is built with mingw and the remaining 
code is built with MSVC.

### Callbacks/VTables

Since CGO can't call function pointers, if you need/want to simulate callback, you to write a C function that accepts 
your vtable/function pointer and calls the function pointer in question. See 
[this header](cgo_headers/etexport_mail_impl.h) for an example.

### Context Cancelled

Since there's no "special" cancelled state in C, this situation needs to be handled internally and communicated via a 
special return value. 