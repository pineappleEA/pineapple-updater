cmake_minimum_required(VERSION 3.10)

#project name
project(Pineapple-Windows VERSION 0.0)

find_package(CURL REQUIRED)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED True)
add_executable(download_test download_test.cpp)
target_include_directories(download_test PUBLIC include ${CURL_INCLUDE_DIR})
target_link_options(download_test PUBLIC "LINKER:-laria2")
target_link_libraries(download_test PUBLIC ${CURL_LIBRARIES})
