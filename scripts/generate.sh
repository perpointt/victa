#!/bin/bash
set -e

# Директории
SPEC_DIR="openapi/specs"
BUNDLES_DIR="openapi/bundles"
GO_OUTPUT_DIR="internal/api/specs"
FINAL_FILE="openapi.yaml"

# Создаем необходимые папки, если их нет
mkdir -p "$BUNDLES_DIR"
mkdir -p "$GO_OUTPUT_DIR"

# Функция для конвертации строки с дефисами (snake-case или kebab-case) в CamelCase
to_camel_case() {
    local input="$1"
    # Заменяем дефисы на подчеркивания
    input="${input//-/_}"
    echo "$input" | awk -F"_" '{
        for(i=1; i<=NF; i++){
            $i = toupper(substr($i,1,1)) substr($i,2)
        }
        printf("%s", $1);
        for(i=2; i<=NF; i++){
            printf("%s", $i)
        }
        print ""
    }'
}

# Функция для конвертации строки в snake_case (замена дефисов на подчеркивания и приведение к нижнему регистру)
to_snake_case() {
    local input="$1"
    input="${input//-/_}"
    echo "$input" | tr '[:upper:]' '[:lower:]'
}

# Функция для переименования общих типов в сгенерированном Go файле
rename_properties() {
    local group="$1"     # группа в snake_case (без дефисов)
    local prefix="$2"    # CamelCase префикс (например, CompanyUsers)
    local file="$GO_OUTPUT_DIR/${group}.gen.go"
    perl -pi -e "s/\bServerInterfaceWrapper\b/${prefix}ServerInterfaceWrapper/g" "$file"
    perl -pi -e "s/\bServerInterface\b/${prefix}ServerInterface/g" "$file"
    perl -pi -e "s/\bMiddlewareFunc\b/${prefix}MiddlewareFunc/g" "$file"
    perl -pi -e "s/\bGinServerOptions\b/${prefix}GinServerOptions/g" "$file"
    perl -pi -e "s/\bRegisterHandlersWithOptions\b/Register${prefix}HandlersWithOptions/g" "$file"
    perl -pi -e "s/\bRegisterHandlers\b/Register${prefix}Handlers/g" "$file"
}

echo "Обработка всех спецификаций из папки $SPEC_DIR..."

# Итерируем по всем YAML-файлам в SPEC_DIR
for SPEC_FILE in "$SPEC_DIR"/*.yaml; do
    # raw_group может содержать дефисы, например "company-users"
    raw_group=$(basename "$SPEC_FILE" .yaml)
    # Для имени файла используем snake_case (замена дефисов на подчеркивания)
    group=$(to_snake_case "$raw_group")
    # Для префикса типов используем CamelCase (без дефисов)
    PREFIX=$(to_camel_case "$raw_group")
    BUNDLE_FILE="$BUNDLES_DIR/openapi-${group}.bundle.yaml"
    OUTPUT_FILE="$GO_OUTPUT_DIR/${group}.gen.go"
    
    if [ ! -f "$SPEC_FILE" ]; then
        echo "Файл спецификации не найден: $SPEC_FILE"
        exit 1
    fi
    
    echo "Бандлинг спецификации для группы '${raw_group}' → $BUNDLE_FILE"
    npx swagger-cli bundle "$SPEC_FILE" --outfile "$BUNDLE_FILE" --type yaml --dereference
    
    echo "Генерация Go кода (типы + сервер) для группы '${raw_group}'..."
    oapi-codegen -generate types,gin -o "$OUTPUT_FILE" -package=api "$BUNDLE_FILE"
    
    echo "Переименование общих типов для группы '${raw_group}' с префиксом ${PREFIX}..."
    rename_properties "$group" "$PREFIX"
done

echo "Объединение всех бандлов в финальный Swagger-файл..."

# Создаем мастер-спецификацию с базовыми данными
MASTER_FILE="$BUNDLES_DIR/master.yaml"
cat > "$MASTER_FILE" <<EOF
openapi: 3.0.0
info:
  title: VICTA API
  version: "1.0.0"
servers:
  - url: /api/v1
paths: {}
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
EOF

# Объединяем раздел paths из каждого бандла с мастер-спецификацией
for BUNDLE in "$BUNDLES_DIR"/*.bundle.yaml; do
    echo "Объединяем пути из $BUNDLE"
    yq eval '.paths' "$BUNDLE" > temp_paths.yaml
    yq eval-all 'select(fileIndex == 0) *+ select(fileIndex == 1)' "$MASTER_FILE" temp_paths.yaml > merged.yaml
    mv merged.yaml "$MASTER_FILE"
    rm temp_paths.yaml
done

# Убираем дублирование серверов, если оно есть
yq eval '.servers |= unique' "$MASTER_FILE" > merged.yaml && mv merged.yaml "$MASTER_FILE"

# Перемещаем мастер-спецификацию в FINAL_FILE (корень проекта)
mv "$MASTER_FILE" "$FINAL_FILE"

echo "Удаление папки $BUNDLES_DIR..."
rm -rf "$BUNDLES_DIR"

echo "Генерация завершена успешно!"
echo "Финальный Swagger файл: $FINAL_FILE"
