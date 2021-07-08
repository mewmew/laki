#include "events.h"
#include "vk.h"
#include "window.h"

#include <stdlib.h>

int main(int argc, char **argv) {
	GLFWwindow *win = init_window();
	VkInstance *instance = init_vulkan();

	event_loop(win);

	cleanup_vulkan(instance);
	cleanup_window(win);

	return EXIT_SUCCESS;
}
