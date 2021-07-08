#include "window.h"
#include "events.h"
#include "vk.h"

int main(int argc, char **argv) {
	GLFWwindow *win = init_window();

	check_vk_extensions();

	event_loop(win);

	cleanup_window(win);

	return 42;
}
