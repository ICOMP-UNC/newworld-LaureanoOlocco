#!/bin/bash

# Ejecuta las pruebas con cobertura
go test -coverprofile=coverage.out ./...

# Extrae el porcentaje de cobertura
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "Coverage is $coverage%"

# Verifica si la cobertura es mayor al 20%
if (( $(echo "$coverage > 20" | bc -l) )); then
  echo "Coverage is greater than 20%"
  exit 0
else
  echo "Coverage is not greater than 20%"
  exit 1
fi