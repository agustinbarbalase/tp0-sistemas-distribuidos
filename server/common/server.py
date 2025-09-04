import socket
import logging
import signal

from .protocol import Protocol, UnexpectedMessage, ConnectionClose
from .utils import store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._client_protocol = None
        self._is_closed = False

        def handle_signal(signum, frame):
            self.__shutdown()

        signal.signal(signal.SIGTERM, handle_signal)

    def run(self):
        """
        Runs the server loop, accepting and handling client connections.
        """

        while not self._is_closed:
            client_sock = self.__accept_new_connection()
            if self._is_closed: 
                break
            self._client_protocol = Protocol(client_sock)
            self.__handle_client_connection()

    def __shutdown(self):
        """
        Gracefull shutdown

        Function used for graceful shutdown of the server
        """
        logging.info('action: shutdown | result: in_progress')
        self._is_closed = True
        self._client_protocol.close()
        self._server_socket.close()
        logging.info('action: shutdown | result: success')

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet = self._client_protocol.recv_bet()
            logging.info(f"action: apuesta_recibida | result: success | dni: {bet.document} | numero: {bet.number}")
            store_bets([bet])
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
            self._client_protocol.send_success_msg()
        except UnexpectedMessage as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            self._client_protocol.send_failure_msg("Invalid message sent")
        except ConnectionClose as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            self._client_protocol.send_failure_msg("Internal server error")
        finally:
            self._client_protocol.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        try:
            # Connection arrived
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError as err:
            if not self._is_closed: 
                raise err
            return None
