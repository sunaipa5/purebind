#ifndef STRESS_TEST_H
#define STRESS_TEST_H

#include <stddef.h>
#include <stdint.h>

typedef int Counter;
typedef Counter UltraCounter;

typedef struct {
  int type;
  float func;
  char *range;
  int select;
} KeywordTest;

typedef struct {
  float transform[16];
  uint8_t shadow_map[1024];
  char tags[8][32];
} ArrayTest;

typedef struct {
  void (*on_event)(int event_id, void *user_data);
  int (*calculate)(float a, float b);
} CallbackTest;

typedef enum {
  STATUS_OK = 0,
  STATUS_ERROR = -1,
  STATUS_PENDING,
  STATUS_UNKNOWN = 99
} SystemStatus;

typedef struct Texture {
  unsigned int id;
  int width;
  int height;
  int mipmaps;
  int format;
} Texture;
typedef Texture Texture2D;
typedef Texture TextureCubemap;

typedef struct {
  char a;
  int b;
  char c;
  double d;
} PaddingTest;

void ComplexFunction(int, float *values, const char *name, Texture2D tex);

#endif
