project(Runtime)
cmake_minimum_required(VERSION 2.8)

include_directories(${PROJECT_SOURCE_DIR}/headers)

file(GLOB RuntimeSources ${PROJECT_SOURCE_DIR}/sources/*.c)
file(GLOB RuntimeHeaders ${PROJECT_SOURCE_DIR}/headers/*.h)

set(SRC_LIST ${RuntimeSources} ${RuntimeHeaders})
add_executable(${PROJECT_NAME} ${SRC_LIST})

