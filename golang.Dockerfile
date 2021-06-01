FROM golang:1.10-alpine
ADD ./script_patch.sh /tmp/script_patch.sh
RUN /tmp/script_patch.sh