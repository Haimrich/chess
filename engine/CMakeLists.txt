cmake_minimum_required(VERSION 3.9)
set (CMAKE_CXX_STANDARD 17)

set(CMAKE_CXX_FLAGS "-Werror -Wall -Wextra -pthread -g")
set(CMAKE_CXX_FLAGS_DEBUG "-Og -DDEBUG")
set(CMAKE_CXX_FLAGS_RELEASE "-O3")

enable_language(CXX)
set(CMAKE_INTERPROCEDURAL_OPTIMIZATION TRUE)
cmake_policy(SET CMP0069 NEW)
include(CheckIPOSupported)
check_ipo_supported()

set(CMAKE_INSTALL_PREFIX "/home/enrico/Desktop/Engine")
set(BUILD_SHARED_LIBS FALSE)

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Release)
endif()

project(engine)

include(GNUInstallDirs)

file(GLOB SOURCES CONFIGURE_DEPENDS "${PROJECT_SOURCE_DIR}/src/*.?pp" )

add_executable(${PROJECT_NAME} ${SOURCES} )

find_package(PkgConfig)
pkg_check_modules(Pistache REQUIRED IMPORTED_TARGET libpistache)
target_link_libraries(${PROJECT_NAME} PRIVATE PkgConfig::Pistache)

find_package(RapidJSON)

install(TARGETS ${PROJECT_NAME} RUNTIME DESTINATION "${PROJECT_SOURCE_DIR}/bin" )
install(TARGETS ${PROJECT_NAME} RUNTIME DESTINATION "${PROJECT_SOURCE_DIR}/.." )
