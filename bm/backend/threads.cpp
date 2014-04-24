/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

#include <vector>

#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/select.h>
#include <sys/socket.h>
#include <unistd.h>

struct nfd
{
	int n, fd;

	nfd(int n, int fd): n(n), fd(fd) {}
	void dbg(const char *s) {
#ifdef DEBUG
		printf("%s thrd %d for sock %d\n", s, n, fd);
#endif
	}
};

void respond(int conn)
{
	write(conn, "HTTP/1.1 200 OK\r\n", 17);
	write(conn, "Content-Type: text/plain; charset=UTF-8\r\n", 41);
	write(conn, "Content-Length: 14\r\n", 20);
	write(conn, "Date: Thu, 1 Jan 1970 00:00:00 GMT\r\n", 36);
	write(conn, "Server-Status: OK\r\n", 19);
	write(conn, "\r\n", 2);
	write(conn, "Hello from C!\n", 14);
}

void *serve(void *args)
{
	auto nfd = (struct nfd*) args;

	nfd->dbg("launched");

	int bytes;
	char request[4096];

	do {
		memset(request, 0, 4096);

		bytes = read(nfd->fd, request, 4096);
		if (bytes < 0) {
			perror("read ");
		}

		if (bytes <= 0) {
			close(nfd->fd);
			nfd->dbg("exiting");
			pthread_exit(args);
		}

		if (strnstr(request, "\r\n\r\n", bytes) != 0) {
			respond(nfd->fd);
		}
	} while(1);
}

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

	for (int n = 0; ;n++) {
		int conn = accept(server, (struct sockaddr*)0, 0);

		auto thread = new pthread_t;
		pthread_create(thread, 0, serve, new nfd(n, conn));
	}
}
