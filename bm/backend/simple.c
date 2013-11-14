#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
	unsigned short port;
	if (sscanf(argv[1], "%hd", &port) < 0) {
		exit(255);
	}

	int server = socket(PF_INET, SOCK_STREAM, 0);
	if (server < 0) {
		perror("socket ");
	}

	struct sockaddr_in addr;
	memset(&addr, 0, sizeof(addr));

	addr.sin_family = PF_INET;
	addr.sin_addr.s_addr = htonl(INADDR_ANY);
	addr.sin_port = htons(port);

	if (bind(server, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
		perror("bind ");
	}

	if (listen(server, 128) < 0) {
		perror("listen ");
	}

	while (1) {
		int conn = accept(server, (struct sockaddr*)0, 0);

		char req[4096];
		read(conn, req, sizeof(req));

		write(conn, "HTTP/1.1 200 OK\r\n", 17);
		write(conn, "Content-Type: text/plain; charset=UTF-8\r\n", 41);
		write(conn, "Content-Length: 14\r\n", 20);
		write(conn, "Connection: close\r\n", 19);
		write(conn, "Date: Thu, 1 Jan 1970 00:00:00 GMT\r\n", 36);
		write(conn, "Server-Status: OK\r\n", 19);
		write(conn, "\r\n", 2);
		write(conn, "Hello from C!\n", 14);

		close(conn);
	}
}
