from .utils import Bet

class ConnectionClose(Exception):
    """
    Custom exception raised when a socket connection is unexpectedly closed.
    """
    pass

class Protocol:
    """
    Protocol implementation for sending and receiving structured messages
    over a socket connection.

    The protocol defines message headers and fixed-size fields that are used
    to parse and construct messages exchanged between client and server.
    """

    # --- Field sizes ---
    SIZE_HEADER_BYTES = 1  # Size (in bytes) of the message header
    SIZE_FIELDS_BYTES = 2  # Size (in bytes) for variable-length field lengths
    SIZE_DOCUMENT_BYTES = 8  # Size (in bytes) of the document number (DNI)
    SIZE_DATE_BYTES = 10  # Size (in bytes) of the birthdate field
    SIZE_NUMBER_BYTES = 4  # Size (in bytes) of the numeric field (e.g., bet number)

    # --- Headers ---
    BET_HEADER = 1  # Header indicating a bet message
    OK_HEADER = 2  # Header indicating a success response
    FAIL_HEADER = 3  # Header indicating a failure response

    def __init__(self, socket):
        self._socket = socket

    def recv_bet(self) -> Bet:
        """
        Receive a bet message from the client.

        The function expects the message to start with a `BET_HEADER`,
        followed by the fields encoded in the following order:
        - First name (length + string)
        - Last name (length + string)
        - Document number (fixed length)
        - Birthdate (fixed length)
        - Number (fixed length)
        """
        code = self.__recv_all(Protocol.SIZE_HEADER_BYTES)
        if code != Protocol.BET_HEADER.to_bytes(Protocol.SIZE_HEADER_BYTES, "big"):
            raise ValueError("Invalid header")

        first_name_length = int.from_bytes(
            self.__recv_all(Protocol.SIZE_FIELDS_BYTES), "big"
        )
        first_name = self.__recv_all(first_name_length).decode("utf-8")

        last_name_length = int.from_bytes(
            self.__recv_all(Protocol.SIZE_FIELDS_BYTES), "big"
        )
        last_name = self.__recv_all(last_name_length).decode("utf-8")

        document = self.__recv_all(Protocol.SIZE_DOCUMENT_BYTES).decode("utf-8")
        birthdate = self.__recv_all(Protocol.SIZE_DATE_BYTES).decode("utf-8")
        number = self.__recv_all(Protocol.SIZE_NUMBER_BYTES).decode("utf-8")

        return Bet("1", first_name, last_name, document, birthdate, number)

    def send_success_msg(self) -> None:
        """
        Send a success message to the client.
        """
        self._socket.sendall(
            Protocol.OK_HEADER.to_bytes(Protocol.SIZE_HEADER_BYTES, "big")
        )

    def send_failure_msg(self) -> None:
        """
        Send a failure message to the client.
        """
        self._socket.sendall(
            Protocol.FAIL_HEADER.to_bytes(Protocol.SIZE_HEADER_BYTES, "big")
        )

    def __recv_all(self, length: int) -> bytes:
        """
        Receive exactly `length` bytes from the socket.

        Keeps reading until the expected amount of data is received.
        If the socket closes before receiving the full data, raises a
        `ConnectionClose` exception.
        """
        data = b""

        while len(data) < length:
            chunk = self._socket.recv(length - len(data))
            if not chunk:
                raise ConnectionClose("Socket closed")
            data += chunk

        return data
