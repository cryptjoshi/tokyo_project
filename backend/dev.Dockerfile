FROM golang:1.23.0-alpine3.20
RUN apk add --update git
RUN apk add --update curl

# ติดตั้ง tzdata เพื่อให้สามารถตั้งค่า Time Zone ได้
RUN apk add --no-cache tzdata

# คัดลอกไฟล์โซนเวลาที่ต้องการ (เช่น Asia/Bangkok)
ENV TZ=Asia/Bangkok

WORKDIR /app

#RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
#    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air


RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

# ติดตั้ง swag เป็น global tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# สร้างเอกสาร Swagger


RUN go mod download

#RUN swag init

#CMD ["air"]
#EXPOSE 3001

CMD ["air", "-c", ".air.toml"]
