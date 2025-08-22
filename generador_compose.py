import sys

def define_server(file):
  file.write(
    "  server:\n"
    "    container_name: server\n"
    "    image: server:latest\n"
    "    entrypoint: python3 /main.py\n"
    "    environment:\n"
    "      - PYTHONUNBUFFERED=1\n"
    "      - LOGGING_LEVEL=DEBUG\n"
    "    volumes:\n"
    "      - ./config.ini:/config.ini\n"
    "    networks:\n"
    "      - testing_net\n"
    "\n"
  )

def define_client(file, number_of_client):
  file.write(
    f"  client{number_of_client}:\n"
    f"    container_name: client{number_of_client}\n"
    "    image: client:latest\n"
    "    entrypoint: /client\n"
    "    environment:\n"
    f"      - CLI_ID={number_of_client}\n"
    "      - CLI_LOG_LEVEL=DEBUG\n"
    "    networks:\n"
    "      - testing_net\n"
    "    depends_on:\n"
    "      - server\n"
    "\n"
  )

def define_network(file):
  file.write(
    "networks:\n"
    "  testing_net:\n"
    "    ipam:\n"
    "      driver: default\n"
    "      config:\n"
    "        - subnet: 172.25.125.0/24\n"
  )

def create_docker_compose(filename, number_of_clients):
  with open(filename, "w") as file:
    file.write(
      "name: tp0\n"
      "services:\n"
    )
    
    define_server(file)
    
    for i in range(1, number_of_clients + 1):
      define_client(file, i)
    
    define_network(file)

def main():
  try:
    if len(sys.argv) != 3:
      print("Usage: python3 mi-generador.py <filename> <number_of_clients>")
      sys.exit(1)

    filename = sys.argv[1]
    number_of_clients = int(sys.argv[2])

    if number_of_clients < 0:
      print("number_of_clients must be not negative")
      sys.exit(1)

    create_docker_compose(filename, number_of_clients)

  except ValueError:
    print("Invalid number of clients. The parameter must be an integer.")
    sys.exit(1)

if __name__ == "__main__":
  main()
