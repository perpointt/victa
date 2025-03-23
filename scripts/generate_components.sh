#!/bin/bash
set -e

# Директории
COMPONENTS_DIR="openapi/components"
BUNDLES_DIR="openapi/bundles"
OUTPUT_DIR="internal/api/components"

mkdir -p "$OUTPUT_DIR"
mkdir -p "$BUNDLES_DIR"

# Функция для конвертации snake_case в CamelCase
to_camel_case() {
    local input="$1"
    echo "$input" | awk -F"_" '{
        for(i=1;i<=NF;i++){
            $i = toupper(substr($i,1,1)) substr($i,2)
        }
        printf("%s", $1);
        for(i=2;i<=NF;i++){
            printf("%s", $i)
        }
        print ""
    }'
}

for SPEC_FILE in "$COMPONENTS_DIR"/*.yaml; do
    base=$(basename "$SPEC_FILE" .yaml)
    COMPONENT_NAME=$(to_camel_case "$base")
    TEMP_SPEC="$BUNDLES_DIR/temp_${base}.yaml"
    OUTPUT_FILE="$OUTPUT_DIR/${base}.gen.go"
    
    echo "Обрабатывается $SPEC_FILE → Component: $COMPONENT_NAME"
    
    # Формируем временную спецификацию, включая фиктивный эндпоинт, ссылающийся на компонент
    cat > "$TEMP_SPEC" <<EOF
openapi: 3.0.0
info:
  title: Dummy API for $COMPONENT_NAME
  version: "1.0.0"
paths:
  /dummy:
    get:
      summary: Dummy endpoint
      operationId: dummy
      responses:
        "200":
          description: Dummy response
          content:
            application/json:
              schema:
                \$ref: "#/components/schemas/$COMPONENT_NAME"
components:
  schemas:
    $COMPONENT_NAME:
EOF

    # Добавляем содержимое исходного файла с отступом (6 пробелов)
    sed 's/^/      /' "$SPEC_FILE" >> "$TEMP_SPEC"
    
    echo "Временный файл $TEMP_SPEC создан"
    
    # Генерируем Go-структуру с помощью oapi-codegen
    oapi-codegen -generate types -o "$OUTPUT_FILE" -package=api "$TEMP_SPEC"
    
    # При необходимости, заменяем map[string]interface{} на interface{}
    sed -i '' 's/map\[string\]interface{}/interface{}/g' "$OUTPUT_FILE"
    
    rm "$TEMP_SPEC"
    echo "Сгенерированный файл: $OUTPUT_FILE"
    echo "------------------------------------"
done

echo "Все компоненты обработаны."
echo "Удаляем папку $BUNDLES_DIR..."
rm -rf "$BUNDLES_DIR"
echo "Папка $BUNDLES_DIR удалена."
