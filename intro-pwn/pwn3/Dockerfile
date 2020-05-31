# docker build -t pwn3 . && docker run -p 1024:1024 --rm -it pwn3

FROM ubuntu:19.10

RUN apt-get update

RUN useradd -d /home/ctf/ -m -p ctf -s /bin/bash ctf
RUN echo "ctf:ctf" | chpasswd

WORKDIR /home/ctf

COPY pwn3 .
COPY flag .
COPY ynetd .

RUN chmod +x ynetd pwn3
RUN chown -R root:root /home/ctf

USER ctf
CMD ./ynetd ./pwn3