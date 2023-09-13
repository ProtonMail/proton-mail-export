#
# Patch Genated CGO header so it's compatible with MSVC C++17
#
# Usage cmake -P <scritp> <path_to_header>
#

set(header "${CMAKE_ARGV3}")

file(READ "${header}" contents)

string(REPLACE "<complex.h>" "<complex>" contents "${contents}")
string(REPLACE "_Fcomplex" "std::complex<float>" contents "${contents}")
string(REPLACE "_Dcomplex" "std::complex<double>" contents "${contents}")

file(WRITE "${header}" "${contents}")
