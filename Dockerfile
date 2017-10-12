FROM scratch

# Copy executable
ADD /build/linux_amd64/convoy /

EXPOSE 8502

CMD ["/convoy"]