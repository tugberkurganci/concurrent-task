# Alpine Linux tabanlı Docker imajını temel al
FROM golang:alpine as builder

# Çalışma dizinini /app olarak belirle
WORKDIR /app

# Docker ana dizinindeki tüm dosyaları /app dizinine kopyala
COPY . .

# Modüllerin tutarlılığını sağlamak için go mod tidy komutunu çalıştır
RUN go mod tidy

# loggerx dizinini oluştur ve dosyayı kopyala
RUN mkdir -p /app/loggerx
RUN touch /app/loggerx/logfile.txt && chmod 666 /app/loggerx/logfile.txt

# Uygulamayı derle ve main adında bir dosya oluştur
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o konzek main.go

RUN chmod +x konzek

# Yeni bir Alpine Linux tabanlı imaj başlat
FROM alpine:latest

# Uygulama dosyasını ve loggerx dizinini kopyala
COPY --from=builder /app/konzek /konzek
COPY --from=builder /app/loggerx /app/loggerx

# Çalışma dizinini /app olarak belirle
WORKDIR /app

# Giriş noktası belirle
ENTRYPOINT ["/konzek"]




