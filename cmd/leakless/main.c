#include <stdio.h>
#include <unistd.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>

int main(int argc, char *argv[]) {
	if (argc < 3) {
		printf("wrong args, usage: leakless <port> <cmd> [args...]\n");
		exit(1);
	}

	printf("parent pid: %d\n", getppid());
	
	char *port = argv[1];
	char *cmd = argv[2];

	int sockfd;
	struct sockaddr_in servaddr;

	sockfd = socket(AF_INET, SOCK_STREAM, 0);
	if (sockfd == -1) {
		printf("failed to create socket\n");
		exit(1);
	}

	bzero(&servaddr, sizeof(servaddr));
	servaddr.sin_family = AF_INET;
	servaddr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
	servaddr.sin_port = htons(atoi(port));

	if (connect(sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) != 0) {
		printf("failed to connect to server\n");
		exit(1);
	}

	char buf[1];
	read(sockfd,buf,1);
}
