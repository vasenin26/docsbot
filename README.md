# docsbot — сервер метрик

Это минимальный HTTP-сервис для проекта docsbot, который предоставляет только метрики Prometheus и health-проверки (liveness/readiness). Сервис не содержит бизнес-логики — только runtime-метрики и простые эндпойнты для мониторинга.

Функции

- GET /healthz — возвращает 200 OK и JSON {"status":"ok"}
- GET /readyz — возвращает 200 OK (готовность)
- GET /metrics — Prometheus endpoint с метриками

Запуск локально

Сборка бинаря и запуск:

    make build
    make run

По умолчанию сервер слушает порт 9090. Можно переопределить порт через переменную окружения METRICS_PORT:

    METRICS_PORT=8080 ./bin/docsbot

Тестирование

Запуск unit-тестов:

    make test

Docker

Сборка образа для архитектуры хоста:

    make docker-build

Запуск контейнера локально:

    make docker-run

Сборка multi-arch образа (требует docker buildx и настроенного реестра):

    make docker-buildx PLATFORMS=linux/amd64,linux/arm64 TAG=yourrepo/docsbot:tag

Подготовка buildx (если ещё не настроен):

    docker buildx create --use --name mybuilder
    docker buildx inspect --bootstrap

Dockerfile и сборка

Dockerfile использует multi-stage сборку: сначала образ golang (builder) компилирует статический бинарь (CGO_ENABLED=0 по умолчанию), затем финальный образ — distroless static:nonroot. При необходимости для отладки можно заменить финальный этап на alpine.

CI

Добавлен пример GitHub Actions workflow (.github/workflows/ci.yml), который выполняет тесты и собирает multi-arch образ с помощью buildx (push в GitHub Container Registry). Настройте секреты и права доступа к реестру при необходимости.

Примечания

- Метрики реализованы с помощью github.com/prometheus/client_golang.
- Если в будущем потребуется собирать дополнительные метрики или автоматически инкрементировать счётчики для каждого HTTP-запроса, можно добавить middleware, который будет обновлять internal/metrics.RequestsTotal.
- Для внешних HTTPS-запросов убедитесь, что в образ добавлены CA certificates, если это требуется в runtime.

Контакты

Для вопросов по реализации и сборке — откройте issue в репозитории.
