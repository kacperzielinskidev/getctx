#!/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
OUTPUT_FILE="$SCRIPT_DIR/get_context.txt"

# Uruchom interaktywny selektor i przechwyć jego wynik do tablicy
echo "Uruchamianie interaktywnego selektora plików..."
# Używamy mapfile (lub read -a w pętli dla starszych wersji basha)
mapfile -t SELECTED_FILES < <(./getctx)

# Sprawdź, czy cokolwiek zostało wybrane
if [ ${#SELECTED_FILES[@]} -eq 0 ]; then
  echo "❌ Nie wybrano żadnych plików. Przerywam."
  exit 1
fi

echo "Wybrano ${#SELECTED_FILES[@]} plików. Budowanie kontekstu..."

# Wyczyść plik wyjściowy
> "$OUTPUT_FILE"

echo "🚀 Starting to build context file: $OUTPUT_FILE"

for file in "${SELECTED_FILES[@]}"; do
  if [[ -f "$file" ]]; then
    echo "   -> Adding content from: $file"
    
    echo "--- START OF FILE: $file ---" >> "$OUTPUT_FILE"
    cat "$file" >> "$OUTPUT_FILE"
    echo -e "\n--- END OF FILE: $file ---\n" >> "$OUTPUT_FILE"
  fi
done

echo "✅ Done! All content has been combined into $OUTPUT_FILE"