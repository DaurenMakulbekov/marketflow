CC=go
CFLAGS=build
SOURCES=./cmd
TARGET=marketflow

all:
	$(CC) $(CFLAGS) -o $(TARGET) $(SOURCES)

clean:
	rm $(TARGET)

fclean: clean

re: fclean all
