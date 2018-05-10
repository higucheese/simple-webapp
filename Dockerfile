FROM heroku/go
MAINTAINER thiguchi <t.higuchi.eeic@gmail.com>
RUN echo "now building..."
ADD ./main.go .
RUN ["go", "get", "github.com/lib/pq"]
RUN ["go", "build", "-o", "webapp", "main.go"]
EXPOSE 80
CMD ["go", "run", "main.go"]
#CMD ["sleep", "5"]
#CMD ["./webapp"]
