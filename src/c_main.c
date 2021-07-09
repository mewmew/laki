// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

#include "app.h"
#include "events.h"
#include "vk.h"
#include "window.h"

#include <stdlib.h>

void c_main() {
	App *app = calloc(1, sizeof(App));
	app->win = init_window();
	init_vulkan(app);

	event_loop(app->win);

	cleanup_vulkan(app);
	cleanup_window(app->win);
}
