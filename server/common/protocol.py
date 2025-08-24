from utils import Bet

class ConnectionClose(Exception):
  pass

class Protocol:
  SIZE_HEADER = 1
  SIZE_FIELDS = 2
  SIZE_DOCUMENT = 8
  SIZE_DATE = 10
  SIZE_NUMBER = 4

  BET_HEADER = 1
  OK_HEADER = 2
  FAIL_HEADER = 3

  def __init__(self, socket):
    self._socket = socket

  def recv_bet(self) -> Bet:
    code = self.__recv_all(Protocol.SIZE_HEADER)
    if code != Protocol.BET_HEADER.to_bytes(Protocol.SIZE_HEADER):
      raise ValueError("Invalid header")

    length_first_name = int.from_bytes(self.__recv_all(Protocol.SIZE_FIELDS), "big")
    first_name = self.__recv_all(length_first_name).decode('utf-8')
    
    length_last_name = int.from_bytes(self.__recv_all(Protocol.SIZE_FIELDS), "big")
    last_name = self.__recv_all(length_last_name).decode('utf-8')

    document = self.__recv_all(Protocol.SIZE_DOCUMENT).decode('utf-8')
    birthdate = self.__recv_all(Protocol.SIZE_DATE).decode('utf-8')
    number = self.__recv_all(Protocol.SIZE_NUMBER).decode('utf-8')

    return Bet('1', first_name, last_name, document, birthdate, number)

  def send_success_msg(self) -> None:
    self._socket.sendall(Protocol.OK_HEADER.to_bytes(Protocol.SIZE_HEADER))

  def send_failure_msg(self) -> None:
    self._socket.sendall(Protocol.FAIL_HEADER.to_bytes(Protocol.SIZE_HEADER))

  def __recv_all(self, length: int) -> bytes:
    data = b''

    while len(data) < length:
      chunk = self._socket.recv(length - len(data))
      if not chunk:
        raise ConnectionClose("Socket closed")
      data += chunk

    return data
