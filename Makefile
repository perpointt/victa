.PHONY: migrate

# Цель для накатывания миграций
migrate:
	chmod +x scripts/migrate.sh
	./scripts/migrate.sh
