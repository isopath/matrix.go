BIN := bin
TARGET := $(BIN)/asmatrix
OBJ := $(BIN)/asmatrix.o

UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
    # macOS
    ASM_FORMAT := macho64
    ASM_FLAGS := --prefix _ -DMACOS
    LDFLAGS := -arch x86_64
    LDLIBS := -lncurses
else ifneq (,$(findstring MINGW,$(UNAME_S)))
    # Windows (MinGW)
    ASM_FORMAT := elf64
    ASM_FLAGS := -DWINDOWS
    LDFLAGS := -no-pie
    LDLIBS := -lpdcurses
else
    # Linux
    ASM_FORMAT := elf64
    ASM_FLAGS :=
    LDFLAGS := -no-pie
    LDLIBS := -lncurses
endif

all: $(TARGET)

$(TARGET): $(OBJ)
	gcc $(LDFLAGS) $(OBJ) -o $(TARGET) $(LDLIBS)

$(OBJ): asmatrix.asm
	mkdir -p $(BIN)
	nasm -f $(ASM_FORMAT) $(ASM_FLAGS) asmatrix.asm -o $(OBJ)

clean:
	rm -rf $(BIN)
