FROM alpine:3.6

# Copy executable
ADD /build/linux_amd64/convoy /

EXPOSE 3500

CMD ["/convoy"]