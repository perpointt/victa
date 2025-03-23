.PHONY: migrate generate

# Цель для накатывания миграций
migrate:
	chmod +x scripts/migrate.sh
	./scripts/migrate.sh

# Цель для генерации кода:
# сначала запускается скрипт generate.sh, затем generate_components.sh
generate:
	chmod +x scripts/generate.sh
	./scripts/generate.sh
	chmod +x scripts/generate_components.sh
	./scripts/generate_components.sh
