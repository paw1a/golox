#include "memory.h"

#define ERR_MEMORY_FATAL 10

void *reallocate(void *pointer, size_t new_size) {
    if (new_size == 0) {
        free(pointer);
        return NULL;
    }

    void *result = realloc(pointer, new_size);
    if (result == NULL)
        exit(ERR_MEMORY_FATAL);

    return result;
}
