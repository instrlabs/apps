#!/bin/bash
echo "🔧 Switching to development mode..."
echo "Adding local shared library replace directives..."

find . -name "go.mod" -not -path "./shared/*" | while read file; do
    if ! grep -q "replace github.com/instrlabs/shared" "$file"; then
        echo "replace github.com/instrlabs/shared => ../../instrlabs-shared" >> "$file"
        echo "✓ Added local shared to $(dirname "$file")"
    else
        echo "- Already in dev mode: $(dirname "$file")"
    fi
done

echo ""
echo "🚀 Development mode enabled!"
echo "💡 Use 'make prod-mode' to return to production mode"