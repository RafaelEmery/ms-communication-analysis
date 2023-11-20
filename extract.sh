# #!/bin/bash

# Verifica se o número correto de argumentos foi fornecido
if [ "$#" -ne 2 ]; then
  echo "wrong arguments"
  exit 1
fi

# make remove-containers-logs
make get-containers-logs service=$1 repeat=$2
sleep 2
make proccess-log-values service=$1 repeat=$2