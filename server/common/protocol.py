from .utils import Bet

class UnexpectedMessage(Exception):
    """
    Custom exception raised when an unexpected message is received.
    """
    pass

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
    SIZE_HEADER_BYTES   =  1   # Size of the message header
    SIZE_LENGTH_BYTES   =  2   # Size for message length

    # --- Headers ---
    BET_HEADER: bytes  = b"\x01"  # Header indicating a bet message
    OK_HEADER: bytes   = b"\x02"  # Header indicating a success response
    FAIL_HEADER: bytes = b"\x03"  # Header indicating a failure response

    # --- Constants for formatting ---
    SEPARATOR: str          = ";"  # Message separator
    NUM_OF_ATTRIBUTES: int  =  6   # Number of attributes in a bet message

    def __init__(self, socket):
        self._socket = socket

    def recv_bet(self) -> Bet:
        """
        Receive a bet message from the client.

        The function expects the message to start with a `BET_HEADER`,
        followed by the fields encoded in separated by the following order:
        - Agency number
        - First name
        - Last name
        - Document number
        - Birthdate
        - Number
        """
        header: bytes = self.__recv_all(Protocol.SIZE_HEADER_BYTES)
        if header != Protocol.BET_HEADER: 
            raise UnexpectedMessage("Invalid header")

        length: int = self.__ntohs(self.__recv_all(Protocol.SIZE_LENGTH_BYTES))
        bet_message: bytes = self.__recv_all(length)

        return self.__deserialize(bet_message)

    def send_success_msg(self) -> None:
        """
        Send a success message to the client.
        """
        self._socket.sendall(Protocol.OK_HEADER)

    def send_failure_msg(self) -> None:
        """
        Send a failure message to the client.
        """
        self._socket.sendall(Protocol.FAIL_HEADER)

    def __deserialize(self, bytes: bytes) -> Bet:
        """
        Deserialize a bet message from the socket.
        """

        attributes = bytes.decode("utf-8").split(Protocol.SEPARATOR)
        if len(attributes) != Protocol.NUM_OF_ATTRIBUTES:
            raise UnexpectedMessage("Invalid bet message format")

        return Bet(*attributes)

    def __ntohs(self, bytes: bytes) -> int:
        """
        Convert a byte sequence to an integer using network byte order.

        This is equivalent to the C function `ntohs`, which converts a
        short integer from network byte order to host byte order.

        If `bytes` is not at least 2 bytes long, the behavior is undefined.
        """
        return int.from_bytes(bytes[:2], "big", signed=False)

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
