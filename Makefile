BIN := bin
TARGET := $(BIN)/asmatrix
OBJ := $(BIN)/asmatrix.o

all: $(TARGET)

$(TARGET): $(OBJ)
	gcc -no-pie $(OBJ) -o $(TARGET) -lncurses

$(OBJ): asmatrix.asm
	mkdir -p $(BIN)
	nasm -f elf64 asmatrix.asm -o $(OBJ)

clean:
	rm -rf $(BIN)
