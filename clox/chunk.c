#include "chunk.h"

void init_chunk(Chunk *chunk) {
    chunk->len = 0;
    chunk->capacity = 0;
    chunk->code = NULL;
}

void append_chunk(Chunk *chunk, uint8_t code) {
    if (chunk->len >= chunk->capacity) {
        int new_capacity = GROW_CAPACITY(chunk->capacity);
        chunk->capacity = new_capacity;
        chunk->code = GROW_ARRAY(uint8_t, chunk->code, new_capacity);
    }
    chunk->code[chunk->len++] = code;
}

void free_chunk(Chunk *chunk) {
    free(chunk->code);
    init_chunk(chunk);
}
