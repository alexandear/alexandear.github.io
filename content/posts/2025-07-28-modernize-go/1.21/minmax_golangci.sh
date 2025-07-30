#! /bin/sh

# nofix

# << snippet begin >>
cat > .golangci.yml << 'EOF'
version: "2"
linters:
  settings:
    revive:
      enable-all-rules: false
      rules:
        - name: redefines-builtin-id
run:
  issues-exit-code: 0
EOF

# No auto fix: redefines-builtin-id: redefinition of the built-in function min
golangci-lint run --config .golangci.yml --enable-only revive ./...
# << snippet end >>
