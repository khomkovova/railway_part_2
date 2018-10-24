#include <stdio.h>
int main(int argc, char *argv[]) {
	char firstswitch = argv[1][0];
	char secondswitch = argv[1][1];
	if ((firstswitch == '0') && (secondswitch == '0')){
        printf("True");
        return 0;
	}
	printf("False");
	return 0;
}