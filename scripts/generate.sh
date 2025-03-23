#!/bin/bash
set -e

# Требуемые утилиты: swagger-cli (через npx), oapi-codegen, yq (версия 4+) и perl
# Директории:
SPEC_DIR="openapi/specs"
BUNDLES_DIR="openapi/bundles"
GO_OUTPUT_DIR="internal/api/specs"

# Создаем необходимые папки, если их нет
mkdir -p "$BUNDLES_DIR"
mkdir -p "$GO_OUTPUT_DIR"

# Функция для получения префикса: первая буква имени (группы) становится заглавной
get_prefix() {
    local group="$1"
    local first_letter=$(echo "${group:0:1}" | tr '[:lower:]' '[:upper:]')
    local rest="${group:1}"
    echo "${first_letter}${rest}"
}

# Функция для переименования дублирующихся свойств во всём сгенерированном файле с помощью perl
rename_properties() {
    local group="$1"
    local prefix="$2"
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
    group=$(basename "$SPEC_FILE" .yaml)
    BUNDLE_FILE="$BUNDLES_DIR/openapi-${group}.bundle.yaml"
    
    if [ ! -f "$SPEC_FILE" ]; then
        echo "Файл спецификации не найден: $SPEC_FILE"
        exit 1
    fi
    
    echo "Бандлинг спецификации для группы ${group}..."
    # Флаг --dereference разворачивает ссылки, чтобы уменьшить глубину вложенности
    npx swagger-cli bundle "$SPEC_FILE" --outfile "$BUNDLE_FILE" --type yaml --dereference
    
    echo "Генерация Go кода (типы + сервер) для группы ${group}..."
    oapi-codegen -generate types,gin -o "$GO_OUTPUT_DIR/${group}.gen.go" -package=api "$BUNDLE_FILE"
    
    prefix=$(get_prefix "$group")
    echo "Переименование дублирующихся свойств для группы ${group} с префиксом ${prefix}..."
    rename_properties "$group" "$prefix"
done

# Объединение всех бандлов в один финальный файл openapi.yaml в корне проекта
FINAL_BUNDLE="openapi.yaml"
echo "Объединение всех бандлов в один финальный файл ($FINAL_BUNDLE) для Swagger..."
yq eval-all 'select(fileIndex == 0) *+ select(fileIndex > 0) | .servers |= unique' "$BUNDLES_DIR"/*.bundle.yaml > "$FINAL_BUNDLE"

echo "Удаление папки $BUNDLES_DIR..."
rm -rf "$BUNDLES_DIR"

echo "Генерация завершена успешно!"
echo "Финальный Swagger файл: $FINAL_BUNDLE"
