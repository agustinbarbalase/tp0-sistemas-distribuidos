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
        self._is_closed = False

        def handle_signal(signum, frame):
            self.__shutdown()

        signal.signal(signal.SIGTERM, handle_signal)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not self._is_closed:
            client_sock = self.__accept_new_connection()
            if self._is_closed: break
            self.__handle_client_connection(client_sock)

    def __shutdown(self):
        """
        Gracefull shutdown

        Function used for graceful shutdown of the server
        """
        logging.info('action: shutdown | result: in_progress')
        self._is_closed = True
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info('action: shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = Protocol(client_sock)
            bets, errors = protocol.recv_batch_bets()
            store_bets(bets)
            if errors > 0:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}")
                protocol.send_failure_msg(f"{errors} bets were invalid")
            else:
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                protocol.send_success_msg()
        except UnexpectedMessage as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            protocol.send_failure_msg("Invalid message sent")
        except ConnectionClose as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            protocol.send_failure_msg("Internal server error")
        finally:
            client_sock.close()

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
            if not self._is_closed: raise err
            return None
