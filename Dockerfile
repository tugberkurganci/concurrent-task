# Alpine Linux tabanlı Docker imajını temel al
FROM golang:alpine as builder

# Çalışma dizinini /app olarak belirle
WORKDIR /app

# Docker ana dizinindeki tüm dosyaları /app dizinine kopyala
COPY . .

# Modüllerin tutarlılığını sağlamak için go mod tidy komutunu çalıştır
RUN go mod tidy

# Uygulamayı derle ve main adında bir dosya oluştur
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demory main.go

RUN chmod +x demory
RUN touch /app/loggerx/logfile.txt && chmod 666 /app/loggerx/logfile.txt

FROM scratch
COPY --from=builder /app/demory /demory
ENTRYPOINT ["/demory"]


