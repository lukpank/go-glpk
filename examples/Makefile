CFLAGS = -std=c99
LDFLAGS = -lglpk

all: glpk_example_c glpk_example

glpk_example_c: glpk_example.c
	gcc $(CFLAGS) $(LDFLAGS) $< -o $@

glpk_example: glpk_example.go
	go build $<

clean:
	rm -f glpk_example glpk_example_c
