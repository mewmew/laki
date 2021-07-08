#include "events.h"
#include "vk.h"
#include "window.h"

int main(int argc, char **argv) {
	GLFWwindow *win = init_window();
	init_vulkan();

	check_vk_extensions();

	event_loop(win);

	cleanup_window(win);

	return 42;
}
