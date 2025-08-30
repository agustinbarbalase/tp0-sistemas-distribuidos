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
    BATCH_HEADER: bytes = b"\x04"  # Header indicating a batch of bets
    FINISH_HEADER: bytes = b"\x05"  # Header indicating the end of transmission
    ID_HEADER: bytes  = b"\x06"  # Header indicating an identification message

    # --- Constants for formatting ---
    SEPARATOR: str          = ";"  # Message separator
    NUM_OF_ATTRIBUTES: int  =  6   # Number of attributes in a bet message

    def __init__(self, socket):
        self._socket = socket

    def wait_identification(self) -> int:
        """
        Waits for an identification message from the client.
        """
        header: bytes = self.__recv_all(Protocol.SIZE_HEADER_BYTES)
        if header != Protocol.ID_HEADER:
            raise UnexpectedMessage("Invalid header")

        id: int = self.__ntohs(self.__recv_all(Protocol.SIZE_LENGTH_BYTES))
        return id
    
    def wait_for_finalization(self) -> None:
        header: bytes = self.__recv_all(Protocol.SIZE_HEADER_BYTES)
        if header != Protocol.FINISH_HEADER:
            raise UnexpectedMessage("Invalid header")

    def send_winner(self, bet: Bet) -> None:
        """
        Sends a winning bet message to the client.
        """
        self._socket.sendall(Protocol.BET_HEADER)

        bet_serialized = self.__serialize_bet(bet)

        self._socket.sendall(self.__htons(len(bet_serialized)))
        self._socket.sendall(bet_serialized)

    def finish_lottery(self) -> None:
        self._socket.sendall(Protocol.FINISH_HEADER)

    def recv_batch_bets(self) -> tuple[list[Bet], int]:
        """
        Receives a batch of bets from the connection.
        This method first reads and validates the batch header from the incoming data.
        It then reads the number of bets to expect, and subsequently receives each bet.
        """
        header: bytes = self.__recv_all(Protocol.SIZE_HEADER_BYTES)
        if header == Protocol.FINISH_HEADER:
            return [], 0
        elif header != Protocol.BATCH_HEADER:
            raise UnexpectedMessage("Invalid header")

        num_bets: int = self.__ntohs(self.__recv_all(Protocol.SIZE_LENGTH_BYTES))
        bets = []
        errors = 0

        for _ in range(num_bets):
            try:
                bet = self.__recv_bet()
                bets.append(bet)
            except UnexpectedMessage as _:
                errors += 1
        
        return bets, errors

    def recv_one_bet(self) -> Bet:
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

        return self.__recv_bet()

    def send_success_msg(self, number_of_bets: int) -> None:
        """
        Send a success message to the client.
        """
        self._socket.sendall(Protocol.OK_HEADER)
        self._socket.sendall(self.__htons(number_of_bets))

    def send_failure_msg(self, number_of_bets: int) -> None:
        """
        Send a failure message to the client.
        """
        self._socket.sendall(Protocol.FAIL_HEADER)
        self._socket.sendall(self.__htons(number_of_bets))

    def __recv_bet(self) -> Bet:
        """
        Receives a bet message from the socket, deserializes it, and returns a Bet object.

        The function first reads the length of the incoming bet message, then reads the message itself,
        and finally deserializes it into a Bet instance.
        """
        length: int = self.__ntohs(self.__recv_all(Protocol.SIZE_LENGTH_BYTES))
        bet_message: bytes = self.__recv_all(length)

        return self.__deserialize_bet(bet_message)

    def __serialize_bet(self, bet: Bet) -> bytes:
        """
        Serialize a Bet object into a byte sequence.
        """
        bet_data = f"{bet.agency}{Protocol.SEPARATOR}{bet.first_name}{Protocol.SEPARATOR}{bet.last_name}{Protocol.SEPARATOR}{bet.document}{Protocol.SEPARATOR}{bet.birthdate}{Protocol.SEPARATOR}{bet.number}"
        return bet_data.encode("utf-8")

    def __deserialize_bet(self, bytes: bytes) -> Bet:
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
    
    def __htons(self, value: int) -> bytes:
        """
        Convert an integer to a byte sequence using network byte order.

        This is equivalent to the C function `htons`, which converts a
        short integer from host byte order to network byte order.

        The returned byte sequence will be exactly 2 bytes long.
        """
        return value.to_bytes(2, "big", signed=False)

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
