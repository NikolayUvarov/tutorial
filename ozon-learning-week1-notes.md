# net/http в разрезе RoundTripper

На примерах

    package main
    
    import (
    	"fmt"
    	"net/http"
    	"time"
    )
    
    // Логирующий транспорт с измерением времени ответа
    type LoggingTransport struct {
    	Transport http.RoundTripper // Вложенный транспорт
    }
    
    func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    	start := time.Now()
    	resp, err := t.Transport.RoundTrip(req) // Выполняем запрос
    	duration := time.Since(start)
    
    	if err != nil {
    		fmt.Printf("Request: %s %s | Error: %v | Duration: %v\n", req.Method, req.URL, err, duration)
    		return nil, err
    	}
    
    	fmt.Printf("Request: %s %s | Duration: %v | Status: %d\n", req.Method, req.URL, duration, resp.StatusCode)
    	return resp, err
    }
    
    // Транспорт, добавляющий заголовки
    type HeaderTransport struct {
    	Transport http.RoundTripper
    }
    
    func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    	start := time.Now()
    	req.Header.Set("X-Custom-Header", "MyValue") // Добавляем заголовок
    	resp, err := t.Transport.RoundTrip(req)      // Выполняем запрос
    	duration := time.Since(start)
    
    	if err != nil {
    		fmt.Printf("Request with header: %s %s | Error: %v | Duration: %v\n", req.Method, req.URL, err, duration)
    		return nil, err
    	}
    
    	fmt.Printf("Request with header: %s %s | Duration: %v | Status: %d\n", req.Method, req.URL, duration, resp.StatusCode)
    	return resp, err
    }
    
    // Фейковый транспорт для тестов
    type FakeTransport struct{}
    
    func (t *FakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    	start := time.Now()
    	time.Sleep(50 * time.Millisecond) // Симулируем задержку ответа
    	duration := time.Since(start)
    
    	fmt.Printf("Fake request: %s %s | Duration: %v | Status: %d\n", req.Method, req.URL, duration, http.StatusOK)
    
    	return &http.Response{
    		StatusCode: http.StatusOK,
    		Body:       http.NoBody, // Пустое тело
    	}, nil
    }
    
    func main() {
    	// HTTP-клиент с логгированием
    	client := &http.Client{
    		Transport: &LoggingTransport{Transport: http.DefaultTransport}, // Подменяем транспорт
    	}
    	_, _ = client.Get("https://example.com")
    
    	// HTTP-клиент с кастомным заголовком
    	clientWithHeader := &http.Client{
    		Transport: &HeaderTransport{Transport: http.DefaultTransport},
    	}
    	_, _ = clientWithHeader.Get("https://example.com")
    
    	// Ненастоящий (моковый, тестовый, фейковый) HTTP-клиент для тестов
    	clientFake := &http.Client{Transport: &FakeTransport{}}
    	_, _ = clientFake.Get("https://example.com")
    }


# net/http - параметрезированный запрос

Пример ниже. На лекции предложен вызов вида http.HandleFunc("Get /post/{id}", hanlerFuncName), но он скорее всего
подойдет для гориллы-мух или для chi.

    package main
    
    import (
    	"encoding/json"
    	"fmt"
    	"log"
    	"net/http"
    
    	"github.com/gorilla/mux"
    )
    
    // Структура для данных запроса
    type User struct {
    	ID   string `json:"id"`
    	Name string `json:"name"`
    	Age  int    `json:"age"`
    }
    
    // Фейковая база данных
    var users = map[string]User{
    	"1": {ID: "1", Name: "Alice", Age: 25},
    	"2": {ID: "2", Name: "Bob", Age: 30},
    }
    
    // Обработчик GET-запроса с параметром в пути (/user/{id})
    func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    	params := mux.Vars(r) // Получаем параметры пути
    	id := params["id"]
    
    	user, exists := users[id]
    	if !exists {
    		http.Error(w, "User not found", http.StatusNotFound)
    		return
    	}
    
    	w.Header().Set("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(user)
    }
    
    // Обработчик POST-запроса (/user)
    func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    	var user User
    	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
    		http.Error(w, "Invalid request body", http.StatusBadRequest)
    		return
    	}
    
    	// Сохраняем нового пользователя
    	users[user.ID] = user
    
    	w.Header().Set("Content-Type", "application/json")
    	w.WriteHeader(http.StatusCreated)
    	json.NewEncoder(w).Encode(user)
    }
    
    func main() {
    	// Создаём роутер
    	router := mux.NewRouter()
    
    	// Роуты с параметрами пути
    	router.HandleFunc("/user/{id}", GetUserHandler).Methods("GET")   // Получение пользователя по ID
    	router.HandleFunc("/user", CreateUserHandler).Methods("POST")    // Создание нового пользователя
    
    	// Запускаем сервер
    	fmt.Println("Server is running on :8080")
    	log.Fatal(http.ListenAndServe(":8080", router))
    }

# net/http сервер с логгированием запроса


    package main
    
    import (
    	"fmt"
    	"log"
    	"net/http"
    	"strings"
    	"time"
    )
    
    // Middleware для логирования времени выполнения запроса
    func LoggingMiddleware(next http.Handler) http.Handler {
    	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    		start := time.Now() // Засекаем время начала запроса
    
    		next.ServeHTTP(w, r) // Вызываем следующий обработчик
    
    		duration := time.Since(start) // Считаем время выполнения
    		log.Printf("Request: %s %s | Duration: %v\n", r.Method, r.URL.Path, duration)
    	})
    }
    
    // Обработчик GET /post/{id}
    func GetPostHandler(w http.ResponseWriter, r *http.Request) {
    	// Разбираем URL вручную
    	parts := strings.Split(r.URL.Path, "/")
    	if len(parts) < 3 || parts[1] != "post" {
    		http.Error(w, "Invalid URL", http.StatusNotFound)
    		return
    	}
    	id := parts[2] // Извлекаем {id}
    
    	time.Sleep(100 * time.Millisecond) // Симулируем задержку обработки запроса
    	fmt.Fprintf(w, "Post ID: %s", id)
    }
    
    // Обработчик GET /
    func HomeHandler(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintln(w, "Welcome to the API!")
    }
    
    func main() {
    	// Создаём базовый маршрутизатор
    	mux := http.NewServeMux()
    
    	// Регистрируем обработчики
    	mux.HandleFunc("/", HomeHandler)
    	mux.HandleFunc("/post/", GetPostHandler)
    
    	// Добавляем middleware
    	handlerWithMiddleware := LoggingMiddleware(mux)
    
    	// Запускаем сервер
    	fmt.Println("Server running on :8080")
    	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
    }



# errors

errors.Is – проверяет, является ли ошибка целевой
errors.As – извлекает ошибку по типу

errors.Is(err, target)	Проверяет, эквивалентна ли ошибка target, даже если она обёрнута.
errors.As(err, &target)	Проверяет, является ли ошибка err нужным типом (по type assertion).

errors.Join(err1, err2) объединит ошибки вместе (цепочки ошибок), например от разных горутин


# interace
Есть функция (все функции) у объекта, описанная в интерфейсе - он реализует интерфейс. Описание не требуется.

var x interface{} = "hello"
но лучше (any - алиаск к пустому интерфейсу, удобнее с дженериками)
var x any = 42

Интерфейс можно передать в фаункцию и вызывать внутри функции функцию-метод интерфейса. 

Можно проверить тип:
var x interface{} = "hello"

    // Приведение типа
    str, ok := x.(string)
    if ok {
        fmt.Println("String:", str)
    }

# interface & generics

    func PrintSlice(slice []interface{}) {
        for _, v := range slice {
            fmt.Println(v)
        }
    }
    
    func PrintSlice[T any](slice []T) {
        for _, v := range slice {
            fmt.Println(v)
        }
    }







# map

## База

from go 1.24 используются swissmaps
Адаптивный рост (до этого удвоение) и реже рехеш.

как std::unordered_map в C++20

## Нюанс
map[int64]struct{} requires 16 bytes per slot #71368

[int64 -> bool] = 16 байт 
[int64 -> struct{}] = 16 байт (из-за выравнивания и метаданных)


go 1.23 было 8 байт
[int64 -> struct{}] = 8 байт (только ключ)
[int64 -> bool] = 9 байт (8 + 1)

потоко-опасно.

Защита от атак на хеш (когда ключи подбираются таким образом, что попадают в один бакет и из-за поиска свободного места 
программа тормозит: вместо o(1) становится o(N). SipHash+random на старте делает атаку почти невозможной, но приводит к
непредзказуемости порядка итерации по порядку (~= сортировка не определена на map)

*sync.Map* быстрее на чтениях (так как содержит два map внутри, один для чтения, другой для чтения записи). 
Потокобезопасно, но может знатно тормозить при копировании внутренних данных на больших объемах.


# Слайсы и сабслайсы в Go

При создании саб-слайса новый получчет новый capacity  от исходного cap(src) - start, len(src) - start, то есть и длина
и капасити уменьшаются на то, на сколько короче становится _от начала_ новый слайс относительно нового.

Копирования данных не происходит

slicen1 := []int{1, 2, 3, 4, 5}
subSlicen := slicen1[1:3]

При измении размера слайса, которое вызывает реаллокацию памяти (on capacity increase) дочерний слайс теряет связь с исходным.
То же самое будет при росте cap слайса - он будет реаллоцирован и связь со старым будет потеряна.

cap увеличивается x2 при малых значениях (до 1024)  и ~1.25x если больше.

При росте до превышения cap длины сабслайса он будет увеличиваться (конец сабслайса будет двигаться по исходному массиву),
и его новые элементы будут совмещены с исходным слайсом.

Перебор слайса for _, el:=range slicen1 - по значению,
Для изменения элементов  использовать по ссылке:  for i, _ :=range slicen1 { slicen1[i]} 



##Подробнее

## 1. Создание сабслайса
При создании **сабслайса** новый срез указывает на тот же массив, что и исходный, но получает свою `len` и `cap`:
- `len(newSlice) = end - start`
- `cap(newSlice) = cap(original) - start`
- **Копирование данных не происходит**.

### Пример:
```go
slicen1 := []int{1, 2, 3, 4, 5}
subSlicen := slicen1[1:3] // {2, 3}
fmt.Println(len(subSlicen)) // 2
fmt.Println(cap(subSlicen)) // 4 (5 - 1)
```

## 2. Реаллокация памяти
- Если размер слайса **изменяется** и это приводит к **реаллокации памяти** (например, `append()` превышает `cap`), дочерний слайс **теряет связь** с исходным.
- Если `cap` слайса увеличивается, создаётся **новый массив**, а старый остаётся без изменений.

## 3. Рост `cap`
**Правило роста `cap` в Go**:
- Если `cap < 1024`, он увеличивается в **2 раза**.
- Если `cap ≥ 1024`, он увеличивается примерно на **1.25x**.

### Пример:
```go
s := make([]int, 2, 2)
fmt.Println(cap(s)) // 2

s = append(s, 1) // Превышает cap=2 → cap увеличится до 4
fmt.Println(cap(s)) // 4

s = append(s, 2, 3, 4) // Превышает cap=4 → cap увеличится до 8
fmt.Println(cap(s)) // 8
```

## 4. Увеличение `len` в пределах `cap`
- Если `append()` увеличивает `len` в пределах `cap`, новые элементы **записываются в исходный массив**.

## 5. Перебор слайса
- `for _, el := range slicen1` – передаёт **по значению**.
- Чтобы изменять элементы слайса, **используйте индексы**:
  ```go
  for i := range slicen1 {
      slicen1[i] = 99 // Изменит сам массив
  }
  ```
- **Ошибка**: `for i, _ := range slicen1 { slicen1[i] }` – тут пропущено присваивание.

---

##  Итог
- Сабслайс **разделяет массив** с оригиналом.  
- При **превышении `cap`** создаётся **новый массив**.  
- `cap` **растёт x2 или ~1.25x**.  
- **Изменения через `range` работают только при индексации**. 🚀  
```

✅ 📌 📄✨
