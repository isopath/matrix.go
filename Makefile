BIN := bin
TARGET := $(BIN)/asmatrix
OBJ := $(BIN)/asmatrix.o

# Detect operating system
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
    # macOS settings
    ASM_FORMAT := macho64
    ASM_FLAGS := --prefix _ -DMACOS
    LDFLAGS := -arch x86_64
else
    # Linux settings
    ASM_FORMAT := elf64
    ASM_FLAGS :=
    LDFLAGS := -no-pie
endif

all: $(TARGET)

$(TARGET): $(OBJ)
	gcc $(LDFLAGS) $(OBJ) -o $(TARGET) -lncurses

$(OBJ): asmatrix.asm
	mkdir -p $(BIN)
	nasm -f $(ASM_FORMAT) $(ASM_FLAGS) asmatrix.asm -o $(OBJ)

clean:
	rm -rf $(BIN)
