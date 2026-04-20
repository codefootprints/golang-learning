# Gunakan image Go yang ringan
FROM golang:1.26-alpine

# Tentukan folder kerja
WORKDIR /app

# Copy file 
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber
COPY . .

# Build aplikasi menjadi file binary (exec)
RUN go build -o main .

# Jalankan aplikasinya
CMD ["./main"]