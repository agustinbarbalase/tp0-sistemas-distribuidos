

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

if __name__ == "__main__":
  pass