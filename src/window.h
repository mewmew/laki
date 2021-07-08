#ifndef __WINDOW_H__
#define __WINDOW_H__

#define GLFW_INCLUDE_VULKAN
#include <GLFW/glfw3.h>

extern GLFWwindow * init_window() ;
extern void cleanup_window(GLFWwindow *win);

#endif // #ifndef __WINDOW_H__
