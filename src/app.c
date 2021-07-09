#include "app.h"

#include <stdlib.h>

App * new_app() {
	return calloc(1, sizeof(App));
}
