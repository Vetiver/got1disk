
# T1Disk Go Library

**T1Disk Go Library** предоставляет интерфейс для работы с файловым хранилищем T1 Disk, включая авторизацию, получение ссылки для загрузки файлов и подтверждение успешной загрузки.

### Шаг 1: Авторизация

Прежде чем загружать файлы, необходимо авторизоваться и получить токен. Для этого используется функция `Login()`.

Пример:

```go
t1disk := got1disk.NewT1Disk(true, logger, "https://base-url-t1.com")
token, err := t1disk.Login("your-email", "your-password")
if err != nil {
    fmt.Println("Ошибка при авторизации:", err)
    return
}
fmt.Println("Токен:", token)
```

### Шаг 2: Загрузка файла

Для загрузки файла используется функция `UploadToT1Disk()`. Она выполняет все шаги загрузки:

* Получение ссылки для загрузки файла;
* Загрузка файла;
* Подтверждение успешной загрузки.

Пример:

```go
filePath := "/path/to/your/file.mp4"        // Путь к файлу на локальной машине
pathOnT1Disk := "T1Disk/timers/10_hours.mp4" // Путь, куда файл будет загружен на T1 Disk

err := t1disk.UploadToT1Disk(pathOnT1Disk, "https://base-url-t1.com", filePath, token, false)
if err != nil {
    fmt.Println("Ошибка при загрузке файла:", err)
    return
}
fmt.Println("Файл успешно загружен!")
```

### API

1. **`NewT1Disk(permanentAuth bool, logger *zap.Logger, BaseUrlT1 string) *T1DiskConstructor`**

   Конструктор для создания нового экземпляра T1Disk.

   **Параметры:**

   * `permanentAuth`: флаг для постоянной авторизации (true/false);
   * `logger`: логгер для логирования операций;
   * `BaseUrlT1`: базовый URL для API T1 Disk.
2. **`Login(login string, password string) (string, error)`**

   Авторизация в системе и получение токена доступа.

   **Возвращает:**

   * `string`: токен авторизации;
   * `error`: ошибка, если авторизация не удалась.
3. **`GetUploadURL(path string, baseURL string, token string, multipart bool) (string, string, string, error)`**

   Получение URL для загрузки файла в T1 Disk.

   **Параметры:**

   * `path`: путь на T1 Disk, куда будет загружен файл;
   * `baseURL`: базовый URL API;
   * `token`: токен авторизации;
   * `multipart`: флаг для использования многочастной загрузки.

     **Возвращает:**
   * URL для загрузки;
   * тип содержимого (`Content-Type`);
   * URL для подтверждения загрузки;
   * ошибка, если не удалось получить URL.
4. **`UploadFile(uploadURL string, contentType string, filePath string, token string) error`**

   Загрузка файла по URL.

   **Параметры:**

   * `uploadURL`: URL для загрузки файла;
   * `contentType`: тип содержимого файла;
   * `filePath`: путь к файлу на локальной машине;
   * `token`: токен авторизации.

     **Возвращает:**
   * ошибка, если загрузка не удалась.
5. **`ConfirmUpload(confirmURL string, token string) error`**

   Подтверждение успешной загрузки файла.

   **Параметры:**

   * `confirmURL`: URL для подтверждения загрузки;
   * `token`: токен авторизации.

     **Возвращает:**
   * ошибка, если подтверждение не удалось.
6. **`UploadToT1Disk(path string, baseURL string, filePath string, token string, multipart bool) error`**

   Полный процесс загрузки файла в T1 Disk, включая получение URL для загрузки, саму загрузку и подтверждение.

   **Параметры:**

   * `path`: путь на T1 Disk;
   * `baseURL`: базовый URL API;
   * `filePath`: путь к файлу на локальной машине;
   * `token`: токен авторизации;
   * `multipart`: флаг для многочастной загрузки.

     **Возвращает:**
   * ошибка, если что-то пошло не так.

### Пример использования

```go
package main

import (
	"fmt"
	"go.uber.org/zap"
	"got1disk"
)

func main() {
	logger, _ := zap.NewProduction()

	// Создаем новый экземпляр клиента
	t1disk := got1disk.NewT1Disk(true, logger, "https://base-url-t1.com")

	// Шаг 1: Авторизация
	token, err := t1disk.Login("your-email", "your-password")
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}
	fmt.Println("Токен авторизации:", token)

	// Шаг 2: Загрузка файла
	filePath := "/path/to/your/file.mp4"
	pathOnT1Disk := "T1Disk/timers/10_hours.mp4"
	err = t1disk.UploadToT1Disk(pathOnT1Disk, "https://base-url-t1.com", filePath, token, false)
	if err != nil {
		fmt.Println("Ошибка при загрузке файла:", err)
		return
	}
	fmt.Println("Файл успешно загружен!")
}
```
