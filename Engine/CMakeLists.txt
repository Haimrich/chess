cmake_minimum_required(VERSION 3.6)
set (CMAKE_CXX_STANDARD 17)

set(CMAKE_CXX_FLAGS "-Werror -Wall -Wextra -pthread -g")
set(CMAKE_CXX_FLAGS_DEBUG "-Og -DDEBUG")
set(CMAKE_CXX_FLAGS_RELEASE "-O3")

set(CMAKE_INSTALL_PREFIX "/home/enrico/Desktop/Engine")

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Release)
endif()

project(engine)

include(GNUInstallDirs)

file(GLOB SOURCES CONFIGURE_DEPENDS 
        "${PROJECT_SOURCE_DIR}/src/*.?pp" 
    )

add_executable(${PROJECT_NAME} ${SOURCES})

find_package(PkgConfig)
pkg_check_modules(Pistache REQUIRED IMPORTED_TARGET libpistache)
target_link_libraries(${PROJECT_NAME} PkgConfig::Pistache)

find_package(RapidJSON)

install(TARGETS ${PROJECT_NAME} RUNTIME DESTINATION "${PROJECT_SOURCE_DIR}/bin")
install(TARGETS ${PROJECT_NAME} RUNTIME DESTINATION "${PROJECT_SOURCE_DIR}/..")
