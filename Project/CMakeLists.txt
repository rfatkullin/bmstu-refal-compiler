project(Project)
cmake_minimum_required(VERSION 2.8)
aux_source_directory(. SRC_LIST)

add_executable(${PROJECT_NAME} ${SRC_LIST})

#Определяем родительскую папку
get_filename_component(PARENT_DIR ${PROJECT_SOURCE_DIR} PATH)

#Подключаем h файлы рантайма
include_directories(${PARENT_DIR}/Runtime/headers)

#Указываем путь до статической библиотеки рантайма
target_link_libraries(${PROJECT_NAME} ${PARENT_DIR}/Runtime-build/libruntime.a)

