#ifndef __CHUNK_H__
#define __CHUNK_H__

#include "common.h"
#include "memory.h"

typedef enum {
    OP_RETURN,
} Opcode;

typedef struct {
    int len;
    int capacity;
    uint8_t *code;
} Chunk;

void init_chunk(chunk *Chunk);
void append_chunk(Chunk *chunk, uint8_t code);
void free_chunk(Chunk *chunk);

#endif
