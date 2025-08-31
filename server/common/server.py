import socket
import logging
import signal
import multiprocessing

from .protocol import Protocol, UnexpectedMessage, ConnectionClose
from .utils import store_bets, load_bets, has_won


class Server:
    def __init__(self, port, listen_backlog, amount_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)
        self._amount_of_clients = amount_of_clients
        self._client_sockets = {}
        self._client = []
        self._barrier_for_winners = multiprocessing.Barrier(amount_of_clients + 1)
        self._barrier_for_finish_connection = multiprocessing.Barrier(
            amount_of_clients + 1
        )
        self._store_lock = multiprocessing.Lock()
        self._load_lock = multiprocessing.Lock()
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
            if self._is_closed:
                break
            thread = multiprocessing.Process(
                target=self.__handle_client_connection,
                args=(
                    client_sock,
                    self._barrier_for_winners,
                    self._barrier_for_finish_connection,
                    self._store_lock,
                    self._load_lock,
                ),
            )
            thread.start()
            self._client.append(thread)
            if len(self._client) == self._amount_of_clients:
                self._barrier_for_winners.wait()
                for t in self._client:
                    t.join()
                self.__announce_winners()

    def __shutdown(self):
        """
        Gracefull shutdown

        Function used for graceful shutdown of the server
        """
        logging.info("action: shutdown | result: in_progress")
        self._is_closed = True
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info("action: shutdown | result: success")

    def __announce_winners(self):
        """
        Announce the winners

        Function used to announce the winners of the game
        """

        for bet in load_bets():
            if has_won(bet):
                try:
                    client_sock = self._client_sockets[bet.agency]
                    protocol = Protocol(client_sock)
                    protocol.send_winner(bet)
                except Exception as e:
                    logging.error(
                        f"action: receive_message | result: fail | error: {e}"
                    )

        for client_socket in self._client_sockets.values():
            try:
                protocol = Protocol(client_socket)
                protocol.finish_lottery()
            except Exception as e:
                logging.error(f"action: receive_message | result: fail | error: {e}")
            finally:
                client_socket.close()

    def __handle_client_connection(
        self,
        client_sock,
        barrier_for_winners,
        store_lock,
    ):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = Protocol(client_sock)
            id = protocol.wait_identification()
            self._client_sockets[id] = client_sock

            while True:
                bets, errors = protocol.recv_batch_bets()
                with store_lock:
                    store_bets(bets)

                if errors > 0:
                    logging.error(
                        f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}"
                    )
                    protocol.send_failure_msg(len(bets))
                elif len(bets) > 0:
                    logging.info(
                        f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
                    )
                    protocol.send_success_msg(len(bets))
                elif len(bets) == 0:
                    logging.info("action: esperando_ganador | result: success")
                    break

            barrier_for_winners.wait()

        except UnexpectedMessage as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            protocol.send_failure_msg("Invalid message sent")
        except ConnectionClose as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            protocol.send_failure_msg("Internal server error")

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        try:
            # Connection arrived
            logging.info("action: accept_connections | result: in_progress")
            c, addr = self._server_socket.accept()
            logging.info(
                f"action: accept_connections | result: success | ip: {addr[0]}"
            )
            return c
        except OSError as err:
            if not self._is_closed:
                logging.error(
                    f"action: accept_connections | result: fail | error: {err}"
                )
                raise err
            return None
