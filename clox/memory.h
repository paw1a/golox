#ifndef __MEMORY_H__
#define __MEMORY_H__

#include "stdlib.h"

#define GROW_CAPACITY(old_capacity) \
    ((old_capacity) < 8 ? 8 : (old_capacity) * 2)

#define GROW_ARRAY(type, pointer, new_len) \
    (type *)reallocate(pointer, new_len * sizeof(type))

void *reallocate(void *pointer, size_t new_size);

#endif