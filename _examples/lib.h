#ifndef EXAMPLE_H
#define EXAMPLE_H

#include <stdint.h>

#define MAX_BUFFER_SIZE 1024

extern const double PI;

typedef enum { COLOR_RED, COLOR_GREEN, COLOR_BLUE } Color;

typedef struct {
  int id;
  char name[50];
  float score;
} Student;

int add(int a, int b);
void print_message(const char *msg);
double compute_area(double radius);
Color get_favorite_color(void);

void fill_buffer(char *buffer, size_t size);

typedef void (*Callback)(int event_code);

void register_callback(Callback cb);

extern int global_counter;

#endif
