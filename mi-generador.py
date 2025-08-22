

def define_server(file):
  file.write(
    "server:\n"
    "  container_name: server\n"
    "  image: server:latest\n"
    "  entrypoint: python3 /main.py\n"
    "  environment:\n"
    "    - PYTHONUNBUFFERED=1\n"
    "    - LOGGING_LEVEL=DEBUG\n"
    "  networks:\n"
    "    - testing_net\n"
    "\n"
  )

def define_client(file, number_of_client):
  file.write(
    f"  client{number_of_client}:\n"
    f"   container_name: client{number_of_client}\n"
    "    image: client:latest\n"
    "    entrypoint: /client\n"
    "    environment:\n"
    f"     - CLI_ID={number_of_client}\n"
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
    "\n"
  )

if __name__ == "__main__":
  pass