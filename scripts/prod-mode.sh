#!/bin/bash
echo "🏭 Switching to production mode..."
echo "Removing local instrlabs-shared library replace directives..."

find . -name "go.mod" -not -path "../../instrlabs-shared/*" | while read file; do
    if grep -q "replace github.com/instrlabs/shared" "$file"; then
        # Create a temporary file without the replace line
        grep -v "replace github.com/instrlabs/shared" "$file" > "$file.tmp"
        mv "$file.tmp" "$file"
        echo "✓ Removed local instrlabs-shared from $(dirname "$file")"
    else
        echo "- Already in prod mode: $(dirname "$file")"
    fi
done

echo ""
echo "📦 Production mode enabled!"
echo "💡 Use 'make dev-mode' to return to development mode"