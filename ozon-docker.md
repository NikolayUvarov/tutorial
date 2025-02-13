
# Docker@Ozon

## Подключение сервиса для тестирования приложения

*Как у меня стало все хорошо:*

- создаем аксесс токен (меню Access Tokens на https://gitlab.ozon.dev/-/profile) с разрешениями на все элементы (там галочки)
- на приглашение команды docker login ... ввел свое имя пользователя и этот токен
- последовательность команд выполнилась успешно:

(выполняется в консоли wsl2/ubuntu)

    docker login gitlab-registry.ozon.dev 
    docker run --rm -it -p 8082:8082 -d gitlab-registry.ozon.dev/go/classroom-16/students/homework-draft/products




