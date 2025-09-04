import socket
import logging
import signal

from multiprocessing import Barrier, Process, Lock
from threading import BrokenBarrierError
from .protocol import Protocol, UnexpectedMessage, ConnectionClose
from .utils import store_bets, load_bets, has_won


class Server:
    def __init__(self, port, listen_backlog, amount_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)
        self._clients = []
        self._barrier_for_winners = Barrier(amount_of_clients)
        self._store_lock = Lock()
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
            if self._is_closed or client_sock is None:
                break
            protocol = Protocol(client_sock)
            process = Process(
                target=self.__handle_client_connection,
                args=(
                    protocol,
                    self._barrier_for_winners,
                    self._store_lock,
                ),
            )
            process.start()
            self._clients.append(process)

        for clients in self._clients:
            clients.terminate()
            clients.join()

    def __shutdown(self):
        """
        Gracefull shutdown

        Function used for graceful shutdown of the server
        """
        logging.info("action: shutdown | result: in_progress")
        self._is_closed = True
        try:
            self._server_socket.close()
        except Exception as e:
            logging.error(f"action: shutdown | result: fail | error: {e}")
        logging.info("action: shutdown | result: success")

    def __send_winners_to_agency(self, agency_id, client_protocol):
        """
        Sends winning bets to their respective agencies.
        """
        winners = []
        for bet in load_bets():
            if has_won(bet) and bet.agency == agency_id:
                winners.append(bet.document)
        try:
            client_protocol.send_winners(winners)
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")

    def __handle_client_connection(
        self,
        protocol,
        barrier_for_winners,
        store_lock,
    ):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        def handle_signal(signum, frame):
            protocol.close()

        signal.signal(signal.SIGTERM, handle_signal)

        try:
            id = protocol.wait_identification()

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
            self.__send_winners_to_agency(id, protocol)
        except BrokenBarrierError as _:
            return
        except UnexpectedMessage as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except ConnectionClose as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            logging.info("action: close_connection | result: in_progress")
            protocol.close()
            logging.info("action: close_connection | result: success")

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
