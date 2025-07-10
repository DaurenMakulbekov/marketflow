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


load:
	docker load -i exchange1_amd64.tar
	docker load -i exchange2_amd64.tar
	docker load -i exchange3_amd64.tar

run:
	docker run -p 40101:40101 --net marketflow_default --ip 172.22.0.5 --name exchange1 -d exchange1
	docker run -p 40102:40102 --net marketflow_default --ip 172.22.0.6 --name exchange2 -d exchange2
	docker run -p 40103:40103 --net marketflow_default --ip 172.22.0.7 --name exchange3 -d exchange3
