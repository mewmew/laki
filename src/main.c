#include "app.h"
#include "events.h"
#include "vk.h"
#include "window.h"

#include <stdlib.h>

int main(int argc, char **argv) {
	App *app = calloc(1, sizeof(App));
	app->win = init_window();
	init_vulkan(app);

	event_loop(app->win);

	cleanup_vulkan(app);
	cleanup_window(app->win);

	return EXIT_SUCCESS;
}
