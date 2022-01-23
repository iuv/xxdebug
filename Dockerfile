FROM busybox
ADD build/linux/registry-hub /app/registry-hub
ADD conf /app/conf
RUN chmod 775 /app/registry-hub
EXPOSE 80
WORKDIR /app
CMD ./registry-hub