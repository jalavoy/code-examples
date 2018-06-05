// jlavoy - 6-20-2017
// will delete the contents of a directory that has grown so large that traditional removal tools no longer work
#define BUF_SIZE 1024*1024*5
#define _GNU_SOURCE
#include <dirent.h>     /* Defines DT_* constants */
#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#define handle_error(msg) \
do { perror(msg); exit(EXIT_FAILURE); } while (0)

struct linux_dirent {
    long           d_ino;
    off_t          d_off;
    unsigned short d_reclen;
    char           d_name[];
};

struct stat statbuf;

int main(int argc, char *argv[]) {
    int fd, nread;
    char buf[BUF_SIZE];
    struct linux_dirent *d;
    int bpos;
    int seconds;
    time_t now = time(NULL);
    if ( argc != 3 ) {
        printf("Usage: %s <directory> <days old>\nEx: %s /tmp 7\nto delete things that are older than 7 days\n", argv[0], argv[0]);
    } else {
        seconds = atoi(argv[2]) * 86400;
        chdir(argv[1]);
        fd = open(".", O_RDONLY|O_DIRECTORY);
        if ( fd == -1 ) {
            handle_error("open");
        }
        for ( ; ; ) {
            nread = syscall(SYS_getdents, fd, buf, BUF_SIZE);
            if ( nread == -1 ) {
                handle_error("get_dents");
            } else if ( nread == 0 ) {
                break;
            }
            for ( bpos = 0 ; bpos < nread ; ) {
                d = (struct linux_dirent *) (buf + bpos);
                if ( stat(d->d_name, &statbuf) == -1 ) {
                    continue;
                }
                if (S_ISREG(statbuf.st_mode)) {
                    if ( now - statbuf.st_mtime >= seconds ) {
                        if ( unlink(d->d_name) == 0 ) {
                            printf("Deleting [%s/%s] - %d\n", argv[1], d->d_name, statbuf.st_mtime);
                        } else {
                            printf("Unable to delete [%s/%s]!\n", argv[1], d->d_name);
                        }
                    }
                }
                bpos += d->d_reclen;
            }
        }
    }
    exit(EXIT_SUCCESS);
}
