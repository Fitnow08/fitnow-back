# MinIO — подключение и использование

## 1. Запуск через Docker

```bash
docker run -d \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

| Адрес | Назначение |
|---|---|
| `localhost:9000` | S3 API (для Go SDK) |
| `localhost:9001` | Веб-консоль |

Открой `http://localhost:9001` → войди с `minioadmin / minioadmin` → создай bucket (например `avatars`).

---

## 2. Установка Go SDK

```bash
go get github.com/minio/minio-go/v7
```

---

## 3. Конфиг

Добавь в `.env` / конфиг приложения:

```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=avatars
```

---

## 4. Клиент

Создай файл `pkg/storage/minio.go`:

```go
package storage

import (
    "context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
    client *minio.Client
    bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        return nil, err
    }

    exists, err := client.BucketExists(context.Background(), bucket)
    if err != nil {
        return nil, err
    }
    if !exists {
        if err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{}); err != nil {
            return nil, err
        }
    }

    return &MinioClient{client: client, bucket: bucket}, nil
}
```

---

## 5. Загрузка файла

```go
func (m *MinioClient) Upload(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error) {
    _, err := m.client.PutObject(ctx, m.bucket, objectName, file, size, minio.PutObjectOptions{
        ContentType: contentType,
    })
    if err != nil {
        return "", err
    }
    // вернуть публичный URL
    return fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucket, objectName), nil
}
```

---

## 6. Получение файла (presigned URL)

Если bucket приватный — генерируй временную ссылку:

```go
func (m *MinioClient) PresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
    url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, expiry, nil)
    if err != nil {
        return "", err
    }
    return url.String(), nil
}
```

---

## 7. Удаление файла

```go
func (m *MinioClient) Delete(ctx context.Context, objectName string) error {
    return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}
```

---

## 8. Пример хендлера загрузки аватара

```go
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // 10 MB max

    file, header, err := r.FormFile("avatar")
    if err != nil {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, api.Error("no file provided"))
        return
    }
    defer file.Close()

    objectName := fmt.Sprintf("users/%s%s", uuid.New().String(), filepath.Ext(header.Filename))

    url, err := h.storage.Upload(r.Context(), objectName, file, header.Size, header.Header.Get("Content-Type"))
    if err != nil {
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, api.Error("failed to upload file"))
        return
    }

    render.Status(r, http.StatusOK)
    render.JSON(w, r, map[string]string{"url": url})
}
```

---

## 9. Публичный доступ к bucket (опционально)

Чтобы файлы были доступны без авторизации, установи политику через веб-консоль:

`Buckets → avatars → Access Policy → Public`

Или через код:

```go
policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::avatars/*"]}]}`
client.SetBucketPolicy(ctx, "avatars", policy)
```
