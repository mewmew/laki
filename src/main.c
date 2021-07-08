#include "events.h"
#include "vk.h"
#include "window.h"

#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
	GLFWwindow *win = init_window();
	VkInstance *instance = init_vulkan();

	check_extensions();
	if (enable_validation_layers) {
		if (!check_validation_layers()) {
			fprintf(stderr, "missing requested validation layer.\n");
			exit(EXIT_FAILURE);
		}
	}

	event_loop(win);

	cleanup_vulkan(instance);
	cleanup_window(win);

	return 42;
}
